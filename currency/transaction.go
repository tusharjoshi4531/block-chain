package currency

import (
	"encoding/gob"
	"io"

	"github.com/tusharjoshi4531/block-chain.git/core"
	"github.com/tusharjoshi4531/block-chain.git/util"
)

type Transaction struct {
	From   string
	To     string
	Amount float64
}

func NewTransaction(from string, to string, amount float64) *Transaction {
	return &Transaction{
		From:   from,
		To:     to,
		Amount: amount,
	}
}

func (tx *Transaction) Encode(w io.Writer) error {
	return gob.NewEncoder(w).Encode(tx)
}

func (tx *Transaction) Decode(r io.Reader) error {
	return gob.NewDecoder(r).Decode(tx)
}

func (tx *Transaction) ToBytes() ([]byte, error) {
	return util.EncodeToBytes(tx)
}

func (tx *Transaction) ToCoreTransaction() (*core.Transaction, error) {
	data, err := tx.ToBytes()
	if err != nil {
		return nil, err
	}
	return core.NewTransaction(data), nil
}
