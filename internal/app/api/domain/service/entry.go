//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=../../../../../test/mock/domain/$GOPACKAGE/$GOFILE
package service

import (
	"context"
	"errors"
	"io"
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
	Create(context.Context, *entity.Volume, *entity.Entry, io.Reader) error
	Update(context.Context, *entity.Volume, *entity.Entry, string) error
	Delete(context.Context, *entity.Volume, *entity.Entry) error
}

type entryService struct {
	entryRepo repository.EntryRepository
	bodyRepo  repository.BodyRepository
}

func NewEntryService(entryRepo repository.EntryRepository, bodyRepo repository.BodyRepository) EntryService {
	return &entryService{
		entryRepo: entryRepo,
		bodyRepo:  bodyRepo,
	}
}

func (s *entryService) Exists(ctx context.Context, entry *entity.Entry) error {
	if entry == nil {
		return ErrRequiredEntry
	}

	ent, err := s.entryRepo.FindOneByKeyAndVolumeID(ctx, entry.Key, entry.VolumeID)
	if err != nil {
		return err
	}
	if ent == nil {
		return nil
	}
	return ErrEntryAlreadyExists
}

func (s *entryService) Create(ctx context.Context, volume *entity.Volume, entry *entity.Entry, body io.Reader) error {
	if volume == nil {
		return ErrRequiredVolume
	}
	if entry == nil {
		return ErrRequiredEntry
	}

	if err := s.createParentEntries(ctx, entry); err != nil {
		return err
	}

	if err := s.entryRepo.Create(ctx, entry); err != nil {
		return err
	}

	path := volume.Name + "/" + entry.Key
	return s.bodyRepo.Create(path, body)
}

func (s *entryService) Update(ctx context.Context, volume *entity.Volume, entry *entity.Entry, src string) error {
	if volume == nil {
		return ErrRequiredVolume
	}
	if entry == nil {
		return ErrRequiredEntry
	}

	if err := s.createParentEntries(ctx, entry); err != nil {
		return err
	}

	if entry.IsFolder() {
		children, err := s.entryRepo.FindByKeyPrefixAndAccountID(ctx, src+"/", entry.AccountID)
		if err != nil {
			return err
		}

		for _, child := range children {
			key := strings.Replace(child.Key, src, entry.Key, 1)
			if err := child.SetKey(key); err != nil {
				return err
			}
			if err := s.entryRepo.Update(ctx, child); err != nil {
				return err
			}
		}
	}

	if err := s.entryRepo.Update(ctx, entry); err != nil {
		return err
	}

	return s.bodyRepo.Update(volume.Name+"/"+src, volume.Name+"/"+entry.Key)
}

func (s *entryService) Delete(ctx context.Context, volume *entity.Volume, entry *entity.Entry) error {
	if volume == nil {
		return ErrRequiredVolume
	}
	if entry == nil {
		return ErrRequiredEntry
	}

	if entry.IsFolder() {
		children, err := s.entryRepo.FindByKeyPrefixAndAccountID(ctx, entry.Key+"/", entry.AccountID)
		if err != nil {
			return err
		}

		for _, child := range children {
			if err := s.entryRepo.Delete(ctx, child); err != nil {
				return err
			}
		}
	}

	if err := s.entryRepo.Delete(ctx, entry); err != nil {
		return err
	}

	return s.bodyRepo.Delete(volume.Name + "/" + entry.Key)
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

func (s *entryService) createParentEntries(ctx context.Context, entry *entity.Entry) error {
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
