package node

import (
	"context"
	"encoding/hex"
	"net"
	"sync"
	"time"

	"github.com/pdrm26/blocker/crypto"
	"github.com/pdrm26/blocker/proto"
	"github.com/pdrm26/blocker/types"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/types/known/emptypb"
)

const blockTime = time.Second * 5

type Mempool struct {
	txx map[string]*proto.Transaction
}

func NewMempool() *Mempool {
	return &Mempool{
		txx: make(map[string]*proto.Transaction),
	}
}

func (pool *Mempool) Has(tx *proto.Transaction) bool {
	hash := hex.EncodeToString(types.HashTransaction(tx))
	_, ok := pool.txx[hash]
	return ok
}
func (pool *Mempool) Add(tx *proto.Transaction) bool {
	if pool.Has(tx) {
		return false
	}
	hash := hex.EncodeToString(types.HashTransaction(tx))
	pool.txx[hash] = tx
	return true
}

type ServerConfig struct {
	Version    int32
	ListenAddr string
	PrivKey    *crypto.PrivateKey
}

type Node struct {
	logger *zap.SugaredLogger

	peerLock sync.RWMutex
	peers    map[proto.NodeClient]*proto.PeerInfo
	mempool  *Mempool
	ServerConfig

	proto.UnimplementedNodeServer
}

func NewNode(serverConfig ServerConfig) *Node {
	logger, _ := zap.NewProduction()
	return &Node{
		peers:        make(map[proto.NodeClient]*proto.PeerInfo),
		logger:       logger.Sugar(),
		mempool:      NewMempool(),
		ServerConfig: serverConfig,
	}
}

func (n *Node) Start(listenAddr string, bootstrapNodes []string) error {
	n.ListenAddr = listenAddr
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

	if n.PrivKey != nil {
		go n.validatorLoop()
	}

	return grpcServer.Serve(ln)
}

func (n *Node) canConnectWith(addr string) bool {
	if n.ListenAddr == addr {
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
		n.logger.Debugw("dialing remote node", "we", n.ListenAddr, "remote", addr)

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
		ListenAddr:      n.ListenAddr,
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
		"we", n.ListenAddr,
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
	peer, ok := peer.FromContext(ctx)
	if !ok {
		panic("Peer not found in context")
	}

	if n.mempool.Add(tx) {
		hash := hex.EncodeToString(types.HashTransaction(tx))
		n.logger.Infow("received tx", "from", peer.Addr, "txHash", hash, "we", n.ListenAddr)
		go func() {
			if err := n.broadcast(tx); err != nil {
				n.logger.Errorw("broadcast error", "error", err)
			}
		}()
	}

	return &emptypb.Empty{}, nil
}

func MakeNodeClient(targetAddr string) (proto.NodeClient, error) {
	conn, err := grpc.NewClient(targetAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return proto.NewNodeClient(conn), nil
}

func (n *Node) broadcast(msg any) error {
	for peer := range n.peers {
		switch v := msg.(type) {
		case *proto.Transaction:
			_, err := peer.HandleTransaction(context.Background(), v)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (n *Node) validatorLoop() {
	n.logger.Infow("starting validator loop", "pubkey", n.PrivKey.Public(), "blockTime", blockTime)
	ticker := time.NewTicker(blockTime)

	for {
		<-ticker.C
		n.logger.Info("time to create a new block")
	}

}
