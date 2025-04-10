package core

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/tusharjoshi4531/block-chain.git/crypto"
)

type Transaction struct {
	Data []byte

	From      *ecdsa.PublicKey
	Signature *crypto.Signature
}

func NewTransaction(data []byte) *Transaction {
	return &Transaction{
		Data: data,
	}
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
