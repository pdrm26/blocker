package types

import (
	"testing"

	"github.com/pdrm26/blocker/utils"
	"github.com/stretchr/testify/assert"
)

func TestHashBlock(t *testing.T) {
	block := utils.RandomBlock()
	hashBlock := HashBlock(block)

	assert.Equal(t, len(hashBlock), 32)

}
