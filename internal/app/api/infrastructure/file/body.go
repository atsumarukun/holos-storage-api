package file

import (
	"io"
	"path/filepath"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/repository"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/pkg/status"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/pkg/status/code"
	"github.com/spf13/afero"
)

var ErrRequiredReader = status.Error(code.Internal, "reader is required")

type bodyRepository struct {
	fs       afero.Fs
	basePath string
}

func NewBodyRepository(fs afero.Fs, basePath string) repository.BodyRepository {
	return bodyRepository{
		fs:       fs,
		basePath: basePath,
	}
}

func (r bodyRepository) Create(path string, reader io.Reader) (err error) {
	if reader == nil {
		return ErrRequiredReader
	}

	if err := r.fs.MkdirAll(r.basePath+filepath.Dir(path), 0755); err != nil {
		return err
	}

	file, err := r.fs.Create(r.basePath + path)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			err = closeErr
		}
	}()

	_, err = io.Copy(file, reader)
	return err
}
