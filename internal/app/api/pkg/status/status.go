package status

import (
	"fmt"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/pkg/status/code"
)

type Status struct {
	code    code.StatusCode
	message string
}

func New(code code.StatusCode, message string) *Status {
	return &Status{
		code:    code,
		message: message,
	}
}

func Error(code code.StatusCode, message string) error {
	return &Status{
		code:    code,
		message: message,
	}
}

func FromError(err error) *Status {
	if err == nil {
		return nil
	}

	if v, ok := err.(*Status); ok {
		return v
	} else {
		return &Status{
			code:    code.Internal,
			message: err.Error(),
		}
	}
}

func (e *Status) Error() string {
	return fmt.Sprintf("code: %s, message: %s", e.code, e.message)
}

func (e *Status) Code() code.StatusCode {
	return e.code
}

func (e *Status) Message() string {
	return e.message
}
