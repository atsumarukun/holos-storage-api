package api

import (
	"context"
	stderr "errors"
	"net/http"

	"github.com/atsumarukun/holos-api-pkg/errors"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/repository"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/infrastructure/api/model"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/infrastructure/api/pkg/client"
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
	const errMessage = "failed to find account by credential"

	req, err := http.NewRequestWithContext(ctx, "GET", r.endpoint, http.NoBody)
	if err != nil {
		return nil, errors.Wrap(err, errors.CodeInternalServerError, errMessage)
	}
	req.Header.Set("Authorization", credential)

	res, err := r.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, errors.CodeInternalServerError, errMessage)
	}
	defer func() {
		// NOTE: errに直接詰めると関数内のエラーがnilで上書きされるためエラー発生時のみ上書きする.
		if e := res.Body.Close(); e != nil {
			err = errors.Wrap(e, errors.CodeInternalServerError, errMessage)
		}
	}()

	var model model.AccountModel
	if err := client.NewDecoder(res).Decode(&model); err != nil {
		var code errors.ErrorCode
		switch {
		case stderr.Is(err, client.ErrUnauthenticated):
			code = errors.CodeUnauthenticated
		case stderr.Is(err, client.ErrUnauthorized):
			code = errors.CodeUnauthorized
		default:
			code = errors.CodeInternalServerError
		}
		return nil, errors.Wrap(err, code, errMessage)
	}

	return transformer.ToAccountEntity(&model), nil
}
