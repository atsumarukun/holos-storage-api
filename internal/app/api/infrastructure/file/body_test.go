package file_test

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/spf13/afero"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/infrastructure/file"
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
		{name: "reader is nil", inputPath: "sample/test", inputReader: nil, expectResult: true, expectError: nil},
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

func TestBody_Update(t *testing.T) {
	tests := []struct {
		name         string
		inputSrc     string
		inputDst     string
		expectResult bool
		expectError  error
		setMockFS    func(fs afero.Fs)
	}{
		{
			name:         "success",
			inputSrc:     "sample/test.txt",
			inputDst:     "update/test.txt",
			expectResult: true,
			expectError:  nil,
			setMockFS: func(fs afero.Fs) {
				if err := afero.WriteFile(fs, "sample/test.txt", []byte("test"), 0o755); err != nil {
					t.Error(err)
				}
			},
		},
		{
			name:         "update error",
			inputSrc:     "sample/test.txt",
			inputDst:     "update/test.txt",
			expectResult: false,
			expectError:  afero.ErrFileNotFound,
			setMockFS: func(fs afero.Fs) {
				if err := afero.WriteFile(fs, "test.txt", []byte("test"), 0o755); err != nil {
					t.Error(err)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			basePath := ""

			tt.setMockFS(fs)

			repo := file.NewBodyRepository(fs, basePath)
			if err := repo.Update(tt.inputSrc, tt.inputDst); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			exists, err := afero.Exists(fs, basePath+tt.inputDst)
			if err != nil {
				t.Error(err)
			}
			if exists != tt.expectResult {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectResult, exists)
			}

			if !errors.Is(tt.expectError, afero.ErrFileNotFound) {
				exists, err = afero.Exists(fs, basePath+tt.inputSrc)
				if err != nil {
					t.Error(err)
				}
				if exists != !tt.expectResult {
					t.Errorf("\nexpect: %v\ngot: %v", !tt.expectResult, exists)
				}
			}
		})
	}
}
