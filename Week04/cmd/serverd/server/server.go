package server

import (
	"context"
	"log"
	"net"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"xll.com/go-000/Week04/api"
)

// GRPCServer is grpc server object
type GRPCServer struct {
	UserService api.UserServer
}

// NewServer is the constructor function
func NewServer(s api.UserServer) *GRPCServer {
	return &GRPCServer{
		UserService: s,
	}
}

// Start is the function to start the GRPC server
func (srv *GRPCServer) Start(ctx context.Context, g *errgroup.Group) error {
	lis, err := net.Listen("tcp", "localhost:9000")
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	g.Go(func() error {
		go func() {
			<-ctx.Done()
			s.GracefulStop()
			log.Printf("Shutdown Server")
		}()
		api.RegisterUserServer(s, srv.UserService)
		log.Printf("Server started at localhost:9000")
		return s.Serve(lis)
	})
	return nil
}
