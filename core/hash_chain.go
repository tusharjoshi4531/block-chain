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

func NewHashChainFromBlocks(blocks []*Block) (*HashChain, error){
	hashes := make([]types.Hash, 0, len(blocks))
	for _, block := range blocks {
		hsh, err := block.Hash()
		if err != nil {
			return nil, err
		}
		hashes = append(hashes, hsh)
	}
	return NewHashChainFromHashes(hashes), nil
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

func (hahsChain *HashChain) GetExcludedBlocks(blockChain BlockChain) []*Block {
	excludedBlocksHash := hahsChain.GetExcludedBlockHashes(blockChain)
	excludedBlocks := make([]*Block, 0, len(excludedBlocksHash))

	for _, hash := range excludedBlocksHash {
		block, err := blockChain.GetBlockWithHash(hash)
		
		if err != nil {
			panic("incorrect behavior while getting excluded blocks from hashchain")
		}

		excludedBlocks = append(excludedBlocks, block)
	}

	return excludedBlocks
}

func (hashChain *HashChain) GetBlockHashes() []types.Hash {
	hashes := make([]types.Hash, 0, len(hashChain.BlockHashes))
	for hash := range hashChain.BlockHashes {
		hashes = append(hashes, hash)
	}
	return hashes
}
