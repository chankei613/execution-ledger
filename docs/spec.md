# Execution Ledger — 仕様書

> 作成: 2026-07-28
> ステータス: 設計フェーズ

---

## 1. 製品概要

**「AIの行動を説明できるようにする監査ログ」** — AIエージェントが何をしたか・なぜそう判断したかを
構造化して記録し、検索・フィルタ・エクスポートできるデスクトップアプリ。

### 解決する問題

AIを業務に使い始めた企業・チームは以下の課題を抱える：

- AIが何を実行したか、後から追跡できない（チャット履歴はaudit logではない）
- 「なぜその判断をしたか」という根拠（rationale）・検討した代替案が残らない
- confidence（自信度）が低い実行が埋もれてしまい、レビューが後手に回る
- コンプライアンス・監査で「AIの行動を説明してください」と言われても答えられない

### ソリューション

任意のAIシステム（comet-taskAI、ai-scheduler、独自スクリプト等）から実行結果を
`POST /api/v1/entries` に送るだけで、イミュータブルな監査ログとして蓄積される。
人間はタイムラインビューで検索・フィルタし、低信頼度のエントリをすぐに見つけられる。

---

## 2. コアコンセプト

### LedgerEntry（受信する1件の実行記録）

comet-taskAI の `ExecutionResult`（schema/types.ts）をベースに、
単体製品として「どのシステムからでも受けられる」ようフィールドを追加する。

```typescript
interface LedgerEntry {
  id: string
  received_at: string          // サーバー側で付与。イミュータブルの起点

  // 送信元の識別（自由記述。comet-taskAIのAgentConfig.idに限定しない）
  source: string                // "comet-taskAI" | "ai-scheduler" | "manual" | 任意の文字列
  agent_id: string               // どのエージェントが実行したか
  subject: string                // 対象タスクの自由記述ID（例: "task#4821"）

  status: "success" | "partial_success" | "failed" | "blocked" | "timed_out" | "token_budget_exceeded" | "generated"
  summary: string                // 次のAIまたは人間が読む前提の1段落

  criteria_results: CriterionResult[]
  outputs: Record<string, unknown>

  confidence: {
    overall: number               // 0.0–1.0
    breakdown: {
      task_understood: number
      execution_complete: number
      correctness: number
      side_effects_clean: number
    }
    low_confidence_areas: string[]
  }

  decisions: Decision[]           // { description, rationale, alternatives_considered }
  actions_taken: Action[]         // { tool, input_summary, timestamp }
  follow_up: FollowUp[]           // { description, suggested_task }

  usage: {
    input_tokens: number
    output_tokens: number
    mcp_calls_by_server: Record<string, number>
  }
}
```

**イミュータブル原則**: `POST /api/v1/entries` のみを提供し、UPDATE/DELETE APIは提供しない。
訂正が必要な場合は新しいエントリを追加し、`subject` で紐付けて追跡する運用とする
（concept.md 6章「AIの全行動はイミュータブルなレコードとして記録。チャット履歴ではなくAuditログ」）。

---

## 3. 機能一覧

### Phase 1-2 (MVP: Ingestion + 検索)

| 機能 | 説明 |
|------|------|
| エントリ追記API | `POST /api/v1/entries`（APIキー認証、追記専用） |
| エントリ検索 | agent_id・source・status・confidence範囲・期間・全文検索 + ページネーション |
| 集計API | ステータス別件数・平均confidence・低信頼度率・エージェント別件数 |
| エクスポート | CSV/JSON（フィルタ条件を引き継ぐ） |
| APIキー管理 | ブートストラップ発行（0件時のみ未認証）・失効 |

### Phase 3 (UI)

| 機能 | 説明 |
|------|------|
| タイムラインビュー | 監査ログ形式の一覧（時刻・エージェント・ステータス・confidence を1行で） |
| フィルタサイドバー | エージェント・source・ステータス・confidence範囲・期間・全文検索 |
| エントリ詳細ドロワー | decisions/actions_taken/follow_up/outputs/usage の全表示 |
| ダッシュボード | 統計カード・ステータス別グラフ・低信頼度率の推移 |
| エクスポートボタン | CSV/JSON ダウンロード |

### Phase 4 (拡張候補・本リリースの対象外)

| 機能 | 説明 |
|------|------|
| Decision Queue 連携 | `status: "blocked"` のエントリを承認待ちインボックスとして扱う（将来のG. AI Decision Reviewerの前提） |
| マスキング | outputs/actions_takenの機密情報を正規表現でマスクしてエクスポート |

---

## 4. タイムラインビューの表示形式

concept.md 6章の例に準拠する:

```
[14:32:07] agent-02  WRITE  task#4821.status → "in_progress"   confidence: 0.94
[14:32:05] agent-02  READ   task#4821
[14:31:58] agent-01  BLOCK  tool:github-mcp  reason: rate_limit → escalating
```

`actions_taken[]` の `tool` + `input_summary` を1行に要約して表示し、
confidenceが閾値（既定0.6）未満のエントリは行全体を警告色でハイライトする。

---

## 5. UX フロー

```
起動
 └── タイムラインビュー（LedgerView, 既定画面）
      ├── フィルタ変更 → 即座に再検索（デバウンス）
      ├── 行クリック → 詳細ドロワー（decisions/actions/outputs/usage）
      ├── エクスポートボタン → CSV/JSON ダウンロード
      └── ダッシュボードタブ → 統計・グラフ

APIキー管理
 └── SettingsView → キー発行（ブートストラップ後は認証必須）・失効
```

---

## 6. データストア

SQLite（ローカル、`~/.execution-ledger/ledger.db`）

```sql
ledger_entries (
  id, received_at, source, agent_id, subject, status, summary,
  criteria_results JSON, outputs JSON,
  confidence_overall, confidence_breakdown JSON, low_confidence_areas JSON,
  decisions JSON, actions_taken JSON, follow_up JSON,
  usage JSON
)
agent_keys (id, name, api_key_hash, created_at, revoked_at)
```

`criteria_results` / `outputs` / `confidence_breakdown` / `low_confidence_areas` /
`decisions` / `actions_taken` / `follow_up` / `usage` はJSONカラムとしてシリアライズする
（harness-manager / agent-config-manager と同じGORM `serializer:json` パターン）。
