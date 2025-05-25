//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=../../../../test/mock/$GOPACKAGE/$GOFILE
package usecase

import (
	"context"

	"github.com/google/uuid"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/repository"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/repository/pkg/transaction"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/service"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/pkg/status"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/pkg/status/code"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/usecase/dto"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/usecase/mapper"
)

var ErrVolumeNotFound = status.Error(code.NotFound, "volume not found")

type VolumeUsecase interface {
	Create(context.Context, uuid.UUID, string, bool) (*dto.VolumeDTO, error)
	Update(context.Context, uuid.UUID, uuid.UUID, string, bool) (*dto.VolumeDTO, error)
	Delete(context.Context, uuid.UUID, uuid.UUID) error
	GetOne(context.Context, uuid.UUID, uuid.UUID) (*dto.VolumeDTO, error)
	GetAll(context.Context, uuid.UUID) ([]*dto.VolumeDTO, error)
}

type volumeUsecase struct {
	transactionObj transaction.TransactionObject
	volumeRepo     repository.VolumeRepository
	bodyRepo       repository.BodyRepository
	volumeServ     service.VolumeService
}

func NewVolumeUsecase(
	transactionObj transaction.TransactionObject,
	volumeRepo repository.VolumeRepository,
	bodyRepo repository.BodyRepository,
	volumeServ service.VolumeService,
) VolumeUsecase {
	return &volumeUsecase{
		transactionObj: transactionObj,
		volumeRepo:     volumeRepo,
		bodyRepo:       bodyRepo,
		volumeServ:     volumeServ,
	}
}

func (u *volumeUsecase) Create(ctx context.Context, accountID uuid.UUID, name string, isPublic bool) (*dto.VolumeDTO, error) {
	volume, err := entity.NewVolume(accountID, name, isPublic)
	if err != nil {
		return nil, err
	}

	if err := u.transactionObj.Transaction(ctx, func(ctx context.Context) error {
		if err := u.volumeServ.Exists(ctx, volume); err != nil {
			return err
		}

		if err := u.volumeRepo.Create(ctx, volume); err != nil {
			return err
		}

		return u.bodyRepo.Create(volume.Name, nil)
	}); err != nil {
		return nil, err
	}

	return mapper.ToVolumeDTO(volume), nil
}

func (u *volumeUsecase) Update(ctx context.Context, accountID, id uuid.UUID, name string, isPublic bool) (*dto.VolumeDTO, error) {
	var volume *entity.Volume

	if err := u.transactionObj.Transaction(ctx, func(ctx context.Context) error {
		var err error
		volume, err = u.volumeRepo.FindOneByIDAndAccountID(ctx, id, accountID)
		if err != nil {
			return err
		}
		if volume == nil {
			return ErrVolumeNotFound
		}

		oldName := volume.Name

		if err := volume.SetName(name); err != nil {
			return err
		}
		volume.SetIsPublic(isPublic)

		if err := u.volumeServ.Exists(ctx, volume); err != nil {
			return err
		}

		if err := u.volumeRepo.Update(ctx, volume); err != nil {
			return err
		}

		return u.bodyRepo.Update(oldName, volume.Name)
	}); err != nil {
		return nil, err
	}

	return mapper.ToVolumeDTO(volume), nil
}

func (u *volumeUsecase) Delete(ctx context.Context, accountID, id uuid.UUID) error {
	return u.transactionObj.Transaction(ctx, func(ctx context.Context) error {
		volume, err := u.volumeRepo.FindOneByIDAndAccountID(ctx, id, accountID)
		if err != nil {
			return err
		}
		if volume == nil {
			return ErrVolumeNotFound
		}

		if err := u.volumeRepo.Delete(ctx, volume); err != nil {
			return err
		}

		return u.bodyRepo.Delete(volume.Name)
	})
}

func (u *volumeUsecase) GetOne(ctx context.Context, accountID, id uuid.UUID) (*dto.VolumeDTO, error) {
	volume, err := u.volumeRepo.FindOneByIDAndAccountID(ctx, id, accountID)
	if err != nil {
		return nil, err
	}
	if volume == nil {
		return nil, ErrVolumeNotFound
	}
	return mapper.ToVolumeDTO(volume), nil
}

func (u *volumeUsecase) GetAll(ctx context.Context, accountID uuid.UUID) ([]*dto.VolumeDTO, error) {
	volumes, err := u.volumeRepo.FindByAccountID(ctx, accountID)
	if err != nil {
		return nil, err
	}
	return mapper.ToVolumeDTOs(volumes), nil
}
