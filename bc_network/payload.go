package bcnetwork

import (
	"bytes"
	"encoding/gob"
	"io"
	"sort"

	"github.com/tusharjoshi4531/block-chain.git/core"
	"github.com/tusharjoshi4531/block-chain.git/util"
)

const (
	MessageTransaction int = iota
	MessageHashChain
	MessageBlocks
	// MessageTXSync
)

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
	encodedBlocks := make([]*core.SerializableBlock, 0, len(blocks))
	for _, block := range blocks {
		encodedBlocks = append(encodedBlocks, core.NewSerializableBlock(block))
	}

	sort.Slice(encodedBlocks, func(i, j int) bool {
		return encodedBlocks[i].Header.Height < encodedBlocks[j].Header.Height
	})

	encoderSlice := util.ToEncoderSlice(encodedBlocks)
	payload, err := util.EncodeSliceToBytes(encoderSlice)
	if err != nil {
		return nil, err
	}

	return &BCPayload{
		MsgType: MessageBlocks,
		Payload: payload,
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
