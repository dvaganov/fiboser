.PHONY:build
build:
	go build -o bin/fiboser ./cmd/main/

.PHONY:run
run:
	go run cmd/main/main.go

.PHONY:test
test:
	go test -v ./pkg/...