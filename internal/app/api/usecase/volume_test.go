package usecase_test

import (
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
		setMockTransactionObj func(*mockTransaction.MockTransactionObject)
		setMockVolumeRepo     func(*mockRepository.MockVolumeRepository)
		setMockBodyRepo       func(*mockRepository.MockBodyRepository)
		setMockVolumeServ     func(*mockService.MockVolumeService)
	}{
		{
			name:           "successfully created",
			inputAccountID: accountID,
			inputName:      "name",
			inputIsPublic:  false,
			expectResult:   volumeDTO,
			expectError:    nil,
			setMockTransactionObj: func(transactionObj *mockTransaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockVolumeRepo: func(volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
			setMockBodyRepo: func(bodyRepo *mockRepository.MockBodyRepository) {
				bodyRepo.
					EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
			setMockVolumeServ: func(volumeServ *mockService.MockVolumeService) {
				volumeServ.
					EXPECT().
					Exists(gomock.Any(), gomock.Any()).
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
			setMockTransactionObj: func(*mockTransaction.MockTransactionObject) {},
			setMockVolumeRepo:     func(*mockRepository.MockVolumeRepository) {},
			setMockBodyRepo:       func(*mockRepository.MockBodyRepository) {},
			setMockVolumeServ:     func(*mockService.MockVolumeService) {},
		},
		{
			name:           "volume already exists",
			inputAccountID: accountID,
			inputName:      "name",
			inputIsPublic:  false,
			expectResult:   nil,
			expectError:    service.ErrVolumeAlreadyExists,
			setMockTransactionObj: func(transactionObj *mockTransaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockVolumeRepo: func(*mockRepository.MockVolumeRepository) {},
			setMockBodyRepo:   func(*mockRepository.MockBodyRepository) {},
			setMockVolumeServ: func(volumeServ *mockService.MockVolumeService) {
				volumeServ.
					EXPECT().
					Exists(gomock.Any(), gomock.Any()).
					Return(service.ErrVolumeAlreadyExists).
					Times(1)
			},
		},
		{
			name:           "create volume error",
			inputAccountID: accountID,
			inputName:      "name",
			inputIsPublic:  false,
			expectResult:   nil,
			expectError:    sql.ErrConnDone,
			setMockTransactionObj: func(transactionObj *mockTransaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockVolumeRepo: func(volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(sql.ErrConnDone).
					Times(1)
			},
			setMockBodyRepo: func(*mockRepository.MockBodyRepository) {},
			setMockVolumeServ: func(volumeServ *mockService.MockVolumeService) {
				volumeServ.
					EXPECT().
					Exists(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name:           "create body error",
			inputAccountID: accountID,
			inputName:      "name",
			inputIsPublic:  false,
			expectResult:   nil,
			expectError:    io.ErrNoProgress,
			setMockTransactionObj: func(transactionObj *mockTransaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockVolumeRepo: func(volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
			setMockBodyRepo: func(bodyRepo *mockRepository.MockBodyRepository) {
				bodyRepo.
					EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(io.ErrNoProgress).
					Times(1)
			},
			setMockVolumeServ: func(volumeServ *mockService.MockVolumeService) {
				volumeServ.
					EXPECT().
					Exists(gomock.Any(), gomock.Any()).
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
			tt.setMockTransactionObj(transactionObj)

			volumeRepo := mockRepository.NewMockVolumeRepository(ctrl)
			tt.setMockVolumeRepo(volumeRepo)

			bodyRepo := mockRepository.NewMockBodyRepository(ctrl)
			tt.setMockBodyRepo(bodyRepo)

			volumeServ := mockService.NewMockVolumeService(ctrl)
			tt.setMockVolumeServ(volumeServ)

			uc := usecase.NewVolumeUsecase(transactionObj, volumeRepo, bodyRepo, volumeServ)
			result, err := uc.Create(ctx, tt.inputAccountID, tt.inputName, tt.inputIsPublic)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			opts := cmp.Options{
				cmpopts.IgnoreFields(dto.VolumeDTO{}, "ID", "CreatedAt", "UpdatedAt"),
			}
			if diff := cmp.Diff(tt.expectResult, result, opts...); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestVolume_Update(t *testing.T) {
	accountID := uuid.New()
	volume := &entity.Volume{
		ID:        uuid.New(),
		AccountID: accountID,
		Name:      "name",
		IsPublic:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	volumeDTO := &dto.VolumeDTO{
		ID:        volume.ID,
		AccountID: volume.AccountID,
		Name:      "update",
		IsPublic:  volume.IsPublic,
		CreatedAt: volume.CreatedAt,
		UpdatedAt: volume.UpdatedAt,
	}

	tests := []struct {
		name                  string
		inputAccountID        uuid.UUID
		inputName             string
		inputNewName          string
		inputIsPublic         bool
		expectResult          *dto.VolumeDTO
		expectError           error
		setMockTransactionObj func(*mockTransaction.MockTransactionObject)
		setMockVolumeRepo     func(*mockRepository.MockVolumeRepository)
		setMockBodyRepo       func(*mockRepository.MockBodyRepository)
		setMockVolumeServ     func(*mockService.MockVolumeService)
	}{
		{
			name:           "successfully updated",
			inputAccountID: accountID,
			inputName:      "name",
			inputNewName:   "update",
			inputIsPublic:  false,
			expectResult:   volumeDTO,
			expectError:    nil,
			setMockTransactionObj: func(transactionObj *mockTransaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockVolumeRepo: func(volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindOneByNameAndAccountID(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(volume, nil).
					Times(1)
				volumeRepo.
					EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
			setMockBodyRepo: func(bodyRepo *mockRepository.MockBodyRepository) {
				bodyRepo.
					EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
			setMockVolumeServ: func(volumeServ *mockService.MockVolumeService) {
				volumeServ.
					EXPECT().
					Exists(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name:           "invalid name",
			inputAccountID: accountID,
			inputName:      "name",
			inputNewName:   "",
			inputIsPublic:  false,
			expectResult:   nil,
			expectError:    entity.ErrShortVolumeName,
			setMockTransactionObj: func(transactionObj *mockTransaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockVolumeRepo: func(volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindOneByNameAndAccountID(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(volume, nil).
					Times(1)
			},
			setMockBodyRepo:   func(*mockRepository.MockBodyRepository) {},
			setMockVolumeServ: func(*mockService.MockVolumeService) {},
		},
		{
			name:           "volume already exists",
			inputAccountID: accountID,
			inputName:      "name",
			inputNewName:   "update",
			inputIsPublic:  false,
			expectResult:   nil,
			expectError:    service.ErrVolumeAlreadyExists,
			setMockTransactionObj: func(transactionObj *mockTransaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockVolumeRepo: func(volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindOneByNameAndAccountID(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(volume, nil).
					Times(1)
			},
			setMockBodyRepo: func(*mockRepository.MockBodyRepository) {},
			setMockVolumeServ: func(volumeServ *mockService.MockVolumeService) {
				volumeServ.
					EXPECT().
					Exists(gomock.Any(), gomock.Any()).
					Return(service.ErrVolumeAlreadyExists).
					Times(1)
			},
		},
		{
			name:           "find volume error",
			inputAccountID: accountID,
			inputName:      "name",
			inputNewName:   "update",
			inputIsPublic:  false,
			expectResult:   nil,
			expectError:    sql.ErrConnDone,
			setMockTransactionObj: func(transactionObj *mockTransaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockVolumeRepo: func(volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindOneByNameAndAccountID(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
			setMockBodyRepo:   func(*mockRepository.MockBodyRepository) {},
			setMockVolumeServ: func(*mockService.MockVolumeService) {},
		},
		{
			name:           "update volume error",
			inputAccountID: accountID,
			inputName:      "name",
			inputNewName:   "update",
			inputIsPublic:  false,
			expectResult:   nil,
			expectError:    sql.ErrConnDone,
			setMockTransactionObj: func(transactionObj *mockTransaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockVolumeRepo: func(volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindOneByNameAndAccountID(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(volume, nil).
					Times(1)
				volumeRepo.
					EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(sql.ErrConnDone).
					Times(1)
			},
			setMockBodyRepo: func(*mockRepository.MockBodyRepository) {},
			setMockVolumeServ: func(volumeServ *mockService.MockVolumeService) {
				volumeServ.
					EXPECT().
					Exists(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name:           "update body error",
			inputAccountID: accountID,
			inputName:      "name",
			inputNewName:   "update",
			inputIsPublic:  false,
			expectResult:   nil,
			expectError:    afero.ErrFileClosed,
			setMockTransactionObj: func(transactionObj *mockTransaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockVolumeRepo: func(volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindOneByNameAndAccountID(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(volume, nil).
					Times(1)
				volumeRepo.
					EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
			setMockBodyRepo: func(bodyRepo *mockRepository.MockBodyRepository) {
				bodyRepo.
					EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(afero.ErrFileClosed).
					Times(1)
			},
			setMockVolumeServ: func(volumeServ *mockService.MockVolumeService) {
				volumeServ.
					EXPECT().
					Exists(gomock.Any(), gomock.Any()).
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
			tt.setMockTransactionObj(transactionObj)

			volumeRepo := mockRepository.NewMockVolumeRepository(ctrl)
			tt.setMockVolumeRepo(volumeRepo)

			bodyRepo := mockRepository.NewMockBodyRepository(ctrl)
			tt.setMockBodyRepo(bodyRepo)

			volumeServ := mockService.NewMockVolumeService(ctrl)
			tt.setMockVolumeServ(volumeServ)

			uc := usecase.NewVolumeUsecase(transactionObj, volumeRepo, bodyRepo, volumeServ)
			result, err := uc.Update(ctx, tt.inputAccountID, tt.inputName, tt.inputNewName, tt.inputIsPublic)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			opts := cmp.Options{
				cmpopts.IgnoreFields(dto.VolumeDTO{}, "UpdatedAt"),
			}
			if diff := cmp.Diff(tt.expectResult, result, opts...); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestVolume_Delete(t *testing.T) {
	accountID := uuid.New()
	volume := &entity.Volume{
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
		expectError           error
		setMockTransactionObj func(*mockTransaction.MockTransactionObject)
		setMockVolumeRepo     func(*mockRepository.MockVolumeRepository)
		setMockBodyRepo       func(*mockRepository.MockBodyRepository)
		setMockVolumeServ     func(*mockService.MockVolumeService)
	}{
		{
			name:           "successfully deleted",
			inputAccountID: accountID,
			inputName:      "name",
			expectError:    nil,
			setMockTransactionObj: func(transactionObj *mockTransaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockVolumeRepo: func(volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindOneByNameAndAccountID(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(volume, nil).
					Times(1)
				volumeRepo.
					EXPECT().
					Delete(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
			setMockBodyRepo: func(bodyRepo *mockRepository.MockBodyRepository) {
				bodyRepo.
					EXPECT().
					Delete(gomock.Any()).
					Return(nil).
					Times(1)
			},
			setMockVolumeServ: func(volumeServ *mockService.MockVolumeService) {
				volumeServ.
					EXPECT().
					CanDelete(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name:           "find volume error",
			inputAccountID: accountID,
			inputName:      "name",
			expectError:    sql.ErrConnDone,
			setMockTransactionObj: func(transactionObj *mockTransaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockVolumeRepo: func(volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindOneByNameAndAccountID(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
			setMockBodyRepo:   func(*mockRepository.MockBodyRepository) {},
			setMockVolumeServ: func(*mockService.MockVolumeService) {},
		},
		{
			name:           "volume has entries",
			inputAccountID: accountID,
			inputName:      "name",
			expectError:    service.ErrVolumeHasEntries,
			setMockTransactionObj: func(transactionObj *mockTransaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockVolumeRepo: func(volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindOneByNameAndAccountID(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(volume, nil).
					Times(1)
			},
			setMockBodyRepo: func(*mockRepository.MockBodyRepository) {},
			setMockVolumeServ: func(volumeServ *mockService.MockVolumeService) {
				volumeServ.
					EXPECT().
					CanDelete(gomock.Any(), gomock.Any()).
					Return(service.ErrVolumeHasEntries).
					Times(1)
			},
		},
		{
			name:           "delete volume error",
			inputAccountID: accountID,
			inputName:      "name",
			expectError:    sql.ErrConnDone,
			setMockTransactionObj: func(transactionObj *mockTransaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockVolumeRepo: func(volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindOneByNameAndAccountID(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(volume, nil).
					Times(1)
				volumeRepo.
					EXPECT().
					Delete(gomock.Any(), gomock.Any()).
					Return(sql.ErrConnDone).
					Times(1)
			},
			setMockBodyRepo: func(*mockRepository.MockBodyRepository) {},
			setMockVolumeServ: func(volumeServ *mockService.MockVolumeService) {
				volumeServ.
					EXPECT().
					CanDelete(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name:           "delete body error",
			inputAccountID: accountID,
			inputName:      "name",
			expectError:    afero.ErrFileClosed,
			setMockTransactionObj: func(transactionObj *mockTransaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockVolumeRepo: func(volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindOneByNameAndAccountID(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(volume, nil).
					Times(1)
				volumeRepo.
					EXPECT().
					Delete(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
			setMockBodyRepo: func(bodyRepo *mockRepository.MockBodyRepository) {
				bodyRepo.
					EXPECT().
					Delete(gomock.Any()).
					Return(afero.ErrFileClosed).
					Times(1)
			},
			setMockVolumeServ: func(volumeServ *mockService.MockVolumeService) {
				volumeServ.
					EXPECT().
					CanDelete(gomock.Any(), gomock.Any()).
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
			tt.setMockTransactionObj(transactionObj)

			volumeRepo := mockRepository.NewMockVolumeRepository(ctrl)
			tt.setMockVolumeRepo(volumeRepo)

			bodyRepo := mockRepository.NewMockBodyRepository(ctrl)
			tt.setMockBodyRepo(bodyRepo)

			volumeServ := mockService.NewMockVolumeService(ctrl)
			tt.setMockVolumeServ(volumeServ)

			uc := usecase.NewVolumeUsecase(transactionObj, volumeRepo, bodyRepo, volumeServ)
			if err := uc.Delete(ctx, tt.inputAccountID, tt.inputName); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
		})
	}
}

func TestVolume_GetOne(t *testing.T) {
	accountID := uuid.New()
	volume := &entity.Volume{
		ID:        uuid.New(),
		AccountID: accountID,
		Name:      "name",
		IsPublic:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	volumeDTO := &dto.VolumeDTO{
		ID:        volume.ID,
		AccountID: volume.AccountID,
		Name:      volume.Name,
		IsPublic:  volume.IsPublic,
		CreatedAt: volume.CreatedAt,
		UpdatedAt: volume.UpdatedAt,
	}

	tests := []struct {
		name              string
		inputAccountID    uuid.UUID
		inputName         string
		expectResult      *dto.VolumeDTO
		expectError       error
		setMockVolumeRepo func(*mockRepository.MockVolumeRepository)
	}{
		{
			name:           "successfully got one",
			inputAccountID: accountID,
			inputName:      "name",
			expectResult:   volumeDTO,
			expectError:    nil,
			setMockVolumeRepo: func(volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindOneByNameAndAccountID(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(volume, nil).
					Times(1)
			},
		},
		{
			name:           "find error",
			inputAccountID: accountID,
			inputName:      "name",
			expectResult:   nil,
			expectError:    sql.ErrConnDone,
			setMockVolumeRepo: func(volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindOneByNameAndAccountID(gomock.Any(), gomock.Any(), gomock.Any()).
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

			uc := usecase.NewVolumeUsecase(nil, volumeRepo, nil, nil)
			result, err := uc.GetOne(ctx, tt.inputAccountID, tt.inputName)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			if diff := cmp.Diff(tt.expectResult, result); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestVolume_GetAll(t *testing.T) {
	accountID := uuid.New()
	volume := &entity.Volume{
		ID:        uuid.New(),
		AccountID: accountID,
		Name:      "name",
		IsPublic:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	volumeDTO := &dto.VolumeDTO{
		ID:        volume.ID,
		AccountID: volume.AccountID,
		Name:      volume.Name,
		IsPublic:  volume.IsPublic,
		CreatedAt: volume.CreatedAt,
		UpdatedAt: volume.UpdatedAt,
	}

	tests := []struct {
		name              string
		inputAccountID    uuid.UUID
		expectResult      []*dto.VolumeDTO
		expectError       error
		setMockVolumeRepo func(*mockRepository.MockVolumeRepository)
	}{
		{
			name:           "successfully got all",
			inputAccountID: accountID,
			expectResult:   []*dto.VolumeDTO{volumeDTO},
			expectError:    nil,
			setMockVolumeRepo: func(volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindByAccountID(gomock.Any(), gomock.Any()).
					Return([]*entity.Volume{volume}, nil).
					Times(1)
			},
		},
		{
			name:           "not found",
			inputAccountID: accountID,
			expectResult:   []*dto.VolumeDTO{},
			expectError:    nil,
			setMockVolumeRepo: func(volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindByAccountID(gomock.Any(), gomock.Any()).
					Return(nil, nil).
					Times(1)
			},
		},
		{
			name:           "find error",
			inputAccountID: accountID,
			expectResult:   nil,
			expectError:    sql.ErrConnDone,
			setMockVolumeRepo: func(volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindByAccountID(gomock.Any(), gomock.Any()).
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

			uc := usecase.NewVolumeUsecase(nil, volumeRepo, nil, nil)
			result, err := uc.GetAll(ctx, tt.inputAccountID)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			if diff := cmp.Diff(tt.expectResult, result); diff != "" {
				t.Error(diff)
			}
		})
	}
}
