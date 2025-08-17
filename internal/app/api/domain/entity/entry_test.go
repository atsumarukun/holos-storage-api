package entity_test

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/entity"
)

func assertEntry(t *testing.T, e *entity.Entry) {
	if e.ID == uuid.Nil {
		t.Error("id is not set")
	}
	if e.AccountID == uuid.Nil {
		t.Error("account_id is not set")
	}
	if e.VolumeID == uuid.Nil {
		t.Error("volume_id is not set")
	}
	if e.Key == "" {
		t.Error("key is not set")
	}
	if e.Type == "" {
		t.Error("type is not set")
	}
	if e.CreatedAt.IsZero() {
		t.Error("created_at is not set")
	}
	if e.UpdatedAt.IsZero() {
		t.Error("updated_at is not set")
	}
	if !e.CreatedAt.Equal(e.UpdatedAt) {
		t.Error("expect created_at and updated_at to be equal")
	}
}

func TestNewEntry(t *testing.T) {
	tests := []struct {
		name           string
		inputAccountID uuid.UUID
		inputVolumeID  uuid.UUID
		inputKey       string
		inputSize      uint64
		inputType      string
		expectError    error
	}{
		{name: "successfully initialized", inputAccountID: uuid.New(), inputVolumeID: uuid.New(), inputKey: "key", inputSize: 1000, inputType: "folder", expectError: nil},
		{name: "account id is nil", inputAccountID: uuid.Nil, inputVolumeID: uuid.New(), inputKey: "key", inputSize: 1000, inputType: "folder", expectError: entity.ErrRequiredEntryAccountID},
		{name: "volume id is nil", inputAccountID: uuid.New(), inputVolumeID: uuid.Nil, inputKey: "key", inputSize: 1000, inputType: "folder", expectError: entity.ErrRequiredEntryVolumeID},
		{name: "invalid key", inputAccountID: uuid.New(), inputVolumeID: uuid.New(), inputKey: "", inputSize: 1000, inputType: "folder", expectError: entity.ErrShortEntryKey},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry, err := entity.NewEntry(tt.inputAccountID, tt.inputVolumeID, tt.inputKey, tt.inputSize, tt.inputType)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			if tt.expectError == nil {
				if entry == nil {
					t.Error("entry is nil")
				} else {
					assertEntry(t, entry)
				}
			}
		})
	}
}

func TestEntry_SetKey(t *testing.T) {
	entry := &entity.Entry{
		ID:        uuid.New(),
		AccountID: uuid.New(),
		VolumeID:  uuid.New(),
		Key:       "test/sample.jpg",
		Size:      10000,
		Type:      "image/jpeg",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name        string
		inputKey    string
		expectError error
	}{
		{name: "mixed lower case and upper case and number", inputKey: "entryKey1234", expectError: nil},
		{name: "valid symbols", inputKey: "!@#$%^&()_-+=[]{};',./~", expectError: nil},
		{name: "include space", inputKey: "entry key", expectError: nil},
		{name: "include backslash", inputKey: "entry\\key", expectError: entity.ErrInvalidEntryKey},
		{name: "include colon", inputKey: "entry:key", expectError: entity.ErrInvalidEntryKey},
		{name: "include asterisk", inputKey: "entry*key", expectError: entity.ErrInvalidEntryKey},
		{name: "include question mark", inputKey: "entry?key", expectError: entity.ErrInvalidEntryKey},
		{name: "include double quotation marks", inputKey: "entry\"key", expectError: entity.ErrInvalidEntryKey},
		{name: "include greater than sign", inputKey: "entry>key", expectError: entity.ErrInvalidEntryKey},
		{name: "include less than sign", inputKey: "entry<key", expectError: entity.ErrInvalidEntryKey},
		{name: "include vertical bar", inputKey: "entry|key", expectError: entity.ErrInvalidEntryKey},
		{name: "full width", inputKey: "エントリーキー", expectError: entity.ErrInvalidEntryKey},
		{name: "0 characters", inputKey: strings.Repeat("a", 0), expectError: entity.ErrShortEntryKey},
		{name: "1 characters", inputKey: strings.Repeat("a", 1), expectError: nil},
		{name: "512 characters", inputKey: strings.Repeat("a/", 255) + "aa", expectError: nil},
		{name: "513 characters", inputKey: strings.Repeat("a", 513), expectError: entity.ErrLongEntryKey},
		{name: "255 characters per element", inputKey: strings.Repeat("a", 255), expectError: nil},
		{name: "255 characters per element", inputKey: strings.Repeat("a", 256), expectError: entity.ErrInvalidEntryKey},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := entry.SetKey(tt.inputKey); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
		})
	}
}
