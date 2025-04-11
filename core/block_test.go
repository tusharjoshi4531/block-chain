package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tusharjoshi4531/block-chain.git/types"
)

func TestBlockValidator(t *testing.T) {
	block := NewBlock()
	tx1 := newSignedTransaction(t, []byte("FOO"))
	tx2 := newSignedTransaction(t, []byte("BAR"))

	block.AddTransaction(tx1)
	block.AddTransaction(tx2)

	assert.NotNil(t, DefaultValidator{}.ValidateBlock(block))
	block.Hash()

	assert.Nil(t, DefaultValidator{}.ValidateBlock(block))

	block.Header.DataHash = types.Hash{}
	assert.NotNil(t, DefaultValidator{}.ValidateBlock(block))
}
