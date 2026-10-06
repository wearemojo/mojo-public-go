package slog

import (
	"fmt"
	"maps"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/wearemojo/mojo-public-go/lib/clog"
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

			clog.SetFields(ctx, clog.Fields{
				"httpRequest": httpRequest,
			})

			// wrap given response writer with one that tracks status code/bytes written
			resWrap := &responseWriter{ResponseWriter: res}

			t1 := time.Now()
			next.ServeHTTP(resWrap, req)
			t2 := time.Now()
			duration := t2.Sub(t1)

			// Goroutines started by the request can still be logging the published map.
			completed := maps.Clone(httpRequest)
			completed["status"] = resWrap.Status
			completed["responseSize"] = resWrap.Bytes
			completed["latency"] = fmt.Sprintf("%.9fs", duration.Seconds())

			clog.SetFields(ctx, clog.Fields{
				"httpRequest": completed,
			})
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
