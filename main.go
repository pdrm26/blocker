package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/pdrm26/blocker/node"
	"github.com/pdrm26/blocker/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	go func() {
		time.Sleep(2 * time.Second)
		makeTx()
	}()

	log.Println("gRPC server listening on :3000")
	if err := grpcServer.Serve(ln); err != nil {
		log.Fatal(err)
	}
}

func makeTx() {
	conn, err := grpc.NewClient(
		":3000",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := proto.NewNodeClient(conn)
	_, err = client.HandleTransaction(context.TODO(), &proto.Transaction{})
	if err != nil {
		log.Fatal("HandleTransaction failed:", err)
	}

	log.Println("Transaction sent successfully")
}
