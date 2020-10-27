default: build

build:
	env GO111MODULE=on GOGC=on go build -o app/eurofxref-ecb ./main.go

test:
	go test -v -cover -covermode=atomic -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

docker:
	docker build -t fahimbagar/eurofxref-ecb:1.0 .

.PHONY: build test docker