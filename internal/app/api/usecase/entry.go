//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=../../../../test/mock/$GOPACKAGE/$GOFILE
package usecase

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/google/uuid"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/repository"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/repository/pkg/transaction"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/service"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/usecase/dto"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/usecase/mapper"
)

type EntryUsecase interface {
	Create(context.Context, uuid.UUID, uuid.UUID, string, uint64, bool, io.Reader) (*dto.EntryDTO, error)
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

func (u *entryUsecase) Create(ctx context.Context, accountID, volumeID uuid.UUID, key string, size uint64, isPublic bool, body io.Reader) (*dto.EntryDTO, error) {
	var entry *entity.Entry

	if err := u.transactionObj.Transaction(ctx, func(ctx context.Context) error {
		volume, err := u.volumeRepo.FindOneByIDAndAccountID(ctx, volumeID, accountID)
		if err != nil {
			return err
		}
		if volume == nil {
			return ErrVolumeNotFound
		}

		entryType := "folder"
		bodyReader := body
		if body != nil {
			buf := make([]byte, 512)
			n, err := body.Read(buf)
			if err != nil && err != io.EOF {
				return err
			}
			entryType = http.DetectContentType(buf[:n])
			bodyReader = io.MultiReader(bytes.NewReader(buf[:n]), body)
		}

		entry, err = entity.NewEntry(accountID, volumeID, key, size, entryType, isPublic)
		if err != nil {
			return err
		}

		if err := u.entryServ.Exists(ctx, entry); err != nil {
			return err
		}

		return u.entryServ.Create(ctx, volume, entry, bodyReader)
	}); err != nil {
		return nil, err
	}

	return mapper.ToEntryDTO(entry), nil
}
