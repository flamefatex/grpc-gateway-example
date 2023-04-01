package inject_ctx

import (
	"net/http"
)

func NewHandler(h http.Handler, opts ...Option) http.Handler {
	o := evaluateOptions(opts)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		newCtx := r.Context()

		if len(o.injectFromRequestCtxFuncs) != 0 {
			for _, fn := range o.injectFromRequestCtxFuncs {
				newCtx = fn(newCtx, r)
			}
		}

		h.ServeHTTP(w, r.WithContext(newCtx))
	})
}
