package file_test

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/infrastructure/file"
	"github.com/spf13/afero"
)

type errReader struct{}

func (e *errReader) Read([]byte) (int, error) {
	return 0, io.ErrNoProgress
}

func TestBody_Create(t *testing.T) {
	tests := []struct {
		name         string
		inputPath    string
		inputReader  io.Reader
		expectResult bool
		expectError  error
	}{
		{name: "success", inputPath: "sample/test.txt", inputReader: bytes.NewBufferString("test"), expectResult: true, expectError: nil},
		{name: "reader is nil", inputPath: "sample/", inputReader: nil, expectResult: true, expectError: nil},
		{name: "create error", inputPath: "sample/test.txt", inputReader: &errReader{}, expectResult: true, expectError: io.ErrNoProgress},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			basePath := ""

			repo := file.NewBodyRepository(fs, basePath)
			if err := repo.Create(tt.inputPath, tt.inputReader); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			exists, err := afero.Exists(fs, basePath+tt.inputPath)
			if err != nil {
				t.Error(err)
			}

			if exists != tt.expectResult {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectResult, exists)
			}
		})
	}
}
