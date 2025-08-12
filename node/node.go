package node

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/pdrm26/blocker/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Node struct {
	version    int32
	listenAddr string
	logger     *zap.SugaredLogger

	peerLock sync.RWMutex
	peers    map[proto.NodeClient]*proto.PeerInfo

	proto.UnimplementedNodeServer
}

func NewNode() *Node {
	logger, _ := zap.NewProduction()
	return &Node{
		peers:   make(map[proto.NodeClient]*proto.PeerInfo),
		version: 531,
		logger:  logger.Sugar(),
	}
}

func (n *Node) Start(listenAddr string, bootstrapNodes []string) error {
	n.listenAddr = listenAddr
	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)

	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}

	proto.RegisterNodeServer(grpcServer, n)

	n.logger.Info("node running on port", listenAddr)
	if len(bootstrapNodes) > 0 {
		go n.bootstrapNetwork(bootstrapNodes)
	}
	return grpcServer.Serve(ln)
}

func (n *Node) canConnectWith(addr string) bool {
	if n.listenAddr == addr {
		return false
	}

	for _, address := range n.getPeerList() {
		if addr == address {
			return false
		}
	}

	return true
}

func (n *Node) bootstrapNetwork(addrs []string) error {
	for _, addr := range addrs {
		if !n.canConnectWith(addr) {
			continue
		}
		client, peer, err := n.dialRemoteNode(addr)
		if err != nil {
			return err
		}
		n.addPeer(client, peer)
	}

	return nil
}

func (n *Node) dialRemoteNode(addr string) (proto.NodeClient, *proto.PeerInfo, error) {
	client, err := MakeNodeClient(addr)
	if err != nil {
		return nil, nil, err
	}
	peer, err := client.Handshake(context.Background(), n.getPeerInfo())
	if err != nil {
		return nil, nil, err
	}

	return client, peer, nil

}

func (n *Node) getPeerInfo() *proto.PeerInfo {
	return &proto.PeerInfo{
		ProtocolVersion: 1,
		BlockHeight:     0,
		ListenAddr:      n.listenAddr,
		PeerList:        n.getPeerList(),
	}
}

func (n *Node) getPeerList() []string {
	n.peerLock.RLock()
	defer n.peerLock.RUnlock()

	peers := []string{}
	for _, peer := range n.peers {
		peers = append(peers, peer.ListenAddr)
	}

	return peers
}

func (n *Node) addPeer(p proto.NodeClient, peerInfo *proto.PeerInfo) {
	n.peerLock.Lock()
	defer n.peerLock.Unlock()

	n.peers[p] = peerInfo

	if len(peerInfo.PeerList) > 0 {
		go n.bootstrapNetwork(peerInfo.PeerList)
	}

	n.logger.Debugw(
		"new peer successfully connected",
		"we", n.listenAddr,
		"remoteNode", peerInfo.ListenAddr,
		"height", peerInfo.BlockHeight,
	)
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

	n.logger.Infof("Received tx from %+v :: incomingTx: %+v\n", remotePeer, tx)
	return &emptypb.Empty{}, nil
}

func MakeNodeClient(targetAddr string) (proto.NodeClient, error) {
	conn, err := grpc.NewClient(targetAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return proto.NewNodeClient(conn), nil
}
