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
	"strconv"
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

func buildMultipartBody(t *testing.T) (body io.Reader, contentType string) {
	buffer := &bytes.Buffer{}
	writer := multipart.NewWriter(buffer)
	defer writer.Close()
	if err := writer.WriteField("volume_name", "volume"); err != nil {
		t.Error(err)
	}
	if err := writer.WriteField("key", "test/sample.txt"); err != nil {
		t.Error(err)
	}
	if err := writer.WriteField("is_public", "false"); err != nil {
		t.Error(err)
	}
	fw, err := writer.CreateFormFile("file", "sample.txt")
	if err != nil {
		t.Error(err)
	}
	_, err = io.Copy(fw, strings.NewReader("test"))
	if err != nil {
		t.Error(err)
	}
	return buffer, writer.FormDataContentType()
}

func TestEntry_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)

	accountID := uuid.New()
	volumeID := uuid.New()
	entryDTO := &dto.EntryDTO{
		ID:        uuid.New(),
		AccountID: accountID,
		VolumeID:  volumeID,
		Key:       "test/sample.txt",
		Size:      4,
		Type:      "text/plain; charset=utf-8",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name           string
		buildBody      func(*testing.T) (io.Reader, string)
		isSetAccountID bool
		expectCode     int
		expectResponse map[string]any
		setMockEntryUC func(context.Context, *mockUsecase.MockEntryUsecase)
	}{
		{
			name:           "success",
			buildBody:      buildMultipartBody,
			isSetAccountID: true,
			expectCode:     http.StatusCreated,
			expectResponse: map[string]any{"id": entryDTO.ID.String(), "volume_id": entryDTO.VolumeID.String(), "key": entryDTO.Key, "size": entryDTO.Size, "type": entryDTO.Type, "created_at": entryDTO.CreatedAt.Format(time.RFC3339Nano), "updated_at": entryDTO.UpdatedAt.Format(time.RFC3339Nano)},
			setMockEntryUC: func(ctx context.Context, entryUC *mockUsecase.MockEntryUsecase) {
				entryUC.
					EXPECT().
					Create(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(entryDTO, nil).
					Times(1)
			},
		},
		{
			name:           "invalid request",
			buildBody:      func(*testing.T) (io.Reader, string) { return http.NoBody, "" },
			isSetAccountID: true,
			expectCode:     http.StatusBadRequest,
			expectResponse: map[string]any{"message": "failed to parse multipart/form-data"},
			setMockEntryUC: func(context.Context, *mockUsecase.MockEntryUsecase) {},
		},
		{
			name:           "account id not found",
			buildBody:      buildMultipartBody,
			isSetAccountID: false,
			expectCode:     http.StatusInternalServerError,
			expectResponse: map[string]any{"message": "internal server error"},
			setMockEntryUC: func(context.Context, *mockUsecase.MockEntryUsecase) {},
		},
		{
			name:           "create error",
			buildBody:      buildMultipartBody,
			isSetAccountID: true,
			expectCode:     http.StatusInternalServerError,
			expectResponse: map[string]any{"message": "internal server error"},
			setMockEntryUC: func(ctx context.Context, entryUC *mockUsecase.MockEntryUsecase) {
				entryUC.
					EXPECT().
					Create(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			w := httptest.NewRecorder()

			body, contentType := tt.buildBody(t)

			c, _ := gin.CreateTestContext(w)
			var err error
			c.Request, err = http.NewRequestWithContext(ctx, "POST", "/entries", body)
			if err != nil {
				t.Error(err)
			}
			c.Request.Header.Add("Content-Type", contentType)
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

			if tt.expectCode == http.StatusCreated {
				if size, ok := response["size"].(float64); ok {
					response["size"] = uint64(size)
				}
			}

			if diff := cmp.Diff(response, tt.expectResponse); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestEntry_Update(t *testing.T) {
	gin.SetMode(gin.TestMode)

	id := uuid.New()
	accountID := uuid.New()
	volumeID := uuid.New()
	entryDTO := &dto.EntryDTO{
		ID:        id,
		AccountID: accountID,
		VolumeID:  volumeID,
		Key:       "test/sample.txt",
		Size:      4,
		Type:      "text/plain; charset=utf-8",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name           string
		requestJSON    []byte
		isSetAccountID bool
		expectCode     int
		expectResponse map[string]any
		setMockEntryUC func(context.Context, *mockUsecase.MockEntryUsecase)
	}{
		{
			name:           "success",
			requestJSON:    []byte(`{"key": "update/sample.txt"}`),
			isSetAccountID: true,
			expectCode:     http.StatusOK,
			expectResponse: map[string]any{"id": entryDTO.ID.String(), "volume_id": entryDTO.VolumeID.String(), "key": entryDTO.Key, "size": entryDTO.Size, "type": entryDTO.Type, "created_at": entryDTO.CreatedAt.Format(time.RFC3339Nano), "updated_at": entryDTO.UpdatedAt.Format(time.RFC3339Nano)},
			setMockEntryUC: func(ctx context.Context, entryUC *mockUsecase.MockEntryUsecase) {
				entryUC.
					EXPECT().
					Update(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(entryDTO, nil).
					Times(1)
			},
		},
		{
			name:           "invalid request",
			requestJSON:    nil,
			isSetAccountID: true,
			expectCode:     http.StatusBadRequest,
			expectResponse: map[string]any{"message": "failed to parse json"},
			setMockEntryUC: func(context.Context, *mockUsecase.MockEntryUsecase) {},
		},
		{
			name:           "account id not found",
			requestJSON:    []byte(`{"key": "update/sample.txt"}`),
			isSetAccountID: false,
			expectCode:     http.StatusInternalServerError,
			expectResponse: map[string]any{"message": "internal server error"},
			setMockEntryUC: func(context.Context, *mockUsecase.MockEntryUsecase) {},
		},
		{
			name:           "update error",
			requestJSON:    []byte(`{"key": "update/sample.txt"}`),
			isSetAccountID: true,
			expectCode:     http.StatusInternalServerError,
			expectResponse: map[string]any{"message": "internal server error"},
			setMockEntryUC: func(ctx context.Context, entryUC *mockUsecase.MockEntryUsecase) {
				entryUC.
					EXPECT().
					Update(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
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
			c.Request, err = http.NewRequestWithContext(ctx, "PUT", "/entries/volume/test/sample.txt", bytes.NewBuffer(tt.requestJSON))
			if err != nil {
				t.Error(err)
			}
			c.Params = append(
				c.Params,
				gin.Param{Key: "volumeName", Value: "volume"},
				gin.Param{Key: "key", Value: "test/sample.txt"},
			)
			if tt.isSetAccountID {
				c.Set("accountID", accountID)
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			entryUC := mockUsecase.NewMockEntryUsecase(ctrl)
			tt.setMockEntryUC(ctx, entryUC)

			hdl := handler.NewEntryHandler(entryUC)
			hdl.Update(c)

			c.Writer.WriteHeaderNow()

			if w.Code != tt.expectCode {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectCode, w.Code)
			}

			var response map[string]any
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Error(err)
			}

			if tt.expectCode == http.StatusOK {
				if size, ok := response["size"].(float64); ok {
					response["size"] = uint64(size)
				}
			}

			if diff := cmp.Diff(response, tt.expectResponse); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestEntry_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)

	accountID := uuid.New()

	tests := []struct {
		name           string
		isSetAccountID bool
		expectCode     int
		expectResponse *map[string]any
		setMockEntryUC func(context.Context, *mockUsecase.MockEntryUsecase)
	}{
		{
			name:           "success",
			isSetAccountID: true,
			expectCode:     http.StatusNoContent,
			expectResponse: nil,
			setMockEntryUC: func(ctx context.Context, entryUC *mockUsecase.MockEntryUsecase) {
				entryUC.
					EXPECT().
					Delete(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name:           "account id not found",
			isSetAccountID: false,
			expectCode:     http.StatusInternalServerError,
			expectResponse: &map[string]any{"message": "internal server error"},
			setMockEntryUC: func(context.Context, *mockUsecase.MockEntryUsecase) {},
		},
		{
			name:           "delete error",
			isSetAccountID: true,
			expectCode:     http.StatusInternalServerError,
			expectResponse: &map[string]any{"message": "internal server error"},
			setMockEntryUC: func(ctx context.Context, entryUC *mockUsecase.MockEntryUsecase) {
				entryUC.
					EXPECT().
					Delete(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(sql.ErrConnDone).
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
			c.Request, err = http.NewRequestWithContext(ctx, "DELETE", "entries/volume/test/sample.txt", http.NoBody)
			if err != nil {
				t.Error(err)
			}
			c.Params = append(
				c.Params,
				gin.Param{Key: "volumeName", Value: "volume"},
				gin.Param{Key: "key", Value: "test/sample.txt"},
			)
			if tt.isSetAccountID {
				c.Set("accountID", accountID)
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			entryUC := mockUsecase.NewMockEntryUsecase(ctrl)
			tt.setMockEntryUC(ctx, entryUC)

			hdl := handler.NewEntryHandler(entryUC)
			hdl.Delete(c)

			c.Writer.WriteHeaderNow()

			if w.Code != tt.expectCode {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectCode, w.Code)
			}

			if tt.expectCode != http.StatusNoContent {
				var response map[string]any
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Error(err)
				}
				if diff := cmp.Diff(&response, tt.expectResponse); diff != "" {
					t.Error(diff)
				}
			}
		})
	}
}

func TestEntry_Head(t *testing.T) {
	gin.SetMode(gin.TestMode)

	id := uuid.New()
	accountID := uuid.New()
	volumeID := uuid.New()
	entryDTO := &dto.EntryDTO{
		ID:        id,
		AccountID: accountID,
		VolumeID:  volumeID,
		Key:       "test/sample.txt",
		Size:      4,
		Type:      "text/plain; charset=utf-8",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name            string
		isSetAccountID  bool
		expectCode      int
		expectSize      string
		expectType      string
		expectUpdatedAt string
		setMockEntryUC  func(context.Context, *mockUsecase.MockEntryUsecase)
	}{
		{
			name:            "success",
			isSetAccountID:  true,
			expectCode:      http.StatusOK,
			expectSize:      strconv.FormatUint(entryDTO.Size, 10),
			expectType:      entryDTO.Type,
			expectUpdatedAt: entryDTO.UpdatedAt.Format(http.TimeFormat),
			setMockEntryUC: func(ctx context.Context, entryUC *mockUsecase.MockEntryUsecase) {
				entryUC.
					EXPECT().
					Head(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(entryDTO, nil).
					Times(1)
			},
		},
		{
			name:            "account id not found",
			isSetAccountID:  false,
			expectCode:      http.StatusInternalServerError,
			expectType:      "",
			expectSize:      "0",
			expectUpdatedAt: "",
			setMockEntryUC:  func(context.Context, *mockUsecase.MockEntryUsecase) {},
		},
		{
			name:            "head error",
			isSetAccountID:  true,
			expectCode:      http.StatusInternalServerError,
			expectType:      "",
			expectSize:      "0",
			expectUpdatedAt: "",
			setMockEntryUC: func(ctx context.Context, entryUC *mockUsecase.MockEntryUsecase) {
				entryUC.
					EXPECT().
					Head(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
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
			c.Request, err = http.NewRequestWithContext(ctx, "HEAD", "/entries/volume/test/sample.txt", http.NoBody)
			if err != nil {
				t.Error(err)
			}
			c.Params = append(
				c.Params,
				gin.Param{Key: "volumeName", Value: "volume"},
				gin.Param{Key: "key", Value: "test/sample.txt"},
			)
			if tt.isSetAccountID {
				c.Set("accountID", accountID)
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			entryUC := mockUsecase.NewMockEntryUsecase(ctrl)
			tt.setMockEntryUC(ctx, entryUC)

			hdl := handler.NewEntryHandler(entryUC)
			hdl.Head(c)

			c.Writer.WriteHeaderNow()

			if w.Code != tt.expectCode {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectCode, w.Code)
			}

			size := w.Header().Get("Content-Length")
			if size != tt.expectSize {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectSize, size)
			}
			contentType := w.Header().Get("Content-Type")
			if contentType != tt.expectType {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectType, contentType)
			}
			updatedAt := w.Header().Get("Last-Modified")
			if updatedAt != tt.expectUpdatedAt {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectUpdatedAt, updatedAt)
			}
		})
	}
}

func TestEntry_GetOne(t *testing.T) {
	gin.SetMode(gin.TestMode)

	id := uuid.New()
	accountID := uuid.New()
	volumeID := uuid.New()
	entryDTO := &dto.EntryDTO{
		ID:        id,
		AccountID: accountID,
		VolumeID:  volumeID,
		Key:       "test/sample.txt",
		Size:      4,
		Type:      "text/plain; charset=utf-8",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name           string
		isSetAccountID bool
		expectCode     int
		expectResponse []byte
		setMockEntryUC func(context.Context, *mockUsecase.MockEntryUsecase)
	}{
		{
			name:           "success",
			isSetAccountID: true,
			expectCode:     http.StatusOK,
			expectResponse: []byte("test"),
			setMockEntryUC: func(ctx context.Context, entryUC *mockUsecase.MockEntryUsecase) {
				entryUC.
					EXPECT().
					GetOne(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(entryDTO, io.NopCloser(bytes.NewReader([]byte("test"))), nil).
					Times(1)
			},
		},
		{
			name:           "account id not found",
			isSetAccountID: false,
			expectCode:     http.StatusInternalServerError,
			expectResponse: []byte(`{"message":"internal server error"}`),
			setMockEntryUC: func(context.Context, *mockUsecase.MockEntryUsecase) {},
		},
		{
			name:           "get error",
			isSetAccountID: true,
			expectCode:     http.StatusInternalServerError,
			expectResponse: []byte(`{"message":"internal server error"}`),
			setMockEntryUC: func(ctx context.Context, entryUC *mockUsecase.MockEntryUsecase) {
				entryUC.
					EXPECT().
					GetOne(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, nil, sql.ErrConnDone).
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
			c.Request, err = http.NewRequestWithContext(ctx, "GET", "/entries/volume/test/sample.txt", http.NoBody)
			if err != nil {
				t.Error(err)
			}
			c.Params = append(
				c.Params,
				gin.Param{Key: "volumeName", Value: "volume"},
				gin.Param{Key: "key", Value: "test/sample.txt"},
			)
			if tt.isSetAccountID {
				c.Set("accountID", accountID)
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			entryUC := mockUsecase.NewMockEntryUsecase(ctrl)
			tt.setMockEntryUC(ctx, entryUC)

			hdl := handler.NewEntryHandler(entryUC)
			hdl.GetOne(c)

			c.Writer.WriteHeaderNow()

			if w.Code != tt.expectCode {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectCode, w.Code)
			}

			resp := w.Result()
			response, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Error(err)
			}

			if diff := cmp.Diff(response, tt.expectResponse); diff != "" {
				t.Error(diff)
			}
		})
	}
}
