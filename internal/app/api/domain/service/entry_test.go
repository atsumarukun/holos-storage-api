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
		IsPublic:  false,
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
					FindOneByKeyAndVolumeIDAndAccountID(ctx, entry.Key, entry.VolumeID, entry.AccountID).
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
					FindOneByKeyAndVolumeIDAndAccountID(ctx, entry.Key, entry.VolumeID, entry.AccountID).
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
					FindOneByKeyAndVolumeIDAndAccountID(ctx, entry.Key, entry.VolumeID, entry.AccountID).
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

			serv := service.NewEntryService(entryRepo, nil)
			if err := serv.Exists(ctx, tt.inputEntry); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
		})
	}
}

func TestEntry_Create(t *testing.T) {
	accountID := uuid.New()
	volume := &entity.Volume{
		ID:        uuid.New(),
		AccountID: accountID,
		Name:      "name",
		IsPublic:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	entry := &entity.Entry{
		ID:        uuid.New(),
		AccountID: accountID,
		VolumeID:  volume.ID,
		Key:       "test/sample.txt",
		Size:      10000,
		Type:      "text/plain",
		IsPublic:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	parentEntry := &entity.Entry{
		ID:        uuid.New(),
		AccountID: accountID,
		VolumeID:  volume.ID,
		Key:       "test",
		Size:      0,
		Type:      "folder",
		IsPublic:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name             string
		inputVolume      *entity.Volume
		inputEntry       *entity.Entry
		inputBody        io.Reader
		expectError      error
		setMockEntryRepo func(context.Context, *mockRepository.MockEntryRepository)
		setMockBodyRepo  func(context.Context, *mockRepository.MockBodyRepository)
	}{
		{
			name:        "success",
			inputVolume: volume,
			inputEntry:  entry,
			inputBody:   bytes.NewBufferString("test"),
			expectError: nil,
			setMockEntryRepo: func(ctx context.Context, entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindOneByKeyAndVolumeIDAndAccountID(ctx, parentEntry.Key, parentEntry.VolumeID, parentEntry.AccountID).
					Return(nil, nil).
					AnyTimes()
				entryRepo.
					EXPECT().
					Create(ctx, gomock.Any()).
					Return(nil).
					AnyTimes()
			},
			setMockBodyRepo: func(_ context.Context, bodyRepo *mockRepository.MockBodyRepository) {
				bodyRepo.
					EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name:             "volume is nil",
			inputVolume:      nil,
			inputEntry:       entry,
			inputBody:        bytes.NewBufferString("test"),
			expectError:      service.ErrRequiredVolume,
			setMockEntryRepo: func(context.Context, *mockRepository.MockEntryRepository) {},
			setMockBodyRepo:  func(context.Context, *mockRepository.MockBodyRepository) {},
		},
		{
			name:             "entry is nil",
			inputVolume:      volume,
			inputEntry:       nil,
			inputBody:        bytes.NewBufferString("test"),
			expectError:      service.ErrRequiredEntry,
			setMockEntryRepo: func(context.Context, *mockRepository.MockEntryRepository) {},
			setMockBodyRepo:  func(context.Context, *mockRepository.MockBodyRepository) {},
		},
		{
			name:        "body is nil",
			inputVolume: volume,
			inputEntry:  entry,
			inputBody:   nil,
			expectError: nil,
			setMockEntryRepo: func(ctx context.Context, entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindOneByKeyAndVolumeIDAndAccountID(ctx, parentEntry.Key, parentEntry.VolumeID, parentEntry.AccountID).
					Return(parentEntry, nil).
					AnyTimes()
				entryRepo.
					EXPECT().
					Create(ctx, gomock.Any()).
					Return(nil).
					AnyTimes()
			},
			setMockBodyRepo: func(_ context.Context, bodyRepo *mockRepository.MockBodyRepository) {
				bodyRepo.
					EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name:        "parent entry already exists",
			inputVolume: volume,
			inputEntry:  entry,
			inputBody:   bytes.NewBufferString("test"),
			expectError: nil,
			setMockEntryRepo: func(ctx context.Context, entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindOneByKeyAndVolumeIDAndAccountID(ctx, parentEntry.Key, parentEntry.VolumeID, parentEntry.AccountID).
					Return(parentEntry, nil).
					AnyTimes()
				entryRepo.
					EXPECT().
					Create(ctx, gomock.Any()).
					Return(nil).
					AnyTimes()
			},
			setMockBodyRepo: func(_ context.Context, bodyRepo *mockRepository.MockBodyRepository) {
				bodyRepo.
					EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name:        "find entry Error",
			inputVolume: volume,
			inputEntry:  entry,
			inputBody:   bytes.NewBufferString("test"),
			expectError: sql.ErrConnDone,
			setMockEntryRepo: func(ctx context.Context, entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindOneByKeyAndVolumeIDAndAccountID(ctx, parentEntry.Key, parentEntry.VolumeID, parentEntry.AccountID).
					Return(nil, sql.ErrConnDone).
					AnyTimes()
			},
			setMockBodyRepo: func(context.Context, *mockRepository.MockBodyRepository) {},
		},
		{
			name:        "create entry error",
			inputVolume: volume,
			inputEntry:  entry,
			inputBody:   bytes.NewBufferString("test"),
			expectError: sql.ErrConnDone,
			setMockEntryRepo: func(ctx context.Context, entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindOneByKeyAndVolumeIDAndAccountID(ctx, parentEntry.Key, parentEntry.VolumeID, parentEntry.AccountID).
					Return(nil, nil).
					AnyTimes()
				entryRepo.
					EXPECT().
					Create(ctx, gomock.Any()).
					Return(sql.ErrConnDone).
					AnyTimes()
			},
			setMockBodyRepo: func(context.Context, *mockRepository.MockBodyRepository) {},
		},
		{
			name:        "create body error",
			inputVolume: volume,
			inputEntry:  entry,
			inputBody:   bytes.NewBufferString("test"),
			expectError: io.ErrNoProgress,
			setMockEntryRepo: func(ctx context.Context, entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindOneByKeyAndVolumeIDAndAccountID(ctx, parentEntry.Key, parentEntry.VolumeID, parentEntry.AccountID).
					Return(nil, nil).
					AnyTimes()
				entryRepo.
					EXPECT().
					Create(ctx, gomock.Any()).
					Return(nil).
					AnyTimes()
			},
			setMockBodyRepo: func(_ context.Context, bodyRepo *mockRepository.MockBodyRepository) {
				bodyRepo.
					EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(io.ErrNoProgress).
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

			bodyRepo := mockRepository.NewMockBodyRepository(ctrl)
			tt.setMockBodyRepo(ctx, bodyRepo)

			serv := service.NewEntryService(entryRepo, bodyRepo)
			if err := serv.Create(ctx, tt.inputVolume, tt.inputEntry, tt.inputBody); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
		})
	}
}
