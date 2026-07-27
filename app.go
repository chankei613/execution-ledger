package main

import (
	"context"
	"net/http"
	"os"
	"path/filepath"

	"github.com/chankei613/execution-ledger/internal/api"
	"github.com/chankei613/execution-ledger/internal/db"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// apiAddr はIngestion APIの待ち受けアドレス。comet-taskAI・ai-scheduler等の
// 外部プロセスがアプリ起動中いつでもPOSTできるよう、ウインドウの表示/非表示に
// 関わらずこのHTTPサーバーは動き続ける（UI自体はこのHTTPを経由せず、下記の
// ネイティブバインディング経由で同じ *api.Server を直接呼ぶ）。
const apiAddr = "127.0.0.1:8421"

// App はWailsのバインディング。実処理は internal/api.Server が持っており、
// ここはWails固有の初期化・エラー通知と、UI向けのネイティブバインディングだけを担当する。
// 同じ Server を cmd/elserve のHTTP APIも使っているので、UIとAPIで挙動がズレない。
type App struct {
	ctx    context.Context
	server *api.Server
	srv    *http.Server
	ready  bool
}

func NewApp() *App { return &App{} }

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	dataDir := appDataDir()
	if err := os.MkdirAll(dataDir, 0o700); err != nil {
		runtime.LogErrorf(ctx, "data dir error: %s", err)
		return
	}

	conn, err := db.Init(filepath.Join(dataDir, "execution-ledger.db"))
	if err != nil {
		runtime.LogErrorf(ctx, "db init error: %s", err)
		return
	}
	a.server = api.New(conn)

	a.srv = &http.Server{Addr: apiAddr, Handler: a.server.Router()}
	go func() {
		if err := a.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			runtime.LogErrorf(ctx, "api server error: %s", err)
		}
	}()

	a.ready = true
	runtime.LogInfof(ctx, "Execution Ledger ready (ingestion api: http://%s, data: %s)", apiAddr, dataDir)
}

func (a *App) shutdown(ctx context.Context) {
	if a.srv != nil {
		_ = a.srv.Close()
	}
	if a.server != nil {
		if sqlDB, err := a.server.DB.DB(); err == nil {
			sqlDB.Close()
		}
	}
}

var errNotReady = &notReadyError{}

type notReadyError struct{}

func (e *notReadyError) Error() string { return "app not ready — check startup logs" }

// ─── フロントエンドへ公開するメソッド ──────────────────────────────────────────

// GetAppVersion はアプリのバージョン文字列を返す。
func (a *App) GetAppVersion() string {
	return AppVersion
}

// GetIngestionURL は外部プロセスがPOSTする先のURLを返す（Settings画面に表示する）。
func (a *App) GetIngestionURL() string {
	return "http://" + apiAddr + "/api/v1/entries"
}

func (a *App) SearchEntries(filters api.EntryFilters, limit int, offset int) (api.SearchResult, error) {
	if !a.ready {
		return api.SearchResult{}, errNotReady
	}
	return a.server.SearchEntries(filters, limit, offset)
}

func (a *App) GetEntry(id string) (db.LedgerEntry, error) {
	if !a.ready {
		return db.LedgerEntry{}, errNotReady
	}
	return a.server.GetEntry(id)
}

func (a *App) GetStats(filters api.EntryFilters) (api.StatsResult, error) {
	if !a.ready {
		return api.StatsResult{}, errNotReady
	}
	return a.server.Stats(filters)
}

func (a *App) ExportEntriesJSON(filters api.EntryFilters) ([]db.LedgerEntry, error) {
	if !a.ready {
		return nil, errNotReady
	}
	return a.server.ExportEntries(filters)
}

func (a *App) ListKeys() ([]db.AgentKey, error) {
	if !a.ready {
		return nil, errNotReady
	}
	return a.server.ListKeys()
}

func (a *App) IssueKey(name string) (api.IssueKeyResult, error) {
	if !a.ready {
		return api.IssueKeyResult{}, errNotReady
	}
	return a.server.IssueKey(name)
}

func (a *App) RevokeKey(id string) error {
	if !a.ready {
		return errNotReady
	}
	return a.server.RevokeKey(id)
}

// Quit はアプリを完全終了する（Settings 画面から呼ぶ）。
func (a *App) Quit() {
	runtime.Quit(a.ctx)
}

func appDataDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return filepath.Join(home, ".execution-ledger")
}
