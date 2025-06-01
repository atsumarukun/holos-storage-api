//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=../../../../../test/mock/domain/$GOPACKAGE/$GOFILE
package service

import (
	"context"
	"errors"
	"path/filepath"
	"strings"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/repository"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/pkg/status"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/pkg/status/code"
)

var (
	ErrRequiredEntry      = status.Error(code.Internal, "entry is required")
	ErrEntryAlreadyExists = status.Error(code.Conflict, "entry key already used")
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
		return ErrRequiredEntry
	}

	_, err := s.entryRepo.FindOneByKeyAndVolumeID(ctx, entry.Key, entry.VolumeID)
	if err != nil {
		if errors.Is(err, repository.ErrEntryNotFound) {
			return nil
		}
		return err
	}
	return ErrEntryAlreadyExists
}

func (s *entryService) CreateAncestors(ctx context.Context, entry *entity.Entry) error {
	if entry == nil {
		return ErrRequiredEntry
	}

	for _, dir := range s.extractDirs(entry.Key) {
		ent, err := entity.NewEntry(entry.AccountID, entry.VolumeID, dir, 0, "folder")
		if err != nil {
			return err
		}
		if err := s.Exists(ctx, ent); err != nil {
			if errors.Is(err, ErrEntryAlreadyExists) {
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
		return ErrRequiredEntry
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
		return ErrRequiredEntry
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
		return nil, ErrRequiredEntry
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
		if errors.Is(err, ErrEntryAlreadyExists) {
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
		return ErrRequiredEntry
	}

	if entry.IsFolder() {
		descendants, err := s.entryRepo.FindByVolumeIDAndAccountID(ctx, entry.VolumeID, entry.AccountID, &entry.Key, nil)
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
