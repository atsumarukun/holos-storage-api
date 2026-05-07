//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=../../../../../test/mock/domain/$GOPACKAGE/$GOFILE
package service

import (
	"context"
	stderr "errors"

	"github.com/atsumarukun/holos-api-pkg/errors"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/repository"
)

var (
	ErrNilVolume              = stderr.New("volume must not be nil")
	ErrVolumeNameAlreadyInUse = stderr.New("volume name already in use")
	ErrVolumeHasEntries       = stderr.New("volume cannot be deleted because it contains entries")
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
		return errors.Wrap(ErrNilVolume, errors.CodeInternalServerError, "failed to check if volume exists")
	}

	vol, err := s.volumeRepo.FindOneByName(ctx, volume.Name)
	if err != nil {
		return err
	}
	if vol != nil {
		return errors.Wrap(ErrVolumeNameAlreadyInUse, errors.CodeDuplicate, "volume already exists")
	}

	return nil
}

func (s *volumeService) CanDelete(ctx context.Context, volume *entity.Volume) error {
	if volume == nil {
		return errors.Wrap(ErrNilVolume, errors.CodeInternalServerError, "failed to check if volume can be deleted")
	}

	entries, err := s.entryRepo.FindByVolumeIDAndAccountID(ctx, volume.ID, volume.AccountID, nil, nil)
	if err != nil {
		return err
	}
	if len(entries) > 0 {
		return errors.Wrap(ErrVolumeHasEntries, errors.CodeConstraintViolation, "volume cannot be deleted")
	}

	return nil
}
