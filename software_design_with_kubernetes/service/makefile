
# ==============================================================================
# Define dependencies

GOLANG          := golang:1.24
ALPINE          := alpine:3.21
KIND            := kindest/node:v1.32.0
POSTGRES        := postgres:17.3
GRAFANA         := grafana/grafana:11.5.0
PROMETHEUS      := prom/prometheus:v3.1.0
TEMPO           := grafana/tempo:2.7.0
LOKI            := grafana/loki:3.4.0
PROMTAIL        := grafana/promtail:3.4.0

KIND_CLUSTER    := dp-cluster
NAMESPACE       := sales-system
SALES_APP       := sales
AUTH_APP        := auth
BASE_IMAGE_NAME := localhost/4925k
VERSION         := 0.0.1
SALES_IMAGE     := $(BASE_IMAGE_NAME)/$(SALES_APP):$(VERSION)
METRICS_IMAGE   := $(BASE_IMAGE_NAME)/metrics:$(VERSION)
AUTH_IMAGE      := $(BASE_IMAGE_NAME)/$(AUTH_APP):$(VERSION)


# VERSION       := "0.0.1-$(shell git rev-parse --short HEAD)"

run: 
	go run app/services/sales/main.go | go run app/tooling/logfmt/main.go

run-help:
	go run app/service/sales/main.go --help | go run app/tooling/logfmt/main.go


# ==============================================================================
# Running from within k8s/kind

dev-up:
	kind create cluster \
		--image $(KIND) \
		--name $(KIND_CLUSTER) \
		--config zarf/k8s/dev/kind-config.yaml

	kubectl wait --timeout=120s --namespace=local-path-storage --for=condition=Available deployment/local-path-provisioner

dev-down:
	kind delete cluster --name $(KIND_CLUSTER)

dev-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

# ==============================================================================
# Building containers

build: sales

sales:
	docker build \
		-f zarf/docker/dockerfile.sales \
		-t $(SALES_IMAGE) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
		.

#===============================================================================

dev-load:
	kind load docker-image $(SALES_IMAGE) --name $(KIND_CLUSTER)

dev-apply:
	kustomize build zarf/k8s/dev/sales | kubectl apply -f -
	kubectl wait pods --namespace=$(NAMESPACE) --selector app=$(SALES_APP) --timeout=120s --for=condition=Ready

dev-restart:
	kubectl rollout restart deployment $(SALES_APP) --namespace=$(NAMESPACE)

dev-update: build dev-load dev-restart

dev-update-apply: build dev-load dev-apply

#===============================================================================

dev-logs: 
	kubectl logs --namespace=$(NAMESPACE) -l app=$(SALES_APP) --all-containers=true -f --tail=100 --max-log-requests=6 | go run app/tooling/logfmt/main.go -service=$(SALES_APP)

dev-describe-deployment:
	kubectl describe deployment --namespace=$(NAMESPACE) $(SALES_APP)

dev-describe-sales:
	kubectl describe pod --namespace=$(NAMESPACE) -l app=$(SALES_APP)


#===============================================================================
# Modules Support

tidy:
	go mod tidy
	go mod vendor	


#===============================================================================
# Metrics and Tracing

metrics-view-sc:
	expvarmon -ports="localhost:4000" -vars="build,requests,goroutiens,errors,panics,mem:memstats.Alloc"


curl:
	curl -il http://localhost:3000/hack

load:
	hey -m GET -c 100 -n 100000 "http://localhost:3000/hack"