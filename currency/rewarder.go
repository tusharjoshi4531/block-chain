package currency

import (
	"crypto/ecdsa"
	"math"

	"github.com/tusharjoshi4531/block-chain.git/core"
)

type Rewarder struct {
	privateKey         *ecdsa.PrivateKey
	blockChain         core.BlockChain
	initAmount         float64
	blocksToHalfAmount uint16
}

func NewRewarder(privKey *ecdsa.PrivateKey, bc core.BlockChain, initAmount float64, blocksToHalfAmount uint16) *Rewarder {
	return &Rewarder{
		privateKey:         privKey,
		blockChain:         bc,
		initAmount:         initAmount,
		blocksToHalfAmount: blocksToHalfAmount,
	}
}

func (rewarder *Rewarder) GenerateReward(winner string) (*core.Transaction, error) {
	bcHeight := rewarder.blockChain.Height()
	rewardAmount := rewarder.initAmount / math.Pow(2, float64(bcHeight/uint32(rewarder.blocksToHalfAmount)))
	
	reward := NewTransaction("::", winner, rewardAmount)
	rewardTx, err := reward.ToCoreTransaction() 
	if err != nil {
		return nil, err
	}

	if err := rewardTx.Sign(rewarder.privateKey); err != nil {
		return nil, err
	}

	return rewardTx, nil
}
