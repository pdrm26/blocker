package types

import (
	"crypto/sha256"
	"encoding/json"

	"github.com/pdrm26/blocker/crypto"
	"github.com/pdrm26/blocker/proto"
)

// HashBlock returns a SHA256 of the block header.
func HashBlock(block *proto.Block) []byte {
	return HashHeader(block.Header)
}

func HashHeader(header *proto.Header) []byte {
	b, err := json.Marshal(header)
	if err != nil {
		panic(err)
	}

	hash := sha256.Sum256(b)
	return hash[:]
}

func SignBlock(privKey *crypto.PrivateKey, block *proto.Block) *crypto.Signature {
	sig := privKey.Sign(HashBlock(block))
	block.PublicKey = privKey.Public().Bytes()
	block.Signature = sig.Bytes()
	return sig
}
