.PHONY: proto
proto:
	protoc -I proto/ \
		-I$(GOPATH)/src \
		-I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--go_out=plugins=grpc:proto/ \
		proto/sponsor.proto
	protoc -I proto/ \
		-I$(GOPATH)/src \
		-I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--grpc-gateway_out=logtostderr=true:proto/ \
		proto/sponsor.proto

.PHONY: build
build:
	@make server
	@make client

.PHONY: server
server:
	go build ./cmd/sponsor-server

.PHONY: client
client:
	go build ./cmd/sponsor-client

.PHONY: clean
clean:
	rm sponsor-server 2> /dev/null
	rm sponsor-client 2> /dev/null
