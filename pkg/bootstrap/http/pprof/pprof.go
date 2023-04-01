package pprof

import (
	"context"
	"net/http"
	"net/http/pprof"
)

func GetPprofMux(ctx context.Context) *http.ServeMux {
	httpMux := http.NewServeMux()
	httpMux.HandleFunc("/debug/pprof/", pprof.Index)
	httpMux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	httpMux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	httpMux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	httpMux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	return httpMux
}
