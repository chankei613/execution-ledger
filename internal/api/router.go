package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type Server struct {
	DB *gorm.DB
}

func NewRouter(conn *gorm.DB) http.Handler {
	s := &Server{DB: conn}
	r := chi.NewRouter()

	r.Route("/api/v1/keys", func(r chi.Router) {
		r.Use(APIKeyAuth(conn, "/api/v1/keys"))
		r.Post("/", s.issueKey)
		r.Get("/", s.listKeys)
		r.Delete("/{id}", s.revokeKey)
	})

	r.Route("/api/v1/entries", func(r chi.Router) {
		r.Use(APIKeyAuth(conn))
		r.Post("/", s.ingestEntry)
		r.Get("/", s.listEntries)
		r.Get("/export", s.exportEntries)
		r.Get("/{id}", s.getEntry)
	})

	r.Route("/api/v1/stats", func(r chi.Router) {
		r.Use(APIKeyAuth(conn))
		r.Get("/", s.stats)
	})

	return r
}
