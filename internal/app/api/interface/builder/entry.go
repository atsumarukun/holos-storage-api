package builder

import (
	"github.com/atsumarukun/holos-storage-api/internal/app/api/interface/schema"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/usecase/dto"
)

func ToEntryResponse(entry *dto.EntryDTO) *schema.EntryResponse {
	return &schema.EntryResponse{
		ID:        entry.ID,
		VolumeID:  entry.VolumeID,
		Key:       entry.Key,
		Size:      entry.Size,
		Type:      entry.Type,
		IsPublic:  entry.IsPublic,
		CreatedAt: entry.CreatedAt,
		UpdatedAt: entry.UpdatedAt,
	}
}
