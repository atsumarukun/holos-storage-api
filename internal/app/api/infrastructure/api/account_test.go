package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/infrastructure/api"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/infrastructure/api/pkg/client"
	"github.com/atsumarukun/holos-storage-api/test/assert"
)

func TestAccount_FindOneByCredential(t *testing.T) {
	account := &entity.Account{
		ID: uuid.New(),
	}

	tests := []struct {
		name            string
		inputCredential string
		expectResult    *entity.Account
		expectError     error
		mockHandlerFunc http.HandlerFunc
	}{
		{
			name:            "successfully found",
			inputCredential: "Session: YNDNun_KFu1uFmS691yJ6eqJ9eczRVKn",
			expectResult:    account,
			expectError:     nil,
			mockHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				if r.Header.Get("Authorization") == "" {
					w.WriteHeader(http.StatusUnauthorized)
				}
				json.NewEncoder(w).Encode(map[string]string{"id": account.ID.String()})
			},
		},
		{
			name:            "unauthenticated",
			inputCredential: "Session: YNDNun_KFu1uFmS691yJ6eqJ9eczRVKn",
			expectResult:    nil,
			expectError:     client.ErrUnauthenticated,
			mockHandlerFunc: func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				body, err := json.Marshal(&map[string]map[string]string{"error": {"code": "UNAUTHENTICATED", "message": "unauthenticated"}})
				if err != nil {
					t.Error(err)
				}
				w.Write(body)
			},
		},
		{
			name:            "unauthorized",
			inputCredential: "Session: YNDNun_KFu1uFmS691yJ6eqJ9eczRVKn",
			expectResult:    nil,
			expectError:     client.ErrUnauthorized,
			mockHandlerFunc: func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				body, err := json.Marshal(&map[string]map[string]string{"error": {"code": "UNAUTHORIZED", "message": "unauthorized"}})
				if err != nil {
					t.Error(err)
				}
				w.Write(body)
			},
		},
		{
			name:            "internal server error",
			inputCredential: "Session: YNDNun_KFu1uFmS691yJ6eqJ9eczRVKn",
			expectResult:    nil,
			expectError:     client.ErrInternalServerError,
			mockHandlerFunc: func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				body, err := json.Marshal(&map[string]map[string]string{"error": {"code": "INTERNAL_SERVER_ERROR", "message": "internal server error"}})
				if err != nil {
					t.Error(err)
				}
				w.Write(body)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(tt.mockHandlerFunc)
			defer srv.Close()

			repo := api.NewAccountRepository(srv.Client(), srv.URL)
			result, err := repo.FindOneByCredential(t.Context(), tt.inputCredential)

			assert.Error(t, err, tt.expectError)

			if diff := cmp.Diff(tt.expectResult, result); diff != "" {
				t.Error(diff)
			}
		})
	}
}
