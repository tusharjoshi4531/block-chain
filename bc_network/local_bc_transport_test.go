package bcnetwork

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/gob"
	"fmt"
	"sort"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tusharjoshi4531/block-chain.git/core"
	"github.com/tusharjoshi4531/block-chain.git/crypto"
)

func TestLocalPeer(t *testing.T) {
	bc := core.NewDefaultBlockChain()
	txPoolA := core.NewDefaultTransactionPool()
	txPoolB := core.NewDefaultTransactionPool()

	pkA := crypto.GeneratePrivateKey()
	pkB := crypto.GeneratePrivateKey()
	ta := NewLocalBlockChainTransport("a", bc, txPoolA)
	tb := NewLocalBlockChainTransport("b", bc, txPoolB)

	ta.LocalTransport.Connect(tb.LocalTransport)
	tb.LocalTransport.Connect(ta.LocalTransport)

	// Create transactions
	numTxA := 10
	for i := 0; i < numTxA; i++ {
		tx := core.NewTransaction([]byte("helloA" + strconv.Itoa(i)))
		assert.Nil(t, tx.Sign(pkA))
		assert.Nil(t, tx.Verify())

		ta.transactionPool.AddTransaction(tx)
	}
	assert.Equal(t, ta.transactionPool.Len(), numTxA)

	numTxB := 5
	for i := 0; i < numTxB; i++ {
		tx := core.NewTransaction([]byte("helloB" + strconv.Itoa(i)))
		assert.Nil(t, tx.Sign(pkB))
		assert.Nil(t, tx.Verify())

		tb.transactionPool.AddTransaction(tx)
	}
	assert.Equal(t, tb.transactionPool.Len(), numTxB)

	// Transmit transactions
	for _, transaction := range ta.transactionPool.Transactions() {
		ta.SendTransaction(tb.Address(), transaction)

		recMsg := <-tb.ReadChan()
		recPayload := &BCPayload{}
		recPayload.Decode(bytes.NewBuffer(recMsg.Payload))

		assert.Equal(t, recPayload.MsgType, MessageTransaction)
		recTx := core.NewTransaction([]byte{})
		recTx.Decode(bytes.NewBuffer(recPayload.Payload))

		doTransactionsMatch(t, transaction, recTx)
		tb.transactionPool.AddTransaction(recTx)
	}
	assert.Equal(t, tb.transactionPool.Len(), numTxA+numTxB)
	for _, transaction := range tb.transactionPool.Transactions() {
		tb.SendTransaction(ta.Address(), transaction)

		recMsg := <-ta.ReadChan()
		recPayload := &BCPayload{}
		recPayload.Decode(bytes.NewBuffer(recMsg.Payload))

		assert.Equal(t, recPayload.MsgType, MessageTransaction)
		recTx := core.NewTransaction([]byte{})
		recTx.Decode(bytes.NewBuffer(recPayload.Payload))

		doTransactionsMatch(t, transaction, recTx)
		ta.transactionPool.AddTransaction(recTx)
	}
	assert.Equal(t, ta.transactionPool.Len(), numTxA+numTxB)
}

