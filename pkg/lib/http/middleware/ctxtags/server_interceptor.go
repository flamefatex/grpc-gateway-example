package ctxtags

import (
	"context"
	"net/http"

	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
)

func CtxTagsHandler(h http.Handler, opts ...Option) http.Handler {
	o := evaluateOptions(opts)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		newCtx := newTagsForCtx(r.Context())

		if o.requestFieldsExtractorFunc != nil {
			setRequestFieldTags(newCtx, o.requestFieldsExtractorFunc, r)
		}

		h.ServeHTTP(w, r.WithContext(newCtx))
	})
}

func newTagsForCtx(ctx context.Context) context.Context {
	t := grpc_ctxtags.NewTags()
	return grpc_ctxtags.SetInContext(ctx, t)
}

func setRequestFieldTags(ctx context.Context, fn RequestFieldExtractorFunc, req *http.Request) {
	if valMap := fn(req); valMap != nil {
		t := grpc_ctxtags.Extract(ctx)
		for k, v := range valMap {
			t.Set("http.request."+k, v)
		}
	}
}
