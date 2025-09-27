.PHONY: proto
proto:
	protoc \
		--go_out=./gen \
		--go_opt=paths=source_relative \
		--go-grpc_out=./gen \
		--go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=./gen \
		--grpc-gateway_opt=paths=source_relative \
		--openapiv2_out=./gen \
		-I ./proto \
		-I ./third_party \
		proto/ateway.proto