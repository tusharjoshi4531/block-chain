package shell

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/tusharjoshi4531/block-chain.git/currency"
	"github.com/tusharjoshi4531/block-chain.git/tcp"
)

const (
	ADD_WALLET = "add_wallet"
	WALLETS    = "wallets"
	TRANSACT   = "transact"
	MINE       = "mine"
	BALANCE    = "balance"
)

type ShellInterface struct {
	server *tcp.TCPServer
}

func NewShellInterface(server *tcp.TCPServer) *ShellInterface {
	return &ShellInterface{server: server}
}

func (sh *ShellInterface) Run() {
	sh.server.Listen()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf(">>> ")
		if scanner.Scan() {
			inp := scanner.Text()
			words := strings.Fields(inp)
			cmd, args := words[0], words[1:]
			msg := sh.processCommand(cmd, args)
			fmt.Print(msg)
		}
	}
}

func (sh *ShellInterface) processCommand(cmd string, args []string) string {
	switch cmd {
	case ADD_WALLET:
		if len(args) < 1 {
			return "ERROR: incomplete args\n"
		}
		return sh.processAddWallet(args[0])
	case WALLETS:
		return sh.processWallets()
	case TRANSACT:
		if len(args) < 3 {
			return "ERROR: incomplet arguments"
		}

		from := args[0]
		to := args[1]
		amt, err := strconv.ParseFloat(args[2], 64)
		if err != nil {
			return fmt.Sprintf("ERROR: %s\n", err.Error())
		}

		return sh.processTransact(from, to, amt)
	case MINE:
		if len(args) < 1 {
			return "ERROR: incomplete args\n"
		}
		return sh.processMine(args[0])
	case BALANCE:
		if len(args) < 1 {
			return "ERROR: incomplet arguments"

		}
		walletId := args[0]
		return sh.processBalance(walletId)
	default:
		return fmt.Sprintf("ERROR: invalid command (%s)\n", cmd)
	}
}

func (sh *ShellInterface) processAddWallet(walletId string) string {
	if err := sh.server.AddWallet(walletId); err != nil {
		return fmt.Sprintf("ERROR: %s\n", err.Error())
	} else {
		return fmt.Sprintf("Added wallet (%s) to block chain\n", walletId)
	}
}

func (sh *ShellInterface) processWallets() string {
	wallets := sh.server.Ledger.GetWallets()
	return fmt.Sprintf("Wallets: %v\n", wallets)
}

func (sh *ShellInterface) processTransact(from, to string, amt float64) string {
	transaction := currency.NewTransaction(from, to, amt)
	tx, err := transaction.ToCoreTransaction()
	if err != nil {
		return fmt.Sprintf("ERROR: %s\n", err.Error())
	}

	if err := tx.Sign(sh.server.PrivKey); err != nil {
		return fmt.Sprintf("ERROR: %s\n", err.Error())

	}

	if err := sh.server.AddTransaction(tx); err != nil {
		return fmt.Sprintf("ERROR: %s\n", err.Error())
	}
	return "Transaction Added\n"
}

func (sh *ShellInterface) processMine(minerWalletId string) string {
	block, err := sh.server.MineBlock(10, minerWalletId)
	if err != nil {
		return fmt.Sprintf("ERROR: %s\n", err.Error())
	}

	err = block.Sign(sh.server.PrivKey)
	if err != nil {
		return fmt.Sprintf("ERROR: %s\n", err.Error())
	}

	err = sh.server.BlockChain.AddBlock(block)
	if err != nil {
		return fmt.Sprintf("ERROR: %s\n", err.Error())
	}

	return "New block created\n"
}

func (sh *ShellInterface) processBalance(walletId string) string {
	balance, err := sh.server.Ledger.GetBalance(walletId)
	if err != nil {
		return fmt.Sprintf("ERROR: %s\n", err.Error())

	}
	return fmt.Sprintf("Wallet (%s) : %f\n", walletId, balance)
}
