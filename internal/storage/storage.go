package storage

import "errors"

type Storage interface {
	CreateAccount(username string) (int, error)
}

var ErrUserExists = errors.New("user already exists")
