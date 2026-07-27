package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

// Server は全ロジックの実体。HTTPハンドラーとWailsネイティブバインディングの
// 両方がこの同じ Server のメソッドを呼ぶことで、UIとAPIの挙動がズレないようにする。
type Server struct {
	DB *gorm.DB
}

func New(conn *gorm.DB) *Server {
	return &Server{DB: conn}
}

func (s *Server) Router() http.Handler {
	r := chi.NewRouter()

	r.Route("/api/v1/keys", func(r chi.Router) {
		r.Use(APIKeyAuth(s.DB, "/api/v1/keys"))
		r.Post("/", s.httpIssueKey)
		r.Get("/", s.httpListKeys)
		r.Delete("/{id}", s.httpRevokeKey)
	})

	r.Route("/api/v1/entries", func(r chi.Router) {
		r.Use(APIKeyAuth(s.DB))
		r.Post("/", s.httpIngestEntry)
		r.Get("/", s.httpListEntries)
		r.Get("/export", s.httpExportEntries)
		r.Get("/{id}", s.httpGetEntry)
	})

	r.Route("/api/v1/stats", func(r chi.Router) {
		r.Use(APIKeyAuth(s.DB))
		r.Get("/", s.httpStats)
	})

	return r
}

// NewRouter はcmd/elserve（単体HTTPサーバー）向けの簡易コンストラクタ。
func NewRouter(conn *gorm.DB) http.Handler {
	return New(conn).Router()
}
