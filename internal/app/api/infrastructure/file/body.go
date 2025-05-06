package file

import (
	"errors"
	"io"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/repository"
	"github.com/spf13/afero"
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

func (r bodyRepository) Create(path string, reader io.Reader) error {
	return errors.New("not implemented")
}
