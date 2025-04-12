//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=../../../../test/mock/$GOPACKAGE/$GOFILE
package usecase

import (
	"context"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/repository"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/usecase/dto"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/usecase/mapper"
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
	authorization, err := u.authorizationRepo.Authorize(ctx, credential)
	if err != nil {
		return nil, err
	}

	return mapper.ToAuthorizationDTO(authorization), nil
}
