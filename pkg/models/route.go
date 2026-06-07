package models

// ExposedRoute represents a flattened logical Gateway API route-to-Service relationship.
// See docs/index.md for the core domain relationship mapping strategy.
type ExposedRoute struct {
	Namespace        string
	GatewayName      string
	RouteName        string
	RouteNamespace   string
	Hostname         string
	ListenerName     string
	HTTPPath         string
	ServiceName      string
	BackendNamespace string
	BackendPort      string
	LBType           string // "internal" or "external"
	IPAddress        string
}
