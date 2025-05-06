package entity

import (
	"regexp"
	"time"

	"github.com/google/uuid"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/pkg/status"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/pkg/status/code"
)

var (
	ErrRequiredEntryAccountID = status.Error(code.Internal, "account id for entry is required")
	ErrRequiredEntryVolumeID  = status.Error(code.Internal, "volume id for entry is required")
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
	entry := Entry{
		Size: size,
		Type: entryType,
	}

	if err := entry.generateID(); err != nil {
		return nil, err
	}
	if err := entry.setAccountID(accountID); err != nil {
		return nil, err
	}
	if err := entry.setVolumeID(volumeID); err != nil {
		return nil, err
	}
	if err := entry.SetKey(key); err != nil {
		return nil, err
	}
	entry.SetIsPublic(isPublic)

	now := time.Now()
	entry.CreatedAt = now
	entry.UpdatedAt = now

	return &entry, nil
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
	if e.IsFolder() && key[len(key)-1:] != "/" {
		key += "/"
	}

	if len(key) < 1 {
		return ErrShortEntryKey
	} else if 255 < len(key) {
		return ErrLongEntryKey
	}

	matched, err := regexp.MatchString(`^[A-Za-z0-9!@#$%^&()_\-+=\[\]{};',./~ ]*$`, key)
	if err != nil {
		return err
	}
	if !matched {
		return ErrInvalidEntryKey
	}

	e.Key = key
	e.UpdatedAt = time.Now()
	return nil
}

func (e *Entry) SetIsPublic(isPublic bool) {
	e.IsPublic = isPublic
	e.UpdatedAt = time.Now()
}

func (e *Entry) IsFolder() bool {
	return e.Type == "folder"
}

func (e *Entry) generateID() error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	e.ID = id
	return nil
}

func (e *Entry) setAccountID(accountID uuid.UUID) error {
	if accountID == uuid.Nil {
		return ErrRequiredEntryAccountID
	}
	e.AccountID = accountID
	return nil
}

func (e *Entry) setVolumeID(volumeID uuid.UUID) error {
	if volumeID == uuid.Nil {
		return ErrRequiredEntryVolumeID
	}
	e.VolumeID = volumeID
	return nil
}
