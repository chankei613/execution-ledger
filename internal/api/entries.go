package api

import (
	"encoding/csv"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/chankei613/execution-ledger/internal/db"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ingestEntryRequest はクライアントから受け取るフィールドのみを列挙する。
// id / received_at はサーバー側で必ず付与し、クライアント指定を無視することで
// イミュータブル・改ざん防止の起点を保証する。
type ingestEntryRequest struct {
	Source  string `json:"source"`
	AgentID string `json:"agent_id"`
	Subject string `json:"subject"`

	Status  db.EntryStatus `json:"status"`
	Summary string         `json:"summary"`

	CriteriaResults []db.CriterionResult `json:"criteria_results"`
	Outputs         map[string]any       `json:"outputs"`

	Confidence struct {
		Overall           float64             `json:"overall"`
		Breakdown         db.ConfidenceBreakdown `json:"breakdown"`
		LowConfidenceAreas []string            `json:"low_confidence_areas"`
	} `json:"confidence"`

	Decisions    []db.Decision `json:"decisions"`
	ActionsTaken []db.Action   `json:"actions_taken"`
	FollowUp     []db.FollowUp `json:"follow_up"`

	Usage db.Usage `json:"usage"`
}

func (s *Server) ingestEntry(w http.ResponseWriter, r *http.Request) {
	var body ingestEntryRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if body.Status == "" || body.AgentID == "" {
		http.Error(w, "agent_id and status are required", http.StatusBadRequest)
		return
	}

	entry := db.LedgerEntry{
		ID:         uuid.NewString(),
		ReceivedAt: time.Now(),

		Source:  body.Source,
		AgentID: body.AgentID,
		Subject: body.Subject,

		Status:  body.Status,
		Summary: body.Summary,

		CriteriaResults: body.CriteriaResults,
		Outputs:         body.Outputs,

		ConfidenceOverall:   body.Confidence.Overall,
		ConfidenceBreakdown: body.Confidence.Breakdown,
		LowConfidenceAreas:  body.Confidence.LowConfidenceAreas,

		Decisions:    body.Decisions,
		ActionsTaken: body.ActionsTaken,
		FollowUp:     body.FollowUp,

		Usage: body.Usage,
	}

	if err := s.DB.Create(&entry).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, entry)
}

func (s *Server) getEntry(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var entry db.LedgerEntry
	if err := s.DB.First(&entry, "id = ?", id).Error; err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, entry)
}

// entryFilters はクエリパラメータから検索条件を組み立てる。search/export/statsで共用する。
func entryFilters(r *http.Request, query *gorm.DB) *gorm.DB {
	q := r.URL.Query()

	if v := q.Get("agent_id"); v != "" {
		query = query.Where("agent_id = ?", v)
	}
	if v := q.Get("source"); v != "" {
		query = query.Where("source = ?", v)
	}
	if v := q.Get("status"); v != "" {
		query = query.Where("status = ?", v)
	}
	if v := q.Get("subject"); v != "" {
		query = query.Where("subject LIKE ?", "%"+v+"%")
	}
	if v := q.Get("q"); v != "" {
		query = query.Where("summary LIKE ?", "%"+v+"%")
	}
	if v := q.Get("min_confidence"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			query = query.Where("confidence_overall >= ?", f)
		}
	}
	if v := q.Get("max_confidence"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			query = query.Where("confidence_overall <= ?", f)
		}
	}
	if v := q.Get("from"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			query = query.Where("received_at >= ?", t)
		}
	}
	if v := q.Get("to"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			query = query.Where("received_at <= ?", t)
		}
	}

	return query
}

type entriesResponse struct {
	Entries []db.LedgerEntry `json:"entries"`
	Total   int64            `json:"total"`
}

func (s *Server) listEntries(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	limit := 50
	offset := 0
	if v, err := strconv.Atoi(q.Get("limit")); err == nil && v > 0 && v <= 500 {
		limit = v
	}
	if v, err := strconv.Atoi(q.Get("offset")); err == nil && v >= 0 {
		offset = v
	}

	base := entryFilters(r, s.DB.Model(&db.LedgerEntry{}))

	var total int64
	base.Count(&total)

	var entries []db.LedgerEntry
	entryFilters(r, s.DB).
		Order("received_at desc").
		Limit(limit).
		Offset(offset).
		Find(&entries)

	writeJSON(w, http.StatusOK, entriesResponse{Entries: entries, Total: total})
}

type statsResponse struct {
	Total             int64            `json:"total"`
	ByStatus          map[string]int64 `json:"by_status"`
	AvgConfidence     float64          `json:"avg_confidence"`
	LowConfidenceRate float64          `json:"low_confidence_rate"` // confidence_overall < 0.6 の割合
	ByAgent           map[string]int64 `json:"by_agent"`
}

const lowConfidenceThreshold = 0.6

func (s *Server) stats(w http.ResponseWriter, r *http.Request) {
	base := entryFilters(r, s.DB.Model(&db.LedgerEntry{}))

	var total int64
	base.Count(&total)

	resp := statsResponse{Total: total, ByStatus: map[string]int64{}, ByAgent: map[string]int64{}}

	if total > 0 {
		var avg float64
		entryFilters(r, s.DB.Model(&db.LedgerEntry{})).Select("avg(confidence_overall)").Row().Scan(&avg)
		resp.AvgConfidence = avg

		var lowCount int64
		entryFilters(r, s.DB.Model(&db.LedgerEntry{})).Where("confidence_overall < ?", lowConfidenceThreshold).Count(&lowCount)
		resp.LowConfidenceRate = float64(lowCount) / float64(total)

		type statusCount struct {
			Status db.EntryStatus
			Count  int64
		}
		var statusCounts []statusCount
		entryFilters(r, s.DB.Model(&db.LedgerEntry{})).Select("status, count(*) as count").Group("status").Scan(&statusCounts)
		for _, sc := range statusCounts {
			resp.ByStatus[string(sc.Status)] = sc.Count
		}

		type agentCount struct {
			AgentID string
			Count   int64
		}
		var agentCounts []agentCount
		entryFilters(r, s.DB.Model(&db.LedgerEntry{})).Select("agent_id, count(*) as count").Group("agent_id").Scan(&agentCounts)
		for _, ac := range agentCounts {
			resp.ByAgent[ac.AgentID] = ac.Count
		}
	}

	writeJSON(w, http.StatusOK, resp)
}

func (s *Server) exportEntries(w http.ResponseWriter, r *http.Request) {
	var entries []db.LedgerEntry
	entryFilters(r, s.DB.Model(&db.LedgerEntry{})).Order("received_at desc").Find(&entries)

	format := r.URL.Query().Get("format")
	if format == "csv" {
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", "attachment; filename=ledger-export.csv")
		cw := csv.NewWriter(w)
		cw.Write([]string{"id", "received_at", "source", "agent_id", "subject", "status", "confidence_overall", "summary"})
		for _, e := range entries {
			cw.Write([]string{
				e.ID,
				e.ReceivedAt.Format(time.RFC3339),
				e.Source,
				e.AgentID,
				e.Subject,
				string(e.Status),
				strconv.FormatFloat(e.ConfidenceOverall, 'f', 3, 64),
				e.Summary,
			})
		}
		cw.Flush()
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", "attachment; filename=ledger-export.json")
	json.NewEncoder(w).Encode(entries)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
