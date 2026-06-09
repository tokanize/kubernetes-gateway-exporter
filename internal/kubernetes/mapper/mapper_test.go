package mapper

import (
	"context"
	"log/slog"
	"strconv"
	"sync"
	"testing"

	"github.com/tokanize/kubernetes-gateway-exporter/pkg/models"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func ptr[T any](v T) *T { return &v }

func TestMapper_GetExposedRoutes(t *testing.T) {
	_ = gwv1.Install(scheme.Scheme)
	_ = corev1.AddToScheme(scheme.Scheme)

	tests := []struct {
		name       string
		namespaces []corev1.Namespace
		services   []corev1.Service
		gateways   []gwv1.Gateway
		routes     []gwv1.HTTPRoute
		grants     []gwv1.ReferenceGrant
		want       []models.ExposedRoute
	}{
		{
			name: "Wildcard hostname intersection",
			gateways: []gwv1.Gateway{
				{
					ObjectMeta: metav1.ObjectMeta{Name: "gw-1", Namespace: "default"},
					Spec: gwv1.GatewaySpec{
						Listeners: []gwv1.Listener{
							{Name: "http", Port: 80, Protocol: gwv1.HTTPProtocolType, Hostname: ptr(gwv1.Hostname("*.example.com"))},
						},
					},
					Status: gwv1.GatewayStatus{
						Conditions: []metav1.Condition{{Type: string(gwv1.GatewayConditionProgrammed), Status: metav1.ConditionTrue}},
						Addresses:  []gwv1.GatewayStatusAddress{{Type: ptr(gwv1.IPAddressType), Value: "10.0.0.1"}},
					},
				},
			},
			routes: []gwv1.HTTPRoute{
				{
					ObjectMeta: metav1.ObjectMeta{Name: "route-1", Namespace: "default"},
					Spec: gwv1.HTTPRouteSpec{
						CommonRouteSpec: gwv1.CommonRouteSpec{ParentRefs: []gwv1.ParentReference{{Name: "gw-1"}}},
						Hostnames:       []gwv1.Hostname{"api.example.com", "other.com", "*.api.example.com"},
						Rules: []gwv1.HTTPRouteRule{
							{
								BackendRefs: []gwv1.HTTPBackendRef{
									{BackendRef: gwv1.BackendRef{BackendObjectReference: gwv1.BackendObjectReference{Name: "svc-1"}}},
								},
							},
						},
					},
					Status: gwv1.HTTPRouteStatus{
						RouteStatus: gwv1.RouteStatus{
							Parents: []gwv1.RouteParentStatus{
								{
									ParentRef: gwv1.ParentReference{Name: "gw-1"},
									Conditions: []metav1.Condition{
										{Type: string(gwv1.RouteConditionAccepted), Status: metav1.ConditionTrue},
									},
								},
							},
						},
					},
				},
			},
			want: []models.ExposedRoute{
				{
					Namespace:        "default",
					GatewayName:      "gw-1",
					RouteName:        "route-1",
					RouteNamespace:   "default",
					Hostname:         "api.example.com", // Intersection
					ListenerName:     "http",
					HTTPPath:         "/",
					ServiceName:      "svc-1",
					BackendNamespace: "default",
					IPAddress:        "10.0.0.1",
					LBType:           "external",
				},
				{
					Namespace:        "default",
					GatewayName:      "gw-1",
					RouteName:        "route-1",
					RouteNamespace:   "default",
					Hostname:         "*.api.example.com", // Intersection
					ListenerName:     "http",
					HTTPPath:         "/",
					ServiceName:      "svc-1",
					BackendNamespace: "default",
					IPAddress:        "10.0.0.1",
					LBType:           "external",
				},
			},
		},
		{
			name: "Multiple path matches",
			gateways: []gwv1.Gateway{
				{
					ObjectMeta: metav1.ObjectMeta{Name: "gw-1", Namespace: "default"},
					Spec: gwv1.GatewaySpec{
						Listeners: []gwv1.Listener{{Name: "http", Port: 80, Protocol: gwv1.HTTPProtocolType}},
					},
					Status: gwv1.GatewayStatus{
						Conditions: []metav1.Condition{{Type: string(gwv1.GatewayConditionProgrammed), Status: metav1.ConditionTrue}},
						Addresses:  []gwv1.GatewayStatusAddress{{Type: ptr(gwv1.IPAddressType), Value: "10.0.0.1"}},
					},
				},
			},
			routes: []gwv1.HTTPRoute{
				{
					ObjectMeta: metav1.ObjectMeta{Name: "route-1", Namespace: "default"},
					Spec: gwv1.HTTPRouteSpec{
						CommonRouteSpec: gwv1.CommonRouteSpec{ParentRefs: []gwv1.ParentReference{{Name: "gw-1"}}},
						Rules: []gwv1.HTTPRouteRule{
							{
								Matches: []gwv1.HTTPRouteMatch{
									{Path: &gwv1.HTTPPathMatch{Value: ptr("/v1")}},
									{Path: &gwv1.HTTPPathMatch{Value: ptr("/v2")}},
								},
								BackendRefs: []gwv1.HTTPBackendRef{
									{BackendRef: gwv1.BackendRef{BackendObjectReference: gwv1.BackendObjectReference{Name: "svc-1"}}},
								},
							},
						},
					},
					Status: gwv1.HTTPRouteStatus{
						RouteStatus: gwv1.RouteStatus{
							Parents: []gwv1.RouteParentStatus{
								{
									ParentRef:  gwv1.ParentReference{Name: "gw-1"},
									Conditions: []metav1.Condition{{Type: string(gwv1.RouteConditionAccepted), Status: metav1.ConditionTrue}},
								},
							},
						},
					},
				},
			},
			want: []models.ExposedRoute{
				{
					Namespace:        "default",
					GatewayName:      "gw-1",
					RouteName:        "route-1",
					RouteNamespace:   "default",
					Hostname:         "",
					ListenerName:     "http",
					HTTPPath:         "/v1",
					ServiceName:      "svc-1",
					BackendNamespace: "default",
					IPAddress:        "10.0.0.1",
					LBType:           "external",
				},
				{
					Namespace:        "default",
					GatewayName:      "gw-1",
					RouteName:        "route-1",
					RouteNamespace:   "default",
					Hostname:         "",
					ListenerName:     "http",
					HTTPPath:         "/v2",
					ServiceName:      "svc-1",
					BackendNamespace: "default",
					IPAddress:        "10.0.0.1",
					LBType:           "external",
				},
			},
		},
		{
			name: "Partial cross-namespace backend via ReferenceGrant",
			gateways: []gwv1.Gateway{
				{
					ObjectMeta: metav1.ObjectMeta{Name: "gw-1", Namespace: "default"},
					Spec: gwv1.GatewaySpec{
						Listeners: []gwv1.Listener{{Name: "http", Port: 80, Protocol: gwv1.HTTPProtocolType}},
					},
					Status: gwv1.GatewayStatus{
						Conditions: []metav1.Condition{{Type: string(gwv1.GatewayConditionProgrammed), Status: metav1.ConditionTrue}},
						Addresses:  []gwv1.GatewayStatusAddress{{Type: ptr(gwv1.IPAddressType), Value: "10.0.0.1"}},
					},
				},
			},
			grants: []gwv1.ReferenceGrant{
				{
					ObjectMeta: metav1.ObjectMeta{Name: "grant-1", Namespace: "backend-ns-1"},
					Spec: gwv1.ReferenceGrantSpec{
						From: []gwv1.ReferenceGrantFrom{{Group: "gateway.networking.k8s.io", Kind: "HTTPRoute", Namespace: "default"}},
						To:   []gwv1.ReferenceGrantTo{{Group: "", Kind: "Service", Name: ptr(gwv1.ObjectName("svc-1"))}},
					},
				},
			},
			routes: []gwv1.HTTPRoute{
				{
					ObjectMeta: metav1.ObjectMeta{Name: "route-1", Namespace: "default"},
					Spec: gwv1.HTTPRouteSpec{
						CommonRouteSpec: gwv1.CommonRouteSpec{ParentRefs: []gwv1.ParentReference{{Name: "gw-1"}}},
						Rules: []gwv1.HTTPRouteRule{
							{
								BackendRefs: []gwv1.HTTPBackendRef{
									{BackendRef: gwv1.BackendRef{BackendObjectReference: gwv1.BackendObjectReference{Name: "svc-1", Namespace: ptr(gwv1.Namespace("backend-ns-1"))}}},
									{BackendRef: gwv1.BackendRef{BackendObjectReference: gwv1.BackendObjectReference{Name: "svc-2", Namespace: ptr(gwv1.Namespace("backend-ns-2"))}}}, // Missing Grant
								},
							},
						},
					},
					Status: gwv1.HTTPRouteStatus{
						RouteStatus: gwv1.RouteStatus{
							Parents: []gwv1.RouteParentStatus{
								{
									ParentRef:  gwv1.ParentReference{Name: "gw-1"},
									Conditions: []metav1.Condition{{Type: string(gwv1.RouteConditionAccepted), Status: metav1.ConditionTrue}},
								},
							},
						},
					},
				},
			},
			want: []models.ExposedRoute{
				{
					Namespace:        "default",
					GatewayName:      "gw-1",
					RouteName:        "route-1",
					RouteNamespace:   "default",
					Hostname:         "",
					ListenerName:     "http",
					HTTPPath:         "/",
					ServiceName:      "svc-1",
					BackendNamespace: "backend-ns-1",
					IPAddress:        "10.0.0.1",
					LBType:           "external",
				},
				// svc-2 is dropped because it has no ReferenceGrant!
			},
		},
		{
			name: "AllowedRoutes By Namespace Selector",
			namespaces: []corev1.Namespace{
				{ObjectMeta: metav1.ObjectMeta{Name: "allowed-ns", Labels: map[string]string{"env": "prod"}}},
				{ObjectMeta: metav1.ObjectMeta{Name: "denied-ns", Labels: map[string]string{"env": "dev"}}},
			},
			gateways: []gwv1.Gateway{
				{
					ObjectMeta: metav1.ObjectMeta{Name: "gw-1", Namespace: "default"},
					Spec: gwv1.GatewaySpec{
						Listeners: []gwv1.Listener{
							{
								Name: "http", Port: 80, Protocol: gwv1.HTTPProtocolType,
								AllowedRoutes: &gwv1.AllowedRoutes{
									Namespaces: &gwv1.RouteNamespaces{
										From:     ptr(gwv1.NamespacesFromSelector),
										Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"env": "prod"}},
									},
								},
							},
						},
					},
					Status: gwv1.GatewayStatus{
						Conditions: []metav1.Condition{{Type: string(gwv1.GatewayConditionProgrammed), Status: metav1.ConditionTrue}},
						Addresses:  []gwv1.GatewayStatusAddress{{Type: ptr(gwv1.IPAddressType), Value: "10.0.0.1"}},
					},
				},
			},
			routes: []gwv1.HTTPRoute{
				{
					ObjectMeta: metav1.ObjectMeta{Name: "route-prod", Namespace: "allowed-ns"},
					Spec: gwv1.HTTPRouteSpec{
						CommonRouteSpec: gwv1.CommonRouteSpec{ParentRefs: []gwv1.ParentReference{{Name: "gw-1", Namespace: ptr(gwv1.Namespace("default"))}}},
						Rules:           []gwv1.HTTPRouteRule{{BackendRefs: []gwv1.HTTPBackendRef{{BackendRef: gwv1.BackendRef{BackendObjectReference: gwv1.BackendObjectReference{Name: "svc-1"}}}}}},
					},
					Status: gwv1.HTTPRouteStatus{
						RouteStatus: gwv1.RouteStatus{Parents: []gwv1.RouteParentStatus{{
							ParentRef:  gwv1.ParentReference{Name: "gw-1", Namespace: ptr(gwv1.Namespace("default"))},
							Conditions: []metav1.Condition{{Type: string(gwv1.RouteConditionAccepted), Status: metav1.ConditionTrue}},
						}}},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{Name: "route-dev", Namespace: "denied-ns"},
					Spec: gwv1.HTTPRouteSpec{
						CommonRouteSpec: gwv1.CommonRouteSpec{ParentRefs: []gwv1.ParentReference{{Name: "gw-1", Namespace: ptr(gwv1.Namespace("default"))}}},
						Rules:           []gwv1.HTTPRouteRule{{BackendRefs: []gwv1.HTTPBackendRef{{BackendRef: gwv1.BackendRef{BackendObjectReference: gwv1.BackendObjectReference{Name: "svc-2"}}}}}},
					},
					Status: gwv1.HTTPRouteStatus{
						RouteStatus: gwv1.RouteStatus{Parents: []gwv1.RouteParentStatus{{
							ParentRef:  gwv1.ParentReference{Name: "gw-1", Namespace: ptr(gwv1.Namespace("default"))},
							Conditions: []metav1.Condition{{Type: string(gwv1.RouteConditionAccepted), Status: metav1.ConditionTrue}},
						}}},
					},
				},
			},
			want: []models.ExposedRoute{
				{
					Namespace:        "allowed-ns",
					GatewayName:      "gw-1",
					RouteName:        "route-prod",
					RouteNamespace:   "allowed-ns",
					Hostname:         "",
					ListenerName:     "http",
					HTTPPath:         "/",
					ServiceName:      "svc-1",
					BackendNamespace: "allowed-ns",
					IPAddress:        "10.0.0.1",
					LBType:           "external",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := fake.NewClientBuilder().WithScheme(scheme.Scheme)
			var objs []client.Object
			services := append([]corev1.Service(nil), tt.services...)
			serviceKeys := make(map[string]struct{}, len(services))
			for _, service := range services {
				serviceKeys[service.Namespace+"/"+service.Name] = struct{}{}
			}
			for _, expected := range tt.want {
				key := expected.BackendNamespace + "/" + expected.ServiceName
				if _, exists := serviceKeys[key]; exists {
					continue
				}
				service := corev1.Service{ObjectMeta: metav1.ObjectMeta{Name: expected.ServiceName, Namespace: expected.BackendNamespace}}
				if expected.BackendPort != "" {
					port, err := strconv.Atoi(expected.BackendPort)
					if err != nil {
						t.Fatalf("invalid expected backend port: %v", err)
					}
					service.Spec.Ports = []corev1.ServicePort{{Port: int32(port)}}
				}
				services = append(services, service)
				serviceKeys[key] = struct{}{}
			}
			for i := range tt.gateways {
				objs = append(objs, &tt.gateways[i])
			}
			for i := range tt.routes {
				objs = append(objs, &tt.routes[i])
			}
			for i := range tt.grants {
				objs = append(objs, &tt.grants[i])
			}
			for i := range tt.namespaces {
				objs = append(objs, &tt.namespaces[i])
			}
			for i := range services {
				objs = append(objs, &services[i])
			}
			builder.WithObjects(objs...)

			c := builder.Build()
			mapper := NewMapper(c, slog.Default())

			if err := mapper.UpdateCache(context.Background()); err != nil {
				t.Fatalf("unexpected error updating cache: %v", err)
			}
			got, err := mapper.GetExposedRoutes(context.Background())
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(got) != len(tt.want) {
				t.Fatalf("got %d routes, want %d", len(got), len(tt.want))
			}

			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("got route %d =\n%+v\nwant:\n%+v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestIntersectHostnames(t *testing.T) {
	tests := []struct {
		name     string
		route    []gwv1.Hostname
		listener *gwv1.Hostname
		want     []string
	}{
		{name: "unrestricted", want: []string{""}},
		{name: "precise mismatch", route: []gwv1.Hostname{"api.example.com"}, listener: ptr(gwv1.Hostname("admin.example.com"))},
		{name: "listener wildcard", route: []gwv1.Hostname{"deep.api.example.com"}, listener: ptr(gwv1.Hostname("*.example.com")), want: []string{"deep.api.example.com"}},
		{name: "wildcard overlap uses more specific value", route: []gwv1.Hostname{"*.example.com"}, listener: ptr(gwv1.Hostname("*.api.example.com")), want: []string{"*.api.example.com"}},
		{name: "wildcard does not match suffix root", route: []gwv1.Hostname{"example.com"}, listener: ptr(gwv1.Hostname("*.example.com"))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := intersectHostnames(tt.route, tt.listener)
			if len(got) != len(tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Fatalf("got %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestSameParentReference(t *testing.T) {
	section := gwv1.SectionName("https")
	port := gwv1.PortNumber(443)
	tests := []struct {
		name     string
		expected gwv1.ParentReference
		actual   gwv1.ParentReference
		want     bool
	}{
		{
			name:     "defaults match",
			expected: gwv1.ParentReference{Name: "gateway"},
			actual:   gwv1.ParentReference{Name: "gateway"},
			want:     true,
		},
		{
			name:     "missing status section does not match",
			expected: gwv1.ParentReference{Name: "gateway", SectionName: &section},
			actual:   gwv1.ParentReference{Name: "gateway"},
		},
		{
			name:     "missing status port does not match",
			expected: gwv1.ParentReference{Name: "gateway", Port: &port},
			actual:   gwv1.ParentReference{Name: "gateway"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sameParentReference("default", tt.expected, tt.actual); got != tt.want {
				t.Fatalf("got %t, want %t", got, tt.want)
			}
		})
	}
}

func TestListenerAllowsHTTPRouteKind(t *testing.T) {
	grpcKind := gwv1.Kind("GRPCRoute")
	httpKind := gwv1.Kind("HTTPRoute")
	tests := []struct {
		name     string
		listener gwv1.Listener
		want     bool
	}{
		{name: "default HTTP listener", listener: gwv1.Listener{Protocol: gwv1.HTTPProtocolType}, want: true},
		{name: "TCP listener", listener: gwv1.Listener{Protocol: gwv1.TCPProtocolType}},
		{name: "HTTPRoute excluded", listener: gwv1.Listener{Protocol: gwv1.HTTPProtocolType, AllowedRoutes: &gwv1.AllowedRoutes{Kinds: []gwv1.RouteGroupKind{{Kind: grpcKind}}}}},
		{name: "HTTPRoute included", listener: gwv1.Listener{Protocol: gwv1.HTTPProtocolType, AllowedRoutes: &gwv1.AllowedRoutes{Kinds: []gwv1.RouteGroupKind{{Kind: httpKind}}}}, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := listenerAllowsHTTPRouteKind(tt.listener); got != tt.want {
				t.Fatalf("got %t, want %t", got, tt.want)
			}
		})
	}
}

func TestGatewayIPAddress(t *testing.T) {
	hostnameType := gwv1.HostnameAddressType
	ipType := gwv1.IPAddressType
	gateway := &gwv1.Gateway{Status: gwv1.GatewayStatus{Addresses: []gwv1.GatewayStatusAddress{
		{Type: nil, Value: "ignored"},
		{Type: &hostnameType, Value: "gateway.example.com"},
		{Type: &ipType, Value: "192.0.2.10"},
	}}}
	if got := gatewayIPAddress(gateway); got != "192.0.2.10" {
		t.Fatalf("got %q, want %q", got, "192.0.2.10")
	}
	if got := gatewayIPAddress(&gwv1.Gateway{}); got != "" {
		t.Fatalf("got %q for Gateway without an IP, want empty string", got)
	}
}

func TestRouteAcceptanceRequiresCurrentGeneration(t *testing.T) {
	route := &gwv1.HTTPRoute{
		ObjectMeta: metav1.ObjectMeta{Name: "route", Namespace: "default", Generation: 2},
		Status: gwv1.HTTPRouteStatus{RouteStatus: gwv1.RouteStatus{Parents: []gwv1.RouteParentStatus{{
			ParentRef: gwv1.ParentReference{Name: "gateway"},
			Conditions: []metav1.Condition{{
				Type:               string(gwv1.RouteConditionAccepted),
				Status:             metav1.ConditionTrue,
				ObservedGeneration: 1,
			}},
		}}}},
	}
	parent := gwv1.ParentReference{Name: "gateway"}
	if isRouteAcceptedByParent(route, parent) {
		t.Fatal("stale Accepted condition must not be treated as current")
	}
	route.Status.Parents[0].Conditions[0].ObservedGeneration = 2
	if !isRouteAcceptedByParent(route, parent) {
		t.Fatal("current Accepted condition should be accepted")
	}
}

func TestServiceBackendExistsRequiresRequestedPort(t *testing.T) {
	testScheme := runtime.NewScheme()
	if err := corev1.AddToScheme(testScheme); err != nil {
		t.Fatalf("add core scheme: %v", err)
	}
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{Name: "api", Namespace: "default"},
		Spec:       corev1.ServiceSpec{Ports: []corev1.ServicePort{{Port: 8080}}},
	}
	client := fake.NewClientBuilder().WithScheme(testScheme).WithObjects(service).Build()
	mapper := NewMapper(client, slog.Default())

	valid := gwv1.HTTPBackendRef{BackendRef: gwv1.BackendRef{BackendObjectReference: gwv1.BackendObjectReference{Name: "api", Port: ptr(gwv1.PortNumber(8080))}}}
	if !mapper.serviceBackendExists(context.Background(), "default", valid) {
		t.Fatal("existing Service port should be resolved")
	}
	invalid := valid
	invalid.Port = ptr(gwv1.PortNumber(9090))
	if mapper.serviceBackendExists(context.Background(), "default", invalid) {
		t.Fatal("missing Service port must not be resolved")
	}
}

func TestMapper_Concurrency(t *testing.T) {
	testScheme := runtime.NewScheme()
	_ = corev1.AddToScheme(testScheme)
	_ = gwv1.Install(testScheme)

	gateway := &gwv1.Gateway{
		ObjectMeta: metav1.ObjectMeta{Name: "gw-1", Namespace: "default"},
		Spec: gwv1.GatewaySpec{
			Listeners: []gwv1.Listener{{Name: "http", Port: 80, Protocol: gwv1.HTTPProtocolType}},
		},
		Status: gwv1.GatewayStatus{
			Conditions: []metav1.Condition{{Type: string(gwv1.GatewayConditionProgrammed), Status: metav1.ConditionTrue}},
			Addresses:  []gwv1.GatewayStatusAddress{{Type: ptr(gwv1.IPAddressType), Value: "10.0.0.1"}},
		},
	}
	route := &gwv1.HTTPRoute{
		ObjectMeta: metav1.ObjectMeta{Name: "route-1", Namespace: "default"},
		Spec: gwv1.HTTPRouteSpec{
			CommonRouteSpec: gwv1.CommonRouteSpec{ParentRefs: []gwv1.ParentReference{{Name: "gw-1"}}},
			Rules: []gwv1.HTTPRouteRule{
				{
					BackendRefs: []gwv1.HTTPBackendRef{
						{BackendRef: gwv1.BackendRef{BackendObjectReference: gwv1.BackendObjectReference{Name: "svc-1"}}},
					},
				},
			},
		},
		Status: gwv1.HTTPRouteStatus{
			RouteStatus: gwv1.RouteStatus{
				Parents: []gwv1.RouteParentStatus{{
					ParentRef:  gwv1.ParentReference{Name: "gw-1"},
					Conditions: []metav1.Condition{{Type: string(gwv1.RouteConditionAccepted), Status: metav1.ConditionTrue}},
				}},
			},
		},
	}
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{Name: "svc-1", Namespace: "default"},
		Spec:       corev1.ServiceSpec{Ports: []corev1.ServicePort{{Port: 80}}},
	}

	c := fake.NewClientBuilder().WithScheme(testScheme).WithObjects(gateway, route, service).Build()
	m := NewMapper(c, slog.Default())

	var wg sync.WaitGroup
	ctx := context.Background()

	// 50 concurrent writers calling UpdateCache
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = m.UpdateCache(ctx)
		}()
	}

	// 100 concurrent readers calling GetExposedRoutes
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = m.GetExposedRoutes(ctx)
		}()
	}

	wg.Wait()
}
