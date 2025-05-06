//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=../../../../../test/mock/domain/$GOPACKAGE/$GOFILE
package service

import (
	"context"
	"errors"
	"io"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/repository"
)

type EntryService interface {
	Exists(context.Context, *entity.Entry) error
	Create(context.Context, *entity.Entry, io.Reader) error
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
	return errors.New("not implemented")
}

func (s entryService) Create(ctx context.Context, entry *entity.Entry, reader io.Reader) error {
	return errors.New("not implemented")
}
