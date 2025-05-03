package main

import (
	"fmt"
	"log"
	"os"

	bcnetwork "github.com/tusharjoshi4531/block-chain.git/bc_network"
	"github.com/tusharjoshi4531/block-chain.git/core"
	"github.com/tusharjoshi4531/block-chain.git/crypto"
	"github.com/tusharjoshi4531/block-chain.git/currency"
	"github.com/tusharjoshi4531/block-chain.git/network"
	"github.com/tusharjoshi4531/block-chain.git/shell"
	"github.com/tusharjoshi4531/block-chain.git/tcp"
)

func main() {
	fmt.Println(os.Args)
	addr, peers := parseArgs(os.Args)

	fmt.Println(addr)

	ledger := currency.NewMemoryLedgerState()
	bc := currency.NewBlockChain(ledger, 10000)
	txPool := core.NewDefaultTransactionPool()
	privKey := crypto.GeneratePrivateKey()
	bcTransport := bcnetwork.NewDefaultBlockChainTransport(
		network.NewDefaultTransport(addr),
		bc,
		txPool,
	)

	server := tcp.NewTcpServer(
		ledger,
		bc,
		txPool,
		privKey,
		bcTransport,
	)

	for _, peer := range peers {
		if err := server.Connect(tcp.NewTcpTransportInterface(peer)); err != nil {
			log.Fatalf("Couldn't connect (%s) to peer (%s), ERROR: (%s)", addr, peer, err.Error())
		}
	}

	sh := shell.NewShellInterface(server)

	sh.Run()
}

func parseArgs(args []string) (string, []string) {
	if len(args) < 1 {
		panic("port not defined")
	}
	addr := args[1]
	peers := args[2:]
	return addr, peers
}
