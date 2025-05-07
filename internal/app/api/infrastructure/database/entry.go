package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/repository"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/infrastructure/database/model"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/infrastructure/database/pkg/transaction"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/infrastructure/database/transformer"
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
	if entry == nil {
		return ErrRequiredEntry
	}

	driver := transaction.GetDriver(ctx, r.db)
	model := transformer.ToEntryModel(entry)
	_, err := driver.NamedExecContext(ctx, "INSERT INTO entries (id, account_id, volume_id, `key`, size, type, created_at, updated_at) VALUES (:id, :account_id, :volume_id, :key, :size, :type, :created_at, :updated_at);", model)
	return err
}

func (r *entryRepository) FindOneByKeyAndVolumeID(ctx context.Context, key string, volumeID uuid.UUID) (*entity.Entry, error) {
	driver := transaction.GetDriver(ctx, r.db)
	var model model.EntryModel
	if err := driver.QueryRowxContext(ctx, "SELECT id, account_id, volume_id, `key`, size, type, created_at, updated_at FROM entries WHERE `key` = ? AND volume_id = ? LIMIT 1;", key, volumeID).StructScan(&model); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return transformer.ToEntryEntity(&model), nil
}
