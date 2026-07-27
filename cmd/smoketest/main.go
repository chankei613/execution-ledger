// cmd/smoketest はExecution LedgerのAPIを一時DBで自前起動し、
// ブートストラップ鍵発行 → ingest → 検索 → 集計 → エクスポート の一連が
// 通しで動くことを確認する。
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/chankei613/execution-ledger/internal/api"
	"github.com/chankei613/execution-ledger/internal/db"
)

func main() {
	dbPath := "smoketest.db"
	os.Remove(dbPath)
	defer os.Remove(dbPath)

	conn, err := db.Init(dbPath)
	if err != nil {
		log.Fatalf("db init: %v", err)
	}

	srv := httptest.NewServer(api.NewRouter(conn))
	defer srv.Close()
	client := srv.Client()

	// 1. bootstrap key issuance (unauthenticated, 0件のときのみ許可)
	keyBody, _ := json.Marshal(map[string]string{"name": "smoketest-agent"})
	var keyResp struct {
		APIKey string `json:"api_key"`
	}
	postJSON(client, srv.URL+"/api/v1/keys", "", keyBody, &keyResp)
	if keyResp.APIKey == "" {
		log.Fatal("FAIL: bootstrap key issuance returned empty key")
	}
	fmt.Println("PASS: bootstrap key issued")

	apiKey := keyResp.APIKey

	// 2回目のキー発行は認証必須になっているはず（未認証だと401）
	req, _ := http.NewRequest(http.MethodPost, srv.URL+"/api/v1/keys", bytes.NewReader(keyBody))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := client.Do(req)
	if resp.StatusCode != http.StatusUnauthorized {
		log.Fatalf("FAIL: expected 401 for unauthenticated 2nd key issuance, got %d", resp.StatusCode)
	}
	fmt.Println("PASS: bootstrap closes after first key (2nd unauthenticated request -> 401)")

	// 2. ingest
	entryBody, _ := json.Marshal(map[string]any{
		"source":   "comet-taskAI",
		"agent_id": "claude-01",
		"subject":  "task#4821",
		"status":   "success",
		"summary":  "Implemented RLS policies",
		"confidence": map[string]any{
			"overall": 0.42,
			"breakdown": map[string]float64{
				"task_understood": 0.9, "execution_complete": 0.9, "correctness": 0.3, "side_effects_clean": 0.6,
			},
			"low_confidence_areas": []string{"correctness"},
		},
		"decisions": []map[string]any{
			{"description": "chose approach X", "rationale": "simplest given constraints", "alternatives_considered": []string{"Y", "Z"}},
		},
		"actions_taken": []map[string]any{
			{"tool": "Bash", "input_summary": "ran migration", "timestamp": "2026-07-28T00:00:00Z"},
		},
	})
	var ingested db.LedgerEntry
	postJSON(client, srv.URL+"/api/v1/entries", apiKey, entryBody, &ingested)
	if ingested.ID == "" || ingested.Status != db.StatusSuccess {
		log.Fatalf("FAIL: ingest returned unexpected entry: %+v", ingested)
	}
	fmt.Printf("PASS: ingested entry %s (confidence=%.2f)\n", ingested.ID, ingested.ConfidenceOverall)

	// 3. search with filter (should find it, low confidence filter)
	var searchResp struct {
		Entries []db.LedgerEntry `json:"entries"`
		Total   int64            `json:"total"`
	}
	getJSON(client, srv.URL+"/api/v1/entries?max_confidence=0.5", apiKey, &searchResp)
	if searchResp.Total != 1 {
		log.Fatalf("FAIL: expected 1 low-confidence entry, got %d", searchResp.Total)
	}
	fmt.Println("PASS: filtered search found the low-confidence entry")

	// 4. stats
	var statsResp struct {
		Total             int64   `json:"total"`
		LowConfidenceRate float64 `json:"low_confidence_rate"`
	}
	getJSON(client, srv.URL+"/api/v1/stats", apiKey, &statsResp)
	if statsResp.Total != 1 || statsResp.LowConfidenceRate != 1.0 {
		log.Fatalf("FAIL: unexpected stats: %+v", statsResp)
	}
	fmt.Println("PASS: stats reflects low_confidence_rate=1.0")

	// 5. export csv
	req, _ = http.NewRequest(http.MethodGet, srv.URL+"/api/v1/entries/export?format=csv", nil)
	req.Header.Set("Authorization", "Bearer "+apiKey)
	resp, err = client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Fatalf("FAIL: csv export failed")
	}
	fmt.Println("PASS: CSV export OK")

	fmt.Println("SMOKE TEST OK")
}

func postJSON(client *http.Client, url, apiKey string, body []byte, out interface{}) {
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		log.Fatalf("FAIL: POST %s -> %s", url, resp.Status)
	}
	json.NewDecoder(resp.Body).Decode(out)
}

func getJSON(client *http.Client, url, apiKey string, out interface{}) {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Authorization", "Bearer "+apiKey)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		log.Fatalf("FAIL: GET %s -> %s", url, resp.Status)
	}
	json.NewDecoder(resp.Body).Decode(out)
}
