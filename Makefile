
# Image URL to use all building/pushing image targets
CLUSTER_NAME ?= encoder-run-local
IMG ?= controller:latest
CONTROLLER_MANAGER_IMG ?= controller
CONTROLLER_MANAGER_IMG_VERSION ?= 0.0.1
GATEWAY_IMG ?= gateway
GATEWAY_IMG_VERSION ?= 0.0.1
CONSOLE_UI_IMG ?= console-ui
CONSOLE_UI_IMG_VERSION ?= 0.0.1
MODELDEPLOYER_IMG ?= codeembedder
MODELDEPLOYER_IMG_VERSION ?= 0.0.1
REPOSITORY_EMBEDDER_IMG ?= repository-embedder
REPOSITORY_EMBEDDER_IMG_VERSION ?= 0.0.1
# ENVTEST_K8S_VERSION refers to the version of kubebuilder assets to be downloaded by envtest binary.
ENVTEST_K8S_VERSION = 1.28.0
ISTIO_VERSION=1.17.2
ISTIO=$(LOCALBIN)/istio-$(ISTIO_VERSION)

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# CONTAINER_TOOL defines the container tool to be used for building images.
# Be aware that the target commands are only tested with Docker which is
# scaffolded by default. However, you might want to replace it to use other
# tools. (i.e. podman)
CONTAINER_TOOL ?= docker

# Setting SHELL to bash allows bash commands to be executed by recipes.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

.PHONY: all
all: build

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk command is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: manifests
manifests: controller-gen ## Generate WebhookConfiguration, ClusterRole and CustomResourceDefinition objects.
	$(CONTROLLER_GEN) rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases

.PHONY: generate
generate: controller-gen ## Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: test
test: manifests generate fmt vet envtest ## Run tests.
	KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use $(ENVTEST_K8S_VERSION) --bin-dir $(LOCALBIN) -p path)" go test ./... -coverprofile cover.out

##@ Build

.PHONY: build
build: manifests generate fmt vet ## Build manager binary.
	go build -o bin/manager cmd/main.go

.PHONY: run
run: manifests generate fmt vet ## Run a controller from your host.
	go run ./cmd/main.go

# If you wish to build the manager image targeting other platforms you can use the --platform flag.
# (i.e. docker build --platform linux/arm64). However, you must enable docker buildKit for it.
# More info: https://docs.docker.com/develop/develop-images/build_enhancements/
.PHONY: docker-build
docker-build: ## Build docker image with the manager.
	$(CONTAINER_TOOL) build -t ${CONTROLLER_MANAGER_IMG}:${CONTROLLER_MANAGER_IMG_VERSION} .

.PHONY: docker-push
docker-push: ## Push docker image with the manager.
	$(CONTAINER_TOOL) push ${CONTROLLER_MANAGER_IMG}:${CONTROLLER_MANAGER_IMG_VERSION}

.PHONY: gateway-build
gateway-build: ## Build the gateway Docker image.
	$(CONTAINER_TOOL) build -t ${GATEWAY_IMG}:${GATEWAY_IMG_VERSION} -f cmd/gateway/Dockerfile .

.PHONY: gateway-push
gateway-push: ## Push the gateway Docker image.
	$(CONTAINER_TOOL) push ${GATEWAY_IMG}:${GATEWAY_IMG_VERSION}

.PHONY: console-ui-build
console-ui-build: ## Build the console-ui Docker image.
	$(CONTAINER_TOOL) build -t ${CONSOLE_UI_IMG}:${CONSOLE_UI_IMG_VERSION} -f frontend/console-ui/Dockerfile frontend/console-ui

.PHONY: modeldeployer-docker-build
modeldeployer-docker-build: ## Build the codeembedder Docker image.
	docker build -t ${MODELDEPLOYER_IMG}:${MODELDEPLOYER_IMG_VERSION} -f cmd/modeldeployer/Dockerfile .

.PHONY: repoembedder-build
repoembedder-build: ## Build the repositoryembedder Docker image.
	$(CONTAINER_TOOL) build -t ${REPOSITORY_EMBEDDER_IMG}:${REPOSITORY_EMBEDDER_IMG_VERSION} -f cmd/repositoryembedder/Dockerfile .


