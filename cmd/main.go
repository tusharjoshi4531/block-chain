package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"github.com/tusharjoshi4531/block-chain.git/currency"
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
				if len(words) < 2 {
					fmt.Println("ERROR: incomplet arguments")
					continue
				}
				walletId := words[1]

				if err := server.AddWallet(walletId); err != nil {
					fmt.Printf("ERROR: %s\n", err.Error())
				} else {
					fmt.Printf("Added wallet (%s) to block chain\n", walletId)
				}
			} else if words[0] == "wallets" {
				wallets := server.Ledger.GetWallets()
				fmt.Printf("Wallets: %v\n", wallets)
			} else if words[0] == "transact" {
				if len(words) < 4 {
					fmt.Println("ERROR: incomplet arguments")
					continue
				}

				from := words[1]
				to := words[2]
				amt, err := strconv.ParseFloat(words[3], 64)
				if err != nil {
					fmt.Printf("ERROR: %s\n", err.Error())
				}

				transaction := currency.NewTransaction(from, to, amt)
				tx, err := transaction.ToCoreTransaction()
				if err != nil {
					fmt.Printf("ERROR: %s\n", err.Error())
					continue
				}

				if err := tx.Sign(server.PrivKey); err != nil {
					fmt.Printf("ERROR: %s\n", err.Error())
					continue
				}

				if err := server.AddTransaction(tx); err != nil {
					fmt.Printf("ERROR: %s\n", err.Error())
					continue
				}
			} else if words[0] == "mine" {
				block, err := server.MineBlock(10)
				if err != nil {
					fmt.Printf("ERROR: %s\n", err.Error())
					continue
				}

				err = block.Sign(server.PrivKey)
				if err != nil {
					fmt.Printf("ERROR: %s\n", err.Error())
					continue
				}

				err = server.BlockChain.AddBlock(block)
				if err != nil {
					fmt.Printf("ERROR: %s\n", err.Error())
					continue
				}

				fmt.Println("new block created")
			} else if words[0] == "balance" {
				if len(words) < 2 {
					fmt.Println("ERROR: incomplet arguments")
					continue
				}
				walletId := words[1]

				balance, err := server.Ledger.GetBalance(walletId)
				if err != nil {
					fmt.Printf("ERROR: %s\n", err.Error())
					continue
				}
				fmt.Printf("Wallet (%s) : %f\n", walletId, balance)
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
