package model

import (
	"time"

	"github.com/google/uuid"
)

type VolumeModel struct {
	ID        uuid.UUID `db:"id"`
	AccountID uuid.UUID `db:"account_id"`
	Name      string    `db:"name"`
	IsPublic  bool      `db:"is_public"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
