package usecase_test

import (
	"database/sql"
	"net/http"
	"testing"
	"time"

	"github.com/atsumarukun/holos-api-pkg/errors"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/infrastructure/api/pkg/client"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/usecase"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/usecase/dto"
	"github.com/atsumarukun/holos-storage-api/test/assert"
	mockRepository "github.com/atsumarukun/holos-storage-api/test/mock/domain/repository"
)

func TestAuthorization_Authorize(t *testing.T) {
	ownerAccount := &entity.Account{
		ID: uuid.New(),
	}
	otherAccount := &entity.Account{
		ID: uuid.New(),
	}
	accountDTO := &dto.AccountDTO{
		ID: ownerAccount.ID,
	}
	publicVolume := &entity.Volume{
		ID:        uuid.New(),
		AccountID: ownerAccount.ID,
		Name:      "name",
		IsPublic:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	privateVolume := &entity.Volume{
		ID:        uuid.New(),
		AccountID: ownerAccount.ID,
		Name:      "name",
		IsPublic:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name               string
		inputCredential    string
		inputVolumeName    string
		inputKey           string
		inputMethod        string
		expectResult       *dto.AccountDTO
		expectError        error
		setMockAccountRepo func(*mockRepository.MockAccountRepository)
		setMockVolumeRepo  func(*mockRepository.MockVolumeRepository)
	}{
		{
			name:            "not get entry",
			inputCredential: "Session: YNDNun_KFu1uFmS691yJ6eqJ9eczRVKn",
			inputVolumeName: "",
			inputKey:        "",
			inputMethod:     "",
			expectResult:    accountDTO,
			expectError:     nil,
			setMockAccountRepo: func(accountRepo *mockRepository.MockAccountRepository) {
				accountRepo.EXPECT().
					FindOneByCredential(gomock.Any(), gomock.Any()).
					Return(ownerAccount, nil).
					Times(1)
			},
			setMockVolumeRepo: func(*mockRepository.MockVolumeRepository) {},
		},
		{
			name:               "get public volume entry",
			inputCredential:    "Session: YNDNun_KFu1uFmS691yJ6eqJ9eczRVKn",
			inputVolumeName:    "name",
			inputKey:           "key/sample.txt",
			inputMethod:        "GET",
			expectResult:       accountDTO,
			expectError:        nil,
			setMockAccountRepo: func(*mockRepository.MockAccountRepository) {},
			setMockVolumeRepo: func(volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindOneByName(gomock.Any(), gomock.Any()).
					Return(publicVolume, nil).
					Times(1)
			},
		},
		{
			name:            "get private volume entry",
			inputCredential: "Session: YNDNun_KFu1uFmS691yJ6eqJ9eczRVKn",
			inputVolumeName: "name",
			inputKey:        "key/sample.txt",
			inputMethod:     "GET",
			expectResult:    accountDTO,
			expectError:     nil,
			setMockAccountRepo: func(accountRepo *mockRepository.MockAccountRepository) {
				accountRepo.EXPECT().
					FindOneByCredential(gomock.Any(), gomock.Any()).
					Return(ownerAccount, nil).
					Times(1)
			},
			setMockVolumeRepo: func(volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindOneByName(gomock.Any(), gomock.Any()).
					Return(privateVolume, nil).
					Times(1)
			},
		},
		{
			name:            "unauthenticated when get entry",
			inputCredential: "",
			inputVolumeName: "name",
			inputKey:        "key/sample.txt",
			inputMethod:     "GET",
			expectResult:    nil,
			expectError:     client.ErrUnauthenticated,
			setMockAccountRepo: func(accountRepo *mockRepository.MockAccountRepository) {
				accountRepo.EXPECT().
					FindOneByCredential(gomock.Any(), gomock.Any()).
					Return(nil, errors.Wrap(client.ErrUnauthenticated, errors.CodeUnauthenticated, "failed to find account by credential")).
					Times(1)
			},
			setMockVolumeRepo: func(volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindOneByName(gomock.Any(), gomock.Any()).
					Return(privateVolume, nil).
					Times(1)
			},
		},
		{
			name:            "authorized account is not owner",
			inputCredential: "Session: YNDNun_KFu1uFmS691yJ6eqJ9eczRVKn",
			inputVolumeName: "name",
			inputKey:        "key/sample.txt",
			inputMethod:     "GET",
			expectResult:    nil,
			expectError:     usecase.ErrAccountIDMismatch,
			setMockAccountRepo: func(accountRepo *mockRepository.MockAccountRepository) {
				accountRepo.EXPECT().
					FindOneByCredential(gomock.Any(), gomock.Any()).
					Return(otherAccount, nil).
					Times(1)
			},
			setMockVolumeRepo: func(volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindOneByName(gomock.Any(), gomock.Any()).
					Return(privateVolume, nil).
					Times(1)
			},
		},
		{
			name:            "authorize error",
			inputCredential: "Session: YNDNun_KFu1uFmS691yJ6eqJ9eczRVKn",
			inputVolumeName: "",
			inputKey:        "",
			inputMethod:     "",
			expectResult:    nil,
			expectError:     http.ErrServerClosed,
			setMockAccountRepo: func(accountRepo *mockRepository.MockAccountRepository) {
				accountRepo.EXPECT().
					FindOneByCredential(gomock.Any(), gomock.Any()).
					Return(nil, errors.Wrap(http.ErrServerClosed, errors.CodeInternalServerError, "failed to find account by credential")).
					Times(1)
			},
			setMockVolumeRepo: func(*mockRepository.MockVolumeRepository) {},
		},
		{
			name:               "find volume error",
			inputCredential:    "Session: YNDNun_KFu1uFmS691yJ6eqJ9eczRVKn",
			inputVolumeName:    "name",
			inputKey:           "key/sample.txt",
			inputMethod:        "GET",
			expectResult:       nil,
			expectError:        sql.ErrConnDone,
			setMockAccountRepo: func(accountRepo *mockRepository.MockAccountRepository) {},
			setMockVolumeRepo: func(volumeRepo *mockRepository.MockVolumeRepository) {
				volumeRepo.
					EXPECT().
					FindOneByName(gomock.Any(), gomock.Any()).
					Return(nil, errors.Wrap(sql.ErrConnDone, errors.CodeInternalServerError, "failed to find volume by name")).
					Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := t.Context()

			accountRepo := mockRepository.NewMockAccountRepository(ctrl)
			tt.setMockAccountRepo(accountRepo)

			volumeRepo := mockRepository.NewMockVolumeRepository(ctrl)
			tt.setMockVolumeRepo(volumeRepo)

			uc := usecase.NewAuthorizationUsecase(accountRepo, volumeRepo)
			result, err := uc.Authorize(ctx, tt.inputCredential, tt.inputVolumeName, tt.inputKey, tt.inputMethod)
			assert.Error(t, err, tt.expectError)

			if diff := cmp.Diff(tt.expectResult, result); diff != "" {
				t.Error(diff)
			}
		})
	}
}
