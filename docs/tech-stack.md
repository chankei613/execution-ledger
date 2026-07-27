# 技術選定

**決定日:** 2026-07-28
**ステータス:** 確定

---

## 決定

| レイヤー | 採用 | 理由 |
|---|---|---|
| Desktop基盤 | Wails v2 | A/B/K/Cの4製品すべてで実績あり。統合時の摩擦が最小 |
| Backend | Go 1.22+ | 同上 |
| Frontend | Vue 3 + Vite + Pinia | 同上（Nuxtは不採用。SSR不要） |
| Styling | UnoCSS + shadcn-vue | 同上 |
| DB | SQLite + GORM | 同上。監査ログ用途でも個人〜小規模チームの規模ではSQLiteで十分 |
| CI | GitHub Actions | **Go 1.23 / macos-14 を最初から採用**（harness-managerで踏んだ`missing LC_UUID`・macos-13ランナー滞留を回避） |
| 配布 | `.app` / `.exe` シングルバイナリ + コード署名・公証 | mcp-server-managerと同じApple Developer認証情報を流用 |

## 却下した案

- **専用ログ収集基盤（Elasticsearch等）**: 個人〜小規模チーム向けのローカルファーストツールという製品全体の方針に反する。オーバースペック
- **Nuxt SSR**: comet-taskAI本体のPhase 0決定と同じ理由で不要

## 他製品からの流用ポイント

- APIキー認証（Bearer token → SHA-256ハッシュ照合 + ブートストラップ認証）: harness-manager / agent-config-manager / comet-taskAI と同じ実装パターンをそのまま踏襲する
- release.yml: harness-managerのv0.1.1で修正済みの構成（Go 1.23、notarytool用zip提出、windows拡張子なし対策、macos-14）をコピーして使う。同じ轍を踏まない
