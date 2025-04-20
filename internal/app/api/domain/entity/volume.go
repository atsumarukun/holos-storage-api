package entity

import (
	"errors"
	"time"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/pkg/status"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/pkg/status/code"
	"github.com/google/uuid"
)

var (
	ErrRequiredVolumeAccountID = status.Error(code.Internal, "account id for volume is required")
	ErrShortVolumeName         = status.Error(code.BadRequest, "volume name is too short")
	ErrLongVolumeName          = status.Error(code.BadRequest, "volume name is too long")
	ErrInvalidVolumeName       = status.Error(code.BadRequest, "volume name contains invalid characters")
)

type Volume struct {
	ID        uuid.UUID
	AccountID uuid.UUID
	Name      string
	IsPublic  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewVolume(accountID uuid.UUID, name string, isPublic bool) (*Volume, error) {
	var volume Volume

	return &volume, nil
}

func RestoreVolume(id uuid.UUID, accountID uuid.UUID, name string, isPublic bool, createdAt time.Time, updatedAt time.Time) *Volume {
	return &Volume{
		ID:        id,
		AccountID: accountID,
		Name:      name,
		IsPublic:  isPublic,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func (v *Volume) SetName(name string) error {
	return errors.New("not implemented")
}

func (v *Volume) SetIsPublic(isPublic bool) {
	v.IsPublic = isPublic
}

func (v *Volume) generateID() error {
	return errors.New("not implemented")
}
