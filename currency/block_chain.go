package currency

import (
	"github.com/tusharjoshi4531/block-chain.git/core"
	"github.com/tusharjoshi4531/block-chain.git/types"
	"github.com/tusharjoshi4531/block-chain.git/util"
)

type BlockChain struct {
	core.DefaultBlockChain
	state       LedgerState
	initBalance float64
}

func NewBlockChain(state LedgerState, initBalance float64) *BlockChain {
	return &BlockChain{
		DefaultBlockChain: *core.NewDefaultBlockChain(),
		state:             state,
		initBalance:       initBalance,
	}
}

func (blockChain *BlockChain) AddBlock(block *core.Block) error {
	prevHighestHash, err := blockChain.DefaultBlockChain.GetHeighestBlock().Hash()
	if err != nil {
		return err
	}

	blockHash, err := block.Hash()
	if err != nil {
		return err
	}

	if err := blockChain.DefaultBlockChain.AddBlock(block); err != nil {
		return err
	}

	currHighestHash, err := blockChain.DefaultBlockChain.GetHeighestBlock().Hash()
	if err != nil {
		return err
	}

	// Update ledger
	if currHighestHash == blockHash {
		return blockChain.updateLedger(prevHighestHash, currHighestHash)
	}

	return nil
}

func (blockChain *BlockChain) AddWallet(walletId string) error {
	return blockChain.state.AddWallet(walletId, blockChain.initBalance)
}

func (blockChain *BlockChain) updateLedger(prevHighestHash, currHighestHash types.Hash) error {
	ancestorHash, err := blockChain.commonAncestor(prevHighestHash, currHighestHash)
	if err != nil {
		return err
	}

	if err := blockChain.revertPath(prevHighestHash, ancestorHash); err != nil {
		return err
	}
	if err := blockChain.commitPath(currHighestHash, ancestorHash); err != nil {
		return err
	}
	return nil
}

func (blockChain *BlockChain) revertPath(child, ancestor types.Hash) error {
	revertedBlocks := make([]*core.Block, 0)
	for child != ancestor {
		currBlock, err := blockChain.GetBlockWithHash(child)
		if err != nil {
			return err
		}

		if err := blockChain.revertBlock(currBlock); err != nil {
			// Undo all reverts
			for i := len(revertedBlocks) - 1; i >= 0; i-- {
				blockChain.commitBlock(revertedBlocks[i])
			}
			return err
		}

		revertedBlocks = append(revertedBlocks, currBlock)
		child = currBlock.Header.PrevBlockHash
	}
	return nil
}

func (blockChain *BlockChain) commitPath(child, ancestor types.Hash) error {
	// Store ancestor to child path in a stack
	path := make([]*core.Block, 0)
	for child != ancestor {
		currBlock, err := blockChain.GetBlockWithHash(child)
		if err != nil {
			return err
		}

		path = append(path, currBlock)
		child = currBlock.Header.PrevBlockHash
	}

	// Commit blocks from ancestor to child
	for i := len(path) - 1; i >= 0; i-- {
		if err := blockChain.commitBlock(path[i]); err != nil {
			// Undo all commits
			for j := i + 1; j < len(path); j++ {
				blockChain.revertBlock(path[j])
			}
			return nil
		}
	}

	return nil
}

func (blockChain *BlockChain) revertBlock(block *core.Block) error {
	return blockChain.processTransacitons(
		block,
		func(tx *Transaction) error {
			return blockChain.state.RevertTransaction(tx)
		},
		func(tx *Transaction) {
			blockChain.state.CommitTransaciton(tx)
		},
	)
}

func (blockChain *BlockChain) commitBlock(block *core.Block) error {
	return blockChain.processTransacitons(
		block,
		func(tx *Transaction) error {
			return blockChain.state.CommitTransaciton(tx)
		},
		func(tx *Transaction) {
			blockChain.state.RevertTransaction(tx)
		},
	)
}

func (blockChain *BlockChain) processTransacitons(
	block *core.Block,
	process func(*Transaction) error,
	undo func(*Transaction),
) error {

	processedTx := make([]*Transaction, 0)
	for _, tx := range block.Transactions {
		transaction, err := NewTransactionFromCoreTransaction(tx)
		if err != nil {
			return err
		}

		if err := process(transaction); err != nil {
			// Undo processed transactions
			for i := len(processedTx) - 1; i >= 0; i-- {
				undo(processedTx[i])
			}
			return err
		}

		processedTx = append(processedTx, transaction)
	}
	return nil
}

func (blockChain *BlockChain) commonAncestor(blockHashA, blockHashB types.Hash) (types.Hash, error) {
	hashChainA, err := blockChain.hashChain(blockHashA)
	if err != nil {
		return types.Hash{}, err
	}
	hashChainB, err := blockChain.hashChain(blockHashB)
	if err != nil {
		return types.Hash{}, err
	}

	ancestor := types.Hash{}
	minLen := min(len(hashChainA), len(hashChainB))
	for i := 0; i < minLen; i++ {
		if hashChainA[i] != hashChainB[i] {
			break
		}
		ancestor = hashChainA[i]
	}

	return ancestor, nil
}

func (blockChain *BlockChain) hashChain(blockHash types.Hash) ([]types.Hash, error) {
	chain := make([]types.Hash, 0)
	currHash := blockHash
	for !currHash.IsZero() {
		currBlock, err := blockChain.GetBlockWithHash(currHash)
		if err != nil {
			return nil, err
		}

		chain = append(chain, currHash)
		currHash = currBlock.Header.PrevBlockHash
	}
	util.RevereseSlice(chain)
	return chain, nil
}
