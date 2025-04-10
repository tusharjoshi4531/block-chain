package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"math/big"
)

type Signature struct {
	R, S *big.Int
}

func (s *Signature) Verify(publicKey *ecdsa.PublicKey, data []byte) bool {
	return ecdsa.Verify(publicKey, data, s.R, s.S)
}

type SerializedPublicKey struct {
	// TODO: Support for other curves
	X, Y *big.Int
}

func SerializePublicKey(publicKey *ecdsa.PublicKey) SerializedPublicKey {
	return SerializedPublicKey{
		X: publicKey.X,
		Y: publicKey.Y,
	}
}

func DecoderPublicKey(serializedKey *SerializedPublicKey) *ecdsa.PublicKey {
	return &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     serializedKey.X,
	}
}

func SignBytes(privateKey *ecdsa.PrivateKey, data []byte) (*Signature, error) {
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, data)
	if err != nil {
		return nil, err
	}
	return &Signature{
		R: r,
		S: s,
	}, nil
}
