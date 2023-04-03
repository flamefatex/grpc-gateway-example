package code

// Code 编码
type Code = string

const (
	OK              Code = "OK"
	AuthError       Code = "AuthError"
	InternalError   Code = "InternalError"
	InvalidArgument Code = "InvalidArgument"
)
