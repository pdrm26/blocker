package types

import (
	"crypto/sha256"

	"github.com/pdrm26/blocker/crypto"
	"github.com/pdrm26/blocker/proto"
	pb "google.golang.org/protobuf/proto"
)

func HashTransaction(tx *proto.Transaction) []byte {
	b, err := pb.Marshal(tx)
	if err != nil {
		panic(err)
	}

	hash := sha256.Sum256(b)
	return hash[:]
}

func SignTransaction(tx *proto.Transaction, privKey *crypto.PrivateKey) *crypto.Signature {
	return privKey.Sign(HashTransaction(tx))
}

func VerifyTransaction(tx *proto.Transaction) bool {
	for _, input := range tx.Inputs {
		if len(input.Signature) == 0 {
			panic("transaction has no signature")
		}

		pubKey := crypto.PublicKeyFromBytes(input.PublicKey)
		sig := crypto.SignatureFromBytes(input.Signature)

		input.Signature = nil
		if !sig.Verify(pubKey, HashTransaction(tx)) {
			return false
		}
	}

	return true
}
