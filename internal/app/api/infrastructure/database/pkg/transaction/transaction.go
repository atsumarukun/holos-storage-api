package transaction

import (
	"context"
	"database/sql"

	"github.com/atsumarukun/holos-api-pkg/errors"
	"github.com/jmoiron/sqlx"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/repository/pkg/transaction"
)

type transactionKey struct{}

type transactionObject struct {
	db *sqlx.DB
}

func NewDBTransactionObject(db *sqlx.DB) transaction.TransactionObject {
	return &transactionObject{
		db: db,
	}
}

func (to *transactionObject) Transaction(ctx context.Context, fn func(context.Context) error) (err error) {
	tx, err := to.db.Beginx()
	if err != nil {
		return errors.Wrap(err, errors.CodeInternalServerError, "failed to begin transaction")
	}

	defer func() {
		if r := recover(); r != nil {
			err = tx.Rollback()
			if err != nil {
				err = errors.Wrap(err, errors.CodeInternalServerError, "failed to rollback transaction")
			}
		}
	}()

	ctx = context.WithValue(ctx, transactionKey{}, tx)

	if err := fn(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, errors.CodeInternalServerError, "failed to commit transaction")
	}

	return nil
}

type driver interface {
	sqlx.Queryer
	sqlx.QueryerContext
	sqlx.Execer
	sqlx.ExecerContext
	NamedExecContext(ctx context.Context, query string, arg any) (sql.Result, error)
}

func GetDriver(ctx context.Context, db *sqlx.DB) driver {
	if tx, ok := ctx.Value(transactionKey{}).(*sqlx.Tx); ok {
		return tx
	}
	return db
}
