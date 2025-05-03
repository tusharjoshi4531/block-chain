package server

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"
	"sync"

	bcnetwork "github.com/tusharjoshi4531/block-chain.git/bc_network"
	"github.com/tusharjoshi4531/block-chain.git/core"
	"github.com/tusharjoshi4531/block-chain.git/network"
	"github.com/tusharjoshi4531/block-chain.git/prot"
)

type BlockChainServer interface {
	prot.Miner
	prot.Comsumer
	prot.Validator
	bcnetwork.BlockChainTransport

	Listen()
	ConnectPeer(network.TransportInterface) error
	Kill()
}

type DefaultBlockChainServer struct {
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

func NewDefaultBlockChainServer(
	blockChain core.BlockChain,
	txPool core.TransactionPool,
	privKey *ecdsa.PrivateKey,
	transport bcnetwork.BlockChainTransport,
	minerFactory func() prot.Miner,
	consumerFactory func() prot.Comsumer,
	validatorFactory func() prot.Validator,
) *DefaultBlockChainServer {
	return &DefaultBlockChainServer{
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

func (server *DefaultBlockChainServer) ConnectPeer(transport network.TransportInterface) error {
	return server.Connect(transport)
}

func (server *DefaultBlockChainServer) Listen() {
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

			err = server.ProcessMessage(recPayload, recMsg.From)
			if err != nil {
				fmt.Println("Error: ", err.Error())
			}
		}
	}()
}

func (server *DefaultBlockChainServer) Kill() {
	server.mu.Lock()
	defer server.mu.Unlock()

	server.running = false
}

func (server *DefaultBlockChainServer) PrivKey() *ecdsa.PrivateKey {
	return server.privKey
}
