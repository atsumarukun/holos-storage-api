package database_test

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/infrastructure/database"
	mockDatabase "github.com/atsumarukun/holos-storage-api/test/mock/database"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func TestEntry_Create(t *testing.T) {
	entry := &entity.Entry{
		ID:        uuid.New(),
		AccountID: uuid.New(),
		VolumeID:  uuid.New(),
		Key:       "test/sample.jpg",
		Size:      10000,
		Type:      "image/jpeg",
		IsPublic:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name        string
		inputEntry  *entity.Entry
		expectError error
		setMockDB   func(mock sqlmock.Sqlmock)
	}{
		{
			name:        "success",
			inputEntry:  entry,
			expectError: nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO entries (id, account_id, volume_id, key, size, type, is_public, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);`)).
					WithArgs(entry.ID, entry.AccountID, entry.VolumeID, entry.Key, entry.Size, entry.Type, entry.IsPublic, entry.CreatedAt, entry.UpdatedAt).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(nil)
			},
		},
		{
			name:        "entry is nil",
			inputEntry:  nil,
			expectError: database.ErrRequiredEntry,
			setMockDB:   func(sqlmock.Sqlmock) {},
		},
		{
			name:        "insert error",
			inputEntry:  entry,
			expectError: sql.ErrConnDone,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO entries (id, account_id, volume_id, key, size, type, is_public, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);`)).
					WithArgs(entry.ID, entry.AccountID, entry.VolumeID, entry.Key, entry.Size, entry.Type, entry.IsPublic, entry.CreatedAt, entry.UpdatedAt).
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

			repo := database.NewEntryRepository(db)
			if err := repo.Create(t.Context(), tt.inputEntry); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestEntry_FindOneByKeyAndVolumeIDAndAccountID(t *testing.T) {
	accountID := uuid.New()
	volumeID := uuid.New()
	entry := &entity.Entry{
		ID:        uuid.New(),
		AccountID: accountID,
		VolumeID:  volumeID,
		Key:       "test/sample.jpg",
		Size:      10000,
		Type:      "image/jpeg",
		IsPublic:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name           string
		inputKey       string
		inputVolumeID  uuid.UUID
		inputAccountID uuid.UUID
		expectResult   *entity.Entry
		expectError    error
		setMockDB      func(mock sqlmock.Sqlmock)
	}{
		{
			name:           "success",
			inputKey:       "key",
			inputVolumeID:  volumeID,
			inputAccountID: accountID,
			expectResult:   entry,
			expectError:    nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, account_id, volume_id, key, size, type, is_public, created_at, updated_at FROM entries WHERE key = ? AND volume_id = ? AND account_id = ? LIMIT 1;`)).
					WithArgs("key", volumeID, accountID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "account_id", "volume_id", "key", "size", "type", "is_public", "created_at", "updated_at"}).AddRow(entry.ID, entry.AccountID, entry.Key, entry.Size, entry.Type, entry.IsPublic, entry.CreatedAt, entry.UpdatedAt)).
					WillReturnError(nil)
			},
		},
		{
			name:           "not found",
			inputKey:       "key",
			inputVolumeID:  volumeID,
			inputAccountID: accountID,
			expectResult:   nil,
			expectError:    nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, account_id, volume_id, key, size, type, is_public, created_at, updated_at FROM entries WHERE key = ? AND volume_id = ? AND account_id = ? LIMIT 1;`)).
					WithArgs("key", volumeID, accountID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "account_id", "volume_id", "key", "size", "type", "is_public", "created_at", "updated_at"})).
					WillReturnError(nil)
			},
		},
		{
			name:           "find error",
			inputKey:       "key",
			inputVolumeID:  volumeID,
			inputAccountID: accountID,
			expectResult:   nil,
			expectError:    sql.ErrConnDone,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, account_id, volume_id, key, size, type, is_public, created_at, updated_at FROM entries WHERE key = ? AND volume_id = ? AND account_id = ? LIMIT 1;`)).
					WithArgs("key", volumeID, accountID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "account_id", "volume_id", "key", "size", "type", "is_public", "created_at", "updated_at"})).
					WillReturnError(sql.ErrConnDone)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := mockDatabase.NewMockDatabase(t)

			tt.setMockDB(mock)

			repo := database.NewEntryRepository(db)
			result, err := repo.FindOneByKeyAndVolumeIDAndAccountID(t.Context(), tt.inputKey, tt.inputVolumeID, tt.inputAccountID)
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
