

build: fmt tidy
	go build

tidy:
	go mod tidy

fmt:
	go fmt */**
