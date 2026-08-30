# Execution Ledger

「AIの行動を説明できるようにする監査ログ」— comet-taskAI ロードマップ Product D。

AIエージェント(comet-taskAI・ai-scheduler・独自スクリプトなど任意のシステム)が何をしたか・なぜそう判断したかを
`confidence` + `decisions[]` + `actions_taken[]` で構造化して記録する。検索・フィルタ・エクスポート可能。

詳細は [docs/spec.md](docs/spec.md) を参照。

## 現在のステータス: v0.1.0 リリース済み

- [x] Phase 0: プロジェクト立ち上げ
- [x] Phase 1: データモデル・Ingestion API（追記専用・APIキー認証・ブートストラップ認証）
- [x] Phase 2: 検索・フィルタ・集計・エクスポートAPI
- [x] Phase 3: Wails + Vue3 UI（台帳ビュー・ダッシュボード・APIキー管理）
- [x] Phase 4: 仕上げ・署名・配布・LP

macOSアプリ（署名・公証済み、Apple Silicon / Intel 共通のUniversalバイナリ）は
[GitHub Releases](https://github.com/chankei613/execution-ledger/releases) から、
ランディングページは https://execution-ledger-virid.vercel.app/ から入手できる。
（アプリ内Help画面は未対応 — 使い方は本READMEを参照）

## 使い方（デスクトップアプリ）

1. [Releases](../../releases) から自分のOS用のビルドをダウンロードして起動する
2. 初回起動時、Settings画面で「Issue key」を押して最初のAPIキーを発行する（**この場でしか表示されないので必ずコピーする**）
3. Settings画面に表示される Ingestion URL（例: `http://localhost:8421/api/v1/entries`）と発行したAPIキーを、記録させたいAIシステム（comet-taskAI・ai-scheduler・独自スクリプトなど）に設定する
4. そのAIシステムが実行結果をPOSTすると、台帳ビューにリアルタイムで並び、ダッシュボードで集計を確認できる

アプリはウインドウを閉じている間もIngestion APIを起動したまま待ち受け続ける。完全に終了するにはSettings画面の「Quit」を使う。

## 使い方（開発・ヘッドレスサーバー）

```bash
make tidy   # 依存解決
make smoke  # bootstrap鍵発行 → ingest → 検索 → 集計 → エクスポート の一連を確認する自己完結テスト
make run    # :8421 でAPIサーバー起動（SQLite: execution-ledger.db、cmd/elserve）
make ui     # frontend/ の vite dev サーバー起動
```

デスクトップアプリとしてビルドするには `wails build`（`wails.json` 参照）。

### APIキー認証

`AgentKey`が0件のときのみ `POST /api/v1/keys` を未認証で許可する（最初の1件を発行するため）。
1件発行された時点で以降は `Authorization: Bearer <key>` が必須になる。
デスクトップアプリのUI自身はネイティブバインディング経由でデータを読むためAPIキーを必要としない。
APIキーが必要なのは外部プロセスからのIngestion（`POST /api/v1/entries`）のみ。

### エントリの記録（イミュータブル）

```bash
curl -X POST localhost:8421/api/v1/entries \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "source": "comet-taskAI",
    "agent_id": "claude-01",
    "subject": "task#4821",
    "status": "success",
    "summary": "Implemented RLS policies",
    "confidence": {"overall": 0.94, "breakdown": {"task_understood":0.95,"execution_complete":0.95,"correctness":0.9,"side_effects_clean":0.98}, "low_confidence_areas":[]},
    "decisions": [{"description":"chose approach X","rationale":"simplest given constraints","alternatives_considered":["Y","Z"]}],
    "actions_taken": [{"tool":"Bash","input_summary":"ran migration","timestamp":"2026-07-28T00:00:00Z"}]
  }'
```

`POST` のみを提供し、UPDATE/DELETE APIは無い（監査ログの信頼性を保つための意図的な設計）。

## API

| メソッド | パス | 用途 |
|---|---|---|
| POST | `/api/v1/keys` | APIキー発行（ブートストラップ時のみ未認証） |
| GET | `/api/v1/keys` | 発行済みキー一覧 |
| DELETE | `/api/v1/keys/{id}` | キー失効 |
| POST | `/api/v1/entries` | エントリ追記（イミュータブル） |
| GET | `/api/v1/entries` | 検索・フィルタ（agent_id/source/status/subject/q/min_confidence/max_confidence/from/to + limit/offset） |
| GET | `/api/v1/entries/{id}` | 単体取得 |
| GET | `/api/v1/entries/export?format=csv\|json` | エクスポート（フィルタ条件を引き継ぐ） |
| GET | `/api/v1/stats` | 集計（ステータス別件数・平均confidence・低信頼度率・エージェント別件数） |

## ディレクトリ構成

```
internal/db/    GORMモデル（LedgerEntry/AgentKey）・SQLite初期化
internal/api/   REST API（keys/entries/stats）+ 認証ミドルウェア
cmd/smoketest/  bootstrap→ingest→search→stats→exportの通しスモークテスト
docs/           設計ドキュメント
```
