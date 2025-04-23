package builder

import (
	"github.com/atsumarukun/holos-storage-api/internal/app/api/interface/schema"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/usecase/dto"
)

func ToVolumeResponse(volume *dto.VolumeDTO) *schema.VolumeResponse {
	return &schema.VolumeResponse{
		ID:        volume.ID,
		Name:      volume.Name,
		IsPublic:  volume.IsPublic,
		CreatedAt: volume.CreatedAt,
		UpdatedAt: volume.UpdatedAt,
	}
}
