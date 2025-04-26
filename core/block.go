package core

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"io"
	"time"

	"github.com/tusharjoshi4531/block-chain.git/crypto"
	"github.com/tusharjoshi4531/block-chain.git/types"
)

type BlockHeader struct {
	Version       uint32
	DataHash      types.Hash
	PrevBlockHash types.Hash
	Timestamp     int64
	Height        uint32
}

func (header *BlockHeader) Encode(w io.Writer) error {
	return gob.NewEncoder(w).Encode(header)
}

func (header *BlockHeader) Bytes() ([]byte, error) {
	buf := &bytes.Buffer{}
	if err := header.Encode(buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (header *BlockHeader) Decode(r io.Reader) error {
	return gob.NewDecoder(r).Decode(header)
}

func (header *BlockHeader) Hash() (types.Hash, error) {
	buf := &bytes.Buffer{}
	if err := header.Encode(buf); err != nil {
		return types.Hash{}, err
	}
	return sha256.Sum256(buf.Bytes()), nil
}

type Block struct {
	Header       BlockHeader
	Transactions []*Transaction
	Validator    *ecdsa.PublicKey
	Signature    *crypto.Signature

	hash types.Hash
}

func NewBlock() *Block {
	return &Block{
		Header: BlockHeader{
			Version:       0,
			DataHash:      types.Hash{},
			PrevBlockHash: types.Hash{},
			Timestamp:     time.Now().UnixNano(),
			Height:        0,
		},
		Transactions: []*Transaction{},
		Validator:    &ecdsa.PublicKey{},
		Signature:    &crypto.Signature{},
	}
}

func NewBlockWithHeaderInfo(height uint32, prevBlockHash types.Hash) *Block {
	block := NewBlock()
	block.Header.Height = height
	block.Header.PrevBlockHash = prevBlockHash
	return block
}

func (block *Block) Hash() (types.Hash, error) {
	if !block.hash.IsZero() {
		return block.hash, nil
	}

	// Hash data
	dataHash, err := block.DataHash()
	if err != nil {
		return types.Hash{}, err
	}
	block.Header.DataHash = dataHash

	// Hash header
	hash, err := block.Header.Hash()
	if err != nil {
		return types.Hash{}, err
	}

	block.hash = hash
	return block.hash, nil
}

func (block *Block) EncodeData(w io.Writer) error {
	len := len(block.Transactions)
	if err := gob.NewEncoder(w).Encode(len); err != nil {
		return err
	}

	for _, transaction := range block.Transactions {
		if err := transaction.Encode(w); err != nil {
			return err
		}
	}
	return nil
}

func (block *Block) DecodeData(r io.Reader) error {
	len := 0
	if err := gob.NewDecoder(r).Decode(&len); err != nil {
		return err
	}

	block.Transactions = make([]*Transaction, 0, len)
	for i := 0; i < len; i++ {
		transaction := NewTransaction([]byte{})
		if err := transaction.Decode(r); err != nil {
			return err
		}
		block.Transactions = append(block.Transactions, transaction)
	}
	return nil
}

func (block *Block) DataHash() (types.Hash, error) {
	buf := &bytes.Buffer{}
	if err := block.EncodeData(buf); err != nil {
		return types.Hash{}, err
	}

	dataHash := sha256.Sum256(buf.Bytes())
	return dataHash, nil
}

func (block *Block) AddTransaction(transaction *Transaction) {
	block.Transactions = append(block.Transactions, transaction)
}

func (block *Block) Sign(privateKey *ecdsa.PrivateKey) error {
	// Hash block
	block.Hash()

	headerBytes, err := block.Header.Bytes()
	if err != nil {
		return err
	}

	sig, err := crypto.SignBytes(privateKey, headerBytes)
	if err != nil {
		return err
	}

	block.Signature = sig
	block.Validator = &privateKey.PublicKey

	return nil
}

func (block *Block) Verify() error {
	if block.Signature == nil {
		return fmt.Errorf("block has no signature")
	}

	headerBytes, err := block.Header.Bytes()
	if err != nil {
		return err
	}

	if !block.Signature.Verify(block.Validator, headerBytes) {
		return fmt.Errorf("incorrect sign in transaction")
	}

	return nil
}

func (block *Block) Encode(w io.Writer) error {
	if err := block.Header.Encode(w); err != nil {
		return err
	}
	if err := block.EncodeData(w); err != nil {
		return err
	}
	if err := crypto.SerializePublicKey(block.Validator).Encode(w); err != nil {
		return err
	}
	if err := block.Signature.Encode(w); err != nil {
		return err
	}

	return nil
}

func (block *Block) Bytes() ([]byte, error) {
	buf := &bytes.Buffer{}
	if err := block.Encode(buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (block *Block) Decode(r io.Reader) error {
	if err := block.Header.Decode(r); err != nil {
		return err
	}

	if err := block.DecodeData(r); err != nil {
		return err
	}

	serializedValidator := &crypto.SerializedPublicKey{}
	if err := serializedValidator.Decode(r); err != nil {
		return err
	}
	block.Validator = crypto.DecodePublicKey(serializedValidator)

	if err := block.Signature.Decode(r); err != nil {
		return err
	}
	return nil
}

func (block *Block) HasTranaction(hash types.Hash) bool {
	for _, transaction := range block.Transactions {
		if transaction.Hash() == hash {
			return true
		}
	}
	return false
}
