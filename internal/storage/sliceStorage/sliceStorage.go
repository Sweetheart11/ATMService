package sliceStorage

import (
	"fmt"

	"github.com/Sweetheart11/ATMService/internal/model"
	"github.com/Sweetheart11/ATMService/internal/storage"
)

type Storage []BankAccount

type BankAccount interface {
	Deposit(amount float64) error
	Withdraw(amount float64) error
	GetBalance() float64
}

func New() (Storage, error) {
	var storage Storage
	return storage, nil
}

func (s *Storage) CreateAccount(username string) (int, error) {
	const op = "storage.sliceStorage.CreateAccount"

	for _, bankAcc := range *s {
		acc := bankAcc.(*model.Account)
		if acc.Username == username {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}
	}

	*s = append(*s, &model.Account{
		Username: username,
	})

	return len(*s) - 1, nil
}
