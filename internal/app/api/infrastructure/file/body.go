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
	return &bodyRepository{
		fs:       fs,
		basePath: basePath,
	}
}

func (r *bodyRepository) Create(path string, reader io.Reader) (err error) {
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

func (r *bodyRepository) Update(src, dst string) error {
	if err := r.fs.MkdirAll(r.basePath+filepath.Dir(dst), 0o755); err != nil {
		return err
	}

	return r.fs.Rename(r.basePath+src, r.basePath+dst)
}

func (r *bodyRepository) Delete(path string) error {
	return r.fs.RemoveAll(r.basePath + path)
}

func (r *bodyRepository) Copy(src, dst string) error {
	info, err := r.fs.Stat(r.basePath + src)
	if err != nil {
		return err
	}

	if info.IsDir() {
		entries, err := afero.ReadDir(r.fs, r.basePath+src)
		if err != nil {
			return err
		}
		for _, entry := range entries {
			if entry.IsDir() {
				if err := r.Copy(src+"/"+entry.Name(), dst+"/"+entry.Name()); err != nil {
					return err
				}
			}
		}
	} else {
		if err := r.copyFile(src, dst); err != nil {
			return err
		}
	}

	return nil
}

func (r *bodyRepository) FindOneByPath(path string) (io.ReadCloser, error) {
	info, err := r.fs.Stat(r.basePath + path)
	if err != nil {
		return nil, err
	}

	if info.IsDir() {
		return nil, nil
	}

	return r.fs.Open(r.basePath + path)
}

func (r *bodyRepository) copyFile(src, dst string) (err error) {
	in, err := r.fs.Open(r.basePath + src)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := in.Close(); closeErr != nil {
			err = closeErr
		}
	}()

	out, err := r.fs.Create(r.basePath + dst)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := out.Close(); closeErr != nil {
			err = closeErr
		}
	}()

	_, err = io.Copy(out, in)
	return err
}
