package tunnel

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	charmlog "github.com/charmbracelet/log"
)

// RequestLog stores the metadata of an intercepted HTTP request.
type RequestLog struct {
	TunnelName string
	Method     string
	Path       string
	Status     int
	Time       time.Time
}

// ProxyStats aggregates traffic metrics.
type ProxyStats struct {
	Requests atomic.Uint64
	Bytes    atomic.Uint64

	mu     sync.RWMutex
	Recent []RequestLog
}

// AddLog appends a request log entry, keeping only the last 5 logs.
func (s *ProxyStats) AddLog(tunnelName, method, path string, status int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Recent = append(s.Recent, RequestLog{
		TunnelName: tunnelName,
		Method:     method,
		Path:       path,
		Status:     status,
		Time:       time.Now(),
	})

	if len(s.Recent) > 10 {
		s.Recent = s.Recent[len(s.Recent)-10:]
	}
}

func (s *ProxyStats) GetRecent() []RequestLog {
	s.mu.RLock()
	defer s.mu.RUnlock()
	cpy := make([]RequestLog, len(s.Recent))
	copy(cpy, s.Recent)
	return cpy
}

// Proxy represents a local reverse proxy that measures traffic and logs requests.
type Proxy struct {
	Stats      ProxyStats
	TargetPort int
	ListenPort int
	Logger     *charmlog.Logger
}

type trackingResponseWriter struct {
	http.ResponseWriter
	stats      *ProxyStats
	statusCode int
}

func (w *trackingResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *trackingResponseWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.stats.Bytes.Add(uint64(n))
	return n, err
}

func (w *trackingResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacker, ok := w.ResponseWriter.(http.Hijacker); ok {
		return hijacker.Hijack()
	}
	return nil, nil, fmt.Errorf("http.Hijacker not implemented")
}

func (w *trackingResponseWriter) Flush() {
	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

func (w *trackingResponseWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

// StartProxy starts a local reverse proxy that tracks requests and bytes.
func StartProxy(ctx context.Context, tunnelName string, targetPort int) (*Proxy, error) {
	targetURL, err := url.Parse(fmt.Sprintf("http://localhost:%d", targetPort))
	if err != nil {
		return nil, fmt.Errorf("parse target URL http://localhost:%d: %w", targetPort, err)
	}

	p := &Proxy{TargetPort: targetPort}

	rp := httputil.NewSingleHostReverseProxy(targetURL)
	// We can use the logger's standard log adapter for errors
	logger := charmlog.New(io.Discard)
	rp.ErrorLog = logger.StandardLog(charmlog.StandardLogOptions{ForceLevel: charmlog.ErrorLevel})



	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Intercept internal ping requests to hide them from logs and avoid hitting the user's local server
		if r.Header.Get("User-Agent") == "Nipo-Ping" {
			w.WriteHeader(http.StatusOK)
			return
		}

		p.Stats.Requests.Add(1)
		
		tw := &trackingResponseWriter{ResponseWriter: w, stats: &p.Stats, statusCode: 200}
		rp.ServeHTTP(tw, r)

		// Record the log after the response is sent
		p.Stats.AddLog(tunnelName, r.Method, r.URL.Path, tw.statusCode)
	})

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, fmt.Errorf("listen on random TCP port: %w", err)
	}

	p.ListenPort = listener.Addr().(*net.TCPAddr).Port

	srv := &http.Server{
		Handler: handler,
	}

	go func() {
		if err := srv.Serve(listener); err != nil && err != http.ErrServerClosed {
			// server exited unexpectedly
		}
	}()

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdownCtx)
	}()

	return p, nil
}
