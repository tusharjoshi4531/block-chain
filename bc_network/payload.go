package bcnetwork

import (
	"bytes"
	"encoding/gob"
	"io"

	"github.com/tusharjoshi4531/block-chain.git/core"
	"github.com/tusharjoshi4531/block-chain.git/util"
)

const (
	MessageTransaction int = iota
	MessageHashChain
	MessageBlocks
	MessageBlocksWithHashChain
	MessageWalletId
	// MessageTXSync
)

func MsgTypeToString(msgType int) string {
	switch msgType {
	case MessageTransaction:
		return "Transaction"
	case MessageHashChain:
		return "HashChain"
	case MessageBlocks:
		return "Blocks"
	case MessageBlocksWithHashChain:
		return "HashChainWithBlocks"
	case MessageWalletId:
		return "WalletId"
	default:
		return "Invalid"
	}
}

type BCPayload struct {
	MsgType int
	Payload []byte
}

func NewBCTransactionPayload(tx *core.Transaction) (*BCPayload, error) {
	transactionBytes, err := tx.Bytes()
	if err != nil {
		return nil, err
	}

	return &BCPayload{
		MsgType: MessageTransaction,
		Payload: transactionBytes,
	}, nil
}

func NewBCHashChain(hashChain *core.HashChain) (*BCPayload, error) {
	payload, err := hashChain.Bytes()
	if err != nil {
		return nil, err
	}

	return &BCPayload{
		MsgType: MessageHashChain,
		Payload: payload,
	}, nil
}

func NewBCBlocks(blocks []*core.Block) (*BCPayload, error) {
	payload, err := encodeBlocksToBytes(blocks)
	if err != nil {
		return nil, err
	}

	return &BCPayload{
		MsgType: MessageBlocks,
		Payload: payload,
	}, nil
}

func NewBCBlocksWithHashChain(blocks []*core.Block, hashChain *core.HashChain) (*BCPayload, error) {
	payload, err := encodeBlocksWithHashChainBytes(blocks, hashChain)
	if err != nil {
		return nil, err
	}

	return &BCPayload{
		MsgType: MessageBlocksWithHashChain,
		Payload: payload,
	}, nil
}

func NewBCWalletId(walletId string) (*BCPayload, error) {
	buf := &bytes.Buffer{}
	if err := gob.NewEncoder(buf).Encode(walletId); err != nil {
		return nil, err
	}

	return &BCPayload{
		MsgType: MessageWalletId,
		Payload: buf.Bytes(),
	}, nil
}

func (payload *BCPayload) Encode(w io.Writer) error {
	return gob.NewEncoder(w).Encode(payload)
}

func (payload *BCPayload) Bytes() ([]byte, error) {
	buf := &bytes.Buffer{}
	if err := payload.Encode(buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (payload *BCPayload) Decode(r io.Reader) error {
	return gob.NewDecoder(r).Decode(payload)
}

func encodeBlocks(w io.Writer, blocks []*core.Block) error {
	return util.EncodeSlice(w, util.ToEncoderSlice(blocks))
}

func encodeBlocksWithHashChain(w io.Writer, blocks []*core.Block, hashChain *core.HashChain) error {
	if err := encodeBlocks(w, blocks); err != nil {
		return err
	}
	if err := hashChain.Encode(w); err != nil {
		return err
	}
	return nil
}

func encodeBlocksToBytes(blocks []*core.Block) ([]byte, error) {
	return util.EncodeSliceToBytes(util.ToEncoderSlice(blocks))
}

func encodeBlocksWithHashChainBytes(blocks []*core.Block, hashChain *core.HashChain) ([]byte, error) {
	return util.EncodeToBytesUsingEncoder(func(w io.Writer) error {
		return encodeBlocksWithHashChain(w, blocks, hashChain)
	})
}
