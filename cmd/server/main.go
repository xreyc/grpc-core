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
