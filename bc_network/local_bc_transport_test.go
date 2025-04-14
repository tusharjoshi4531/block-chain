package bcnetwork

import (
	"bytes"
	"crypto/ecdsa"
	"strconv"
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

	ta.LocalTransport.Connect(&tb.LocalTransport)
	tb.LocalTransport.Connect(&ta.LocalTransport)

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
	assert.Equal(t, tb.transactionPool.Len(), numTxA+numTxB)
}

func TestLocalNetwork(t *testing.T) {
	connSize := 10
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
			assert.Nil(t, ts[i].Connect(&ts[j].LocalTransport))
			assert.Nil(t, ts[j].Connect(&ts[i].LocalTransport))
		}
	}

	// create tx
	numTx := 100
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

func doTransactionsMatch(t *testing.T, tx1, tx2 *core.Transaction) {
	assert.Equal(t, tx1.From, tx2.From)
	assert.Equal(t, tx1.Signature, tx2.Signature)
	assert.Equal(t, tx1.Data, tx2.Data)
}
