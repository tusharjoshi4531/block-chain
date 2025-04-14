package core

import (
	"bytes"
	"encoding/gob"
	"fmt"
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

func TestSuccessiveEncoding(t *testing.T) {
	a := "HASEAE"
	b := 145
	c := 12.3

	buf := &bytes.Buffer{}
	assert.Nil(t, gob.NewEncoder(buf).Encode(a))
	assert.Nil(t, gob.NewEncoder(buf).Encode(b))
	assert.Nil(t, gob.NewEncoder(buf).Encode(c))

	ra := ""
	rb := 0
	rc := 0.1

	assert.Nil(t, gob.NewDecoder(buf).Decode(&ra))
	assert.Nil(t, gob.NewDecoder(buf).Decode(&rb))
	assert.Nil(t, gob.NewDecoder(buf).Decode(&rc))

	assert.Equal(t, ra, a)
	assert.Equal(t, rb, b)
	assert.Equal(t, rc, c)

	tx := newSignedTransaction(t, []byte("FOOO"))
	fmt.Println(tx.Signature)
	buf = &bytes.Buffer{}

	assert.Nil(t, gob.NewEncoder(buf).Encode(tx.Data))
	assert.Nil(t, gob.NewEncoder(buf).Encode(*crypto.SerializePublicKey(tx.From)))
	assert.Nil(t, tx.Signature.Encode(buf))

	ntx := NewTransaction([]byte{})
	assert.Nil(t, gob.NewDecoder(buf).Decode(&ntx.Data))
	pk := &crypto.SerializedPublicKey{}
	assert.Nil(t, gob.NewDecoder(buf).Decode(pk))
	ntx.From = crypto.DecodePublicKey(pk)

	assert.Nil(t, ntx.Signature.Decode(buf))

	assert.Equal(t, tx.Signature, ntx.Signature)
	assert.Equal(t, tx, ntx)
}

func TestTransactionEncodeDecode(t *testing.T) {
	tx := newSignedTransaction(t, []byte("FOO"))

	txBytes, err := tx.Bytes()
	assert.Nil(t, err)

	txDecoded := NewTransaction([]byte{})
	assert.Nil(t, txDecoded.Decode(bytes.NewBuffer(txBytes)))

	fmt.Println(tx.Signature)
	fmt.Println(txDecoded.Signature)
	assert.Equal(t, txDecoded, tx)
}


func newSignedTransaction(t *testing.T, data []byte) *Transaction {
	tx := NewTransaction(data)
	privateKey := crypto.GeneratePrivateKey()

	assert.Nil(t, tx.Sign(privateKey))
	assert.Nil(t, tx.Verify())

	return tx
}

