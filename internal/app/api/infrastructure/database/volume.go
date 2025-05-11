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

func (r *volumeRepository) Update(ctx context.Context, volume *entity.Volume) error {
	if volume == nil {
		return ErrRequiredVolume
	}

	driver := transaction.GetDriver(ctx, r.db)
	model := transformer.ToVolumeModel(volume)
	_, err := driver.NamedExecContext(ctx, "UPDATE volumes SET account_id = :account_id, name = :name, is_public = :is_public, updated_at = :updated_at WHERE id = :id LIMIT 1;", model)
	return err
}

func (r *volumeRepository) Delete(ctx context.Context, volume *entity.Volume) error {
	if volume == nil {
		return ErrRequiredVolume
	}

	driver := transaction.GetDriver(ctx, r.db)
	model := transformer.ToVolumeModel(volume)
	_, err := driver.NamedExecContext(ctx, "DELETE FROM volumes WHERE id = :id LIMIT 1;", model)
	return err
}

func (r *volumeRepository) FindOneByName(ctx context.Context, name string) (*entity.Volume, error) {
	driver := transaction.GetDriver(ctx, r.db)
	var model model.VolumeModel
	if err := driver.QueryRowxContext(ctx, `SELECT id, account_id, name, is_public, created_at, updated_at FROM volumes WHERE name = ? LIMIT 1;`, name).StructScan(&model); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return transformer.ToVolumeEntity(&model), nil
}

func (r *volumeRepository) FindOneByIDAndAccountID(ctx context.Context, id, accountID uuid.UUID) (*entity.Volume, error) {
	driver := transaction.GetDriver(ctx, r.db)
	var model model.VolumeModel
	if err := driver.QueryRowxContext(ctx, `SELECT id, account_id, name, is_public, created_at, updated_at FROM volumes WHERE id = ? AND account_id = ? LIMIT 1;`, id, accountID).StructScan(&model); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return transformer.ToVolumeEntity(&model), nil
}

func (r *volumeRepository) FindByAccountID(ctx context.Context, accountID uuid.UUID) (volumes []*entity.Volume, err error) {
	driver := transaction.GetDriver(ctx, r.db)
	rows, err := driver.QueryxContext(ctx, `SELECT id, account_id, name, is_public, created_at, updated_at FROM volumes WHERE account_id = ?;`, accountID)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
	}()

	var models []*model.VolumeModel
	for rows.Next() {
		var model model.VolumeModel
		if err := rows.StructScan(&model); err != nil {
			return nil, err
		}
		models = append(models, &model)
	}
	return transformer.ToVolumeEntities(models), nil
}
