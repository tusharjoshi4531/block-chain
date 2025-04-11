package core

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"fmt"

	"github.com/tusharjoshi4531/block-chain.git/crypto"
	"github.com/tusharjoshi4531/block-chain.git/types"
)

type Transaction struct {
	Data []byte

	From      *ecdsa.PublicKey
	Signature *crypto.Signature

	hash      types.Hash
	firstSeen int64
}

func NewTransaction(data []byte) *Transaction {
	return &Transaction{
		Data: data,
	}
}

func (tx *Transaction) SetFirstSeen(t int64) {
	tx.firstSeen = t
}

func (tx *Transaction) FirstSeen() int64 {
	return tx.firstSeen
}

func (tx *Transaction) Sign(privateKey *ecdsa.PrivateKey) error {
	sig, err := crypto.SignBytes(privateKey, tx.Data)
	if err != nil {
		return err
	}

	tx.Signature = sig
	tx.From = &privateKey.PublicKey

	return nil
}

func (tx *Transaction) Verify() error {
	if tx.Signature == nil {
		return fmt.Errorf("transaction has no signature")
	}

	if !tx.Signature.Verify(tx.From, tx.Data) {
		return fmt.Errorf("incorrect sign in transaction")
	}

	return nil
}

func (tx *Transaction) Hash() types.Hash {
	if tx.hash.IsZero() {
		tx.hash = sha256.Sum256(tx.Data)
	}
	return tx.hash
}
