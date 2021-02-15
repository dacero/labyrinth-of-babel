run:
	go run src/hello.go
build:
	go mod download
	go build -o bin/lob src/lob.go
lint:
	golangci-lint run
# Docker
dbuild:
	docker build -t lob .
drun:
	docker-compose up -d
