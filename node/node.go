package node

import (
	"context"
	"fmt"

	"github.com/pdrm26/blocker/proto"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Node struct {
	proto.UnimplementedNodeServer
}

func NewNode() *Node {
	return &Node{}
}

func (n *Node) HandleTransaction(ctx context.Context, tx *proto.Transaction) (*emptypb.Empty, error) {
	peer, ok := peer.FromContext(ctx)
	if !ok {
		fmt.Println("NOT OK")
	}
	fmt.Println("Received transaction ::", peer)
	return nil, nil
}
