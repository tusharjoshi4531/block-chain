package bcnetwork

import (
	"bytes"
	"fmt"
	"sort"
	"time"

	"github.com/tusharjoshi4531/block-chain.git/core"
	"github.com/tusharjoshi4531/block-chain.git/network"
	"github.com/tusharjoshi4531/block-chain.git/util"
)

type BlockChainTransport interface {
	SendTransaction(string, *core.Transaction) error
	BroadcastTransaction(*core.Transaction) error
	SyncBlockChain(to string) error
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

func (tr *LocalBlockChainTransport) SyncBlockChain(tp string) error {
	return nil
}

func (tr *LocalBlockChainTransport) SendBlockChainHash(to string) error {
	payload, err := NewBCHashChain(core.NewHashChain(tr.blockChain))
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

func (tr *LocalBlockChainTransport) SendBlocks(to string, blocks []*core.Block) error {
	payload, err := NewBCBlocks(blocks)
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
		transaction, err := tr.decodeTransaction(payload.Payload)
		if err != nil {
			return err
		}

		return tr.addTransaction(transaction)
	case MessageBlocks:
		blocks, err := tr.decodeBlocks(payload.Payload)
		if err != nil {
			return err
		}

		return tr.addBlocks(blocks)
	default:
		return fmt.Errorf("incorrect message type (%d)", payload.MsgType)
	}
}

func (tr *LocalBlockChainTransport) decodeTransaction(payload []byte) (*core.Transaction, error) {
	transaction := core.NewTransaction([]byte{})
	err := transaction.Decode(bytes.NewBuffer(payload))
	return transaction, err
}

func (tr *LocalBlockChainTransport) decodeBlocks(payload []byte) ([]*core.SerializableBlock, error) {
	return util.DecodeSlice(bytes.NewBuffer(payload), func() *core.SerializableBlock {
		return &core.SerializableBlock{}
	})
}

func (tr *LocalBlockChainTransport) addTransaction(transaction *core.Transaction) error {
	transaction.SetFirstSeen(time.Now().UnixNano())
	return tr.transactionPool.AddTransaction(transaction)
}

func (tr *LocalBlockChainTransport) addBlocks(blocks []*core.SerializableBlock) error {
	sort.Slice(blocks, func(i, j int) bool {
		return blocks[i].Header.Height < blocks[j].Header.Height
	})
	
	for _, encodedBlock := range blocks {
		block, err := encodedBlock.Reconstruct(tr.transactionPool)
		if err != nil {
			return err
		}
		if err := tr.blockChain.AddBlock(block); err != nil {
			return err
		}
	}
	return nil
}
