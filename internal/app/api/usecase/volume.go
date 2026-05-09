//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=../../../../test/mock/$GOPACKAGE/$GOFILE
package usecase

import (
	"context"
	stderr "errors"

	"github.com/atsumarukun/holos-api-pkg/errors"
	"github.com/google/uuid"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/repository"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/repository/pkg/transaction"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/service"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/usecase/dto"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/usecase/mapper"
)

var ErrVolumeNotFound = stderr.New("volume not found")

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
	bodyServ       service.BodyService
}

func NewVolumeUsecase(
	transactionObj transaction.TransactionObject,
	volumeRepo repository.VolumeRepository,
	bodyRepo repository.BodyRepository,
	volumeServ service.VolumeService,
	bodyServ service.BodyService,
) VolumeUsecase {
	return &volumeUsecase{
		transactionObj: transactionObj,
		volumeRepo:     volumeRepo,
		bodyRepo:       bodyRepo,
		volumeServ:     volumeServ,
		bodyServ:       bodyServ,
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

		path, err := u.bodyServ.BuildPath(volume, nil)
		if err != nil {
			return err
		}
		return u.bodyRepo.Create(path, nil)
	}); err != nil {
		return nil, err
	}

	return mapper.ToVolumeDTO(volume), nil
}

func (u *volumeUsecase) Update(ctx context.Context, accountID uuid.UUID, name, newName string, isPublic bool) (*dto.VolumeDTO, error) {
	const errMessage = "failed to update volume"

	var volume *entity.Volume
	if err := u.transactionObj.Transaction(ctx, func(ctx context.Context) error {
		var err error
		volume, err = u.getOne(ctx, accountID, name, errMessage)
		if err != nil {
			return err
		}

		srcPath, err := u.bodyServ.BuildPath(volume, nil)
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

		dstPath, err := u.bodyServ.BuildPath(volume, nil)
		if err != nil {
			return err
		}
		return u.bodyRepo.Update(srcPath, dstPath)
	}); err != nil {
		return nil, err
	}

	return mapper.ToVolumeDTO(volume), nil
}

func (u *volumeUsecase) Delete(ctx context.Context, accountID uuid.UUID, name string) error {
	const errMessage = "failed to delete volume"

	return u.transactionObj.Transaction(ctx, func(ctx context.Context) error {
		volume, err := u.getOne(ctx, accountID, name, errMessage)
		if err != nil {
			return err
		}

		if err := u.volumeServ.CanDelete(ctx, volume); err != nil {
			return err
		}
		if err := u.volumeRepo.Delete(ctx, volume); err != nil {
			return err
		}

		path, err := u.bodyServ.BuildPath(volume, nil)
		if err != nil {
			return err
		}
		return u.bodyRepo.Delete(path)
	})
}

func (u *volumeUsecase) GetOne(ctx context.Context, accountID uuid.UUID, name string) (*dto.VolumeDTO, error) {
	const errMessage = "failed to get volume"

	volume, err := u.getOne(ctx, accountID, name, errMessage)
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

func (u *volumeUsecase) getOne(ctx context.Context, accountID uuid.UUID, name, errMessage string) (*entity.Volume, error) {
	volume, err := u.volumeRepo.FindOneByNameAndAccountID(ctx, name, accountID)
	if err != nil {
		return nil, err
	}
	if volume == nil {
		return nil, errors.Wrap(ErrVolumeNotFound, errors.CodeNotFound, errMessage)
	}
	return volume, nil
}
