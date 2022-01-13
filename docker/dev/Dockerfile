FROM golang:1.17 AS go
# TODO: add openapitools/openapi-generator-cli

FROM go AS protoc
RUN VERSION=3.19.3 && \
    apt-get update && \
    apt-get install unzip && \
    wget -O x.zip https://github.com/protocolbuffers/protobuf/releases/download/v$VERSION/protoc-$VERSION-linux-x86_64.zip && \
    unzip x.zip -d protoc

FROM go AS protobuf
RUN go install github.com/golang/protobuf/protoc-gen-go@v1.5.2
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0

FROM go AS oapi-codegen
# TODO: update to v1.5.1
RUN go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.3.9

FROM go AS cleanarch
RUN go install github.com/roblaszczak/go-cleanarch@v1.2.1

FROM go
COPY --from=protoc       /go/protoc/                /usr/
COPY --from=protobuf     /go/bin/protoc-gen-go      /usr/bin/
COPY --from=protobuf     /go/bin/protoc-gen-go-grpc /usr/bin/
COPY --from=oapi-codegen /go/bin/oapi-codegen       /usr/bin/
COPY --from=cleanarch    /go/bin/go-cleanarch       /usr/bin/