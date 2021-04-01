run:
	go run src/main.go
lint:
	golangci-lint run
# Docker
dbuild:
	docker build -t lob . --target build 
drun:
	docker-compose up -d
dtest:
	docker build -t lob-test . --target base
	docker-compose -f docker-compose.test.yml up
ddeploy:
	docker build -t lob . --target build
	docker-compose -f docker-compose.prod.yml up -d