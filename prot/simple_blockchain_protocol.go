package prot

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/tusharjoshi4531/block-chain.git/core"
)

type SimpleValidator struct {
	blockChain core.BlockChain
	privateKey *ecdsa.PrivateKey
}

func NewSimpleValidator(blockChain core.BlockChain, privateKey *ecdsa.PrivateKey) *SimpleValidator {
	return &SimpleValidator{
		blockChain: blockChain,
		privateKey: privateKey,
	}
}

func (valdator *SimpleValidator) ValidateBlock(block *core.Block) error {
	bc := valdator.blockChain
	if block.Header.Height != bc.Height()+1 {
		return fmt.Errorf("expected block height (%d); founc (%d)", bc.Height()+1, block.Header.Height)
	}
	prevBlock, err := bc.GetPrevBlock(block)
	if err != nil {
		return err
	}
	if prevBlock.Header.Height+1 != block.Header.Height {
		return fmt.Errorf("previous block height (%d) does not match header height (%d)", prevBlock.Header.Height, block.Header.Height)
	}

	if err := block.Sign(valdator.privateKey); err != nil {
		return err;
	}
	return nil
}

type SimpleMiner struct {
	blockChain      core.BlockChain
	transactionPool core.TransactionPool
	privateKey      *ecdsa.PrivateKey
}

func (miner *SimpleMiner) MineBlock(transactionsLimit uint32) (*core.Block, error) {
	bc := miner.blockChain
	txPool := miner.transactionPool

	prevBloack := bc.GetHeighestBlock()
	prevHash, err := prevBloack.Hash()
	if err != nil {
		return nil, err
	}

	block := core.NewBlock()
	block.Header.PrevBlockHash = prevHash
	block.Header.Height = bc.Height() + 1

	transactions := txPool.Transactions()
	numTx := uint32(0)
	for _, transaction := range transactions {
		if numTx == transactionsLimit {
			break
		}

		if bc.HasTransactionInChain(transaction.Hash(), prevHash) == nil {
			continue
		}

		numTx++
		block.AddTransaction(transaction)
	}
	reward := core.NewTransaction([]byte("Reward"))
	if err := reward.Sign(miner.privateKey); err != nil {
		return nil, err
	}

	block.AddTransaction(reward)

	return block, nil
}
