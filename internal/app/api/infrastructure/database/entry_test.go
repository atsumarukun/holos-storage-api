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
	"github.com/atsumarukun/holos-storage-api/internal/app/api/pkg/types"
	mockDatabase "github.com/atsumarukun/holos-storage-api/test/mock/database"
)

func TestEntry_Create(t *testing.T) {
	entry := &entity.Entry{
		ID:        uuid.New(),
		AccountID: uuid.New(),
		VolumeID:  uuid.New(),
		Key:       "test/sample.jpg",
		Size:      10000,
		Type:      "image/jpeg",
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
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO entries (id, account_id, volume_id, `key`, size, type, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?);")).
					WithArgs(entry.ID, entry.AccountID, entry.VolumeID, entry.Key, entry.Size, entry.Type, entry.CreatedAt, entry.UpdatedAt).
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
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO entries (id, account_id, volume_id, `key`, size, type, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?);")).
					WithArgs(entry.ID, entry.AccountID, entry.VolumeID, entry.Key, entry.Size, entry.Type, entry.CreatedAt, entry.UpdatedAt).
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

func TestEntry_Update(t *testing.T) {
	entry := &entity.Entry{
		ID:        uuid.New(),
		AccountID: uuid.New(),
		VolumeID:  uuid.New(),
		Key:       "test/sample.jpg",
		Size:      10000,
		Type:      "image/jpeg",
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
				mock.ExpectExec(regexp.QuoteMeta("UPDATE entries SET account_id = ?, volume_id = ?, `key` = ?, size = ?, type = ?, updated_at = ? WHERE id = ? LIMIT 1;")).
					WithArgs(entry.AccountID, entry.VolumeID, entry.Key, entry.Size, entry.Type, entry.UpdatedAt, entry.ID).
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
			name:        "update error",
			inputEntry:  entry,
			expectError: sql.ErrConnDone,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("UPDATE entries SET account_id = ?, volume_id = ?, `key` = ?, size = ?, type = ?, updated_at = ? WHERE id = ? LIMIT 1;")).
					WithArgs(entry.AccountID, entry.VolumeID, entry.Key, entry.Size, entry.Type, entry.UpdatedAt, entry.ID).
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
			if err := repo.Update(t.Context(), tt.inputEntry); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestEntry_Delete(t *testing.T) {
	entry := &entity.Entry{
		ID:        uuid.New(),
		AccountID: uuid.New(),
		VolumeID:  uuid.New(),
		Key:       "test/sample.jpg",
		Size:      10000,
		Type:      "image/jpeg",
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
				mock.ExpectExec(regexp.QuoteMeta("DELETE FROM entries WHERE id = ? LIMIT 1;")).
					WithArgs(entry.ID).
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
			name:        "delete error",
			inputEntry:  entry,
			expectError: sql.ErrConnDone,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("DELETE FROM entries WHERE id = ? LIMIT 1;")).
					WithArgs(entry.ID).
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
			if err := repo.Delete(t.Context(), tt.inputEntry); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestEntry_FindOneByKeyAndVolumeID(t *testing.T) {
	accountID := uuid.New()
	volumeID := uuid.New()
	entry := &entity.Entry{
		ID:        uuid.New(),
		AccountID: accountID,
		VolumeID:  volumeID,
		Key:       "test/sample.jpg",
		Size:      10000,
		Type:      "image/jpeg",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name          string
		inputKey      string
		inputVolumeID uuid.UUID
		expectResult  *entity.Entry
		expectError   error
		setMockDB     func(mock sqlmock.Sqlmock)
	}{
		{
			name:          "success",
			inputKey:      "key",
			inputVolumeID: volumeID,
			expectResult:  entry,
			expectError:   nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, account_id, volume_id, `key`, size, type, created_at, updated_at FROM entries WHERE `key` = ? AND volume_id = ? LIMIT 1;")).
					WithArgs("key", volumeID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "account_id", "volume_id", "key", "size", "type", "created_at", "updated_at"}).AddRow(entry.ID, entry.AccountID, entry.VolumeID, entry.Key, entry.Size, entry.Type, entry.CreatedAt, entry.UpdatedAt)).
					WillReturnError(nil)
			},
		},
		{
			name:          "not found",
			inputKey:      "key",
			inputVolumeID: volumeID,
			expectResult:  nil,
			expectError:   nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, account_id, volume_id, `key`, size, type, created_at, updated_at FROM entries WHERE `key` = ? AND volume_id = ? LIMIT 1;")).
					WithArgs("key", volumeID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "account_id", "volume_id", "key", "size", "type", "created_at", "updated_at"})).
					WillReturnError(nil)
			},
		},
		{
			name:          "find error",
			inputKey:      "key",
			inputVolumeID: volumeID,
			expectResult:  nil,
			expectError:   sql.ErrConnDone,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, account_id, volume_id, `key`, size, type, created_at, updated_at FROM entries WHERE `key` = ? AND volume_id = ? LIMIT 1;")).
					WithArgs("key", volumeID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "account_id", "volume_id", "key", "size", "type", "created_at", "updated_at"})).
					WillReturnError(sql.ErrConnDone)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := mockDatabase.NewMockDatabase(t)

			tt.setMockDB(mock)

			repo := database.NewEntryRepository(db)
			result, err := repo.FindOneByKeyAndVolumeID(t.Context(), tt.inputKey, tt.inputVolumeID)
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
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, account_id, volume_id, `key`, size, type, created_at, updated_at FROM entries WHERE `key` = ? AND volume_id = ? AND account_id = ? LIMIT 1;")).
					WithArgs("key", volumeID, accountID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "account_id", "volume_id", "key", "size", "type", "created_at", "updated_at"}).AddRow(entry.ID, entry.AccountID, entry.VolumeID, entry.Key, entry.Size, entry.Type, entry.CreatedAt, entry.UpdatedAt)).
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
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, account_id, volume_id, `key`, size, type, created_at, updated_at FROM entries WHERE `key` = ? AND volume_id = ? AND account_id = ? LIMIT 1;")).
					WithArgs("key", volumeID, accountID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "account_id", "volume_id", "key", "size", "type", "created_at", "updated_at"})).
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
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, account_id, volume_id, `key`, size, type, created_at, updated_at FROM entries WHERE `key` = ? AND volume_id = ? AND account_id = ? LIMIT 1;")).
					WithArgs("key", volumeID, accountID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "account_id", "volume_id", "key", "size", "type", "created_at", "updated_at"})).
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

