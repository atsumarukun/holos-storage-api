package client

import (
	"encoding/json"
	stderr "errors"
	"net/http"

	"github.com/atsumarukun/holos-api-pkg/errors"
)

var (
	ErrUnauthenticated     = stderr.New("unauthenticated")
	ErrUnauthorized        = stderr.New("unauthorized")
	ErrInternalServerError = stderr.New("internal server error")
)

type Decoder interface {
	Decode(any) error
}

type decoder struct {
	res *http.Response
}

func NewDecoder(res *http.Response) Decoder {
	return &decoder{res: res}
}

func (d *decoder) Decode(v any) error {
	if d.res.StatusCode < 300 {
		return json.NewDecoder(d.res.Body).Decode(v)
	}

	var errorResponse struct {
		Error struct {
			Code    errors.ErrorCode `json:"code"`
			Message string           `json:"message"`
		} `json:"error"`
	}
	if err := json.NewDecoder(d.res.Body).Decode(&errorResponse); err != nil {
		return err
	}

	switch errorResponse.Error.Code {
	case errors.CodeUnauthenticated:
		return ErrUnauthenticated
	case errors.CodeUnauthorized:
		return ErrUnauthorized
	default:
		return ErrInternalServerError
	}
}
