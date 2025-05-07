package mapper

import (
	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/usecase/dto"
)

func ToEntryDTO(entry *entity.Entry) *dto.EntryDTO {
	return &dto.EntryDTO{
		ID:        entry.ID,
		AccountID: entry.AccountID,
		VolumeID:  entry.VolumeID,
		Key:       entry.Key,
		Size:      entry.Size,
		Type:      entry.Type,
		CreatedAt: entry.CreatedAt,
		UpdatedAt: entry.UpdatedAt,
	}
}
