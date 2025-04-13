package errors

import (
	"log"
	"net/http"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/pkg/status"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/pkg/status/code"
	"github.com/gin-gonic/gin"
)

var codeMap = map[code.StatusCode]int{
	code.BadRequest:   http.StatusBadRequest,
	code.Unauthorized: http.StatusUnauthorized,
	code.Forbidden:    http.StatusForbidden,
	code.Conflict:     http.StatusConflict,
	code.Internal:     http.StatusInternalServerError,
}

var messageMap = map[code.StatusCode]string{
	code.BadRequest:   "bad request",
	code.Unauthorized: "unauthorized",
	code.Forbidden:    "forbidden",
	code.Conflict:     "conflict",
	code.Internal:     "internal server error",
}

func Handle(c *gin.Context, err error) {
	log.Println(err)

	if v, ok := err.(*status.Status); ok {
		if v.Code() == code.BadRequest {
			c.JSON(codeMap[v.Code()], v.Message())
			return
		}
		c.JSON(codeMap[v.Code()], map[string]string{"message": messageMap[v.Code()]})
		return
	}
	c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error"})
}
