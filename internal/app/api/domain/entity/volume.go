package entity

import (
	stderr "errors"
	"regexp"
	"time"

	"github.com/atsumarukun/holos-api-pkg/errors"
	"github.com/google/uuid"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/pkg/status"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/pkg/status/code"
)

var (
	ErrVolumeNameInvalidLength = stderr.New("volume name must be between 1 and 512 characters")
	ErrVolumeNameInvalidChars  = stderr.New("volume name contains invalid characters")
	ErrVolumeNilAccountID      = stderr.New("account id must not be nil")

	ErrRequiredVolumeAccountID = status.Error(code.Internal, "account id for volume is required")
	ErrShortVolumeName         = status.Error(code.UnprocessableContent, "volume name is too short")
	ErrLongVolumeName          = status.Error(code.UnprocessableContent, "volume name is too long")
	ErrInvalidVolumeName       = status.Error(code.UnprocessableContent, "volume name contains invalid characters")
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

func RestoreVolume(id, accountID uuid.UUID, name string, isPublic bool, createdAt, updatedAt time.Time) *Volume {
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
	const errMessage = "failed to set volume name"

	if len(name) < 1 || 255 < len(name) {
		return errors.Wrap(ErrVolumeNameInvalidLength, errors.CodeInvalidInput, errMessage)
	}
	matched, err := regexp.MatchString(`^[A-Za-z0-9!@#$%^&()_\-+=\[\]{};',.~ ]*$`, name)
	if err != nil {
		return errors.Wrap(err, errors.CodeInternalServerError, errMessage)
	} else if !matched {
		return errors.Wrap(ErrVolumeNameInvalidChars, errors.CodeInvalidInput, errMessage)
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
		return errors.Wrap(err, errors.CodeInternalServerError, "failed to generate volume id")
	}
	v.ID = id
	return nil
}

func (v *Volume) setAccountID(accountID uuid.UUID) error {
	if accountID == uuid.Nil {
		return errors.Wrap(ErrVolumeNilAccountID, errors.CodeInternalServerError, "failed to set account id on volume")
	}
	v.AccountID = accountID
	return nil
}
