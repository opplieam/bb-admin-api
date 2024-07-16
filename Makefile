# Check to see if we can use ash, in Alpine images, or default to BASH.
SHELL_PATH = /bin/ash
SHELL = $(if $(wildcard $(SHELL_PATH)),/bin/ash,/bin/bash)

dev-up:
	minikube start
	eval $(minikube docker-env)

dev-down:
	minikube delete
