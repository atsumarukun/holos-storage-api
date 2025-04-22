//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=../../../../../test/mock/domain/$GOPACKAGE/$GOFILE
package service

import (
	"context"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/repository"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/pkg/status"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/pkg/status/code"
)

var (
	ErrRequiredVolume      = status.Error(code.Internal, "volume is required")
	ErrVolumeAlreadyExists = status.Error(code.Conflict, "volume name already used")
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
	if volume == nil {
		return ErrRequiredVolume
	}

	vol, err := s.volumeRepo.FindOneByNameAndAccountID(ctx, volume.Name, volume.AccountID)
	if err != nil {
		return err
	}
	if vol == nil {
		return nil
	}
	return ErrVolumeAlreadyExists
}
