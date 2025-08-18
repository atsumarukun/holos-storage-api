package schema

import "time"

type CreateEntryRequest struct {
	Key string `form:"key" binding:"required"`
}

type UpdateEntryRequest struct {
	Key string `json:"key"`
}

type EntryResponse struct {
	Key       string    `json:"key"`
	Size      uint64    `json:"size"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
