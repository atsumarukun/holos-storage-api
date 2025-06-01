package file_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/spf13/afero"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/infrastructure/file"
)

const basePath = "storage/"

func checkExists(fs afero.Fs, paths []string, expect bool) error {
	for _, path := range paths {
		exists, err := afero.Exists(fs, basePath+path)
		if err != nil {
			return err
		}
		if exists != expect {
			if expect {
				return fmt.Errorf("%s is not exists", path)
			}
			return fmt.Errorf("%s is exists", path)
		}
	}
	return nil
}

type errReader struct{}

func (e *errReader) Read([]byte) (int, error) {
	return 0, io.ErrNoProgress
}

func TestBody_Create(t *testing.T) {
	tests := []struct {
		name        string
		inputPath   string
		inputReader io.Reader
		expectPaths []string
		expectError error
	}{
		{name: "create file", inputPath: "key/sample.txt", inputReader: bytes.NewBufferString("test"), expectPaths: []string{"key", "key/sample.txt"}, expectError: nil},
		{name: "create folder", inputPath: "key", inputReader: nil, expectPaths: []string{"key"}, expectError: nil},
		{name: "create error", inputPath: "key/sample.txt", inputReader: &errReader{}, expectPaths: []string{}, expectError: io.ErrNoProgress},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()

			repo := file.NewBodyRepository(fs, basePath)
			if err := repo.Create(tt.inputPath, tt.inputReader); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			if err := checkExists(fs, tt.expectPaths, true); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestBody_Update(t *testing.T) {
	tests := []struct {
		name          string
		inputSrc      string
		inputDst      string
		expectPaths   []string
		unexpectPaths []string
		expectError   error
		setMockFS     func(fs afero.Fs)
	}{
		{
			name:          "update file",
			inputSrc:      "key/sample.txt",
			inputDst:      "key/update.txt",
			expectPaths:   []string{"key", "key/update.txt"},
			unexpectPaths: []string{"key/sample.txt"},
			expectError:   nil,
			setMockFS: func(fs afero.Fs) {
				if err := afero.WriteFile(fs, basePath+"key/sample.txt", []byte("test"), 0o755); err != nil {
					t.Error(err)
				}
			},
		},
		{
			name:          "update folder",
			inputSrc:      "key/sample.txt",
			inputDst:      "update/sample.txt",
			expectPaths:   []string{"key", "update", "update/sample.txt"},
			unexpectPaths: []string{"key/sample.txt"},
			expectError:   nil,
			setMockFS: func(fs afero.Fs) {
				if err := afero.WriteFile(fs, basePath+"key/sample.txt", []byte("test"), 0o755); err != nil {
					t.Error(err)
				}
			},
		},
		{
			name:          "nonexistent folder",
			inputSrc:      "key/sample.txt",
			inputDst:      "sample/update.txt",
			expectPaths:   []string{"key", "sample", "sample/update.txt"},
			unexpectPaths: []string{"key/sample.txt"},
			expectError:   nil,
			setMockFS: func(fs afero.Fs) {
				if err := afero.WriteFile(fs, basePath+"key/sample.txt", []byte("test"), 0o755); err != nil {
					t.Error(err)
				}
			},
		},
		{
			name:          "not found",
			inputSrc:      "key/sample.txt",
			inputDst:      "key/update.txt",
			expectPaths:   []string{},
			unexpectPaths: []string{"key", "key/sample.txt", "sample/update.txt"},
			expectError:   afero.ErrFileNotFound,
			setMockFS:     func(afero.Fs) {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()

			tt.setMockFS(fs)

			repo := file.NewBodyRepository(fs, basePath)
			if err := repo.Update(tt.inputSrc, tt.inputDst); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			if err := checkExists(fs, tt.expectPaths, true); err != nil {
				t.Error(err)
			}
			if err := checkExists(fs, tt.unexpectPaths, false); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestBody_Delete(t *testing.T) {
	tests := []struct {
		name          string
		inputPath     string
		expectPaths   []string
		unexpectPaths []string
		expectError   error
		setMockFS     func(fs afero.Fs)
	}{
		{
			name:          "delete file",
			inputPath:     "key/sample.txt",
			expectPaths:   []string{"key"},
			unexpectPaths: []string{"key/sample.txt"},
			expectError:   nil,
			setMockFS: func(fs afero.Fs) {
				if err := afero.WriteFile(fs, basePath+"key/sample.txt", []byte("test"), 0o755); err != nil {
					t.Error(err)
				}
			},
		},
		{
			name:          "delete folder",
			inputPath:     "key",
			expectPaths:   []string{},
			unexpectPaths: []string{"key", "key/sample.txt"},
			expectError:   nil,
			setMockFS: func(fs afero.Fs) {
				if err := afero.WriteFile(fs, basePath+"key/sample.txt", []byte("test"), 0o755); err != nil {
					t.Error(err)
				}
			},
		},
		{
			name:          "not found",
			inputPath:     "key/sample.txt",
			expectPaths:   []string{},
			unexpectPaths: []string{"key/sample.txt"},
			expectError:   nil,
			setMockFS:     func(afero.Fs) {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()

			tt.setMockFS(fs)

			repo := file.NewBodyRepository(fs, basePath)
			if err := repo.Delete(tt.inputPath); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			if err := checkExists(fs, tt.expectPaths, true); err != nil {
				t.Error(err)
			}
			if err := checkExists(fs, tt.unexpectPaths, false); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestBody_Copy(t *testing.T) {
	tests := []struct {
		name          string
		inputSrc      string
		inputDst      string
		expectPaths   []string
		unexpectPaths []string
		expectError   error
		setMockFS     func(fs afero.Fs)
	}{
		{
			name:          "copy file",
			inputSrc:      "key/sample.txt",
			inputDst:      "key/sample copy.txt",
			expectPaths:   []string{"key", "key/sample.txt", "key/sample copy.txt"},
			unexpectPaths: []string{},
			expectError:   nil,
			setMockFS: func(fs afero.Fs) {
				if err := afero.WriteFile(fs, basePath+"key/sample.txt", []byte("test"), 0o755); err != nil {
					t.Error(err)
				}
			},
		},
		{
			name:          "copy folder",
			inputSrc:      "key",
			inputDst:      "key copy",
			expectPaths:   []string{"key", "key/sample.txt", "key copy", "key copy/sample.txt"},
			unexpectPaths: []string{},
			expectError:   nil,
			setMockFS: func(fs afero.Fs) {
				if err := afero.WriteFile(fs, basePath+"key/sample.txt", []byte("test"), 0o755); err != nil {
					t.Error(err)
				}
			},
		},
		{
			name:          "not found",
			inputSrc:      "key/sample.txt",
			inputDst:      "key/sample copy.txt",
			expectPaths:   []string{},
			unexpectPaths: []string{"key", "key/sample.txt", "key/sample copy.txt"},
			expectError:   afero.ErrFileNotFound,
			setMockFS:     func(afero.Fs) {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()

			tt.setMockFS(fs)

			repo := file.NewBodyRepository(fs, basePath)
			if err := repo.Copy(tt.inputSrc, tt.inputDst); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			if err := checkExists(fs, tt.expectPaths, true); err != nil {
				t.Error(err)
			}
			if err := checkExists(fs, tt.unexpectPaths, false); err != nil {
				t.Error(err)
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
			name:         "find file",
			inputPath:    "key/sample.txt",
			expectResult: []byte("test"),
			expectError:  nil,
			setMockFS: func(fs afero.Fs) {
				if err := afero.WriteFile(fs, basePath+"key/sample.txt", []byte("test"), 0o755); err != nil {
					t.Error(err)
				}
			},
		},
		{
			name:         "find folder",
			inputPath:    "key",
			expectResult: nil,
			expectError:  nil,
			setMockFS: func(fs afero.Fs) {
				if err := afero.WriteFile(fs, basePath+"key/sample.txt", []byte("test"), 0o755); err != nil {
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
