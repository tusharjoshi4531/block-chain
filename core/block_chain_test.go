package core

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tusharjoshi4531/block-chain.git/crypto"
	"github.com/tusharjoshi4531/block-chain.git/network"
	"github.com/tusharjoshi4531/block-chain.git/types"
)

func TestInsertBlock(t *testing.T) {
	bc := NewDefaultBlockChain()

	tx1 := newSignedTransaction(t, []byte("Foo"))
	tx2 := newSignedTransaction(t, []byte("Bar"))

	currBlock := bc.GetGenesis()
	lenBlocks := 100
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

func TestLongestChain(t *testing.T) {
	bc := NewDefaultBlockChain()

	tx1 := newSignedTransaction(t, []byte("Foo"))
	tx2 := newSignedTransaction(t, []byte("Bar"))

	currBlock := bc.GetGenesis()
	lenBlocks := 100
	for i := 1; i <= lenBlocks; i++ {
		prevHash, err := currBlock.Hash()
		assert.Nil(t, err)

		block := newSignedBlock(t, uint32(i), prevHash, []*Transaction{tx1, tx2})
		assert.Nil(t, bc.AddBlock(block))

		currBlock = block
	}

	currBlock = bc.GetGenesis()
	lenBlocks = 150
	for i := 1; i <= lenBlocks; i++ {
		prevHash, err := currBlock.Hash()
		assert.Nil(t, err)

		block := newSignedBlock(t, uint32(i), prevHash, []*Transaction{tx1, tx2})
		assert.Nil(t, bc.AddBlock(block))

		currBlock = block
	}

	assert.Equal(t, uint32(lenBlocks), bc.Height())
	heighestBlock := bc.GetHeighestBlock()

	assert.Equal(t, currBlock, heighestBlock)
}

func TestHasTransaction(t *testing.T) {
	numBlocks := 5
	numTransactionsPerBlock := 10

	txx := make([][]*Transaction, numBlocks)
	bhash := make([]types.Hash, numBlocks)

	bc := NewDefaultBlockChain()

	prevHsh, err := bc.genesis.Hash()
	assert.Nil(t, err)

	// Create block chain
	for i := 0; i < numBlocks; i++ {
		txx[i] = make([]*Transaction, numTransactionsPerBlock)

		for j := 0; j < numTransactionsPerBlock; j++ {
			txx[i][j] = newSignedTransaction(t, []byte(strconv.Itoa(i)+":"+strconv.Itoa(j)))
		}

		block := newSignedBlock(t, uint32(i+1), prevHsh, txx[i])
		assert.Nil(t, bc.AddBlock(block))

		prevHsh, err = block.Hash()
		assert.Nil(t, err)

		bhash[i] = prevHsh
	}

	// Check for transactions
	for i := numBlocks - 1; i >= 0; i-- {
		assert.Nil(t, bc.HasTransactionInChain(txx[i][0].Hash(), bhash[i]))
		if i < numBlocks-1 {
			assert.NotNil(t, bc.HasTransactionInChain(txx[i+1][0].Hash(), bhash[i]))
		}
		assert.NotNil(t, bc.HasTransactionInChain(types.Hash{}, bhash[i]))
	}
}

func TestBlockChainNetwork(t *testing.T) {
	ta := network.NewLocalTransport("A")
	tb := network.NewLocalTransport("B")

	// privKeyA := crypto.GeneratePrivateKey()
	// privKeyB := crypto.GeneratePrivateKey()

	ta.Connect(tb)
	tb.Connect(ta)

	// A creates transaction

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
