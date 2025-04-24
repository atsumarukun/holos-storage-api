package usecase_test

import (
	"context"
	"database/sql"
	"errors"
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
		setMockTransactionObj func(context.Context, *mockTransaction.MockTransactionObject)
		setMockVolumeRepo     func(context.Context, *mockRepository.MockVolumeRepository)
		setMockVolumeServ     func(context.Context, *mockService.MockVolumeService)
	}{
		{
			name:           "success",
			inputAccountID: accountID,
			inputName:      "name",
			inputIsPublic:  false,
			expectResult:   volumeDTO,
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
			setMockTransactionObj: func(context.Context, *mockTransaction.MockTransactionObject) {},
			setMockVolumeRepo:     func(context.Context, *mockRepository.MockVolumeRepository) {},
			setMockVolumeServ:     func(context.Context, *mockService.MockVolumeService) {},
		},
		{
			name:           "volume already exists",
			inputAccountID: accountID,
			inputName:      "name",
			inputIsPublic:  false,
			expectResult:   nil,
			expectError:    service.ErrVolumeAlreadyExists,
			setMockTransactionObj: func(ctx context.Context, transactionObj *mockTransaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockVolumeRepo: func(context.Context, *mockRepository.MockVolumeRepository) {},
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

			transactionObj := mockTransaction.NewMockTransactionObject(ctrl)
			tt.setMockTransactionObj(ctx, transactionObj)

			volumeRepo := mockRepository.NewMockVolumeRepository(ctrl)
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

func TestVolume_Update(t *testing.T) {
	id := uuid.New()
	accountID := uuid.New()
	volume := &entity.Volume{
		ID:        id,
		AccountID: accountID,
		Name:      "name",
		IsPublic:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	volumeDTO := &dto.VolumeDTO{
		ID:        id,
		AccountID: accountID,
		Name:      "name",
		IsPublic:  false,
		CreatedAt: volume.CreatedAt,
		UpdatedAt: volume.UpdatedAt,
	}

	tests := []struct {
		name                  string
		inputAccountID        uuid.UUID
		inputID               uuid.UUID
		inputName             string
		inputIsPublic         bool
		expectResult          *dto.VolumeDTO
		expectError           error
		setMockTransactionObj func(context.Context, *mockTransaction.MockTransactionObject)
		setMockVolumeRepo     func(context.Context, *mockRepository.MockVolumeRepository)
		setMockVolumeServ     func(context.Context, *mockService.MockVolumeService)
	}{
		{
			name:           "success",
			inputAccountID: accountID,
			inputID:        id,
			inputName:      "name",
			inputIsPublic:  false,
			expectResult:   volumeDTO,
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
				volumeRepo.
					EXPECT().
					Update(ctx, gomock.Any()).
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
			name:           "invalid name",
			inputAccountID: accountID,
			inputID:        id,
			inputName:      "",
			inputIsPublic:  false,
			expectResult:   nil,
			expectError:    entity.ErrShortVolumeName,
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
			setMockVolumeServ: func(context.Context, *mockService.MockVolumeService) {},
		},
		{
			name:           "volume already exists",
			inputAccountID: accountID,
			inputID:        id,
			inputName:      "name",
			inputIsPublic:  false,
			expectResult:   nil,
			expectError:    service.ErrVolumeAlreadyExists,
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
			setMockVolumeServ: func(ctx context.Context, volumeServ *mockService.MockVolumeService) {
				volumeServ.
					EXPECT().
					Exists(ctx, gomock.Any()).
					Return(service.ErrVolumeAlreadyExists).
					Times(1)
			},
		},
		{
			name:           "find error",
			inputAccountID: accountID,
			inputID:        id,
			inputName:      "name",
			inputIsPublic:  false,
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
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
			setMockVolumeServ: func(context.Context, *mockService.MockVolumeService) {},
		},
		{
			name:           "update error",
			inputAccountID: accountID,
			inputID:        id,
			inputName:      "name",
			inputIsPublic:  false,
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
				volumeRepo.
					EXPECT().
					Update(ctx, gomock.Any()).
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

			transactionObj := mockTransaction.NewMockTransactionObject(ctrl)
			tt.setMockTransactionObj(ctx, transactionObj)

			volumeRepo := mockRepository.NewMockVolumeRepository(ctrl)
			tt.setMockVolumeRepo(ctx, volumeRepo)

			volumeServ := mockService.NewMockVolumeService(ctrl)
			tt.setMockVolumeServ(ctx, volumeServ)

			uc := usecase.NewVolumeUsecase(transactionObj, volumeRepo, volumeServ)
			result, err := uc.Update(ctx, tt.inputAccountID, tt.inputID, tt.inputName, tt.inputIsPublic)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			opts := cmp.Options{
				cmpopts.IgnoreFields(dto.VolumeDTO{}, "ID", "UpdatedAt"),
			}
			if diff := cmp.Diff(result, tt.expectResult, opts...); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestVolume_Delete(t *testing.T) {
	id := uuid.New()
	accountID := uuid.New()
	volume := &entity.Volume{
		ID:        id,
		AccountID: accountID,
		Name:      "name",
		IsPublic:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name                  string
		inputAccountID        uuid.UUID
		inputID               uuid.UUID
		expectError           error
		setMockTransactionObj func(context.Context, *mockTransaction.MockTransactionObject)
		setMockVolumeRepo     func(context.Context, *mockRepository.MockVolumeRepository)
	}{
		{
			name:           "success",
			inputAccountID: accountID,
			inputID:        id,
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
				volumeRepo.
					EXPECT().
					Delete(ctx, gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name:           "find error",
			inputAccountID: accountID,
			inputID:        id,
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
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
		},
		{
			name:           "delete error",
			inputAccountID: accountID,
			inputID:        id,
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
				volumeRepo.
					EXPECT().
					Delete(ctx, gomock.Any()).
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

			uc := usecase.NewVolumeUsecase(transactionObj, volumeRepo, nil)
			if err := uc.Delete(ctx, tt.inputAccountID, tt.inputID); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
		})
	}
}

func TestVolume_GetOne(t *testing.T) {
	id := uuid.New()
	accountID := uuid.New()
	volume := &entity.Volume{
		ID:        id,
		AccountID: accountID,
		Name:      "name",
		IsPublic:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	volumeDTO := &dto.VolumeDTO{
		ID:        id,
		AccountID: accountID,
		Name:      "name",
		IsPublic:  false,
		CreatedAt: volume.CreatedAt,
		UpdatedAt: volume.UpdatedAt,
	}

	tests := []struct {
		name              string
		inputAccountID    uuid.UUID
		inputID           uuid.UUID
		expectResult      *dto.VolumeDTO
		expectError       error
		setMockVolumeRepo func(context.Context, *mockRepository.MockVolumeRepository)
	}{
		{
			name:           "success",
			inputAccountID: accountID,
			inputID:        id,
			expectResult:   volumeDTO,
			expectError:    nil,
			setMockVolumeRepo: func(ctx context.Context, volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindOneByIDAndAccountID(ctx, gomock.Any(), gomock.Any()).
					Return(volume, nil).
					Times(1)
			},
		},
		{
			name:           "not found",
			inputAccountID: accountID,
			inputID:        id,
			expectResult:   nil,
			expectError:    usecase.ErrVolumeNotFound,
			setMockVolumeRepo: func(ctx context.Context, volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindOneByIDAndAccountID(ctx, gomock.Any(), gomock.Any()).
					Return(nil, nil).
					Times(1)
			},
		},
		{
			name:           "find error",
			inputAccountID: accountID,
			inputID:        id,
			expectResult:   nil,
			expectError:    sql.ErrConnDone,
			setMockVolumeRepo: func(ctx context.Context, volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindOneByIDAndAccountID(ctx, gomock.Any(), gomock.Any()).
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

			uc := usecase.NewVolumeUsecase(nil, volumeRepo, nil)
			result, err := uc.GetOne(ctx, tt.inputAccountID, tt.inputID)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			opts := cmp.Options{
				cmpopts.IgnoreFields(dto.VolumeDTO{}, "ID", "UpdatedAt"),
			}
			if diff := cmp.Diff(result, tt.expectResult, opts...); diff != "" {
				t.Error(diff)
			}
		})
	}
}
