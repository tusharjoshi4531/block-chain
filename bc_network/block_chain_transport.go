package bcnetwork

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"sort"
	"time"

	"github.com/tusharjoshi4531/block-chain.git/core"
	"github.com/tusharjoshi4531/block-chain.git/network"
	"github.com/tusharjoshi4531/block-chain.git/util"
)

type BlockChainTransportSender interface {
	network.Transport
	SendTransaction(string, *core.Transaction) error
	SendHashChain(to string) error
	SendBlocks(to string, blocks []*core.Block) error
	SendBlocksWithHashChain(to string, blocks []*core.Block) error
	SendWalletId(to string, walletId string) error
	BroadcastTransaction(*core.Transaction) error
	BroadcastHashChain() error
	BroadcastWalletId(walletId string) error
}

type BlockChainTransportProcessor interface {
	AddWallet(walletId string) error
	ProcessMessage(*BCPayload, string) error
}

type BlockChainTransport interface {
	BlockChainTransportSender
	BlockChainTransportProcessor
}

type DefaultBlockChainTransport struct {
	network.Transport
	blockChain      core.BlockChain
	transactionPool core.TransactionPool
}

func NewDefaultBlockChainTransport(transport network.Transport, blockChain core.BlockChain, transactionPool core.TransactionPool) *DefaultBlockChainTransport {
	return &DefaultBlockChainTransport{
		Transport:       transport,
		blockChain:      blockChain,
		transactionPool: transactionPool,
	}
}

func (tr *DefaultBlockChainTransport) SendTransaction(to string, transaction *core.Transaction) error {
	payload, err := NewBCTransactionPayload(transaction)
	if err != nil {
		return err
	}
	payloadBytes, err := payload.Bytes()
	if err != nil {
		return err
	}

	msg := network.NewMessage(tr.Address(), payloadBytes)
	return tr.SendMessageTo(to, msg)
}

func (tr *DefaultBlockChainTransport) SyncBlockChain(tp string) error {
	return nil
}

func (tr *DefaultBlockChainTransport) SendHashChain(to string) error {
	payload, err := NewBCHashChain(core.NewHashChain(tr.blockChain))
	if err != nil {
		return err
	}
	payloadBytes, err := payload.Bytes()
	if err != nil {
		return err
	}
	msg := network.NewMessage(tr.Address(), payloadBytes)
	return tr.SendMessageTo(to, msg)
}

func (tr *DefaultBlockChainTransport) SendBlocks(to string, blocks []*core.Block) error {
	payload, err := NewBCBlocks(blocks)
	if err != nil {
		return err
	}
	payloadBytes, err := payload.Bytes()
	if err != nil {
		return err
	}
	msg := network.NewMessage(tr.Address(), payloadBytes)
	return tr.SendMessageTo(to, msg)
}

func (tr *DefaultBlockChainTransport) SendBlocksWithHashChain(to string, blocks []*core.Block) error {
	payload, err := NewBCBlocksWithHashChain(blocks, core.NewHashChain(tr.blockChain))
	if err != nil {
		return err
	}
	payloadBytes, err := payload.Bytes()
	if err != nil {
		return err
	}
	msg := network.NewMessage(tr.Address(), payloadBytes)
	return tr.SendMessageTo(to, msg)
}

func (tr *DefaultBlockChainTransport) SendWalletId(to, walletId string) error {
	payload, err := NewBCWalletId(walletId)
	if err != nil {
		return err
	}
	payloadBytes, err := payload.Bytes()
	if err != nil {
		return err
	}
	msg := network.NewMessage(tr.Address(), payloadBytes)
	return tr.SendMessageTo(to, msg)
}

func (tr *DefaultBlockChainTransport) BroadcastTransaction(transaction *core.Transaction) error {
	payload, err := NewBCTransactionPayload(transaction)
	if err != nil {
		return err
	}
	payloadBytes, err := payload.Bytes()
	if err != nil {
		return err
	}

	msg := network.NewMessage(tr.Address(), payloadBytes)
	return tr.BroadCastMessage(msg)
}

func (tr *DefaultBlockChainTransport) BroadcastHashChain() error {
	payload, err := NewBCHashChain(core.NewHashChain(tr.blockChain))
	if err != nil {
		return err
	}
	payloadBytes, err := payload.Bytes()
	if err != nil {
		return err
	}
	msg := network.NewMessage(tr.Address(), payloadBytes)
	return tr.BroadCastMessage(msg)
}

