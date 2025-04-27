package core

import (
	"fmt"

	"github.com/tusharjoshi4531/block-chain.git/types"
)

type BlockChain interface {
	AddBlock(*Block) error
	GetHeighestBlock() *Block
	GetGenesis() *Block
	GetPrevBlock(*Block) (*Block, error)
	GetBlockWithHash(types.Hash) (*Block, error)
	HasTransactionInChain(transactionHash types.Hash, blockHash types.Hash) error
	GetTransactionsInChain(blockHash types.Hash) ([]*Transaction, error)
	Height() uint32

	GetBlockHashes() []types.Hash
}

type DefaultBlockChain struct {
	height         uint32
	blocks         map[types.Hash]*Block
	blocksAtHeight map[uint32][]*Block
	genesis        *Block
	heighestBlock  *Block
}

func NewDefaultBlockChain() *DefaultBlockChain {
	chain := &DefaultBlockChain{
		height:         0,
		blocks:         make(map[types.Hash]*Block),
		blocksAtHeight: make(map[uint32][]*Block),
		heighestBlock:  nil,
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
		return fmt.Errorf("block (%s) has incorrect height; Required = (%d); Founc = (%d)", blockHash.String(), prevHeight+1, blockHeight)
	}

	blockChain.addBlockWithoutValidation(blockHash, block)

	return nil
}

func (blockChain *DefaultBlockChain) GetHeighestBlock() *Block {
	if blockChain.heighestBlock == nil {
		panic("No heighest block exists in blockchain")
	}
	return blockChain.heighestBlock
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

func (blockChain *DefaultBlockChain) GetTransactionsInChain(tailBlockHash types.Hash) ([]*Transaction, error) {
	currBlockHash := tailBlockHash
	blocks := make([]*Block, 0)
	for {
		currBlock, ok := blockChain.blocks[currBlockHash]
		if !ok {
			return nil, fmt.Errorf("block with hash (%s) is not present in the block chain", tailBlockHash)
		}
		if currBlock.Header.Height == 0 {
			break
		}

		// transactions = append(transactions, currBlock.Transactions...)
		blocks = append(blocks, currBlock)
		currBlockHash = currBlock.Header.PrevBlockHash
	}

	transactions := make([]*Transaction, 0)
	for i := len(blocks) - 1; i >= 0; i-- {
		currBlock := blocks[i]
		transactions = append(transactions, currBlock.Transactions...)
	}
	return transactions, nil
}

func (blockChain *DefaultBlockChain) addBlockWithoutValidation(blockHash types.Hash, block *Block) {
	blockHeight := block.Header.Height

	blockChain.blocks[blockHash] = block
	if _, ok := blockChain.blocksAtHeight[blockHeight]; !ok {
		blockChain.blocksAtHeight = make(map[uint32][]*Block)
	}
	blockChain.blocksAtHeight[blockHeight] = append(blockChain.blocksAtHeight[blockHeight], block)

	if blockHeight >= blockChain.height {
		blockChain.height = blockHeight
		blockChain.heighestBlock = block
	}
}

func (blockChain *DefaultBlockChain) GetBlockHashes() []types.Hash {
	chain := make([]types.Hash, 0, len(blockChain.blocks))
	for hash := range blockChain.blocks {
		// fmt.Printf("HSH: %s, \nBLCK: %v\n\n", hash.String(), block)
		chain = append(chain, hash)
	}
	// fmt.Println("SZ: ", len(chain))
	return chain
}

func (blockChain *DefaultBlockChain) Copy() *DefaultBlockChain {
	newBlocks := make(map[types.Hash]*Block)
	for k, v := range blockChain.blocks {
		newBlocks[k] = v
	}

	newBlocksAtHeight := make(map[uint32][]*Block)
	for k, vs := range blockChain.blocksAtHeight {
		for _, v := range vs {
			newBlocksAtHeight[k] = append(newBlocksAtHeight[k], v)
		}
	}
	return &DefaultBlockChain{
		height:         blockChain.height,
		blocks:         newBlocks,
		blocksAtHeight: newBlocksAtHeight,
		genesis:        blockChain.genesis,
		heighestBlock:  blockChain.heighestBlock,
	}
}
