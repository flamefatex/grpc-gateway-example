// nolint:gomnd
package errorx

import "fmt"

// ErrorBadRequest new BadRequest error that is mapped to a 400 response.
func ErrorBadRequest(reason, message string) *Error {
	return New(400, reason, message)
}

// ErrorfBadRequest New(code fmt.Sprintf(format, a...))
func ErrorfBadRequest(reason, format string, a ...interface{}) *Error {
	return New(400, reason, fmt.Sprintf(format, a...))
}

// IsBadRequest determines if err is an error which indicates a BadRequest error.
// It supports wrapped errors.
func IsBadRequest(err error) bool {
	return Code(err) == 400
}

// ErrorUnauthorized new Unauthorized error that is mapped to a 401 response.
func ErrorUnauthorized(reason, message string) *Error {
	return New(401, reason, message)
}

// ErrorfUnauthorized New(code fmt.Sprintf(format, a...))
func ErrorfUnauthorized(reason, format string, a ...interface{}) *Error {
	return New(401, reason, fmt.Sprintf(format, a...))
}

// IsUnauthorized determines if err is an error which indicates an Unauthorized error.
// It supports wrapped errors.
func IsUnauthorized(err error) bool {
	return Code(err) == 401
}

// ErrorForbidden new Forbidden error that is mapped to a 403 response.
func ErrorForbidden(reason, message string) *Error {
	return New(403, reason, message)
}

// ErrorfForbidden New(code fmt.Sprintf(format, a...))
func ErrorfForbidden(reason, format string, a ...interface{}) *Error {
	return New(403, reason, fmt.Sprintf(format, a...))
}

// IsForbidden determines if err is an error which indicates a Forbidden error.
// It supports wrapped errors.
func IsForbidden(err error) bool {
	return Code(err) == 403
}

// ErrorNotFound new NotFound error that is mapped to a 404 response.
func ErrorNotFound(reason, message string) *Error {
	return New(404, reason, message)
}

// ErrorfNotFound New(code fmt.Sprintf(format, a...))
func ErrorfNotFound(reason, format string, a ...interface{}) *Error {
	return New(404, reason, fmt.Sprintf(format, a...))
}

// IsNotFound determines if err is an error which indicates an NotFound error.
// It supports wrapped errors.
func IsNotFound(err error) bool {
	return Code(err) == 404
}

// ErrorConflict new Conflict error that is mapped to a 409 response.
func ErrorConflict(reason, message string) *Error {
	return New(409, reason, message)
}

// ErrorfConflict New(code fmt.Sprintf(format, a...))
func ErrorfConflict(reason, format string, a ...interface{}) *Error {
	return New(409, reason, fmt.Sprintf(format, a...))
}

// IsConflict determines if err is an error which indicates a Conflict error.
// It supports wrapped errors.
func IsConflict(err error) bool {
	return Code(err) == 409
}

// ErrorInternalServer new InternalServer error that is mapped to a 500 response.
func ErrorInternalServer(reason, message string) *Error {
	return New(500, reason, message)
}

// ErrorfInternalServer New(code fmt.Sprintf(format, a...))
func ErrorfInternalServer(reason, format string, a ...interface{}) *Error {
	return New(500, reason, fmt.Sprintf(format, a...))
}

// IsInternalServer determines if err is an error which indicates an Internal error.
// It supports wrapped errors.
func IsInternalServer(err error) bool {
	return Code(err) == 500
}

// ErrorServiceUnavailable new ServiceUnavailable error that is mapped to an HTTP 503 response.
func ErrorServiceUnavailable(reason, message string) *Error {
	return New(503, reason, message)
}

// ErrorfServiceUnavailable New(code fmt.Sprintf(format, a...))
func ErrorfServiceUnavailable(reason, format string, a ...interface{}) *Error {
	return New(503, reason, fmt.Sprintf(format, a...))
}

// IsServiceUnavailable determines if err is an error which indicates an Unavailable error.
// It supports wrapped errors.
func IsServiceUnavailable(err error) bool {
	return Code(err) == 503
}

// ErrorGatewayTimeout new GatewayTimeout error that is mapped to an HTTP 504 response.
func ErrorGatewayTimeout(reason, message string) *Error {
	return New(504, reason, message)
}

// ErrorfGatewayTimeout New(code fmt.Sprintf(format, a...))
func ErrorfGatewayTimeout(reason, format string, a ...interface{}) *Error {
	return New(504, reason, fmt.Sprintf(format, a...))
}

// IsGatewayTimeout determines if err is an error which indicates a GatewayTimeout error.
// It supports wrapped errors.
func IsGatewayTimeout(err error) bool {
	return Code(err) == 504
}

// ErrorClientClosed new ClientClosed error that is mapped to an HTTP 499 response.
func ErrorClientClosed(reason, message string) *Error {
	return New(499, reason, message)
}

// ErrorfClientClosed New(code fmt.Sprintf(format, a...))
func ErrorfClientClosed(reason, format string, a ...interface{}) *Error {
	return New(499, reason, fmt.Sprintf(format, a...))
}

// IsClientClosed determines if err is an error which indicates a IsClientClosed error.
// It supports wrapped errors.
func IsClientClosed(err error) bool {
	return Code(err) == 499
}
