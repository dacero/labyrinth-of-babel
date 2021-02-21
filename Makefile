run:
	go run src/main.go
build:
	go mod download
	go build -o bin/lob src/main.go
lint:
	golangci-lint run
# Docker
dbuild:
	docker build -t lob .
drun:
	docker-compose up -d
