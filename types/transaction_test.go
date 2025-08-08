package types

import (
	"testing"

	"github.com/pdrm26/blocker/crypto"
	"github.com/pdrm26/blocker/proto"
	"github.com/pdrm26/blocker/utils"
	"github.com/stretchr/testify/assert"
)

func TestTransaction(t *testing.T) {
	fromPrivKey := crypto.NewPrivateKey()
	toPrivKey := crypto.NewPrivateKey()

	input := &proto.TxInput{
		PrevTxHash:   utils.RandomHash(),
		PrevOutIndex: 0,
		PublicKey:    fromPrivKey.Public().Bytes(),
	}

	output1 := &proto.TxOutput{
		Amount:  10,
		Address: toPrivKey.Public().Address().Bytes(),
	}
	output2 := &proto.TxOutput{
		Amount:  90,
		Address: fromPrivKey.Public().Address().Bytes(),
	}

	tx := &proto.Transaction{
		Version: 1,
		Inputs:  []*proto.TxInput{input},
		Outputs: []*proto.TxOutput{output1, output2},
	}

	sign := SignTransaction(tx, fromPrivKey)

	input.Signature = sign.Bytes()

	assert.True(t, VerifyTransaction(tx))
}