func TestLocalPeerWithReceive(t *testing.T) {
	bc := core.NewDefaultBlockChain()
	txPoolA := core.NewDefaultTransactionPool()
	txPoolB := core.NewDefaultTransactionPool()

	pkA := crypto.GeneratePrivateKey()
	pkB := crypto.GeneratePrivateKey()
	ta := NewLocalBlockChainTransport("a", bc, txPoolA)
	tb := NewLocalBlockChainTransport("b", bc, txPoolB)

	ta.Connect(tb.LocalTransport)
	tb.Connect(ta.LocalTransport)

	// Create transactions
	numTxA := 10
	for i := 0; i < numTxA; i++ {
		tx := core.NewTransaction([]byte("helloA" + strconv.Itoa(i)))
		assert.Nil(t, tx.Sign(pkA))
		assert.Nil(t, tx.Verify())

		ta.transactionPool.AddTransaction(tx)
	}
	assert.Equal(t, ta.transactionPool.Len(), numTxA)

	numTxB := 5
	for i := 0; i < numTxB; i++ {
		tx := core.NewTransaction([]byte("helloB" + strconv.Itoa(i)))
		assert.Nil(t, tx.Sign(pkB))
		assert.Nil(t, tx.Verify())

		tb.transactionPool.AddTransaction(tx)
	}
	assert.Equal(t, tb.transactionPool.Len(), numTxB)

	// Transmit transactions
	for _, transaction := range ta.transactionPool.Transactions() {
		ta.SendTransaction(tb.Address(), transaction)
	}
	for _, transaction := range tb.transactionPool.Transactions() {
		tb.SendTransaction(ta.Address(), transaction)
	}

	for i := 0; i < numTxA; i++ {
		recMsg := <-tb.ReadChan()
		recPayload := &BCPayload{}
		recPayload.Decode(bytes.NewBuffer(recMsg.Payload))

		assert.Equal(t, recPayload.MsgType, MessageTransaction)
		err := tb.ReceiveMessage(recPayload)
		assert.Nil(t, err)
	}

	for i := 0; i < numTxB; i++ {
		recMsg := <-ta.ReadChan()
		recPayload := &BCPayload{}
		recPayload.Decode(bytes.NewBuffer(recMsg.Payload))

		assert.Equal(t, recPayload.MsgType, MessageTransaction)
		err := ta.ReceiveMessage(recPayload)
		assert.Nil(t, err)
	}

	assert.Equal(t, tb.transactionPool.Len(), numTxA+numTxB)
	assert.Equal(t, ta.transactionPool.Len(), numTxA+numTxB)
}

func TestLocalNetwork(t *testing.T) {
	connSize := 3
	bc := make([]core.BlockChain, connSize)

	txPools := make([]*core.DefaultTransactionPool, connSize)
	pks := make([]*ecdsa.PrivateKey, connSize)
	ts := make([]*LocalBlockChainTransport, connSize)

	for i := 0; i < connSize; i++ {
		bc[i] = core.NewDefaultBlockChain()
		txPools[i] = core.NewDefaultTransactionPool()
		pks[i] = crypto.GeneratePrivateKey()
		ts[i] = NewLocalBlockChainTransport(strconv.Itoa(i), bc[i], txPools[i])
	}

	for i := 0; i < connSize; i++ {
		for j := i + 1; j < connSize; j++ {
			assert.Nil(t, ts[i].Connect(ts[j].LocalTransport))
			assert.Nil(t, ts[j].Connect(ts[i].LocalTransport))
		}
	}

	// create tx
	numTx := 10
	for i := 0; i < numTx; i++ {
		tx := core.NewTransaction([]byte("helloA" + strconv.Itoa(i)))
		assert.Nil(t, tx.Sign(pks[0]))
		assert.Nil(t, tx.Verify())

		assert.Nil(t, ts[0].transactionPool.AddTransaction(tx))
	}
	assert.Equal(t, ts[0].transactionPool.Len(), numTx)

	// Send and receive tx
	for _, tx := range ts[0].transactionPool.Transactions() {
		assert.Nil(t, ts[0].BroadcastTransaction(tx))
		for j := 1; j < connSize; j++ {
			recMsg := <-ts[j].ReadChan()
			recPayload := &BCPayload{}
			recPayload.Decode(bytes.NewBuffer(recMsg.Payload))

			assert.Equal(t, recPayload.MsgType, MessageTransaction)
			recTx := core.NewTransaction([]byte{})
			recTx.Decode(bytes.NewBuffer(recPayload.Payload))

			doTransactionsMatch(t, tx, recTx)
			ts[j].transactionPool.AddTransaction(recTx)
		}
	}
}

