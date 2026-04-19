package code

type StatusCode string

var (
	BadRequest           StatusCode = "BAD_REQUEST"
	Unauthorized         StatusCode = "UNAUTHORIZED"
	Forbidden            StatusCode = "FORBIDDEN"
	NotFound             StatusCode = "NOT_FOUND"
	Conflict             StatusCode = "CONFLICT"
	UnprocessableContent StatusCode = "UNPROCESSABLE_CONTENT"
	Internal             StatusCode = "INTERNAL"
)
