//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=../../../../../test/mock/domain/$GOPACKAGE/$GOFILE
package repository

import (
	"context"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/entity"
)

type AccountRepository interface {
	FindOneByCredential(context.Context, string) (*entity.Account, error)
}
