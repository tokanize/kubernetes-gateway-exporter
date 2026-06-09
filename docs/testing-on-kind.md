# Testing Locally on kind

This guide verifies Kubernetes Gateway Exporter on a local `kind` cluster without installing a Gateway controller. The controller behavior required by the exporter is simulated by writing an `Accepted=True` condition to the HTTPRoute status.

The test validates a logical `Gateway -> Listener -> HTTPRoute -> Service` relationship. It does not create workload Pods and does not test Pod or EndpointSlice enumeration, because the exporter does not expose those objects as metric series.

## Prerequisites

- Docker with BuildKit
- kind
- Helm
- kubectl
- curl

The repository currently uses Gateway API `v1.5.1`. Keep the CRD version below aligned with `sigs.k8s.io/gateway-api` in `go.mod`.

## 1. Create the cluster and install Gateway API

```bash
kind create cluster --name gateway-test

kubectl apply --server-side -f \
  https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.5.1/standard-install.yaml

kubectl wait --for=condition=Established \
  crd/gateways.gateway.networking.k8s.io \
  crd/httproutes.gateway.networking.k8s.io \
  crd/referencegrants.gateway.networking.k8s.io \
  --timeout=60s
```

## 2. Build and load the image

The Dockerfile uses BuildKit's target architecture, so this works with both amd64 and arm64 kind nodes.

```bash
docker build -t my-registry/kubernetes-gateway-exporter:latest .
kind load docker-image my-registry/kubernetes-gateway-exporter:latest \
  --name gateway-test
```

## 3. Deploy the exporter

The local cluster does not include the Prometheus Operator CRDs, so disable `ServiceMonitor` creation.

```bash
helm upgrade --install exporter ./deploy/chart/kubernetes-gateway-exporter \
  --namespace monitoring \
  --create-namespace \
  --set serviceMonitor.enabled=false

kubectl rollout status deployment/exporter-kubernetes-gateway-exporter \
  -n monitoring \
  --timeout=90s
```

If the Pod does not start, inspect it before continuing:

```bash
kubectl get pods -n monitoring
kubectl describe pod -n monitoring -l app=kubernetes-gateway-exporter
kubectl logs -n monitoring -l app=kubernetes-gateway-exporter
```

## 4. Create the test resources

The repository contains [`examples/dummy-resources.yaml`](https://github.com/tokanize/kubernetes-gateway-exporter/blob/main/examples/dummy-resources.yaml), which creates namespace `test-ns`, Service `test-svc`, GatewayClass `my-gc`, Gateway `my-gateway`, and HTTPRoute `my-route` matching `/api`.

```bash
kubectl apply -f examples/dummy-resources.yaml
```

No Gateway controller handles `example.com/gateway-controller`, so the HTTPRoute remains unaccepted until its status is patched.

## 5. Simulate Gateway controller acceptance

The exporter requires an `Accepted=True` condition whose `observedGeneration` matches the current HTTPRoute generation. Generate the patch dynamically instead of hard-coding generation `1`:

```bash
ROUTE_GENERATION=$(kubectl get httproute my-route -n test-ns \
  -o jsonpath='{.metadata.generation}')

TRANSITION_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ)

kubectl patch httproute my-route \
  -n test-ns \
  --type=merge \
  --subresource=status \
  -p "{\"status\":{\"parents\":[{\"controllerName\":\"example.com/gateway-controller\",\"parentRef\":{\"name\":\"my-gateway\"},\"conditions\":[{\"type\":\"Accepted\",\"status\":\"True\",\"reason\":\"Accepted\",\"message\":\"Manually patched for kind testing\",\"observedGeneration\":${ROUTE_GENERATION},\"lastTransitionTime\":\"${TRANSITION_TIME}\"}]}]}}"
```

Verify the status:

```bash
kubectl get httproute my-route -n test-ns -o yaml
```

The tracked [`examples/patch-route.yaml`](https://github.com/tokanize/kubernetes-gateway-exporter/blob/main/examples/patch-route.yaml) is an example status payload. The dynamic command is preferred because Route generation may change after edits.

## 6. Verify the metric

Port-forward the exporter Service:

```bash
kubectl port-forward \
  -n monitoring \
  service/exporter-kubernetes-gateway-exporter \
  8080:8080
```

In another terminal:

```bash
curl -fsS http://127.0.0.1:8080/metrics | grep '^exposed_route_info'
```

Expected series:

```text
exposed_route_info{backend_namespace="test-ns",backend_port="8080",gateway_name="my-gateway",hostname="",http_path="/api",ip_address="",lb_type="external",listener_name="http",namespace="test-ns",route_name="my-route",service_name="test-svc"} 1
```

The empty `ip_address` is expected because no Gateway controller allocates an address in this test.

## 7. Verify that acceptance is required

Remove the status and confirm that the series disappears:

```bash
kubectl patch httproute my-route \
  -n test-ns \
  --type=merge \
  --subresource=status \
  -p '{"status":{"parents":[]}}'

curl -fsS http://127.0.0.1:8080/metrics | grep '^exposed_route_info' || true
```

## 8. Clean up

Stop the port-forward with `Ctrl+C`, then delete the cluster:

```bash
kind delete cluster --name gateway-test
```
