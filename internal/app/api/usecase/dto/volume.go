package dto

import (
	"time"

	"github.com/google/uuid"
)

type VolumeDTO struct {
	ID        uuid.UUID
	AccountID uuid.UUID
	Name      string
	IsPublic  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
