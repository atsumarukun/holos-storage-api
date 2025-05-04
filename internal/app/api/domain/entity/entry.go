package entity

import (
	"errors"
	"time"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/pkg/status"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/pkg/status/code"
	"github.com/google/uuid"
)

var (
	ErrRequiredentryAccountID = status.Error(code.Internal, "account id for entry is required")
	ErrRequiredentryVolumeID  = status.Error(code.Internal, "volume id for entry is required")
	ErrShortEntryKey          = status.Error(code.BadRequest, "entry key is too short")
	ErrLongEntryKey           = status.Error(code.BadRequest, "entry key is too long")
	ErrInvalidEntryKey        = status.Error(code.BadRequest, "entry key contains invalid characters")
)

type Entry struct {
	ID        uuid.UUID
	AccountID uuid.UUID
	VolumeID  uuid.UUID
	Key       string
	Size      uint64
	Type      string
	IsPublic  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewEntry(accountID, volumeID uuid.UUID, key string, size uint64, entryType string, isPublic bool) (*Entry, error) {
	return nil, errors.New("not implemented")
}

func RestoreEntry(id, accountID, volumeID uuid.UUID, key string, size uint64, entryType string, isPublic bool, createdAt, updatedAt time.Time) *Entry {
	return &Entry{
		ID:        id,
		AccountID: accountID,
		VolumeID:  volumeID,
		Key:       key,
		Size:      size,
		Type:      entryType,
		IsPublic:  isPublic,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func (e *Entry) SetKey(key string) error {
	return errors.New("not implemented")
}

func (e *Entry) SetIsPublic(isPublic bool) {}

func (e *Entry) setAccountID(accountID uuid.UUID) error {
	return errors.New("not implemented")
}

func (e *Entry) setVolumeID(volumeID uuid.UUID) error {
	return errors.New("not implemented")
}
