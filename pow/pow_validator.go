package pow

import (
	"fmt"

	"github.com/tusharjoshi4531/block-chain.git/core"
	"github.com/tusharjoshi4531/block-chain.git/types"
)

type PowValidator struct {
	RequiredPrefixZerosInHex uint8
}

func NewPowValidator(prefZerosInHex uint8) *PowValidator {
	return &PowValidator{
		RequiredPrefixZerosInHex: prefZerosInHex,
	}
}

func (validator *PowValidator) ValidateBlock(block *core.Block) error {
	hash, err := block.Hash()
	if err != nil {
		return err
	}
	if !validateHash(hash, validator.RequiredPrefixZerosInHex) {
		return fmt.Errorf(
			"block hash (%s) does not contain (%d) zeros in its prefix",
			hash.String(),
			validator.RequiredPrefixZerosInHex,
		)
	}
	return nil
}

func validateHash(hash types.Hash, prefZerosiInHex uint8) bool {
	hashStr := hash.String()
	for i := uint8(0); i < prefZerosiInHex; i++ {
		if hashStr[i] != '0' {
			return false
		}
	}
	return true
}
