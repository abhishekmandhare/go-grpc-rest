
run:
	go run .\cmd\startup.go

docker-build:
	docker-compose build server

docker-run:
	docker-clean docker-build 
	docker-compose up server

vendor:
	go mod tidy
	go mod vendor