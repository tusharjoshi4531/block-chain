package server

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"
	"sync"

	bcnetwork "github.com/tusharjoshi4531/block-chain.git/bc_network"
	"github.com/tusharjoshi4531/block-chain.git/core"
	"github.com/tusharjoshi4531/block-chain.git/crypto"
	"github.com/tusharjoshi4531/block-chain.git/network"
	"github.com/tusharjoshi4531/block-chain.git/prot"
)

type LocalBlockChainServer struct {
	prot.Miner
	prot.Comsumer
	prot.Validator
	bcnetwork.BlockChainTransport
	blockChain      core.BlockChain
	transactionPool core.TransactionPool
	running         bool
	privKey         *ecdsa.PrivateKey
	mu              sync.RWMutex
}

func NewSimpleLocalBlockChainServer(address string) *LocalBlockChainServer {
	bc := core.NewDefaultBlockChain()
	txPool := core.NewDefaultTransactionPool()
	privKey := crypto.GeneratePrivateKey()
	transport := bcnetwork.NewLocalBlockChainTransport(address, bc, txPool)

	return NewLocalBlockChainServer(
		bc,
		txPool,
		privKey,
		transport,
		func() prot.Miner { return prot.NewSimpleMiner(bc, txPool, privKey) },
		func() prot.Comsumer { return prot.NewSimpleConsumer(bc, txPool, transport) },
		func() prot.Validator { return prot.NewSimpleValidator(bc, privKey) },
	)
}

func NewLocalBlockChainServer(
	blockChain core.BlockChain,
	txPool core.TransactionPool,
	privKey *ecdsa.PrivateKey,
	transport bcnetwork.BlockChainTransport,
	minerFactory func() prot.Miner,
	consumerFactory func() prot.Comsumer,
	validatorFactory func() prot.Validator,
) *LocalBlockChainServer {
	return &LocalBlockChainServer{
		Miner:               minerFactory(),
		Comsumer:            consumerFactory(),
		Validator:           validatorFactory(),
		BlockChainTransport: transport,
		blockChain:          blockChain,
		transactionPool:     txPool,
		privKey:             privKey,
		running:             false,
	}
}

func (server *LocalBlockChainServer) ConnectPeer(transport network.Transport) error {
	return server.Connect(transport)
}

func (server *LocalBlockChainServer) Listen() {
	server.running = true

	go func() {
		for {
			server.mu.RLock()
			if !server.running {
				break
			}
			server.mu.RUnlock()

			recMsg := <-server.ReadChan()
			recPayload := &bcnetwork.BCPayload{}
			err := recPayload.Decode(bytes.NewBuffer(recMsg.Payload))

			if err != nil {
				fmt.Println("couldn't decode message")
				continue
			}

			err = server.ReceiveMessage(recPayload, recMsg.From)
			if err != nil {
				fmt.Println("Error: ", err.Error())
			}
		}
	}()
}

func (server *LocalBlockChainServer) Kill() {
	server.mu.Lock()
	defer server.mu.Unlock()

	server.running = false
}

func (server *LocalBlockChainServer) PrivKey() *ecdsa.PrivateKey {
	return server.privKey
}
