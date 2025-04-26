package prot

import (
	"crypto/ecdsa"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	bcnetwork "github.com/tusharjoshi4531/block-chain.git/bc_network"
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
		tx.SetFirstSeen(int64(i))
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
		assert.Equal(t, block.Header.Height, uint32(i+1))
		fmt.Println(len(block.Transactions))
		assert.Equal(t, len(block.Transactions), int(blockSz+1))
	}

}

func TestGetTransactions(t *testing.T) {
	txPool := core.NewDefaultTransactionPool()
	bc := core.NewDefaultBlockChain()
	consumer, validator, miner, privKey := createSimpleParticipant("a", txPool, bc)

	numTx := 100
	txx := make([]*core.Transaction, 0)
	for i := 0; i < numTx; i++ {
		tx := core.NewTransaction([]byte(strconv.Itoa(i)))
		tx.SetFirstSeen(int64(i))
		assert.Nil(t, tx.Sign(privKey))
		assert.Nil(t, txPool.AddTransaction(tx))
		assert.Nil(t, tx.Verify())
		txx = append(txx, tx)
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
		assert.Equal(t, block.Header.Height, uint32(i+1))
		assert.Equal(t, len(block.Transactions), int(blockSz+1))
	}

	txx2, err := consumer.GetTransactions()
	assert.Nil(t, err)

	assert.Equal(t, numBlocks*(int(blockSz)+1), len(txx2))
	j := 0
	for i := 0; i < len(txx2); i++ {
		if i%11 == 10 {
			j++
		} else {
			fmt.Println(i, j, string(txx[i-j].Data), string(txx2[i].Data))
			assert.Equal(t, txx[i-j], txx2[i])
		}
	}
}

func createSimpleParticipant(address string, txPool core.TransactionPool, bc core.BlockChain) (*SimpleConsumer, *SimpleValidator, *SimpleMiner, *ecdsa.PrivateKey) {
	transport := bcnetwork.NewLocalBlockChainTransport(address, bc, txPool)

	privKey := crypto.GeneratePrivateKey()
	miner := NewSimpleMiner(bc, txPool, privKey)
	validator := NewSimpleValidator(bc, privKey)
	consumer := NewSimpleConsumer(bc, transport)
	return consumer, validator, miner, privKey
}
