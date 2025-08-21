package node

import (
	"testing"

	"github.com/pdrm26/blocker/types"
	"github.com/pdrm26/blocker/utils"
	"github.com/stretchr/testify/assert"
)

func TestAddBlock(t *testing.T) {
	chain := NewChain(NewMemoryBlockStore())

	for i := 0; i < 100; i++ {
		block := utils.RandomBlock()
		blockHash := types.HashBlock(block)

		assert.Nil(t, chain.AddBlock(block))

		fetchedBlockByHash, err := chain.GetBlockByHash(blockHash)
		assert.Nil(t, err)
		assert.Equal(t, fetchedBlockByHash, block)

		fetchedBlockByHeight, err := chain.GetBlockByHeight(i)
		assert.Nil(t, err)
		assert.Equal(t, fetchedBlockByHeight, block)
	}

}

func TestChainHeight(t *testing.T) {
	chain := NewChain(NewMemoryBlockStore())

	for i := 0; i < 100; i++ {
		block := utils.RandomBlock()
		prevBlock, err := chain.GetBlockByHeight(chain.Height())
		assert.Nil(t, err)
		block.Header.PrevHash = types.HashBlock(prevBlock)
		assert.Nil(t, chain.AddBlock(block))
		assert.Equal(t, chain.Height(), i + 1)
	}
}

func TestNewChain(t *testing.T) {
	chain := NewChain(NewMemoryBlockStore())
	assert.Equal(t, 0, chain.Height())
	_, err := chain.GetBlockByHeight(0)
	assert.Nil(t, err)
}
