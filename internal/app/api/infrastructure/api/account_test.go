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
			inputCredential: "Session: token",
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
			inputCredential: "Session: token",
			expectResult:    nil,
			expectError:     api.ErrUnauthorized,
			mockHandlerFunc: func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
			},
		},
		{
			name:            "authorization faild",
			inputCredential: "Session: token",
			expectResult:    nil,
			expectError:     api.ErrAuthorizationFaild,
			mockHandlerFunc: func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(tt.mockHandlerFunc)
			defer srv.Close()

			repo := api.NewAccountRepository(srv.Client(), srv.URL)
			result, err := repo.FindOneByCredential(t.Context(), tt.inputCredential)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			if diff := cmp.Diff(result, tt.expectResult); diff != "" {
				t.Error(diff)
			}
		})
	}
}
