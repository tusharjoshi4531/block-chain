package core

import (
	"bytes"
	"encoding/gob"
	"io"

	"github.com/tusharjoshi4531/block-chain.git/types"
)

type HashChain struct {
	BlockHashes map[types.Hash]bool
}

func NewHashChain(blockChain BlockChain) *HashChain {
	hashes := blockChain.GetBlockHashes()
	return NewHashChainFromHashes(hashes)
}

func NewHashChainFromHashes(hashes []types.Hash) *HashChain {
	hashChain := &HashChain{
		BlockHashes: make(map[types.Hash]bool),
	}
	for _, hash := range hashes {
		hashChain.BlockHashes[hash] = true
	}

	return hashChain
}

func (hashChain *HashChain) Encode(w io.Writer) error {
	return gob.NewEncoder(w).Encode(hashChain.BlockHashes)
}

func (hashChain *HashChain) Bytes() ([]byte, error) {
	buf := &bytes.Buffer{}
	if err := hashChain.Encode(buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (hashChain *HashChain) Decode(r io.Reader) error {
	return gob.NewDecoder(r).Decode(&hashChain.BlockHashes)
}

func (hashChain *HashChain) GetExcludedBlockHashes(blockChain BlockChain) []types.Hash {
	excludedBlockHashes := make([]types.Hash, 0)

	hashes := blockChain.GetBlockHashes()
	for _, hash := range hashes {
		_, ok := hashChain.BlockHashes[hash]
		if !ok {
			excludedBlockHashes = append(excludedBlockHashes, hash)
		}
	}
	return excludedBlockHashes
}

func (hashChain *HashChain) GetBlockHashes() []types.Hash {
	hashes := make([]types.Hash, 0, len(hashChain.BlockHashes))
	for hash := range hashChain.BlockHashes {
		hashes = append(hashes, hash)
	}
	return hashes
}
