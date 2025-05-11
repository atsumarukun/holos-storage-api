package transformer

import (
	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/infrastructure/database/model"
)

func ToEntryModel(entry *entity.Entry) *model.EntryModel {
	return &model.EntryModel{
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

func ToEntryEntity(entry *model.EntryModel) *entity.Entry {
	return entity.RestoreEntry(
		entry.ID,
		entry.AccountID,
		entry.VolumeID,
		entry.Key,
		entry.Size,
		entry.Type,
		entry.CreatedAt,
		entry.UpdatedAt,
	)
}

func ToEntryEntities(entries []*model.EntryModel) []*entity.Entry {
	entities := make([]*entity.Entry, len(entries))
	for i, entry := range entries {
		entities[i] = ToEntryEntity(entry)
	}
	return entities
}
