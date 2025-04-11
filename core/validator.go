package core

import "fmt"

type Validator interface {
	ValidateBlock(block *Block) error
	ValidateTransactions(transactions []*Transaction) error
}

type DefaultValidator struct{}

func (DefaultValidator) ValidateBlock(block *Block) error {
	// Verify DataHash
	dataHash, err := block.DataHash()
	if err != nil {
		return err
	}


	if dataHash != block.Header.DataHash {
		return fmt.Errorf("data hash of block does not match transactions")
	}

	// Verify Block Hash
	headerHash, err := block.Header.Hash()
	if err != nil {
		return err
	}
	blockHash, err := block.Hash()
	if err != nil {
		return err
	}

	if headerHash != blockHash {
		return fmt.Errorf("block hash does not match the block")
	}
	
	return nil
}

func (DefaultValidator) ValidateTransactions(transactions []*Transaction) error {
	for _, transaction := range transactions {
		if err := transaction.Verify(); err != nil {
			return err
		}
	}
	return nil
}