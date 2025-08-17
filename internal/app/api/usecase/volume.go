//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=../../../../test/mock/$GOPACKAGE/$GOFILE
package usecase

import (
	"context"

	"github.com/google/uuid"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/repository"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/repository/pkg/transaction"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/service"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/usecase/dto"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/usecase/mapper"
)

type VolumeUsecase interface {
	Create(context.Context, uuid.UUID, string, bool) (*dto.VolumeDTO, error)
	Update(context.Context, uuid.UUID, string, string, bool) (*dto.VolumeDTO, error)
	Delete(context.Context, uuid.UUID, string) error
	GetOne(context.Context, uuid.UUID, string) (*dto.VolumeDTO, error)
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

func (u *volumeUsecase) Update(ctx context.Context, accountID uuid.UUID, name, newName string, isPublic bool) (*dto.VolumeDTO, error) {
	var volume *entity.Volume

	if err := u.transactionObj.Transaction(ctx, func(ctx context.Context) error {
		var err error
		volume, err = u.volumeRepo.FindOneByNameAndAccountID(ctx, name, accountID)
		if err != nil {
			return err
		}

		volume.SetIsPublic(isPublic)
		if volume.Name == newName {
			return u.volumeRepo.Update(ctx, volume)
		}

		if err := volume.SetName(newName); err != nil {
			return err
		}
		if err := u.volumeServ.Exists(ctx, volume); err != nil {
			return err
		}
		if err := u.volumeRepo.Update(ctx, volume); err != nil {
			return err
		}

		return u.bodyRepo.Update(name, volume.Name)
	}); err != nil {
		return nil, err
	}

	return mapper.ToVolumeDTO(volume), nil
}

func (u *volumeUsecase) Delete(ctx context.Context, accountID uuid.UUID, name string) error {
	return u.transactionObj.Transaction(ctx, func(ctx context.Context) error {
		volume, err := u.volumeRepo.FindOneByNameAndAccountID(ctx, name, accountID)
		if err != nil {
			return err
		}

		if err := u.volumeServ.CanDelete(ctx, volume); err != nil {
			return err
		}

		if err := u.volumeRepo.Delete(ctx, volume); err != nil {
			return err
		}

		return u.bodyRepo.Delete(volume.Name)
	})
}

func (u *volumeUsecase) GetOne(ctx context.Context, accountID uuid.UUID, name string) (*dto.VolumeDTO, error) {
	volume, err := u.volumeRepo.FindOneByNameAndAccountID(ctx, name, accountID)
	if err != nil {
		return nil, err
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
