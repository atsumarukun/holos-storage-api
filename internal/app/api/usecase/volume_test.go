package usecase_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/service"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/usecase"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/usecase/dto"
	"github.com/atsumarukun/holos-storage-api/test/mock/domain/repository"
	"github.com/atsumarukun/holos-storage-api/test/mock/domain/repository/pkg/transaction"
	mockService "github.com/atsumarukun/holos-storage-api/test/mock/domain/service"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func TestVolume_Create(t *testing.T) {
	accountID := uuid.New()
	volumeDTO := &dto.VolumeDTO{
		ID:        uuid.New(),
		AccountID: accountID,
		Name:      "name",
		IsPublic:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name                  string
		inputAccountID        uuid.UUID
		inputName             string
		inputIsPublic         bool
		expectResult          *dto.VolumeDTO
		expectError           error
		setMockTransactionObj func(context.Context, *transaction.MockTransactionObject)
		setMockVolumeRepo     func(context.Context, *repository.MockVolumeRepository)
		setMockVolumeServ     func(context.Context, *mockService.MockVolumeService)
	}{
		{
			name:           "success",
			inputAccountID: accountID,
			inputName:      "name",
			inputIsPublic:  false,
			expectResult:   volumeDTO,
			expectError:    nil,
			setMockTransactionObj: func(ctx context.Context, transactionObj *transaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockVolumeRepo: func(ctx context.Context, volumeRepo *repository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					Create(ctx, gomock.Any()).
					Return(nil).
					Times(1)
			},
			setMockVolumeServ: func(ctx context.Context, volumeServ *mockService.MockVolumeService) {
				volumeServ.
					EXPECT().
					Exists(ctx, gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name:                  "invalid name",
			inputAccountID:        accountID,
			inputName:             "",
			inputIsPublic:         false,
			expectResult:          nil,
			expectError:           entity.ErrShortVolumeName,
			setMockTransactionObj: func(context.Context, *transaction.MockTransactionObject) {},
			setMockVolumeRepo:     func(context.Context, *repository.MockVolumeRepository) {},
			setMockVolumeServ:     func(context.Context, *mockService.MockVolumeService) {},
		},
		{
			name:           "volume already exists",
			inputAccountID: accountID,
			inputName:      "name",
			inputIsPublic:  false,
			expectResult:   nil,
			expectError:    service.ErrVolumeAlreadyExists,
			setMockTransactionObj: func(ctx context.Context, transactionObj *transaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockVolumeRepo: func(context.Context, *repository.MockVolumeRepository) {},
			setMockVolumeServ: func(ctx context.Context, volumeServ *mockService.MockVolumeService) {
				volumeServ.
					EXPECT().
					Exists(ctx, gomock.Any()).
					Return(service.ErrVolumeAlreadyExists).
					Times(1)
			},
		},
		{
			name:           "create error",
			inputAccountID: accountID,
			inputName:      "name",
			inputIsPublic:  false,
			expectResult:   nil,
			expectError:    sql.ErrConnDone,
			setMockTransactionObj: func(ctx context.Context, transactionObj *transaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockVolumeRepo: func(ctx context.Context, volumeRepo *repository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					Create(ctx, gomock.Any()).
					Return(sql.ErrConnDone).
					Times(1)
			},
			setMockVolumeServ: func(ctx context.Context, volumeServ *mockService.MockVolumeService) {
				volumeServ.
					EXPECT().
					Exists(ctx, gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := t.Context()

			transactionObj := transaction.NewMockTransactionObject(ctrl)
			tt.setMockTransactionObj(ctx, transactionObj)

			volumeRepo := repository.NewMockVolumeRepository(ctrl)
			tt.setMockVolumeRepo(ctx, volumeRepo)

			volumeServ := mockService.NewMockVolumeService(ctrl)
			tt.setMockVolumeServ(ctx, volumeServ)

			uc := usecase.NewVolumeUsecase(transactionObj, volumeRepo, volumeServ)
			result, err := uc.Create(ctx, tt.inputAccountID, tt.inputName, tt.inputIsPublic)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			opts := cmp.Options{
				cmpopts.IgnoreFields(dto.VolumeDTO{}, "ID", "CreatedAt", "UpdatedAt"),
			}
			if diff := cmp.Diff(result, tt.expectResult, opts...); diff != "" {
				t.Error(diff)
			}
		})
	}
}
