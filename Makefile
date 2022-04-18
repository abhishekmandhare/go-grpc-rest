
run:
	go run .\cmd\main.go

docker-build:
	docker-compose build server

docker-run:
	docker-clean docker-build 
	docker-compose up server

.PHONY: vendor
vendor:
	go mod tidy
	go mod vendor

.PHONY: gen
gen:
	buf mod update
	buf generate