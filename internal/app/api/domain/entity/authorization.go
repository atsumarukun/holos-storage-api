package entity

import "github.com/google/uuid"

type Authorization struct {
	AccountID uuid.UUID
}

func RestoreAuthorization(accountID uuid.UUID) *Authorization {
	return &Authorization{
		AccountID: accountID,
	}
}