func TestConcurrentLocalNetworl(t *testing.T) {
	numNodes := 10
	transports := make([]*LocalBlockChainTransport, numNodes)
	pks := make([]*ecdsa.PrivateKey, numNodes)

	for i := 0; i < numNodes; i++ {
		transports[i], pks[i] = createLocalBlockchainTransport(fmt.Sprintf("Node: %d", i))
		for j := 0; j < i; j++ {
			transports[i].Connect(transports[j].LocalTransport)
			transports[j].Connect(transports[i].LocalTransport)
		}
	}

	numTx := 100
	var wg sync.WaitGroup
	wg.Add(numNodes)

	iter := int64(0)
	for i := 0; i < numNodes; i++ {
		// Transaction creation
		go func(i int) {
			for j := 0; j < numTx; j++ {
				tx := core.NewTransaction([]byte(fmt.Sprintf("TX: (%d, %d)", i, j)))
				assert.Nil(t, tx.Sign(pks[i]))
				assert.Nil(t, tx.Verify())

				tx.SetFirstSeen(iter)
				iter++

				assert.Nil(t, transports[i].transactionPool.AddTransaction(tx))

				// send transaction
				for k := 0; k < numNodes; k++ {
					if k != i {
						assert.Nil(t, transports[i].SendTransaction(transports[k].Address(), tx))
					}
				}
			}
		}(i)

		// Transaction receiver
		go func(i int) {
			defer wg.Done()
			for j := 0; j < numTx*(numNodes-1); j++ {
				recMsg := <-transports[i].ReadChan()
				recPayload := &BCPayload{}
				recPayload.Decode(bytes.NewBuffer(recMsg.Payload))

				assert.Equal(t, recPayload.MsgType, MessageTransaction)
				assert.Nil(t, transports[i].ReceiveMessage(recPayload))
			}
		}(i)
	}
	wg.Wait()
	for i := 0; i < numNodes; i++ {
		assert.Equal(t, transports[i].transactionPool.Len(), numTx*numNodes)
	}

	txx := transports[0].transactionPool.Transactions()
	fmt.Println(len(txx))
	sort.Slice(txx, func(i, j int) bool {
		return string(txx[i].Data) < string(txx[j].Data)
	})
	for _, transport := range transports {
		txx2 := transport.transactionPool.Transactions()
		sort.Slice(txx2, func(i, j int) bool {
			return string(txx2[i].Data) < string(txx2[j].Data)
		})
		for i := 0; i < len(txx); i++ {
			doTransactionsMatch(t, txx[i], txx2[i])
		}
	}
	// assert.True(t, false)
}

func TestSendHashChain(t *testing.T) {
	ta, pka := createLocalBlockchainTransport("A")
	tb, pkb := createLocalBlockchainTransport("B")

	assert.Nil(t, ta.Connect(tb))
	assert.Nil(t, tb.Connect(ta))

	ta.blockChain = createDummyBlockcahin(t, 100, 5, pka)
	tb.blockChain = createDummyBlockcahin(t, 100, 5, pkb)

	assert.Equal(t, int(ta.blockChain.Height()), 20)
	assert.Equal(t, int(tb.blockChain.Height()), 20)

	for i := 0; i < 2; i++ {
		var tr, otr *LocalBlockChainTransport
		// var pk, tpk *ecdsa.PrivateKey

		if i == 0 {
			tr, otr = ta, tb
		} else {
			tr, otr = tb, ta
		}

		// Send blocks
		assert.Nil(t, tr.SendBlockChainHash(otr.Address()))

		recMsg := <-otr.ReadChan()

		recPayload := &BCPayload{}
		assert.Nil(t, recPayload.Decode(bytes.NewBuffer(recMsg.Payload)))
		assert.Equal(t, recPayload.MsgType, MessageHashChain)

		// Rec
		chain := &core.HashChain{}
		chain.Decode(bytes.NewBuffer(recPayload.Payload))

		extraBlocks := chain.GetExcludedBlockHashes(tr.blockChain)
		assert.Equal(t, len(extraBlocks), 0)
	}
}