# PLATFORMS defines the target platforms for the manager image be built to provide support to multiple
# architectures. (i.e. make docker-buildx IMG=myregistry/mypoperator:0.0.1). To use this option you need to:
# - be able to use docker buildx. More info: https://docs.docker.com/build/buildx/
# - have enabled BuildKit. More info: https://docs.docker.com/develop/develop-images/build_enhancements/
# - be able to push the image to your registry (i.e. if you do not set a valid value via IMG=<myregistry/image:<tag>> then the export will fail)
# To adequately provide solutions that are compatible with multiple platforms, you should consider using this option.
PLATFORMS ?= linux/arm64,linux/amd64,linux/s390x,linux/ppc64le
.PHONY: docker-buildx
docker-buildx: ## Build and push docker image for the manager for cross-platform support
	# copy existing Dockerfile and insert --platform=${BUILDPLATFORM} into Dockerfile.cross, and preserve the original Dockerfile
	sed -e '1 s/\(^FROM\)/FROM --platform=\$$\{BUILDPLATFORM\}/; t' -e ' 1,// s//FROM --platform=\$$\{BUILDPLATFORM\}/' Dockerfile > Dockerfile.cross
	- $(CONTAINER_TOOL) buildx create --name project-v3-builder
	$(CONTAINER_TOOL) buildx use project-v3-builder
	- $(CONTAINER_TOOL) buildx build --push --platform=$(PLATFORMS) --tag ${IMG} -f Dockerfile.cross .
	- $(CONTAINER_TOOL) buildx rm project-v3-builder
	rm Dockerfile.cross

##@ Deployment

ifndef ignore-not-found
  ignore-not-found = false
endif

.PHONY: install
install: manifests kustomize ## Install CRDs into the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/crd | $(KUBECTL) apply -f -

.PHONY: uninstall
uninstall: manifests kustomize ## Uninstall CRDs from the K8s cluster specified in ~/.kube/config. Call with ignore-not-found=true to ignore resource not found errors during deletion.
	$(KUSTOMIZE) build config/crd | $(KUBECTL) delete --ignore-not-found=$(ignore-not-found) -f -

.PHONY: deploy
deploy: kind-install cert-manager-install manifests istio kustomize docker-build gateway-build console-ui-build modeldeployer-docker-build repoembedder-build ## Deploy controller to the K8s cluster specified in ~/.kube/config.
ifeq ("$(wildcard /tmp/encoder-run)", "")
	@echo "kind storage folder does not exists, please execute the following command: mkdir -p /tmp/encoder-run";
else
	make kind-delete
	make kind
	$(ISTIO)/bin/istioctl x precheck
	$(ISTIO)/bin/istioctl install --set profile=default -y
	make install
	$(CMCTL) x install --set prometheus.enabled=false
	$(KUSTOMIZE) build config/kserve | $(KUBECTL) apply -f -
	kind load docker-image $(CONTROLLER_MANAGER_IMG):$(CONTROLLER_MANAGER_IMG_VERSION) --name=$(CLUSTER_NAME)
	kind load docker-image $(GATEWAY_IMG):$(GATEWAY_IMG_VERSION) --name=$(CLUSTER_NAME)
	kind load docker-image $(CONSOLE_UI_IMG):$(CONSOLE_UI_IMG_VERSION) --name=$(CLUSTER_NAME)
	kind load docker-image $(MODELDEPLOYER_IMG):$(MODELDEPLOYER_IMG_VERSION) --name=$(CLUSTER_NAME)
	kind load docker-image $(REPOSITORY_EMBEDDER_IMG):$(REPOSITORY_EMBEDDER_IMG_VERSION) --name=$(CLUSTER_NAME)
	cd config/manager && $(KUSTOMIZE) edit set image controller=$(CONTROLLER_MANAGER_IMG):$(CONTROLLER_MANAGER_IMG_VERSION)
	$(KUSTOMIZE) build config/default | $(KUBECTL) apply -f -
	make default-admin
	make default-redis-db
endif

.PHONY: undeploy
undeploy: ## Undeploy controller from the K8s cluster specified in ~/.kube/config. Call with ignore-not-found=true to ignore resource not found errors during deletion.
	$(KUSTOMIZE) build config/default | $(KUBECTL) delete --ignore-not-found=$(ignore-not-found) -f -

##@ Build Dependencies

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
KUBECTL ?= kubectl
KUSTOMIZE ?= $(LOCALBIN)/kustomize
CMCTL ?= $(LOCALBIN)/cmctl
CONTROLLER_GEN ?= $(LOCALBIN)/controller-gen
ENVTEST ?= $(LOCALBIN)/setup-envtest

