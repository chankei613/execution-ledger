package api

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/chankei613/execution-ledger/internal/db"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// IngestEntryInput はクライアントから受け取るフィールドのみを列挙する。
// id / received_at はサーバー側で必ず付与し、クライアント指定を無視することで
// イミュータブル・改ざん防止の起点を保証する。
type IngestEntryInput struct {
	Source  string `json:"source"`
	AgentID string `json:"agent_id"`
	Subject string `json:"subject"`

	Status  db.EntryStatus `json:"status"`
	Summary string         `json:"summary"`

	CriteriaResults []db.CriterionResult `json:"criteria_results"`
	Outputs         map[string]any       `json:"outputs"`

	Confidence struct {
		Overall            float64                `json:"overall"`
		Breakdown          db.ConfidenceBreakdown `json:"breakdown"`
		LowConfidenceAreas []string               `json:"low_confidence_areas"`
	} `json:"confidence"`

	Decisions    []db.Decision `json:"decisions"`
	ActionsTaken []db.Action   `json:"actions_taken"`
	FollowUp     []db.FollowUp `json:"follow_up"`

	Usage db.Usage `json:"usage"`
}

// IngestEntry はエントリを1件追記する（HTTP・ネイティブバインディング共用）。
func (s *Server) IngestEntry(body IngestEntryInput) (db.LedgerEntry, error) {
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

	err := s.DB.Create(&entry).Error
	return entry, err
}

func (s *Server) GetEntry(id string) (db.LedgerEntry, error) {
	var entry db.LedgerEntry
	err := s.DB.First(&entry, "id = ?", id).Error
	return entry, err
}

// EntryFilters は検索・集計・エクスポートで共用する絞り込み条件。
type EntryFilters struct {
	AgentID       string
	Source        string
	Status        string
	Subject       string
	Query         string
	MinConfidence *float64
	MaxConfidence *float64
	From          *time.Time
	To            *time.Time
}

func FiltersFromQuery(q url.Values) EntryFilters {
	f := EntryFilters{
		AgentID: q.Get("agent_id"),
		Source:  q.Get("source"),
		Status:  q.Get("status"),
		Subject: q.Get("subject"),
		Query:   q.Get("q"),
	}
	if v := q.Get("min_confidence"); v != "" {
		if parsed, err := strconv.ParseFloat(v, 64); err == nil {
			f.MinConfidence = &parsed
		}
	}
	if v := q.Get("max_confidence"); v != "" {
		if parsed, err := strconv.ParseFloat(v, 64); err == nil {
			f.MaxConfidence = &parsed
		}
	}
	if v := q.Get("from"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			f.From = &t
		}
	}
	if v := q.Get("to"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			f.To = &t
		}
	}
	return f
}

func (f EntryFilters) apply(query *gorm.DB) *gorm.DB {
	if f.AgentID != "" {
		query = query.Where("agent_id = ?", f.AgentID)
	}
	if f.Source != "" {
		query = query.Where("source = ?", f.Source)
	}
	if f.Status != "" {
		query = query.Where("status = ?", f.Status)
	}
	if f.Subject != "" {
		query = query.Where("subject LIKE ?", "%"+f.Subject+"%")
	}
	if f.Query != "" {
		query = query.Where("summary LIKE ?", "%"+f.Query+"%")
	}
	if f.MinConfidence != nil {
		query = query.Where("confidence_overall >= ?", *f.MinConfidence)
	}
	if f.MaxConfidence != nil {
		query = query.Where("confidence_overall <= ?", *f.MaxConfidence)
	}
	if f.From != nil {
		query = query.Where("received_at >= ?", *f.From)
	}
	if f.To != nil {
		query = query.Where("received_at <= ?", *f.To)
	}
	return query
}

type SearchResult struct {
	Entries []db.LedgerEntry `json:"entries"`
	Total   int64            `json:"total"`
}

