package api_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/infrastructure/api"
)

func TestAuthorization_Authorize(t *testing.T) {
	authorization := &entity.Authorization{
		AccountID: uuid.New(),
	}

	tests := []struct {
		name            string
		expectResult    *entity.Authorization
		expectError     error
		mockHandlerFunc http.HandlerFunc
	}{
		{
			name:         "success",
			expectResult: authorization,
			expectError:  nil,
			mockHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]string{"account_id": authorization.AccountID.String()})
			},
		},
		{
			name:         "unauthorized",
			expectResult: nil,
			expectError:  api.ErrUnauthorized,
			mockHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(tt.mockHandlerFunc)
			defer srv.Close()

			repo := api.NewAuthorizationRepository(srv.Client(), srv.URL)
			result, err := repo.Authorize(t.Context(), "Session: token")
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			if diff := cmp.Diff(result, tt.expectResult); diff != "" {
				t.Error(diff)
			}
		})
	}
}
