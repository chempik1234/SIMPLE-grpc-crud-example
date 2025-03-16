generate:
	protoc -I . \
		--go_out=. \
		--go-grpc_out=. \
		--grpc-gateway_out . \
		--grpc-gateway_opt generate_unbound_methods=true \
		.\api\order.proto

build: generate
	go build .\cmd\main.go

run: build
	./main.exe

build_orders:
	docker build -t orders_service:latest -f ./build/backend/Dockerfile ./

build_tests:
	docker build -t tests:latest -f ./build/integrationTests/Dockerfile ./

test:
	docker rm -f test-container
	docker compose -f ./deployments/integrationTests/docker-compose.yml down
	docker compose -f ./deployments/integrationTests/docker-compose.yml up -d
	docker run --network=integrationtests_deployments_service_network --name test-container tests:latest
	docker compose -f ./deployments/integrationTests/docker-compose.yml down
	docker rm -f test-container