func (s *Server) SearchEntries(f EntryFilters, limit, offset int) (SearchResult, error) {
	if limit <= 0 || limit > 500 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	var total int64
	if err := f.apply(s.DB.Model(&db.LedgerEntry{})).Count(&total).Error; err != nil {
		return SearchResult{}, err
	}

	var entries []db.LedgerEntry
	err := f.apply(s.DB.Model(&db.LedgerEntry{})).
		Order("received_at desc").
		Limit(limit).
		Offset(offset).
		Find(&entries).Error

	return SearchResult{Entries: entries, Total: total}, err
}

type StatsResult struct {
	Total             int64            `json:"total"`
	ByStatus          map[string]int64 `json:"by_status"`
	AvgConfidence     float64          `json:"avg_confidence"`
	LowConfidenceRate float64          `json:"low_confidence_rate"` // confidence_overall < 0.6 の割合
	ByAgent           map[string]int64 `json:"by_agent"`
}

const lowConfidenceThreshold = 0.6

func (s *Server) Stats(f EntryFilters) (StatsResult, error) {
	var total int64
	if err := f.apply(s.DB.Model(&db.LedgerEntry{})).Count(&total).Error; err != nil {
		return StatsResult{}, err
	}

	resp := StatsResult{Total: total, ByStatus: map[string]int64{}, ByAgent: map[string]int64{}}
	if total == 0 {
		return resp, nil
	}

	var avg float64
	f.apply(s.DB.Model(&db.LedgerEntry{})).Select("avg(confidence_overall)").Row().Scan(&avg)
	resp.AvgConfidence = avg

	var lowCount int64
	f.apply(s.DB.Model(&db.LedgerEntry{})).Where("confidence_overall < ?", lowConfidenceThreshold).Count(&lowCount)
	resp.LowConfidenceRate = float64(lowCount) / float64(total)

	type statusCount struct {
		Status db.EntryStatus
		Count  int64
	}
	var statusCounts []statusCount
	f.apply(s.DB.Model(&db.LedgerEntry{})).Select("status, count(*) as count").Group("status").Scan(&statusCounts)
	for _, sc := range statusCounts {
		resp.ByStatus[string(sc.Status)] = sc.Count
	}

	type agentCount struct {
		AgentID string
		Count   int64
	}
	var agentCounts []agentCount
	f.apply(s.DB.Model(&db.LedgerEntry{})).Select("agent_id, count(*) as count").Group("agent_id").Scan(&agentCounts)
	for _, ac := range agentCounts {
		resp.ByAgent[ac.AgentID] = ac.Count
	}

	return resp, nil
}

func (s *Server) ExportEntries(f EntryFilters) ([]db.LedgerEntry, error) {
	var entries []db.LedgerEntry
	err := f.apply(s.DB.Model(&db.LedgerEntry{})).Order("received_at desc").Find(&entries).Error
	return entries, err
}

func exportCSV(entries []db.LedgerEntry) string {
	var buf bytes.Buffer
	cw := csv.NewWriter(&buf)
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
	return buf.String()
}

// ─── HTTPハンドラー（外部プロセスからのIngestion用。薄いラッパー） ────────────

func (s *Server) httpIngestEntry(w http.ResponseWriter, r *http.Request) {
	var body IngestEntryInput
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if body.Status == "" || body.AgentID == "" {
		http.Error(w, "agent_id and status are required", http.StatusBadRequest)
		return
	}

	entry, err := s.IngestEntry(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, entry)
}

func (s *Server) httpGetEntry(w http.ResponseWriter, r *http.Request) {
	entry, err := s.GetEntry(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, entry)
}

func (s *Server) httpListEntries(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))

	result, err := s.SearchEntries(FiltersFromQuery(q), limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (s *Server) httpStats(w http.ResponseWriter, r *http.Request) {
	result, err := s.Stats(FiltersFromQuery(r.URL.Query()))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (s *Server) httpExportEntries(w http.ResponseWriter, r *http.Request) {
	entries, err := s.ExportEntries(FiltersFromQuery(r.URL.Query()))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.URL.Query().Get("format") == "csv" {
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", "attachment; filename=ledger-export.csv")
		w.Write([]byte(exportCSV(entries)))
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
