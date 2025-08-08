package main

import (
	"log"
	"net"

	"github.com/pdrm26/blocker/node"
	"github.com/pdrm26/blocker/proto"
	"google.golang.org/grpc"
)

func main() {
	node := node.NewNode()

	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)

	ln, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatal(err)
	}

	proto.RegisterNodeServer(grpcServer, node)
	grpcServer.Serve(ln)
}
