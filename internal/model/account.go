package model

import (
	"errors"
	"sync"
)

type Account struct {
	Username string
	mu       sync.Mutex
	balance  float64
}

func (a *Account) Deposit(amount float64) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.balance += amount
	return nil
}

func (a *Account) Withdraw(amount float64) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.balance < amount {
		return errors.New("insufficient funds")
	}
	a.balance -= amount
	return nil
}

func (a *Account) GetBalance() float64 {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.balance
}
