//
// This is just demo code that happens to compile and run.
//
// This software is just plain non-sense, non-idomatic, bugged,
// unpolished and unfinished.
// It's just my studying notes hacked in a pkg!
//
// This stuff is intended to be used to explain things while looking
// decent and simple on the slides.
//
package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
)

import (
	"crypto/rand"
	sha "crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/big"
)

type PrivateKey interface {
	Sign(hash []byte) Signature
	Public() PublicKey
}

type privateKey struct {
	*ecdsa.PrivateKey
}

var encoder = base64.StdEncoding.WithPadding(base64.NoPadding)

func NewPrivateKey() PrivateKey {
	c := elliptic.P521()
	priv, err := ecdsa.GenerateKey(c, rand.Reader)
	if err != nil {
		panic(err)
	}
	return &privateKey{priv}
}

func (k *privateKey) Sign(hash []byte) Signature {
	r, s, err := ecdsa.Sign(rand.Reader, k.PrivateKey, hash)
	if err != nil {
		panic(err)
	}
	return Signature{r, s}
}

func (k *privateKey) Public() PublicKey {
	pub := k.PrivateKey.Public()
	return &publicKey{
		data: pub.(*ecdsa.PublicKey),
	}
}

type Signature struct {
	r, s *big.Int
}

func (sign Signature) IsNull() bool {
	return sign.r == nil || sign.s == nil
}

type PublicKey interface {
	Address() string
	Verify(hash []byte, sign Signature) bool
}

type publicKey struct {
	data *ecdsa.PublicKey
}

func (k *publicKey) Verify(hash []byte, sign Signature) bool {
	return ecdsa.Verify(k.data, hash, sign.r, sign.s)
}

func (k *publicKey) Address() string {
	str := fmt.Sprint(*k.data)
	hash := Hash([]byte(str))
	return encoder.EncodeToString(hash)
}

const HashSize = 32

func Hash(data []byte) []byte {
	hash := sha.Sum256(data)
	return hash[:]
}
