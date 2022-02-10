.PHONY: docker-make
docker-make proto lint test test_v test_short test_race test_stress test_reconnect test_all openapi_http:
	@docker-compose run --rm dev make $(CMD)

.PHONY: openapi_http
openapi_http: CMD="_openapi_http"

.PHONY: _openapi_http
_openapi_http:
	oapi-codegen -generate types -o "internal/store/ports/openapi_types.gen.go" -package "ports" "api/openapi/store.yml"
	oapi-codegen -generate chi-server -o "internal/store/ports/openapi_api.gen.go" -package "ports" "api/openapi/store.yml"
	oapi-codegen -generate types -o "internal/common/client/store/openapi_types.gen.go" -package "store" "api/openapi/store.yml"
	oapi-codegen -generate client -o "internal/common/client/store/openapi_client_gen.go" -package "store" "api/openapi/store.yml"

.PHONY: proto
proto: CMD="_proto"

.PHONY: _proto
_proto:
	protoc \
      --proto_path=api/protobuf "api/protobuf/store.proto" \
      "--go_out=internal/common/genproto/store" --go_opt=paths=source_relative \
      --go-grpc_opt=require_unimplemented_servers=false \
      "--go-grpc_out=internal/common/genproto/store" --go-grpc_opt=paths=source_relative

.PHONY: lint
lint: CMD="_lint"

.PHONY: _lint
_lint:
	go-cleanarch -infrastructure ports ./internal

.PHONY: test_all
test_all: CMD="_test_all"

.PHONY: _test_all
_test_all: _test _test_v _test_short _test_race _test_stress _test_reconnect

.PHONY: test
test: CMD="_test"

.PHONY: _test
_test:
	find . -name go.mod -execdir go test ./... \;

.PHONY: test_v
test_v: CMD="_test_v"

.PHONY: _test_v
_test_v:
	find . -name go.mod -execdir go test -v ./... \;

.PHONY: test_short
test_short: CMD="_test_short"

.PHONY: _test_short
_test_short:
	find . -name go.mod -execdir go test ./... -short \;

.PHONY: test_race
test_race: CMD="_test_race"

.PHONY: _test_race
_test_race:
	find . -name go.mod -execdir go test ./... -short -race \;

.PHONY: test_stress
test_stress: CMD="_test_stress"

.PHONY: _test_stress
_test_stress:
	find . -name go.mod -execdir go test -tags=stress -timeout=30m ./... \;

.PHONY: test_reconnect
test_reconnect: CMD="_test_reconnect"

.PHONY: _test_reconnect
_test_reconnect:
	find . -name go.mod -execdir go test -tags=reconnect ./... \;
