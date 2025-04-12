package transformer

import (
	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/infrastructure/model"
)

func ToAuthorizationEntity(authorization *model.AuthorizationModel) *entity.Authorization {
	return entity.RestoreAuthorization(authorization.AccountID)
}
