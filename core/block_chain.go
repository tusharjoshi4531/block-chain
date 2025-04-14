package core

import (
	"fmt"
	"math/rand"

	"github.com/tusharjoshi4531/block-chain.git/types"
)

type BlockChain interface {
	AddBlock(*Block) error
	GetHeighestBlock() *Block
	GetGenesis() *Block
	GetPrevBlock(*Block) (*Block, error)
	GetBlockWithHash(types.Hash) (*Block, error)
	HasTransactionInChain(transactionHash types.Hash, blockHash types.Hash) error
	Height() uint32
}

type DefaultBlockChain struct {
	height         uint32
	blocks         map[types.Hash]*Block
	blocksAtHeight map[uint32][]*Block
	genesis        *Block
}

func NewDefaultBlockChain() *DefaultBlockChain {
	chain := &DefaultBlockChain{
		height:         0,
		blocks:         make(map[types.Hash]*Block),
		blocksAtHeight: make(map[uint32][]*Block),
	}
	// TODO: Add genesis block
	genesisBlock := NewBlock()
	genesisBlock.Header.Height = 0
	hash, err := genesisBlock.Hash()
	if err != nil {
		panic(err)
	}

	chain.addBlockWithoutValidation(hash, genesisBlock)
	chain.genesis = genesisBlock

	return chain
}

func (blockChain *DefaultBlockChain) GetGenesis() *Block {
	return blockChain.genesis
}

func (blockChain *DefaultBlockChain) AddBlock(block *Block) error {
	prevBlock, ok := blockChain.blocks[block.Header.PrevBlockHash]
	if !ok {
		return fmt.Errorf("previous block of hash (%s) doesnot exist", block.Header.PrevBlockHash)
	}
	blockHeight := block.Header.Height
	prevHeight := prevBlock.Header.Height

	blockHash, err := block.Hash()
	if err != nil {
		return err
	}

	if prevBlock.Header.Height != block.Header.Height-1 {
		return fmt.Errorf("block (%s) has incorrect size; Required = (%d); Founc = (%d)", blockHash, prevHeight+1, blockHeight)
	}

	blockChain.addBlockWithoutValidation(blockHash, block)

	return nil
}

func (blockChain *DefaultBlockChain) GetHeighestBlock() *Block {
	blocks, ok := blockChain.blocksAtHeight[blockChain.height]
	if !ok {
		panic("block chain height is longer than existing blocks heights")
	}

	numBlocks := len(blocks)
	idx := rand.Int31n(int32(numBlocks))

	return blocks[idx]
}

func (blockChain *DefaultBlockChain) GetPrevBlock(block *Block) (*Block, error) {
	return blockChain.GetBlockWithHash(block.Header.PrevBlockHash)
}

func (blockChain *DefaultBlockChain) GetBlockWithHash(hash types.Hash) (*Block, error) {
	block, ok := blockChain.blocks[hash]
	if !ok {
		return nil, fmt.Errorf("couldnot find block with hash (%s)", hash)
	}
	return block, nil
}

func (blockChain *DefaultBlockChain) Height() uint32 {
	return blockChain.height
}

func (blockChain *DefaultBlockChain) HasTransactionInChain(transactionHash types.Hash, tailBlockHash types.Hash) error {
	currBlockHash := tailBlockHash
	for {
		currBlock, ok := blockChain.blocks[currBlockHash]
		if !ok {
			return fmt.Errorf("block with hash (%s) is not present in the block chain", tailBlockHash)
		}
		if currBlock.Header.Height == 0 {
			break
		}

		if currBlock.HasTranaction(transactionHash) {
			return nil
		}
		currBlockHash = currBlock.Header.PrevBlockHash
	}
	return fmt.Errorf("transaction with hash (%s) is not present in the block chain", transactionHash)
}

func (blockChain *DefaultBlockChain) addBlockWithoutValidation(blockHash types.Hash, block *Block) {
	blockHeight := block.Header.Height

	blockChain.blocks[blockHash] = block
	if _, ok := blockChain.blocksAtHeight[blockHeight]; !ok {
		blockChain.blocksAtHeight = make(map[uint32][]*Block)
	}
	blockChain.blocksAtHeight[blockHeight] = append(blockChain.blocksAtHeight[blockHeight], block)
	blockChain.height = max(blockChain.height, blockHeight)
}
