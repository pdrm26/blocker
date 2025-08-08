package node

import (
	"context"

	"github.com/pdrm26/blocker/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Node struct {
	proto.UnimplementedNodeServer
}

func NewNode() *Node {
	return &Node{}
}

func (n *Node) HandleTransaction(ctx context.Context, tx *proto.Transaction) (*emptypb.Empty, error) {
	return nil, nil
}
