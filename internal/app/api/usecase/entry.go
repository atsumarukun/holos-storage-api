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
	Create(context.Context, uuid.UUID, string, string, uint64, io.Reader) (*dto.EntryDTO, error)
	Update(context.Context, uuid.UUID, string, string, string) (*dto.EntryDTO, error)
	Delete(context.Context, uuid.UUID, string, string) error
	Head(context.Context, uuid.UUID, string, string) (*dto.EntryDTO, error)
	GetOne(context.Context, uuid.UUID, string, string) (*dto.EntryDTO, io.ReadCloser, error)
}

type entryUsecase struct {
	transactionObj transaction.TransactionObject
	entryRepo      repository.EntryRepository
	bodyRepo       repository.BodyRepository
	volumeRepo     repository.VolumeRepository
	entryServ      service.EntryService
}

func NewEntryUsecase(
	transactionObj transaction.TransactionObject,
	entryRepo repository.EntryRepository,
	bodyRepo repository.BodyRepository,
	volumeRepo repository.VolumeRepository,
	entryServ service.EntryService,
) EntryUsecase {
	return &entryUsecase{
		transactionObj: transactionObj,
		entryRepo:      entryRepo,
		bodyRepo:       bodyRepo,
		volumeRepo:     volumeRepo,
		entryServ:      entryServ,
	}
}

func (u *entryUsecase) Create(ctx context.Context, accountID uuid.UUID, volumeName, key string, size uint64, body io.Reader) (*dto.EntryDTO, error) {
	var entry *entity.Entry

	if err := u.transactionObj.Transaction(ctx, func(ctx context.Context) error {
		volume, err := u.volumeRepo.FindOneByNameAndAccountID(ctx, volumeName, accountID)
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

		entry, err = entity.NewEntry(accountID, volume.ID, key, size, entryType)
		if err != nil {
			return err
		}

		if err := u.entryServ.Exists(ctx, entry); err != nil {
			return err
		}

		if err := u.entryServ.Create(ctx, entry, bodyReader); err != nil {
			return err
		}

		path := volume.Name + "/" + entry.Key
		return u.bodyRepo.Create(path, bodyReader)
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

		if err := u.entryServ.Update(ctx, entry, key); err != nil {
			return err
		}

		src := volume.Name + "/" + key
		dst := volume.Name + "/" + entry.Key
		return u.bodyRepo.Update(src, dst)
	}); err != nil {
		return nil, err
	}

	return mapper.ToEntryDTO(entry), nil
}

func (u *entryUsecase) Delete(ctx context.Context, accountID uuid.UUID, volumeName, key string) error {
	return u.transactionObj.Transaction(ctx, func(ctx context.Context) error {
		volume, err := u.volumeRepo.FindOneByNameAndAccountID(ctx, volumeName, accountID)
		if err != nil {
			return err
		}
		if volume == nil {
			return ErrVolumeNotFound
		}

		entry, err := u.entryRepo.FindOneByKeyAndVolumeIDAndAccountID(ctx, key, volume.ID, accountID)
		if err != nil {
			return err
		}
		if entry == nil {
			return ErrEntryNotFound
		}

		if err := u.entryServ.Delete(ctx, entry); err != nil {
			return err
		}

		path := volume.Name + "/" + entry.Key
		return u.bodyRepo.Delete(path)
	})
}

func (u *entryUsecase) Head(ctx context.Context, accountID uuid.UUID, volumeName, key string) (*dto.EntryDTO, error) {
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

		return nil
	}); err != nil {
		return nil, err
	}

	return mapper.ToEntryDTO(entry), nil
}

func (u *entryUsecase) GetOne(ctx context.Context, accountID uuid.UUID, volumeName, key string) (*dto.EntryDTO, io.ReadCloser, error) {
	var entry *entity.Entry
	var body io.ReadCloser

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

		path := volume.Name + "/" + entry.Key
		body, err = u.bodyRepo.FindOneByPath(path)
		return err
	}); err != nil {
		return nil, nil, err
	}

	return mapper.ToEntryDTO(entry), body, nil
}
