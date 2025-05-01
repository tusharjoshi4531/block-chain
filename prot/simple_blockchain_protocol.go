package prot

import (
	"crypto/ecdsa"
	"fmt"

	bcnetwork "github.com/tusharjoshi4531/block-chain.git/bc_network"
	"github.com/tusharjoshi4531/block-chain.git/core"
)

type SimpleRewarder struct {
	privateKey *ecdsa.PrivateKey
}

func NewSimpleRewarder(privKey *ecdsa.PrivateKey) *SimpleRewarder {
	return &SimpleRewarder{
		privateKey: privKey,
	}
}

func (rewarder *SimpleRewarder) GenerateReward(winner string) (*core.Transaction, error) {
	tx := core.NewTransaction([]byte("REWARD"))
	err := tx.Sign(rewarder.privateKey)
	return tx, err
}

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
		return err
	}
	return nil
}

type SimpleMiner struct {
	blockChain      core.BlockChain
	transactionPool core.TransactionPool
	rewarder        Rewarder
	privateKey      *ecdsa.PrivateKey
}

func NewSimpleMiner(blockChain core.BlockChain, transactionPool core.TransactionPool, privateKey *ecdsa.PrivateKey) *SimpleMiner {
	return &SimpleMiner{
		blockChain:      blockChain,
		transactionPool: transactionPool,
		privateKey:      privateKey,
		rewarder:        NewSimpleRewarder(privateKey),
	}
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
	reward, err := miner.rewarder.GenerateReward("")
	if err != nil {
		return nil, err
	}

	block.AddTransaction(reward)

	return block, nil
}

type SimpleConsumer struct {
	blockChain      core.BlockChain
	transactionPool core.TransactionPool
	// privateKey      *ecdsa.PrivateKey
	transport bcnetwork.BlockChainTransport
}

func NewSimpleConsumer(blockChain core.BlockChain, transactionPool core.TransactionPool, transport bcnetwork.BlockChainTransport) *SimpleConsumer {
	return &SimpleConsumer{
		blockChain:      blockChain,
		transactionPool: transactionPool,
		transport:       transport,
	}
}

func (consumer *SimpleConsumer) AddTransaction(transaction *core.Transaction) error {
	if err := consumer.transactionPool.AddTransaction(transaction); err != nil {
		return err
	}
	return consumer.transport.BroadcastTransaction(transaction)
}

func (consumer *SimpleConsumer) GetTransactions() ([]*core.Transaction, error) {
	// TODO: implement get transaction in longest chain
	blockChain := consumer.blockChain
	block := blockChain.GetHeighestBlock()
	blockHash, err := block.Hash()
	if err != nil {
		return nil, err
	}
	return blockChain.GetTransactionsInChain(blockHash)
}
