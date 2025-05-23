package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/gob"
	"io"
	"math/big"
)

type Signature struct {
	R, S *big.Int
}

func NewNilSignature() *Signature {
	return &Signature{
		R: big.NewInt(-1),
		S: big.NewInt(-1),
	}
}

func (s *Signature) Verify(publicKey *ecdsa.PublicKey, data []byte) bool {
	return ecdsa.Verify(publicKey, data, s.R, s.S)
}

func (s *Signature) Encode(w io.Writer) error {
	return gob.NewEncoder(w).Encode(s)
}

func (s *Signature) Decode(r io.Reader) error {
	return gob.NewDecoder(r).Decode(s)
}

type SerializablePublicKey struct {
	// TODO: Support for other curves
	X, Y *big.Int
}

func (key *SerializablePublicKey) Encode(w io.Writer) error {
	return gob.NewEncoder(w).Encode(key)
}

func (key *SerializablePublicKey) Decode(r io.Reader) error {
	return gob.NewDecoder(r).Decode(key)
}

func SerializePublicKey(publicKey *ecdsa.PublicKey) *SerializablePublicKey {
	return &SerializablePublicKey{
		X: publicKey.X,
		Y: publicKey.Y,
	}
}

func DecodePublicKey(serializedKey *SerializablePublicKey) *ecdsa.PublicKey {
	return &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     serializedKey.X,
		Y:     serializedKey.Y,
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
