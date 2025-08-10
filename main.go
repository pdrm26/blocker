package main

import (
	"context"
	"log"
	"time"

	"github.com/pdrm26/blocker/node"
	"github.com/pdrm26/blocker/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	node := node.NewNode()

	go func() {
		time.Sleep(2 * time.Second)
		makeTx()
	}()

	log.Fatal(node.Start(":3000"))
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
	client.ExchangeNodeInfo(context.TODO(), &proto.PeerInfo{ProtocolVersion: 1, BlockHeight: 10})
	_, err = client.HandleTransaction(context.TODO(), &proto.Transaction{})
	if err != nil {
		log.Fatal("HandleTransaction failed:", err)
	}

	log.Println("Transaction sent successfully")
}
