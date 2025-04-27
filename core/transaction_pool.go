package core

import (
	"fmt"
	"sort"
	"sync"

	"github.com/tusharjoshi4531/block-chain.git/types"
)

type TransactionPool interface {
	AddTransaction(tx *Transaction) error
	Len() int
	Transactions() []*Transaction
	GetTransaction(types.Hash) (*Transaction, error)
	HasTransaction(types.Hash) bool
}

type DefaultTransactionPool struct {
	mu           sync.RWMutex
	transacitons map[types.Hash]*Transaction
}

func NewDefaultTransactionPool() *DefaultTransactionPool {
	return &DefaultTransactionPool{
		transacitons: make(map[types.Hash]*Transaction),
	}
}

func (txPool *DefaultTransactionPool) AddTransaction(tx *Transaction) error {
	txPool.mu.Lock()
	defer txPool.mu.Unlock()

	hash := tx.Hash()
	if _, ok := txPool.transacitons[hash]; ok {
		return fmt.Errorf("transaction (%s) already present in pool", hash)
	}
	txPool.transacitons[hash] = tx
	return nil
}

func (txPool *DefaultTransactionPool) Transactions() []*Transaction {
	txPool.mu.RLock()
	defer txPool.mu.RUnlock()

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
	txPool.mu.RLock()
	defer txPool.mu.RUnlock()

	return len(txPool.transacitons)
}

func (txPool *DefaultTransactionPool) GetTransaction(hash types.Hash) (*Transaction, error) {
	transaction, ok := txPool.transacitons[hash]
	if !ok {
		return nil, fmt.Errorf("transaction with hash (%s) does not exist in the pool", hash.String())
	}
	return transaction, nil
}

func (txPool *DefaultTransactionPool) HasTransaction(hash types.Hash) bool {
	_, ok := txPool.transacitons[hash]
	return ok
}
