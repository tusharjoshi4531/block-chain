package currency

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommitTransaction(t *testing.T) {
	state := NewMemoryLedgerState()

	assert.Nil(t, state.AddWallet("A", 100))
	assert.Nil(t, state.AddWallet("B", 100))

	tx1 := NewTransaction(RewardSymbol, "A", 50)
	tx2 := NewTransaction("A", "B", 13)

	assert.Nil(t, state.CommitTransaciton(tx1))
	assert.Nil(t, state.CommitTransaciton(tx2))

	balance, err := state.GetBalance("A")
	assert.Nil(t, err)
	fmt.Println(balance)
	assert.Equal(t, balance, float64(100+50-13))

	balance, err = state.GetBalance("B")
	assert.Nil(t, err)
	assert.Equal(t, balance, float64(100+13))
}

func TestRevertTransaction(t *testing.T) {
	state := NewMemoryLedgerState()

	assert.Nil(t, state.AddWallet("A", 100))
	assert.Nil(t, state.AddWallet("B", 100))

	tx1 := NewTransaction(RewardSymbol, "A", 50)
	tx2 := NewTransaction("A", "B", 13)

	assert.Nil(t, state.CommitTransaciton(tx1))
	assert.Nil(t, state.CommitTransaciton(tx2))

	balance, err := state.GetBalance("A")
	assert.Nil(t, err)
	fmt.Println(balance)
	assert.Equal(t, balance, float64(100+50-13))

	balance, err = state.GetBalance("B")
	assert.Nil(t, err)
	assert.Equal(t, balance, float64(100+13))

	assert.Nil(t, state.RevertTransaction(tx2))

	balance, err = state.GetBalance("A")
	assert.Nil(t, err)
	fmt.Println(balance)
	assert.Equal(t, balance, float64(100+50))

	balance, err = state.GetBalance("B")
	assert.Nil(t, err)
	assert.Equal(t, balance, float64(100))
}

func TestInvalidTransaction(t *testing.T) {
	state := NewMemoryLedgerState()

	assert.Nil(t, state.AddWallet("A", 100))
	assert.Nil(t, state.AddWallet("B", 100))

	tx1 := NewTransaction(RewardSymbol, "A", 50)
	tx2 := NewTransaction("A", "B", 200)

	assert.Nil(t, state.CommitTransaciton(tx1))
	assert.NotNil(t, state.CommitTransaciton(tx2))

	balance, err := state.GetBalance("A")
	assert.Nil(t, err)
	fmt.Println(balance)
	assert.Equal(t, balance, float64(100+50))

	balance, err = state.GetBalance("B")
	assert.Nil(t, err)
	assert.Equal(t, balance, float64(100))
}

func TestAddWallet(t *testing.T) {
	state := NewMemoryLedgerState()

	assert.Nil(t, state.AddWallet("A", 100))
	assert.Nil(t, state.AddWallet("B", 100))

	tx1 := NewTransaction(RewardSymbol, "A", 50)
	tx2 := NewTransaction("A", "B", 200)

	assert.Nil(t, state.CommitTransaciton(tx1))
	assert.NotNil(t, state.CommitTransaciton(tx2))

	balance, err := state.GetBalance("A")
	assert.Nil(t, err)
	fmt.Println(balance)
	assert.Equal(t, balance, float64(100+50))

	balance, err = state.GetBalance("B")
	assert.Nil(t, err)
	assert.Equal(t, balance, float64(100))

	assert.Nil(t, state.AddWallet("C", 1000))
	tx3 := NewTransaction("C", "A", 500)

	assert.Nil(t, state.CommitTransaciton(tx3))

	balance, err = state.GetBalance("A")
	assert.Nil(t, err)
	fmt.Println(balance)
	assert.Equal(t, balance, float64(100+50+500))

	balance, err = state.GetBalance("B")
	assert.Nil(t, err)
	assert.Equal(t, balance, float64(100))

	balance, err = state.GetBalance("C")
	assert.Nil(t, err)
	assert.Equal(t, balance, float64(1000-500))

	assert.Equal(t, len(state.GetWallets()), 3)
}
