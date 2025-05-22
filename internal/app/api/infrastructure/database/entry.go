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

func (r *entryRepository) Update(ctx context.Context, entry *entity.Entry) error {
	if entry == nil {
		return ErrRequiredEntry
	}

	driver := transaction.GetDriver(ctx, r.db)
	model := transformer.ToEntryModel(entry)
	_, err := driver.NamedExecContext(ctx, "UPDATE entries SET account_id = :account_id, volume_id = :volume_id, `key` = :key, size = :size, type = :type, updated_at = :updated_at WHERE id = :id LIMIT 1;", model)
	return err
}

func (r *entryRepository) Delete(ctx context.Context, entry *entity.Entry) error {
	if entry == nil {
		return ErrRequiredEntry
	}

	driver := transaction.GetDriver(ctx, r.db)
	model := transformer.ToEntryModel(entry)
	_, err := driver.NamedExecContext(ctx, "DELETE FROM entries WHERE id = :id LIMIT 1;", model)
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

func (r *entryRepository) FindOneByKeyAndVolumeIDAndAccountID(ctx context.Context, key string, volumeID, accountID uuid.UUID) (*entity.Entry, error) {
	driver := transaction.GetDriver(ctx, r.db)
	var model model.EntryModel
	if err := driver.QueryRowxContext(ctx, "SELECT id, account_id, volume_id, `key`, size, type, created_at, updated_at FROM entries WHERE `key` = ? AND volume_id = ? AND account_id = ? LIMIT 1;", key, volumeID, accountID).StructScan(&model); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return transformer.ToEntryEntity(&model), nil
}

func (r *entryRepository) FindByKeyPrefixAndAccountID(ctx context.Context, keyword string, accountID uuid.UUID) (entries []*entity.Entry, err error) {
	driver := transaction.GetDriver(ctx, r.db)

	rows, err := driver.QueryxContext(ctx, "SELECT id, account_id, volume_id, `key`, size, type, created_at, updated_at FROM entries WHERE `key` LIKE ? AND account_id = ?;", keyword+"%", accountID)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
	}()

	var models []*model.EntryModel
	for rows.Next() {
		var model model.EntryModel
		if err := rows.StructScan(&model); err != nil {
			return nil, err
		}
		models = append(models, &model)
	}

	return transformer.ToEntryEntities(models), nil
}

func (r *entryRepository) FindByVolumeIDAndAccountID(ctx context.Context, volumeID, accountID uuid.UUID, prefix *string, depth *uint64) ([]*entity.Entry, error) {
	return nil, errors.New("not implemented")
}
