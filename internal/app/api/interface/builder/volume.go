package builder

import (
	"github.com/atsumarukun/holos-storage-api/internal/app/api/interface/schema"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/usecase/dto"
)

func ToVolumeResponse(volume *dto.VolumeDTO) *schema.VolumeResponse {
	return &schema.VolumeResponse{
		Name:      volume.Name,
		IsPublic:  volume.IsPublic,
		CreatedAt: volume.CreatedAt,
		UpdatedAt: volume.UpdatedAt,
	}
}

func ToVolumeResponses(volumes []*dto.VolumeDTO) []*schema.VolumeResponse {
	responses := make([]*schema.VolumeResponse, len(volumes))
	for i, volume := range volumes {
		responses[i] = ToVolumeResponse(volume)
	}
	return responses
}
