package file

import (
	"io"
	"path/filepath"

	"github.com/atsumarukun/holos-api-pkg/errors"
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
	const errMessage = "failed to create body"

	if err := r.fs.MkdirAll(r.basePath+filepath.Dir(path), 0o755); err != nil {
		return errors.Wrap(err, errors.CodeInternalServerError, errMessage)
	}

	if reader == nil {
		return r.fs.Mkdir(r.basePath+path, 0o755)
	} else {
		file, err := r.fs.Create(r.basePath + path)
		if err != nil {
			return errors.Wrap(err, errors.CodeInternalServerError, errMessage)
		}
		defer func() {
			// NOTE: errに直接詰めると関数内のエラーがnilで上書きされるためエラー発生時のみ上書きする.
			if e := file.Close(); e != nil {
				err = errors.Wrap(e, errors.CodeInternalServerError, errMessage)
			}
		}()

		if _, err := io.Copy(file, reader); err != nil {
			return errors.Wrap(err, errors.CodeInternalServerError, errMessage)
		}

		return nil
	}
}

func (r *bodyRepository) Update(src, dst string) error {
	const errMessage = "failed to update body"

	// NOTE: srcで指定されたpathが存在するか判定する.
	_, err := r.fs.Stat(r.basePath + src)
	if err != nil {
		return errors.Wrap(err, errors.CodeInternalServerError, errMessage)
	}

	if err := r.fs.MkdirAll(r.basePath+filepath.Dir(dst), 0o755); err != nil {
		return errors.Wrap(err, errors.CodeInternalServerError, errMessage)
	}

	if err := r.fs.Rename(r.basePath+src, r.basePath+dst); err != nil {
		return errors.Wrap(err, errors.CodeInternalServerError, errMessage)
	}

	return nil
}

func (r *bodyRepository) Delete(path string) error {
	const errMessage = "failed to delete body"

	if err := r.fs.RemoveAll(r.basePath + path); err != nil {
		return errors.Wrap(err, errors.CodeInternalServerError, errMessage)
	}

	return nil
}

func (r *bodyRepository) Copy(src, dst string) error {
	const errMessage = "failed to copy body"

	info, err := r.fs.Stat(r.basePath + src)
	if err != nil {
		return errors.Wrap(err, errors.CodeInternalServerError, errMessage)
	}

	if info.IsDir() {
		if err := r.fs.Mkdir(r.basePath+dst, 0o755); err != nil {
			return errors.Wrap(err, errors.CodeInternalServerError, errMessage)
		}
		entries, err := afero.ReadDir(r.fs, r.basePath+src)
		if err != nil {
			return errors.Wrap(err, errors.CodeInternalServerError, errMessage)
		}
		for _, entry := range entries {
			if err := r.Copy(src+"/"+entry.Name(), dst+"/"+entry.Name()); err != nil {
				return errors.Wrap(err, errors.CodeInternalServerError, errMessage)
			}
		}
	} else {
		if err := r.copyFile(src, dst, errMessage); err != nil {
			return err
		}
	}

	return nil
}

func (r *bodyRepository) FindOneByPath(path string) (io.ReadCloser, error) {
	const errMessage = "failed to find body by path"

	info, err := r.fs.Stat(r.basePath + path)
	if err != nil {
		return nil, errors.Wrap(err, errors.CodeInternalServerError, errMessage)
	}

	if info.IsDir() {
		return nil, nil
	}

	body, err := r.fs.Open(r.basePath + path)
	if err != nil {
		return nil, errors.Wrap(err, errors.CodeInternalServerError, errMessage)
	}

	return body, nil
}

func (r *bodyRepository) copyFile(src, dst, errMessage string) (err error) {
	in, err := r.fs.Open(r.basePath + src)
	if err != nil {
		return errors.Wrap(err, errors.CodeInternalServerError, errMessage)
	}
	defer func() {
		if e := in.Close(); e != nil {
			err = errors.Wrap(e, errors.CodeInternalServerError, errMessage)
		}
	}()

	out, err := r.fs.Create(r.basePath + dst)
	if err != nil {
		return errors.Wrap(err, errors.CodeInternalServerError, errMessage)
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = errors.Wrap(e, errors.CodeInternalServerError, errMessage)
		}
	}()

	if _, err = io.Copy(out, in); err != nil {
		return errors.Wrap(err, errors.CodeInternalServerError, errMessage)
	}

	return nil
}
