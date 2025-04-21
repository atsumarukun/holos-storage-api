package database

import (
	"context"
	"errors"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/repository"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type volumeRepository struct {
	db *sqlx.DB
}

func NewVolumeRepository(db *sqlx.DB) repository.VolumeRepository {
	return &volumeRepository{
		db: db,
	}
}

func (r *volumeRepository) Create(ctx context.Context, volume *entity.Volume) error {
	return errors.New("not implemented")
}

func (r *volumeRepository) FindOneByNameAndAccountID(ctx context.Context, name string, accountID uuid.UUID) (*entity.Volume, error) {
	return nil, errors.New("not implemented")
}
