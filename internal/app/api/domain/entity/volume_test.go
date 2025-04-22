package entity_test

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/entity"
)

func assertVolume(t *testing.T, v *entity.Volume) {
	if v.ID == uuid.Nil {
		t.Error("id is not set")
	}
	if v.AccountID == uuid.Nil {
		t.Error("account_id is not set")
	}
	if v.Name == "" {
		t.Error("name is not set")
	}
	if v.CreatedAt.IsZero() {
		t.Error("created_at is not set")
	}
	if v.UpdatedAt.IsZero() {
		t.Error("updated_at is not set")
	}
	if !v.CreatedAt.Equal(v.UpdatedAt) {
		t.Error("expect created_at and updated_at to be equal")
	}
}

func TestNewVolume(t *testing.T) {
	tests := []struct {
		name           string
		inputAccountID uuid.UUID
		inputName      string
		inputIsPublic  bool
		expectError    error
	}{
		{name: "success", inputAccountID: uuid.New(), inputName: "name", inputIsPublic: false, expectError: nil},
		{name: "account id is nil", inputAccountID: uuid.Nil, inputName: "name", inputIsPublic: false, expectError: entity.ErrRequiredVolumeAccountID},
		{name: "invalid name", inputAccountID: uuid.New(), inputName: "", inputIsPublic: false, expectError: entity.ErrShortVolumeName},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			volume, err := entity.NewVolume(tt.inputAccountID, tt.inputName, tt.inputIsPublic)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			if tt.expectError == nil {
				if volume == nil {
					t.Error("volume is nil")
				} else {
					assertVolume(t, volume)
				}
			}
		})
	}
}

func TestVolume_SetName(t *testing.T) {
	volume := &entity.Volume{
		ID:        uuid.New(),
		AccountID: uuid.New(),
		Name:      "name",
		IsPublic:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name        string
		inputName   string
		expectError error
	}{
		{name: "lower case only", inputName: "name", expectError: nil},
		{name: "upper case only", inputName: "NAME", expectError: nil},
		{name: "number only", inputName: "1234", expectError: nil},
		{name: "mixed lower case and upper case and number", inputName: "volumeName1234", expectError: nil},
		{name: "valid symbols", inputName: "!@#$%^&()_-+=[]{};',.~", expectError: nil},
		{name: "include space", inputName: "volume name", expectError: nil},
		{name: "include slash", inputName: "volume/name", expectError: entity.ErrInvalidVolumeName},
		{name: "include backslash", inputName: "volume\\name", expectError: entity.ErrInvalidVolumeName},
		{name: "include colon", inputName: "volume:name", expectError: entity.ErrInvalidVolumeName},
		{name: "include asterisk", inputName: "volume*name", expectError: entity.ErrInvalidVolumeName},
		{name: "include question mark", inputName: "volume?name", expectError: entity.ErrInvalidVolumeName},
		{name: "include double quotation marks", inputName: "volume\"name", expectError: entity.ErrInvalidVolumeName},
		{name: "include greater than sign", inputName: "volume>name", expectError: entity.ErrInvalidVolumeName},
		{name: "include less than sign", inputName: "volume<name", expectError: entity.ErrInvalidVolumeName},
		{name: "include vertical bar", inputName: "volume|name", expectError: entity.ErrInvalidVolumeName},
		{name: "hiragana", inputName: "なまえ", expectError: entity.ErrInvalidVolumeName},
		{name: "katakana", inputName: "ナマエ", expectError: entity.ErrInvalidVolumeName},
		{name: "kanji", inputName: "ナマエ", expectError: entity.ErrInvalidVolumeName},
		{name: "0 characters", inputName: strings.Repeat("a", 0), expectError: entity.ErrShortVolumeName},
		{name: "1 characters", inputName: strings.Repeat("a", 1), expectError: nil},
		{name: "255 characters", inputName: strings.Repeat("a", 255), expectError: nil},
		{name: "256 characters", inputName: strings.Repeat("a", 256), expectError: entity.ErrLongVolumeName},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := volume.SetName(tt.inputName); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
		})
	}
}
