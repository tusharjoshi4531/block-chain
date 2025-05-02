package bcnetwork

import (
	"github.com/tusharjoshi4531/block-chain.git/core"
	"github.com/tusharjoshi4531/block-chain.git/network"
)

type LocalBlockChainTransport struct {
	*DefaultBlockChainTransport
	network.TransportInterface
}

func NewLocalBlockChainTransport(address string, blockChain core.BlockChain, transactionPool core.TransactionPool) *LocalBlockChainTransport {
	transport := network.NewLocalTransport(address)
	return &LocalBlockChainTransport{
		DefaultBlockChainTransport: NewDefaultBlockChainTransport(transport, blockChain, transactionPool),
		TransportInterface:         transport,
	}
}
