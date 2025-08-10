package node

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/pdrm26/blocker/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Node struct {
	version int32
	proto.UnimplementedNodeServer
}

func NewNode() *Node {
	return &Node{
		version: 531,
	}
}

func (n *Node) Start(listenAddr string) error {
	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)

	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}

	proto.RegisterNodeServer(grpcServer, n)

	log.Println("gRPC server listening on", listenAddr)
	return grpcServer.Serve(ln)
}

func (n *Node) ExchangeNodeInfo(ctx context.Context, incomingPeerInfo *proto.PeerInfo) (*proto.PeerInfo, error) {
	localPeerInfo := &proto.PeerInfo{
		ProtocolVersion: n.version,
		BlockHeight:     1000,
	}

	remotePeer, ok := peer.FromContext(ctx)
	if !ok {
		fmt.Println("Peer not found in context")
	}

	fmt.Printf(
		"EXCHANGE NODE INFO :: incoming: %+v, local: %+v, remotePeer: %+v\n",
		incomingPeerInfo,
		localPeerInfo,
		remotePeer,
	)

	return localPeerInfo, nil

}

func (n *Node) HandleTransaction(ctx context.Context, tx *proto.Transaction) (*emptypb.Empty, error) {
	remotePeer, ok := peer.FromContext(ctx)
	if !ok {
		fmt.Println("Peer not found in context")
	}

	fmt.Printf("Received transaction from %+v :: incomingTx: %+v\n", remotePeer, tx)

	return &emptypb.Empty{}, nil
}
