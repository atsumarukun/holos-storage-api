package parameter

import (
	"fmt"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/pkg/status"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/pkg/status/code"
	"github.com/gin-gonic/gin"
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
