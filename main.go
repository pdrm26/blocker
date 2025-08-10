package main

import (
	"context"
	"log"
	"time"

	"github.com/pdrm26/blocker/node"
	"github.com/pdrm26/blocker/proto"
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
	client, err := node.MakeNodeClient(":3000")
	if err != nil {
		log.Fatal(err)
	}
	client.ExchangeNodeInfo(
		context.TODO(),
		//TODO: ListenAddr is incorrect i think because at the end it goes to targetAddr!!! check it later
		&proto.PeerInfo{ProtocolVersion: 1, BlockHeight: 10, ListenAddr: ":4000"},
	)
	_, err = client.HandleTransaction(context.TODO(), &proto.Transaction{})
	if err != nil {
		log.Fatal("HandleTransaction failed:", err)
	}

	log.Println("Transaction sent successfully")
}
