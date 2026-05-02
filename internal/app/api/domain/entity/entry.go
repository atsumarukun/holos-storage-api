package entity

import (
	stderr "errors"
	"regexp"
	"strings"
	"time"

	"github.com/atsumarukun/holos-api-pkg/errors"
	"github.com/google/uuid"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/pkg/status"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/pkg/status/code"
)

var (
	ErrEntryKeyInvalidLength        = stderr.New("entry key must be between 1 and 512 characters")
	ErrEntryKeyElementInvalidLength = stderr.New("entry key element must be between 1 and 255 characters")
	ErrEntryKeyInvalidChars         = stderr.New("entry key contains invalid characters")
	ErrEntryNilAccountID            = stderr.New("account id must not be nil")
	ErrEntryNilVolumeID             = stderr.New("volume id must not be nil")

	ErrRequiredEntryAccountID = status.Error(code.Internal, "account id for entry is required")
	ErrRequiredEntryVolumeID  = status.Error(code.Internal, "volume id for entry is required")
	ErrShortEntryKey          = status.Error(code.UnprocessableContent, "entry key is too short")
	ErrLongEntryKey           = status.Error(code.UnprocessableContent, "entry key is too long")
	ErrInvalidEntryKey        = status.Error(code.UnprocessableContent, "entry key contains invalid characters")
)

type Entry struct {
	ID        uuid.UUID
	AccountID uuid.UUID
	VolumeID  uuid.UUID
	Key       string
	Size      uint64
	Type      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewEntry(accountID, volumeID uuid.UUID, key string, size uint64, entryType string) (*Entry, error) {
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

	now := time.Now()
	entry.CreatedAt = now
	entry.UpdatedAt = now

	return &entry, nil
}

func RestoreEntry(id, accountID, volumeID uuid.UUID, key string, size uint64, entryType string, createdAt, updatedAt time.Time) *Entry {
	return &Entry{
		ID:        id,
		AccountID: accountID,
		VolumeID:  volumeID,
		Key:       key,
		Size:      size,
		Type:      entryType,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func (e *Entry) SetKey(key string) error {
	const errMessage = "failed to set entry key"

	key = strings.Trim(key, "/")
	if len(key) < 1 || 512 < len(key) {
		return errors.Wrap(ErrEntryKeyInvalidLength, errors.CodeInvalidInput, errMessage)
	}

	for k := range strings.SplitSeq(key, "/") {
		if len(k) < 1 || 255 < len(k) {
			return errors.Wrap(ErrEntryKeyElementInvalidLength, errors.CodeInvalidInput, errMessage)
		}
	}

	matched, err := regexp.MatchString(`^[A-Za-z0-9!@#$%^&()_\-+=\[\]{};',./~ ]*$`, key)
	if err != nil {
		return errors.Wrap(err, errors.CodeInternalServerError, errMessage)
	} else if !matched {
		return errors.Wrap(ErrEntryKeyInvalidChars, errors.CodeInvalidInput, errMessage)
	}

	e.Key = key
	e.UpdatedAt = time.Now()
	return nil
}

func (e *Entry) IsFolder() bool {
	return e.Type == "folder"
}

func (e *Entry) generateID() error {
	id, err := uuid.NewRandom()
	if err != nil {
		return errors.Wrap(err, errors.CodeInternalServerError, "failed to generate entry id")
	}
	e.ID = id
	return nil
}

func (e *Entry) setAccountID(accountID uuid.UUID) error {
	if accountID == uuid.Nil {
		return errors.Wrap(ErrEntryNilAccountID, errors.CodeInternalServerError, "failed to set account id on entry")
	}
	e.AccountID = accountID
	return nil
}

func (e *Entry) setVolumeID(volumeID uuid.UUID) error {
	if volumeID == uuid.Nil {
		return errors.Wrap(ErrEntryNilVolumeID, errors.CodeInternalServerError, "failed to set volume id on entry")
	}
	e.VolumeID = volumeID
	return nil
}
