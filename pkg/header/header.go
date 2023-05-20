package header

import (
	"context"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	// 客户端侧
	XAppIdHeader = "x-app-id"

	// 服务端侧
	XHttpCode            = "x-http-code"
	XForwardedForHeader  = "x-forwarded-for"
	XForwardedHostHeader = "x-forwarded-host"
)

// GetHeader 获取header
func GetHeader(ctx context.Context, key string) string {
	s := metadata.ValueFromIncomingContext(ctx, key)
	if len(s) < 1 {
		return ""
	}
	return s[0]
}

// SetHeader 设置header
func SetHeader(ctx context.Context, key string, value string) {
	_ = grpc.SetHeader(ctx, metadata.Pairs(key, value))
}

// GetHttpHeader 获取http header
func GetHttpHeader(ctx context.Context, key string) string {
	return GetHeader(ctx, runtime.MetadataPrefix+key)
}

// SetHttpHeader 设置http header, key为http header的key,如SetCookie
func SetHttpHeader(ctx context.Context, key string, value string) {
	SetHeader(ctx, runtime.MetadataPrefix+key, value)
}

// SetHttpCode 设置http code
func SetHttpCode(ctx context.Context, code string) {
	SetHeader(ctx, XHttpCode, code)
}

// GetAppId 获取应用Id
func GetAppId(ctx context.Context) string {
	return GetHeader(ctx, XAppIdHeader)
}

// GetDomain 获取域名
func GetDomain(ctx context.Context) string {
	host := GetHeader(ctx, XForwardedHostHeader)
	if host != "" {
		s := strings.Split(host, ":")
		return s[0]
	}
	return ""
}

// GetDomainIdentifier 获取域名标识
func GetDomainIdentifier(ctx context.Context) string {
	domain := GetDomain(ctx)

	if domain != "" {
		s := strings.Split(domain, ".")
		return s[0]
	}
	return ""
}
