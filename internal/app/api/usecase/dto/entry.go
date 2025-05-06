package dto

import (
	"time"

	"github.com/google/uuid"
)

type EntryDTO struct {
	ID        uuid.UUID
	AccountID uuid.UUID
	VolumeID  uuid.UUID
	Key       string
	Size      uint64
	Type      string
	IsPublic  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
