package tcp

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"
	"io"
	"log"
	"net"

	bcnetwork "github.com/tusharjoshi4531/block-chain.git/bc_network"
	"github.com/tusharjoshi4531/block-chain.git/core"
	"github.com/tusharjoshi4531/block-chain.git/crypto"
	"github.com/tusharjoshi4531/block-chain.git/currency"
	"github.com/tusharjoshi4531/block-chain.git/network"
	"github.com/tusharjoshi4531/block-chain.git/pow"
	"github.com/tusharjoshi4531/block-chain.git/prot"
	"github.com/tusharjoshi4531/block-chain.git/server"
)

type TCPServer struct {
	*server.DefaultBlockChainServer
	BlockChain core.BlockChain
	Ledger     currency.LedgerState
	PrivKey    *ecdsa.PrivateKey
}

func NewTcpServer(address string) *TCPServer {
	ledger := currency.NewMemoryLedgerState()
	bc := currency.NewBlockChain(ledger, 10000)
	txPool := core.NewDefaultTransactionPool()
	privKey := crypto.GeneratePrivateKey()
	bcTransport := bcnetwork.NewDefaultBlockChainTransport(
		network.NewDefaultTransport(address),
		bc,
		txPool,
	)
	return &TCPServer{
		BlockChain: bc,
		PrivKey:    privKey,
		DefaultBlockChainServer: server.NewDefaultBlockChainServer(
			bc,
			txPool,
			privKey,
			bcTransport,
			func() prot.Miner {
				return pow.NewPowMiner(
					1,
					bc,
					txPool, privKey,
					address,
					currency.NewRewarder(privKey, bc, 100, 10),
				)
			},
			func() prot.Comsumer {
				return prot.NewSimpleConsumer(
					bc,
					txPool,
					bcTransport,
				)
			},
			func() prot.Validator {
				return prot.NewSimpleValidator(
					bc,
					privKey,
				)
			},
		),
		Ledger: ledger,
	}
}

func (server *TCPServer) Listen() {
	listener, err := net.Listen("tcp", server.Address())
	if err != nil {
		log.Fatal(err)
	}

	server.DefaultBlockChainServer.Listen()

	fmt.Printf("Server running at -> %s\n", server.Address())
	go func() {
		defer listener.Close()
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Println(err)
				continue
			}

			go server.handleConnection(conn.(*net.TCPConn))
		}
	}()
}

func (server *TCPServer) handleConnection(conn *net.TCPConn) error {
	defer conn.Close()

	buf := &bytes.Buffer{}
	if _, err := io.Copy(buf, conn); err != nil {
		return err
	}

	msg := &network.Message{}
	if err := msg.Decode(buf); err != nil {
		return err
	}

	server.WriteChan() <- *msg
	return nil
}
