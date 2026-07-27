# 開発計画

**予測期間:** 2〜3週間相当（ロードマップ見積もり）。過去実績（A〜C）に基づき短縮を狙う。

| Phase | 内容 |
|---|---|
| Phase 0 | プロジェクト立ち上げ（Go初期化・docs・GitHub repo） |
| Phase 1 | データモデル・Ingestion API（追記専用、APIキー認証） |
| Phase 2 | 検索・フィルタ・集計・エクスポートAPI |
| Phase 3 | Wails + Vue3 UI（タイムライン・フィルタ・詳細・ダッシュボード） |
| Phase 4 | 仕上げ・署名・配布・LP |

## 優先順位の根拠

concept.md 9章の指示通り、Execution LedgerはUI・分析より前に「記録が正しく蓄積される」ことが最重要。
Phase 1（Ingestion）とPhase 2（検索）をUIより先に固め、動くAPIができた時点で
`curl`ベースの手動テストとcomet-taskAIからの実データ投入で正しさを検証してからUIに進む。
