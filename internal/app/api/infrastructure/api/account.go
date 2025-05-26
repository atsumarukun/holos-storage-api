package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/repository"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/infrastructure/api/model"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/infrastructure/api/pkg/errors"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/infrastructure/api/transformer"
)

type accountRepository struct {
	client   *http.Client
	endpoint string
}

func NewAccountRepository(client *http.Client, endpoint string) repository.AccountRepository {
	return &accountRepository{
		client:   client,
		endpoint: endpoint,
	}
}

func (r *accountRepository) FindOneByCredential(ctx context.Context, credential string) (account *entity.Account, err error) {
	req, err := http.NewRequestWithContext(ctx, "GET", r.endpoint, http.NoBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", credential)

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
		return nil, repository.ErrUnauthorized
	} else if resp.StatusCode != http.StatusOK {
		return nil, errors.Decode(resp)
	}

	var model model.AccountModel
	if err := json.NewDecoder(resp.Body).Decode(&model); err != nil {
		return nil, err
	}

	return transformer.ToAccountEntity(&model), nil
}