func (tr *DefaultBlockChainTransport) BroadcastWalletId(walletId string) error {
	payload, err := NewBCWalletId(walletId)
	if err != nil {
		return err
	}
	payloadBytes, err := payload.Bytes()
	if err != nil {
		return err
	}
	msg := network.NewMessage(tr.Address(), payloadBytes)
	return tr.BroadCastMessage(msg)
}

func (tr *DefaultBlockChainTransport) AddWallet(walletId string) error {
	if err := tr.blockChain.AddWallet(walletId); err != nil {
		return err
	}
	return tr.BroadcastWalletId(walletId)
}

func (tr *DefaultBlockChainTransport) ProcessMessage(payload *BCPayload, from string) error {
	switch payload.MsgType {
	case MessageTransaction:
		return tr.handleTransactionMessage(payload.Payload)
	case MessageBlocks:
		return tr.handleBlocksMessage(payload.Payload)
	case MessageHashChain:
		return tr.handleHashChainMessage(payload.Payload, from)
	case MessageBlocksWithHashChain:
		return tr.handleBlocksWithHashChainMessage(payload.Payload, from)
	case MessageWalletId:
		return tr.handleWalletId(payload.Payload)
	default:
		return fmt.Errorf("incorrect message type (%d)", payload.MsgType)
	}
}

func (tr *DefaultBlockChainTransport) handleTransactionMessage(payload []byte) error {
	transaction, err := tr.decodeTransactionFromBytes(payload)
	if err != nil {
		return err
	}

	return tr.addTransaction(transaction)
}

func (tr *DefaultBlockChainTransport) handleBlocksMessage(payload []byte) error {
	blocks, err := tr.decodeBlocksFromBytes(payload)
	if err != nil {
		return err
	}

	return tr.addBlocks(blocks)
}

func (tr *DefaultBlockChainTransport) handleHashChainMessage(payload []byte, from string) error {
	hashChain, err := tr.decodeHashChainFromBytes(payload)
	if err != nil {
		return err
	}

	extraBlocks := hashChain.GetExcludedBlocks(tr.blockChain)
	tr.SendBlocksWithHashChain(from, extraBlocks)
	return nil
}

func (tr *DefaultBlockChainTransport) handleBlocksWithHashChainMessage(payload []byte, from string) error {
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

func (tr *DefaultBlockChainTransport) handleWalletId(payload []byte) error {
	walletId, err := decodeWalletIdFromBytes(payload)
	if err != nil {
		return err
	}

	return tr.blockChain.AddWallet(walletId)
}

func (tr *DefaultBlockChainTransport) decodeTransactionFromBytes(payload []byte) (*core.Transaction, error) {
	transaction := core.NewTransaction([]byte{})
	err := transaction.Decode(bytes.NewBuffer(payload))
	return transaction, err
}

func (tr *DefaultBlockChainTransport) decodeBlocksFromBytes(payload []byte) ([]*core.Block, error) {
	return util.DecodeSlice(bytes.NewBuffer(payload), func() *core.Block {
		return core.NewBlock()
	})
}

func (tr *DefaultBlockChainTransport) decodeHashChainFromBytes(payload []byte) (*core.HashChain, error) {
	hashChain := &core.HashChain{}
	err := hashChain.Decode(bytes.NewBuffer(payload))
	return hashChain, err
}

func (tr *DefaultBlockChainTransport) decodeBlocksAndHashchainFromBytes(payload []byte) ([]*core.Block, *core.HashChain, error) {
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

func decodeWalletIdFromBytes(payload []byte) (string, error) {
	buf := bytes.NewBuffer(payload)
	walletId := ""
	if err := gob.NewDecoder(buf).Decode(&walletId); err != nil {
		return "", err
	}
	return walletId, nil
}

func (tr *DefaultBlockChainTransport) addTransaction(transaction *core.Transaction) error {
	transaction.SetFirstSeen(time.Now().UnixNano())
	err := tr.transactionPool.AddTransaction(transaction)
	fmt.Println(tr.transactionPool.Len())
	return err
}

func (tr *DefaultBlockChainTransport) addBlocks(blocks []*core.Block) error {
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
