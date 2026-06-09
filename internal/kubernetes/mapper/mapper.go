package mapper

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/tokanize/kubernetes-gateway-exporter/pkg/models"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
)

const gatewayAPIGroup = "gateway.networking.k8s.io"

// Mapper maps relationships between Gateways, HTTPRoutes, and Services.
// See docs/adrs/001-k8s-informers.md for the caching strategy.
type Mapper struct {
	client client.Client
	logger *slog.Logger

	mu            sync.RWMutex
	exposedRoutes []models.ExposedRoute
}

// NewMapper creates a new relationship mapper using the provided cached client.
func NewMapper(c client.Client, logger *slog.Logger) *Mapper {
	return &Mapper{
		client: c,
		logger: logger.With(slog.String("component", "mapper")),
	}
}

// Start satisfies the manager.Runnable interface to run the background periodic cache updater.
func (m *Mapper) Start(ctx context.Context) error {
	m.logger.Info("Starting Mapper background cache updater")

	// Perform initial synchronous update so metrics are populated immediately
	if err := m.UpdateCache(ctx); err != nil {
		m.logger.Error("Initial cache update failed", slog.Any("error", err))
	}

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			m.logger.Info("Stopping Mapper background cache updater")
			return nil
		case <-ticker.C:
			if err := m.UpdateCache(ctx); err != nil {
				m.logger.Error("Failed to update cache in background", slog.Any("error", err))
			}
		}
	}
}

// UpdateCache calculates all exposed routes and updates the thread-safe in-memory cache.
func (m *Mapper) UpdateCache(ctx context.Context) error {
	routes, err := m.calculateExposedRoutes(ctx)
	if err != nil {
		return err
	}

	m.mu.Lock()
	m.exposedRoutes = routes
	m.mu.Unlock()

	m.logger.Debug("In-memory routing cache updated", slog.Int("routes_count", len(routes)))
	return nil
}

func parentGroup(ref gwv1.ParentReference) string {
	if ref.Group == nil {
		return gatewayAPIGroup
	}
	return string(*ref.Group)
}

func parentKind(ref gwv1.ParentReference) string {
	if ref.Kind == nil {
		return "Gateway"
	}
	return string(*ref.Kind)
}

func parentNamespace(routeNamespace string, ref gwv1.ParentReference) string {
	if ref.Namespace == nil {
		return routeNamespace
	}
	return string(*ref.Namespace)
}

func optionalSectionEqual(left, right *gwv1.SectionName) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}
	return *left == *right
}

func optionalPortEqual(left, right *gwv1.PortNumber) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}
	return *left == *right
}

func sameParentReference(routeNamespace string, expected, actual gwv1.ParentReference) bool {
	return expected.Name == actual.Name &&
		parentNamespace(routeNamespace, expected) == parentNamespace(routeNamespace, actual) &&
		parentGroup(expected) == parentGroup(actual) &&
		parentKind(expected) == parentKind(actual) &&
		optionalSectionEqual(expected.SectionName, actual.SectionName) &&
		optionalPortEqual(expected.Port, actual.Port)
}

func isRouteAcceptedByParent(route *gwv1.HTTPRoute, parentRef gwv1.ParentReference) bool {
	for _, parentStatus := range route.Status.Parents {
		if !sameParentReference(route.Namespace, parentRef, parentStatus.ParentRef) {
			continue
		}

		for _, condition := range parentStatus.Conditions {
			if condition.Type == string(gwv1.RouteConditionAccepted) &&
				condition.Status == metav1.ConditionTrue &&
				condition.ObservedGeneration == route.Generation {
				return true
			}
		}
	}
	return false
}

func wildcardSuffix(hostname string) (string, bool) {
	if !strings.HasPrefix(hostname, "*.") {
		return "", false
	}
	return hostname[1:], true
}

func wildcardMatches(wildcard, hostname string) bool {
	suffix, ok := wildcardSuffix(wildcard)
	if !ok || !strings.HasSuffix(hostname, suffix) {
		return false
	}
	return len(hostname) > len(suffix)
}

func intersectHostname(routeHostname, listenerHostname string) (string, bool) {
	if routeHostname == listenerHostname {
		return routeHostname, true
	}
	if wildcardMatches(listenerHostname, routeHostname) {
		return routeHostname, true
	}
	if wildcardMatches(routeHostname, listenerHostname) {
		return listenerHostname, true
	}

	routeSuffix, routeWildcard := wildcardSuffix(routeHostname)
	listenerSuffix, listenerWildcard := wildcardSuffix(listenerHostname)
	if routeWildcard && listenerWildcard {
		switch {
		case strings.HasSuffix(routeSuffix, listenerSuffix):
			return routeHostname, true
		case strings.HasSuffix(listenerSuffix, routeSuffix):
			return listenerHostname, true
		}
	}

	return "", false
}

