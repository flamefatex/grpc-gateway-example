package grpc_gateway

import (
	"fmt"
	"net/textproto"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

const XHeaderPrefix = "X-"

func OutgoingHeaderMatcher(key string) (string, bool) {
	if strings.HasPrefix(key, runtime.MetadataPrefix) {
		return key[len(runtime.MetadataPrefix):], true
	}

	return fmt.Sprintf("%s%s", runtime.MetadataHeaderPrefix, key), true
}

func IncomingHeaderMatcher(key string) (string, bool) {
	key = textproto.CanonicalMIMEHeaderKey(key)
	if strings.HasPrefix(key, XHeaderPrefix) {
		return key, true
	}
	return runtime.DefaultHeaderMatcher(key)
}
