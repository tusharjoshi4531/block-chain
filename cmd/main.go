package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/tusharjoshi4531/block-chain.git/tcp"
)

func main() {
	fmt.Println(os.Args)
	addr, peers := parseArgs(os.Args)

	fmt.Println(addr)

	server := tcp.NewTcpServer(addr)

	for _, peer := range peers {
		if err := server.Connect(tcp.NewTcpTransportInterface(peer)); err != nil {
			log.Fatalf("Couldn't connect (%s) to peer (%s), ERROR: (%s)", addr, peer, err.Error())
		}
	}

	server.Listen()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf(">>> ")
		if scanner.Scan() {
			cmd := scanner.Text()
			words := strings.Fields(cmd)

			if words[0] == "add_wallet" {
				walletId := words[1]

				if err := server.AddWallet(walletId); err != nil {
					fmt.Printf("ERROR: %s\n", err.Error())
				} else {
					fmt.Printf("Added wallet (%s) to block chain\n", walletId)
				}
			} else if words[0] == "print_wallets" {
				wallets := server.Ledger.GetWallets()
				fmt.Printf("Wallets: %v\n", wallets)
			}
		}
	}
}

func parseArgs(args []string) (string, []string) {
	if len(args) < 1 {
		panic("port not defined")
	}
	addr := args[1]
	peers := args[2:]
	return addr, peers
}
