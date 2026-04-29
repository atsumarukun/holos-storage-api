package assert

import (
	stderr "errors"
	"testing"

	"github.com/atsumarukun/holos-api-pkg/errors"
)

func Error(t *testing.T, got, expect error) {
	t.Helper()

	if !stderr.Is(got, expect) {
		t.Errorf("\nexpect: %v\ngot: %v", expect, got)
	}

	if got != nil {
		if _, ok := got.(interface {
			Code() errors.ErrorCode
			Message() string
		}); !ok {
			t.Errorf("error is not wrapped")
		}
	}
}