func TestSendHashChainIncorrect(t *testing.T) {
	ta, pka := createLocalBlockchainTransport("A")
	tb, pkb := createLocalBlockchainTransport("B")

	assert.Nil(t, ta.Connect(tb))
	assert.Nil(t, tb.Connect(ta))

	bc := createDummyBlockcahin(t, 100, 5, pka)
	ta.blockChain = bc
	tb.blockChain = bc.Copy()

	hshA, err := ta.blockChain.GetGenesis().Hash()
	assert.Nil(t, err)
	hshB, err := tb.blockChain.GetGenesis().Hash()
	assert.Nil(t, err)
	assert.Equal(t, hshA, hshB)

	assert.Equal(t, int(ta.blockChain.Height()), 20)
	assert.Equal(t, int(tb.blockChain.Height()), 20)

	assert.Equal(t, len(ta.blockChain.GetBlockHashes()), 21)
	assert.Equal(t, len(tb.blockChain.GetBlockHashes()), 21)

	txx := []*core.Transaction{
		core.NewTransaction([]byte("NewA")),
	}
	assert.Nil(t, txx[0].Sign(pka))
	extendBlockChain(t, ta.blockChain, txx, pka)

	txx = []*core.Transaction{
		core.NewTransaction([]byte("NewB")),
	}
	assert.Nil(t, txx[0].Sign(pkb))
	extendBlockChain(t, tb.blockChain, txx, pkb)

	assert.Equal(t, int(ta.blockChain.Height()), 21)
	assert.Equal(t, int(tb.blockChain.Height()), 21)

	assert.Equal(t, len(ta.blockChain.GetBlockHashes()), 22)
	assert.Equal(t, len(tb.blockChain.GetBlockHashes()), 22)

	for i := 0; i < 2; i++ {
		var tr, otr *LocalBlockChainTransport
		// var pk, tpk *ecdsa.PrivateKey

		if i == 0 {
			fmt.Println("A")
			tr, otr = ta, tb
		} else {
			fmt.Println("B")
			tr, otr = tb, ta
		}

		// Send blocks
		assert.Nil(t, tr.SendBlockChainHash(otr.Address()))

		recMsg := <-otr.ReadChan()

		recPayload := &BCPayload{}
		assert.Nil(t, recPayload.Decode(bytes.NewBuffer(recMsg.Payload)))
		assert.Equal(t, recPayload.MsgType, MessageHashChain)

		// Rec
		chain := &core.HashChain{}
		chain.Decode(bytes.NewBuffer(recPayload.Payload))

		assert.Equal(t, len(chain.GetBlockHashes()), len(otr.blockChain.GetBlockHashes()))
		assert.Equal(t, len(chain.GetBlockHashes()), len(tr.blockChain.GetBlockHashes()))

		extraBlocks := chain.GetExcludedBlockHashes(otr.blockChain)
		assert.Equal(t, len(extraBlocks), 1)
	}
}

