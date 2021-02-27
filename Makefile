run:
	go run src/main.go
build:
	go mod download
	go build -o bin/lob
lint:
	golangci-lint run
# Docker
dbuild:
	docker build -t lob . --target build 
drun:
	docker-compose up -d
dtest:
	docker-compose -f docker-compose.test.yml up
