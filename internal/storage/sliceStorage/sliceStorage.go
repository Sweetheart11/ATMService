package sliceStorage

import (
	"github.com/Sweetheart11/ATMService/internal/models"
	"github.com/Sweetheart11/ATMService/internal/storage"
)

type SliceStorage []BankAccount

type BankAccount interface {
	Deposit(amount float64) error
	Withdraw(amount float64) error
	GetBalance() float64
}

func New() (SliceStorage, error) {
	var sliceStorage SliceStorage
	return sliceStorage, nil
}

func (s *SliceStorage) CreateAccount(username string) error {
	for _, bankAcc := range *s {
		acc := bankAcc.(*models.Account)
		if acc.Username == username {
			return storage.ErrUserExists
		}
	}

	*s = append(*s, &models.Account{
		Username: username,
	})

	return nil
}
