//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=../../../../../test/mock/domain/$GOPACKAGE/$GOFILE
package service

import (
	"context"
	stderr "errors"
	"path/filepath"
	"strings"

	"github.com/atsumarukun/holos-api-pkg/errors"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/repository"
)

var (
	ErrNilEntry             = stderr.New("entry must not be nil")
	ErrEntryKeyAlreadyInUse = stderr.New("entry key already in use")
)

type EntryService interface {
	Exists(context.Context, *entity.Entry) error
	CreateAncestors(context.Context, *entity.Entry) error
	UpdateDescendants(context.Context, *entity.Entry, string) error
	DeleteDescendants(context.Context, *entity.Entry) error
	Copy(context.Context, *entity.Entry) (*entity.Entry, error)
	CopyDescendants(context.Context, *entity.Entry, string) error
}

type entryService struct {
	entryRepo repository.EntryRepository
}

func NewEntryService(entryRepo repository.EntryRepository) EntryService {
	return &entryService{
		entryRepo: entryRepo,
	}
}

func (s *entryService) Exists(ctx context.Context, entry *entity.Entry) error {
	if entry == nil {
		return errors.Wrap(ErrNilEntry, errors.CodeInternalServerError, "failed to check if entry exists")
	}

	ent, err := s.entryRepo.FindOneByKeyAndVolumeID(ctx, entry.Key, entry.VolumeID)
	if err != nil {
		return err
	}
	if ent != nil {
		return errors.Wrap(ErrEntryKeyAlreadyInUse, errors.CodeDuplicate, "entry already exists")
	}

	return nil
}

func (s *entryService) CreateAncestors(ctx context.Context, entry *entity.Entry) error {
	if entry == nil {
		return errors.Wrap(ErrNilEntry, errors.CodeInternalServerError, "failed to create ancestor entries")
	}

	for _, dir := range s.extractDirs(entry.Key) {
		ent, err := entity.NewEntry(entry.AccountID, entry.VolumeID, dir, 0, "folder")
		if err != nil {
			return err
		}
		if err := s.Exists(ctx, ent); err != nil {
			if stderr.Is(err, ErrEntryKeyAlreadyInUse) {
				continue
			} else {
				return err
			}
		}
		if err := s.entryRepo.Create(ctx, ent); err != nil {
			return err
		}
	}

	return nil
}

func (s *entryService) UpdateDescendants(ctx context.Context, entry *entity.Entry, src string) error {
	if entry == nil {
		return errors.Wrap(ErrNilEntry, errors.CodeInternalServerError, "failed to update descendant entries")
	}

	if entry.IsFolder() {
		descendants, err := s.entryRepo.FindByVolumeIDAndAccountID(ctx, entry.VolumeID, entry.AccountID, &src, nil)
		if err != nil {
			return err
		}

		for _, descendant := range descendants {
			key := strings.Replace(descendant.Key, src, entry.Key, 1)
			if err := descendant.SetKey(key); err != nil {
				return err
			}
			if err := s.entryRepo.Update(ctx, descendant); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *entryService) DeleteDescendants(ctx context.Context, entry *entity.Entry) error {
	if entry == nil {
		return errors.Wrap(ErrNilEntry, errors.CodeInternalServerError, "failed to delete descendant entries")
	}

	if entry.IsFolder() {
		descendants, err := s.entryRepo.FindByVolumeIDAndAccountID(ctx, entry.VolumeID, entry.AccountID, &entry.Key, nil)
		if err != nil {
			return err
		}

		for _, descendant := range descendants {
			if err := s.entryRepo.Delete(ctx, descendant); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *entryService) Copy(ctx context.Context, entry *entity.Entry) (*entity.Entry, error) {
	if entry == nil {
		return nil, errors.Wrap(ErrNilEntry, errors.CodeInternalServerError, "failed to copy entry")
	}

	name := filepath.Base(entry.Key)
	ext := filepath.Ext(name)
	base := strings.TrimSuffix(name, ext)
	key := strings.Replace(entry.Key, name, base+" copy"+ext, 1)

	copied, err := entity.NewEntry(entry.AccountID, entry.VolumeID, key, entry.Size, entry.Type)
	if err != nil {
		return nil, err
	}

	if err := s.Exists(ctx, copied); err != nil {
		if stderr.Is(err, ErrEntryKeyAlreadyInUse) {
			copied, err = s.Copy(ctx, copied)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return copied, nil
}

func (s *entryService) CopyDescendants(ctx context.Context, entry *entity.Entry, src string) error {
	if entry == nil {
		return errors.Wrap(ErrNilEntry, errors.CodeInternalServerError, "failed to copy descendant entries")
	}

	if entry.IsFolder() {
		descendants, err := s.entryRepo.FindByVolumeIDAndAccountID(ctx, entry.VolumeID, entry.AccountID, &src, nil)
		if err != nil {
			return err
		}

		for _, descendant := range descendants {
			key := strings.Replace(descendant.Key, src, entry.Key, 1)
			copied, err := entity.NewEntry(descendant.AccountID, descendant.VolumeID, key, descendant.Size, descendant.Type)
			if err != nil {
				return err
			}
			if err := s.entryRepo.Create(ctx, copied); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *entryService) extractDirs(key string) []string {
	dirKey := filepath.Dir(key)
	if dirKey == "." {
		return nil
	}

	dirs := make([]string, strings.Count(dirKey, "/")+1)
	var current string

	for i, part := range strings.Split(dirKey, "/") {
		current += part + "/"
		dirs[i] = current
	}

	return dirs
}