func TestEntry_FindByKeyPrefixAndAccountID(t *testing.T) {
	entry := &entity.Entry{
		ID:        uuid.New(),
		AccountID: uuid.New(),
		VolumeID:  uuid.New(),
		Key:       "test/sample.jpg",
		Size:      10000,
		Type:      "image/jpeg",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name           string
		inputKeyword   string
		inputAccountID uuid.UUID
		expectResult   []*entity.Entry
		expectError    error
		setMockDB      func(mock sqlmock.Sqlmock)
	}{
		{
			name:           "success",
			inputKeyword:   "test/",
			inputAccountID: entry.AccountID,
			expectResult:   []*entity.Entry{entry},
			expectError:    nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, account_id, volume_id, `key`, size, type, created_at, updated_at FROM entries WHERE `key` LIKE ? AND account_id = ?;")).
					WithArgs("test/%", entry.AccountID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "account_id", "volume_id", "key", "size", "type", "created_at", "updated_at"}).AddRow(entry.ID, entry.AccountID, entry.VolumeID, entry.Key, entry.Size, entry.Type, entry.CreatedAt, entry.UpdatedAt)).
					WillReturnError(nil)
			},
		},
		{
			name:           "not found",
			inputKeyword:   "test/",
			inputAccountID: entry.AccountID,
			expectResult:   []*entity.Entry{},
			expectError:    nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, account_id, volume_id, `key`, size, type, created_at, updated_at FROM entries WHERE `key` LIKE ? AND account_id = ?;")).
					WithArgs("test/%", entry.AccountID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "account_id", "volume_id", "key", "size", "type", "created_at", "updated_at"})).
					WillReturnError(nil)
			},
		},
		{
			name:           "find error",
			inputKeyword:   "test/",
			inputAccountID: entry.AccountID,
			expectResult:   nil,
			expectError:    sql.ErrConnDone,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, account_id, volume_id, `key`, size, type, created_at, updated_at FROM entries WHERE `key` LIKE ? AND account_id = ?;")).
					WithArgs("test/%", entry.AccountID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "account_id", "volume_id", "key", "size", "type", "created_at", "updated_at"})).
					WillReturnError(sql.ErrConnDone)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := mockDatabase.NewMockDatabase(t)

			tt.setMockDB(mock)

			repo := database.NewEntryRepository(db)
			result, err := repo.FindByKeyPrefixAndAccountID(t.Context(), tt.inputKeyword, tt.inputAccountID)
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

