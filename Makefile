.PHONY: proto
proto:
	protoc \
      --proto_path=api/protobuf "api/protobuf/store.proto" \
      "--go_out=internal/common/genproto/store" --go_opt=paths=source_relative \
      --go-grpc_opt=require_unimplemented_servers=false \
      "--go-grpc_out=internal/common/genproto/store" --go-grpc_opt=paths=source_relative