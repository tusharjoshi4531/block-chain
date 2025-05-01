package pow

import (
	"crypto/ecdsa"
	"fmt"
	"math/rand"

	"github.com/tusharjoshi4531/block-chain.git/core"
	"github.com/tusharjoshi4531/block-chain.git/types"
)

type PowValidator struct {
	RequiredPrefixZerosInHex uint8
}

func NewPowValidator(prefZerosInHex uint8) *PowValidator {
	return &PowValidator{
		RequiredPrefixZerosInHex: prefZerosInHex,
	}
}

func (validator *PowValidator) ValidateBlock(block *core.Block) error {
	hash, err := block.Hash()
	if err != nil {
		return err
	}
	if !validateHash(hash, validator.RequiredPrefixZerosInHex) {
		return fmt.Errorf(
			"block hash (%s) does not contain (%d) zeros in its prefix",
			hash.String(),
			validator.RequiredPrefixZerosInHex,
		)
	}
	return nil
}

type PowMiner struct {
	validator       *PowValidator
	blockChain      core.BlockChain
	transactionPool core.TransactionPool
	privateKey      *ecdsa.PrivateKey
}

func NewPowMiner(validator *PowValidator, bc core.BlockChain, txPool core.TransactionPool, privKey *ecdsa.PrivateKey) *PowMiner {
	return &PowMiner{
		validator:       validator,
		blockChain:      bc,
		transactionPool: txPool,
		privateKey:      privKey,
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
	reward := core.NewTransaction([]byte("Reward"))
	if err := reward.Sign(miner.privateKey); err != nil {
		return nil, err
	}

	block.AddTransaction(reward)
	currNonceVal := rand.Uint64()
	block.SetNonce(NewPowNonce(currNonceVal))

	for miner.validator.ValidateBlock(block) != nil {
		currNonceVal++
		block.SetNonce(NewPowNonce(currNonceVal))
	}

	return block, nil
}

func validateHash(hash types.Hash, prefZerosiInHex uint8) bool {
	hashStr := hash.String()
	for i := uint8(0); i < prefZerosiInHex; i++ {
		if hashStr[i] != '0' {
			return false
		}
	}
	return true
}
