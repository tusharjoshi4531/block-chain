package types

import "encoding/hex"

type Hash [32]byte

func (h *Hash) IsZero() bool {
	for i := 0; i < 32; i++ {
		if h[i] != 0 {
			return false
		}
	}
	return true
}

func (h *Hash) String() string {
	return hex.EncodeToString(h[:])
}
