package core

import (
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
