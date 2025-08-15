package main

import (
	"context"
	"log"
	"time"

	"github.com/pdrm26/blocker/crypto"
	"github.com/pdrm26/blocker/node"
	"github.com/pdrm26/blocker/proto"
	"github.com/pdrm26/blocker/utils"
)

func main() {
	makeNode(":3000", []string{})
	time.Sleep(time.Second)
	makeNode(":4000", []string{":3000"})
	time.Sleep(time.Second)
	makeNode(":5000", []string{":4000"})

	for {
		time.Sleep(time.Second)
		makeTx()
	}
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

	privKey := crypto.NewPrivateKey()
	tx := &proto.Transaction{
		Version: 1,
		Inputs: []*proto.TxInput{
			{
				PrevTxHash:   utils.RandomHash(),
				PrevOutIndex: 0,
				PublicKey:    privKey.Public().Bytes(),
			},
		},
		Outputs: []*proto.TxOutput{
			{
				Amount:  99,
				Address: privKey.Public().Address().Bytes(),
			},
		},
	}
	_, err = client.HandleTransaction(context.TODO(), tx)
	if err != nil {
		log.Fatal("HandleTransaction failed:", err)
	}
}
