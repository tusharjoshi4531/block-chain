package bcnetwork

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"

	"github.com/tusharjoshi4531/block-chain.git/core"
)

const (
	MessageTransaction int = iota
	MessageHashChain
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

func NewBCHashChain(blockChain core.BlockChain) (*BCPayload, error) {
	hashChain := blockChain.GetHashChain()

	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)

	if err := enc.Encode(hashChain); err != nil {
		return nil, err
	}

	return &BCPayload{
		MsgType: MessageHashChain,
		Payload: buf.Bytes(),
	}, nil
}

func DecodeTransactionMessage(message *BCPayload) (*core.Transaction, error) {
	if message.MsgType != MessageTransaction {
		return nil, fmt.Errorf("invalid message type: Expected (%d) - Found(%d)", MessageTransaction, message.MsgType)
	}

	transaction := core.NewTransaction([]byte{})
	if err := transaction.Decode(bytes.NewBuffer(message.Payload)); err != nil {
		return nil, err
	}

	return transaction, nil
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
