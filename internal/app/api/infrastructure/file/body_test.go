package file_test

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/google/go-cmp/cmp"
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
		{name: "create file", inputPath: "key/sample.txt", inputReader: bytes.NewBufferString("test"), expectResult: true, expectError: nil},
		{name: "create folder", inputPath: "key", inputReader: nil, expectResult: true, expectError: nil},
		{name: "create error", inputPath: "key/sample.txt", inputReader: &errReader{}, expectResult: true, expectError: io.ErrNoProgress},
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
			name:         "successfully updated",
			inputSrc:     "key/sample.txt",
			inputDst:     "key/update.txt",
			expectResult: true,
			expectError:  nil,
			setMockFS: func(fs afero.Fs) {
				if err := afero.WriteFile(fs, "key/sample.txt", []byte("test"), 0o755); err != nil {
					t.Error(err)
				}
			},
		},
		{
			name:         "not found",
			inputSrc:     "key/sample.txt",
			inputDst:     "key/update.txt",
			expectResult: false,
			expectError:  afero.ErrFileNotFound,
			setMockFS:    func(afero.Fs) {},
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

func TestBody_Delete(t *testing.T) {
	tests := []struct {
		name         string
		inputPath    string
		expectResult bool
		expectError  error
		setMockFS    func(fs afero.Fs)
	}{
		{
			name:         "successfully deleted",
			inputPath:    "key/sample.txt",
			expectResult: false,
			expectError:  nil,
			setMockFS: func(fs afero.Fs) {
				if err := afero.WriteFile(fs, "key/sample.txt", []byte("test"), 0o755); err != nil {
					t.Error(err)
				}
			},
		},
		{
			name:         "not found",
			inputPath:    "key/sample.txt",
			expectResult: false,
			expectError:  nil,
			setMockFS:    func(afero.Fs) {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			basePath := ""

			tt.setMockFS(fs)

			repo := file.NewBodyRepository(fs, basePath)
			if err := repo.Delete(tt.inputPath); !errors.Is(err, tt.expectError) {
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

func TestBody_Copy(t *testing.T) {
	tests := []struct {
		name         string
		inputSrc     string
		inputDst     string
		expectResult bool
		expectError  error
		setMockFS    func(fs afero.Fs)
	}{
		{
			name:         "successfully copied",
			inputSrc:     "key/sample.txt",
			inputDst:     "key/sample copy.txt",
			expectResult: true,
			expectError:  nil,
			setMockFS: func(fs afero.Fs) {
				if err := afero.WriteFile(fs, "key/sample.txt", []byte("test"), 0o755); err != nil {
					t.Error(err)
				}
			},
		},
		{
			name:         "not found",
			inputSrc:     "key/sample.txt",
			inputDst:     "key/sample copy.txt",
			expectResult: false,
			expectError:  afero.ErrFileNotFound,
			setMockFS:    func(afero.Fs) {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			basePath := ""

			tt.setMockFS(fs)

			repo := file.NewBodyRepository(fs, basePath)
			if err := repo.Copy(tt.inputSrc, tt.inputDst); !errors.Is(err, tt.expectError) {
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
				if exists != tt.expectResult {
					t.Errorf("\nexpect: %v\ngot: %v", tt.expectResult, exists)
				}
			}
		})
	}
}

func TestBody_FindOneByPath(t *testing.T) {
	tests := []struct {
		name         string
		inputPath    string
		expectResult []byte
		expectError  error
		setMockFS    func(fs afero.Fs)
	}{
		{
			name:         "found",
			inputPath:    "key/sample.txt",
			expectResult: []byte("test"),
			expectError:  nil,
			setMockFS: func(fs afero.Fs) {
				if err := afero.WriteFile(fs, "key/sample.txt", []byte("test"), 0o755); err != nil {
					t.Error(err)
				}
			},
		},
		{
			name:         "not found",
			inputPath:    "key/sample.txt",
			expectResult: nil,
			expectError:  afero.ErrFileNotFound,
			setMockFS:    func(afero.Fs) {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			basePath := ""

			tt.setMockFS(fs)

			repo := file.NewBodyRepository(fs, basePath)
			body, err := repo.FindOneByPath(tt.inputPath)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			if tt.expectResult != nil {
				result, err := io.ReadAll(body)
				if err != nil {
					t.Error(err)
				}
				if diff := cmp.Diff(tt.expectResult, result); diff != "" {
					t.Error(diff)
				}
			}
		})
	}
}