func TestBlockchainSyncManual(t *testing.T) {
	numTx := 100
	blockSz := 5
	numBlocks := numTx / blockSz
	numDivergeTx := 10
	numDivergeBlocks := numDivergeTx / blockSz

	trs, privKeys := createNetworkWithSameBlocks(t, 2, numTx, blockSz)
	ta, tb := trs[0], trs[1]
	pka, pkb := privKeys[0], privKeys[1]

	assert.Nil(t, ta.Connect(tb))
	assert.Nil(t, tb.Connect(ta))

	bc := createDummyBlockcahin(t, numTx, blockSz, pka)
	ta.blockChain = bc
	tb.blockChain = bc.Copy()

	hshA, err := ta.blockChain.GetGenesis().Hash()
	assert.Nil(t, err)
	hshB, err := tb.blockChain.GetGenesis().Hash()
	assert.Nil(t, err)
	assert.Equal(t, hshA, hshB)

	assert.Equal(t, int(ta.blockChain.Height()), numBlocks)
	assert.Equal(t, int(tb.blockChain.Height()), numBlocks)

	assert.Equal(t, len(ta.blockChain.GetBlockHashes()), numBlocks+1)
	assert.Equal(t, len(tb.blockChain.GetBlockHashes()), numBlocks+1)

	txxa := extendBlockChainAuto(t, ta.blockChain, "TA_", numDivergeTx, blockSz, pka)
	txxb := extendBlockChainAuto(t, tb.blockChain, "TB_", numDivergeTx, blockSz, pkb)

	for _, tx := range txxa {
		ta.transactionPool.AddTransaction(tx)
		tb.transactionPool.AddTransaction(tx)
	}

	for _, tx := range txxb {
		ta.transactionPool.AddTransaction(tx)
		tb.transactionPool.AddTransaction(tx)
	}

	assert.Equal(t, ta.transactionPool.Len(), numTx+2*numDivergeTx)
	assert.Equal(t, tb.transactionPool.Len(), numTx+2*numDivergeTx)

	assert.Equal(t, int(ta.blockChain.Height()), numBlocks+numDivergeBlocks)
	assert.Equal(t, int(tb.blockChain.Height()), numBlocks+numDivergeBlocks)

	assert.Equal(t, len(ta.blockChain.GetBlockHashes()), numBlocks+numDivergeBlocks+1)
	assert.Equal(t, len(tb.blockChain.GetBlockHashes()), numBlocks+numDivergeBlocks+1)

	for i := 0; i < 2; i++ {
		var tr, otr *LocalBlockChainTransport
		// var pk, tpk *ecdsa.PrivateKey

		if i == 0 {
			tr, otr = ta, tb
		} else {
			tr, otr = tb, ta
		}

		// Send block hash
		assert.Nil(t, tr.SendBlockChainHash(otr.Address()))

		recMsg := <-otr.ReadChan()

		recPayload := &BCPayload{}
		assert.Nil(t, recPayload.Decode(bytes.NewBuffer(recMsg.Payload)))
		assert.Equal(t, recPayload.MsgType, MessageHashChain)

		// Rec
		chain := &core.HashChain{}
		chain.Decode(bytes.NewBuffer(recPayload.Payload))

		// assert.Equal(t, len(chain.GetBlockHashes()), len(otr.blockChain.GetBlockHashes()))
		// assert.Equal(t, len(chain.GetBlockHashes()), len(tr.blockChain.GetBlockHashes()))

		extraBlocks := chain.GetExcludedBlockHashes(otr.blockChain)
		assert.Equal(t, len(extraBlocks), numDivergeBlocks)

		// Send blocks
		blocks := make([]*core.Block, 0, len(extraBlocks))
		for _, hash := range extraBlocks {
			block, err := otr.blockChain.GetBlockWithHash(hash)
			assert.Nil(t, err)
			blocks = append(blocks, block)
		}

		assert.Nil(t, otr.SendBlocks(tr.Address(), blocks))

		// Rec blocks
		recMsg = <-tr.ReadChan()

		recPayload = &BCPayload{}
		assert.Nil(t, recPayload.Decode(bytes.NewBuffer(recMsg.Payload)))
		assert.Equal(t, recPayload.MsgType, MessageBlocks)

		len := 0
		buf := bytes.NewBuffer(recPayload.Payload)
		assert.Nil(t, gob.NewDecoder(buf).Decode(&len))
		fmt.Printf("LEN: %d\n", len)

		bl := &core.SerializableBlock{}
		assert.Nil(t, bl.Decode(buf))

		assert.Nil(t, tr.ReceiveMessage(recPayload))
	}

	assert.Equal(t, len(ta.blockChain.GetBlockHashes()), len(tb.blockChain.GetBlockHashes()))
	hashesA := ta.blockChain.GetBlockHashes()
	hashesB := tb.blockChain.GetBlockHashes()
	sort.Slice(hashesA, func(i, j int) bool {
		return hashesA[i].String() < hashesA[j].String()
	})
	sort.Slice(hashesB, func(i, j int) bool {
		return hashesB[i].String() < hashesB[j].String()
	})
	assert.Equal(t, hashesA, hashesB)
}

func createNetworkWithSameBlocks(t *testing.T, numNodes, numTx, blockSz int) ([]*LocalBlockChainTransport, []*ecdsa.PrivateKey) {
	privKeys := make([]*ecdsa.PrivateKey, numNodes)
	for i := 0; i < numNodes; i++ {
		privKeys[i] = crypto.GeneratePrivateKey()
	}

	numBlocks := numTx / blockSz
	txx := make([]*core.Transaction, numTx)

	txPools := make([]*core.DefaultTransactionPool, numNodes)
	for i := 0; i < numNodes; i++ {
		txPools[i] = core.NewDefaultTransactionPool()
	}

	for i := 0; i < numTx; i++ {
		txx[i] = core.NewTransaction([]byte(fmt.Sprintf("DATA: %d", i)))
		assert.Nil(t, txx[i].Sign(privKeys[0]))

		for j := 0; j < numNodes; j++ {
			assert.Nil(t, txPools[j].AddTransaction(txx[i]))
		}
	}

	bc := core.NewDefaultBlockChain()
	prevHash, err := bc.GetGenesis().Hash()
	assert.Nil(t, err)

	j := 0
	for i := 0; i < numBlocks; i++ {
		block := core.NewBlockWithHeaderInfo(bc.Height()+1, prevHash)
		for k := 0; k < blockSz; k++ {
			block.AddTransaction(txx[j])
			j++
		}
		assert.Nil(t, block.Sign(privKeys[0]))
		bc.AddBlock(block)

		prevHash, err = block.Hash()
		assert.Nil(t, err)
	}

	trs := make([]*LocalBlockChainTransport, numTx)
	for i := 0; i < numNodes; i++ {
		trs[i] = NewLocalBlockChainTransport(strconv.Itoa(i), bc.Copy(), txPools[i])
		for j := 0; j < i; j++ {
			trs[i].Connect(trs[j])
			trs[j].Connect(trs[i])
		}
	}

	return trs, privKeys
}

