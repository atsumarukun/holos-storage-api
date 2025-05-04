package database

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/repository"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/pkg/status"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/pkg/status/code"
)

var ErrRequiredEntry = status.Error(code.Internal, "entry is required")

type entryRepository struct {
	db *sqlx.DB
}

func NewEntryRepository(db *sqlx.DB) repository.EntryRepository {
	return &entryRepository{
		db: db,
	}
}

func (r *entryRepository) Create(ctx context.Context, entry *entity.Entry) error {
	return errors.New("not implemented")
}

func (r *entryRepository) FindOneByKeyAndVolumeIDAndAccountID(ctx context.Context, key string, volumeID uuid.UUID, accountID uuid.UUID) (*entity.Entry, error) {
	return nil, errors.New("not implemented")
}
