package server

import (
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	bcnetwork "github.com/tusharjoshi4531/block-chain.git/bc_network"
	"github.com/tusharjoshi4531/block-chain.git/core"
	"github.com/tusharjoshi4531/block-chain.git/crypto"
	"github.com/tusharjoshi4531/block-chain.git/prot"
)

func TestLocalBlockchainServer(t *testing.T) {
	serverA := NewSimpleLocalBlockChainServer("A")
	serverB := NewSimpleLocalBlockChainServer("B")

	fmt.Println(serverA)
	fmt.Println(serverB)

	assert.Nil(t, serverA.Connect(serverB))
	assert.Nil(t, serverB.Connect(serverA))

	serverA.Listen()
	serverB.Listen()

	numTxA := 10
	numTxB := 20

	for i := 0; i < numTxA; i++ {
		tx := core.NewTransaction([]byte(fmt.Sprintf("A: %d", i)))
		tx.SetFirstSeen(time.Now().UnixNano())
		assert.Nil(t, tx.Sign(serverA.privKey))

		assert.Nil(t, serverA.AddTransaction(tx))
	}

	for i := 0; i < numTxB; i++ {
		tx := core.NewTransaction([]byte(fmt.Sprintf("B: %d", i)))
		assert.Nil(t, tx.Sign(serverB.privKey))

		serverA.SendTransaction(serverA.Address(), tx)
	}

	time.Sleep(1 * time.Second)
	serverA.Kill()
	serverB.Kill()

	txxb := serverB.transactionPool.Transactions()
	txxa := serverA.transactionPool.Transactions()
	assert.Equal(t, len(txxa), len(txxb))

	sort.Slice(txxa, func(i, j int) bool {
		return string(txxa[i].Data) < string(txxa[j].Data)
	})
	sort.Slice(txxb, func(i, j int) bool {
		return string(txxb[i].Data) < string(txxb[j].Data)
	})

	for i := 0; i < len(txxa); i++ {
		assert.Equal(t, txxa[i].Data, txxb[i].Data)
		assert.Equal(t, txxa[i].From, txxb[i].From)
		assert.Equal(t, txxa[i].Signature, txxb[i].Signature)
		assert.Equal(t, txxa[i].Hash(), txxb[i].Hash())
	}
}

func TestLocalServerSync(t *testing.T) {
	serverA := NewSimpleLocalBlockChainServer("A")
	serverB := NewSimpleLocalBlockChainServer("B")

	fmt.Println(serverA)
	fmt.Println(serverB)

	assert.Nil(t, serverA.Connect(serverB))
	assert.Nil(t, serverB.Connect(serverA))

	serverA.Listen()
	serverB.Listen()

	numTxA := 10
	numTxB := 20

	for i := 0; i < numTxA; i++ {
		tx := core.NewTransaction([]byte(fmt.Sprintf("A: %d", i)))
		tx.SetFirstSeen(time.Now().UnixNano())
		assert.Nil(t, tx.Sign(serverA.privKey))

		assert.Nil(t, serverA.AddTransaction(tx))
	}

	for i := 0; i < numTxB; i++ {
		tx := core.NewTransaction([]byte(fmt.Sprintf("B: %d", i)))
		tx.SetFirstSeen(time.Now().UnixNano())
		assert.Nil(t, tx.Sign(serverB.privKey))

		assert.Nil(t, serverB.AddTransaction(tx))
	}

	time.Sleep(1 * time.Second / 2)

	txxb := serverB.transactionPool.Transactions()
	txxa := serverA.transactionPool.Transactions()
	assert.Equal(t, len(txxa), len(txxb))

	sort.Slice(txxa, func(i, j int) bool {
		return string(txxa[i].Data) < string(txxa[j].Data)
	})
	sort.Slice(txxb, func(i, j int) bool {
		return string(txxb[i].Data) < string(txxb[j].Data)
	})

	fmt.Println("TXX:")
	for i := 0; i < len(txxa); i++ {
		hsh := txxa[i].Hash()
		str := hsh.String()
		fmt.Println(str)
		assert.Equal(t, txxa[i].Data, txxb[i].Data)
		assert.Equal(t, txxa[i].From, txxb[i].From)
		assert.Equal(t, txxa[i].Signature, txxb[i].Signature)
		assert.Equal(t, txxa[i].Hash(), txxb[i].Hash())
	}
	numBlocks := 3
	blockSz := 5
	for i := 0; i < numBlocks; i++ {
		block, err := serverA.MineBlock(uint32(blockSz))
		assert.Nil(t, err)

		assert.Nil(t, serverA.blockChain.AddBlock(block))
		fmt.Printf("Block%d: \n", i+1)
		for _, tx := range block.Transactions {
			hsh := tx.Hash()
			fmt.Println(hsh.String())
		}
	}

	assert.Nil(t, serverA.BroadcastHashChain())

	time.Sleep(1 * time.Second / 2)
	serverA.Kill()
	serverB.Kill()

	compareBlockchains(t, serverA.blockChain, serverB.blockChain)

}

func compareBlockchains(t *testing.T, bc1 core.BlockChain, bc2 core.BlockChain) {
	hashes1 := bc1.GetBlockHashes()
	hashes2 := bc2.GetBlockHashes()

	sort.Slice(hashes1, func(i, j int) bool {
		return hashes1[i].String() < hashes1[j].String()
	})
	sort.Slice(hashes2, func(i, j int) bool {
		return hashes2[i].String() < hashes2[j].String()
	})
	assert.Equal(t, len(hashes1), len(hashes2))

	for idx, hash := range hashes1 {
		assert.Equal(t, hash.String(), hashes2[idx].String())
	}
}

func createLocalServer(t *testing.T, numTx, blockSz int, pref, addr string) *LocalBlockChainServer {
	numBlocks := numTx / blockSz

	bc := core.NewDefaultBlockChain()
	txPool := core.NewDefaultTransactionPool()
	privKey := crypto.GeneratePrivateKey()

	// Create transactions
	txx := make([]*core.Transaction, numTx)
	for i := 0; i < numTx; i++ {
		txx[i] = core.NewTransaction([]byte(fmt.Sprintf("%s: %d", pref, i)))
		assert.Nil(t, txx[i].Sign(privKey))

		// Update pool
		assert.Nil(t, txPool.AddTransaction(txx[i]))
	}

	prevHash, err := bc.GetGenesis().Hash()
	assert.Nil(t, err)

	// Create blocks
	k := 0
	for i := 0; i < numBlocks; i++ {
		block := core.NewBlockWithHeaderInfo(bc.Height()+1, prevHash)
		for j := 0; j < blockSz; j++ {
			block.AddTransaction(txx[k])
		}
		assert.Nil(t, block.Sign(privKey))
		assert.Nil(t, bc.AddBlock(block))

		prevHash, err = block.Hash()
		assert.Nil(t, err)
	}

	transport := bcnetwork.NewLocalBlockChainTransport(addr, bc, txPool)

	return NewLocalBlockChainServer(
		bc,
		txPool,
		privKey,
		transport,
		func() prot.Miner { return prot.NewSimpleMiner(bc, txPool, privKey) },
		func() prot.Comsumer { return prot.NewSimpleConsumer(bc, txPool, transport) },
		func() prot.Validator { return prot.NewSimpleValidator(bc, privKey) },
	)
}
