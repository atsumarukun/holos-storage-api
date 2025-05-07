package model

import (
	"time"

	"github.com/google/uuid"
)

type EntryModel struct {
	ID        uuid.UUID `db:"id"`
	AccountID uuid.UUID `db:"account_id"`
	VolumeID  uuid.UUID `db:"volume_id"`
	Key       string    `db:"key"`
	Size      uint64    `db:"size"`
	Type      string    `db:"type"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
