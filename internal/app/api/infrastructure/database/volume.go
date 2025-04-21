package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/repository"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/infrastructure/database/model"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/infrastructure/database/pkg/transaction"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/infrastructure/database/transformer"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/pkg/status"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/pkg/status/code"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var ErrRequiredVolume = status.Error(code.Internal, "volume is required")

type volumeRepository struct {
	db *sqlx.DB
}

func NewVolumeRepository(db *sqlx.DB) repository.VolumeRepository {
	return &volumeRepository{
		db: db,
	}
}

func (r *volumeRepository) Create(ctx context.Context, volume *entity.Volume) error {
	if volume == nil {
		return ErrRequiredVolume
	}

	driver := transaction.GetDriver(ctx, r.db)
	model := transformer.ToVolumeModel(volume)
	_, err := driver.NamedExecContext(ctx, "INSERT INTO volumes (id, account_id, name, is_public, created_at, updated_at) VALUES (:id, :account_id, :name, :is_public, :created_at, :updated_at);", model)
	return err
}

func (r *volumeRepository) FindOneByNameAndAccountID(ctx context.Context, name string, accountID uuid.UUID) (*entity.Volume, error) {
	driver := transaction.GetDriver(ctx, r.db)
	var model model.VolumeModel
	if err := driver.QueryRowxContext(ctx, `SELECT id, account_id, name, is_public, created_at, updated_at FROM volumes WHERE name = ? AND account_id = ? LIMIT 1;`, name, accountID).StructScan(&model); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return transformer.ToVolumeEntity(&model), nil
}
