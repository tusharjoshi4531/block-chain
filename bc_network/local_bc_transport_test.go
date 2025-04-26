package bcnetwork

import (
	"bytes"
	"crypto/ecdsa"
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

		recMsg := <-tb.Receive()
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

		recMsg := <-ta.Receive()
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
		recMsg := <-tb.Receive()
		recPayload := &BCPayload{}
		recPayload.Decode(bytes.NewBuffer(recMsg.Payload))

		assert.Equal(t, recPayload.MsgType, MessageTransaction)
		err := tb.ReceiveMessage(recPayload)
		assert.Nil(t, err)
	}

	for i := 0; i < numTxB; i++ {
		recMsg := <-ta.Receive()
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
			recMsg := <-ts[j].Receive()
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
				recMsg := <-transports[i].Receive()
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

	ta.Connect(tb)
	tb.Connect(ta)

	ta.blockChain = createDummyBlockcahin(t, 100, 5, pka)
	tb.blockChain = createDummyBlockcahin(t, 100, 5, pkb)

	assert.Equal(t, int(ta.blockChain.Height()), 20)
	assert.Equal(t, int(tb.blockChain.Height()), 20)
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
	bc := core.NewDefaultBlockChain()
	j := 0

	prevHash, err := bc.GetGenesis().Hash()
	assert.Nil(t, err)
	currBlock := core.NewBlockWithHeaderInfo(1, prevHash)

	for i := 0; i < numTx; i++ {
		tx := core.NewTransaction([]byte(fmt.Sprintf("%d", i)))
		currBlock.AddTransaction(tx)
		j++

		if j == blockSz {
			assert.Nil(t, currBlock.Sign(privKey))
			assert.Nil(t, bc.AddBlock(currBlock))
			prevHash, err := currBlock.Hash()
			assert.Nil(t, err)

			currBlock = core.NewBlockWithHeaderInfo(currBlock.Header.Height+1, prevHash)
			j = 0;
		}
	}
	return bc
}
