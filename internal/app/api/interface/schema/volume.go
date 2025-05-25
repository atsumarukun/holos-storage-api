package schema

import (
	"time"
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
	Name      string    `json:"name"`
	IsPublic  bool      `json:"is_public"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
