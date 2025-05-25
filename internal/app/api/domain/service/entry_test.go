package service_test

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/service"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/pkg/types"
	mockRepository "github.com/atsumarukun/holos-storage-api/test/mock/domain/repository"
)

func TestEntry_Exists(t *testing.T) {
	entry := &entity.Entry{
		ID:        uuid.New(),
		AccountID: uuid.New(),
		VolumeID:  uuid.New(),
		Key:       "test/sample.txt",
		Size:      10000,
		Type:      "text/plain",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name             string
		inputEntry       *entity.Entry
		expectError      error
		setMockEntryRepo func(context.Context, *mockRepository.MockEntryRepository)
	}{
		{
			name:        "not exists",
			inputEntry:  entry,
			expectError: nil,
			setMockEntryRepo: func(ctx context.Context, entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindOneByKeyAndVolumeID(ctx, entry.Key, entry.VolumeID).
					Return(nil, nil).
					Times(1)
			},
		},
		{
			name:        "exists",
			inputEntry:  entry,
			expectError: service.ErrEntryAlreadyExists,
			setMockEntryRepo: func(ctx context.Context, entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindOneByKeyAndVolumeID(ctx, entry.Key, entry.VolumeID).
					Return(entry, nil).
					Times(1)
			},
		},
		{
			name:             "entry is nil",
			inputEntry:       nil,
			expectError:      service.ErrRequiredEntry,
			setMockEntryRepo: func(context.Context, *mockRepository.MockEntryRepository) {},
		},
		{
			name:        "find error",
			inputEntry:  entry,
			expectError: sql.ErrConnDone,
			setMockEntryRepo: func(ctx context.Context, entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindOneByKeyAndVolumeID(ctx, entry.Key, entry.VolumeID).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := t.Context()

			entryRepo := mockRepository.NewMockEntryRepository(ctrl)
			tt.setMockEntryRepo(ctx, entryRepo)

			serv := service.NewEntryService(entryRepo)
			if err := serv.Exists(ctx, tt.inputEntry); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
		})
	}
}

