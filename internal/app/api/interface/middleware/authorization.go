package middleware

import (
	"strings"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/interface/pkg/errors"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/pkg/status"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/pkg/status/code"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/usecase"
	"github.com/gin-gonic/gin"
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
	if len(strings.Split(credential, " ")) != 2 {
		err := status.Error(code.Unauthorized, "unauthorized")
		errors.Handle(c, err)
		c.Abort()
		return
	}

	ctx := c.Request.Context()

	account, err := m.authorizationUC.Authorize(ctx, credential)
	if err != nil {
		errors.Handle(c, err)
		c.Abort()
		return
	}

	c.Set("accountID", account.ID)
	c.Next()
}
