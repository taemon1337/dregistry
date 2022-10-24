IMAGE := $(DOCKER_USER)/dregistry
VERSION := 0.0.1

tidy:
	go mod tidy

fmt:
	go fmt */**

build:
	docker build -t $(IMAGE):$(VERSION) .
