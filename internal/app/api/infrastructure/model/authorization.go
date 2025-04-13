package model

import "github.com/google/uuid"

type AuthorizationModel struct {
	AccountID uuid.UUID `json:"id"`
}
