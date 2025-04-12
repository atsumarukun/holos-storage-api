package code

type StatusCode string

var (
	BadRequest   StatusCode = "BAD_REQUEST"
	Unauthorized StatusCode = "UNAUTHORIZED"
	Forbidden    StatusCode = "FORBIDDEN"
	Conflict     StatusCode = "CONFLICT"
	Internal     StatusCode = "INTERNAL"
)
