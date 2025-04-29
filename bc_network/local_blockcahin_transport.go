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
	network.Transport
	SendTransaction(string, *core.Transaction) error
	BroadcastTransaction(*core.Transaction) error
	SyncBlockChain(to string) error
	SendHashChain(to string) error
	SendBlocks(to string, blocks []*core.Block) error
	SendBlocksWithHashChain(to string, blocks []*core.Block) error
	BroadcastHashChain() error

	ReceiveMessage(*BCPayload, string) error
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

func (tr *LocalBlockChainTransport) SendHashChain(to string) error {
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

func (tr *LocalBlockChainTransport) SendBlocksWithHashChain(to string, blocks []*core.Block) error {
	payload, err := NewBCBlocksWithHashChain(blocks, core.NewHashChain(tr.blockChain))
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

func (tr *LocalBlockChainTransport) BroadcastHashChain() error {
	payload, err := NewBCHashChain(core.NewHashChain(tr.blockChain))
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

func (tr *LocalBlockChainTransport) ReceiveMessage(payload *BCPayload, from string) error {
	fmt.Printf("REC - FROM: %s, TO: %s, TYPE: %s\n", from, tr.Address(), MsgTypeToString(payload.MsgType))
	switch payload.MsgType {
	case MessageTransaction:
		return tr.handleTransactionMessage(payload.Payload)
	case MessageBlocks:
		return tr.handleBlocksMessage(payload.Payload)
	case MessageHashChain:
		return tr.handleHashChainMessage(payload.Payload, from)
	case MessageBlocksWithHashChain:
		return tr.handleBlocksWithHashChainMessage(payload.Payload, from)
	default:
		return fmt.Errorf("incorrect message type (%d)", payload.MsgType)
	}
}

func (tr *LocalBlockChainTransport) handleTransactionMessage(payload []byte) error {
	transaction, err := tr.decodeTransactionFromBytes(payload)
	if err != nil {
		return err
	}

	return tr.addTransaction(transaction)
}

func (tr *LocalBlockChainTransport) handleBlocksMessage(payload []byte) error {
	blocks, err := tr.decodeBlocksFromBytes(payload)
	if err != nil {
		return err
	}

	return tr.addBlocks(blocks)
}

func (tr *LocalBlockChainTransport) handleHashChainMessage(payload []byte, from string) error {
	hashChain, err := tr.decodeHashChainFromBytes(payload)
	if err != nil {
		return err
	}

	extraBlocks := hashChain.GetExcludedBlocks(tr.blockChain)
	tr.SendBlocksWithHashChain(from, extraBlocks)
	return nil
}

func (tr *LocalBlockChainTransport) handleBlocksWithHashChainMessage(payload []byte, from string) error {
	blocks, hashChain, err := tr.decodeBlocksAndHashchainFromBytes(payload)
	if err != nil {
		return err
	}

	if err := tr.addBlocks(blocks); err != nil {
		return err
	}

	extraBlocks := hashChain.GetExcludedBlocks(tr.blockChain)
	return tr.SendBlocks(from, extraBlocks)
}

func (tr *LocalBlockChainTransport) decodeTransactionFromBytes(payload []byte) (*core.Transaction, error) {
	transaction := core.NewTransaction([]byte{})
	err := transaction.Decode(bytes.NewBuffer(payload))
	return transaction, err
}

func (tr *LocalBlockChainTransport) decodeBlocksFromBytes(payload []byte) ([]*core.Block, error) {
	return util.DecodeSlice(bytes.NewBuffer(payload), func() *core.Block {
		return core.NewBlock()
	})
}

func (tr *LocalBlockChainTransport) decodeHashChainFromBytes(payload []byte) (*core.HashChain, error) {
	hashChain := &core.HashChain{}
	err := hashChain.Decode(bytes.NewBuffer(payload))
	return hashChain, err
}

func (tr *LocalBlockChainTransport) decodeBlocksAndHashchainFromBytes(payload []byte) ([]*core.Block, *core.HashChain, error) {
	buf := bytes.NewBuffer(payload)
	blocks, err := util.DecodeSlice(buf, func() *core.Block {
		return core.NewBlock()
	})
	if err != nil {
		return nil, nil, err
	}
	hashChain := &core.HashChain{}
	if err := hashChain.Decode(buf); err != nil {
		return nil, nil, err
	}
	return blocks, hashChain, err
}

func (tr *LocalBlockChainTransport) addTransaction(transaction *core.Transaction) error {
	transaction.SetFirstSeen(time.Now().UnixNano())
	err := tr.transactionPool.AddTransaction(transaction)
	fmt.Println(tr.transactionPool.Len())
	return err
}

func (tr *LocalBlockChainTransport) addBlocks(blocks []*core.Block) error {
	sort.Slice(blocks, func(i, j int) bool {
		return blocks[i].Header.Height < blocks[j].Header.Height
	})

	for _, block := range blocks {
		if err := tr.blockChain.AddBlock(block); err != nil {
			return err
		}
	}
	return nil
}
