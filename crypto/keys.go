package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"io"
)

const (
	PrivateKeySize = 64
	SignatureLen   = 64
	PublicKeySize  = 32
	SeedLen        = 32
	AddressSize    = 20
)

type PrivateKey struct {
	key ed25519.PrivateKey
}

func NewPrivateKeyFromString(s string) *PrivateKey {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}

	return NewPrivateKeyFromSeed(b)
}

func NewPrivateKeyFromSeed(seed []byte) *PrivateKey {
	if len(seed) != SeedLen {
		panic("invalid seed length, must be 32")
	}

	return &PrivateKey{
		key: ed25519.NewKeyFromSeed(seed),
	}
}

func NewPrivateKey() *PrivateKey {
	seed := make([]byte, SeedLen)

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

func PublicKeyFromBytes(pubKeyBytes []byte) *PublicKey {
	if len(pubKeyBytes) != PublicKeySize {
		panic("length of the bytes not equal to 32")
	}
	return &PublicKey{
		key: ed25519.PublicKey(pubKeyBytes),
	}
}

func (k *PublicKey) Bytes() []byte {
	return k.key
}

func (k *PublicKey) Verify(rawMsg, sig []byte) bool {
	return ed25519.Verify(k.key, rawMsg, sig)
}

func (k *PublicKey) Address() Address {
	return Address{
		value: k.key[PublicKeySize-AddressSize:],
	}
}

type Signature struct {
	value []byte
}

func SignatureFromBytes(sigByte []byte) *Signature {
	if len(sigByte) != SignatureLen {
		panic("length of the bytes not equal to 64")
	}
	return &Signature{
		value: sigByte,
	}

}

func (s *Signature) Bytes() []byte {
	return s.value
}

func (s *Signature) Verify(pubKey *PublicKey, msg []byte) bool {
	return ed25519.Verify(pubKey.key, msg, s.value)
}

type Address struct {
	value []byte
}

func (s Address) Bytes() []byte {
	return s.value
}

func (s Address) String() string {
	return hex.EncodeToString(s.value)
}
