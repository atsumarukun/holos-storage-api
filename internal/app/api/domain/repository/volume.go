//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=../../../../../test/mock/domain/$GOPACKAGE/$GOFILE
package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/entity"
)

type VolumeRepository interface {
	Create(context.Context, *entity.Volume) error
	FindOneByNameAndAccountID(context.Context, string, uuid.UUID) (*entity.Volume, error)
}
