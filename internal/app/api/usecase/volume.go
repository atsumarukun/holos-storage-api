//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=../../../../test/mock/$GOPACKAGE/$GOFILE
package usecase

import (
	"context"
	"errors"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/repository"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/repository/pkg/transaction"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/service"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/usecase/dto"
	"github.com/google/uuid"
)

type VolumeUsecase interface {
	Create(context.Context, uuid.UUID, string, bool) (*dto.VolumeDTO, error)
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
	return nil, errors.New("not implemented")
}
