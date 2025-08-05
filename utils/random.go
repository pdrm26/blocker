package utils

import (
	"crypto/rand"
	intrand "math/rand"
	"time"

	"github.com/pdrm26/blocker/proto"
)

func RandomHash() []byte {
	hash := make([]byte, 32)
	rand.Read(hash)
	return hash

}

func RandomBlock() *proto.Block {
	header := &proto.Header{
		Version:   1,
		Height:    int32(intrand.Intn(1000)),
		PrevHash:  RandomHash(),
		RootHash:  RandomHash(),
		Timestamp: time.Now().Unix(),
	}

	return &proto.Block{Header: header}

}
