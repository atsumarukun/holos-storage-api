//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=../../../../test/mock/$GOPACKAGE/$GOFILE
package usecase

import (
	"context"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/repository"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/usecase/dto"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/usecase/mapper"
)

type AuthorizationUsecase interface {
	Authorize(context.Context, string) (*dto.AccountDTO, error)
}

type authorizationUsecase struct {
	accountRepo repository.AccountRepository
}

func NewAuthorizationUsecase(accountRepo repository.AccountRepository) AuthorizationUsecase {
	return &authorizationUsecase{
		accountRepo: accountRepo,
	}
}

func (u *authorizationUsecase) Authorize(ctx context.Context, credential string) (*dto.AccountDTO, error) {
	account, err := u.accountRepo.FindOneByCredential(ctx, credential)
	if err != nil {
		return nil, err
	}

	return mapper.ToAccountDTO(account), nil
}
