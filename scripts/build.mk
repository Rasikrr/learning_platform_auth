PKG=./pkg
SWAG_PATH=internal/ports/http/server.go

update_core:
	go get github.com/Rasikrr/learning_platform_core@latest

coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go fmt ./...
	golangci-lint run

tests:
	go test ./... -v


build:
	mkdir -p ./bin
	go build -o ./bin ./cmd/app/main.go

gen_docs:
	swag init -g internal/ports/http/server.go -d .

proto:
	@echo "Generating proto grpc files..."
	@rm -rf "${PKG}/api/grpc"
	@mkdir -p "${PKG}"
	@protoc \
		-I="./" \
		--go_out="${PKG}" \
		--go-grpc_out="${PKG}" \
		--go_opt="paths=source_relative" \
		--go-grpc_opt="paths=source_relative" \
		./api/proto/*/*/*.proto
	@echo "Files generated"

proto_deps:
	@echo "Generating proto deps..."
	@rm -rf "${PKG}/deps"
	@mkdir -p "${PKG}/deps"
	@protoc \
		-I="./deps" \
		--go_out="${PKG}/deps" \
		--go-grpc_out="${PKG}/deps" \
		--go_opt="paths=source_relative" \
		--go-grpc_opt="paths=source_relative" \
		./deps/api/proto/*/*.proto
	@echo "Files generated"