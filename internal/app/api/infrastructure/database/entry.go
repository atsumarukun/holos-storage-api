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
	const errMessage = "failed to create entry"

	if entry == nil {
		return errors.Wrap(repository.ErrNilEntry, errors.CodeInternalServerError, errMessage)
	}

	driver := transaction.GetDriver(ctx, r.db)
	model := transformer.ToEntryModel(entry)

	if _, err := driver.NamedExecContext(ctx, "INSERT INTO entries (id, account_id, volume_id, `key`, size, type, created_at, updated_at) VALUES (:id, :account_id, :volume_id, :key, :size, :type, :created_at, :updated_at);", model); err != nil {
		return errors.Wrap(err, errors.CodeInternalServerError, errMessage)
	}

	return nil
}

func (r *entryRepository) Update(ctx context.Context, entry *entity.Entry) error {
	const errMessage = "failed to update entry"

	if entry == nil {
		return errors.Wrap(repository.ErrNilEntry, errors.CodeInternalServerError, errMessage)
	}

	driver := transaction.GetDriver(ctx, r.db)
	model := transformer.ToEntryModel(entry)

	if _, err := driver.NamedExecContext(ctx, "UPDATE entries SET account_id = :account_id, volume_id = :volume_id, `key` = :key, size = :size, type = :type, updated_at = :updated_at WHERE id = :id LIMIT 1;", model); err != nil {
		return errors.Wrap(err, errors.CodeInternalServerError, errMessage)
	}

	return nil
}

func (r *entryRepository) Delete(ctx context.Context, entry *entity.Entry) error {
	const errMessage = "failed to delete entry"

	if entry == nil {
		return errors.Wrap(repository.ErrNilEntry, errors.CodeInternalServerError, errMessage)
	}

	driver := transaction.GetDriver(ctx, r.db)
	model := transformer.ToEntryModel(entry)

	if _, err := driver.NamedExecContext(ctx, "DELETE FROM entries WHERE id = :id LIMIT 1;", model); err != nil {
		return errors.Wrap(err, errors.CodeInternalServerError, errMessage)
	}

	return nil
}

func (r *entryRepository) FindOneByKeyAndVolumeID(ctx context.Context, key string, volumeID uuid.UUID) (*entity.Entry, error) {
	const errMessage = "failed to find entry by key and volume_id"

	return r.findOne(
		ctx,
		"SELECT id, account_id, volume_id, `key`, size, type, created_at, updated_at FROM entries WHERE `key` = ? AND volume_id = ? LIMIT 1;",
		[]any{key, volumeID},
		errMessage,
	)
}

func (r *entryRepository) FindOneByKeyAndVolumeIDAndAccountID(ctx context.Context, key string, volumeID, accountID uuid.UUID) (*entity.Entry, error) {
	const errMessage = "failed to find entry by key and volume_id and account_id"

	return r.findOne(
		ctx,
		"SELECT id, account_id, volume_id, `key`, size, type, created_at, updated_at FROM entries WHERE `key` = ? AND volume_id = ? AND account_id = ? LIMIT 1;",
		[]any{key, volumeID, accountID},
		errMessage,
	)
}

func (r *entryRepository) FindByVolumeIDAndAccountID(ctx context.Context, volumeID, accountID uuid.UUID, prefix *string, depth *uint64) (entries []*entity.Entry, err error) {
	const errMessage = "failed to find entry by volume_id and account_id"

	driver := transaction.GetDriver(ctx, r.db)

	filterQuery := " volume_id = ? AND account_id = ?"
	filterArguments := []any{volumeID, accountID}

	if prefix != nil {
		filterQuery += " AND `key` LIKE ?"
		filterArguments = append(filterArguments, *prefix+"/%")
	}
	if depth != nil {
		filterQuery += " AND LENGTH(`key`) - LENGTH(REPLACE(`key`, '/', '')) <= LENGTH(?) - LENGTH(REPLACE(?, '/', '')) + ?"
		if prefix == nil {
			filterArguments = append(filterArguments, "", "", *depth-1)
		} else {
			filterArguments = append(filterArguments, *prefix, *prefix, *depth)
		}
	}

	rows, err := driver.QueryxContext(ctx, "SELECT id, account_id, volume_id, `key`, size, type, created_at, updated_at FROM entries WHERE"+filterQuery+";", filterArguments...)
	if err != nil {
		return nil, errors.Wrap(err, errors.CodeInternalServerError, errMessage)
	}
	defer func() {
		// NOTE: errに直接詰めると関数内のエラーがnilで上書きされるためエラー発生時のみ上書きする.
		if e := rows.Close(); e != nil {
			err = errors.Wrap(e, errors.CodeInternalServerError, errMessage)
		}
	}()

	var models []*model.EntryModel
	for rows.Next() {
		var model model.EntryModel
		if err := rows.StructScan(&model); err != nil {
			return nil, errors.Wrap(err, errors.CodeInternalServerError, errMessage)
		}
		models = append(models, &model)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, errors.CodeInternalServerError, errMessage)
	}

	return transformer.ToEntryEntities(models), nil
}

// nolint:dupl // 集約単位のrepository実装. 集約境界を保つためrepository間での共通化は行わず重複を許容.
func (r *entryRepository) findOne(ctx context.Context, query string, args []any, errMessage string) (*entity.Entry, error) {
	driver := transaction.GetDriver(ctx, r.db)
	var model model.EntryModel

	if err := driver.QueryRowxContext(ctx, query, args...).StructScan(&model); err != nil {
		if stderr.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.Wrap(err, errors.CodeInternalServerError, errMessage)
	}

	return transformer.ToEntryEntity(&model), nil
}
