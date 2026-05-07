package database

import (
	"context"
	"database/sql"
	stderr "errors"

	"github.com/atsumarukun/holos-api-pkg/errors"
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
	const errMessage = "failed to create volume"

	if volume == nil {
		return errors.Wrap(repository.ErrNilVolume, errors.CodeInternalServerError, errMessage)
	}

	driver := transaction.GetDriver(ctx, r.db)
	model := transformer.ToVolumeModel(volume)

	if _, err := driver.NamedExecContext(ctx, "INSERT INTO volumes (id, account_id, name, is_public, created_at, updated_at) VALUES (:id, :account_id, :name, :is_public, :created_at, :updated_at);", model); err != nil {
		return errors.Wrap(err, errors.CodeInternalServerError, errMessage)
	}

	return nil
}

func (r *volumeRepository) Update(ctx context.Context, volume *entity.Volume) error {
	const errMessage = "failed to update volume"

	if volume == nil {
		return errors.Wrap(repository.ErrNilVolume, errors.CodeInternalServerError, errMessage)
	}

	driver := transaction.GetDriver(ctx, r.db)
	model := transformer.ToVolumeModel(volume)

	if _, err := driver.NamedExecContext(ctx, "UPDATE volumes SET account_id = :account_id, name = :name, is_public = :is_public, updated_at = :updated_at WHERE id = :id LIMIT 1;", model); err != nil {
		return errors.Wrap(err, errors.CodeInternalServerError, errMessage)
	}

	return nil
}

func (r *volumeRepository) Delete(ctx context.Context, volume *entity.Volume) error {
	const errMessage = "failed to update volume"

	if volume == nil {
		return errors.Wrap(repository.ErrNilVolume, errors.CodeInternalServerError, errMessage)
	}

	driver := transaction.GetDriver(ctx, r.db)
	model := transformer.ToVolumeModel(volume)

	if _, err := driver.NamedExecContext(ctx, "DELETE FROM volumes WHERE id = :id LIMIT 1;", model); err != nil {
		return errors.Wrap(err, errors.CodeInternalServerError, errMessage)
	}

	return nil
}

func (r *volumeRepository) FindOneByName(ctx context.Context, name string) (*entity.Volume, error) {
	const errMessage = "failed to find volume by name"

	return r.findOne(
		ctx,
		`SELECT id, account_id, name, is_public, created_at, updated_at FROM volumes WHERE name = ? LIMIT 1;`,
		[]any{name},
		errMessage,
	)
}

func (r *volumeRepository) FindOneByNameAndAccountID(ctx context.Context, name string, accountID uuid.UUID) (*entity.Volume, error) {
	const errMessage = "failed to find volume by name and account_id"

	return r.findOne(
		ctx,
		`SELECT id, account_id, name, is_public, created_at, updated_at FROM volumes WHERE name = ? AND account_id = ? LIMIT 1;`,
		[]any{name, accountID},
		errMessage,
	)
}

func (r *volumeRepository) FindOneByIDAndAccountID(ctx context.Context, id, accountID uuid.UUID) (*entity.Volume, error) {
	const errMessage = "failed to find volume by id and account_id"

	return r.findOne(
		ctx,
		`SELECT id, account_id, name, is_public, created_at, updated_at FROM volumes WHERE id = ? AND account_id = ? LIMIT 1;`,
		[]any{id, accountID},
		errMessage,
	)
}

func (r *volumeRepository) FindByAccountID(ctx context.Context, accountID uuid.UUID) (volumes []*entity.Volume, err error) {
	const errMessage = "failed to find volume by account_id"

	driver := transaction.GetDriver(ctx, r.db)
	rows, err := driver.QueryxContext(ctx, `SELECT id, account_id, name, is_public, created_at, updated_at FROM volumes WHERE account_id = ?;`, accountID)
	if err != nil {
		return nil, errors.Wrap(err, errors.CodeInternalServerError, errMessage)
	}
	defer func() {
		// NOTE: errに直接詰めると関数内のエラーがnilで上書きされるためエラー発生時のみ上書きする.
		if e := rows.Close(); e != nil {
			err = errors.Wrap(e, errors.CodeInternalServerError, errMessage)
		}
	}()

	var models []*model.VolumeModel
	for rows.Next() {
		var model model.VolumeModel
		if err := rows.StructScan(&model); err != nil {
			return nil, errors.Wrap(err, errors.CodeInternalServerError, errMessage)
		}
		models = append(models, &model)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, errors.CodeInternalServerError, errMessage)
	}

	return transformer.ToVolumeEntities(models), nil
}

// nolint:dupl // 集約単位のrepository実装. 集約境界を保つためrepository間での共通化は行わず重複を許容.
func (r *volumeRepository) findOne(ctx context.Context, query string, args []any, errMessage string) (*entity.Volume, error) {
	driver := transaction.GetDriver(ctx, r.db)
	var model model.VolumeModel

	if err := driver.QueryRowxContext(ctx, query, args...).StructScan(&model); err != nil {
		if stderr.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.Wrap(err, errors.CodeInternalServerError, errMessage)
	}

	return transformer.ToVolumeEntity(&model), nil
}
