package main

import (
	"github.com/pytorchtw/go-db-services-starter/handlers"
	"github.com/pytorchtw/go-db-services-starter/proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	Address = "0.0.0.0:9090"
)

func main() {
	listen, err := net.Listen("tcp", Address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	proto.RegisterGrpcDBServer(s, &handlers.GrpcHandler{})

	log.Println("Listen on " + Address)

	if err := s.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
