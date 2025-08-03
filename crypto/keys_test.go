package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPrivateKey(t *testing.T) {
	privKey := NewPrivateKey()
	assert.Equal(t, len(privKey.Bytes()), PrivateKeySize)

	pubKey := privKey.Public()
	assert.Equal(t, len(pubKey.Bytes()), PublicKeySize)
}

func TestPrivateKeySign(t *testing.T) {
	privKey := NewPrivateKey()
	pubKey := privKey.Public()

	msg := []byte("Hello")
	signature := privKey.Sign(msg)

	// test with valid message and pubkey
	assert.True(t, signature.Verify(pubKey, msg))

	// test with "in"valid message
	assert.False(t, signature.Verify(pubKey, []byte("Hi")))

	// test with "in"valid pubkey
	anotherPrivKey := NewPrivateKey()
	anotherPubKey := anotherPrivKey.Public()
	assert.False(t, signature.Verify(anotherPubKey, msg))
}
