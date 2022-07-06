package http

import "net/http"

func (s *Server) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}
