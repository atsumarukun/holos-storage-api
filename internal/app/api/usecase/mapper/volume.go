package mapper

import (
	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/usecase/dto"
)

func ToVolumeDTO(volume *entity.Volume) *dto.VolumeDTO {
	return &dto.VolumeDTO{
		ID:        volume.ID,
		AccountID: volume.AccountID,
		Name:      volume.Name,
		IsPublic:  volume.IsPublic,
		CreatedAt: volume.CreatedAt,
		UpdatedAt: volume.UpdatedAt,
	}
}
