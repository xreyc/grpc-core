## Implementation steps

#### Clone repository
```bash
git clone git@github.com:xreyc/grpc-core.git
cd grpc-core
```

#### Initialize module
```bash
go mod init github.com/xreyc/grpc-core
```

#### Add grpc-contract as submodule
```bash
git submodule add https://github.com/xreyc/grpc-contract.git contract
git submodule update --init --recursive
```

These will generate

```
grpc-core/contract/auth/v1/user.proto
```

#### Install require tools
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

#### Generate code for proto
```bash
mkdir -p internal/gen/go
```

Then run

```bash
protoc \
  --go_out=internal/gen/go \
  --go-grpc_out=internal/gen/go \
  --go_opt=paths=source_relative \
  --go-grpc_opt=paths=source_relative \
  -I contract \
  contract/auth/v1/user.proto
```

This will generate
```
internal/gen/go/auth/v1/user.pb.go
internal/gen/go/auth/v1/user_grpc.pb.go
```

#### Install dependencies
```bash
go get google.golang.org/grpc
go get github.com/gin-gonic/gin
```

#### Project structure
```
grpc-core/
├── cmd/
│   └── server/
│       └── main.go
├── contract/                   # submodule
├── internal/
│   ├── app/
│   ├── config/
│   ├── domain/
│   ├── gen/go/                # generated gRPC files
│   ├── handler/
│   │   ├── rest/              # Gin REST handlers
│   │   └── grpc/              # gRPC client
│   ├── route/                 # Gin router setup
│   ├── service/               # business logic layer
│   └── grpc/                  # grpc client setup
├── go.mod
├── go.sum
└── Makefile
```

#### Implement grpc client
`internal/grpc/client.go`
```go
package grpc

import (
    "log"

    authv1 "github.com/xreyc/grpc-core/internal/gen/go/auth/v1"
    "google.golang.org/grpc"
)

var AuthClient authv1.UserServiceClient

func InitGRPCClients() {
    conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("Failed to connect to auth service: %v", err)
    }

    AuthClient = authv1.NewUserServiceClient(conn)
}
```

#### Implement rest handler
`internal/handler/rest/user_handler.go`
```go
package rest

import (
    "net/http"

    "github.com/gin-gonic/gin"
    grpcClient "github.com/xreyc/grpc-core/internal/grpc"
    authv1 "github.com/xreyc/grpc-core/internal/gen/go/auth/v1"
)

func GetUserDetails(c *gin.Context) {
    username := c.Query("username")
    if username == "" {
        username = "xreyc" // default for testing
    }

    resp, err := grpcClient.AuthClient.GetUserDetails(c, &authv1.GetUserRequest{
        Username: username,
    })
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "username":  resp.GetUsername(),
        "email":     resp.GetEmail(),
        "full_name": resp.GetFullName(),
    })
}
```

#### Implement router
`internal/route/router.go`
```go
package route

import (
    "github.com/gin-gonic/gin"
    "github.com/xreyc/grpc-core/internal/handler/rest"
)

func SetupRouter() *gin.Engine {
    r := gin.Default()
    r.GET("/get-user-details", rest.GetUserDetails)
    return r
}
```

#### Implement entry point
`cmd/server/main.go`
```go
package main

import (
    "github.com/xreyc/grpc-core/internal/grpc"
    "github.com/xreyc/grpc-core/internal/route"
)

func main() {
    grpc.InitGRPCClients()

    r := route.SetupRouter()
    r.Run(":8080") // REST server
}
```

#### Run server
```bash
go run cmd/server/main.go
```

#### Create a Makefile
```makefile
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
```

Usage
```bash
make run       # Run the gRPC server
make build     # Build binary to ./bin/grpc-auth
make proto     # Regenerate .pb.go files
make tidy      # Clean up go.mod/go.sum
make clean     # Remove ./bin folder
```