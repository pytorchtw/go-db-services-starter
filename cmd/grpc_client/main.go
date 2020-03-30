package main

import (
	"github.com/pytorchtw/go-db-services-starter/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
)

const (
	Address = "0.0.0.0:9090"
)

func main() {
	conn, err := grpc.Dial(Address, grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	c := proto.NewGrpcDBClient(conn)

	res, err := c.SayHello(context.Background(), &proto.SimpleRequest{Message: "Hello World"})
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(res.Message)
}
