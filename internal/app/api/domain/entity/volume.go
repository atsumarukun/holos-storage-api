package entity

import (
	"regexp"
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

	if err := volume.generateID(); err != nil {
		return nil, err
	}
	if err := volume.setAccountID(accountID); err != nil {
		return nil, err
	}
	if err := volume.SetName(name); err != nil {
		return nil, err
	}
	volume.SetIsPublic(isPublic)

	now := time.Now()
	volume.CreatedAt = now
	volume.UpdatedAt = now

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
	if len(name) < 1 {
		return ErrShortVolumeName
	} else if 255 < len(name) {
		return ErrLongVolumeName
	}
	matched, err := regexp.MatchString(`^[A-Za-z0-9!@#$%^&()_\-+=\[\]{};',.~ ]*$`, name)
	if err != nil {
		return err
	}
	if !matched {
		return ErrInvalidVolumeName
	}
	v.Name = name
	v.UpdatedAt = time.Now()
	return nil
}

func (v *Volume) SetIsPublic(isPublic bool) {
	v.IsPublic = isPublic
	v.UpdatedAt = time.Now()
}

func (v *Volume) generateID() error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	v.ID = id
	return nil
}

func (v *Volume) setAccountID(accountID uuid.UUID) error {
	if accountID == uuid.Nil {
		return ErrRequiredVolumeAccountID
	}
	v.AccountID = accountID
	return nil
}
