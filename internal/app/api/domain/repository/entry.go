//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=../../../../../test/mock/domain/$GOPACKAGE/$GOFILE
package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/pkg/status"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/pkg/status/code"
)

var ErrEntryNotFound = status.Error(code.NotFound, "entry not found")

type EntryRepository interface {
	Create(context.Context, *entity.Entry) error
	Update(context.Context, *entity.Entry) error
	Delete(context.Context, *entity.Entry) error
	FindOneByKeyAndVolumeID(context.Context, string, uuid.UUID) (*entity.Entry, error)
	FindOneByKeyAndVolumeIDAndAccountID(context.Context, string, uuid.UUID, uuid.UUID) (*entity.Entry, error)
	FindByVolumeIDAndAccountID(context.Context, uuid.UUID, uuid.UUID, *string, *uint64) ([]*entity.Entry, error)
}
