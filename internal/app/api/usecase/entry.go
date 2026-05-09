//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=../../../../test/mock/$GOPACKAGE/$GOFILE
package usecase

import (
	"bytes"
	"context"
	stderr "errors"
	"io"
	"net/http"

	"github.com/atsumarukun/holos-api-pkg/errors"
	"github.com/google/uuid"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/repository"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/repository/pkg/transaction"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/service"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/usecase/dto"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/usecase/mapper"
)

var ErrEntryNotFound = stderr.New("entry not found")

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
	bodyServ       service.BodyService
}

func NewEntryUsecase(
	transactionObj transaction.TransactionObject,
	entryRepo repository.EntryRepository,
	bodyRepo repository.BodyRepository,
	volumeRepo repository.VolumeRepository,
	entryServ service.EntryService,
	bodyServ service.BodyService,
) EntryUsecase {
	return &entryUsecase{
		transactionObj: transactionObj,
		entryRepo:      entryRepo,
		bodyRepo:       bodyRepo,
		volumeRepo:     volumeRepo,
		entryServ:      entryServ,
		bodyServ:       bodyServ,
	}
}

func (u *entryUsecase) Create(ctx context.Context, accountID uuid.UUID, volumeName, key string, size uint64, body io.Reader) (*dto.EntryDTO, error) {
	const errMessage = "failed to update entry"

	var entry *entity.Entry
	if err := u.transactionObj.Transaction(ctx, func(ctx context.Context) error {
		volume, err := u.volumeRepo.FindOneByNameAndAccountID(ctx, volumeName, accountID)
		if err != nil {
			return err
		}
		if volume == nil {
			return errors.Wrap(ErrVolumeNotFound, errors.CodeNotFound, errMessage)
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

		path, err := u.bodyServ.BuildPath(volume, entry)
		if err != nil {
			return err
		}
		return u.bodyRepo.Create(path, bodyReader)
	}); err != nil {
		return nil, err
	}

	return mapper.ToEntryDTO(entry), nil
}

func (u *entryUsecase) Update(ctx context.Context, accountID uuid.UUID, volumeName, key, newKey string) (*dto.EntryDTO, error) {
	const errMessage = "failed to update entry"

	var entry *entity.Entry
	if err := u.transactionObj.Transaction(ctx, func(ctx context.Context) error {
		volume, err := u.volumeRepo.FindOneByNameAndAccountID(ctx, volumeName, accountID)
		if err != nil {
			return err
		}
		if volume == nil {
			return errors.Wrap(ErrVolumeNotFound, errors.CodeNotFound, errMessage)
		}

		entry, err = u.entryRepo.FindOneByKeyAndVolumeIDAndAccountID(ctx, key, volume.ID, accountID)
		if err != nil {
			return err
		}
		if entry == nil {
			return errors.Wrap(ErrEntryNotFound, errors.CodeNotFound, errMessage)
		}

		srcPath, err := u.bodyServ.BuildPath(volume, entry)
		if err != nil {
			return err
		}

		if err := u.update(ctx, entry, key, newKey); err != nil {
			return err
		}

		dstPath, err := u.bodyServ.BuildPath(volume, entry)
		if err != nil {
			return err
		}
		return u.bodyRepo.Update(srcPath, dstPath)
	}); err != nil {
		return nil, err
	}

	return mapper.ToEntryDTO(entry), nil
}

func (u *entryUsecase) Delete(ctx context.Context, accountID uuid.UUID, volumeName, key string) error {
	const errMessage = "failed to delete entry"

	return u.transactionObj.Transaction(ctx, func(ctx context.Context) error {
		volume, err := u.volumeRepo.FindOneByNameAndAccountID(ctx, volumeName, accountID)
		if err != nil {
			return err
		}
		if volume == nil {
			return errors.Wrap(ErrVolumeNotFound, errors.CodeNotFound, errMessage)
		}

		entry, err := u.entryRepo.FindOneByKeyAndVolumeIDAndAccountID(ctx, key, volume.ID, accountID)
		if err != nil {
			return err
		}
		if entry == nil {
			return nil
		}

		if err := u.entryServ.DeleteDescendants(ctx, entry); err != nil {
			return err
		}
		if err := u.entryRepo.Delete(ctx, entry); err != nil {
			return err
		}

		path, err := u.bodyServ.BuildPath(volume, entry)
		if err != nil {
			return err
		}
		return u.bodyRepo.Delete(path)
	})
}

func (u *entryUsecase) Copy(ctx context.Context, accountID uuid.UUID, volumeName, key string) (*dto.EntryDTO, error) {
	const errMessage = "failed to copy entry"

	var entry *entity.Entry
	if err := u.transactionObj.Transaction(ctx, func(ctx context.Context) error {
		volume, err := u.volumeRepo.FindOneByNameAndAccountID(ctx, volumeName, accountID)
		if err != nil {
			return err
		}
		if volume == nil {
			return errors.Wrap(ErrVolumeNotFound, errors.CodeNotFound, errMessage)
		}

		srcEntry, err := u.entryRepo.FindOneByKeyAndVolumeIDAndAccountID(ctx, key, volume.ID, accountID)
		if err != nil {
			return err
		}
		if srcEntry == nil {
			return errors.Wrap(ErrEntryNotFound, errors.CodeNotFound, errMessage)
		}

		srcPath, err := u.bodyServ.BuildPath(volume, srcEntry)
		if err != nil {
			return err
		}

		dstEntry, err := u.copy(ctx, srcEntry)
		if err != nil {
			return err
		}
		entry = dstEntry

		dstPath, err := u.bodyServ.BuildPath(volume, dstEntry)
		if err != nil {
			return err
		}
		return u.bodyRepo.Copy(srcPath, dstPath)
	}); err != nil {
		return nil, err
	}

	return mapper.ToEntryDTO(entry), nil
}

func (u *entryUsecase) GetMeta(ctx context.Context, accountID uuid.UUID, volumeName, key string) (*dto.EntryDTO, error) {
	const errMessage = "failed to get entry meta"

	var entry *entity.Entry
	if err := u.transactionObj.Transaction(ctx, func(ctx context.Context) error {
		volume, err := u.volumeRepo.FindOneByNameAndAccountID(ctx, volumeName, accountID)
		if err != nil {
			return err
		}
		if volume == nil {
			return errors.Wrap(ErrVolumeNotFound, errors.CodeNotFound, errMessage)
		}

		entry, err = u.entryRepo.FindOneByKeyAndVolumeIDAndAccountID(ctx, key, volume.ID, accountID)
		if err != nil {
			return err
		}
		if entry == nil {
			return errors.Wrap(ErrEntryNotFound, errors.CodeNotFound, errMessage)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return mapper.ToEntryDTO(entry), nil
}

func (u *entryUsecase) GetOne(ctx context.Context, accountID uuid.UUID, volumeName, key string) (*dto.EntryDTO, io.ReadCloser, error) {
	const errMessage = "failed to get entry"

	var entry *entity.Entry
	var body io.ReadCloser
	if err := u.transactionObj.Transaction(ctx, func(ctx context.Context) error {
		volume, err := u.volumeRepo.FindOneByNameAndAccountID(ctx, volumeName, accountID)
		if err != nil {
			return err
		}
		if volume == nil {
			return errors.Wrap(ErrVolumeNotFound, errors.CodeNotFound, errMessage)
		}

		entry, err = u.entryRepo.FindOneByKeyAndVolumeIDAndAccountID(ctx, key, volume.ID, accountID)
		if err != nil {
			return err
		}
		if entry == nil {
			return errors.Wrap(ErrEntryNotFound, errors.CodeNotFound, errMessage)
		}

		path, err := u.bodyServ.BuildPath(volume, entry)
		if err != nil {
			return err
		}
		body, err = u.bodyRepo.FindOneByPath(path)
		return err
	}); err != nil {
		return nil, nil, err
	}

	return mapper.ToEntryDTO(entry), body, nil
}

func (u *entryUsecase) Search(ctx context.Context, accountID uuid.UUID, volumeName string, prefix *string, depth *uint64) ([]*dto.EntryDTO, error) {
	const errMessage = "failed to search entry"

	var entries []*entity.Entry
	if err := u.transactionObj.Transaction(ctx, func(ctx context.Context) error {
		volume, err := u.volumeRepo.FindOneByNameAndAccountID(ctx, volumeName, accountID)
		if err != nil {
			return err
		}
		if volume == nil {
			return errors.Wrap(ErrVolumeNotFound, errors.CodeNotFound, errMessage)
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
		return "", nil, errors.Wrap(err, errors.CodeInternalServerError, "failed to get body info")
	}

	entryType := http.DetectContentType(buf[:n])
	bodyReader := io.MultiReader(bytes.NewReader(buf[:n]), body)

	return entryType, bodyReader, nil
}

func (u *entryUsecase) update(ctx context.Context, entry *entity.Entry, key, newKey string) error {
	if entry.Key == newKey {
		return nil
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

	return u.entryRepo.Update(ctx, entry)
}

func (u *entryUsecase) copy(ctx context.Context, src *entity.Entry) (*entity.Entry, error) {
	srcKey := src.Key

	dst, err := u.entryServ.Copy(ctx, src)
	if err != nil {
		return nil, err
	}
	if err := u.entryServ.CopyDescendants(ctx, dst, srcKey); err != nil {
		return nil, err
	}

	if err := u.entryRepo.Create(ctx, dst); err != nil {
		return nil, err
	}

	return dst, nil
}
