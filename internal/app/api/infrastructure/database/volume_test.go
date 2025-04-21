package database_test

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/infrastructure/database"
	mockDatabase "github.com/atsumarukun/holos-storage-api/test/mock/database"
)

func TestVolume_Create(t *testing.T) {
	volume := &entity.Volume{
		ID:        uuid.New(),
		AccountID: uuid.New(),
		Name:      "name",
		IsPublic:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name        string
		inputVolume *entity.Volume
		expectError error
		setMockDB   func(mock sqlmock.Sqlmock)
	}{
		{
			name:        "success",
			inputVolume: volume,
			expectError: nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO volumes (id, account_id, name, is_public, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?);`)).
					WithArgs(volume.ID, volume.AccountID, volume.Name, volume.IsPublic, volume.CreatedAt, volume.UpdatedAt).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(nil)
			},
		},
		{
			name:        "volume is nil",
			inputVolume: nil,
			expectError: database.ErrRequiredVolume,
			setMockDB:   func(sqlmock.Sqlmock) {},
		},
		{
			name:        "insert error",
			inputVolume: volume,
			expectError: sql.ErrConnDone,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO volumes (id, account_id, name, is_public, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?);`)).
					WithArgs(volume.ID, volume.AccountID, volume.Name, volume.IsPublic, volume.CreatedAt, volume.UpdatedAt).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(sql.ErrConnDone)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := mockDatabase.NewMockDatabase(t)
			defer db.Close()

			tt.setMockDB(mock)

			repo := database.NewVolumeRepository(db)
			if err := repo.Create(t.Context(), tt.inputVolume); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestVolume_FindOneByNameAndAccountID(t *testing.T) {
	accountID := uuid.New()
	volume := &entity.Volume{
		ID:        uuid.New(),
		AccountID: accountID,
		Name:      "name",
		IsPublic:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name           string
		inputName      string
		inputAccountID uuid.UUID
		expectResult   *entity.Volume
		expectError    error
		setMockDB      func(mock sqlmock.Sqlmock)
	}{
		{
			name:         "success",
			inputName:    "name",
			expectResult: volume,
			expectError:  nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, account_id, name, is_public, created_at, updated_at FROM volumes WHERE name = ? AND account_id = ? LIMIT 1;`)).
					WithArgs("name", accountID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "account_id", "name", "is_public", "created_at", "updated_at"}).AddRow(volume.ID, volume.AccountID, volume.Name, volume.IsPublic, volume.CreatedAt, volume.UpdatedAt)).
					WillReturnError(nil)
			},
		},
		{
			name:         "not found",
			inputName:    "name",
			expectResult: nil,
			expectError:  nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, account_id, name, is_public, created_at, updated_at FROM volumes WHERE name = ? AND account_id = ? LIMIT 1;`)).
					WithArgs("name", accountID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "account_id", "name", "is_public", "created_at", "updated_at"})).
					WillReturnError(nil)
			},
		},
		{
			name:         "find error",
			inputName:    "name",
			expectResult: nil,
			expectError:  sql.ErrConnDone,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, account_id, name, is_public, created_at, updated_at FROM volumes WHERE name = ? AND account_id = ? LIMIT 1;`)).
					WithArgs("name", accountID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "account_id", "name", "is_public", "created_at", "updated_at"})).
					WillReturnError(sql.ErrConnDone)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := mockDatabase.NewMockDatabase(t)

			tt.setMockDB(mock)

			repo := database.NewVolumeRepository(db)
			result, err := repo.FindOneByNameAndAccountID(t.Context(), tt.inputName, tt.inputAccountID)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			if diff := cmp.Diff(result, tt.expectResult); diff != "" {
				t.Error(diff)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Error(err)
			}
		})
	}
}