func TestEntry_FindByVolumeIDAndAccountID(t *testing.T) {
	entry := &entity.Entry{
		ID:        uuid.New(),
		AccountID: uuid.New(),
		VolumeID:  uuid.New(),
		Key:       "test/sample.jpg",
		Size:      10000,
		Type:      "image/jpeg",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name           string
		inputVolumeID  uuid.UUID
		inputAccountID uuid.UUID
		inputPrefix    *string
		inputDepth     *uint64
		expectResult   []*entity.Entry
		expectError    error
		setMockDB      func(mock sqlmock.Sqlmock)
	}{
		{
			name:           "find all",
			inputVolumeID:  entry.VolumeID,
			inputAccountID: entry.AccountID,
			inputPrefix:    nil,
			inputDepth:     nil,
			expectResult:   []*entity.Entry{entry},
			expectError:    nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, account_id, volume_id, `key`, size, type, created_at, updated_at FROM entries WHERE volume_id = ? AND account_id = ?;")).
					WithArgs(entry.VolumeID, entry.AccountID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "account_id", "volume_id", "key", "size", "type", "created_at", "updated_at"}).AddRow(entry.ID, entry.AccountID, entry.VolumeID, entry.Key, entry.Size, entry.Type, entry.CreatedAt, entry.UpdatedAt)).
					WillReturnError(nil)
			},
		},
		{
			name:           "find by prefix",
			inputVolumeID:  entry.VolumeID,
			inputAccountID: entry.AccountID,
			inputPrefix:    types.ToPointer("test"),
			inputDepth:     nil,
			expectResult:   []*entity.Entry{entry},
			expectError:    nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, account_id, volume_id, `key`, size, type, created_at, updated_at FROM entries WHERE volume_id = ? AND account_id = ? AND `key` LIKE ?;")).
					WithArgs(entry.VolumeID, entry.AccountID, "test/%").
					WillReturnRows(sqlmock.NewRows([]string{"id", "account_id", "volume_id", "key", "size", "type", "created_at", "updated_at"}).AddRow(entry.ID, entry.AccountID, entry.VolumeID, entry.Key, entry.Size, entry.Type, entry.CreatedAt, entry.UpdatedAt)).
					WillReturnError(nil)
			},
		},
		{
			name:           "find by prefix with depth",
			inputVolumeID:  entry.VolumeID,
			inputAccountID: entry.AccountID,
			inputPrefix:    types.ToPointer("test"),
			inputDepth:     types.ToPointer(uint64(1)),
			expectResult:   []*entity.Entry{entry},
			expectError:    nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, account_id, volume_id, `key`, size, type, created_at, updated_at FROM entries WHERE volume_id = ? AND account_id = ? AND `key` LIKE ? AND LENGTH(`key`) - LENGTH(REPLACE(`key`, '/', '')) <= LENGTH(?) - LENGTH(REPLACE(?, '/', '')) + ?;")).
					WithArgs(entry.VolumeID, entry.AccountID, "test/%", "test", "test", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "account_id", "volume_id", "key", "size", "type", "created_at", "updated_at"}).AddRow(entry.ID, entry.AccountID, entry.VolumeID, entry.Key, entry.Size, entry.Type, entry.CreatedAt, entry.UpdatedAt)).
					WillReturnError(nil)
			},
		},
		{
			name:           "invalid argument",
			inputVolumeID:  entry.VolumeID,
			inputAccountID: entry.AccountID,
			inputPrefix:    nil,
			inputDepth:     types.ToPointer(uint64(1)),
			expectResult:   nil,
			expectError:    database.ErrInvalidArguments,
			setMockDB:      func(sqlmock.Sqlmock) {},
		},
		{
			name:           "not found",
			inputVolumeID:  entry.VolumeID,
			inputAccountID: entry.AccountID,
			inputPrefix:    nil,
			inputDepth:     nil,
			expectResult:   []*entity.Entry{},
			expectError:    nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, account_id, volume_id, `key`, size, type, created_at, updated_at FROM entries WHERE volume_id = ? AND account_id = ?;")).
					WithArgs(entry.VolumeID, entry.AccountID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "account_id", "volume_id", "key", "size", "type", "created_at", "updated_at"})).
					WillReturnError(nil)
			},
		},
		{
			name:           "find error",
			inputVolumeID:  entry.VolumeID,
			inputAccountID: entry.AccountID,
			inputPrefix:    nil,
			inputDepth:     nil,
			expectResult:   nil,
			expectError:    sql.ErrConnDone,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, account_id, volume_id, `key`, size, type, created_at, updated_at FROM entries WHERE volume_id = ? AND account_id = ?;")).
					WithArgs(entry.VolumeID, entry.AccountID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "account_id", "volume_id", "key", "size", "type", "created_at", "updated_at"})).
					WillReturnError(sql.ErrConnDone)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := mockDatabase.NewMockDatabase(t)

			tt.setMockDB(mock)

			repo := database.NewEntryRepository(db)
			result, err := repo.FindByVolumeIDAndAccountID(t.Context(), tt.inputVolumeID, tt.inputAccountID, tt.inputPrefix, tt.inputDepth)
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
