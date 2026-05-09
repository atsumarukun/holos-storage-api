package service_test

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/service"
	"github.com/atsumarukun/holos-storage-api/test/assert"
)

func TestBody_BuildPath(t *testing.T) {
	accountID := uuid.New()
	volume := &entity.Volume{
		ID:        uuid.New(),
		AccountID: accountID,
		Name:      "name",
		IsPublic:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	entry := &entity.Entry{
		ID:        uuid.New(),
		AccountID: accountID,
		VolumeID:  volume.ID,
		Key:       "key/sample.txt",
		Size:      4,
		Type:      "text/plain; charset=utf-8",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name         string
		inputVolume  *entity.Volume
		inputEntry   *entity.Entry
		expectResult string
		expectError  error
	}{
		{name: "successfully built path with volume and entry", inputVolume: volume, inputEntry: entry, expectResult: "name/key/sample.txt", expectError: nil},
		{name: "successfully built path with volume only", inputVolume: volume, inputEntry: nil, expectResult: "name", expectError: nil},
		{name: "failed to build path because volume is nil", inputVolume: nil, inputEntry: entry, expectResult: "", expectError: service.ErrNilVolume},
		{name: "failed to build path because volume and entry are nil", inputVolume: nil, inputEntry: nil, expectResult: "", expectError: service.ErrNilVolume},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			serv := service.NewBodyService()
			result, err := serv.BuildPath(tt.inputVolume, tt.inputEntry)
			assert.Error(t, err, tt.expectError)

			if result != tt.expectResult {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectResult, result)
			}
		})
	}
}
