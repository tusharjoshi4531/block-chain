package prot

import "github.com/tusharjoshi4531/block-chain.git/core"

type Rewarder interface {
	GenerateReward(winner string) (*core.Transaction, error)
}

type Validator interface {
	ValidateBlock(block *core.Block) error
}

type Miner interface {
	MineBlock(transactionsLimit uint32, minerWalletId string) (*core.Block, error)
}

type Comsumer interface {
	AddTransaction(transaction *core.Transaction) error
	GetTransactions() ([]*core.Transaction, error)
}
