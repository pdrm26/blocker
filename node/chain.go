package node

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/pdrm26/blocker/crypto"
	"github.com/pdrm26/blocker/proto"
	"github.com/pdrm26/blocker/types"
)

const seed = "68c21e93b509d6de263c61b9754f9285fd8c3709e579f5baf4a83d874164c937"

type HeaderList struct {
	headers []*proto.Header
}
type UTXO struct {
	Hash     string
	OutIndex int
	Amount   int64
	Spent    bool
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
	txStore    TXStorer
	utxoStore  UTXOStorer
	blockStore BlockStorer
	headers    *HeaderList
}

func NewChain(blockStore BlockStorer, txStore TXStorer) *Chain {
	chain := &Chain{
		txStore:    txStore,
		utxoStore:  NewMemoryUTXOStore(),
		blockStore: blockStore,
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

	for _, tx := range block.Transactions {
		// for getting the hash of the genesis transaction and use it in the tests like: TestAddBlockWithTX
		// fmt.Println("NEW TX:", hex.EncodeToString(types.HashTransaction(tx)))
		if err := c.txStore.Put(tx); err != nil {
			return err
		}

		hash := hex.EncodeToString(types.HashTransaction(tx))
		for index, output := range tx.Outputs {
			utxo := &UTXO{
				Hash:     hash,
				Amount:   output.Amount,
				OutIndex: index,
				Spent:    false,
			}
			if err := c.utxoStore.Put(utxo); err != nil {
				return err
			}
		}
	}

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
	privKey := crypto.NewPrivateKeyFromString(seed)

	block := &proto.Block{
		Header: &proto.Header{
			Version: 1,
		},
	}
	tx := &proto.Transaction{
		Version: 1,
		Inputs:  []*proto.TxInput{},
		Outputs: []*proto.TxOutput{
			{
				Amount:  1000,
				Address: privKey.Public().Address().Bytes(),
			},
		},
	}

	block.Transactions = append(block.Transactions, tx)
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

	for _, tx := range b.Transactions {
		if !types.VerifyTransaction(tx) {
			return fmt.Errorf("invalid tx signature")
		}
	}

	return nil
}
