package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/repository"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/infrastructure/model"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/infrastructure/transformer"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/pkg/status"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/pkg/status/code"
)

var (
	ErrUnauthorized       = status.Error(code.Unauthorized, "unauthorized")
	ErrAuthorizationFaild = status.Error(code.Internal, "authorization faild")
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

func (r *authorizationRepository) Authorize(ctx context.Context, credential string) (authorization *entity.Authorization, err error) {
	req, err := http.NewRequestWithContext(ctx, "GET", r.endpoint, http.NoBody)
	if err != nil {
		return nil, err
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		// NOTE: errに直接詰めると関数内のエラーがnilで上書きされるためエラー発生時のみ上書きする.
		if e := resp.Body.Close(); e != nil {
			err = e
		}
	}()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, ErrUnauthorized
	} else if resp.StatusCode != http.StatusOK {
		return nil, ErrAuthorizationFaild
	}

	var model model.AuthorizationModel
	if err := json.NewDecoder(resp.Body).Decode(&model); err != nil {
		return nil, err
	}

	return transformer.ToAuthorizationEntity(&model), nil
}