func intersectHostnames(routeHostnames []gwv1.Hostname, listenerHostname *gwv1.Hostname) []string {
	if len(routeHostnames) == 0 {
		if listenerHostname == nil {
			return []string{""}
		}
		return []string{string(*listenerHostname)}
	}

	intersections := make([]string, 0, len(routeHostnames))
	seen := make(map[string]struct{})
	for _, routeHostname := range routeHostnames {
		result := string(routeHostname)
		if listenerHostname != nil {
			var ok bool
			result, ok = intersectHostname(result, string(*listenerHostname))
			if !ok {
				continue
			}
		}
		if _, exists := seen[result]; exists {
			continue
		}
		seen[result] = struct{}{}
		intersections = append(intersections, result)
	}
	return intersections
}

func listenerAllowsHTTPRouteKind(listener gwv1.Listener) bool {
	if listener.Protocol != gwv1.HTTPProtocolType && listener.Protocol != gwv1.HTTPSProtocolType {
		return false
	}
	if listener.AllowedRoutes == nil || len(listener.AllowedRoutes.Kinds) == 0 {
		return true
	}

	for _, allowedKind := range listener.AllowedRoutes.Kinds {
		group := gatewayAPIGroup
		if allowedKind.Group != nil {
			group = string(*allowedKind.Group)
		}
		if group == gatewayAPIGroup && allowedKind.Kind == "HTTPRoute" {
			return true
		}
	}
	return false
}

func (m *Mapper) listenerAllowsNamespace(ctx context.Context, listener gwv1.Listener, gatewayNamespace, routeNamespace string) bool {
	from := gwv1.NamespacesFromSame
	var selector *metav1.LabelSelector
	if listener.AllowedRoutes != nil && listener.AllowedRoutes.Namespaces != nil {
		if listener.AllowedRoutes.Namespaces.From != nil {
			from = *listener.AllowedRoutes.Namespaces.From
		}
		selector = listener.AllowedRoutes.Namespaces.Selector
	}

	switch from {
	case gwv1.NamespacesFromSame:
		return routeNamespace == gatewayNamespace
	case gwv1.NamespacesFromAll:
		return true
	case gwv1.NamespacesFromSelector:
		if selector == nil {
			return false
		}
		var namespace corev1.Namespace
		if err := m.client.Get(ctx, client.ObjectKey{Name: routeNamespace}, &namespace); err != nil {
			return false
		}
		compiledSelector, err := metav1.LabelSelectorAsSelector(selector)
		if err != nil {
			return false
		}
		return compiledSelector.Matches(labels.Set(namespace.Labels))
	default:
		return false
	}
}

func listenerSelected(parentRef gwv1.ParentReference, listener gwv1.Listener) bool {
	if parentRef.SectionName != nil && *parentRef.SectionName != listener.Name {
		return false
	}
	return parentRef.Port == nil || *parentRef.Port == listener.Port
}

func (m *Mapper) listenerAcceptsRoute(ctx context.Context, listener gwv1.Listener, gatewayNamespace string, route *gwv1.HTTPRoute) bool {
	return listenerAllowsHTTPRouteKind(listener) &&
		m.listenerAllowsNamespace(ctx, listener, gatewayNamespace, route.Namespace) &&
		len(intersectHostnames(route.Spec.Hostnames, listener.Hostname)) > 0
}

func (m *Mapper) hasReferenceGrant(ctx context.Context, routeNamespace, backendNamespace, backendName string) bool {
	var grants gwv1.ReferenceGrantList
	if err := m.client.List(ctx, &grants, client.InNamespace(backendNamespace)); err != nil {
		return false
	}

	for _, grant := range grants.Items {
		fromMatches := false
		for _, from := range grant.Spec.From {
			if string(from.Group) == gatewayAPIGroup && from.Kind == "HTTPRoute" && string(from.Namespace) == routeNamespace {
				fromMatches = true
				break
			}
		}
		if !fromMatches {
			continue
		}

		for _, to := range grant.Spec.To {
			if string(to.Group) == "" && to.Kind == "Service" && (to.Name == nil || string(*to.Name) == backendName) {
				return true
			}
		}
	}
	return false
}

func backendIdentity(routeNamespace string, backendRef gwv1.HTTPBackendRef) (namespace, port string, supported bool) {
	group := ""
	if backendRef.Group != nil {
		group = string(*backendRef.Group)
	}
	kind := "Service"
	if backendRef.Kind != nil {
		kind = string(*backendRef.Kind)
	}
	if group != "" || kind != "Service" {
		return "", "", false
	}

	namespace = routeNamespace
	if backendRef.Namespace != nil {
		namespace = string(*backendRef.Namespace)
	}
	if backendRef.Port != nil {
		port = strconv.Itoa(int(*backendRef.Port))
	}
	return namespace, port, true
}

