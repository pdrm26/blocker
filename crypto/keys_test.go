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

	assert.True(t, signature.Verify(pubKey, msg))
	assert.False(t, signature.Verify(pubKey, []byte("Hi")))
}
