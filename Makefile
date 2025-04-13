# # Set Docker registry and version
# REGISTRY=hsri
# VERSION=v1.0.0

# # Kubernetes namespace (change if needed)
# NAMESPACE=default

# # 🔹 Build Docker Images
# build-rest:
# 	docker build -t $(REGISTRY)/rest-api:$(VERSION) -f deploy/docker/Dockerfile.rest-api .

# build-grpc:
# 	docker build -t $(REGISTRY)/grpc-server:$(VERSION) -f deploy/docker/Dockerfile.grpc .

# build-frontend:
# 	docker build -t $(REGISTRY)/frontend:$(VERSION) -f deploy/docker/Dockerfile.frontend .

# # 🔹 Push Images to Registry
# push-rest: build-rest
# 	docker push $(REGISTRY)/rest-api:$(VERSION)

# push-grpc: build-grpc
# 	docker push $(REGISTRY)/grpc-server:$(VERSION)

# push-frontend: build-frontend
# 	docker push $(REGISTRY)/frontend:$(VERSION)

# # 🔹 Apply Kubernetes ConfigMaps and Secrets
# deploy-config:
# 	kubectl apply -f deploy/kubernetes/config/config-db.yaml -n $(NAMESPACE)
# 	kubectl apply -f deploy/kubernetes/config/secret-db.yaml -n $(NAMESPACE)

# # 🔹 Deploy Services to Kubernetes
# deploy-rest: push-rest
# 	kubectl apply -f deploy/kubernetes/deployment/deployment-rest.yaml -n $(NAMESPACE)

# deploy-grpc: push-grpc
# 	kubectl apply -f deploy/kubernetes/deployment/deployment-grpc.yaml -n $(NAMESPACE)

# deploy-frontend: push-frontend
# 	kubectl apply -f deploy/kubernetes/deployment/deployment-frontend.yaml -n $(NAMESPACE)

# # 🔹 Deploy Everything (Full Deployment Pipeline)
# deploy-all: deploy-config deploy-rest deploy-grpc deploy-frontend

# # 🔹 Check Kubernetes Status
# status:
# 	kubectl get pods -n $(NAMESPACE)
# 	kubectl get services -n $(NAMESPACE)
# 	kubectl logs -l app=rest-api -n $(NAMESPACE) -f
# 	kubectl logs -l app=grpc-server -n $(NAMESPACE) -f

# # 🔹 Delete Deployments, Services, and Configs
# clean:
# 	kubectl delete --ignore-not-found=true -f deploy/kubernetes/deployment/ -n $(NAMESPACE)
# 	kubectl delete --ignore-not-found=true -f deploy/kubernetes/config/ -n $(NAMESPACE)


# Set Docker registry and version
REGISTRY=hsri
VERSION=v2.0.0
NAMESPACE=default

# 🔹 Build, Push, and Deploy in One Step
all: build push deploy

# 🔹 Build Docker Images
build:
	docker build -t $(REGISTRY)/rest-api:$(VERSION) -f deploy/docker/Dockerfile.rest-api .
	docker build -t $(REGISTRY)/grpc-server:$(VERSION) -f deploy/docker/Dockerfile.grpc .
	docker build -t $(REGISTRY)/frontend:$(VERSION) -f deploy/docker/Dockerfile.frontend .

# 🔹 Push Images to Registry
push: build
	docker push $(REGISTRY)/rest-api:$(VERSION)
	docker push $(REGISTRY)/grpc-server:$(VERSION)
	docker push $(REGISTRY)/frontend:$(VERSION)

# 🔹 Deploy Kubernetes Configs and Services
deploy: deploy-config deploy-apps

deploy-config:
	kubectl apply -f deploy/kubernetes/config/ -n $(NAMESPACE)

deploy-apps:
	kubectl apply -f deploy/kubernetes/deployment/ -n $(NAMESPACE)

# 🔹 Check Kubernetes Status
status:
	kubectl get pods,services -n $(NAMESPACE)
	kubectl logs -l app=rest-api -n $(NAMESPACE) -f &
	kubectl logs -l app=grpc-server -n $(NAMESPACE) -f &

# 🔹 Cleanup Everything (Deployments, Services, Configs, Images)
clean:
	kubectl delete --ignore-not-found=true -f deploy/kubernetes/deployment/ -n $(NAMESPACE)
	kubectl delete --ignore-not-found=true -f deploy/kubernetes/config/ -n $(NAMESPACE)
	kubectl delete namespace $(NAMESPACE) --ignore-not-found=true
	docker rmi -f $(REGISTRY)/rest-api:$(VERSION) $(REGISTRY)/grpc-server:$(VERSION) $(REGISTRY)/frontend:$(VERSION) || true
