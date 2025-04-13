package mapper

import (
	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/usecase/dto"
)

func ToAccountDTO(account *entity.Account) *dto.AccountDTO {
	return &dto.AccountDTO{
		ID: account.ID,
	}
}
