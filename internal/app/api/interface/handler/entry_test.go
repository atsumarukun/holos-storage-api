package handler_test

import (
	"bytes"
	"database/sql"
	"fmt"
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
	if err := writer.WriteField("key", "key/sample.txt"); err != nil {
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
	entryDTO := &dto.EntryDTO{
		ID:        uuid.New(),
		AccountID: accountID,
		VolumeID:  uuid.New(),
		Key:       "key/sample.txt",
		Size:      4,
		Type:      "text/plain; charset=utf-8",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name                  string
		buildRequestBody      func(*testing.T) (io.Reader, string)
		hasAccountIDInContext bool
		expectCode            int
		expectResponse        []byte
		setMockEntryUC        func(*mockUsecase.MockEntryUsecase)
	}{
		{
			name:                  "successfully created",
			buildRequestBody:      buildMultipartBody,
			hasAccountIDInContext: true,
			expectCode:            http.StatusCreated,
			expectResponse:        fmt.Appendf(nil, `{"key":"%s","size":%d,"type":"%s","created_at":"%s","updated_at":"%s"}`, entryDTO.Key, entryDTO.Size, entryDTO.Type, entryDTO.CreatedAt.Format(time.RFC3339Nano), entryDTO.UpdatedAt.Format(time.RFC3339Nano)),
			setMockEntryUC: func(entryUC *mockUsecase.MockEntryUsecase) {
				entryUC.
					EXPECT().
					Create(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(entryDTO, nil).
					Times(1)
			},
		},
		{
			name:                  "invalid request",
			buildRequestBody:      func(*testing.T) (io.Reader, string) { return http.NoBody, "" },
			hasAccountIDInContext: true,
			expectCode:            http.StatusBadRequest,
			expectResponse:        []byte(`{"message":"failed to parse multipart/form-data"}`),
			setMockEntryUC:        func(*mockUsecase.MockEntryUsecase) {},
		},
		{
			name:                  "account id not set",
			buildRequestBody:      buildMultipartBody,
			hasAccountIDInContext: false,
			expectCode:            http.StatusInternalServerError,
			expectResponse:        []byte(`{"message":"internal server error"}`),
			setMockEntryUC:        func(*mockUsecase.MockEntryUsecase) {},
		},
		{
			name:                  "create error",
			buildRequestBody:      buildMultipartBody,
			hasAccountIDInContext: true,
			expectCode:            http.StatusInternalServerError,
			expectResponse:        []byte(`{"message":"internal server error"}`),
			setMockEntryUC: func(entryUC *mockUsecase.MockEntryUsecase) {
				entryUC.
					EXPECT().
					Create(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			w := httptest.NewRecorder()

			body, contentType := tt.buildRequestBody(t)

			c, _ := gin.CreateTestContext(w)
			var err error
			c.Request, err = http.NewRequestWithContext(ctx, "POST", "/entries/volume", body)
			if err != nil {
				t.Error(err)
			}
			c.Request.Header.Add("Content-Type", contentType)
			if tt.hasAccountIDInContext {
				c.Set("accountID", accountID)
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			entryUC := mockUsecase.NewMockEntryUsecase(ctrl)
			tt.setMockEntryUC(entryUC)

			hdl := handler.NewEntryHandler(entryUC)
			hdl.Create(c)

			c.Writer.WriteHeaderNow()

			if w.Code != tt.expectCode {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectCode, w.Code)
			}

			if diff := cmp.Diff(tt.expectResponse, w.Body.Bytes()); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestEntry_Update(t *testing.T) {
	gin.SetMode(gin.TestMode)

	accountID := uuid.New()
	entryDTO := &dto.EntryDTO{
		ID:        uuid.New(),
		AccountID: accountID,
		VolumeID:  uuid.New(),
		Key:       "key/sample.txt",
		Size:      4,
		Type:      "text/plain; charset=utf-8",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name                  string
		requestBody           []byte
		hasAccountIDInContext bool
		expectCode            int
		expectResponse        []byte
		setMockEntryUC        func(*mockUsecase.MockEntryUsecase)
	}{
		{
			name:                  "successfully updated",
			requestBody:           []byte(`{"key": "update/sample.txt"}`),
			hasAccountIDInContext: true,
			expectCode:            http.StatusOK,
			expectResponse:        fmt.Appendf(nil, `{"key":"%s","size":%d,"type":"%s","created_at":"%s","updated_at":"%s"}`, entryDTO.Key, entryDTO.Size, entryDTO.Type, entryDTO.CreatedAt.Format(time.RFC3339Nano), entryDTO.UpdatedAt.Format(time.RFC3339Nano)),
			setMockEntryUC: func(entryUC *mockUsecase.MockEntryUsecase) {
				entryUC.
					EXPECT().
					Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(entryDTO, nil).
					Times(1)
			},
		},
		{
			name:                  "invalid request",
			requestBody:           nil,
			hasAccountIDInContext: true,
			expectCode:            http.StatusBadRequest,
			expectResponse:        []byte(`{"message":"failed to parse json"}`),
			setMockEntryUC:        func(*mockUsecase.MockEntryUsecase) {},
		},
		{
			name:                  "account id not set",
			requestBody:           []byte(`{"key": "update/sample.txt"}`),
			hasAccountIDInContext: false,
			expectCode:            http.StatusInternalServerError,
			expectResponse:        []byte(`{"message":"internal server error"}`),
			setMockEntryUC:        func(*mockUsecase.MockEntryUsecase) {},
		},
		{
			name:                  "update error",
			requestBody:           []byte(`{"key": "update/sample.txt"}`),
			hasAccountIDInContext: true,
			expectCode:            http.StatusInternalServerError,
			expectResponse:        []byte(`{"message":"internal server error"}`),
			setMockEntryUC: func(entryUC *mockUsecase.MockEntryUsecase) {
				entryUC.
					EXPECT().
					Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
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
			c.Request, err = http.NewRequestWithContext(ctx, "PUT", "/entries/volume/key/sample.txt", bytes.NewBuffer(tt.requestBody))
			if err != nil {
				t.Error(err)
			}
			c.Params = append(
				c.Params,
				gin.Param{Key: "volumeName", Value: "volume"},
				gin.Param{Key: "key", Value: "key/sample.txt"},
			)
			if tt.hasAccountIDInContext {
				c.Set("accountID", accountID)
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			entryUC := mockUsecase.NewMockEntryUsecase(ctrl)
			tt.setMockEntryUC(entryUC)

			hdl := handler.NewEntryHandler(entryUC)
			hdl.Update(c)

			c.Writer.WriteHeaderNow()

			if w.Code != tt.expectCode {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectCode, w.Code)
			}

			if diff := cmp.Diff(tt.expectResponse, w.Body.Bytes()); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestEntry_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)

	accountID := uuid.New()

	tests := []struct {
		name                  string
		hasAccountIDInContext bool
		expectCode            int
		expectResponse        []byte
		setMockEntryUC        func(*mockUsecase.MockEntryUsecase)
	}{
		{
			name:                  "successfully deleted",
			hasAccountIDInContext: true,
			expectCode:            http.StatusNoContent,
			expectResponse:        nil,
			setMockEntryUC: func(entryUC *mockUsecase.MockEntryUsecase) {
				entryUC.
					EXPECT().
					Delete(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name:                  "account id not set",
			hasAccountIDInContext: false,
			expectCode:            http.StatusInternalServerError,
			expectResponse:        []byte(`{"message":"internal server error"}`),
			setMockEntryUC:        func(*mockUsecase.MockEntryUsecase) {},
		},
		{
			name:                  "delete error",
			hasAccountIDInContext: true,
			expectCode:            http.StatusInternalServerError,
			expectResponse:        []byte(`{"message":"internal server error"}`),
			setMockEntryUC: func(entryUC *mockUsecase.MockEntryUsecase) {
				entryUC.
					EXPECT().
					Delete(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
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
			c.Request, err = http.NewRequestWithContext(ctx, "DELETE", "entries/volume/key/sample.txt", http.NoBody)
			if err != nil {
				t.Error(err)
			}
			c.Params = append(
				c.Params,
				gin.Param{Key: "volumeName", Value: "volume"},
				gin.Param{Key: "key", Value: "key/sample.txt"},
			)
			if tt.hasAccountIDInContext {
				c.Set("accountID", accountID)
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			entryUC := mockUsecase.NewMockEntryUsecase(ctrl)
			tt.setMockEntryUC(entryUC)

			hdl := handler.NewEntryHandler(entryUC)
			hdl.Delete(c)

			c.Writer.WriteHeaderNow()

			if w.Code != tt.expectCode {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectCode, w.Code)
			}

			if diff := cmp.Diff(tt.expectResponse, w.Body.Bytes()); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestEntry_Copy(t *testing.T) {
	gin.SetMode(gin.TestMode)

	accountID := uuid.New()
	entryDTO := &dto.EntryDTO{
		ID:        uuid.New(),
		AccountID: accountID,
		VolumeID:  uuid.New(),
		Key:       "key/sample copy.txt",
		Size:      4,
		Type:      "text/plain; charset=utf-8",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name                  string
		hasAccountIDInContext bool
		expectCode            int
		expectResponse        []byte
		setMockEntryUC        func(*mockUsecase.MockEntryUsecase)
	}{
		{
			name:                  "successfully copied",
			hasAccountIDInContext: true,
			expectCode:            http.StatusOK,
			expectResponse:        fmt.Appendf(nil, `{"key":"%s","size":%d,"type":"%s","created_at":"%s","updated_at":"%s"}`, entryDTO.Key, entryDTO.Size, entryDTO.Type, entryDTO.CreatedAt.Format(time.RFC3339Nano), entryDTO.UpdatedAt.Format(time.RFC3339Nano)),
			setMockEntryUC: func(entryUC *mockUsecase.MockEntryUsecase) {
				entryUC.
					EXPECT().
					Copy(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(entryDTO, nil).
					Times(1)
			},
		},
		{
			name:                  "account id not set",
			hasAccountIDInContext: false,
			expectCode:            http.StatusInternalServerError,
			expectResponse:        []byte(`{"message":"internal server error"}`),
			setMockEntryUC:        func(*mockUsecase.MockEntryUsecase) {},
		},
		{
			name:                  "copy error",
			hasAccountIDInContext: true,
			expectCode:            http.StatusInternalServerError,
			expectResponse:        []byte(`{"message":"internal server error"}`),
			setMockEntryUC: func(entryUC *mockUsecase.MockEntryUsecase) {
				entryUC.
					EXPECT().
					Copy(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
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
			c.Request, err = http.NewRequestWithContext(ctx, "POST", "entries/volume/key/sample.txt", http.NoBody)
			if err != nil {
				t.Error(err)
			}
			c.Params = append(
				c.Params,
				gin.Param{Key: "volumeName", Value: "volume"},
				gin.Param{Key: "key", Value: "key/sample.txt"},
			)
			if tt.hasAccountIDInContext {
				c.Set("accountID", accountID)
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			entryUC := mockUsecase.NewMockEntryUsecase(ctrl)
			tt.setMockEntryUC(entryUC)

			hdl := handler.NewEntryHandler(entryUC)
			hdl.Copy(c)

			c.Writer.WriteHeaderNow()

			if w.Code != tt.expectCode {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectCode, w.Code)
			}

			if diff := cmp.Diff(tt.expectResponse, w.Body.Bytes()); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestEntry_GetMeta(t *testing.T) {
	gin.SetMode(gin.TestMode)

	accountID := uuid.New()
	volumeID := uuid.New()
	fileEntryDTO := &dto.EntryDTO{
		ID:        uuid.New(),
		AccountID: accountID,
		VolumeID:  volumeID,
		Key:       "key/sample.txt",
		Size:      4,
		Type:      "text/plain; charset=utf-8",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	folderEntryDTO := &dto.EntryDTO{
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
		name                  string
		hasAccountIDInContext bool
		expectCode            int
		expectHeader          http.Header
		setMockEntryUC        func(*mockUsecase.MockEntryUsecase)
	}{
		{
			name:                  "successfully got a file meta",
			hasAccountIDInContext: true,
			expectCode:            http.StatusOK,
			expectHeader:          http.Header{"Content-Length": {strconv.FormatUint(fileEntryDTO.Size, 10)}, "Content-Type": {fileEntryDTO.Type}, "Holos-Entry-Type": {fileEntryDTO.Type}, "Last-Modified": {fileEntryDTO.UpdatedAt.Format(http.TimeFormat)}},
			setMockEntryUC: func(entryUC *mockUsecase.MockEntryUsecase) {
				entryUC.
					EXPECT().
					GetMeta(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(fileEntryDTO, nil).
					Times(1)
			},
		},
		{
			name:                  "successfully got a folder meta",
			hasAccountIDInContext: true,
			expectCode:            http.StatusOK,
			expectHeader:          http.Header{"Content-Length": {strconv.FormatUint(folderEntryDTO.Size, 10)}, "Content-Type": {"application/octet-stream"}, "Holos-Entry-Type": {folderEntryDTO.Type}, "Last-Modified": {folderEntryDTO.UpdatedAt.Format(http.TimeFormat)}},
			setMockEntryUC: func(entryUC *mockUsecase.MockEntryUsecase) {
				entryUC.
					EXPECT().
					GetMeta(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(folderEntryDTO, nil).
					Times(1)
			},
		},
		{
			name:                  "account id not set",
			hasAccountIDInContext: false,
			expectCode:            http.StatusInternalServerError,
			expectHeader:          http.Header{},
			setMockEntryUC:        func(*mockUsecase.MockEntryUsecase) {},
		},
		{
			name:                  "get error",
			hasAccountIDInContext: true,
			expectCode:            http.StatusInternalServerError,
			expectHeader:          http.Header{},
			setMockEntryUC: func(entryUC *mockUsecase.MockEntryUsecase) {
				entryUC.
					EXPECT().
					GetMeta(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
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
			c.Request, err = http.NewRequestWithContext(ctx, "HEAD", "/entries/volume/key/sample.txt", http.NoBody)
			if err != nil {
				t.Error(err)
			}
			c.Params = append(
				c.Params,
				gin.Param{Key: "volumeName", Value: "volume"},
				gin.Param{Key: "key", Value: "key/sample.txt"},
			)
			if tt.hasAccountIDInContext {
				c.Set("accountID", accountID)
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			entryUC := mockUsecase.NewMockEntryUsecase(ctrl)
			tt.setMockEntryUC(entryUC)

			hdl := handler.NewEntryHandler(entryUC)
			hdl.GetMeta(c)

			c.Writer.WriteHeaderNow()

			if w.Code != tt.expectCode {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectCode, w.Code)
			}

			if diff := cmp.Diff(tt.expectHeader, w.Header()); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestEntry_GetOne(t *testing.T) {
	gin.SetMode(gin.TestMode)

	accountID := uuid.New()
	volumeID := uuid.New()
	fileEntryDTO := &dto.EntryDTO{
		ID:        uuid.New(),
		AccountID: accountID,
		VolumeID:  volumeID,
		Key:       "key/sample.txt",
		Size:      4,
		Type:      "text/plain; charset=utf-8",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	folderEntryDTO := &dto.EntryDTO{
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
		name                  string
		hasAccountIDInContext bool
		expectCode            int
		expectHeader          http.Header
		expectResponse        []byte
		setMockEntryUC        func(*mockUsecase.MockEntryUsecase)
	}{
		{
			name:                  "successfully got a file",
			hasAccountIDInContext: true,
			expectCode:            http.StatusOK,
			expectHeader:          http.Header{"Content-Length": {strconv.FormatUint(fileEntryDTO.Size, 10)}, "Content-Type": {fileEntryDTO.Type}, "Holos-Entry-Type": {fileEntryDTO.Type}, "Last-Modified": {fileEntryDTO.UpdatedAt.Format(http.TimeFormat)}},
			expectResponse:        []byte("test"),
			setMockEntryUC: func(entryUC *mockUsecase.MockEntryUsecase) {
				entryUC.
					EXPECT().
					GetOne(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(fileEntryDTO, io.NopCloser(bytes.NewReader([]byte("test"))), nil).
					Times(1)
			},
		},
		{
			name:                  "successfully got a folder",
			hasAccountIDInContext: true,
			expectCode:            http.StatusOK,
			expectHeader:          http.Header{"Content-Length": {strconv.FormatUint(folderEntryDTO.Size, 10)}, "Content-Type": {"application/octet-stream"}, "Holos-Entry-Type": {folderEntryDTO.Type}, "Last-Modified": {folderEntryDTO.UpdatedAt.Format(http.TimeFormat)}},
			expectResponse:        nil,
			setMockEntryUC: func(entryUC *mockUsecase.MockEntryUsecase) {
				entryUC.
					EXPECT().
					GetOne(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(folderEntryDTO, nil, nil).
					Times(1)
			},
		},
		{
			name:                  "account id not set",
			hasAccountIDInContext: false,
			expectCode:            http.StatusInternalServerError,
			expectHeader:          http.Header{"Content-Type": {"application/json; charset=utf-8"}},
			expectResponse:        []byte(`{"message":"internal server error"}`),
			setMockEntryUC:        func(*mockUsecase.MockEntryUsecase) {},
		},
		{
			name:                  "get error",
			hasAccountIDInContext: true,
			expectCode:            http.StatusInternalServerError,
			expectHeader:          http.Header{"Content-Type": {"application/json; charset=utf-8"}},
			expectResponse:        []byte(`{"message":"internal server error"}`),
			setMockEntryUC: func(entryUC *mockUsecase.MockEntryUsecase) {
				entryUC.
					EXPECT().
					GetOne(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
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
			c.Request, err = http.NewRequestWithContext(ctx, "GET", "/entries/volume/key/sample.txt", http.NoBody)
			if err != nil {
				t.Error(err)
			}
			c.Params = append(
				c.Params,
				gin.Param{Key: "volumeName", Value: "volume"},
				gin.Param{Key: "key", Value: "key/sample.txt"},
			)
			if tt.hasAccountIDInContext {
				c.Set("accountID", accountID)
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			entryUC := mockUsecase.NewMockEntryUsecase(ctrl)
			tt.setMockEntryUC(entryUC)

			hdl := handler.NewEntryHandler(entryUC)
			hdl.GetOne(c)

			c.Writer.WriteHeaderNow()

			if w.Code != tt.expectCode {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectCode, w.Code)
			}

			if diff := cmp.Diff(tt.expectHeader, w.Header()); diff != "" {
				t.Error(diff)
			}

			if diff := cmp.Diff(tt.expectResponse, w.Body.Bytes()); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestEntry_Search(t *testing.T) {
	gin.SetMode(gin.TestMode)

	accountID := uuid.New()
	entryDTO := &dto.EntryDTO{
		ID:        uuid.New(),
		AccountID: accountID,
		VolumeID:  uuid.New(),
		Key:       "key/sample.txt",
		Size:      4,
		Type:      "text/plain; charset=utf-8",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name                  string
		hasAccountIDInContext bool
		expectCode            int
		expectResponse        []byte
		setMockEntryUC        func(*mockUsecase.MockEntryUsecase)
	}{
		{
			name:                  "successfully searched",
			hasAccountIDInContext: true,
			expectCode:            http.StatusOK,
			expectResponse:        fmt.Appendf(nil, `{"entries":[{"key":"%s","size":%d,"type":"%s","created_at":"%s","updated_at":"%s"}]}`, entryDTO.Key, entryDTO.Size, entryDTO.Type, entryDTO.CreatedAt.Format(time.RFC3339Nano), entryDTO.UpdatedAt.Format(time.RFC3339Nano)),
			setMockEntryUC: func(entryUC *mockUsecase.MockEntryUsecase) {
				entryUC.
					EXPECT().
					Search(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]*dto.EntryDTO{entryDTO}, nil).
					Times(1)
			},
		},
		{
			name:                  "not found",
			hasAccountIDInContext: true,
			expectCode:            http.StatusOK,
			expectResponse:        []byte(`{"entries":[]}`),
			setMockEntryUC: func(entryUC *mockUsecase.MockEntryUsecase) {
				entryUC.
					EXPECT().
					Search(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]*dto.EntryDTO{}, nil).
					Times(1)
			},
		},
		{
			name:                  "account id not set",
			hasAccountIDInContext: false,
			expectCode:            http.StatusInternalServerError,
			expectResponse:        []byte(`{"message":"internal server error"}`),
			setMockEntryUC:        func(*mockUsecase.MockEntryUsecase) {},
		},
		{
			name:                  "search error",
			hasAccountIDInContext: true,
			expectCode:            http.StatusInternalServerError,
			expectResponse:        []byte(`{"message":"internal server error"}`),
			setMockEntryUC: func(entryUC *mockUsecase.MockEntryUsecase) {
				entryUC.
					EXPECT().
					Search(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
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
			c.Request, err = http.NewRequestWithContext(ctx, "GET", "/entries/volume", http.NoBody)
			if err != nil {
				t.Error(err)
			}
			c.Params = append(
				c.Params,
				gin.Param{Key: "volumeName", Value: "volume"},
			)
			if tt.hasAccountIDInContext {
				c.Set("accountID", accountID)
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			entryUC := mockUsecase.NewMockEntryUsecase(ctrl)
			tt.setMockEntryUC(entryUC)

			hdl := handler.NewEntryHandler(entryUC)
			hdl.Search(c)

			c.Writer.WriteHeaderNow()

			if w.Code != tt.expectCode {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectCode, w.Code)
			}

			if diff := cmp.Diff(tt.expectResponse, w.Body.Bytes()); diff != "" {
				t.Error(diff)
			}
		})
	}
}
