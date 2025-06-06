//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=../../../../../test/mock/domain/$GOPACKAGE/$GOFILE
package service

import (
	"context"
	"errors"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/repository"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/pkg/status"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/pkg/status/code"
)

var (
	ErrRequiredVolume      = status.Error(code.Internal, "volume is required")
	ErrVolumeAlreadyExists = status.Error(code.Conflict, "volume name already used")
	ErrVolumeHasEntries    = status.Error(code.Conflict, "volume has entries")
)

type VolumeService interface {
	Exists(context.Context, *entity.Volume) error
	CanDelete(context.Context, *entity.Volume) error
}

type volumeService struct {
	volumeRepo repository.VolumeRepository
	entryRepo  repository.EntryRepository
}

func NewVolumeService(volumeRepo repository.VolumeRepository, entryRepo repository.EntryRepository) VolumeService {
	return &volumeService{
		volumeRepo: volumeRepo,
		entryRepo:  entryRepo,
	}
}

func (s *volumeService) Exists(ctx context.Context, volume *entity.Volume) error {
	if volume == nil {
		return ErrRequiredVolume
	}
	_, err := s.volumeRepo.FindOneByName(ctx, volume.Name)
	if err != nil {
		if errors.Is(err, repository.ErrVolumeNotFound) {
			return nil
		}
		return err
	}
	return ErrVolumeAlreadyExists
}

func (s *volumeService) CanDelete(ctx context.Context, volume *entity.Volume) error {
	if volume == nil {
		return ErrRequiredVolume
	}

	entries, err := s.entryRepo.FindByVolumeIDAndAccountID(ctx, volume.ID, volume.AccountID, nil, nil)
	if err != nil {
		return err
	}
	if 0 < len(entries) {
		return ErrVolumeHasEntries
	}
	return nil
}
