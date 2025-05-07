package usecase_test

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/service"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/usecase"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/usecase/dto"
	mockRepository "github.com/atsumarukun/holos-storage-api/test/mock/domain/repository"
	mockTransaction "github.com/atsumarukun/holos-storage-api/test/mock/domain/repository/pkg/transaction"
	mockService "github.com/atsumarukun/holos-storage-api/test/mock/domain/service"
)

func TestEntry_Create(t *testing.T) {
	accountID := uuid.New()
	volumeID := uuid.New()
	volume := &entity.Volume{
		ID:        volumeID,
		AccountID: accountID,
		Name:      "name",
		IsPublic:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	entryDTO := &dto.EntryDTO{
		ID:        uuid.New(),
		AccountID: accountID,
		VolumeID:  volumeID,
		Key:       "test/sample.txt",
		Size:      4,
		Type:      "text/plain; charset=utf-8",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	folderEntryDTO := &dto.EntryDTO{
		ID:        uuid.New(),
		AccountID: accountID,
		VolumeID:  volumeID,
		Key:       "test/sample",
		Size:      0,
		Type:      "folder",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name                  string
		inputAccountID        uuid.UUID
		inputVolumeID         uuid.UUID
		inputKey              string
		inputSize             uint64
		inputBody             io.Reader
		expectResult          *dto.EntryDTO
		expectError           error
		setMockTransactionObj func(context.Context, *mockTransaction.MockTransactionObject)
		setMockVolumeRepo     func(context.Context, *mockRepository.MockVolumeRepository)
		setMockEntryServ      func(context.Context, *mockService.MockEntryService)
	}{
		{
			name:           "success",
			inputAccountID: accountID,
			inputVolumeID:  volumeID,
			inputKey:       "test/sample.txt",
			inputSize:      4,
			inputBody:      bytes.NewBufferString("test"),
			expectResult:   entryDTO,
			expectError:    nil,
			setMockTransactionObj: func(ctx context.Context, transactionObj *mockTransaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockVolumeRepo: func(ctx context.Context, volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindOneByIDAndAccountID(ctx, gomock.Any(), gomock.Any()).
					Return(volume, nil).
					Times(1)
			},
			setMockEntryServ: func(ctx context.Context, entryServ *mockService.MockEntryService) {
				entryServ.
					EXPECT().
					Exists(ctx, gomock.Any()).
					Return(nil).
					Times(1)
				entryServ.
					EXPECT().
					Create(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name:           "invalid key",
			inputAccountID: accountID,
			inputVolumeID:  volumeID,
			inputKey:       "",
			inputSize:      4,
			inputBody:      bytes.NewBufferString("test"),
			expectResult:   nil,
			expectError:    entity.ErrShortEntryKey,
			setMockTransactionObj: func(ctx context.Context, transactionObj *mockTransaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockVolumeRepo: func(ctx context.Context, volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindOneByIDAndAccountID(ctx, gomock.Any(), gomock.Any()).
					Return(volume, nil).
					Times(1)
			},
			setMockEntryServ: func(context.Context, *mockService.MockEntryService) {},
		},
		{
			name:           "body is nil",
			inputAccountID: accountID,
			inputVolumeID:  volumeID,
			inputKey:       "test/sample",
			inputSize:      0,
			inputBody:      nil,
			expectResult:   folderEntryDTO,
			expectError:    nil,
			setMockTransactionObj: func(ctx context.Context, transactionObj *mockTransaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockVolumeRepo: func(ctx context.Context, volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindOneByIDAndAccountID(ctx, gomock.Any(), gomock.Any()).
					Return(volume, nil).
					Times(1)
			},
			setMockEntryServ: func(ctx context.Context, entryServ *mockService.MockEntryService) {
				entryServ.
					EXPECT().
					Exists(ctx, gomock.Any()).
					Return(nil).
					Times(1)
				entryServ.
					EXPECT().
					Create(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name:           "entry already exists",
			inputAccountID: accountID,
			inputVolumeID:  volumeID,
			inputKey:       "test/sample.txt",
			inputSize:      4,
			inputBody:      bytes.NewBufferString("test"),
			expectResult:   nil,
			expectError:    service.ErrEntryAlreadyExists,
			setMockTransactionObj: func(ctx context.Context, transactionObj *mockTransaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockVolumeRepo: func(ctx context.Context, volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindOneByIDAndAccountID(ctx, gomock.Any(), gomock.Any()).
					Return(volume, nil).
					Times(1)
			},
			setMockEntryServ: func(ctx context.Context, entryServ *mockService.MockEntryService) {
				entryServ.
					EXPECT().
					Exists(ctx, gomock.Any()).
					Return(service.ErrEntryAlreadyExists).
					Times(1)
			},
		},
		{
			name:           "volume not found",
			inputAccountID: accountID,
			inputVolumeID:  volumeID,
			inputKey:       "test/sample.txt",
			inputSize:      4,
			inputBody:      bytes.NewBufferString("test"),
			expectResult:   nil,
			expectError:    usecase.ErrVolumeNotFound,
			setMockTransactionObj: func(ctx context.Context, transactionObj *mockTransaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockVolumeRepo: func(ctx context.Context, volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindOneByIDAndAccountID(ctx, gomock.Any(), gomock.Any()).
					Return(nil, nil).
					Times(1)
			},
			setMockEntryServ: func(context.Context, *mockService.MockEntryService) {},
		},
		{
			name:           "create error",
			inputAccountID: accountID,
			inputVolumeID:  volumeID,
			inputKey:       "test/sample.txt",
			inputSize:      4,
			inputBody:      bytes.NewBufferString("test"),
			expectResult:   nil,
			expectError:    sql.ErrConnDone,
			setMockTransactionObj: func(ctx context.Context, transactionObj *mockTransaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockVolumeRepo: func(ctx context.Context, volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindOneByIDAndAccountID(ctx, gomock.Any(), gomock.Any()).
					Return(volume, nil).
					Times(1)
			},
			setMockEntryServ: func(ctx context.Context, entryServ *mockService.MockEntryService) {
				entryServ.
					EXPECT().
					Exists(ctx, gomock.Any()).
					Return(nil).
					Times(1)
				entryServ.
					EXPECT().
					Create(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(sql.ErrConnDone).
					Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := t.Context()

			transactionObj := mockTransaction.NewMockTransactionObject(ctrl)
			tt.setMockTransactionObj(ctx, transactionObj)

			volumeRepo := mockRepository.NewMockVolumeRepository(ctrl)
			tt.setMockVolumeRepo(ctx, volumeRepo)

			entryServ := mockService.NewMockEntryService(ctrl)
			tt.setMockEntryServ(ctx, entryServ)

			uc := usecase.NewEntryUsecase(transactionObj, nil, volumeRepo, entryServ)
			result, err := uc.Create(ctx, tt.inputAccountID, tt.inputVolumeID, tt.inputKey, tt.inputSize, tt.inputBody)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			opts := cmp.Options{
				cmpopts.IgnoreFields(dto.EntryDTO{}, "ID", "CreatedAt", "UpdatedAt"),
			}
			if diff := cmp.Diff(result, tt.expectResult, opts...); diff != "" {
				t.Error(diff)
			}
		})
	}
}
