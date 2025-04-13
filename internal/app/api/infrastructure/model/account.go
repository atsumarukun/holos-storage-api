package model

import "github.com/google/uuid"

type AccountModel struct {
	ID uuid.UUID `json:"id"`
}
