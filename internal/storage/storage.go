package storage

import "errors"

type Storage interface {
	CreateAccount(username string)
}

var ErrUserExists = errors.New("user already exists")
