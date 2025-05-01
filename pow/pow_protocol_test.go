package pow

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tusharjoshi4531/block-chain.git/core"
	"github.com/tusharjoshi4531/block-chain.git/crypto"
	"github.com/tusharjoshi4531/block-chain.git/prot"
	"github.com/tusharjoshi4531/block-chain.git/types"
)

func TestValidateHash(t *testing.T) {
	hash := types.Hash{}
	assert.True(t, hash.IsZero())

	assert.True(t, validateHash(hash, uint8(1)))
}

func TestMineBlock(t *testing.T) {
	bc := core.NewDefaultBlockChain()
	txPool := core.NewDefaultTransactionPool()
	privKey := crypto.GeneratePrivateKey()

	numTx := 100
	for i := 0; i < numTx; i++ {
		tx := core.NewTransaction([]byte(fmt.Sprintf("DATA: %d", i)))
		assert.Nil(t, tx.Sign(privKey))
		assert.Nil(t, txPool.AddTransaction(tx))
	}

	numBlocks := 4
	blockSz := 20

	prefZeros := uint8(2)
	validator := NewPowValidator(prefZeros)
	rewarder := prot.NewSimpleRewarder(privKey)
	miner := NewPowMiner(prefZeros, bc, txPool, privKey, rewarder)

	for i := 0; i < numBlocks; i++ {
		block, err := miner.MineBlock(uint32(blockSz))
		assert.Nil(t, err)

		assert.Nil(t, validator.ValidateBlock(block))
	}
}