## Tool Versions
KUSTOMIZE_VERSION ?= v5.1.1
CONTROLLER_TOOLS_VERSION ?= v0.14.0

.PHONY: kustomize
kustomize: $(KUSTOMIZE) ## Download kustomize locally if necessary. If wrong version is installed, it will be removed before downloading.
$(KUSTOMIZE): $(LOCALBIN)
	@if test -x $(LOCALBIN)/kustomize && ! $(LOCALBIN)/kustomize version | grep -q $(KUSTOMIZE_VERSION); then \
		echo "$(LOCALBIN)/kustomize version is not expected $(KUSTOMIZE_VERSION). Removing it before installing."; \
		rm -rf $(LOCALBIN)/kustomize; \
	fi
	test -s $(LOCALBIN)/kustomize || GOBIN=$(LOCALBIN) GO111MODULE=on go install sigs.k8s.io/kustomize/kustomize/v5@$(KUSTOMIZE_VERSION)

.PHONY: controller-gen
controller-gen: $(CONTROLLER_GEN) ## Download controller-gen locally if necessary. If wrong version is installed, it will be overwritten.
$(CONTROLLER_GEN): $(LOCALBIN)
	test -s $(LOCALBIN)/controller-gen && $(LOCALBIN)/controller-gen --version | grep -q $(CONTROLLER_TOOLS_VERSION) || \
	GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_TOOLS_VERSION)

.PHONY: envtest
envtest: $(ENVTEST) ## Download envtest-setup locally if necessary.
$(ENVTEST): $(LOCALBIN)
	test -s $(LOCALBIN)/setup-envtest || GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-runtime/tools/setup-envtest@latest

.PHONY: istio
istio: ## Download istio if it doesn't already exist
	@if [ ! -d "$(ISTIO)" ]; then \
		echo "Istio $(ISTIO_VERSION) not found. Downloading and installing..."; \
		curl -L https://istio.io/downloadIstio | ISTIO_VERSION=$(ISTIO_VERSION) sh -; \
		rm -rf $(ISTIO); \
		mv istio-$(ISTIO_VERSION) $(ISTIO); \
	else \
		echo "Istio $(ISTIO_VERSION) already installed at $(ISTIO)."; \
	fi

.PHONY: cert-manager-install
cert-manager-install: ## Install cert-manager
	@if [ ! -x $(LOCALBIN)/cmctl ]; then \
		echo "cmctl not installed. Installing..."; \
		OS=$$(go env GOOS); ARCH=$$(go env GOARCH); \
		URL="https://github.com/cert-manager/cmctl/releases/download/v2.0.0/cmctl_$${OS}_$${ARCH}"; \
		echo "Downloading cmctl from $${URL}"; \
		curl -fsSL -o cmctl $${URL} || { echo "Download failed"; exit 1; }; \
		chmod +x ./cmctl; \
		sudo mv ./cmctl $(LOCALBIN)/cmctl; \
	else \
		echo "cmctl is already installed."; \
	fi

.PHONY: kind
kind: # Create a kind cluster
	kind create cluster --config=config/kind/cluster.yaml
	kubectl apply -f config/kind/ingress-nginx.yaml

.PHONY: kind-install
kind-install: # Install the controller in the kind cluster
	@if [ ! -x $(LOCALBIN)/kind ]; then \
		echo "kind not installed. Installing..."; \
		if [ $$(uname) = "Linux" ]; then \
			[ $$(uname -m) = x86_64 ] && curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.22.0/kind-linux-amd64; \
			[ $$(uname -m) = aarch64 ] && curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.22.0/kind-linux-arm64; \
		elif [ $$(uname) = "Darwin" ]; then \
			[ $$(uname -m) = x86_64 ] && curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.22.0/kind-darwin-amd64; \
			[ $$(uname -m) = arm64 ] && curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.22.0/kind-darwin-arm64; \
		fi; \
		chmod +x ./kind; \
		sudo mv ./kind $(LOCALBIN)/kind; \
	fi

.PHONY: kind-delete
kind-delete: # Delete the kind cluster
	kind delete cluster --name $(CLUSTER_NAME)

.PHONY: default-admin
default-admin: # Create a default admin user
	kubectl apply -f config/samples/admin_rbac.yaml

.PHONY: default-redis-db
default-redis-db: # Create a default redis db
	kubectl apply -f config/samples/cloud_v1alpha1_storage.yaml
