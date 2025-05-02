package currency

import (
	"crypto/ecdsa"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tusharjoshi4531/block-chain.git/core"
	"github.com/tusharjoshi4531/block-chain.git/crypto"
)

func TestBlockChainLedger(t *testing.T) {
	state := NewMemoryLedgerState()
	initBal := float64(1000)
	privKey := crypto.GeneratePrivateKey()
	bc := NewBlockChain(state, initBal)

	bc.AddMember("A")
	bc.AddMember("B")

	// commonBlocks := make([]*core.Block, 0)

	prevHash, err := bc.GetGenesis().Hash()
	assert.Nil(t, err)

	block := core.NewBlockWithHeaderInfo(1, prevHash)
	block.AddTransaction(createTransaction(t, "A", "B", 10, privKey))
	assert.Nil(t, bc.AddBlock(block))
	prevHash, err = block.Hash()
	assert.Nil(t, err)

	block = core.NewBlockWithHeaderInfo(2, prevHash)
	block.AddTransaction(createTransaction(t, "A", "B", 10, privKey))
	assert.Nil(t, bc.AddBlock(block))
	prevHash, err = block.Hash()
	assert.Nil(t, err)

	block = core.NewBlockWithHeaderInfo(3, prevHash)
	block.AddTransaction(createTransaction(t, "A", "B", 10, privKey))
	assert.Nil(t, bc.AddBlock(block))
	ancestorHash, err := block.Hash()
	assert.Nil(t, err)

	balanceA, err := state.GetBalance("A")
	assert.Nil(t, err)
	assert.Equal(t, balanceA, float64(initBal-30))

	balanceB, err := state.GetBalance("B")
	assert.Nil(t, err)
	assert.Equal(t, balanceB, float64(initBal+30))

	// Branching
	// Branch A
	prevHash = ancestorHash

	block = core.NewBlockWithHeaderInfo(4, prevHash)
	block.AddTransaction(createTransaction(t, RewardSymbol, "B", 100, privKey))
	assert.Nil(t, bc.AddBlock(block))
	prevHash, err = block.Hash()
	assert.Nil(t, err)

	block = core.NewBlockWithHeaderInfo(5, prevHash)
	block.AddTransaction(createTransaction(t, "B", "A", 10, privKey))
	assert.Nil(t, bc.AddBlock(block))
	prevHash, err = block.Hash()
	assert.Nil(t, err)

	balanceA, err = state.GetBalance("A")
	assert.Nil(t, err)
	assert.Equal(t, balanceA, float64(initBal-20))

	balanceB, err = state.GetBalance("B")
	assert.Nil(t, err)
	assert.Equal(t, balanceB, float64(initBal+120))

	// Branch B
	prevHash = ancestorHash

	block = core.NewBlockWithHeaderInfo(4, prevHash)
	block.AddTransaction(createTransaction(t, RewardSymbol, "B", 500, privKey))
	assert.Nil(t, bc.AddBlock(block))
	prevHash, err = block.Hash()
	assert.Nil(t, err)

	block = core.NewBlockWithHeaderInfo(5, prevHash)
	block.AddTransaction(createTransaction(t, RewardSymbol, "A", 500, privKey))
	assert.Nil(t, bc.AddBlock(block))
	prevHash, err = block.Hash()
	assert.Nil(t, err)

	block = core.NewBlockWithHeaderInfo(6, prevHash)
	block.AddTransaction(createTransaction(t, "B", "A", 50, privKey))
	assert.Nil(t, bc.AddBlock(block))
	_, err = block.Hash()
	assert.Nil(t, err)

	balanceA, err = state.GetBalance("A")
	assert.Nil(t, err)
	assert.Equal(t, balanceA, float64(initBal-30+500+50))

	balanceB, err = state.GetBalance("B")
	assert.Nil(t, err)
	assert.Equal(t, balanceB, float64(initBal+30+500-50))

}

func createTransaction(t *testing.T, from, to string, val float64, privKey *ecdsa.PrivateKey) *core.Transaction {
	tx := NewTransaction(from, to, val)
	_tx, err := tx.ToCoreTransaction()
	assert.Nil(t, err)
	assert.Nil(t, _tx.Sign(privKey))
	return _tx
}
