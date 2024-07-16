# Check to see if we can use ash, in Alpine images, or default to BASH.
SHELL_PATH = /bin/ash
SHELL = $(if $(wildcard $(SHELL_PATH)),/bin/ash,/bin/bash)

dev-up:
	minikube start
	eval $(minikube docker-env)

dev-down:
	minikube delete


BASE_IMAGE_NAME := opplieam
SERVICE_NAME    := bb-admin-api
#VERSION         := "0.0.1-$(shell git rev-parse --short HEAD)"
VERSION         := "local-dev"
SERVICE_IMAGE   := $(BASE_IMAGE_NAME)/$(SERVICE_NAME):$(VERSION)

docker-build:
	docker build \
		-t $(SERVICE_IMAGE) \
    	--build-arg BUILD_REF=$(VERSION) \
    	--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
    	.
