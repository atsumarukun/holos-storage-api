package file

import (
	"io"
	"path/filepath"

	"github.com/spf13/afero"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/repository"
)

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
	if err := r.fs.MkdirAll(r.basePath+filepath.Dir(path), 0o755); err != nil {
		return err
	}

	if reader == nil {
		return r.fs.Mkdir(r.basePath+path, 0o755)
	} else {
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
}

func (r bodyRepository) Update(src string, dst string) error {
	return r.fs.Rename(src, dst)
}
