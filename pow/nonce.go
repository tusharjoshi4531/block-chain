package pow

import (
	"encoding/gob"
	"io"
)

type PowNonce struct {
	value uint64
}

func NewPowNonce(value uint64) *PowNonce {
	return &PowNonce{
		value: value,
	}
}

func (nonce *PowNonce) Encode(w io.Writer) error {
	return gob.NewEncoder(w).Encode(nonce.value)
}

func (nonce *PowNonce) Decode(r io.Reader) error {
	return gob.NewDecoder(r).Decode(&nonce.value)
}
