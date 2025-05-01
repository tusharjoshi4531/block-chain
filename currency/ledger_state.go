package currency

import "fmt"

type LedgerState interface {
	CommitTransaciton(transaction *Transaction) error
	RevertTransaction(transaction *Transaction) error
	HasMember(id string) bool
	AddMember(id string, balance float64) error
	GetBalance(id string) (float64, error)
	GetMembers() []string
}

type MemoryLedgerState struct {
	balance map[string]float64
}

func NewMemoryLedgerState() *MemoryLedgerState {
	return &MemoryLedgerState{
		balance: make(map[string]float64),
	}
}

func (state *MemoryLedgerState) HasMember(id string) bool {
	_, ok := state.balance[id]
	return ok
}

func (state *MemoryLedgerState) CommitTransaciton(transaction *Transaction) error {
	from, to := transaction.From, transaction.To
	if !state.HasMember(from) {
		return fmt.Errorf("no member with id (%s) is present in ledger", from)
	}
	if !state.HasMember(to) {
		return fmt.Errorf("no member with id (%s) is present in ledger", to)
	}

	amt := transaction.Amount
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

func (state *MemoryLedgerState) AddMember(id string, balance float64) error {
	if state.HasMember(id) {
		return fmt.Errorf("member with id (%s) is already present in ledger", id)
	}
	state.balance[id] = balance
	return nil
}

func (state *MemoryLedgerState) GetBalance(id string) (float64, error) {
	balance, ok := state.balance[id]
	if !ok {
		return 0, fmt.Errorf("member with id (%s) is not present in the ledger", id)
	}
	return balance, nil
}

func (state *MemoryLedgerState) GetMembers() []string {
	members := make([]string, 0, len(state.balance))
	for member := range state.balance {
		members = append(members, member)
	}
	return members
}
