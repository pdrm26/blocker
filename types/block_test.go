package types

import (
	"testing"

	"github.com/pdrm26/blocker/crypto"
	"github.com/pdrm26/blocker/utils"
	"github.com/stretchr/testify/assert"
)

func TestHashBlock(t *testing.T) {
	block := utils.RandomBlock()
	hashBlock := HashBlock(block)

	assert.Equal(t, len(hashBlock), 32)

}

func TestSignBlock(t *testing.T) {
	privKey := crypto.NewPrivateKey()
	pubKey := privKey.Public()

	block := utils.RandomBlock()
	sig := SignBlock(privKey, block)

	assert.Equal(t, len(sig.Bytes()), 64)
	assert.True(t, sig.Verify(pubKey, HashBlock(block)))
}
