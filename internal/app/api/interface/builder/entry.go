package builder

import (
	"github.com/atsumarukun/holos-storage-api/internal/app/api/interface/schema"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/usecase/dto"
)

func ToEntryResponse(entry *dto.EntryDTO) *schema.EntryResponse {
	return &schema.EntryResponse{
		Key:       entry.Key,
		Size:      entry.Size,
		Type:      entry.Type,
		CreatedAt: entry.CreatedAt,
		UpdatedAt: entry.UpdatedAt,
	}
}

func ToEntryResponses(entries []*dto.EntryDTO) []*schema.EntryResponse {
	responses := make([]*schema.EntryResponse, len(entries))
	for i, entry := range entries {
		responses[i] = ToEntryResponse(entry)
	}
	return responses
}
