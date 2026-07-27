package api

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"github.com/chankei613/execution-ledger/internal/db"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type issueKeyRequest struct {
	Name string `json:"name"`
}

type issueKeyResponse struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	APIKey  string `json:"api_key"` // 発行時にしか見せない生のキー
}

func generateAPIKey() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

func (s *Server) issueKey(w http.ResponseWriter, r *http.Request) {
	var body issueKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	raw, err := generateAPIKey()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ak := db.AgentKey{
		ID:         uuid.NewString(),
		Name:       body.Name,
		APIKeyHash: HashAPIKey(raw),
		CreatedAt:  time.Now(),
	}
	if err := s.DB.Create(&ak).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, issueKeyResponse{ID: ak.ID, Name: ak.Name, APIKey: raw})
}

func (s *Server) listKeys(w http.ResponseWriter, r *http.Request) {
	var keys []db.AgentKey
	s.DB.Order("created_at asc").Find(&keys)
	writeJSON(w, http.StatusOK, keys)
}

func (s *Server) revokeKey(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	now := time.Now()
	res := s.DB.Model(&db.AgentKey{}).Where("id = ? AND revoked_at IS NULL", id).Update("revoked_at", now)
	if res.Error != nil {
		http.Error(w, res.Error.Error(), http.StatusInternalServerError)
		return
	}
	if res.RowsAffected == 0 {
		http.Error(w, "key not found or already revoked", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
