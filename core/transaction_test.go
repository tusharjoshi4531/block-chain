package core

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tusharjoshi4531/block-chain.git/crypto"
)

func TestSignTransaction(t *testing.T) {
	tx := NewTransaction([]byte("FOOO"))
	privateKey := crypto.GeneratePrivateKey()

	assert.Nil(t, tx.Sign(privateKey))
	assert.Nil(t, tx.Verify())
}

func TestTranscactionPool(t *testing.T) {
	txPool := NewDefaultTransactionPool()
	poolSize := 1000
	for i := 0; i < poolSize; i++ {
		tx := NewTransaction([]byte(strconv.Itoa(i)))
		tx.SetFirstSeen(rand.Int63n(1000000000000))

		assert.Nil(t, txPool.AddTransaction(tx))
	}
	assert.Equal(t, txPool.Len(), poolSize)

	txx := txPool.Transactions()
	assert.Equal(t, len(txx), txPool.Len())

	for i := 0; i < txPool.Len()-1; i++ {
		assert.True(t, txx[i].FirstSeen() < txx[i+1].FirstSeen())
	}
}
