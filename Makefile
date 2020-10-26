default: build

build: sync
	env GO111MODULE=on GOGC=off go build -mod=vendor -v ./cmd/main.go

sync:
	go mod tidy
	go mod vendor