func createLocalBlockchainTransport(address string) (*LocalBlockChainTransport, *ecdsa.PrivateKey) {
	bc := core.NewDefaultBlockChain()
	txPool := core.NewDefaultTransactionPool()
	pk := crypto.GeneratePrivateKey()

	tr := NewLocalBlockChainTransport(address, bc, txPool)
	return tr, pk
}

func doTransactionsMatch(t *testing.T, tx1, tx2 *core.Transaction) bool {
	if !assert.Equal(t, tx1.From, tx2.From) {
		return false
	}
	if !assert.Equal(t, tx1.Signature, tx2.Signature) {
		return false
	}
	if !assert.Equal(t, tx1.Data, tx2.Data) {
		return false
	}
	return true
}

func createDummyBlockcahin(t *testing.T, numTx, blockSz int, privKey *ecdsa.PrivateKey) *core.DefaultBlockChain {
	return createDummyBlockcahinWithPool(t, numTx, blockSz, privKey, core.NewDefaultTransactionPool())
}

func createDummyBlockcahinWithPool(t *testing.T, numTx, blockSz int, privKey *ecdsa.PrivateKey, txPool core.TransactionPool) *core.DefaultBlockChain {
	bc := core.NewDefaultBlockChain()
	j := 0

	prevHash, err := bc.GetGenesis().Hash()
	assert.Nil(t, err)
	currBlock := core.NewBlockWithHeaderInfo(1, prevHash)

	for i := 0; i < numTx; i++ {
		tx := core.NewTransaction([]byte(fmt.Sprintf("%d", i)))
		assert.Nil(t, txPool.AddTransaction(tx))
		currBlock.AddTransaction(tx)
		j++

		if j == blockSz {
			assert.Nil(t, currBlock.Sign(privKey))
			assert.Nil(t, bc.AddBlock(currBlock))
			prevHash, err := currBlock.Hash()
			assert.Nil(t, err)

			currBlock = core.NewBlockWithHeaderInfo(currBlock.Header.Height+1, prevHash)
			j = 0
		}
	}
	return bc
}

func extendBlockChain(t *testing.T, bc core.BlockChain, txx []*core.Transaction, privKey *ecdsa.PrivateKey) {
	height := bc.Height()
	prevHash, err := bc.GetHeighestBlock().Header.Hash()
	assert.Nil(t, err)

	block := core.NewBlockWithHeaderInfo(height+1, prevHash)
	for _, tx := range txx {
		block.AddTransaction(tx)
	}
	block.Sign(privKey)
	bc.AddBlock(block)
}

func extendBlockChainAuto(t *testing.T, bc core.BlockChain, pref string, numTx, blockSz int, privKey *ecdsa.PrivateKey) []*core.Transaction {
	block := bc.GetHeighestBlock()
	prevHash, err := block.Hash()
	assert.Nil(t, err)
	currBlock := core.NewBlockWithHeaderInfo(block.Header.Height+1, prevHash)

	txx := make([]*core.Transaction, numTx)
	j := 0
	for i := 0; i < numTx; i++ {
		tx := core.NewTransaction([]byte(fmt.Sprintf("%s%d", pref, i)))
		currBlock.AddTransaction(tx)
		txx[i] = tx
		j++

		if j == blockSz {
			assert.Nil(t, currBlock.Sign(privKey))
			assert.Nil(t, bc.AddBlock(currBlock))
			prevHash, err := currBlock.Hash()
			assert.Nil(t, err)

			currBlock = core.NewBlockWithHeaderInfo(currBlock.Header.Height+1, prevHash)
			j = 0
		}
	}
	if j != 0 {
		assert.Nil(t, currBlock.Sign(privKey))
		assert.Nil(t, bc.AddBlock(currBlock))
		_, err := currBlock.Hash()
		assert.Nil(t, err)
	}
	return txx
}
