package bcnetwork

import (
	"github.com/tusharjoshi4531/block-chain.git/core"
	"github.com/tusharjoshi4531/block-chain.git/network"
)

type BlockChainTransport interface {
	SendTransaction(string, *core.Transaction) error
	BroadcastTransaction(*core.Transaction) error
}

type LocalBlockChainTransport struct {
	network.LocalTransport
	blockChain      core.BlockChain
	transactionPool core.TransactionPool
}

func NewLocalBlockChainTransport(address string, blockChain core.BlockChain, transactionPool core.TransactionPool) *LocalBlockChainTransport {
	return &LocalBlockChainTransport{
		LocalTransport:  *network.NewLocalTransport(address),
		blockChain:      blockChain,
		transactionPool: transactionPool,
	}
}

func (tr *LocalBlockChainTransport) SendTransaction(to string, transaction *core.Transaction) error {
	payload, err := NewBCTransactionPayload(transaction)
	if err != nil {
		return err
	}
	payloadBytes, err := payload.Bytes()
	if err != nil {
		return err
	}

	msg := network.NewMessage(tr.Address(), payloadBytes)
	return tr.LocalTransport.SendMessage(to, msg)
}

func (tr *LocalBlockChainTransport) BroadcastTransaction(transaction *core.Transaction) error {
	payload, err := NewBCTransactionPayload(transaction)
	if err != nil {
		return err
	}
	payloadBytes, err := payload.Bytes()
	if err != nil {
		return err
	}

	msg := network.NewMessage(tr.Address(), payloadBytes)
	return tr.LocalTransport.BroadCastMessage(msg)
}
