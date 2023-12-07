package http

import "net/http"

func (s *Server) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("ok"))
}
