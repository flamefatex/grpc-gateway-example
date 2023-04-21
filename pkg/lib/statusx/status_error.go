package statusx

import (
	"context"
	"fmt"

	"github.com/flamefatex/grpc-gateway-example/pkg/lib/statusx/code"
	proto_status "github.com/flamefatex/grpc-gateway-example/proto/gen/go/common/status"
	"google.golang.org/grpc/codes"
	grpc_status "google.golang.org/grpc/status"
)

// FrdInfo 友好信息，用于前端返回
type FrdInfo struct {
	frdCode code.Code // 友好code，用于接口返回
	frdMsg  string    // 友好信息，用于接口返回
}

type StatusError struct {
	grpcCode         codes.Code // grpc code
	msg              string     // 原始错误信息,用于本地日志，避免国际化影响
	curErr           error      // 包装err
	frdInfo          FrdInfo    // 友好信息
	innermostFrdInfo FrdInfo    // 最内层的友好信息，用于透传最内层的友好信息，返回前端
}

func (s *StatusError) Error() string {
	return s.msg
}

// Unwrap 把内层的error 暴露出去
func (s *StatusError) Unwrap() error {
	u, ok := s.curErr.(interface {
		Unwrap() error
	})
	// 如果不是，返回nil
	if !ok {
		return nil
	}
	// 否则则调用该error的Unwrap方法返回被嵌套的error
	return u.Unwrap()
}

// GRPCStatus grpc status convert使用的接口,用于grpc通信
func (s *StatusError) GRPCStatus() *grpc_status.Status {
	rs := s.FrdProto()
	st, _ := grpc_status.New(s.grpcCode, s.msg).WithDetails(rs)
	return st
}

func (s *StatusError) Code() codes.Code {
	return s.grpcCode
}
func (s *StatusError) FrdCode() string {
	return s.frdInfo.frdCode
}

// FrdMsg 返回当前错误的友好信息
func (s *StatusError) FrdMsg() string {
	msg := s.frdInfo.frdMsg
	if msg == "" {
		msg = s.msg
	}
	return msg
}

// FrdProto 获取友好信息响应的proto,用于前端友好显示
func (s *StatusError) FrdProto() (rs *proto_status.ResponseStatus) {
	msg := s.FrdMsg()
	if s.innermostFrdInfo.frdMsg != "" {
		msg = s.innermostFrdInfo.frdMsg
	}
	rs = &proto_status.ResponseStatus{
		Code:    uint32(s.grpcCode),
		Reason:  s.frdInfo.frdCode,
		Message: msg,
	}
	return
}

// WithFrdMsg 自定义友好信息 用于http返回json给前端
func (s *StatusError) WithFrdMsg(frdMsg string) *StatusError {
	s.frdInfo.frdMsg = frdMsg

	// 如果最内层的友好信息为空，则赋值
	if s.innermostFrdInfo.frdMsg == "" {
		s.innermostFrdInfo = s.frdInfo
	}
	return s
}

func (s *StatusError) WithFrdMsgf(formatFrdMsg string, a ...interface{}) *StatusError {
	return s.WithFrdMsg(fmt.Sprintf(formatFrdMsg, a...))
}

func OK(ctx context.Context) *proto_status.ResponseStatus {
	return &proto_status.ResponseStatus{
		Code:    uint32(codes.OK),
		Reason:  code.OK,
		Message: "OK",
	}
}
