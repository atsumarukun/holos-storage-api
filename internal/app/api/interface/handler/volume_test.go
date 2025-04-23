package handler_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/interface/handler"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/usecase/dto"
	mockUsecase "github.com/atsumarukun/holos-storage-api/test/mock/usecase"
)

func TestVolume_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)

	volumeDTO := &dto.VolumeDTO{
		ID:        uuid.New(),
		AccountID: uuid.New(),
		Name:      "name",
		IsPublic:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name            string
		requestJSON     []byte
		isSetAccountID  bool
		expectCode      int
		expectResponse  map[string]any
		setMockVolumeUC func(context.Context, *mockUsecase.MockVolumeUsecase)
	}{
		{
			name:           "success",
			requestJSON:    []byte(`{"name": "name", "is_public": false}`),
			isSetAccountID: true,
			expectCode:     http.StatusCreated,
			expectResponse: map[string]any{"id": volumeDTO.ID.String(), "account_id": volumeDTO.AccountID.String(), "name": volumeDTO.Name, "is_public": volumeDTO.IsPublic, "created_at": volumeDTO.CreatedAt.Format(time.RFC3339Nano), "updated_at": volumeDTO.UpdatedAt.Format(time.RFC3339Nano)},
			setMockVolumeUC: func(ctx context.Context, volumeUC *mockUsecase.MockVolumeUsecase) {
				volumeUC.
					EXPECT().
					Create(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(volumeDTO, nil).
					Times(1)
			},
		},
		{
			name:            "invalid request",
			requestJSON:     nil,
			isSetAccountID:  true,
			expectCode:      http.StatusBadRequest,
			expectResponse:  map[string]any{"message": "bad request"},
			setMockVolumeUC: func(context.Context, *mockUsecase.MockVolumeUsecase) {},
		},
		{
			name:            "account id not found",
			requestJSON:     []byte(`{"name": "name", "is_public": false}`),
			isSetAccountID:  false,
			expectCode:      http.StatusInternalServerError,
			expectResponse:  map[string]any{"message": "internal server error"},
			setMockVolumeUC: func(context.Context, *mockUsecase.MockVolumeUsecase) {},
		},
		{
			name:           "create error",
			requestJSON:    []byte(`{"name": "name", "is_public": false}`),
			isSetAccountID: true,
			expectCode:     http.StatusInternalServerError,
			expectResponse: map[string]any{"message": "internal server error"},
			setMockVolumeUC: func(ctx context.Context, volumeUC *mockUsecase.MockVolumeUsecase) {
				volumeUC.
					EXPECT().
					Create(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, sql.ErrConnDone).
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
			c.Request, err = http.NewRequestWithContext(ctx, "POST", "/volumes", bytes.NewBuffer(tt.requestJSON))
			if err != nil {
				t.Error(err)
			}
			if tt.isSetAccountID {
				c.Set("accountID", uuid.New())
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			volumeUC := mockUsecase.NewMockVolumeUsecase(ctrl)
			tt.setMockVolumeUC(ctx, volumeUC)

			hdl := handler.NewVolumeHandler(volumeUC)
			hdl.Create(c)

			c.Writer.WriteHeaderNow()

			if w.Code != tt.expectCode {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectCode, w.Code)
			}

			var response map[string]any
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Error(err)
			}
			if diff := cmp.Diff(response, tt.expectResponse); diff != "" {
				t.Error(diff)
			}
		})
	}
}
