package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/repository"
)

type authorizationRepository struct {
	client   *http.Client
	endpoint string
}

func NewAuthorizationRepository(client *http.Client, endpoint string) repository.AuthorizationRepository {
	return &authorizationRepository{
		client:   client,
		endpoint: endpoint,
	}
}

func (r *authorizationRepository) Authorize(ctx context.Context, credential string) (*entity.Authorization, error) {
	return nil, errors.New("not implemented")
}
