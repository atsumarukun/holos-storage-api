//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=../../../../../test/mock/domain/$GOPACKAGE/$GOFILE
package service

import (
	"context"
	"errors"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/repository"
)

type VolumeService interface {
	Exists(context.Context, *entity.Volume) error
}

type volumeService struct {
	volumeRepo repository.VolumeRepository
}

func NewVolumeService(volumeRepo repository.VolumeRepository) VolumeService {
	return &volumeService{
		volumeRepo: volumeRepo,
	}
}

func (s *volumeService) Exists(ctx context.Context, volume *entity.Volume) error {
	return errors.New("not implemented")
}
