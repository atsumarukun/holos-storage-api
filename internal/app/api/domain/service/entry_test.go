package service_test

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/repository"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/service"
	mockRepository "github.com/atsumarukun/holos-storage-api/test/mock/domain/repository"
)

func TestEntry_Exists(t *testing.T) {
	entry := &entity.Entry{
		ID:        uuid.New(),
		AccountID: uuid.New(),
		VolumeID:  uuid.New(),
		Key:       "key/sample.txt",
		Size:      4,
		Type:      "text/plain; charset=utf-8",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name             string
		inputEntry       *entity.Entry
		expectError      error
		setMockEntryRepo func(*mockRepository.MockEntryRepository)
	}{
		{
			name:        "not exists",
			inputEntry:  entry,
			expectError: nil,
			setMockEntryRepo: func(entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindOneByKeyAndVolumeID(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, repository.ErrEntryNotFound).
					Times(1)
			},
		},
		{
			name:        "exists",
			inputEntry:  entry,
			expectError: service.ErrEntryAlreadyExists,
			setMockEntryRepo: func(entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindOneByKeyAndVolumeID(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(entry, nil).
					Times(1)
			},
		},
		{
			name:             "entry is nil",
			inputEntry:       nil,
			expectError:      service.ErrRequiredEntry,
			setMockEntryRepo: func(*mockRepository.MockEntryRepository) {},
		},
		{
			name:        "find error",
			inputEntry:  entry,
			expectError: sql.ErrConnDone,
			setMockEntryRepo: func(entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindOneByKeyAndVolumeID(gomock.Any(), gomock.Any(), gomock.Any()).
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

			serv := service.NewEntryService(entryRepo)
			if err := serv.Exists(ctx, tt.inputEntry); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
		})
	}
}

func TestEntry_CreateAncestors(t *testing.T) {
	accountID := uuid.New()
	volumeID := uuid.New()
	entry := &entity.Entry{
		ID:        uuid.New(),
		AccountID: accountID,
		VolumeID:  volumeID,
		Key:       "key/sample.txt",
		Size:      4,
		Type:      "text/plain; charset=utf-8",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	ancestorEntry := &entity.Entry{
		ID:        uuid.New(),
		AccountID: accountID,
		VolumeID:  volumeID,
		Key:       "key",
		Size:      0,
		Type:      "folder",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name             string
		inputEntry       *entity.Entry
		expectError      error
		setMockEntryRepo func(*mockRepository.MockEntryRepository)
	}{
		{
			name:        "successfully created",
			inputEntry:  entry,
			expectError: nil,
			setMockEntryRepo: func(entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindOneByKeyAndVolumeID(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, repository.ErrEntryNotFound).
					Times(1)
				entryRepo.
					EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name:             "entry is nil",
			inputEntry:       nil,
			expectError:      service.ErrRequiredEntry,
			setMockEntryRepo: func(*mockRepository.MockEntryRepository) {},
		},
		{
			name:        "ancestor entry already exists",
			inputEntry:  entry,
			expectError: nil,
			setMockEntryRepo: func(entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindOneByKeyAndVolumeID(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(ancestorEntry, nil).
					Times(1)
			},
		},
		{
			name:        "find entry error",
			inputEntry:  entry,
			expectError: sql.ErrConnDone,
			setMockEntryRepo: func(entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindOneByKeyAndVolumeID(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
		},
		{
			name:        "create entry error",
			inputEntry:  entry,
			expectError: sql.ErrConnDone,
			setMockEntryRepo: func(entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindOneByKeyAndVolumeID(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, repository.ErrEntryNotFound).
					Times(1)
				entryRepo.
					EXPECT().
					Create(gomock.Any(), gomock.Any()).
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

			entryRepo := mockRepository.NewMockEntryRepository(ctrl)
			tt.setMockEntryRepo(entryRepo)

			serv := service.NewEntryService(entryRepo)
			if err := serv.CreateAncestors(ctx, tt.inputEntry); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
		})
	}
}

func TestEntry_UpdateDescendants(t *testing.T) {
	accountID := uuid.New()
	volumeID := uuid.New()
	fileEntry := &entity.Entry{
		ID:        uuid.New(),
		AccountID: accountID,
		VolumeID:  volumeID,
		Key:       "key/sample.txt",
		Size:      4,
		Type:      "text/plain; charset=utf-8",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	folderEntry := &entity.Entry{
		ID:        uuid.New(),
		AccountID: accountID,
		VolumeID:  volumeID,
		Key:       "update",
		Size:      0,
		Type:      "folder",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	descendantEntry := &entity.Entry{
		ID:        uuid.New(),
		AccountID: accountID,
		VolumeID:  volumeID,
		Key:       "key/sample.txt",
		Size:      4,
		Type:      "text/plain; charset=utf-8",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name             string
		inputEntry       *entity.Entry
		inputSrc         string
		expectError      error
		setMockEntryRepo func(*mockRepository.MockEntryRepository)
	}{
		{
			name:             "update file entry",
			inputEntry:       fileEntry,
			inputSrc:         "key",
			expectError:      nil,
			setMockEntryRepo: func(*mockRepository.MockEntryRepository) {},
		},
		{
			name:        "update folder entry",
			inputEntry:  folderEntry,
			inputSrc:    "key",
			expectError: nil,
			setMockEntryRepo: func(entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindByVolumeIDAndAccountID(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]*entity.Entry{descendantEntry}, nil).
					Times(1)
				entryRepo.
					EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name:             "entry is nil",
			inputEntry:       nil,
			inputSrc:         "update",
			expectError:      service.ErrRequiredEntry,
			setMockEntryRepo: func(*mockRepository.MockEntryRepository) {},
		},
		{
			name:        "find entry error",
			inputEntry:  folderEntry,
			inputSrc:    "key",
			expectError: sql.ErrConnDone,
			setMockEntryRepo: func(entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindByVolumeIDAndAccountID(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
		},
		{
			name:        "update entry error",
			inputEntry:  folderEntry,
			inputSrc:    "key",
			expectError: sql.ErrConnDone,
			setMockEntryRepo: func(entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindByVolumeIDAndAccountID(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]*entity.Entry{descendantEntry}, nil).
					Times(1)
				entryRepo.
					EXPECT().
					Update(gomock.Any(), gomock.Any()).
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

			entryRepo := mockRepository.NewMockEntryRepository(ctrl)
			tt.setMockEntryRepo(entryRepo)

			serv := service.NewEntryService(entryRepo)
			if err := serv.UpdateDescendants(ctx, tt.inputEntry, tt.inputSrc); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
		})
	}
}

func TestEntry_DeleteDescendants(t *testing.T) {
	accountID := uuid.New()
	volumeID := uuid.New()
	fileEntry := &entity.Entry{
		ID:        uuid.New(),
		AccountID: accountID,
		VolumeID:  volumeID,
		Key:       "key/sample.txt",
		Size:      4,
		Type:      "text/plain; charset=utf-8",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	folderEntry := &entity.Entry{
		ID:        uuid.New(),
		AccountID: accountID,
		VolumeID:  volumeID,
		Key:       "key",
		Size:      0,
		Type:      "folder",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	descendantEntry := &entity.Entry{
		ID:        uuid.New(),
		AccountID: accountID,
		VolumeID:  volumeID,
		Key:       "key/sample.txt",
		Size:      4,
		Type:      "text/plain; charset=utf-8",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name             string
		inputEntry       *entity.Entry
		expectError      error
		setMockEntryRepo func(*mockRepository.MockEntryRepository)
	}{
		{
			name:             "delete file entry",
			inputEntry:       fileEntry,
			expectError:      nil,
			setMockEntryRepo: func(*mockRepository.MockEntryRepository) {},
		},
		{
			name:        "delete folder entry",
			inputEntry:  folderEntry,
			expectError: nil,
			setMockEntryRepo: func(entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindByVolumeIDAndAccountID(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]*entity.Entry{descendantEntry}, nil).
					Times(1)
				entryRepo.
					EXPECT().
					Delete(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name:             "entry is nil",
			inputEntry:       nil,
			expectError:      service.ErrRequiredEntry,
			setMockEntryRepo: func(*mockRepository.MockEntryRepository) {},
		},
		{
			name:        "find entry error",
			inputEntry:  folderEntry,
			expectError: sql.ErrConnDone,
			setMockEntryRepo: func(entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindByVolumeIDAndAccountID(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
		},
		{
			name:        "delete entry error",
			inputEntry:  folderEntry,
			expectError: sql.ErrConnDone,
			setMockEntryRepo: func(entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindByVolumeIDAndAccountID(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]*entity.Entry{descendantEntry}, nil).
					Times(1)
				entryRepo.
					EXPECT().
					Delete(gomock.Any(), gomock.Any()).
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

			entryRepo := mockRepository.NewMockEntryRepository(ctrl)
			tt.setMockEntryRepo(entryRepo)

			serv := service.NewEntryService(entryRepo)
			if err := serv.DeleteDescendants(ctx, tt.inputEntry); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
		})
	}
}

func TestEntry_Copy(t *testing.T) {
	accountID := uuid.New()
	volumeID := uuid.New()
	fileEntry := &entity.Entry{
		ID:        uuid.New(),
		AccountID: accountID,
		VolumeID:  volumeID,
		Key:       "key/sample.txt",
		Size:      4,
		Type:      "text/plain; charset=utf-8",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	folderEntry := &entity.Entry{
		ID:        uuid.New(),
		AccountID: accountID,
		VolumeID:  volumeID,
		Key:       "key",
		Size:      0,
		Type:      "folder",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	copiedFileEntry := &entity.Entry{
		ID:        uuid.New(),
		AccountID: fileEntry.AccountID,
		VolumeID:  fileEntry.VolumeID,
		Key:       "key/sample copy.txt",
		Size:      fileEntry.Size,
		Type:      fileEntry.Type,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	copiedFolderEntry := &entity.Entry{
		ID:        uuid.New(),
		AccountID: folderEntry.AccountID,
		VolumeID:  folderEntry.VolumeID,
		Key:       "key copy",
		Size:      folderEntry.Size,
		Type:      folderEntry.Type,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	copiedFileEntryTwice := &entity.Entry{
		ID:        uuid.New(),
		AccountID: fileEntry.AccountID,
		VolumeID:  fileEntry.VolumeID,
		Key:       "key/sample copy copy.txt",
		Size:      fileEntry.Size,
		Type:      fileEntry.Type,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name             string
		inputEntry       *entity.Entry
		expectResult     *entity.Entry
		expectError      error
		setMockEntryRepo func(*mockRepository.MockEntryRepository)
	}{
		{
			name:         "copy file entry",
			inputEntry:   fileEntry,
			expectResult: copiedFileEntry,
			expectError:  nil,
			setMockEntryRepo: func(entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindOneByKeyAndVolumeID(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, repository.ErrEntryNotFound).
					Times(1)
			},
		},
		{
			name:         "copy folder entry",
			inputEntry:   folderEntry,
			expectResult: copiedFolderEntry,
			expectError:  nil,
			setMockEntryRepo: func(entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindOneByKeyAndVolumeID(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, repository.ErrEntryNotFound).
					Times(1)
			},
		},
		{
			name:         "duplicate key",
			inputEntry:   fileEntry,
			expectResult: copiedFileEntryTwice,
			expectError:  nil,
			setMockEntryRepo: func(entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindOneByKeyAndVolumeID(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(copiedFileEntry, nil).
					Times(1)
				entryRepo.
					EXPECT().
					FindOneByKeyAndVolumeID(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, repository.ErrEntryNotFound).
					Times(1)
			},
		},
		{
			name:             "entry is nil",
			inputEntry:       nil,
			expectResult:     nil,
			expectError:      service.ErrRequiredEntry,
			setMockEntryRepo: func(*mockRepository.MockEntryRepository) {},
		},
		{
			name:         "find entry error",
			inputEntry:   fileEntry,
			expectResult: nil,
			expectError:  sql.ErrConnDone,
			setMockEntryRepo: func(entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindOneByKeyAndVolumeID(gomock.Any(), gomock.Any(), gomock.Any()).
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

			serv := service.NewEntryService(entryRepo)

			result, err := serv.Copy(ctx, tt.inputEntry)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			opts := cmp.Options{
				cmpopts.IgnoreFields(entity.Entry{}, "ID", "CreatedAt", "UpdatedAt"),
			}
			if diff := cmp.Diff(tt.expectResult, result, opts...); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func CopyDescendants(t *testing.T) {
	accountID := uuid.New()
	volumeID := uuid.New()
	fileEntry := &entity.Entry{
		ID:        uuid.New(),
		AccountID: accountID,
		VolumeID:  volumeID,
		Key:       "key/sample.txt",
		Size:      4,
		Type:      "text/plain; charset=utf-8",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	copiedFileEntry := &entity.Entry{
		ID:        uuid.New(),
		AccountID: fileEntry.AccountID,
		VolumeID:  fileEntry.VolumeID,
		Key:       "key/sample copy.txt",
		Size:      fileEntry.Size,
		Type:      fileEntry.Type,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	copiedFolderEntry := &entity.Entry{
		ID:        uuid.New(),
		AccountID: accountID,
		VolumeID:  volumeID,
		Key:       "key copy",
		Size:      0,
		Type:      "folder",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name             string
		inputEntry       *entity.Entry
		inputSrc         string
		expectError      error
		setMockEntryRepo func(*mockRepository.MockEntryRepository)
	}{
		{
			name:             "copy file entry",
			inputEntry:       copiedFileEntry,
			inputSrc:         "key/sample.txt",
			expectError:      nil,
			setMockEntryRepo: func(*mockRepository.MockEntryRepository) {},
		},
		{
			name:        "copy folder entry",
			inputEntry:  copiedFolderEntry,
			inputSrc:    "key",
			expectError: nil,
			setMockEntryRepo: func(entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindByVolumeIDAndAccountID(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]*entity.Entry{fileEntry}, nil).
					Times(1)
				entryRepo.
					EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name:             "entry is nil",
			inputEntry:       nil,
			inputSrc:         "key",
			expectError:      service.ErrRequiredEntry,
			setMockEntryRepo: func(*mockRepository.MockEntryRepository) {},
		},
		{
			name:        "find entry error",
			inputEntry:  copiedFolderEntry,
			inputSrc:    "key",
			expectError: sql.ErrConnDone,
			setMockEntryRepo: func(entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindByVolumeIDAndAccountID(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
		},
		{
			name:        "create entry error",
			inputEntry:  copiedFolderEntry,
			inputSrc:    "key",
			expectError: sql.ErrConnDone,
			setMockEntryRepo: func(entryRepo *mockRepository.MockEntryRepository) {
				entryRepo.
					EXPECT().
					FindByVolumeIDAndAccountID(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]*entity.Entry{fileEntry}, nil).
					Times(1)
				entryRepo.
					EXPECT().
					Create(gomock.Any(), gomock.Any()).
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

			entryRepo := mockRepository.NewMockEntryRepository(ctrl)
			tt.setMockEntryRepo(entryRepo)

			serv := service.NewEntryService(entryRepo)
			if err := serv.CopyDescendants(ctx, tt.inputEntry, tt.inputSrc); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
		})
	}
}
