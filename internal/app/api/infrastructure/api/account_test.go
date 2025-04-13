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
	"github.com/atsumarukun/holos-storage-api/internal/app/api/pkg/status"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/pkg/status/code"
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
			name:            "success",
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
			name:            "unauthorized",
			inputCredential: "Session: YNDNun_KFu1uFmS691yJ6eqJ9eczRVKn",
			expectResult:    nil,
			expectError:     status.Error(code.Unauthorized, "unauthorized"),
			mockHandlerFunc: func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				body, err := json.Marshal(&map[string]string{"message": "unauthorized"})
				if err != nil {
					t.Error(err)
				}
				w.Write(body)
			},
		},
		{
			name:            "authorization faild",
			inputCredential: "Session: YNDNun_KFu1uFmS691yJ6eqJ9eczRVKn",
			expectResult:    nil,
			expectError:     status.Error(code.Internal, "internal server error"),
			mockHandlerFunc: func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				body, err := json.Marshal(&map[string]string{"message": "internal server error"})
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
			if !status.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			if diff := cmp.Diff(result, tt.expectResult); diff != "" {
				t.Error(diff)
			}
		})
	}
}
