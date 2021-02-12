run:
	go run src/hello.go
build:
	go mod download
	go build -o bin/hello src/hello.go
lint:
	golangci-lint run