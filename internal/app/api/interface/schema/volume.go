package schema

import (
	"time"

	"github.com/google/uuid"
)

type CreateVolumeRequest struct {
	Name     string `json:"name"`
	IsPublic bool   `json:"is_public"`
}

type UpdateVolumeRequest struct {
	Name     string `json:"name"`
	IsPublic bool   `json:"is_public"`
}

type VolumeResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	IsPublic  bool      `json:"is_public"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
