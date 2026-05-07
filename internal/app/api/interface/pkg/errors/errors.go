package errors

import (
	stderr "errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/atsumarukun/holos-api-pkg/errors"
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Code    errors.ErrorCode `json:"code"`
	Message string           `json:"message"`
}

func Handle(c *gin.Context, err error) {
	if err == nil {
		return
	}

	slog.ErrorContext(c.Request.Context(), err.Error())

	status := http.StatusInternalServerError
	res := ErrorResponse{errors.CodeUnknown, "internal server error"}

	if v, ok := err.(interface{ Code() errors.ErrorCode }); ok {
		status = StatusCode[v.Code()]
		switch v.Code() {
		case errors.CodeUnknown:
			res = ErrorResponse{v.Code(), "internal server error"}
		case errors.CodeDuplicate, errors.CodeConstraintViolation, errors.CodeInvalidInput:
			res = ErrorResponse{v.Code(), stderr.Unwrap(err).Error()}
		default:
			res = ErrorResponse{v.Code(), strings.ToLower(strings.ReplaceAll(v.Code().String(), "_", " "))}
		}
	}

	c.JSON(status, map[string]ErrorResponse{"error": res})
}
