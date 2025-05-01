package currency

import (
	"github.com/tusharjoshi4531/block-chain.git/core"
	"github.com/tusharjoshi4531/block-chain.git/types"
)

type BlockChain struct {
	core.DefaultBlockChain
	blockBranch map[types.Hash]string
	state       LedgerState
}

func NewBlockChain(state LedgerState) *BlockChain {
	return &BlockChain{
		DefaultBlockChain: *core.NewDefaultBlockChain(),
		blockBranch:       make(map[types.Hash]string),
		state:             state,
	}
}
