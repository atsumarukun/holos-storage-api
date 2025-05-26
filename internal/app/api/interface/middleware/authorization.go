package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/interface/pkg/errors"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/usecase"
)

type AuthorizationMiddleware interface {
	Authorize(*gin.Context)
}

type authorizationMiddleware struct {
	authorizationUC usecase.AuthorizationUsecase
}

func NewAuthorizationMiddleware(authorizationUC usecase.AuthorizationUsecase) AuthorizationMiddleware {
	return &authorizationMiddleware{
		authorizationUC: authorizationUC,
	}
}

func (m *authorizationMiddleware) Authorize(c *gin.Context) {
	credential := c.Request.Header.Get("Authorization")
	volumeName := c.Param("volumeName")
	key := c.Param("key")
	method := c.Request.Method

	ctx := c.Request.Context()

	account, err := m.authorizationUC.Authorize(ctx, credential, volumeName, key, method)
	if err != nil {
		errors.Handle(c, err)
		c.Abort()
		return
	}

	c.Set("accountID", account.ID)
	c.Next()
}
