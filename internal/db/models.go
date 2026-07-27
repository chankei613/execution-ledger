// Package db はExecution LedgerのGORMモデルとSQLite初期化を提供する。
// スキーマは comet-taskAI の schema/types.ts (ExecutionResult) を単体製品向けに
// source/subject/agent_id を自由記述にして緩めたもの。docs/spec.md 参照。
package db

import "time"

type EntryStatus string

const (
	StatusSuccess            EntryStatus = "success"
	StatusPartialSuccess      EntryStatus = "partial_success"
	StatusFailed              EntryStatus = "failed"
	StatusBlocked             EntryStatus = "blocked"
	StatusTimedOut            EntryStatus = "timed_out"
	StatusTokenBudgetExceeded EntryStatus = "token_budget_exceeded"
	StatusGenerated           EntryStatus = "generated"
)

type CriterionResult struct {
	Description string `json:"description"`
	Met         bool   `json:"met"`
	Evidence    string `json:"evidence,omitempty"`
}

type ConfidenceBreakdown struct {
	TaskUnderstood    float64 `json:"task_understood"`
	ExecutionComplete float64 `json:"execution_complete"`
	Correctness       float64 `json:"correctness"`
	SideEffectsClean  float64 `json:"side_effects_clean"`
}

type Decision struct {
	Description            string   `json:"description"`
	Rationale               string   `json:"rationale"`
	AlternativesConsidered  []string `json:"alternatives_considered"`
}

type Action struct {
	Tool         string    `json:"tool"`
	InputSummary string    `json:"input_summary"`
	Timestamp    time.Time `json:"timestamp"`
}

type FollowUp struct {
	Description    string         `json:"description"`
	SuggestedTask  map[string]any `json:"suggested_task,omitempty"`
}

type Usage struct {
	InputTokens      int            `json:"input_tokens"`
	OutputTokens     int            `json:"output_tokens"`
	MCPCallsByServer map[string]int `json:"mcp_calls_by_server,omitempty"`
}

// LedgerEntry は追記専用（append-only）。UPDATE/DELETEは提供しない。
type LedgerEntry struct {
	ID         string    `gorm:"primaryKey" json:"id"`
	ReceivedAt time.Time `gorm:"index" json:"received_at"`

	Source  string `gorm:"index" json:"source"`
	AgentID string `gorm:"index" json:"agent_id"`
	Subject string `gorm:"index" json:"subject"`

	Status  EntryStatus `gorm:"index" json:"status"`
	Summary string      `json:"summary"`

	CriteriaResults []CriterionResult `gorm:"serializer:json" json:"criteria_results"`
	Outputs         map[string]any    `gorm:"serializer:json" json:"outputs"`

	ConfidenceOverall    float64             `gorm:"index" json:"confidence_overall"`
	ConfidenceBreakdown  ConfidenceBreakdown `gorm:"serializer:json" json:"confidence_breakdown"`
	LowConfidenceAreas   []string            `gorm:"serializer:json" json:"low_confidence_areas"`

	Decisions    []Decision `gorm:"serializer:json" json:"decisions"`
	ActionsTaken []Action   `gorm:"serializer:json" json:"actions_taken"`
	FollowUp     []FollowUp `gorm:"serializer:json" json:"follow_up"`

	Usage Usage `gorm:"serializer:json" json:"usage"`
}

// AgentKey — Ingestion API を叩くためのAPIキー。ハッシュのみ保存する。
type AgentKey struct {
	ID         string     `gorm:"primaryKey" json:"id"`
	Name       string     `json:"name"`
	APIKeyHash string     `json:"-"`
	CreatedAt  time.Time  `json:"created_at"`
	RevokedAt  *time.Time `json:"revoked_at"`
}
