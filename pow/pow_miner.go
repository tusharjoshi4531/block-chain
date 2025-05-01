package pow

import (
	"crypto/ecdsa"
	"math/rand"

	"github.com/tusharjoshi4531/block-chain.git/core"
	"github.com/tusharjoshi4531/block-chain.git/prot"
)

type PowMiner struct {
	RequiredPrefixZerosInHex uint8
	blockChain               core.BlockChain
	transactionPool          core.TransactionPool
	privateKey               *ecdsa.PrivateKey
	rewarder                 prot.Rewarder
}

func NewPowMiner(prefixZerosInHex uint8, bc core.BlockChain, txPool core.TransactionPool, privKey *ecdsa.PrivateKey, rewarder prot.Rewarder) *PowMiner {
	return &PowMiner{
		RequiredPrefixZerosInHex: prefixZerosInHex,
		blockChain:               bc,
		transactionPool:          txPool,
		privateKey:               privKey,
		rewarder:                 rewarder,
	}
}

func (miner *PowMiner) MineBlock(transactionsLimit uint32) (*core.Block, error) {
	bc := miner.blockChain
	txPool := miner.transactionPool

	prevBloack := bc.GetHeighestBlock()
	prevHash, err := prevBloack.Hash()
	if err != nil {
		return nil, err
	}

	block := core.NewBlockWithHeaderInfo(bc.Height()+1, prevHash)

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
	reward, err := miner.rewarder.GenerateReward()
	if err != nil {
		return nil, err
	}

	block.AddTransaction(reward)

	currNonceVal := rand.Uint64()
	for {
		block.SetNonce(NewPowNonce(currNonceVal))

		hash, err := block.Hash()
		if err != nil {
			return nil, err
		}

		if validateHash(hash, miner.RequiredPrefixZerosInHex) {
			break
		}
		currNonceVal++
	}

	return block, nil
}
