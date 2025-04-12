//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=../../../../test/mock/$GOPACKAGE/$GOFILE
package usecase

import (
	"context"
	"errors"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/repository"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/usecase/dto"
)

type AuthorizationUsecase interface {
	Authorize(context.Context, string) (*dto.AuthorizationDTO, error)
}

type authorizationUsecase struct {
	authorizationRepo repository.AuthorizationRepository
}

func NewAuthorizationUsecase(authorizationRepo repository.AuthorizationRepository) AuthorizationUsecase {
	return &authorizationUsecase{
		authorizationRepo: authorizationRepo,
	}
}

func (u *authorizationUsecase) Authorize(ctx context.Context, credential string) (*dto.AuthorizationDTO, error) {
	return nil, errors.New("not implemented")
}
