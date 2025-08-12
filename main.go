package main

import (
	"context"
	"log"
	"time"

	"github.com/pdrm26/blocker/node"
	"github.com/pdrm26/blocker/proto"
)

func main() {
	makeNode(":3000", []string{})
	time.Sleep(time.Second)
	makeNode(":4000", []string{":3000"})
	time.Sleep(time.Second)
	makeNode(":5000", []string{":4000"})
	time.Sleep(time.Second)

	select {}
}

func makeNode(listenAddr string, bootstrapNodes []string) *node.Node {
	n := node.NewNode()
	go n.Start(listenAddr, bootstrapNodes)
	return n
}

func makeTx() {
	client, err := node.MakeNodeClient(":3000")
	if err != nil {
		log.Fatal(err)
	}
	client.Handshake(
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
