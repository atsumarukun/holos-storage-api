package errors

import (
	"log"
	"net/http"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/pkg/status"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/pkg/status/code"
	"github.com/gin-gonic/gin"
)

type response struct {
	code    int
	message string
}

var responseMap = map[code.StatusCode]response{
	code.BadRequest:   {code: http.StatusBadRequest, message: "bad request"},
	code.Unauthorized: {code: http.StatusUnauthorized, message: "unauthorized"},
	code.Forbidden:    {code: http.StatusForbidden, message: "forbidden"},
	code.Conflict:     {code: http.StatusConflict, message: "conflict"},
	code.Internal:     {code: http.StatusInternalServerError, message: "internal server error"},
}

func Handle(c *gin.Context, err error) {
	log.Println(err)

	if v, ok := err.(*status.Status); ok {
		if v.Code() == code.BadRequest {
			c.JSON(http.StatusBadRequest, v.Message())
			return
		}
		resp := responseMap[v.Code()]
		c.JSON(resp.code, map[string]string{"message": resp.message})
		return
	}
	c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error"})
}
