package core

import (
	"encoding/gob"
	"github.com/tusharjoshi4531/block-chain.git/crypto"
	"github.com/tusharjoshi4531/block-chain.git/types"
	"github.com/tusharjoshi4531/block-chain.git/util"
	"io"
)

type SerializableBlock struct {
	Header       BlockHeader
	Transactions []types.Hash
	Validator    *crypto.SerializablePublicKey
	Signature    *crypto.Signature
}

func NewSerializableBlock(block *Block) *SerializableBlock {
	return &SerializableBlock{
		Header:       block.Header,
		Transactions: getTransactionHashes(block.Transactions),
		Validator:    crypto.SerializePublicKey(block.Validator),
		Signature:    block.Signature,
	}
}

func (block *SerializableBlock) Encode(w io.Writer) error {
	return gob.NewEncoder(w).Encode(block)
}

func (block *SerializableBlock) Bytes() ([]byte, error) {
	return util.EncodeToBytes(block)
}

func (block *SerializableBlock) Decode(r io.Reader) error {
	return gob.NewDecoder(r).Decode(block)
}

func (block *SerializableBlock) Reconstruct(txPool TransactionPool) (*Block, error) {
	newBlock := NewBlock()
	newBlock.Header = block.Header
	newBlock.Validator = crypto.DecodePublicKey(block.Validator)
	newBlock.Signature = block.Signature

	for _, txHash := range block.Transactions {
		transaction, err := txPool.GetTransaction(txHash)
		if err != nil {
			return nil, err
		}
		newBlock.AddTransaction(transaction)
	}
	return newBlock, nil
}

func getTransactionHashes(transactions []*Transaction) []types.Hash {
	hashes := make([]types.Hash, 0, len(transactions))
	for _, transaction := range transactions {
		hashes = append(hashes, transaction.Hash())
	}
	return hashes
}
