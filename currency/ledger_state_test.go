package currency

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	bcnetwork "github.com/tusharjoshi4531/block-chain.git/bc_network"
	"github.com/tusharjoshi4531/block-chain.git/core"
	"github.com/tusharjoshi4531/block-chain.git/crypto"
	"github.com/tusharjoshi4531/block-chain.git/util"
)



func TestEncodeBlock(t *testing.T) {

	bc := createBlockChain(t, 100, 5)
	assert.Equal(t, len(bc.GetBlockHashes()), 6)
	assert.Equal(t, bc.Height(), uint32(5))
	hc := core.NewHashChain(bc)

	extendBlockChain(t, bc, "ext", 10, 2)

	extBlocks := hc.GetExcludedBlocks(bc)
	assert.Equal(t, len(extBlocks), 2)

	payload, err := bcnetwork.NewBCBlocks(extBlocks)

	// data, err := util.EncodeSliceToBytes(util.ToEncoderSlice(blocks))
	assert.Nil(t, err)

	nblocks, err := util.DecodeSlice(bytes.NewBuffer(payload.Payload), func() *core.Block { return core.NewBlock() })
	assert.Nil(t, err)

	ha, err := extBlocks[0].Hash()
	assert.Nil(t, err)
	hb, err := nblocks[0].Hash()
	assert.Nil(t, err)
	assert.Equal(t, ha, hb)

}

func TestCommitTransaction(t *testing.T) {
	state := NewMemoryLedgerState()

	assert.Nil(t, state.AddWallet("A", 100))
	assert.Nil(t, state.AddWallet("B", 100))

	tx1 := NewTransaction(RewardSymbol, "A", 50)
	tx2 := NewTransaction("A", "B", 13)

	assert.Nil(t, state.CommitTransaciton(tx1))
	assert.Nil(t, state.CommitTransaciton(tx2))

	balance, err := state.GetBalance("A")
	assert.Nil(t, err)
	fmt.Println(balance)
	assert.Equal(t, balance, float64(100+50-13))

	balance, err = state.GetBalance("B")
	assert.Nil(t, err)
	assert.Equal(t, balance, float64(100+13))
}

func TestRevertTransaction(t *testing.T) {
	state := NewMemoryLedgerState()

	assert.Nil(t, state.AddWallet("A", 100))
	assert.Nil(t, state.AddWallet("B", 100))

	tx1 := NewTransaction(RewardSymbol, "A", 50)
	tx2 := NewTransaction("A", "B", 13)

	assert.Nil(t, state.CommitTransaciton(tx1))
	assert.Nil(t, state.CommitTransaciton(tx2))

	balance, err := state.GetBalance("A")
	assert.Nil(t, err)
	fmt.Println(balance)
	assert.Equal(t, balance, float64(100+50-13))

	balance, err = state.GetBalance("B")
	assert.Nil(t, err)
	assert.Equal(t, balance, float64(100+13))

	assert.Nil(t, state.RevertTransaction(tx2))

	balance, err = state.GetBalance("A")
	assert.Nil(t, err)
	fmt.Println(balance)
	assert.Equal(t, balance, float64(100+50))

	balance, err = state.GetBalance("B")
	assert.Nil(t, err)
	assert.Equal(t, balance, float64(100))
}

func TestInvalidTransaction(t *testing.T) {
	state := NewMemoryLedgerState()

	assert.Nil(t, state.AddWallet("A", 100))
	assert.Nil(t, state.AddWallet("B", 100))

	tx1 := NewTransaction(RewardSymbol, "A", 50)
	tx2 := NewTransaction("A", "B", 200)

	assert.Nil(t, state.CommitTransaciton(tx1))
	assert.NotNil(t, state.CommitTransaciton(tx2))

	balance, err := state.GetBalance("A")
	assert.Nil(t, err)
	fmt.Println(balance)
	assert.Equal(t, balance, float64(100+50))

	balance, err = state.GetBalance("B")
	assert.Nil(t, err)
	assert.Equal(t, balance, float64(100))
}

func TestAddWallet(t *testing.T) {
	state := NewMemoryLedgerState()

	assert.Nil(t, state.AddWallet("A", 100))
	assert.Nil(t, state.AddWallet("B", 100))

	tx1 := NewTransaction(RewardSymbol, "A", 50)
	tx2 := NewTransaction("A", "B", 200)

	assert.Nil(t, state.CommitTransaciton(tx1))
	assert.NotNil(t, state.CommitTransaciton(tx2))

	balance, err := state.GetBalance("A")
	assert.Nil(t, err)
	fmt.Println(balance)
	assert.Equal(t, balance, float64(100+50))

	balance, err = state.GetBalance("B")
	assert.Nil(t, err)
	assert.Equal(t, balance, float64(100))

	assert.Nil(t, state.AddWallet("C", 1000))
	tx3 := NewTransaction("C", "A", 500)

	assert.Nil(t, state.CommitTransaciton(tx3))

	balance, err = state.GetBalance("A")
	assert.Nil(t, err)
	fmt.Println(balance)
	assert.Equal(t, balance, float64(100+50+500))

	balance, err = state.GetBalance("B")
	assert.Nil(t, err)
	assert.Equal(t, balance, float64(100))

	balance, err = state.GetBalance("C")
	assert.Nil(t, err)
	assert.Equal(t, balance, float64(1000-500))

	assert.Equal(t, len(state.GetWallets()), 3)
}

func createBlockChain(t *testing.T, numTx, numBlocks int) core.BlockChain {
	bc := NewBlockChain(NewMemoryLedgerState(), 100)
	extendBlockChain(t, bc, "init", numTx, numBlocks)
	return bc
}

func extendBlockChain(t *testing.T, bc core.BlockChain, pref string, numTx, numBlocks int) {
	privKey := crypto.GeneratePrivateKey()
	blockSz := numTx / numBlocks

	txx := make([]*core.Transaction, numTx)
	for i := 0; i < numTx; i++ {
		txx[i] = core.NewTransaction([]byte(pref + strconv.Itoa(i)))
		assert.Nil(t, txx[i].Sign(privKey))
	}

	phash, err := bc.GetHeighestBlock().Hash()
	assert.Nil(t, err)
	k := 0
	for i := 0; i < numBlocks; i++ {
		fmt.Println("H: ", bc.Height())
		block := core.NewBlockWithHeaderInfo(bc.Height()+1, phash)
		for j := 0; j < blockSz; j++ {
			block.AddTransaction(txx[k])
			k++
		}
		assert.Nil(t, block.Sign(privKey))
		assert.Nil(t, bc.AddBlock(block))
		phash, err = block.Hash()
		assert.Nil(t, err)
	}
}
