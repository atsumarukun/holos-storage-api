//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=../../../../../test/mock/domain/$GOPACKAGE/$GOFILE
package repository

import (
	"context"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/entity"
	"github.com/google/uuid"
)

type EntryRepository interface {
	Create(context.Context, *entity.Entry) error
	FindOneByKeyAndVolumeIDAndAccountID(context.Context, string, uuid.UUID, uuid.UUID) (*entity.Entry, error)
}
