package grpc_gateway

import (
	"context"
	"net/http"
	"strconv"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/proto"
)

// HttpResponseModifier http响应重写
func HttpResponseModifier(ctx context.Context, w http.ResponseWriter, p proto.Message) error {
	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		return nil
	}

	// set http status code
	if vals := md.HeaderMD.Get("x-http-code"); len(vals) > 0 {
		code, err := strconv.Atoi(vals[0])
		if err != nil {
			return err
		}
		// delete the headers to not expose any grpc-metadata in http response
		delete(md.HeaderMD, "x-http-code")
		delete(w.Header(), "Grpc-Metadata-X-Http-Code")
		w.WriteHeader(code)
	}

	return nil
}
