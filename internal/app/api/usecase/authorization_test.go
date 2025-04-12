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
	authorization := &entity.Authorization{
		AccountID: uuid.New(),
	}
	authorizationDTO := &dto.AuthorizationDTO{
		AccountID: authorization.AccountID,
	}

	tests := []struct {
		name                     string
		inputCredential          string
		expectResult             *dto.AuthorizationDTO
		expectError              error
		setMockAuthorizationRepo func(context.Context, *repository.MockAuthorizationRepository)
	}{
		{
			name:            "success",
			inputCredential: "Session: YNDNun_KFu1uFmS691yJ6eqJ9eczRVKn",
			expectResult:    authorizationDTO,
			expectError:     nil,
			setMockAuthorizationRepo: func(ctx context.Context, authorizationRepo *repository.MockAuthorizationRepository) {
				authorizationRepo.EXPECT().
					Authorize(ctx, gomock.Any()).
					Return(authorization, nil).
					Times(1)
			},
		},
		{
			name:            "authorize error",
			inputCredential: "Session: YNDNun_KFu1uFmS691yJ6eqJ9eczRVKn",
			expectResult:    nil,
			expectError:     http.ErrServerClosed,
			setMockAuthorizationRepo: func(ctx context.Context, authorizationRepo *repository.MockAuthorizationRepository) {
				authorizationRepo.EXPECT().
					Authorize(ctx, gomock.Any()).
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

			authorizationRepo := repository.NewMockAuthorizationRepository(ctrl)
			tt.setMockAuthorizationRepo(ctx, authorizationRepo)

			uc := usecase.NewAuthorizationUsecase(authorizationRepo)
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
