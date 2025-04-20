package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
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
