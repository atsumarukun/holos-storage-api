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
	ErrEntryAlreadyExists = status.Error(code.Conflict, "entry name already used")
)

type EntryService interface {
	Exists(context.Context, *entity.Entry) error
	Create(context.Context, *entity.Volume, *entity.Entry, io.Reader) error
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

	ent, err := s.entryRepo.FindOneByKeyAndVolumeIDAndAccountID(ctx, entry.Key, entry.VolumeID, entry.AccountID)
	if err != nil {
		return err
	}
	if ent == nil {
		return nil
	}
	return ErrEntryAlreadyExists
}

func (s entryService) Create(ctx context.Context, volume *entity.Volume, entry *entity.Entry, body io.Reader) error {
	if volume == nil {
		return ErrRequiredVolume
	}
	if entry == nil {
		return ErrRequiredEntry
	}

	for _, dir := range s.extractDirs(entry.Key) {
		ent, err := entity.NewEntry(entry.AccountID, volume.ID, dir, 0, "folder", entry.IsPublic)
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

	if err := s.entryRepo.Create(ctx, entry); err != nil {
		return err
	}

	if body != nil {
		path := volume.Name + "/" + entry.Key
		return s.bodyRepo.Create(path, body)
	}
	return nil
}

func (s entryService) extractDirs(key string) []string {
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
