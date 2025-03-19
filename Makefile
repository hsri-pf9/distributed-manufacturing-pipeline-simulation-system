# Set Docker registry and version
REGISTRY=hsri
VERSION=v1.0.0

# Kubernetes namespace (change if needed)
NAMESPACE=default

# 🔹 Build Docker Images
build-rest:
	docker build -t $(REGISTRY)/rest-api:$(VERSION) -f deploy/docker/Dockerfile.rest-api .

build-grpc:
	docker build -t $(REGISTRY)/grpc-server:$(VERSION) -f deploy/docker/Dockerfile.grpc .

build-frontend:
	docker build -t $(REGISTRY)/frontend:$(VERSION) -f deploy/docker/Dockerfile.frontend .

build-democtl:
	docker build -t $(REGISTRY)/democtl:$(VERSION) -f deploy/docker/Dockerfile.democtl .

# 🔹 Push Images to Registry
push-rest: build-rest
	docker push $(REGISTRY)/rest-api:$(VERSION)

push-grpc: build-grpc
	docker push $(REGISTRY)/grpc-server:$(VERSION)

push-frontend: build-frontend
	docker push $(REGISTRY)/frontend:$(VERSION)

push-democtl: build-democtl
	docker push $(REGISTRY)/democtl:$(VERSION)

# 🔹 Apply Kubernetes ConfigMaps and Secrets
deploy-config:
	kubectl apply -f deploy/kubernetes/config/config-db.yaml -n $(NAMESPACE)
	kubectl apply -f deploy/kubernetes/config/secret-db.yaml -n $(NAMESPACE)

# 🔹 Deploy Services to Kubernetes
deploy-rest: push-rest
	kubectl apply -f deploy/kubernetes/deployment/deployment-rest.yaml -n $(NAMESPACE)

deploy-grpc: push-grpc
	kubectl apply -f deploy/kubernetes/deployment/deployment-grpc.yaml -n $(NAMESPACE)

deploy-frontend: push-frontend
	kubectl apply -f deploy/kubernetes/deployment/deployment-frontend.yaml -n $(NAMESPACE)

deploy-democtl: push-democtl
	kubectl apply -f deploy/kubernetes/jobs/democtl-job.yaml -n $(NAMESPACE)

# 🔹 Deploy Everything (Full Deployment Pipeline)
deploy-all: deploy-config deploy-rest deploy-grpc deploy-frontend deploy-democtl

# 🔹 Check Kubernetes Status
status:
	kubectl get pods -n $(NAMESPACE)
	kubectl get services -n $(NAMESPACE)
	kubectl logs -l app=rest-api -n $(NAMESPACE) --tail=20
	kubectl logs -l app=grpc-server -n $(NAMESPACE) --tail=20

# 🔹 Delete Deployments, Services, and Configs
clean:
	kubectl delete --ignore-not-found=true -f deploy/kubernetes/deployment/ -n $(NAMESPACE)
	kubectl delete --ignore-not-found=true -f deploy/kubernetes/config/ -n $(NAMESPACE)
	kubectl delete --ignore-not-found=true job democtl-job -n $(NAMESPACE)
