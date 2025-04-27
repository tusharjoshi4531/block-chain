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
	*prot.SimpleMiner
	*prot.SimpleConsumer
	*prot.SimpleValidator
	*bcnetwork.LocalBlockChainTransport
	blockChain      core.BlockChain
	transactionPool core.TransactionPool
	running         bool
	privKey         *ecdsa.PrivateKey
	mu              sync.RWMutex
}

func NewLocalBlockChainTransport(address string) *LocalBlockChainServer {
	bc := core.NewDefaultBlockChain()
	txPool := core.NewDefaultTransactionPool()
	privKey := crypto.GeneratePrivateKey()
	transport := bcnetwork.NewLocalBlockChainTransport(address, bc, txPool)

	return &LocalBlockChainServer{
		SimpleMiner:              prot.NewSimpleMiner(bc, txPool, privKey),
		SimpleConsumer:           prot.NewSimpleConsumer(bc, txPool, transport),
		SimpleValidator:          prot.NewSimpleValidator(bc, privKey),
		LocalBlockChainTransport: transport,
		blockChain:               bc,
		transactionPool:          txPool,
		privKey:                  privKey,
		running:                  false,
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

			err = server.ReceiveMessage(recPayload)
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
