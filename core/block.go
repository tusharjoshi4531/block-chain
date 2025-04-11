package core

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/gob"
	"io"

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

func (header *BlockHeader) Bytes(w io.Writer) error {
	return gob.NewEncoder(w).Encode(header)
}

func (header *BlockHeader) Hash() (types.Hash, error) {
	buf := &bytes.Buffer{}
	if err := header.Bytes(buf); err != nil {
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

func (block *Block) DataHash() (types.Hash, error) {
	buf := &bytes.Buffer{}
	for _, transaction := range block.Transactions {
		if err := transaction.Bytes(buf); err != nil {
			return types.Hash{}, err
		}
	}

	dataHash := sha256.Sum256(buf.Bytes())
	return dataHash, nil
}
