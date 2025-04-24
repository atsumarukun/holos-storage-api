package transformer

import (
	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/infrastructure/database/model"
)

func ToVolumeModel(volume *entity.Volume) *model.VolumeModel {
	return &model.VolumeModel{
		ID:        volume.ID,
		AccountID: volume.AccountID,
		Name:      volume.Name,
		IsPublic:  volume.IsPublic,
		CreatedAt: volume.CreatedAt,
		UpdatedAt: volume.UpdatedAt,
	}
}

func ToVolumeEntity(volume *model.VolumeModel) *entity.Volume {
	return entity.RestoreVolume(
		volume.ID,
		volume.AccountID,
		volume.Name,
		volume.IsPublic,
		volume.CreatedAt,
		volume.UpdatedAt,
	)
}

func ToVolumeEntities(volumes []*model.VolumeModel) []*entity.Volume {
	entities := make([]*entity.Volume, len(volumes))
	for i, volume := range volumes {
		entities[i] = ToVolumeEntity(volume)
	}
	return entities
}
