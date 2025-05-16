//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=../../../../../test/mock/domain/$GOPACKAGE/$GOFILE
package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/entity"
)

type VolumeRepository interface {
	Create(context.Context, *entity.Volume) error
	Update(context.Context, *entity.Volume) error
	Delete(context.Context, *entity.Volume) error
	FindOneByName(context.Context, string) (*entity.Volume, error)
	FindOneByNameAndAccountID(context.Context, string, uuid.UUID) (*entity.Volume, error)
	FindOneByIDAndAccountID(context.Context, uuid.UUID, uuid.UUID) (*entity.Volume, error)
	FindByAccountID(context.Context, uuid.UUID) ([]*entity.Volume, error)
}
