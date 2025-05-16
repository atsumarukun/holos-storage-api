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
	"github.com/atsumarukun/holos-storage-api/internal/app/api/pkg/status"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/pkg/status/code"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/usecase/dto"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/usecase/mapper"
)

var ErrEntryNotFound = status.Error(code.NotFound, "entry not found")

type EntryUsecase interface {
	Create(context.Context, uuid.UUID, uuid.UUID, string, uint64, io.Reader) (*dto.EntryDTO, error)
	Update(context.Context, uuid.UUID, string, string, string) (*dto.EntryDTO, error)
	Delete(context.Context, uuid.UUID, uuid.UUID) error
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

func (u *entryUsecase) Create(ctx context.Context, accountID, volumeID uuid.UUID, key string, size uint64, body io.Reader) (*dto.EntryDTO, error) {
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

		entry, err = entity.NewEntry(accountID, volumeID, key, size, entryType)
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

func (u *entryUsecase) Update(ctx context.Context, accountID uuid.UUID, volumeName, key, newKey string) (*dto.EntryDTO, error) {
	var entry *entity.Entry

	if err := u.transactionObj.Transaction(ctx, func(ctx context.Context) error {
		volume, err := u.volumeRepo.FindOneByNameAndAccountID(ctx, volumeName, accountID)
		if err != nil {
			return err
		}
		if volume == nil {
			return ErrVolumeNotFound
		}

		entry, err = u.entryRepo.FindOneByKeyAndVolumeIDAndAccountID(ctx, key, volume.ID, accountID)
		if err != nil {
			return err
		}
		if entry == nil {
			return ErrEntryNotFound
		}

		if err := entry.SetKey(newKey); err != nil {
			return err
		}

		if err := u.entryServ.Exists(ctx, entry); err != nil {
			return err
		}

		return u.entryServ.Update(ctx, volume, entry, key)
	}); err != nil {
		return nil, err
	}

	return mapper.ToEntryDTO(entry), nil
}

func (u *entryUsecase) Delete(ctx context.Context, accountID, id uuid.UUID) error {
	return u.transactionObj.Transaction(ctx, func(ctx context.Context) error {
		var err error
		entry, err := u.entryRepo.FindOneByIDAndAccountID(ctx, id, accountID)
		if err != nil {
			return err
		}
		if entry == nil {
			return ErrEntryNotFound
		}

		volume, err := u.volumeRepo.FindOneByIDAndAccountID(ctx, entry.VolumeID, accountID)
		if err != nil {
			return err
		}
		if volume == nil {
			return ErrVolumeNotFound
		}

		return u.entryServ.Delete(ctx, volume, entry)
	})
}
