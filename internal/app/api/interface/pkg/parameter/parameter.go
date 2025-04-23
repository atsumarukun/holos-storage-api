package parameter

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/pkg/status"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/pkg/status/code"
)

func GetContextParameter[T any](c *gin.Context, name string) (T, error) {
	var zero T

	param, exists := c.Get(name)
	if !exists {
		return zero, status.Error(code.Internal, fmt.Sprintf("context does not have %s", name))
	}

	v, ok := param.(T)
	if !ok {
		return zero, status.Error(code.Internal, "invalid context parameter type")
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
			return zero, status.Error(code.BadRequest, err.Error())
		}
		return any(v).(T), nil
	default:
		return zero, status.Error(code.Internal, "invalid path parameter type")
	}
}
