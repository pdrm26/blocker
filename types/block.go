package types

import (
	"crypto/sha256"

	"github.com/pdrm26/blocker/crypto"
	"github.com/pdrm26/blocker/proto"
	pb "google.golang.org/protobuf/proto"
)

// HashBlock returns a SHA256 of the block header.
func HashBlock(block *proto.Block) []byte {
	b, err := pb.Marshal(block)
	if err != nil {
		panic(err)
	}

	hash := sha256.Sum256(b)
	return hash[:]
}

func SignBlock(privKey *crypto.PrivateKey, block *proto.Block) *crypto.Signature {
	return privKey.Sign(HashBlock(block))
}