func TestEntry_Create(t *testing.T) {
	accountID := uuid.New()
	volumeID := uuid.New()
	entry := &entity.Entry{
		ID:        uuid.New(),
		AccountID: accountID,
		VolumeID:  volumeID,
		Key:       "test/sample.txt",
		Size:      10000,
		Type:      "text/plain",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	parentEntry := &entity.Entry{
		ID:        uuid.New(),
		AccountID: accountID,
		VolumeID:  volumeID,
		Key:       "test",
		Size:      0,
		Type:      "folder",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name             string
		inputEntry       *entity.Entry
		inputBody        io.Reader
		expectError      error
		setMockEntryRepo func(context.Context, *mockRepository.MockEntryRepository)
	}{
		{
			name:        "success",
			inputEntry:  entry,
			inputBody:   bytes.NewBufferString("test"),
			expectError: nil,
			setMockEntryRepo: func(ctx context.Context, entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindOneByKeyAndVolumeID(ctx, parentEntry.Key, parentEntry.VolumeID).
					Return(nil, nil).
					AnyTimes()
				entryRepo.
					EXPECT().
					Create(ctx, gomock.Any()).
					Return(nil).
					AnyTimes()
			},
		},
		{
			name:             "entry is nil",
			inputEntry:       nil,
			inputBody:        bytes.NewBufferString("test"),
			expectError:      service.ErrRequiredEntry,
			setMockEntryRepo: func(context.Context, *mockRepository.MockEntryRepository) {},
		},
		{
			name:        "parent entry already exists",
			inputEntry:  entry,
			inputBody:   bytes.NewBufferString("test"),
			expectError: nil,
			setMockEntryRepo: func(ctx context.Context, entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindOneByKeyAndVolumeID(ctx, parentEntry.Key, parentEntry.VolumeID).
					Return(parentEntry, nil).
					AnyTimes()
				entryRepo.
					EXPECT().
					Create(ctx, gomock.Any()).
					Return(nil).
					AnyTimes()
			},
		},
		{
			name:        "find entry error",
			inputEntry:  entry,
			inputBody:   bytes.NewBufferString("test"),
			expectError: sql.ErrConnDone,
			setMockEntryRepo: func(ctx context.Context, entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindOneByKeyAndVolumeID(ctx, parentEntry.Key, parentEntry.VolumeID).
					Return(nil, sql.ErrConnDone).
					AnyTimes()
			},
		},
		{
			name:        "create entry error",
			inputEntry:  entry,
			inputBody:   bytes.NewBufferString("test"),
			expectError: sql.ErrConnDone,
			setMockEntryRepo: func(ctx context.Context, entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindOneByKeyAndVolumeID(ctx, parentEntry.Key, parentEntry.VolumeID).
					Return(nil, nil).
					AnyTimes()
				entryRepo.
					EXPECT().
					Create(ctx, gomock.Any()).
					Return(sql.ErrConnDone).
					AnyTimes()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := t.Context()

			entryRepo := mockRepository.NewMockEntryRepository(ctrl)
			tt.setMockEntryRepo(ctx, entryRepo)

			serv := service.NewEntryService(entryRepo)
			if err := serv.Create(ctx, tt.inputEntry, tt.inputBody); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
		})
	}
}

func TestEntry_Update(t *testing.T) {
	accountID := uuid.New()
	volumeID := uuid.New()
	entry := &entity.Entry{
		ID:        uuid.New(),
		AccountID: accountID,
		VolumeID:  volumeID,
		Key:       "test",
		Size:      0,
		Type:      "folder",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	childEntry := &entity.Entry{
		ID:        uuid.New(),
		AccountID: accountID,
		VolumeID:  volumeID,
		Key:       "test/sample.txt",
		Size:      10000,
		Type:      "text/plain",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name             string
		inputEntry       *entity.Entry
		inputSrc         string
		expectError      error
		setMockEntryRepo func(context.Context, *mockRepository.MockEntryRepository)
	}{
		{
			name:        "success",
			inputEntry:  entry,
			inputSrc:    "update",
			expectError: nil,
			setMockEntryRepo: func(ctx context.Context, entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindByVolumeIDAndAccountID(ctx, entry.VolumeID, entry.AccountID, types.ToPointer("update"), nil).
					Return([]*entity.Entry{childEntry}, nil).
					AnyTimes()
				entryRepo.
					EXPECT().
					Update(ctx, gomock.Any()).
					Return(nil).
					AnyTimes()
			},
		},
		{
			name:             "entry is nil",
			inputEntry:       nil,
			inputSrc:         "update",
			expectError:      service.ErrRequiredEntry,
			setMockEntryRepo: func(context.Context, *mockRepository.MockEntryRepository) {},
		},
		{
			name:        "find entry error",
			inputEntry:  entry,
			inputSrc:    "update",
			expectError: sql.ErrConnDone,
			setMockEntryRepo: func(ctx context.Context, entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindByVolumeIDAndAccountID(ctx, entry.VolumeID, entry.AccountID, types.ToPointer("update"), nil).
					Return(nil, sql.ErrConnDone).
					AnyTimes()
			},
		},
		{
			name:        "update entry error",
			inputEntry:  entry,
			inputSrc:    "update",
			expectError: sql.ErrConnDone,
			setMockEntryRepo: func(ctx context.Context, entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindByVolumeIDAndAccountID(ctx, entry.VolumeID, entry.AccountID, types.ToPointer("update"), nil).
					Return([]*entity.Entry{childEntry}, nil).
					AnyTimes()
				entryRepo.
					EXPECT().
					Update(ctx, gomock.Any()).
					Return(sql.ErrConnDone).
					AnyTimes()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := t.Context()

			entryRepo := mockRepository.NewMockEntryRepository(ctrl)
			tt.setMockEntryRepo(ctx, entryRepo)

			serv := service.NewEntryService(entryRepo)
			if err := serv.Update(ctx, tt.inputEntry, tt.inputSrc); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
		})
	}
}

func TestEntry_Delete(t *testing.T) {
	accountID := uuid.New()
	volumeID := uuid.New()
	entry := &entity.Entry{
		ID:        uuid.New(),
		AccountID: accountID,
		VolumeID:  volumeID,
		Key:       "test",
		Size:      0,
		Type:      "folder",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	childEntry := &entity.Entry{
		ID:        uuid.New(),
		AccountID: accountID,
		VolumeID:  volumeID,
		Key:       "test/sample.txt",
		Size:      10000,
		Type:      "text/plain",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name             string
		inputEntry       *entity.Entry
		expectError      error
		setMockEntryRepo func(context.Context, *mockRepository.MockEntryRepository)
	}{
		{
			name:        "success",
			inputEntry:  entry,
			expectError: nil,
			setMockEntryRepo: func(ctx context.Context, entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindByVolumeIDAndAccountID(ctx, entry.VolumeID, entry.AccountID, types.ToPointer("test"), nil).
					Return([]*entity.Entry{childEntry}, nil).
					AnyTimes()
				entryRepo.
					EXPECT().
					Delete(ctx, gomock.Any()).
					Return(nil).
					AnyTimes()
			},
		},
		{
			name:             "entry is nil",
			inputEntry:       nil,
			expectError:      service.ErrRequiredEntry,
			setMockEntryRepo: func(context.Context, *mockRepository.MockEntryRepository) {},
		},
		{
			name:        "find entry error",
			inputEntry:  entry,
			expectError: sql.ErrConnDone,
			setMockEntryRepo: func(ctx context.Context, entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindByVolumeIDAndAccountID(ctx, entry.VolumeID, entry.AccountID, types.ToPointer("test"), nil).
					Return(nil, sql.ErrConnDone).
					AnyTimes()
			},
		},
		{
			name:        "delete entry error",
			inputEntry:  entry,
			expectError: sql.ErrConnDone,
			setMockEntryRepo: func(ctx context.Context, entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindByVolumeIDAndAccountID(ctx, entry.VolumeID, entry.AccountID, types.ToPointer("test"), nil).
					Return([]*entity.Entry{childEntry}, nil).
					AnyTimes()
				entryRepo.
					EXPECT().
					Delete(ctx, gomock.Any()).
					Return(sql.ErrConnDone).
					AnyTimes()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := t.Context()

			entryRepo := mockRepository.NewMockEntryRepository(ctrl)
			tt.setMockEntryRepo(ctx, entryRepo)

			serv := service.NewEntryService(entryRepo)
			if err := serv.Delete(ctx, tt.inputEntry); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
		})
	}
}
