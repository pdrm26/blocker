package node

import (
	"encoding/hex"

	"github.com/pdrm26/blocker/proto"
)


type Chain struct {
	blockStore BlockStorer
}

func NewChain(bs BlockStorer) *Chain {
	return &Chain{
		blockStore: bs,
	}
}

func (c *Chain) AddBlock(block *proto.Block) error {
	return c.blockStore.Put(block)
}

func (c *Chain) GetBlockByHash(hash []byte) (*proto.Block, error) {
	hashHex := hex.EncodeToString(hash)
	return c.blockStore.Get(hashHex)
}

func (c *Chain) GetBlockByHeight(height int) (*proto.Block, error) {
	return nil, nil
}
