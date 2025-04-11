package core

import "testing"

func TestDefaultBlockChain(t *testing.T) {

}

func createBlock(t *testing.T, height uint32, txx []*Transaction) {
	block := NewBlock()
	block.Header.Height = height
	
}