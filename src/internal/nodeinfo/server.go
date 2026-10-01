package nodeinfo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"syscall"
)

type Server struct {
	provider StatusProvider
	httpSrv  *http.Server
}

func NewServer(p StatusProvider) *Server {
	s := &Server{provider: p}
	mux := s.Routes()
	s.httpSrv = &http.Server{Handler: mux}
	return s
}

func (s *Server) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", s.handleHealth)
	mux.HandleFunc("GET /status", s.handleStatus)
	return mux
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	status := s.provider.CurrentStatus(r.Context())
	body, err := json.Marshal(status)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}

func (s *Server) Start(addr string) error {
	if addr == "" || addr[0] != ':' {
		return fmt.Errorf("listen address must be :port (e.g. :12345), got %q", addr)
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		if isAddrInUse(err) {
			return fmt.Errorf("HTTP server port %s is already in use: %w", addr, err)
		}
		return fmt.Errorf("HTTP server listen on %s: %w", addr, err)
	}
	s.httpSrv.Addr = ln.Addr().String()
	go func() {
		_ = s.httpSrv.Serve(ln)
	}()
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.httpSrv == nil {
		return nil
	}
	return s.httpSrv.Shutdown(ctx)
}

func isAddrInUse(err error) bool {
	var opErr *net.OpError
	if errors.As(err, &opErr) {
		var sysErr syscall.Errno
		if errors.As(opErr.Err, &sysErr) {
			return sysErr == syscall.EADDRINUSE
		}
	}
	return errors.Is(err, syscall.EADDRINUSE)
}
