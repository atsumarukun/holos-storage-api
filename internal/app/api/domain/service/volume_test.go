package service_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/entity"
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
		setMockVolumeRepo func(context.Context, *mockRepository.MockVolumeRepository)
	}{
		{
			name:        "not exists",
			inputVolume: volume,
			expectError: nil,
			setMockVolumeRepo: func(ctx context.Context, volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindOneByName(ctx, volume.Name).
					Return(nil, nil).
					Times(1)
			},
		},
		{
			name:        "exists",
			inputVolume: volume,
			expectError: service.ErrVolumeAlreadyExists,
			setMockVolumeRepo: func(ctx context.Context, volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindOneByName(ctx, volume.Name).
					Return(volume, nil).
					Times(1)
			},
		},
		{
			name:              "volume is nil",
			inputVolume:       nil,
			expectError:       service.ErrRequiredVolume,
			setMockVolumeRepo: func(context.Context, *mockRepository.MockVolumeRepository) {},
		},
		{
			name:        "find error",
			inputVolume: volume,
			expectError: sql.ErrConnDone,
			setMockVolumeRepo: func(ctx context.Context, volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindOneByName(ctx, volume.Name).
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
			tt.setMockVolumeRepo(ctx, volumeRepo)

			serv := service.NewVolumeService(volumeRepo, nil)
			if err := serv.Exists(ctx, tt.inputVolume); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
		})
	}
}
