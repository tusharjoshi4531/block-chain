package currency

import (
	bcnetwork "github.com/tusharjoshi4531/block-chain.git/bc_network"
	"github.com/tusharjoshi4531/block-chain.git/core"
)

type Consumer struct {
	blockChain      core.BlockChain
	transactionPool core.TransactionPool
	transport       bcnetwork.BlockChainTransport
}

func NewSimpleConsumer(
	blockChain core.BlockChain,
	transactionPool core.TransactionPool,
	transport bcnetwork.BlockChainTransport,
) *Consumer {
	return &Consumer{
		blockChain:      blockChain,
		transactionPool: transactionPool,
		transport:       transport,
	}
}

func (consumer *Consumer) AddTransaction(transaction Transaction) error {
	tx, err := transaction.ToCoreTransaction()
	if err != nil {
		return err
	}
	if err := consumer.transactionPool.AddTransaction(tx); err != nil {
		return err
	}
	return consumer.transport.BroadcastTransaction(tx)
}
