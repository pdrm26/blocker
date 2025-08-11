package node

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/pdrm26/blocker/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Node struct {
	version    int32
	listenAddr string

	peerLock sync.RWMutex
	peers    map[proto.NodeClient]*proto.PeerInfo

	proto.UnimplementedNodeServer
}

func NewNode() *Node {
	return &Node{
		peers:   make(map[proto.NodeClient]*proto.PeerInfo),
		version: 531,
	}
}

func (n *Node) Start(listenAddr string) error {
	n.listenAddr = listenAddr
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

func (n *Node) BootstrapNetwork(addrs []string) error {
	for _, addr := range addrs {
		client, err := MakeNodeClient(addr)
		if err != nil {
			return err
		}

		peerInfo, err := client.Handshake(context.Background(), n.getPeerInfo())
		if err != nil {
			fmt.Println("handshake error", err)
			continue
		}

		n.addPeer(client, peerInfo)
	}

	return nil
}

func (n *Node) getPeerInfo() *proto.PeerInfo {
	return &proto.PeerInfo{
		ProtocolVersion: 1,
		BlockHeight:     0,
		ListenAddr:      n.listenAddr,
	}
}

func (n *Node) addPeer(p proto.NodeClient, peerInfo *proto.PeerInfo) {
	n.peerLock.Lock()
	defer n.peerLock.Unlock()

	fmt.Printf("[%s] new peer connected (%s) - height (%d)\n", n.listenAddr, peerInfo.ListenAddr, peerInfo.BlockHeight)

	n.peers[p] = peerInfo
}

func (n *Node) removePeer(p proto.NodeClient) {
	n.peerLock.Lock()
	defer n.peerLock.Unlock()

	delete(n.peers, p)
}

func (n *Node) Handshake(ctx context.Context, incomingPeerInfo *proto.PeerInfo) (*proto.PeerInfo, error) {
	client, err := MakeNodeClient(incomingPeerInfo.ListenAddr)
	if err != nil {
		return nil, err
	}

	n.addPeer(client, incomingPeerInfo)

	return n.getPeerInfo(), nil

}

func (n *Node) HandleTransaction(ctx context.Context, tx *proto.Transaction) (*emptypb.Empty, error) {
	remotePeer, ok := peer.FromContext(ctx)
	if !ok {
		fmt.Println("Peer not found in context")
	}

	fmt.Printf("Received transaction from %+v :: incomingTx: %+v\n", remotePeer, tx)

	return &emptypb.Empty{}, nil
}

func MakeNodeClient(targetAddr string) (proto.NodeClient, error) {
	conn, err := grpc.NewClient(targetAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return proto.NewNodeClient(conn), nil
}
