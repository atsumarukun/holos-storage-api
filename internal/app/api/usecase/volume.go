//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=../../../../test/mock/$GOPACKAGE/$GOFILE
package usecase

import (
	"context"
	"errors"

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
	Update(context.Context, uuid.UUID, uuid.UUID, string, bool) (*dto.VolumeDTO, error)
	Delete(context.Context, uuid.UUID, uuid.UUID) error
}

type volumeUsecase struct {
	transactionObj transaction.TransactionObject
	volumeRepo     repository.VolumeRepository
	volumeServ     service.VolumeService
}

func NewVolumeUsecase(
	transactionObj transaction.TransactionObject,
	volumeRepo repository.VolumeRepository,
	volumeServ service.VolumeService,
) VolumeUsecase {
	return &volumeUsecase{
		transactionObj: transactionObj,
		volumeRepo:     volumeRepo,
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

		return u.volumeRepo.Create(ctx, volume)
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

		if err := volume.SetName(name); err != nil {
			return err
		}
		volume.SetIsPublic(isPublic)

		if err := u.volumeServ.Exists(ctx, volume); err != nil {
			return err
		}

		return u.volumeRepo.Update(ctx, volume)
	}); err != nil {
		return nil, err
	}

	return mapper.ToVolumeDTO(volume), nil
}

func (u *volumeUsecase) Delete(ctx context.Context, accountID, id uuid.UUID) error {
	return errors.New("not implemented")
}
