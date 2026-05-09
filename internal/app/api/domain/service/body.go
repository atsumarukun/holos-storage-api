//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=../../../../../test/mock/domain/$GOPACKAGE/$GOFILE
package service

import (
	"github.com/atsumarukun/holos-api-pkg/errors"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/domain/entity"
)

type BodyService interface {
	BuildPath(*entity.Volume, *entity.Entry) (string, error)
}

type bodyService struct{}

func NewBodyService() BodyService {
	return &bodyService{}
}

func (s *bodyService) BuildPath(volume *entity.Volume, entry *entity.Entry) (string, error) {
	if volume == nil {
		return "", errors.Wrap(ErrNilVolume, errors.CodeInternalServerError, "failed to build body path")
	}

	if entry == nil {
		return volume.Name, nil
	}

	return volume.Name + "/" + entry.Key, nil
}
