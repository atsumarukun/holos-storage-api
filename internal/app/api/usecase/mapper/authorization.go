package mapper

import (
	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/usecase/dto"
)

func ToAuthorizationDTO(authorization *entity.Authorization) *dto.AuthorizationDTO {
	return &dto.AuthorizationDTO{
		AccountID: authorization.AccountID,
	}
}
