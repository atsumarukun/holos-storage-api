package handler_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestEntry_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)

	accountID := uuid.New()
	entryDTO := &dto.EntryDTO{
		ID:        uuid.New(),
		AccountID: accountID,
		VolumeID:  uuid.New(),
		Key:       "test/sample.txt",
		Size:      4,
		Type:      "text/plain; charset=utf-8",
		IsPublic:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	if err := w.WriteField("key", "test/sample.txt"); err != nil {
		t.Error(err)
	}
	if err := w.WriteField("is_public", "false"); err != nil {
		t.Error(err)
	}
	fw, err := w.CreateFormFile("file", "sample.txt")
	if err != nil {
		t.Error(err)
	}
	_, err = io.Copy(fw, strings.NewReader("test"))
	if err != nil {
		t.Error(err)
	}
	w.Close()

	tests := []struct {
		name           string
		requestBody    *bytes.Buffer
		isSetAccountID bool
		expectCode     int
		expectResponse map[string]any
		setMockEntryUC func(context.Context, *mockUsecase.MockEntryUsecase)
	}{
		{
			name:           "success",
			requestBody:    body,
			isSetAccountID: true,
			expectCode:     http.StatusCreated,
			expectResponse: map[string]any{"id": entryDTO.ID.String(), "volume_id": entryDTO.VolumeID.String(), "key": entryDTO.Key, "size": entryDTO.Size, "type": entryDTO.Type, "is_public": entryDTO.IsPublic, "created_at": entryDTO.CreatedAt.Format(time.RFC3339Nano), "updated_at": entryDTO.UpdatedAt.Format(time.RFC3339Nano)},
			setMockEntryUC: func(ctx context.Context, entryUC *mockUsecase.MockEntryUsecase) {
				entryUC.
					EXPECT().
					Create(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(entryDTO, nil).
					Times(1)
			},
		},
		{
			name:           "invalid request",
			requestBody:    nil,
			isSetAccountID: true,
			expectCode:     http.StatusBadRequest,
			expectResponse: map[string]any{"message": "faild to parse multipart/form-data"},
			setMockEntryUC: func(context.Context, *mockUsecase.MockEntryUsecase) {},
		},
		{
			name:           "account id not found",
			requestBody:    body,
			isSetAccountID: false,
			expectCode:     http.StatusInternalServerError,
			expectResponse: map[string]any{"message": "internal server error"},
			setMockEntryUC: func(context.Context, *mockUsecase.MockEntryUsecase) {},
		},
		{
			name:           "create error",
			requestBody:    body,
			isSetAccountID: true,
			expectCode:     http.StatusInternalServerError,
			expectResponse: map[string]any{"message": "internal server error"},
			setMockEntryUC: func(ctx context.Context, entryUC *mockUsecase.MockEntryUsecase) {
				entryUC.
					EXPECT().
					Create(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
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
			c.Request, err = http.NewRequestWithContext(ctx, "POST", "/entries", tt.requestBody)
			if err != nil {
				t.Error(err)
			}
			if tt.isSetAccountID {
				c.Set("accountID", accountID)
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			entryUC := mockUsecase.NewMockEntryUsecase(ctrl)
			tt.setMockEntryUC(ctx, entryUC)

			hdl := handler.NewEntryHandler(entryUC)
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
