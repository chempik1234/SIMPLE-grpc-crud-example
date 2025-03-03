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