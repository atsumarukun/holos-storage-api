package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/repository"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/interface/middleware"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/usecase/dto"
	mockUsecase "github.com/atsumarukun/holos-storage-api/test/mock/usecase"
)

func TestAuthorization_Authorize(t *testing.T) {
	gin.SetMode(gin.TestMode)

	accountDTO := &dto.AccountDTO{
		ID: uuid.New(),
	}

	tests := []struct {
		name                   string
		authorizationHeader    string
		expectResult           uuid.UUID
		expectError            []byte
		setMockAuthorizationUC func(*mockUsecase.MockAuthorizationUsecase)
	}{
		{
			name:                "session token is set",
			authorizationHeader: "Session 1Ty1HKTPKTt8xEi-_3HTbWf2SCHOdqOS",
			expectResult:        accountDTO.ID,
			expectError:         nil,
			setMockAuthorizationUC: func(authorizationUC *mockUsecase.MockAuthorizationUsecase) {
				authorizationUC.EXPECT().
					Authorize(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(accountDTO, nil).
					Times(1)
			},
		},
		{
			name:                "session token not set",
			authorizationHeader: "",
			expectResult:        uuid.Nil,
			expectError:         []byte(`{"message":"unauthorized"}`),
			setMockAuthorizationUC: func(authorizationUC *mockUsecase.MockAuthorizationUsecase) {
				authorizationUC.EXPECT().
					Authorize(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, repository.ErrUnauthorized).
					Times(1)
			},
		},
		{
			name:                "authorize error",
			authorizationHeader: "Session 1Ty1HKTPKTt8xEi-_3HTbWf2SCHOdqOS",
			expectResult:        uuid.Nil,
			expectError:         []byte(`{"message":"internal server error"}`),
			setMockAuthorizationUC: func(authorizationUC *mockUsecase.MockAuthorizationUsecase) {
				authorizationUC.EXPECT().
					Authorize(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
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
			c.Request, err = http.NewRequestWithContext(ctx, "GET", "/volumes", http.NoBody)
			if err != nil {
				t.Error(err)
			}
			c.Request.Header.Add("Authorization", tt.authorizationHeader)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			authorizationUC := mockUsecase.NewMockAuthorizationUsecase(ctrl)
			tt.setMockAuthorizationUC(authorizationUC)

			mw := middleware.NewAuthorizationMiddleware(authorizationUC)
			mw.Authorize(c)

			accountID, _ := c.Get("accountID")
			result, _ := accountID.(uuid.UUID)
			if diff := cmp.Diff(result, tt.expectResult); diff != "" {
				t.Error(diff)
			}

			if diff := cmp.Diff(tt.expectError, w.Body.Bytes()); diff != "" {
				t.Error(diff)
			}
		})
	}
}
