//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=../../../../../test/mock/domain/$GOPACKAGE/$GOFILE
package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/entity"
)

type EntryRepository interface {
	Create(context.Context, *entity.Entry) error
	Update(context.Context, *entity.Entry) error
	FindOneByKeyAndVolumeID(context.Context, string, uuid.UUID) (*entity.Entry, error)
}
