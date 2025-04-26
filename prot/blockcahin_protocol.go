package prot

import "github.com/tusharjoshi4531/block-chain.git/core"

type Validator interface {
	ValidateBlock(*core.Block) error
}

type Miner interface {
	MineBlock(uint32) (*core.Block, error)
}

type Comsumer interface {
	AddTransaction(*core.Transaction) error
	GetTransactions() ([]*core.Transaction, error)
}
