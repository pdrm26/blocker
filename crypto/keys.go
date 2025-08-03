package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"io"
)

const (
	PrivateKeySize = 64
	PublicKeySize  = 32
	SeedSize       = 32
)

type PrivateKey struct {
	key ed25519.PrivateKey
}

func NewPrivateKey() *PrivateKey {
	seed := make([]byte, SeedSize)

	_, err := io.ReadFull(rand.Reader, seed)
	if err != nil {
		panic(err)
	}

	return &PrivateKey{
		key: ed25519.NewKeyFromSeed(seed),
	}
}

func (k *PrivateKey) Bytes() []byte {
	return k.key
}

func (k *PrivateKey) Sign(message []byte) *Signature {
	return &Signature{value: ed25519.Sign(k.key, message)}
}

func (k *PrivateKey) Public() *PublicKey {
	publicKey := make([]byte, PublicKeySize)
	copy(publicKey, k.key[32:])
	return &PublicKey{key: publicKey}
}

type PublicKey struct {
	key ed25519.PublicKey
}

func (k *PublicKey) Bytes() []byte {
	return k.key
}

func (k *PublicKey) Verify(rawMsg, sig []byte) bool {
	return ed25519.Verify(k.key, rawMsg, sig)
}

type Signature struct {
	value []byte
}

func (s *Signature) Bytes() []byte {
	return s.value
}

func (s *Signature) Verify(pubKey *PublicKey, msg []byte) bool {
	return ed25519.Verify(pubKey.key, msg, s.value)
}
