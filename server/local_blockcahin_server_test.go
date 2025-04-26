package server

import (
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tusharjoshi4531/block-chain.git/core"
)

func TestLocalBlockchainServer(t *testing.T) {
	serverA := NewLocalBlockChainTransport("A")
	serverB := NewLocalBlockChainTransport("B")

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
		assert.Nil(t, tx.Sign(serverA.privKey))

		serverA.SendTransaction(serverB.Address(), tx)
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
	assert.Empty(t, len(txxa), len(txxb))

	sort.Slice(txxa, func(i, j int) bool {
		return string(txxa[i].Data) < string(txxa[j].Data)
	})
	sort.Slice(txxb, func(i, j int) bool {
		return string(txxb[i].Data) < string(txxb[j].Data)
	})

	for i := 0; i < len(txxa); i++ {
		assert.Equal(t, txxa[i], txxb[i])
	}
}
