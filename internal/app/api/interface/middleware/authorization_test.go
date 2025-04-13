package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/interface/middleware"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/usecase/dto"
	"github.com/atsumarukun/holos-storage-api/test/mock/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func TestAuthorization_Authorize(t *testing.T) {
	gin.SetMode(gin.TestMode)

	authorizationDTO := &dto.AuthorizationDTO{
		AccountID: uuid.New(),
	}

	tests := []struct {
		name                   string
		authorizationHeader    string
		expectResult           uuid.UUID
		setMockAuthorizationUC func(context.Context, *usecase.MockAuthorizationUsecase)
	}{
		{
			name:                "success",
			authorizationHeader: "Session 1Ty1HKTPKTt8xEi-_3HTbWf2SCHOdqOS",
			expectResult:        authorizationDTO.AccountID,
			setMockAuthorizationUC: func(ctx context.Context, authorizationUC *usecase.MockAuthorizationUsecase) {
				authorizationUC.EXPECT().
					Authorize(ctx, gomock.Any()).
					Return(authorizationDTO, nil).
					Times(1)
			},
		},
		{
			name:                   "invalid authorization header",
			authorizationHeader:    "",
			expectResult:           uuid.Nil,
			setMockAuthorizationUC: func(context.Context, *usecase.MockAuthorizationUsecase) {},
		},
		{
			name:                "authorize error",
			authorizationHeader: "Session 1Ty1HKTPKTt8xEi-_3HTbWf2SCHOdqOS",
			expectResult:        uuid.Nil,
			setMockAuthorizationUC: func(ctx context.Context, authorizationUC *usecase.MockAuthorizationUsecase) {
				authorizationUC.EXPECT().
					Authorize(ctx, gomock.Any()).
					Return(nil, http.ErrServerClosed).
					Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			w := httptest.NewRecorder()

			c, _ := gin.CreateTestContext(w)
			var err error
			c.Request, err = http.NewRequestWithContext(ctx, "GET", "/folders", http.NoBody)
			if err != nil {
				t.Error(err)
			}
			c.Request.Header.Add("Authorization", tt.authorizationHeader)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			authorizationUC := usecase.NewMockAuthorizationUsecase(ctrl)
			tt.setMockAuthorizationUC(ctx, authorizationUC)

			mw := middleware.NewAuthorizationMiddleware(authorizationUC)
			mw.Authorize(c)

			accountID, exists := c.Get("accountID")
			if exists && tt.expectResult == uuid.Nil {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectResult, accountID)
			} else {
				result, _ := accountID.(uuid.UUID)
				if diff := cmp.Diff(result, tt.expectResult); diff != "" {
					t.Error(diff)
				}
			}
		})
	}
}
