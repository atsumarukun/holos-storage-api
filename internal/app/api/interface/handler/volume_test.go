package handler_test

import (
	"bytes"
	"database/sql"
	"fmt"
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

	accountID := uuid.New()
	volumeDTO := &dto.VolumeDTO{
		ID:        uuid.New(),
		AccountID: accountID,
		Name:      "name",
		IsPublic:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name                  string
		requestBody           []byte
		hasAccountIDInContext bool
		expectCode            int
		expectResponse        []byte
		setMockVolumeUC       func(*mockUsecase.MockVolumeUsecase)
	}{
		{
			name:                  "successfully created",
			requestBody:           []byte(`{"name":"name","is_public":false}`),
			hasAccountIDInContext: true,
			expectCode:            http.StatusCreated,
			expectResponse:        fmt.Appendf(nil, `{"name":"%s","is_public":%t,"created_at":"%s","updated_at":"%s"}`, volumeDTO.Name, volumeDTO.IsPublic, volumeDTO.CreatedAt.Format(time.RFC3339Nano), volumeDTO.UpdatedAt.Format(time.RFC3339Nano)),
			setMockVolumeUC: func(volumeUC *mockUsecase.MockVolumeUsecase) {
				volumeUC.
					EXPECT().
					Create(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(volumeDTO, nil).
					Times(1)
			},
		},
		{
			name:                  "invalid request",
			requestBody:           nil,
			hasAccountIDInContext: true,
			expectCode:            http.StatusBadRequest,
			expectResponse:        []byte(`{"message":"failed to parse json"}`),
			setMockVolumeUC:       func(*mockUsecase.MockVolumeUsecase) {},
		},
		{
			name:                  "account id not set",
			requestBody:           []byte(`{"name":"name","is_public":false}`),
			hasAccountIDInContext: false,
			expectCode:            http.StatusInternalServerError,
			expectResponse:        []byte(`{"message":"internal server error"}`),
			setMockVolumeUC:       func(*mockUsecase.MockVolumeUsecase) {},
		},
		{
			name:                  "create error",
			requestBody:           []byte(`{"name":"name","is_public":false}`),
			hasAccountIDInContext: true,
			expectCode:            http.StatusInternalServerError,
			expectResponse:        []byte(`{"message":"internal server error"}`),
			setMockVolumeUC: func(volumeUC *mockUsecase.MockVolumeUsecase) {
				volumeUC.
					EXPECT().
					Create(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
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
			c.Request, err = http.NewRequestWithContext(ctx, "POST", "/volumes", bytes.NewBuffer(tt.requestBody))
			if err != nil {
				t.Error(err)
			}
			if tt.hasAccountIDInContext {
				c.Set("accountID", accountID)
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			volumeUC := mockUsecase.NewMockVolumeUsecase(ctrl)
			tt.setMockVolumeUC(volumeUC)

			hdl := handler.NewVolumeHandler(volumeUC)
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

func TestVolume_Update(t *testing.T) {
	gin.SetMode(gin.TestMode)

	accountID := uuid.New()
	volumeDTO := &dto.VolumeDTO{
		ID:        uuid.New(),
		AccountID: accountID,
		Name:      "name",
		IsPublic:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name                  string
		requestBody           []byte
		hasAccountIDInContext bool
		expectCode            int
		expectResponse        []byte
		setMockVolumeUC       func(*mockUsecase.MockVolumeUsecase)
	}{
		{
			name:                  "successfully updated",
			requestBody:           []byte(`{"name": "name", "is_public": false}`),
			hasAccountIDInContext: true,
			expectCode:            http.StatusOK,
			expectResponse:        fmt.Appendf(nil, `{"name":"%s","is_public":%t,"created_at":"%s","updated_at":"%s"}`, volumeDTO.Name, volumeDTO.IsPublic, volumeDTO.CreatedAt.Format(time.RFC3339Nano), volumeDTO.UpdatedAt.Format(time.RFC3339Nano)),
			setMockVolumeUC: func(volumeUC *mockUsecase.MockVolumeUsecase) {
				volumeUC.
					EXPECT().
					Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(volumeDTO, nil).
					Times(1)
			},
		},
		{
			name:                  "invalid request",
			requestBody:           nil,
			hasAccountIDInContext: true,
			expectCode:            http.StatusBadRequest,
			expectResponse:        []byte(`{"message":"failed to parse json"}`),
			setMockVolumeUC:       func(*mockUsecase.MockVolumeUsecase) {},
		},
		{
			name:                  "account id not set",
			requestBody:           []byte(`{"name": "name", "is_public": false}`),
			hasAccountIDInContext: false,
			expectCode:            http.StatusInternalServerError,
			expectResponse:        []byte(`{"message":"internal server error"}`),
			setMockVolumeUC:       func(*mockUsecase.MockVolumeUsecase) {},
		},
		{
			name:                  "update error",
			requestBody:           []byte(`{"name": "name", "is_public": false}`),
			hasAccountIDInContext: true,
			expectCode:            http.StatusInternalServerError,
			expectResponse:        []byte(`{"message":"internal server error"}`),
			setMockVolumeUC: func(volumeUC *mockUsecase.MockVolumeUsecase) {
				volumeUC.
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
			c.Request, err = http.NewRequestWithContext(ctx, "PUT", "/volumes/name", bytes.NewBuffer(tt.requestBody))
			if err != nil {
				t.Error(err)
			}
			c.Params = append(c.Params, gin.Param{Key: "name", Value: "name"})
			if tt.hasAccountIDInContext {
				c.Set("accountID", accountID)
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			volumeUC := mockUsecase.NewMockVolumeUsecase(ctrl)
			tt.setMockVolumeUC(volumeUC)

			hdl := handler.NewVolumeHandler(volumeUC)
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

func TestVolume_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)

	accountID := uuid.New()

	tests := []struct {
		name                  string
		hasAccountIDInContext bool
		expectCode            int
		expectResponse        []byte
		setMockVolumeUC       func(*mockUsecase.MockVolumeUsecase)
	}{
		{
			name:                  "successfully deleted",
			hasAccountIDInContext: true,
			expectCode:            http.StatusNoContent,
			expectResponse:        nil,
			setMockVolumeUC: func(volumeUC *mockUsecase.MockVolumeUsecase) {
				volumeUC.
					EXPECT().
					Delete(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name:                  "account id not set",
			hasAccountIDInContext: false,
			expectCode:            http.StatusInternalServerError,
			expectResponse:        []byte(`{"message":"internal server error"}`),
			setMockVolumeUC:       func(*mockUsecase.MockVolumeUsecase) {},
		},
		{
			name:                  "delete error",
			hasAccountIDInContext: true,
			expectCode:            http.StatusInternalServerError,
			expectResponse:        []byte(`{"message":"internal server error"}`),
			setMockVolumeUC: func(volumeUC *mockUsecase.MockVolumeUsecase) {
				volumeUC.
					EXPECT().
					Delete(gomock.Any(), gomock.Any(), gomock.Any()).
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
			c.Request, err = http.NewRequestWithContext(ctx, "DELETE", "/volumes/name", http.NoBody)
			if err != nil {
				t.Error(err)
			}
			c.Params = append(c.Params, gin.Param{Key: "name", Value: "name"})
			if tt.hasAccountIDInContext {
				c.Set("accountID", accountID)
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			volumeUC := mockUsecase.NewMockVolumeUsecase(ctrl)
			tt.setMockVolumeUC(volumeUC)

			hdl := handler.NewVolumeHandler(volumeUC)
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

func TestVolume_GetOne(t *testing.T) {
	gin.SetMode(gin.TestMode)

	accountID := uuid.New()
	volumeDTO := &dto.VolumeDTO{
		ID:        uuid.New(),
		AccountID: accountID,
		Name:      "name",
		IsPublic:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name                  string
		hasAccountIDInContext bool
		expectCode            int
		expectResponse        []byte
		setMockVolumeUC       func(*mockUsecase.MockVolumeUsecase)
	}{
		{
			name:                  "successfully got one",
			hasAccountIDInContext: true,
			expectCode:            http.StatusOK,
			expectResponse:        fmt.Appendf(nil, `{"name":"%s","is_public":%t,"created_at":"%s","updated_at":"%s"}`, volumeDTO.Name, volumeDTO.IsPublic, volumeDTO.CreatedAt.Format(time.RFC3339Nano), volumeDTO.UpdatedAt.Format(time.RFC3339Nano)),
			setMockVolumeUC: func(volumeUC *mockUsecase.MockVolumeUsecase) {
				volumeUC.
					EXPECT().
					GetOne(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(volumeDTO, nil).
					Times(1)
			},
		},
		{
			name:                  "account id not set",
			hasAccountIDInContext: false,
			expectCode:            http.StatusInternalServerError,
			expectResponse:        []byte(`{"message":"internal server error"}`),
			setMockVolumeUC:       func(*mockUsecase.MockVolumeUsecase) {},
		},
		{
			name:                  "get error",
			hasAccountIDInContext: true,
			expectCode:            http.StatusInternalServerError,
			expectResponse:        []byte(`{"message":"internal server error"}`),
			setMockVolumeUC: func(volumeUC *mockUsecase.MockVolumeUsecase) {
				volumeUC.
					EXPECT().
					GetOne(gomock.Any(), gomock.Any(), gomock.Any()).
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
			c.Request, err = http.NewRequestWithContext(ctx, "GET", "/volumes/name", http.NoBody)
			if err != nil {
				t.Error(err)
			}
			c.Params = append(c.Params, gin.Param{Key: "name", Value: "name"})
			if tt.hasAccountIDInContext {
				c.Set("accountID", accountID)
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			volumeUC := mockUsecase.NewMockVolumeUsecase(ctrl)
			tt.setMockVolumeUC(volumeUC)

			hdl := handler.NewVolumeHandler(volumeUC)
			hdl.GetOne(c)

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

func TestVolume_GetAll(t *testing.T) {
	gin.SetMode(gin.TestMode)

	accountID := uuid.New()
	volumeDTO := &dto.VolumeDTO{
		ID:        uuid.New(),
		AccountID: accountID,
		Name:      "name",
		IsPublic:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name                  string
		hasAccountIDInContext bool
		expectCode            int
		expectResponse        []byte
		setMockVolumeUC       func(*mockUsecase.MockVolumeUsecase)
	}{
		{
			name:                  "successfully got all",
			hasAccountIDInContext: true,
			expectCode:            http.StatusOK,
			expectResponse:        fmt.Appendf(nil, `{"volumes":[{"name":"%s","is_public":%t,"created_at":"%s","updated_at":"%s"}]}`, volumeDTO.Name, volumeDTO.IsPublic, volumeDTO.CreatedAt.Format(time.RFC3339Nano), volumeDTO.UpdatedAt.Format(time.RFC3339Nano)),
			setMockVolumeUC: func(volumeUC *mockUsecase.MockVolumeUsecase) {
				volumeUC.
					EXPECT().
					GetAll(gomock.Any(), gomock.Any()).
					Return([]*dto.VolumeDTO{volumeDTO}, nil).
					Times(1)
			},
		},
		{
			name:                  "not found",
			hasAccountIDInContext: true,
			expectCode:            http.StatusOK,
			expectResponse:        []byte(`{"volumes":[]}`),
			setMockVolumeUC: func(volumeUC *mockUsecase.MockVolumeUsecase) {
				volumeUC.
					EXPECT().
					GetAll(gomock.Any(), gomock.Any()).
					Return([]*dto.VolumeDTO{}, nil).
					Times(1)
			},
		},
		{
			name:                  "account id not set",
			hasAccountIDInContext: false,
			expectCode:            http.StatusInternalServerError,
			expectResponse:        []byte(`{"message":"internal server error"}`),
			setMockVolumeUC:       func(*mockUsecase.MockVolumeUsecase) {},
		},
		{
			name:                  "get error",
			hasAccountIDInContext: true,
			expectCode:            http.StatusInternalServerError,
			expectResponse:        []byte(`{"message":"internal server error"}`),
			setMockVolumeUC: func(volumeUC *mockUsecase.MockVolumeUsecase) {
				volumeUC.
					EXPECT().
					GetAll(gomock.Any(), gomock.Any()).
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
			c.Request, err = http.NewRequestWithContext(ctx, "GET", "/volumes", http.NoBody)
			if err != nil {
				t.Error(err)
			}
			if tt.hasAccountIDInContext {
				c.Set("accountID", accountID)
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			volumeUC := mockUsecase.NewMockVolumeUsecase(ctrl)
			tt.setMockVolumeUC(volumeUC)

			hdl := handler.NewVolumeHandler(volumeUC)
			hdl.GetAll(c)

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
