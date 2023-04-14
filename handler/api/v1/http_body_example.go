package v1

import (
	"context"

	proto_v1_hbe "github.com/flamefatex/grpc-gateway-example/proto/gen/go/api/v1/http_body_example"
	"google.golang.org/genproto/googleapis/api/httpbody"
	"google.golang.org/protobuf/types/known/emptypb"
)

type httpBodyExampleHandler struct {
	proto_v1_hbe.UnimplementedHttpBodyExampleServiceServer
}

func NewHttpBodyExampleHandler() *httpBodyExampleHandler {
	return &httpBodyExampleHandler{}
}

func (h *httpBodyExampleHandler) HelloWorld(ctx context.Context, empty *emptypb.Empty) (*httpbody.HttpBody, error) {
	return &httpbody.HttpBody{
		ContentType: "text/html",
		Data:        []byte("Hello World"),
	}, nil
}
func (h *httpBodyExampleHandler) Download(empty *emptypb.Empty, stream proto_v1_hbe.HttpBodyExampleService_DownloadServer) error {
	msgs := []*httpbody.HttpBody{
		{
			ContentType: "text/html",
			Data:        []byte("Hello 1"),
		},
		{
			ContentType: "text/html",
			Data:        []byte("Hello 2"),
		},
	}

	for _, msg := range msgs {
		if err := stream.Send(msg); err != nil {
			return err
		}
	}

	return nil
}
