package errors

import (
	"encoding/json"
	"net/http"

	"github.com/atsumarukun/holos-storage-api/internal/app/api/pkg/status"
	"github.com/atsumarukun/holos-storage-api/internal/app/api/pkg/status/code"
)

var codeMap = map[int]code.StatusCode{
	http.StatusBadRequest:          code.BadRequest,
	http.StatusUnauthorized:        code.Unauthorized,
	http.StatusForbidden:           code.Forbidden,
	http.StatusConflict:            code.Conflict,
	http.StatusInternalServerError: code.Internal,
}

type apiError struct {
	Message string `json:"message"`
}

func Decode(resp *http.Response) error {
	var apiErr apiError
	if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
		return err
	}
	return status.Error(codeMap[resp.StatusCode], apiErr.Message)
}
