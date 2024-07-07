package pkg

import (
	"backend/explore"
	"backend/pkg/explore_server"
	"backend/pkg/store"
	"context"
	"log"

	"google.golang.org/grpc"
)

type ServerImpl struct {
	*grpc.Server
	store.Store
}

func New(st store.Store) *ServerImpl {
	// Create a new gRPC server.
	s := newServer()
	explore.RegisterExploreServiceServer(s, &explore_server.ExploreServerImpl{Store: st})
	return &ServerImpl{Server: s, Store: st}
}

func newServer() *grpc.Server {
	s := grpc.NewServer(grpc.UnaryInterceptor(loggingInterceptor))
	return s
}

// gRPC logs.
func loggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("Received request: %v", req)
	resp, err := handler(ctx, req)
	return resp, err
}
