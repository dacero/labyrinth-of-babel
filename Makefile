run:
	go run src/hello.go
build:
	go build -o bin/hello src/hello.go
lint:
	golangci-lint run