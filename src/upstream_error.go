package src

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
)

// isUpstreamErrorStatus reports whether statusCode is one Traefik's own reverse proxy returns
// when the backend behind this route is unreachable, has no available server, or times out —
// as opposed to a status code the app itself chose to return.
func isUpstreamErrorStatus(statusCode int) bool {
	switch statusCode {
	case http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return true
	default:
		return false
	}
}

// upstreamErrorInterceptor wraps the real ResponseWriter passed to next.ServeHTTP so the
// plugin can inspect the upstream's status code before any of its body reaches the client.
//
// If the first WriteHeader call reports a healthy status, everything is forwarded to the
// underlying ResponseWriter untouched and unbuffered from that point on, so normal responses
// (including large or streamed ones) pass through with no extra copying.
//
// If it reports one of isUpstreamErrorStatus's codes, the interceptor switches into
// "intercepting" mode: the header is never sent and the body is discarded, so the caller can
// render the fallback page onto the real ResponseWriter afterwards as if nothing had been
// written yet.
type upstreamErrorInterceptor struct {
	rw http.ResponseWriter

	headerWritten bool
	intercepting  bool
}

func (w *upstreamErrorInterceptor) Header() http.Header {
	return w.rw.Header()
}

func (w *upstreamErrorInterceptor) WriteHeader(statusCode int) {
	if w.headerWritten || w.intercepting {
		return
	}

	if isUpstreamErrorStatus(statusCode) {
		w.intercepting = true
		return
	}

	w.rw.WriteHeader(statusCode)
	w.headerWritten = true
}

func (w *upstreamErrorInterceptor) Write(p []byte) (int, error) {
	if w.intercepting {
		return len(p), nil
	}

	if !w.headerWritten {
		w.WriteHeader(http.StatusOK)
	}

	return w.rw.Write(p)
}

// Flush lets a healthy, streamed response (e.g. SSE or chunked output) reach the client
// incrementally instead of being buffered until the handler returns. net/http/httputil's
// reverse proxy type-asserts for http.Flusher and silently skips flushing without it.
func (w *upstreamErrorInterceptor) Flush() {
	if w.intercepting {
		return
	}

	if !w.headerWritten {
		w.WriteHeader(http.StatusOK)
	}

	if f, ok := w.rw.(http.Flusher); ok {
		f.Flush()
	}
}

// Hijack passes through to the underlying ResponseWriter's http.Hijacker, required for
// protocol upgrades (WebSockets) proxied through this route to work at all. An upgrade
// response is never one of isUpstreamErrorStatus's codes, so there's nothing to intercept
// here regardless.
func (w *upstreamErrorInterceptor) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := w.rw.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("underlying ResponseWriter does not support hijacking")
	}

	return hijacker.Hijack()
}
