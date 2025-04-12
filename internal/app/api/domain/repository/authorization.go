//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=../../../../../test/mock/domain/$GOPACKAGE/$GOFILE
package repository

import (
	"context"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/entity"
)

type AuthorizationRepository interface {
	Authorize(context.Context, string) (*entity.Authorization, error)
}
