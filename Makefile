APP_NAME := grpc-core
MAIN := cmd/server/main.go
PROTO_DIR := contract
GEN_DIR := internal/gen/go
PROTO_FILES := $(PROTO_DIR)/auth/v1/user.proto

.PHONY: run build clean proto tidy

# Run the REST server
run:
	go run $(MAIN)

# Build the binary
build:
	go build -o bin/$(APP_NAME) $(MAIN)

# Clean up binaries and generated files
clean:
	rm -rf bin $(GEN_DIR)

# Generate gRPC and Protobuf code
proto:
	mkdir -p $(GEN_DIR)
	protoc \
		--go_out=$(GEN_DIR) \
		--go-grpc_out=$(GEN_DIR) \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		-I $(PROTO_DIR) \
		$(PROTO_FILES)

# Tidy dependencies
tidy:
	go mod tidy
