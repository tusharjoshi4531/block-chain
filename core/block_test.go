package core

import (
	"bytes"
	"strconv"
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

func TestSignBlock(t *testing.T) {
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

func TestEncodeDecodeBlock(t *testing.T) {
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

	buf := &bytes.Buffer{}
	assert.Nil(t, block.Encode(buf))

	block2 := NewBlock()
	assert.Nil(t, block2.Decode(buf))

	assert.Equal(t, block.Header, block2.Header)
	assert.Equal(t, block.Validator, block2.Validator)
	assert.Equal(t, block.Signature, block2.Signature)
	assert.Equal(t, block.Transactions, block2.Transactions)
}

func TestHastTransaction(t *testing.T) {
	numTx := 5
	txx := make([]*Transaction, numTx)
	for i := 0; i < numTx; i++ {
		txx[i] = newSignedTransaction(t, []byte(strconv.Itoa(i)))
	}

	block := newSignedBlock(t, 1, types.Hash{}, txx)

	for i := 0; i < numTx; i++ {
		assert.True(t, block.HasTranaction(txx[i].Hash()))
	}
}
