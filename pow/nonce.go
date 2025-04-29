package pow

import (
	"encoding/gob"
	"io"

	"github.com/tusharjoshi4531/block-chain.git/types"
)

type PowNonce struct {
	value types.Hash
}

func (nonce *PowNonce) Encode(w io.Writer) error {
	return gob.NewEncoder(w).Encode(nonce.value)
}

func (nonce *PowNonce) Decoder(r io.Reader) error {
	return gob.NewDecoder(r).Decode(&nonce.value)
}

