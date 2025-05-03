package pow

import "encoding/gob"

type PowNonce struct {
	Value uint64
}

func NewPowNonce(value uint64) *PowNonce {
	return &PowNonce{
		Value: value,
	}
}

func init() {
	gob.Register(&PowNonce{})
}
