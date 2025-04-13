package entity

import "github.com/google/uuid"

type Account struct {
	ID uuid.UUID
}

func RestoreAccount(id uuid.UUID) *Account {
	return &Account{
		ID: id,
	}
}
