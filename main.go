package main

import (
	"backend/pkg"
	"backend/pkg/store"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

func main() {
	flag.Parse()

	// Create store.
	st := store.New()

	// Create a new gRPC server.
	s := pkg.New(st)
	// Run the gRPC server
	go func() {
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
		if err != nil {
			log.Printf("%v", err)
			os.Exit(1)
		}

		log.Println("gRPC server starting")
		if err := s.Server.Serve(listener); err != nil {
			log.Printf("%v", err)
			os.Exit(1)
		}
	}()

	// Gracefully shut down gRPC server after receiving an interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	slog.Info("Shutting down gRPC server")
	s.Server.GracefulStop()
}
