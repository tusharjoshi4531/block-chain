package server

import (
	"crypto/ecdsa"

	bcnetwork "github.com/tusharjoshi4531/block-chain.git/bc_network"
	"github.com/tusharjoshi4531/block-chain.git/core"
	"github.com/tusharjoshi4531/block-chain.git/crypto"
	"github.com/tusharjoshi4531/block-chain.git/network"
	"github.com/tusharjoshi4531/block-chain.git/prot"
)

type LocalBlockChainServer struct {
	*DefaultBlockChainServer
	network.TransportInterface
}

func NewSimpleLocalBlockChainServer(address string) *LocalBlockChainServer {
	bc := core.NewDefaultBlockChain()
	txPool := core.NewDefaultTransactionPool()
	privKey := crypto.GeneratePrivateKey()
	transport := bcnetwork.NewLocalBlockChainTransport(address, bc, txPool)

	return &LocalBlockChainServer{
		DefaultBlockChainServer: NewDefaultBlockChainServer(
			bc,
			txPool,
			privKey,
			transport,
			func() prot.Miner { return prot.NewSimpleMiner(bc, txPool, privKey) },
			func() prot.Comsumer { return prot.NewSimpleConsumer(bc, txPool, transport) },
			func() prot.Validator { return prot.NewSimpleValidator(bc, privKey) },
		),
		TransportInterface: transport,
	}
}

func NewLocalBlockChainServer(
	blockChain core.BlockChain,
	txPool core.TransactionPool,
	privKey *ecdsa.PrivateKey,
	transport *bcnetwork.LocalBlockChainTransport,
	minerFactory func() prot.Miner,
	consumerFactory func() prot.Comsumer,
	validatorFactory func() prot.Validator,
) *LocalBlockChainServer {
	return &LocalBlockChainServer{
		DefaultBlockChainServer: NewDefaultBlockChainServer(
			blockChain,
			txPool,
			privKey,
			transport,
			minerFactory,
			consumerFactory,
			validatorFactory,
		),
		TransportInterface: transport,
	}
}
