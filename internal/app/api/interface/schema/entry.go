package schema

import (
	"time"

	"github.com/google/uuid"
)

type CreateEntryRequest struct {
	VolumeID string `form:"volume_id" binding:"required"`
	Key      string `form:"key" binding:"required"`
	IsPublic string `form:"is_public" binding:"required"`
}

type EntryResponse struct {
	ID        uuid.UUID `json:"id"`
	VolumeID  uuid.UUID `json:"volume_id"`
	Key       string    `json:"key"`
	Size      uint64    `json:"size"`
	Type      string    `json:"type"`
	IsPublic  bool      `json:"is_public"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
