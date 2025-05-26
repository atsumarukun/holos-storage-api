package usecase_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/usecase"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/usecase/dto"
	mockRepository "github.com/atsumarukun/holos-storage-api/test/mock/domain/repository"
)

func TestAuthorization_Authorize(t *testing.T) {
	account := &entity.Account{
		ID: uuid.New(),
	}
	accountDTO := &dto.AccountDTO{
		ID: account.ID,
	}

	tests := []struct {
		name               string
		inputCredential    string
		inputVolumeName    string
		inputKey           string
		inputMethod        string
		expectResult       *dto.AccountDTO
		expectError        error
		setMockAccountRepo func(context.Context, *mockRepository.MockAccountRepository)
		setMockVolumeRepo  func(context.Context, *mockRepository.MockVolumeRepository)
	}{
		{
			name:            "success",
			inputCredential: "Session: YNDNun_KFu1uFmS691yJ6eqJ9eczRVKn",
			inputVolumeName: "",
			inputKey:        "",
			inputMethod:     "",
			expectResult:    accountDTO,
			expectError:     nil,
			setMockAccountRepo: func(ctx context.Context, accountRepo *mockRepository.MockAccountRepository) {
				accountRepo.EXPECT().
					FindOneByCredential(ctx, gomock.Any()).
					Return(account, nil).
					Times(1)
			},
			setMockVolumeRepo: func(context.Context, *mockRepository.MockVolumeRepository) {},
		},
		{
			name:            "authorize error",
			inputCredential: "Session: YNDNun_KFu1uFmS691yJ6eqJ9eczRVKn",
			inputVolumeName: "",
			inputKey:        "",
			inputMethod:     "",
			expectResult:    nil,
			expectError:     http.ErrServerClosed,
			setMockAccountRepo: func(ctx context.Context, accountRepo *mockRepository.MockAccountRepository) {
				accountRepo.EXPECT().
					FindOneByCredential(ctx, gomock.Any()).
					Return(nil, http.ErrServerClosed).
					Times(1)
			},
			setMockVolumeRepo: func(context.Context, *mockRepository.MockVolumeRepository) {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := t.Context()

			accountRepo := mockRepository.NewMockAccountRepository(ctrl)
			tt.setMockAccountRepo(ctx, accountRepo)

			volumeRepo := mockRepository.NewMockVolumeRepository(ctrl)
			tt.setMockVolumeRepo(ctx, volumeRepo)

			uc := usecase.NewAuthorizationUsecase(accountRepo, volumeRepo)
			result, err := uc.Authorize(ctx, tt.inputCredential, tt.inputVolumeName, tt.inputKey, tt.inputMethod)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			if diff := cmp.Diff(result, tt.expectResult); diff != "" {
				t.Error(diff)
			}
		})
	}
}
