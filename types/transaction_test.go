package types

import (
	"testing"

	"github.com/pdrm26/blocker/crypto"
	"github.com/pdrm26/blocker/proto"
	"github.com/pdrm26/blocker/utils"
	"github.com/stretchr/testify/assert"
)

func TestTransaction(t *testing.T) {
	senderPrivKey := crypto.NewPrivateKey()
	receiverPrivKey := crypto.NewPrivateKey()

	input := &proto.TxInput{
		PrevTxHash:   utils.RandomHash(), // points to a previous transaction that created a UTXO
		PrevOutIndex: 0, // selects which output of that transaction (e.g. output 0)
		PublicKey:    senderPrivKey.Public().Bytes(), // Who owns the coin
	}

	output1 := &proto.TxOutput{
		Amount:  10,
		Address: receiverPrivKey.Public().Address().Bytes(),
	}
	output2 := &proto.TxOutput{
		Amount:  90,
		Address: senderPrivKey.Public().Address().Bytes(),
	}

	tx := &proto.Transaction{
		Version: 1,
		Inputs:  []*proto.TxInput{input},
		Outputs: []*proto.TxOutput{output1, output2},
	}

	sign := SignTransaction(tx, senderPrivKey)

	input.Signature = sign.Bytes()

	assert.True(t, VerifyTransaction(tx))
}
