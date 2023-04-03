package statusx

import (
	"errors"
	"fmt"

	"google.golang.org/grpc/codes"
)

// Error 封装error
func Error(grpcCode codes.Code, err error) (s *StatusError) {
	if err == nil {
		return nil
	}

	s = &StatusError{
		grpcCode: grpcCode,
		msg:      err.Error(),
		curErr:   err, // 包装 curErr
	}
	s.frdInfo.frdCode = grpcCode.String()

	// 透传最内层友好信息
	var preS *StatusError
	if errors.As(err, &preS) {
		if preS.innermostFrdInfo.frdMsg != "" {
			s.innermostFrdInfo = preS.innermostFrdInfo
		}
	}
	return
}

// Errorf returns New(grpcCode, fmt.Errorf(format, a...)).
func Errorf(grpcCode codes.Code, format string, a ...interface{}) *StatusError {
	return Error(grpcCode, fmt.Errorf(format, a...))
}

// FromError 从error提取StatusError
func FromError(err error) (s *StatusError, ok bool) {
	if err == nil {
		return nil, true
	}
	if errors.As(err, &s) {
		return s, true
	}
	return Error(codes.Unknown, err), false
}

// Convert is a convenience function which removes the need to handle the
// boolean return value from FromError.
func Convert(err error) *StatusError {
	s, _ := FromError(err)
	return s
}
