package node

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/pdrm26/blocker/crypto"
	"github.com/pdrm26/blocker/proto"
	"github.com/pdrm26/blocker/types"
)

type HeaderList struct {
	headers []*proto.Header
}

func NewHeaderList() *HeaderList {
	return &HeaderList{
		headers: []*proto.Header{},
	}
}

func (h *HeaderList) Add(header *proto.Header) {
	h.headers = append(h.headers, header)
}

func (h *HeaderList) Get(height int) *proto.Header {
	if height > h.Height() {
		panic("height is too high!")
	}
	return h.headers[height]
}

func (h *HeaderList) Len() int {
	return len(h.headers)
}

func (h *HeaderList) Height() int {
	return h.Len() - 1
}

type Chain struct {
	blockStore BlockStorer
	headers    *HeaderList
}

func NewChain(bs BlockStorer) *Chain {
	chain := &Chain{
		blockStore: bs,
		headers:    NewHeaderList(),
	}
	chain.addBlock(chain.createGenesisBlock())
	return chain
}

func (c *Chain) Height() int {
	return c.headers.Height()
}

func (c *Chain) AddBlock(block *proto.Block) error {
	if err := c.ValidateBlock(block); err != nil {
		return err
	}
	return c.addBlock(block)
}

func (c *Chain) addBlock(block *proto.Block) error {
	c.headers.Add(block.Header)
	return c.blockStore.Put(block)
}

func (c *Chain) GetBlockByHash(hash []byte) (*proto.Block, error) {
	hashHex := hex.EncodeToString(hash)
	return c.blockStore.Get(hashHex)
}

func (c *Chain) GetBlockByHeight(height int) (*proto.Block, error) {
	if height > c.Height() {
		return nil, fmt.Errorf("given height (%d) too heigh - height (%d)", height, c.Height())
	}

	header := c.headers.Get(height)
	headerHash := types.HashHeader(header)
	return c.GetBlockByHash(headerHash)
}

func (c *Chain) createGenesisBlock() *proto.Block {
	privKey := crypto.NewPrivateKey()
	block := &proto.Block{
		Header: &proto.Header{
			Version: 1,
		},
	}
	types.SignBlock(privKey, block)

	return block
}

func (c *Chain) ValidateBlock(b *proto.Block) error {
	if !types.VerifyBlock(b) {
		return fmt.Errorf("invalid block signature")
	}

	currentBlock, err := c.GetBlockByHeight(c.Height())
	if err != nil {
		return err
	}

	hash := types.HashBlock(currentBlock)
	if !bytes.Equal(hash, b.Header.PrevHash) {
		return fmt.Errorf("invalid previous hash block")
	}
	return nil
}
