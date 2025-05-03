package tcp

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	bcnetwork "github.com/tusharjoshi4531/block-chain.git/bc_network"
	"github.com/tusharjoshi4531/block-chain.git/core"
	"github.com/tusharjoshi4531/block-chain.git/crypto"
	"github.com/tusharjoshi4531/block-chain.git/currency"
	"github.com/tusharjoshi4531/block-chain.git/network"
	"github.com/tusharjoshi4531/block-chain.git/pow"
	"github.com/tusharjoshi4531/block-chain.git/util"
)

func TestEncodeBlocksWithNonce(t *testing.T) {
	block1 := core.NewBlock()
	block1.SetNonce(pow.NewPowNonce(5))

	block2 := core.NewBlock()

	blocks := []*core.Block{block1, block2}

	payload, err := util.EncodeSliceToBytes(util.ToEncoderSlice(blocks))
	assert.Nil(t, err)

	nblocks, err := util.DecodeSlice(bytes.NewBuffer(payload), func() *core.Block { return core.NewBlock() })
	assert.Nil(t, err)

	assert.Equal(t, block1.Header, nblocks[0].Header)
	assert.Equal(t, block2.Header, nblocks[1].Header)
}

func TestEncodeBlockFromServer(t *testing.T) {
	ledger := currency.NewMemoryLedgerState()
	bc := currency.NewBlockChain(ledger, 10000)
	txPool := core.NewDefaultTransactionPool()
	privKey := crypto.GeneratePrivateKey()
	bcTransport := bcnetwork.NewDefaultBlockChainTransport(
		network.NewDefaultTransport("net"),
		bc,
		txPool,
	)

	server := NewTcpServer(
		ledger,
		bc,
		txPool,
		privKey,
		bcTransport,
	)

	server.AddWallet("A")
	server.ConnectPeer(network.NewLocalTransport("B"))
	for i := 0; i < 100; i++ {
		tx := core.NewTransaction([]byte("SDFa" + strconv.Itoa(i)))
		assert.Nil(t, tx.Sign(privKey))
		assert.Nil(t, server.AddTransaction(tx))
	}

	for i := 0; i < 10; i++ {
		block, err := server.MineBlock(3, "A")
		assert.Nil(t, err)
		server.BlockChain.AddBlock(block)
	}

	// extendBlockChain(t, server.BlockChain, "init", 100, 10)
	hc := core.NewHashChain(server.BlockChain)

	for i := 0; i < 2; i++ {
		block, err := server.MineBlock(3, "A")
		assert.Nil(t, err)
		server.BlockChain.AddBlock(block)
	}

	extBlocks := hc.GetExcludedBlocks(server.BlockChain)
	assert.Equal(t, len(extBlocks), 2)

	assert.Nil(t, server.SendBlocks("B", extBlocks))
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
