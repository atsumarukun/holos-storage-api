//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=../../../../test/mock/$GOPACKAGE/$GOFILE
package usecase

import (
	"context"
	"errors"
	"io"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/repository"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/repository/pkg/transaction"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/service"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/usecase/dto"
	"github.com/google/uuid"
)

type EntryUsecase interface {
	Create(context.Context, uuid.UUID, uuid.UUID, string, bool, io.Reader) (*dto.EntryDTO, error)
}

type entryUsecase struct {
	transactionObj transaction.TransactionObject
	entryRepo      repository.EntryRepository
	volumeRepo     repository.VolumeRepository
	entryServ      service.EntryService
}

func NewEntryUsecase(
	transactionObj transaction.TransactionObject,
	entryRepo repository.EntryRepository,
	volumeRepo repository.VolumeRepository,
	entryServ service.EntryService,
) EntryUsecase {
	return &entryUsecase{
		transactionObj: transactionObj,
		entryRepo:      entryRepo,
		volumeRepo:     volumeRepo,
		entryServ:      entryServ,
	}
}

func (u *entryUsecase) Create(ctx context.Context, accountID, volumeID uuid.UUID, key string, isPublic bool, body io.Reader) (*dto.EntryDTO, error) {
	return nil, errors.New("not implemented")
}
