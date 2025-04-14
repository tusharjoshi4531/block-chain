package prot

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tusharjoshi4531/block-chain.git/core"
	"github.com/tusharjoshi4531/block-chain.git/crypto"
)

func TestSimpleMining(t *testing.T) {
	bc := core.NewDefaultBlockChain()
	txPool := core.NewDefaultTransactionPool()
	privKey := crypto.GeneratePrivateKey()

	miner := NewSimpleMiner(bc, txPool, privKey)
	validator := NewSimpleValidator(bc, privKey)

	numTx := 100
	for i := 0; i < numTx; i++ {
		tx := core.NewTransaction([]byte(strconv.Itoa(i)))
		assert.Nil(t, tx.Sign(privKey))
		assert.Nil(t, txPool.AddTransaction(tx))
		assert.Nil(t, tx.Verify())
	}

	numBlocks, blockSz := 3, uint32(10)
	for i := 0; i < numBlocks; i++ {
		block, err := miner.MineBlock(blockSz)
		assert.Nil(t, err)

		// Correct block
		assert.Nil(t, validator.ValidateBlock(block))

		// Incorrect blocks``
		// Incorrect prev Hash
		block1 := core.NewBlock()
		block1.Header.Height = block.Header.Height
		block1.Transactions = block.Transactions
		assert.NotNil(t, validator.ValidateBlock(block1))

		// Incorrect Height
		block2 := core.NewBlock()
		block2.Header.Height = block.Header.Height - 1
		block2.Header.PrevBlockHash = block.Header.PrevBlockHash
		block2.Transactions = block.Transactions
		assert.NotNil(t, validator.ValidateBlock(block2))

		assert.Nil(t, bc.AddBlock(block))
		assert.Equal(t, block.Header.Height, uint32(i + 1))
	}
}
