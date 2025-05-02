package currency

import "fmt"

const RewardSymbol = "::"

type LedgerState interface {
	CommitTransaciton(transaction *Transaction) error
	RevertTransaction(transaction *Transaction) error
	HasWallet(id string) bool
	AddWallet(id string, balance float64) error
	GetBalance(id string) (float64, error)
	GetWallets() []string
}

type MemoryLedgerState struct {
	balance map[string]float64
}

func NewMemoryLedgerState() *MemoryLedgerState {
	return &MemoryLedgerState{
		balance: make(map[string]float64),
	}
}

func (state *MemoryLedgerState) HasWallet(id string) bool {
	_, ok := state.balance[id]
	return ok
}

func (state *MemoryLedgerState) CommitTransaciton(transaction *Transaction) error {
	from, to := transaction.From, transaction.To
	amt := transaction.Amount

	if from == RewardSymbol {
		state.balance[to] += amt
		return nil
	}

	if to == RewardSymbol {
		state.balance[from] -= amt
		return nil
	}

	if !state.HasWallet(to) {
		return fmt.Errorf("no member with id (%s) is present in ledger", to)
	}

	if !state.HasWallet(from) {
		return fmt.Errorf("no member with id (%s) is present in ledger", from)
	}

	fromAmt := state.balance[from]
	if fromAmt < amt {
		return fmt.Errorf("sender (%s) does not have enough balance", from)
	}

	state.balance[from] -= amt
	state.balance[to] += amt

	return nil
}

func (state *MemoryLedgerState) RevertTransaction(transaction *Transaction) error {
	revTransaction := NewTransaction(transaction.To, transaction.From, transaction.Amount)
	return state.CommitTransaciton(revTransaction)
}

func (state *MemoryLedgerState) AddWallet(walletId string, balance float64) error {
	if state.HasWallet(walletId) {
		return fmt.Errorf("member with id (%s) is already present in ledger", walletId)
	}
	state.balance[walletId] = balance
	return nil
}

func (state *MemoryLedgerState) GetBalance(id string) (float64, error) {
	balance, ok := state.balance[id]
	if !ok {
		return 0, fmt.Errorf("member with id (%s) is not present in the ledger", id)
	}
	return balance, nil
}

func (state *MemoryLedgerState) GetWallets() []string {
	members := make([]string, 0, len(state.balance))
	for member := range state.balance {
		members = append(members, member)
	}
	return members
}
