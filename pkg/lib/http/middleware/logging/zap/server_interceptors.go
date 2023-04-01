package zap

import (
	"context"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// SystemField is used in every log statement made through http_zap. Can be overwritten before any initialization code.
	SystemField = zap.String("system", "http")

	// ServerField is used in every server-side log statement made through http_zap.Can be overwritten before initialization.
	ServerField = zap.String("span.kind", "server")
)

func LoggingHandler(h http.Handler, logger *zap.Logger, opts ...Option) http.Handler {
	o := evaluateServerOpt(opts)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		newCtx := r.Context()

		// 注入基本数据到ctx
		newCtx = newLoggerForCall(newCtx, logger, r, startTime)

		h.ServeHTTP(w, r.WithContext(newCtx))

		// 完成access log 输出
		finishedExtFields := getAccessLogExtFields(r)
		finishedExtFields = append(finishedExtFields, o.durationFunc(time.Since(startTime))) // 加上耗时
		ctxzap.Extract(newCtx).With(finishedExtFields...).Info("request finished")

	})
}

func newLoggerForCall(ctx context.Context, logger *zap.Logger, req *http.Request, start time.Time) context.Context {
	var f []zapcore.Field
	f = append(f, zap.String("http.start_time", start.Format(time.RFC3339)))
	if d, ok := ctx.Deadline(); ok {
		f = append(f, zap.String("http.request.deadline", d.Format(time.RFC3339)))
	}
	callLog := logger.With(append(f, serverCallFields(req)...)...)
	return ctxzap.ToContext(ctx, callLog)
}

func serverCallFields(req *http.Request) []zapcore.Field {
	return []zapcore.Field{
		SystemField,
		ServerField,
		zap.String("http.method", req.Method),
		zap.String("http.url", req.RequestURI),
	}
}

func getAccessLogExtFields(req *http.Request) []zapcore.Field {
	return []zapcore.Field{
		zap.String("http.method", req.Method),
		zap.String("http.url", req.RequestURI),
		zap.String("http.proto", req.Proto),
		zap.String("http.remote_addr", req.RemoteAddr),
		zap.String("http.user_agent", req.UserAgent()),
	}
}
