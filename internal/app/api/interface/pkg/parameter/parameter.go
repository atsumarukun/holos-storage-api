package parameter

import (
	stderr "errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	ErrParameterMissing     = stderr.New("parameter is missing")
	ErrInvalidParameterType = stderr.New("invalid parameter type")
)

func GetContextParameter[T any](c *gin.Context, name string) (T, error) {
	var zero T

	param, exists := c.Get(name)
	if !exists {
		return zero, ErrParameterMissing
	}

	v, ok := param.(T)
	if !ok {
		return zero, ErrInvalidParameterType
	}

	return v, nil
}

func GetPathParameter[T any](c *gin.Context, name string) (T, error) {
	var zero T
	param := c.Param(name)

	switch any(zero).(type) {
	case uuid.UUID:
		v, err := uuid.Parse(param)
		if err != nil {
			return zero, err
		}
		result, ok := any(v).(T)
		if !ok {
			return zero, ErrInvalidParameterType
		}
		return result, nil
	default:
		return zero, ErrInvalidParameterType
	}
}
