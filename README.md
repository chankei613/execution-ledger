# Execution Ledger

「AIの行動を説明できるようにする監査ログ」— comet-taskAI ロードマップ Product D。

AIエージェント(comet-taskAI・ai-scheduler・独自スクリプトなど任意のシステム)が何をしたか・なぜそう判断したかを
`confidence` + `decisions[]` + `actions_taken[]` で構造化して記録する。検索・フィルタ・エクスポート可能。

詳細は [docs/spec.md](docs/spec.md) を参照。

## 現在のステータス: Phase 1-2（Ingestion + 検索/集計/エクスポートAPI）完了

- [x] Phase 0: プロジェクト立ち上げ
- [x] Phase 1: データモデル・Ingestion API（追記専用・APIキー認証・ブートストラップ認証）
- [x] Phase 2: 検索・フィルタ・集計・エクスポートAPI
- [ ] Phase 3: Wails + Vue3 UI
- [ ] Phase 4: 仕上げ・署名・配布・LP

## 使い方

```bash
make tidy   # 依存解決
make smoke  # bootstrap鍵発行 → ingest → 検索 → 集計 → エクスポート の一連を確認する自己完結テスト
make run    # :8421 でAPIサーバー起動（SQLite: execution-ledger.db）
```

### APIキー認証

`AgentKey`が0件のときのみ `POST /api/v1/keys` を未認証で許可する（最初の1件を発行するため）。
1件発行された時点で以降は `Authorization: Bearer <key>` が必須になる。

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
