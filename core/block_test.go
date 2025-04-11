package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tusharjoshi4531/block-chain.git/crypto"
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

func TestSignBloc(t *testing.T) {
	block := NewBlock()
	tx1 := newSignedTransaction(t, []byte("FOOO"))
	tx2 := newSignedTransaction(t, []byte("BARR"))
	block.Transactions = append(block.Transactions, tx1)
	block.Transactions = append(block.Transactions, tx2)

	block.Hash()
	assert.Nil(t, DefaultValidator{}.ValidateBlock(block))

	privateKey := crypto.GeneratePrivateKey()
	assert.Nil(t, block.Sign(privateKey))
	assert.Nil(t, block.Verify())

	otherPrivKey := crypto.GeneratePrivateKey()
	block.Validator = &otherPrivKey.PublicKey

	assert.NotNil(t, block.Verify())
}
