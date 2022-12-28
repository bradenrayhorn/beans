package http

import (
	"mime"
	"net/http"
	"path"
	"path/filepath"

	"github.com/bradenrayhorn/beans/web"
	"github.com/rs/zerolog/log"
)

func (s *Server) handleServeFrontend() http.HandlerFunc {
	errHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}

	index, err := web.FrontendFS.ReadFile("dist/index.html")
	if err != nil {
		log.Error().Msg(err.Error())
		return errHandler
	}

	return func(w http.ResponseWriter, r *http.Request) {
		file, err := web.FrontendFS.ReadFile(path.Join("dist", r.URL.Path))
		if err == nil {
			w.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(r.URL.Path)))
			w.Write(file)
			return
		}

		w.Header().Set("Content-Type", mime.TypeByExtension(".html"))
		w.Write(index)
	}
}
