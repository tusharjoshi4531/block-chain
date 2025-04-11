package core

import (
	"fmt"
	"sort"

	"github.com/tusharjoshi4531/block-chain.git/types"
)

type TransactionPool interface {
	AddTransaction(tx *Transaction) error
	Len() int
}

type DefaultTransactionPool struct {
	transacitons map[types.Hash]*Transaction
}

func NewDefaultTransactionPool() *DefaultTransactionPool {
	return &DefaultTransactionPool{
		transacitons: make(map[types.Hash]*Transaction),
	}
}

func (txPool *DefaultTransactionPool) AddTransaction(tx *Transaction) error {

	hash := tx.Hash()
	if _, ok := txPool.transacitons[hash]; ok {
		return fmt.Errorf("transaction (%s) already present in pool", hash)
	}
	txPool.transacitons[hash] = tx
	return nil
}

func (txPool *DefaultTransactionPool) Transactions() []*Transaction {
	transactions := make([]*Transaction, 0, txPool.Len())
	for _, transaction := range txPool.transacitons {
		transactions = append(transactions, transaction)
	}

	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].FirstSeen() < transactions[j].FirstSeen()
	})

	return transactions
}

func (txPool *DefaultTransactionPool) Len() int {
	return len(txPool.transacitons)
}
