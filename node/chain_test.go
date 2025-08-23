package node

import (
	"encoding/hex"
	"testing"

	"github.com/pdrm26/blocker/crypto"
	"github.com/pdrm26/blocker/proto"
	"github.com/pdrm26/blocker/types"
	"github.com/pdrm26/blocker/utils"
	"github.com/stretchr/testify/assert"
)

func randomBlock(t *testing.T, chain *Chain) *proto.Block {
	privKey := crypto.NewPrivateKey()
	block := utils.RandomBlock()
	prevBlock, err := chain.GetBlockByHeight(chain.Height())
	assert.Nil(t, err)
	block.Header.PrevHash = types.HashBlock(prevBlock)
	types.SignBlock(privKey, block)

	return block
}

func TestAddBlock(t *testing.T) {
	chain := NewChain(NewMemoryBlockStore(), NewMemoryTXStore())

	for i := 0; i < 100; i++ {
		block := randomBlock(t, chain)
		blockHash := types.HashBlock(block)

		assert.Nil(t, chain.AddBlock(block))

		fetchedBlockByHash, err := chain.GetBlockByHash(blockHash)
		assert.Nil(t, err)
		assert.Equal(t, fetchedBlockByHash, block)

		fetchedBlockByHeight, err := chain.GetBlockByHeight(i + 1)
		assert.Nil(t, err)
		assert.Equal(t, fetchedBlockByHeight, block)
	}

}

func TestChainHeight(t *testing.T) {
	chain := NewChain(NewMemoryBlockStore(), NewMemoryTXStore())

	for i := 0; i < 100; i++ {
		block := randomBlock(t, chain)
		assert.Nil(t, chain.AddBlock(block))
		assert.Equal(t, chain.Height(), i+1)
	}
}

func TestNewChain(t *testing.T) {
	chain := NewChain(NewMemoryBlockStore(), NewMemoryTXStore())
	assert.Equal(t, 0, chain.Height())
	_, err := chain.GetBlockByHeight(0)
	assert.Nil(t, err)
}

func TestAddBlockWithTX(t *testing.T) {
	var (
		chain     = NewChain(NewMemoryBlockStore(), NewMemoryTXStore())
		block     = randomBlock(t, chain)
		privKey   = crypto.NewPrivateKey()
		recipient = crypto.NewPrivateKey().Public().Address()
	)

	inputs := []*proto.TxInput{
		{
			PublicKey: privKey.Public().Bytes(),
		},
	}
	outputs := []*proto.TxOutput{
		{
			Amount:  99,
			Address: recipient.Bytes(),
		},
	}

	tx := &proto.Transaction{Version: 1, Inputs: inputs, Outputs: outputs}
	block.Transactions = append(block.Transactions, tx)

	assert.Nil(t, chain.AddBlock(block))

	txHash := hex.EncodeToString(types.HashTransaction(tx))
	fetchedTx, err := chain.txStore.Get(txHash)

	assert.Equal(t, tx, fetchedTx)
	assert.Nil(t, err)

}
