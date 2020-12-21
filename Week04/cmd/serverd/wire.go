//+build wireinject

package main

import (
	"github.com/google/wire"
	"xll.com/go-000/Week04/cmd/serverd/server"
	"xll.com/go-000/Week04/internal/data"
	"xll.com/go-000/Week04/internal/service"
)

// InitializeGRPCServer is builder of gprc server
func InitializeGRPCServer() *server.GRPCServer {
	wire.Build(server.NewServer, service.NewUserService, data.NewUserMemoryRepository)
	return &server.GRPCServer{}
}
