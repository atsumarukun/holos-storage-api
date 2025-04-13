package usecase_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/usecase"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/usecase/dto"
	"github.com/atsumarukun/holos-storage-api/test/mock/domain/repository"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
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
		expectResult       *dto.AccountDTO
		expectError        error
		setMockAccountRepo func(context.Context, *repository.MockAccountRepository)
	}{
		{
			name:            "success",
			inputCredential: "Session: YNDNun_KFu1uFmS691yJ6eqJ9eczRVKn",
			expectResult:    accountDTO,
			expectError:     nil,
			setMockAccountRepo: func(ctx context.Context, accountRepo *repository.MockAccountRepository) {
				accountRepo.EXPECT().
					FindOneByCredential(ctx, gomock.Any()).
					Return(account, nil).
					Times(1)
			},
		},
		{
			name:            "authorize error",
			inputCredential: "Session: YNDNun_KFu1uFmS691yJ6eqJ9eczRVKn",
			expectResult:    nil,
			expectError:     http.ErrServerClosed,
			setMockAccountRepo: func(ctx context.Context, accountRepo *repository.MockAccountRepository) {
				accountRepo.EXPECT().
					FindOneByCredential(ctx, gomock.Any()).
					Return(nil, http.ErrServerClosed).
					Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := t.Context()

			accountRepo := repository.NewMockAccountRepository(ctrl)
			tt.setMockAccountRepo(ctx, accountRepo)

			uc := usecase.NewAuthorizationUsecase(accountRepo)
			result, err := uc.Authorize(ctx, tt.inputCredential)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			if diff := cmp.Diff(result, tt.expectResult); diff != "" {
				t.Error(diff)
			}
		})
	}
}
