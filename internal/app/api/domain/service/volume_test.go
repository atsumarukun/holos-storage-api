package service_test

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/repository"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/service"
	mockRepository "github.com/atsumarukun/holos-storage-api/test/mock/domain/repository"
)

func TestVolume_Exists(t *testing.T) {
	volume := &entity.Volume{
		ID:        uuid.New(),
		AccountID: uuid.New(),
		Name:      "name",
		IsPublic:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name              string
		inputVolume       *entity.Volume
		expectError       error
		setMockVolumeRepo func(*mockRepository.MockVolumeRepository)
	}{
		{
			name:        "not exists",
			inputVolume: volume,
			expectError: nil,
			setMockVolumeRepo: func(volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindOneByName(gomock.Any(), gomock.Any()).
					Return(nil, repository.ErrVolumeNotFound).
					Times(1)
			},
		},
		{
			name:        "exists",
			inputVolume: volume,
			expectError: service.ErrVolumeAlreadyExists,
			setMockVolumeRepo: func(volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindOneByName(gomock.Any(), gomock.Any()).
					Return(volume, nil).
					Times(1)
			},
		},
		{
			name:              "volume is nil",
			inputVolume:       nil,
			expectError:       service.ErrRequiredVolume,
			setMockVolumeRepo: func(*mockRepository.MockVolumeRepository) {},
		},
		{
			name:        "find error",
			inputVolume: volume,
			expectError: sql.ErrConnDone,
			setMockVolumeRepo: func(volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindOneByName(gomock.Any(), gomock.Any()).
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

			volumeRepo := mockRepository.NewMockVolumeRepository(ctrl)
			tt.setMockVolumeRepo(volumeRepo)

			serv := service.NewVolumeService(volumeRepo, nil)
			if err := serv.Exists(ctx, tt.inputVolume); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
		})
	}
}

func TestVolume_CanDelete(t *testing.T) {
	volume := &entity.Volume{
		ID:        uuid.New(),
		AccountID: uuid.New(),
		Name:      "name",
		IsPublic:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	entry := &entity.Entry{
		ID:        uuid.New(),
		AccountID: volume.ID,
		VolumeID:  volume.AccountID,
		Key:       "test/sample.txt",
		Size:      10000,
		Type:      "text/plain",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name             string
		inputVolume      *entity.Volume
		expectError      error
		setMockEntryRepo func(*mockRepository.MockEntryRepository)
	}{
		{
			name:        "volume has not entries",
			inputVolume: volume,
			expectError: nil,
			setMockEntryRepo: func(entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindByVolumeIDAndAccountID(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, nil).
					Times(1)
			},
		},
		{
			name:        "volume has entries",
			inputVolume: volume,
			expectError: service.ErrVolumeHasEntries,
			setMockEntryRepo: func(entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindByVolumeIDAndAccountID(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]*entity.Entry{entry}, nil).
					Times(1)
			},
		},
		{
			name:             "volume is nil",
			inputVolume:      nil,
			expectError:      service.ErrRequiredVolume,
			setMockEntryRepo: func(*mockRepository.MockEntryRepository) {},
		},
		{
			name:        "find error",
			inputVolume: volume,
			expectError: sql.ErrConnDone,
			setMockEntryRepo: func(entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindByVolumeIDAndAccountID(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
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
			tt.setMockEntryRepo(entryRepo)

			serv := service.NewVolumeService(nil, entryRepo)
			if err := serv.CanDelete(ctx, tt.inputVolume); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
		})
	}
}
