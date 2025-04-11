package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tusharjoshi4531/block-chain.git/crypto"
	"github.com/tusharjoshi4531/block-chain.git/types"
)

func TestInsertBlock(t *testing.T) {
	bc := NewDefaultBlockChain()

	tx1 := newSignedTransaction(t, []byte("Foo"))
	tx2 := newSignedTransaction(t, []byte("Bar"))

	currBlock := bc.GetGenesis()
	lenBlocks := 2
	for i := 1; i <= lenBlocks; i++ {
		prevHash, err := currBlock.Hash()
		assert.Nil(t, err)

		block := newSignedBlock(t, uint32(i), prevHash, []*Transaction{tx1, tx2})
		assert.Nil(t, bc.AddBlock(block))

		currBlock = block
	}

	assert.Equal(t, uint32(lenBlocks), bc.Height())

	block := newSignedBlock(t, uint32(lenBlocks+1), types.Hash{}, []*Transaction{})
	assert.NotNil(t, bc.AddBlock(block))

	prevHash, err := currBlock.Hash()
	assert.Nil(t, err)

	block = newSignedBlock(t, uint32(lenBlocks+10), prevHash, []*Transaction{})
	assert.NotNil(t, bc.AddBlock(block))
}

func newSignedBlock(t *testing.T, height uint32, prevHash types.Hash, txx []*Transaction) *Block {
	block := NewBlock()
	block.Header.Height = height
	block.Header.PrevBlockHash = prevHash
	for _, tx := range txx {
		block.AddTransaction(tx)
	}
	privKey := crypto.GeneratePrivateKey()
	assert.Nil(t, block.Sign(privKey))
	assert.Nil(t, block.Verify())
	assert.Nil(t, DefaultValidator{}.ValidateBlock(block))

	return block
}
