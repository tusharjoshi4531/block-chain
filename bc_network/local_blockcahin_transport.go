package bcnetwork

import (
	"bytes"
	"fmt"
	"time"

	"github.com/tusharjoshi4531/block-chain.git/core"
	"github.com/tusharjoshi4531/block-chain.git/network"
)

type BlockChainTransport interface {
	SendTransaction(string, *core.Transaction) error
	BroadcastTransaction(*core.Transaction) error
	SendNewBlocks(to string) error
	ReceiveMessage(*BCPayload) error
}

type LocalBlockChainTransport struct {
	*network.LocalTransport
	blockChain      core.BlockChain
	transactionPool core.TransactionPool
}

func NewLocalBlockChainTransport(address string, blockChain core.BlockChain, transactionPool core.TransactionPool) *LocalBlockChainTransport {
	return &LocalBlockChainTransport{
		LocalTransport:  network.NewLocalTransport(address),
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

func (tr *LocalBlockChainTransport) SendNewBlocks(to string) error {
	payload, err := NewBCHashChain(tr.blockChain)
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

func (tr *LocalBlockChainTransport) ReceiveMessage(payload *BCPayload) error {
	switch payload.MsgType {
	case MessageTransaction:
		transaction, err := tr.decodeTransaction(payload)
		if err != nil {
			return err
		}

		return tr.receiveTransaction(transaction)
	default:
		return fmt.Errorf("incorrect message type (%d)", payload.MsgType)
	}
}

func (tr *LocalBlockChainTransport) decodeTransaction(payload *BCPayload) (*core.Transaction, error) {
	transaction := core.NewTransaction([]byte{})
	err := transaction.Decode(bytes.NewBuffer(payload.Payload))
	return transaction, err
}

func (tr *LocalBlockChainTransport) receiveTransaction(transaction *core.Transaction) error {
	transaction.SetFirstSeen(time.Now().UnixNano())
	return tr.transactionPool.AddTransaction(transaction)
}