func (m *Mapper) serviceBackendExists(ctx context.Context, namespace string, backendRef gwv1.HTTPBackendRef) bool {
	var service corev1.Service
	if err := m.client.Get(ctx, client.ObjectKey{Namespace: namespace, Name: string(backendRef.Name)}, &service); err != nil {
		return false
	}
	if backendRef.Port == nil {
		return true
	}
	for _, servicePort := range service.Spec.Ports {
		if servicePort.Port == int32(*backendRef.Port) {
			return true
		}
	}
	return false
}

func rulePaths(rule gwv1.HTTPRouteRule) []string {
	if len(rule.Matches) == 0 {
		return []string{"/"}
	}

	paths := make([]string, 0, len(rule.Matches))
	for _, match := range rule.Matches {
		path := "/"
		if match.Path != nil && match.Path.Value != nil {
			path = *match.Path.Value
		}
		paths = append(paths, path)
	}
	return paths
}

func gatewayIPAddress(gateway *gwv1.Gateway) string {
	for _, address := range gateway.Status.Addresses {
		if address.Type != nil && *address.Type == gwv1.IPAddressType {
			return address.Value
		}
	}
	return ""
}

// GetExposedRoutes returns a copy of pre-computed exposed routes in O(1) time.
func (m *Mapper) GetExposedRoutes(ctx context.Context) ([]models.ExposedRoute, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Shallow-copy to prevent external mutation or data races
	routesCopy := make([]models.ExposedRoute, len(m.exposedRoutes))
	copy(routesCopy, m.exposedRoutes)
	return routesCopy, nil
}

// calculateExposedRoutes resolves accepted HTTPRoute relationships before traffic is observed.
func (m *Mapper) calculateExposedRoutes(ctx context.Context) ([]models.ExposedRoute, error) {
	var routes gwv1.HTTPRouteList
	if err := m.client.List(ctx, &routes); err != nil {
		return nil, fmt.Errorf("failed to list HTTPRoutes: %w", err)
	}

	var exposedRoutes []models.ExposedRoute
	for i := range routes.Items {
		route := &routes.Items[i]
		for _, parentRef := range route.Spec.ParentRefs {
			if parentGroup(parentRef) != gatewayAPIGroup || parentKind(parentRef) != "Gateway" || !isRouteAcceptedByParent(route, parentRef) {
				continue
			}

			gatewayNamespace := parentNamespace(route.Namespace, parentRef)
			var gateway gwv1.Gateway
			if err := m.client.Get(ctx, client.ObjectKey{Namespace: gatewayNamespace, Name: string(parentRef.Name)}, &gateway); err != nil {
				continue
			}

			ipAddress := gatewayIPAddress(&gateway)

			lbType := "external"
			gatewayClass := strings.ToLower(string(gateway.Spec.GatewayClassName))
			if strings.Contains(gatewayClass, "internal") || strings.Contains(gatewayClass, "ilb") {
				lbType = "internal"
			}

			for _, listener := range gateway.Spec.Listeners {
				if !listenerSelected(parentRef, listener) || !m.listenerAcceptsRoute(ctx, listener, gatewayNamespace, route) {
					continue
				}

				for _, hostname := range intersectHostnames(route.Spec.Hostnames, listener.Hostname) {
					for _, rule := range route.Spec.Rules {
						for _, path := range rulePaths(rule) {
							for _, backendRef := range rule.BackendRefs {
								backendNamespace, backendPort, supported := backendIdentity(route.Namespace, backendRef)
								if !supported {
									continue
								}
								if backendNamespace != route.Namespace && !m.hasReferenceGrant(ctx, route.Namespace, backendNamespace, string(backendRef.Name)) {
									continue
								}
								if !m.serviceBackendExists(ctx, backendNamespace, backendRef) {
									continue
								}

								exposedRoutes = append(exposedRoutes, models.ExposedRoute{
									Namespace:        route.Namespace,
									GatewayName:      gateway.Name,
									RouteName:        route.Name,
									RouteNamespace:   route.Namespace,
									Hostname:         hostname,
									ListenerName:     string(listener.Name),
									HTTPPath:         path,
									ServiceName:      string(backendRef.Name),
									BackendNamespace: backendNamespace,
									BackendPort:      backendPort,
									LBType:           lbType,
									IPAddress:        ipAddress,
								})
							}
						}
					}
				}
			}
		}
	}

	return deduplicate(exposedRoutes), nil
}

func deduplicate(routes []models.ExposedRoute) []models.ExposedRoute {
	seen := make(map[models.ExposedRoute]struct{}, len(routes))
	deduplicated := make([]models.ExposedRoute, 0, len(routes))
	for _, route := range routes {
		if _, exists := seen[route]; exists {
			continue
		}
		seen[route] = struct{}{}
		deduplicated = append(deduplicated, route)
	}
	return deduplicated
}
