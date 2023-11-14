package slog

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/wearemojo/mojo-public-go/lib/clog"
	"github.com/wearemojo/mojo-public-go/lib/merr"
	"github.com/wearemojo/mojo-public-go/lib/mlog"
)

func SetCLogFieldsForGCP() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			ctx := req.Context()

			// resolve full URL
			scheme := "http"
			if req.TLS != nil || req.Header.Get("X-Forwarded-Proto") == "https" {
				scheme = "https"
			}
			baseURL := url.URL{Scheme: scheme, Host: req.Host}
			fullURL := baseURL.ResolveReference(req.URL)

			// https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#httprequest
			httpRequest := map[string]any{
				"requestMethod": req.Method,
				"requestUrl":    fullURL.String(),
				"requestSize":   strconv.FormatInt(req.ContentLength, 10),
				"userAgent":     req.UserAgent(),
				"remoteIp":      req.RemoteAddr,
				"referer":       req.Referer(),
				"protocol":      req.Proto,
				// "serverIp":      "",
				// cache keys are not applicable
			}

			err := clog.SetFields(ctx, clog.Fields{
				"httpRequest": httpRequest,
			})
			if err != nil {
				mlog.Warn(ctx, merr.New(ctx, "clog_set_fields_failed", nil, err))
			}

			// wrap given response writer with one that tracks status code/bytes written
			resWrap := &responseWriter{ResponseWriter: res}

			t1 := time.Now()
			next.ServeHTTP(resWrap, req)
			t2 := time.Now()
			duration := t2.Sub(t1)

			httpRequest["status"] = resWrap.Status
			httpRequest["responseSize"] = resWrap.Bytes
			httpRequest["latency"] = fmt.Sprintf("%.9fs", duration.Seconds())
		})
	}
}

type responseWriter struct {
	http.ResponseWriter

	Status int
	Bytes  int64
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.Status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(data []byte) (int, error) {
	if rw.Status == 0 {
		rw.Status = http.StatusOK
	}

	rw.Bytes += int64(len(data))

	return rw.ResponseWriter.Write(data)
}
