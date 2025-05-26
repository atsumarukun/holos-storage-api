//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=../../../../test/mock/$GOPACKAGE/$GOFILE
package usecase

import (
	"context"
	"errors"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/repository"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/pkg/status"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/pkg/status/code"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/usecase/dto"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/usecase/mapper"
)

var ErrForbidden = status.Error(code.Forbidden, "forbidden")

type AuthorizationUsecase interface {
	Authorize(context.Context, string, string, string, string) (*dto.AccountDTO, error)
}

type authorizationUsecase struct {
	accountRepo repository.AccountRepository
	volumeRepo  repository.VolumeRepository
}

func NewAuthorizationUsecase(accountRepo repository.AccountRepository, volumeRepo repository.VolumeRepository) AuthorizationUsecase {
	return &authorizationUsecase{
		accountRepo: accountRepo,
		volumeRepo:  volumeRepo,
	}
}

func (u *authorizationUsecase) Authorize(ctx context.Context, credential, volumeName, key, method string) (*dto.AccountDTO, error) {
	isGetEntry := volumeName != "" && key != "" && (method == "GET" || method == "HEAD")

	if isGetEntry {
		volume, err := u.volumeRepo.FindOneByName(ctx, volumeName)
		if err != nil {
			return nil, err
		}
		if volume.IsPublic {
			account := entity.NewAccount(volume.AccountID)
			return mapper.ToAccountDTO(account), nil
		}
	}

	account, err := u.accountRepo.FindOneByCredential(ctx, credential)
	if err != nil {
		if isGetEntry && credential == "" && errors.Is(err, repository.ErrUnauthorized) {
			return nil, ErrForbidden
		}
		return nil, err
	}

	return mapper.ToAccountDTO(account), nil
}
