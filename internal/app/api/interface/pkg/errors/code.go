package errors

import (
	"net/http"

	"github.com/atsumarukun/holos-api-pkg/errors"
)

var StatusCode = map[errors.ErrorCode]int{
	errors.CodeBadRequest:          http.StatusBadRequest,
	errors.CodeUnauthenticated:     http.StatusUnauthorized,
	errors.CodeUnauthorized:        http.StatusForbidden,
	errors.CodeNotFound:            http.StatusNotFound,
	errors.CodeDuplicate:           http.StatusConflict,
	errors.CodeConstraintViolation: http.StatusConflict,
	errors.CodeInvalidInput:        http.StatusUnprocessableEntity,
	errors.CodeInternalServerError: http.StatusInternalServerError,
	errors.CodeUnknown:             http.StatusInternalServerError,
}
