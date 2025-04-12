package middleware

import (
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

func (u *authorizationMiddleware) Authorize(c *gin.Context) {}
