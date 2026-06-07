package metrics

import (
	"log/slog"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/tokanize/kubernetes-gateway-exporter/internal/kubernetes/mapper"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func TestExporterGatherIncludesCompleteLabelSet(t *testing.T) {
	scheme := runtime.NewScheme()
	if err := corev1.AddToScheme(scheme); err != nil {
		t.Fatalf("install core Kubernetes scheme: %v", err)
	}
	if err := gwv1.Install(scheme); err != nil {
		t.Fatalf("install Gateway API scheme: %v", err)
	}

	ipType := gwv1.IPAddressType
	gateway := &gwv1.Gateway{
		ObjectMeta: metav1.ObjectMeta{Name: "public", Namespace: "default"},
		Spec: gwv1.GatewaySpec{
			GatewayClassName: "external",
			Listeners:        []gwv1.Listener{{Name: "http", Protocol: gwv1.HTTPProtocolType, Port: 80}},
		},
		Status: gwv1.GatewayStatus{Addresses: []gwv1.GatewayStatusAddress{{Type: &ipType, Value: "192.0.2.10"}}},
	}
	route := &gwv1.HTTPRoute{
		ObjectMeta: metav1.ObjectMeta{Name: "api", Namespace: "default", Generation: 1},
		Spec: gwv1.HTTPRouteSpec{
			CommonRouteSpec: gwv1.CommonRouteSpec{ParentRefs: []gwv1.ParentReference{{Name: "public"}}},
			Rules:           []gwv1.HTTPRouteRule{{BackendRefs: []gwv1.HTTPBackendRef{{BackendRef: gwv1.BackendRef{BackendObjectReference: gwv1.BackendObjectReference{Name: "api", Port: ptr(gwv1.PortNumber(8080))}}}}}},
		},
		Status: gwv1.HTTPRouteStatus{RouteStatus: gwv1.RouteStatus{Parents: []gwv1.RouteParentStatus{{
			ParentRef:  gwv1.ParentReference{Name: "public"},
			Conditions: []metav1.Condition{{Type: string(gwv1.RouteConditionAccepted), Status: metav1.ConditionTrue, ObservedGeneration: 1}},
		}}}},
	}

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{Name: "api", Namespace: "default"},
		Spec:       corev1.ServiceSpec{Ports: []corev1.ServicePort{{Port: 8080}}},
	}
	client := fake.NewClientBuilder().WithScheme(scheme).WithObjects(gateway, route, service).Build()
	collector := NewExporter(mapper.NewMapper(client, slog.Default()), slog.Default())
	registry := prometheus.NewPedanticRegistry()
	if err := registry.Register(collector); err != nil {
		t.Fatalf("register collector: %v", err)
	}

	families, err := registry.Gather()
	if err != nil {
		t.Fatalf("gather metrics: %v", err)
	}
	if len(families) != 1 || len(families[0].Metric) != 1 {
		t.Fatalf("unexpected metric families: %+v", families)
	}

	wantLabels := map[string]string{
		"namespace": "default", "gateway_name": "public", "route_name": "api",
		"route_namespace": "default", "hostname": "", "listener_name": "http",
		"http_path": "/", "service_name": "api", "backend_namespace": "default",
		"backend_port": "8080", "lb_type": "external", "ip_address": "192.0.2.10",
	}
	for _, pair := range families[0].Metric[0].Label {
		if want, ok := wantLabels[pair.GetName()]; !ok || want != pair.GetValue() {
			t.Fatalf("unexpected label %s=%q", pair.GetName(), pair.GetValue())
		}
		delete(wantLabels, pair.GetName())
	}
	if len(wantLabels) != 0 {
		t.Fatalf("missing labels: %v", wantLabels)
	}
}

func ptr[T any](value T) *T { return &value }
