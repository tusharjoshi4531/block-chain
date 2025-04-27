package core

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tusharjoshi4531/block-chain.git/types"
	"github.com/tusharjoshi4531/block-chain.git/util"
)

func TestSerializeBlock(t *testing.T) {
	block := newSignedBlock(t, 1, types.Hash{}, []*Transaction{})
	txPool := NewDefaultTransactionPool()

	numTx := 100
	for i := 0; i < numTx; i++ {
		tx := newSignedTransaction(t, []byte(fmt.Sprintf("DATA: %d", i)))
		assert.Nil(t, txPool.AddTransaction(tx))
		block.AddTransaction(tx)
	}

	serBlock := NewSerializableBlock(block)
	decBlock, err := serBlock.Reconstruct(txPool)
	assert.Nil(t, err)

	assert.Equal(t, decBlock.Header, block.Header)
	assert.Equal(t, decBlock.Transactions, block.Transactions)
	assert.Equal(t, decBlock.Signature, block.Signature)
	assert.Equal(t, decBlock.Validator, block.Validator)
}

func TestSerializeBlockSlice(t *testing.T) {
	numTx := 100
	blockSz := 4
	numBlocks := numTx / blockSz

	txPool := NewDefaultTransactionPool()
	txx := make([]*Transaction, numTx)

	for i := 0; i < numTx; i++ {
		tx := newSignedTransaction(t, []byte(fmt.Sprintf("DATA: %d", i)))
		txx[i] = tx
		assert.Nil(t, txPool.AddTransaction(tx))
	}

	blocks := make([]*Block, numBlocks)
	j := 0
	for i := 0; i < numBlocks; i++ {
		btxx := make([]*Transaction, blockSz)
		for k := 0; k < blockSz; k++ {
			btxx[k] = txx[j]
			j++
		}
		blocks[i] = newSignedBlock(t, 1, types.Hash{}, btxx)
	}

	encBlocks := make([]*SerializableBlock, numBlocks)
	for i := 0; i < numBlocks; i++ {
		encBlocks[i] = NewSerializableBlock(blocks[i])
	}

	buf := &bytes.Buffer{}
	assert.Nil(t, util.EncodeSlice(buf, util.ToEncoderSlice(encBlocks)))

	// testDecodeSlice(t, buf)

	dec, err := util.DecodeSlice[*SerializableBlock](buf, func() *SerializableBlock {
		return &SerializableBlock{};
	})
	assert.Nil(t, err)
	assert.Equal(t, len(encBlocks), numBlocks)
	assert.Equal(t, len(dec), numBlocks)

	for i := 0; i < numBlocks; i++ {
		assert.Equal(t, encBlocks[i].Header, dec[i].Header)
	}
}

func testDecodeSlice(t *testing.T, r io.Reader) []*SerializableBlock {
	// ln := 0
	// assert.Nil(t, gob.NewDecoder(r).Decode(&ln))
	// fmt.Println(ln)

	// blck := &SerializableBlock{}
	// assert.Nil(t, blck.Decode(r))
	
	// fmt.Println(blck.Header.Height)
	// fmt.Println(len(blck.Transactions))

	// assert.False(t, true)

	var length int
	assert.Nil(t, gob.NewDecoder(r).Decode(&length))
	fmt.Println(length)

	items := make([]*SerializableBlock, length)
	for i := 0; i < length - 2; i++ {
		assert.Nil(t, items[i].Decode(r))
	}

	return items
}
