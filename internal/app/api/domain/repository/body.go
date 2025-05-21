//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=../../../../../test/mock/domain/$GOPACKAGE/$GOFILE
package repository

import "io"

type BodyRepository interface {
	Create(string, io.Reader) error
	Update(string, string) error
	Delete(string) error
	FindOneByPath(string) (io.ReadCloser, error)
}
