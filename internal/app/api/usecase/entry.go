//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=../../../../test/mock/$GOPACKAGE/$GOFILE
package usecase

import (
	"bytes"
	"context"
	"errors"
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
	Create(context.Context, uuid.UUID, string, string, uint64, io.Reader) (*dto.EntryDTO, error)
	Update(context.Context, uuid.UUID, string, string, string) (*dto.EntryDTO, error)
	Delete(context.Context, uuid.UUID, string, string) error
	Copy(context.Context, uuid.UUID, string, string) (*dto.EntryDTO, error)
	GetMeta(context.Context, uuid.UUID, string, string) (*dto.EntryDTO, error)
	GetOne(context.Context, uuid.UUID, string, string) (*dto.EntryDTO, io.ReadCloser, error)
	Search(context.Context, uuid.UUID, string, *string, *uint64) ([]*dto.EntryDTO, error)
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

		entryType, bodyReader, err := u.getBodyInfo(body)
		if err != nil {
			return err
		}

		entry, err = entity.NewEntry(accountID, volume.ID, key, size, entryType)
		if err != nil {
			return err
		}

		if err := u.entryServ.Exists(ctx, entry); err != nil {
			return err
		}
		if err := u.entryServ.CreateAncestors(ctx, entry); err != nil {
			return err
		}

		if err := u.entryRepo.Create(ctx, entry); err != nil {
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

		entry, err = u.entryRepo.FindOneByKeyAndVolumeIDAndAccountID(ctx, key, volume.ID, accountID)
		if err != nil {
			return err
		}

		if err := entry.SetKey(newKey); err != nil {
			return err
		}

		if err := u.entryServ.Exists(ctx, entry); err != nil {
			return err
		}
		if err := u.entryServ.CreateAncestors(ctx, entry); err != nil {
			return err
		}
		if err := u.entryServ.UpdateDescendants(ctx, entry, key); err != nil {
			return err
		}

		if err := u.entryRepo.Update(ctx, entry); err != nil {
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

		entry, err := u.entryRepo.FindOneByKeyAndVolumeIDAndAccountID(ctx, key, volume.ID, accountID)
		if err != nil {
			return err
		}

		if err := u.entryServ.DeleteDescendants(ctx, entry); err != nil {
			return err
		}

		if err := u.entryRepo.Delete(ctx, entry); err != nil {
			return err
		}

		path := volume.Name + "/" + entry.Key
		return u.bodyRepo.Delete(path)
	})
}

func (u *entryUsecase) Copy(ctx context.Context, accountID uuid.UUID, volumeName, key string) (*dto.EntryDTO, error) {
	return nil, errors.New("not implemented")
}

func (u *entryUsecase) GetMeta(ctx context.Context, accountID uuid.UUID, volumeName, key string) (*dto.EntryDTO, error) {
	var entry *entity.Entry

	if err := u.transactionObj.Transaction(ctx, func(ctx context.Context) error {
		volume, err := u.volumeRepo.FindOneByNameAndAccountID(ctx, volumeName, accountID)
		if err != nil {
			return err
		}

		entry, err = u.entryRepo.FindOneByKeyAndVolumeIDAndAccountID(ctx, key, volume.ID, accountID)
		if err != nil {
			return err
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

		entry, err = u.entryRepo.FindOneByKeyAndVolumeIDAndAccountID(ctx, key, volume.ID, accountID)
		if err != nil {
			return err
		}

		path := volume.Name + "/" + entry.Key
		body, err = u.bodyRepo.FindOneByPath(path)
		return err
	}); err != nil {
		return nil, nil, err
	}

	return mapper.ToEntryDTO(entry), body, nil
}

func (u *entryUsecase) Search(ctx context.Context, accountID uuid.UUID, volumeName string, prefix *string, depth *uint64) ([]*dto.EntryDTO, error) {
	var entries []*entity.Entry

	if err := u.transactionObj.Transaction(ctx, func(ctx context.Context) error {
		volume, err := u.volumeRepo.FindOneByNameAndAccountID(ctx, volumeName, accountID)
		if err != nil {
			return err
		}

		entries, err = u.entryRepo.FindByVolumeIDAndAccountID(ctx, volume.ID, accountID, prefix, depth)
		return err
	}); err != nil {
		return nil, err
	}

	return mapper.ToEntryDTOs(entries), nil
}

func (u *entryUsecase) getBodyInfo(body io.Reader) (string, io.Reader, error) {
	if body == nil {
		return "folder", nil, nil
	}

	buf := make([]byte, 512)
	n, err := body.Read(buf)
	if err != nil && err != io.EOF {
		return "", nil, err
	}

	entryType := http.DetectContentType(buf[:n])
	bodyReader := io.MultiReader(bytes.NewReader(buf[:n]), body)

	return entryType, bodyReader, nil
}
