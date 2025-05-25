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
	"github.com/spf13/afero"
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
		inputVolumeName       string
		inputKey              string
		inputSize             uint64
		inputBody             io.Reader
		expectResult          *dto.EntryDTO
		expectError           error
		setMockTransactionObj func(context.Context, *mockTransaction.MockTransactionObject)
		setMockEntryRepo      func(context.Context, *mockRepository.MockEntryRepository)
		setMockBodyRepo       func(context.Context, *mockRepository.MockBodyRepository)
		setMockVolumeRepo     func(context.Context, *mockRepository.MockVolumeRepository)
		setMockEntryServ      func(context.Context, *mockService.MockEntryService)
	}{
		{
			name:            "success",
			inputAccountID:  accountID,
			inputVolumeName: volume.Name,
			inputKey:        "test/sample.txt",
			inputSize:       4,
			inputBody:       bytes.NewBufferString("test"),
			expectResult:    entryDTO,
			expectError:     nil,
			setMockTransactionObj: func(ctx context.Context, transactionObj *mockTransaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockEntryRepo: func(ctx context.Context, entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					Create(ctx, gomock.Any()).
					Return(nil).
					Times(1)
			},
			setMockBodyRepo: func(_ context.Context, bodyRepo *mockRepository.MockBodyRepository) {
				bodyRepo.
					EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
			setMockVolumeRepo: func(ctx context.Context, volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindOneByNameAndAccountID(ctx, gomock.Any(), gomock.Any()).
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
					CreateAncestors(ctx, gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name:            "invalid key",
			inputAccountID:  accountID,
			inputVolumeName: volume.Name,
			inputKey:        "",
			inputSize:       4,
			inputBody:       bytes.NewBufferString("test"),
			expectResult:    nil,
			expectError:     entity.ErrShortEntryKey,
			setMockTransactionObj: func(ctx context.Context, transactionObj *mockTransaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockEntryRepo: func(context.Context, *mockRepository.MockEntryRepository) {},
			setMockBodyRepo:  func(context.Context, *mockRepository.MockBodyRepository) {},
			setMockVolumeRepo: func(ctx context.Context, volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindOneByNameAndAccountID(ctx, gomock.Any(), gomock.Any()).
					Return(volume, nil).
					Times(1)
			},
			setMockEntryServ: func(context.Context, *mockService.MockEntryService) {},
		},
		{
			name:            "body is nil",
			inputAccountID:  accountID,
			inputVolumeName: volume.Name,
			inputKey:        "test/sample",
			inputSize:       0,
			inputBody:       nil,
			expectResult:    folderEntryDTO,
			expectError:     nil,
			setMockTransactionObj: func(ctx context.Context, transactionObj *mockTransaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockEntryRepo: func(ctx context.Context, entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					Create(ctx, gomock.Any()).
					Return(nil).
					Times(1)
			},
			setMockBodyRepo: func(_ context.Context, bodyRepo *mockRepository.MockBodyRepository) {
				bodyRepo.
					EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
			setMockVolumeRepo: func(ctx context.Context, volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindOneByNameAndAccountID(ctx, gomock.Any(), gomock.Any()).
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
					CreateAncestors(ctx, gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name:            "entry already exists",
			inputAccountID:  accountID,
			inputVolumeName: volume.Name,
			inputKey:        "test/sample.txt",
			inputSize:       4,
			inputBody:       bytes.NewBufferString("test"),
			expectResult:    nil,
			expectError:     service.ErrEntryAlreadyExists,
			setMockTransactionObj: func(ctx context.Context, transactionObj *mockTransaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockEntryRepo: func(context.Context, *mockRepository.MockEntryRepository) {},
			setMockBodyRepo:  func(context.Context, *mockRepository.MockBodyRepository) {},
			setMockVolumeRepo: func(ctx context.Context, volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindOneByNameAndAccountID(ctx, gomock.Any(), gomock.Any()).
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
			name:            "create ancestors entries error",
			inputAccountID:  accountID,
			inputVolumeName: volume.Name,
			inputKey:        "test/sample.txt",
			inputSize:       4,
			inputBody:       bytes.NewBufferString("test"),
			expectResult:    nil,
			expectError:     sql.ErrConnDone,
			setMockTransactionObj: func(ctx context.Context, transactionObj *mockTransaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockEntryRepo: func(context.Context, *mockRepository.MockEntryRepository) {},
			setMockBodyRepo:  func(context.Context, *mockRepository.MockBodyRepository) {},
			setMockVolumeRepo: func(ctx context.Context, volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindOneByNameAndAccountID(ctx, gomock.Any(), gomock.Any()).
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
					CreateAncestors(ctx, gomock.Any()).
					Return(sql.ErrConnDone).
					Times(1)
			},
		},
		{
			name:            "create entry error",
			inputAccountID:  accountID,
			inputVolumeName: volume.Name,
			inputKey:        "test/sample.txt",
			inputSize:       4,
			inputBody:       bytes.NewBufferString("test"),
			expectResult:    nil,
			expectError:     sql.ErrConnDone,
			setMockTransactionObj: func(ctx context.Context, transactionObj *mockTransaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockEntryRepo: func(ctx context.Context, entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					Create(ctx, gomock.Any()).
					Return(sql.ErrConnDone).
					Times(1)
			},
			setMockBodyRepo: func(context.Context, *mockRepository.MockBodyRepository) {},
			setMockVolumeRepo: func(ctx context.Context, volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindOneByNameAndAccountID(ctx, gomock.Any(), gomock.Any()).
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
					CreateAncestors(ctx, gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name:            "create body error",
			inputAccountID:  accountID,
			inputVolumeName: volume.Name,
			inputKey:        "test/sample.txt",
			inputSize:       4,
			inputBody:       bytes.NewBufferString("test"),
			expectResult:    nil,
			expectError:     io.ErrNoProgress,
			setMockTransactionObj: func(ctx context.Context, transactionObj *mockTransaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockEntryRepo: func(ctx context.Context, entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					Create(ctx, gomock.Any()).
					Return(nil).
					Times(1)
			},
			setMockBodyRepo: func(_ context.Context, bodyRepo *mockRepository.MockBodyRepository) {
				bodyRepo.
					EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(io.ErrNoProgress).
					Times(1)
			},
			setMockVolumeRepo: func(ctx context.Context, volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindOneByNameAndAccountID(ctx, gomock.Any(), gomock.Any()).
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
					CreateAncestors(ctx, gomock.Any()).
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

			transactionObj := mockTransaction.NewMockTransactionObject(ctrl)
			tt.setMockTransactionObj(ctx, transactionObj)

			entryRepo := mockRepository.NewMockEntryRepository(ctrl)
			tt.setMockEntryRepo(ctx, entryRepo)

			bodyRepo := mockRepository.NewMockBodyRepository(ctrl)
			tt.setMockBodyRepo(ctx, bodyRepo)

			volumeRepo := mockRepository.NewMockVolumeRepository(ctrl)
			tt.setMockVolumeRepo(ctx, volumeRepo)

			entryServ := mockService.NewMockEntryService(ctrl)
			tt.setMockEntryServ(ctx, entryServ)

			uc := usecase.NewEntryUsecase(transactionObj, entryRepo, bodyRepo, volumeRepo, entryServ)
			result, err := uc.Create(ctx, tt.inputAccountID, tt.inputVolumeName, tt.inputKey, tt.inputSize, tt.inputBody)
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

func TestEntry_Update(t *testing.T) {
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
		AccountID: volume.AccountID,
		VolumeID:  volume.ID,
		Key:       "test/sample.txt",
		Size:      4,
		Type:      "text/plain; charset=utf-8",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	entryDTO := &dto.EntryDTO{
		ID:        entry.ID,
		AccountID: entry.AccountID,
		VolumeID:  entry.VolumeID,
		Key:       "update/sample.txt",
		Size:      entry.Size,
		Type:      entry.Type,
		CreatedAt: entry.CreatedAt,
		UpdatedAt: entry.UpdatedAt,
	}

	tests := []struct {
		name                  string
		inputAccountID        uuid.UUID
		inputVolumeName       string
		inputKey              string
		inputNewKey           string
		expectResult          *dto.EntryDTO
		expectError           error
		setMockTransactionObj func(context.Context, *mockTransaction.MockTransactionObject)
		setMockEntryRepo      func(context.Context, *mockRepository.MockEntryRepository)
		setMockBodyRepo       func(context.Context, *mockRepository.MockBodyRepository)
		setMockVolumeRepo     func(context.Context, *mockRepository.MockVolumeRepository)
		setMockEntryServ      func(context.Context, *mockService.MockEntryService)
	}{
		{
			name:            "success",
			inputAccountID:  entry.AccountID,
			inputVolumeName: "volume",
			inputKey:        "test/sample.txt",
			inputNewKey:     "update/sample.txt",
			expectResult:    entryDTO,
			expectError:     nil,
			setMockTransactionObj: func(ctx context.Context, transactionObj *mockTransaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockEntryRepo: func(ctx context.Context, entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindOneByKeyAndVolumeIDAndAccountID(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(entry, nil).
					Times(1)
				entryRepo.
					EXPECT().
					Update(ctx, gomock.Any()).
					Return(nil).
					Times(1)
			},
			setMockBodyRepo: func(_ context.Context, bodyRepo *mockRepository.MockBodyRepository) {
				bodyRepo.
					EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
			setMockVolumeRepo: func(ctx context.Context, volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindOneByNameAndAccountID(ctx, gomock.Any(), gomock.Any()).
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
					CreateAncestors(ctx, gomock.Any()).
					Return(nil).
					Times(1)
				entryServ.
					EXPECT().
					UpdateDescendants(ctx, gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name:            "invalid update key",
			inputAccountID:  entry.AccountID,
			inputVolumeName: "volume",
			inputKey:        "test/sample.txt",
			inputNewKey:     "",
			expectResult:    nil,
			expectError:     entity.ErrShortEntryKey,
			setMockTransactionObj: func(ctx context.Context, transactionObj *mockTransaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockEntryRepo: func(ctx context.Context, entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindOneByKeyAndVolumeIDAndAccountID(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(entry, nil).
					Times(1)
			},
			setMockBodyRepo: func(context.Context, *mockRepository.MockBodyRepository) {},
			setMockVolumeRepo: func(ctx context.Context, volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindOneByNameAndAccountID(ctx, gomock.Any(), gomock.Any()).
					Return(volume, nil).
					Times(1)
			},
			setMockEntryServ: func(context.Context, *mockService.MockEntryService) {},
		},
		{
			name:            "entry already exists",
			inputAccountID:  entry.AccountID,
			inputVolumeName: "volume",
			inputKey:        "test/sample.txt",
			inputNewKey:     "update/sample.txt",
			expectResult:    nil,
			expectError:     service.ErrEntryAlreadyExists,
			setMockTransactionObj: func(ctx context.Context, transactionObj *mockTransaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockEntryRepo: func(ctx context.Context, entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindOneByKeyAndVolumeIDAndAccountID(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(entry, nil).
					Times(1)
			},
			setMockBodyRepo: func(context.Context, *mockRepository.MockBodyRepository) {},
			setMockVolumeRepo: func(ctx context.Context, volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindOneByNameAndAccountID(ctx, gomock.Any(), gomock.Any()).
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
			name:            "create ancestor entries error",
			inputAccountID:  entry.AccountID,
			inputVolumeName: "volume",
			inputKey:        "test/sample.txt",
			inputNewKey:     "update/sample.txt",
			expectResult:    nil,
			expectError:     sql.ErrConnDone,
			setMockTransactionObj: func(ctx context.Context, transactionObj *mockTransaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockEntryRepo: func(ctx context.Context, entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindOneByKeyAndVolumeIDAndAccountID(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(entry, nil).
					Times(1)
			},
			setMockBodyRepo: func(context.Context, *mockRepository.MockBodyRepository) {},
			setMockVolumeRepo: func(ctx context.Context, volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindOneByNameAndAccountID(ctx, gomock.Any(), gomock.Any()).
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
					CreateAncestors(ctx, gomock.Any()).
					Return(sql.ErrConnDone).
					Times(1)
			},
		},
		{
			name:            "update descendant entries error",
			inputAccountID:  entry.AccountID,
			inputVolumeName: "volume",
			inputKey:        "test/sample.txt",
			inputNewKey:     "update/sample.txt",
			expectResult:    nil,
			expectError:     sql.ErrConnDone,
			setMockTransactionObj: func(ctx context.Context, transactionObj *mockTransaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockEntryRepo: func(ctx context.Context, entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindOneByKeyAndVolumeIDAndAccountID(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(entry, nil).
					Times(1)
			},
			setMockBodyRepo: func(context.Context, *mockRepository.MockBodyRepository) {},
			setMockVolumeRepo: func(ctx context.Context, volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindOneByNameAndAccountID(ctx, gomock.Any(), gomock.Any()).
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
					CreateAncestors(ctx, gomock.Any()).
					Return(nil).
					Times(1)
				entryServ.
					EXPECT().
					UpdateDescendants(ctx, gomock.Any(), gomock.Any()).
					Return(sql.ErrConnDone).
					Times(1)
			},
		},
		{
			name:            "update entry error",
			inputAccountID:  entry.AccountID,
			inputVolumeName: "volume",
			inputKey:        "test/sample.txt",
			inputNewKey:     "update/sample.txt",
			expectResult:    nil,
			expectError:     sql.ErrConnDone,
			setMockTransactionObj: func(ctx context.Context, transactionObj *mockTransaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockEntryRepo: func(ctx context.Context, entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindOneByKeyAndVolumeIDAndAccountID(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(entry, nil).
					Times(1)
				entryRepo.
					EXPECT().
					Update(ctx, gomock.Any()).
					Return(sql.ErrConnDone).
					Times(1)
			},
			setMockBodyRepo: func(context.Context, *mockRepository.MockBodyRepository) {},
			setMockVolumeRepo: func(ctx context.Context, volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindOneByNameAndAccountID(ctx, gomock.Any(), gomock.Any()).
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
					CreateAncestors(ctx, gomock.Any()).
					Return(nil).
					Times(1)
				entryServ.
					EXPECT().
					UpdateDescendants(ctx, gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name:            "update body error",
			inputAccountID:  entry.AccountID,
			inputVolumeName: "volume",
			inputKey:        "test/sample.txt",
			inputNewKey:     "update/sample.txt",
			expectResult:    nil,
			expectError:     afero.ErrFileClosed,
			setMockTransactionObj: func(ctx context.Context, transactionObj *mockTransaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockEntryRepo: func(ctx context.Context, entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindOneByKeyAndVolumeIDAndAccountID(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(entry, nil).
					Times(1)
				entryRepo.
					EXPECT().
					Update(ctx, gomock.Any()).
					Return(nil).
					Times(1)
			},
			setMockBodyRepo: func(_ context.Context, bodyRepo *mockRepository.MockBodyRepository) {
				bodyRepo.
					EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(afero.ErrFileClosed).
					Times(1)
			},
			setMockVolumeRepo: func(ctx context.Context, volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindOneByNameAndAccountID(ctx, gomock.Any(), gomock.Any()).
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
					CreateAncestors(ctx, gomock.Any()).
					Return(nil).
					Times(1)
				entryServ.
					EXPECT().
					UpdateDescendants(ctx, gomock.Any(), gomock.Any()).
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

			transactionObj := mockTransaction.NewMockTransactionObject(ctrl)
			tt.setMockTransactionObj(ctx, transactionObj)

			entryRepo := mockRepository.NewMockEntryRepository(ctrl)
			tt.setMockEntryRepo(ctx, entryRepo)

			bodyRepo := mockRepository.NewMockBodyRepository(ctrl)
			tt.setMockBodyRepo(ctx, bodyRepo)

			volumeRepo := mockRepository.NewMockVolumeRepository(ctrl)
			tt.setMockVolumeRepo(ctx, volumeRepo)

			entryServ := mockService.NewMockEntryService(ctrl)
			tt.setMockEntryServ(ctx, entryServ)

			uc := usecase.NewEntryUsecase(transactionObj, entryRepo, bodyRepo, volumeRepo, entryServ)
			result, err := uc.Update(ctx, tt.inputAccountID, tt.inputVolumeName, tt.inputKey, tt.inputNewKey)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			opts := cmp.Options{
				cmpopts.IgnoreFields(dto.EntryDTO{}, "UpdatedAt"),
			}
			if diff := cmp.Diff(result, tt.expectResult, opts...); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestEntry_Delete(t *testing.T) {
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
		AccountID: volume.AccountID,
		VolumeID:  volume.ID,
		Key:       "test/sample.txt",
		Size:      4,
		Type:      "text/plain; charset=utf-8",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name                  string
		inputAccountID        uuid.UUID
		inputVolumeName       string
		inputKey              string
		expectError           error
		setMockTransactionObj func(context.Context, *mockTransaction.MockTransactionObject)
		setMockEntryRepo      func(context.Context, *mockRepository.MockEntryRepository)
		setMockBodyRepo       func(context.Context, *mockRepository.MockBodyRepository)
		setMockVolumeRepo     func(context.Context, *mockRepository.MockVolumeRepository)
		setMockEntryServ      func(context.Context, *mockService.MockEntryService)
	}{
		{
			name:            "success",
			inputAccountID:  entry.AccountID,
			inputVolumeName: "volume",
			inputKey:        "test/sample.txt",
			expectError:     nil,
			setMockTransactionObj: func(ctx context.Context, transactionObj *mockTransaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockEntryRepo: func(ctx context.Context, entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindOneByKeyAndVolumeIDAndAccountID(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(entry, nil).
					Times(1)
				entryRepo.
					EXPECT().
					Delete(ctx, gomock.Any()).
					Return(nil).
					Times(1)
			},
			setMockBodyRepo: func(_ context.Context, bodyRepo *mockRepository.MockBodyRepository) {
				bodyRepo.
					EXPECT().
					Delete(gomock.Any()).
					Return(nil).
					Times(1)
			},
			setMockVolumeRepo: func(ctx context.Context, volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindOneByNameAndAccountID(ctx, gomock.Any(), gomock.Any()).
					Return(volume, nil).
					Times(1)
			},
			setMockEntryServ: func(ctx context.Context, entryServ *mockService.MockEntryService) {
				entryServ.
					EXPECT().
					DeleteDescendants(ctx, gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name:            "delete descendant entries error",
			inputAccountID:  entry.AccountID,
			inputVolumeName: "volume",
			inputKey:        "test/sample.txt",
			expectError:     sql.ErrConnDone,
			setMockTransactionObj: func(ctx context.Context, transactionObj *mockTransaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockEntryRepo: func(ctx context.Context, entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindOneByKeyAndVolumeIDAndAccountID(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(entry, nil).
					Times(1)
			},
			setMockBodyRepo: func(context.Context, *mockRepository.MockBodyRepository) {},
			setMockVolumeRepo: func(ctx context.Context, volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindOneByNameAndAccountID(ctx, gomock.Any(), gomock.Any()).
					Return(volume, nil).
					Times(1)
			},
			setMockEntryServ: func(ctx context.Context, entryServ *mockService.MockEntryService) {
				entryServ.
					EXPECT().
					DeleteDescendants(ctx, gomock.Any()).
					Return(sql.ErrConnDone).
					Times(1)
			},
		},
		{
			name:            "delete entry error",
			inputAccountID:  entry.AccountID,
			inputVolumeName: "volume",
			inputKey:        "test/sample.txt",
			expectError:     sql.ErrConnDone,
			setMockTransactionObj: func(ctx context.Context, transactionObj *mockTransaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockEntryRepo: func(ctx context.Context, entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindOneByKeyAndVolumeIDAndAccountID(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(entry, nil).
					Times(1)
				entryRepo.
					EXPECT().
					Delete(ctx, gomock.Any()).
					Return(sql.ErrConnDone).
					Times(1)
			},
			setMockBodyRepo: func(context.Context, *mockRepository.MockBodyRepository) {},
			setMockVolumeRepo: func(ctx context.Context, volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindOneByNameAndAccountID(ctx, gomock.Any(), gomock.Any()).
					Return(volume, nil).
					Times(1)
			},
			setMockEntryServ: func(ctx context.Context, entryServ *mockService.MockEntryService) {
				entryServ.
					EXPECT().
					DeleteDescendants(ctx, gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name:            "delete body error",
			inputAccountID:  entry.AccountID,
			inputVolumeName: "volume",
			inputKey:        "test/sample.txt",
			expectError:     afero.ErrFileClosed,
			setMockTransactionObj: func(ctx context.Context, transactionObj *mockTransaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockEntryRepo: func(ctx context.Context, entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindOneByKeyAndVolumeIDAndAccountID(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(entry, nil).
					Times(1)
				entryRepo.
					EXPECT().
					Delete(ctx, gomock.Any()).
					Return(nil).
					Times(1)
			},
			setMockBodyRepo: func(_ context.Context, bodyRepo *mockRepository.MockBodyRepository) {
				bodyRepo.
					EXPECT().
					Delete(gomock.Any()).
					Return(afero.ErrFileClosed).
					Times(1)
			},
			setMockVolumeRepo: func(ctx context.Context, volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindOneByNameAndAccountID(ctx, gomock.Any(), gomock.Any()).
					Return(volume, nil).
					Times(1)
			},
			setMockEntryServ: func(ctx context.Context, entryServ *mockService.MockEntryService) {
				entryServ.
					EXPECT().
					DeleteDescendants(ctx, gomock.Any()).
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

			transactionObj := mockTransaction.NewMockTransactionObject(ctrl)
			tt.setMockTransactionObj(ctx, transactionObj)

			entryRepo := mockRepository.NewMockEntryRepository(ctrl)
			tt.setMockEntryRepo(ctx, entryRepo)

			bodyRepo := mockRepository.NewMockBodyRepository(ctrl)
			tt.setMockBodyRepo(ctx, bodyRepo)

			volumeRepo := mockRepository.NewMockVolumeRepository(ctrl)
			tt.setMockVolumeRepo(ctx, volumeRepo)

			entryServ := mockService.NewMockEntryService(ctrl)
			tt.setMockEntryServ(ctx, entryServ)

			uc := usecase.NewEntryUsecase(transactionObj, entryRepo, bodyRepo, volumeRepo, entryServ)
			if err := uc.Delete(ctx, tt.inputAccountID, tt.inputVolumeName, tt.inputKey); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
		})
	}
}

func TestEntry_Head(t *testing.T) {
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
		AccountID: volume.AccountID,
		VolumeID:  volume.ID,
		Key:       "test/sample.txt",
		Size:      4,
		Type:      "text/plain; charset=utf-8",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	entryDTO := &dto.EntryDTO{
		ID:        entry.ID,
		AccountID: entry.AccountID,
		VolumeID:  entry.VolumeID,
		Key:       "test/sample.txt",
		Size:      entry.Size,
		Type:      entry.Type,
		CreatedAt: entry.CreatedAt,
		UpdatedAt: entry.UpdatedAt,
	}

	tests := []struct {
		name                  string
		inputAccountID        uuid.UUID
		inputVolumeName       string
		inputKey              string
		expectResult          *dto.EntryDTO
		expectError           error
		setMockTransactionObj func(context.Context, *mockTransaction.MockTransactionObject)
		setMockEntryRepo      func(context.Context, *mockRepository.MockEntryRepository)
		setMockVolumeRepo     func(context.Context, *mockRepository.MockVolumeRepository)
	}{
		{
			name:            "success",
			inputAccountID:  entry.AccountID,
			inputVolumeName: "volume",
			inputKey:        "test/sample.txt",
			expectResult:    entryDTO,
			expectError:     nil,
			setMockTransactionObj: func(ctx context.Context, transactionObj *mockTransaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockEntryRepo: func(ctx context.Context, entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindOneByKeyAndVolumeIDAndAccountID(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(entry, nil).
					Times(1)
			},
			setMockVolumeRepo: func(ctx context.Context, volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindOneByNameAndAccountID(ctx, gomock.Any(), gomock.Any()).
					Return(volume, nil).
					Times(1)
			},
		},
		{
			name:            "find entry error",
			inputAccountID:  entry.AccountID,
			inputVolumeName: "volume",
			inputKey:        "test/sample.txt",
			expectResult:    nil,
			expectError:     sql.ErrConnDone,
			setMockTransactionObj: func(ctx context.Context, transactionObj *mockTransaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockEntryRepo: func(ctx context.Context, entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindOneByKeyAndVolumeIDAndAccountID(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
			setMockVolumeRepo: func(ctx context.Context, volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindOneByNameAndAccountID(ctx, gomock.Any(), gomock.Any()).
					Return(volume, nil).
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

			entryRepo := mockRepository.NewMockEntryRepository(ctrl)
			tt.setMockEntryRepo(ctx, entryRepo)

			volumeRepo := mockRepository.NewMockVolumeRepository(ctrl)
			tt.setMockVolumeRepo(ctx, volumeRepo)

			uc := usecase.NewEntryUsecase(transactionObj, entryRepo, nil, volumeRepo, nil)
			result, err := uc.Head(ctx, tt.inputAccountID, tt.inputVolumeName, tt.inputKey)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			if diff := cmp.Diff(result, tt.expectResult); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestEntry_GetOne(t *testing.T) {
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
		AccountID: volume.AccountID,
		VolumeID:  volume.ID,
		Key:       "test/sample.txt",
		Size:      4,
		Type:      "text/plain; charset=utf-8",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	entryDTO := &dto.EntryDTO{
		ID:        entry.ID,
		AccountID: entry.AccountID,
		VolumeID:  entry.VolumeID,
		Key:       "test/sample.txt",
		Size:      entry.Size,
		Type:      entry.Type,
		CreatedAt: entry.CreatedAt,
		UpdatedAt: entry.UpdatedAt,
	}

	tests := []struct {
		name                  string
		inputAccountID        uuid.UUID
		inputVolumeName       string
		inputKey              string
		expectEntry           *dto.EntryDTO
		expectBody            io.ReadCloser
		expectError           error
		setMockTransactionObj func(context.Context, *mockTransaction.MockTransactionObject)
		setMockEntryRepo      func(context.Context, *mockRepository.MockEntryRepository)
		setMockBodyRepo       func(context.Context, *mockRepository.MockBodyRepository)
		setMockVolumeRepo     func(context.Context, *mockRepository.MockVolumeRepository)
	}{
		{
			name:            "success",
			inputAccountID:  entry.AccountID,
			inputVolumeName: "volume",
			inputKey:        "test/sample.txt",
			expectEntry:     entryDTO,
			expectBody:      nil,
			expectError:     nil,
			setMockTransactionObj: func(ctx context.Context, transactionObj *mockTransaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockEntryRepo: func(ctx context.Context, entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindOneByKeyAndVolumeIDAndAccountID(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(entry, nil).
					Times(1)
			},
			setMockBodyRepo: func(_ context.Context, bodyRepo *mockRepository.MockBodyRepository) {
				bodyRepo.
					EXPECT().
					FindOneByPath(gomock.Any()).
					Return(nil, nil).
					Times(1)
			},
			setMockVolumeRepo: func(ctx context.Context, volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindOneByNameAndAccountID(ctx, gomock.Any(), gomock.Any()).
					Return(volume, nil).
					Times(1)
			},
		},
		{
			name:            "find entry error",
			inputAccountID:  entry.AccountID,
			inputVolumeName: "volume",
			inputKey:        "test/sample.txt",
			expectEntry:     nil,
			expectBody:      nil,
			expectError:     sql.ErrConnDone,
			setMockTransactionObj: func(ctx context.Context, transactionObj *mockTransaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockEntryRepo: func(ctx context.Context, entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindOneByKeyAndVolumeIDAndAccountID(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
			setMockBodyRepo: func(context.Context, *mockRepository.MockBodyRepository) {},
			setMockVolumeRepo: func(ctx context.Context, volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindOneByNameAndAccountID(ctx, gomock.Any(), gomock.Any()).
					Return(volume, nil).
					Times(1)
			},
		},
		{
			name:            "find body error",
			inputAccountID:  entry.AccountID,
			inputVolumeName: "volume",
			inputKey:        "test/sample.txt",
			expectEntry:     nil,
			expectBody:      nil,
			expectError:     afero.ErrFileNotFound,
			setMockTransactionObj: func(ctx context.Context, transactionObj *mockTransaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockEntryRepo: func(ctx context.Context, entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindOneByKeyAndVolumeIDAndAccountID(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(entry, nil).
					Times(1)
			},
			setMockBodyRepo: func(_ context.Context, bodyRepo *mockRepository.MockBodyRepository) {
				bodyRepo.
					EXPECT().
					FindOneByPath(gomock.Any()).
					Return(nil, afero.ErrFileNotFound).
					Times(1)
			},
			setMockVolumeRepo: func(ctx context.Context, volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindOneByNameAndAccountID(ctx, gomock.Any(), gomock.Any()).
					Return(volume, nil).
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

			entryRepo := mockRepository.NewMockEntryRepository(ctrl)
			tt.setMockEntryRepo(ctx, entryRepo)

			bodyRepo := mockRepository.NewMockBodyRepository(ctrl)
			tt.setMockBodyRepo(ctx, bodyRepo)

			volumeRepo := mockRepository.NewMockVolumeRepository(ctrl)
			tt.setMockVolumeRepo(ctx, volumeRepo)

			uc := usecase.NewEntryUsecase(transactionObj, entryRepo, bodyRepo, volumeRepo, nil)
			entry, body, err := uc.GetOne(ctx, tt.inputAccountID, tt.inputVolumeName, tt.inputKey)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			if diff := cmp.Diff(entry, tt.expectEntry); diff != "" {
				t.Error(diff)
			}

			if diff := cmp.Diff(body, tt.expectBody); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestEntry_Search(t *testing.T) {
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
		AccountID: volume.AccountID,
		VolumeID:  volume.ID,
		Key:       "test/sample.txt",
		Size:      4,
		Type:      "text/plain; charset=utf-8",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	entryDTO := &dto.EntryDTO{
		ID:        entry.ID,
		AccountID: entry.AccountID,
		VolumeID:  entry.VolumeID,
		Key:       "test/sample.txt",
		Size:      entry.Size,
		Type:      entry.Type,
		CreatedAt: entry.CreatedAt,
		UpdatedAt: entry.UpdatedAt,
	}

	tests := []struct {
		name                  string
		inputAccountID        uuid.UUID
		inputVolumeName       string
		inputPrefix           *string
		inputDepth            *uint64
		expectResult          []*dto.EntryDTO
		expectError           error
		setMockTransactionObj func(context.Context, *mockTransaction.MockTransactionObject)
		setMockEntryRepo      func(context.Context, *mockRepository.MockEntryRepository)
		setMockVolumeRepo     func(context.Context, *mockRepository.MockVolumeRepository)
	}{
		{
			name:            "success",
			inputAccountID:  entry.AccountID,
			inputVolumeName: "volume",
			inputPrefix:     nil,
			inputDepth:      nil,
			expectResult:    []*dto.EntryDTO{entryDTO},
			expectError:     nil,
			setMockTransactionObj: func(ctx context.Context, transactionObj *mockTransaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockEntryRepo: func(ctx context.Context, entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindByVolumeIDAndAccountID(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]*entity.Entry{entry}, nil).
					Times(1)
			},
			setMockVolumeRepo: func(ctx context.Context, volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindOneByNameAndAccountID(ctx, gomock.Any(), gomock.Any()).
					Return(volume, nil).
					Times(1)
			},
		},
		{
			name:            "entry not found",
			inputAccountID:  entry.AccountID,
			inputVolumeName: "volume",
			inputPrefix:     nil,
			inputDepth:      nil,
			expectResult:    []*dto.EntryDTO{},
			expectError:     nil,
			setMockTransactionObj: func(ctx context.Context, transactionObj *mockTransaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockEntryRepo: func(ctx context.Context, entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindByVolumeIDAndAccountID(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, nil).
					Times(1)
			},
			setMockVolumeRepo: func(ctx context.Context, volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindOneByNameAndAccountID(ctx, gomock.Any(), gomock.Any()).
					Return(volume, nil).
					Times(1)
			},
		},
		{
			name:            "find entry error",
			inputAccountID:  entry.AccountID,
			inputVolumeName: "volume",
			inputPrefix:     nil,
			inputDepth:      nil,
			expectResult:    nil,
			expectError:     sql.ErrConnDone,
			setMockTransactionObj: func(ctx context.Context, transactionObj *mockTransaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockEntryRepo: func(ctx context.Context, entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindByVolumeIDAndAccountID(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
			setMockVolumeRepo: func(ctx context.Context, volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindOneByNameAndAccountID(ctx, gomock.Any(), gomock.Any()).
					Return(volume, nil).
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

			entryRepo := mockRepository.NewMockEntryRepository(ctrl)
			tt.setMockEntryRepo(ctx, entryRepo)

			volumeRepo := mockRepository.NewMockVolumeRepository(ctrl)
			tt.setMockVolumeRepo(ctx, volumeRepo)

			uc := usecase.NewEntryUsecase(transactionObj, entryRepo, nil, volumeRepo, nil)
			result, err := uc.Search(ctx, tt.inputAccountID, tt.inputVolumeName, tt.inputPrefix, tt.inputDepth)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			if diff := cmp.Diff(result, tt.expectResult); diff != "" {
				t.Error(diff)
			}
		})
	}
}
