package transformer

import (
	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/infrastructure/api/model"
)

func ToAccountEntity(account *model.AccountModel) *entity.Account {
	return entity.RestoreAccount(account.ID)
}
