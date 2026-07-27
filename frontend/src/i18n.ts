import { ref } from 'vue'

export type Locale = 'en' | 'ja'

// localStorage に保存して再起動後も維持する
const saved = window.localStorage.getItem('locale') as Locale | null
const locale = ref<Locale>(saved === 'en' || saved === 'ja' ? saved : 'ja')

const messages: Record<Locale, Record<string, string>> = {
  en: {
    'app.subtitle': 'Execution Ledger',
    'lang.toggle': 'JA',
    'nav.ledger': 'Ledger',
    'nav.dashboard': 'Dashboard',
    'nav.settings': 'Settings',

    'ledger.title': 'Ledger',
    'ledger.filters.agent': 'Agent',
    'ledger.filters.source': 'Source',
    'ledger.filters.status': 'Status',
    'ledger.filters.subject': 'Subject',
    'ledger.filters.query': 'Search summary',
    'ledger.filters.minConfidence': 'Min confidence',
    'ledger.filters.maxConfidence': 'Max confidence',
    'ledger.filters.reset': 'Reset filters',
    'ledger.empty': 'No entries match these filters',
    'ledger.loading': 'Loading…',
    'ledger.export.json': 'Export JSON',
    'ledger.export.csv': 'Export CSV',
    'ledger.page.prev': 'Prev',
    'ledger.page.next': 'Next',
    'ledger.page.of': '{from}–{to} of {total}',
    'ledger.lowConfidence': 'low confidence',

    'detail.title': 'Entry detail',
    'detail.close': 'Close',
    'detail.summary': 'Summary',
    'detail.confidence': 'Confidence',
    'detail.decisions': 'Decisions',
    'detail.actions': 'Actions taken',
    'detail.followUp': 'Follow-up',
    'detail.outputs': 'Outputs',
    'detail.usage': 'Usage',
    'detail.criteria': 'Acceptance criteria',
    'detail.none': 'None recorded',

    'dashboard.title': 'Dashboard',
    'dashboard.total': 'Total entries',
    'dashboard.avgConfidence': 'Avg confidence',
    'dashboard.lowConfidenceRate': 'Low confidence rate',
    'dashboard.byStatus': 'By status',
    'dashboard.byAgent': 'By agent',
    'dashboard.other': 'Other',
    'dashboard.empty': 'No entries yet',

    'settings.title': 'Settings',
    'settings.version': 'Version',
    'settings.ingestion.title': 'Ingestion endpoint',
    'settings.ingestion.desc': 'POST execution results here from any AI system (comet-taskAI, ai-scheduler, custom scripts).',
    'settings.keys.title': 'API keys',
    'settings.keys.name': 'Key name',
    'settings.keys.issue': 'Issue key',
    'settings.keys.issued': 'Key issued — copy it now, it will not be shown again',
    'settings.keys.copy': 'Copy',
    'settings.keys.revoke': 'Revoke',
    'settings.keys.revoked': 'Revoked',
    'settings.keys.empty': 'No keys issued yet. Issue one to allow an AI system to ingest entries.',
    'settings.quit': 'Quit',
    'settings.quit.confirm': 'Quit the app? Ingestion will stop until you reopen it.',

    'status.success': 'success',
    'status.partial_success': 'partial success',
    'status.failed': 'failed',
    'status.blocked': 'blocked',
    'status.timed_out': 'timed out',
    'status.token_budget_exceeded': 'budget exceeded',
    'status.generated': 'generated',
  },
  ja: {
    'app.subtitle': 'Execution Ledger',
    'lang.toggle': 'EN',
    'nav.ledger': '台帳',
    'nav.dashboard': 'ダッシュボード',
    'nav.settings': '設定',

    'ledger.title': '実行台帳',
    'ledger.filters.agent': 'エージェント',
    'ledger.filters.source': 'ソース',
    'ledger.filters.status': 'ステータス',
    'ledger.filters.subject': '対象(subject)',
    'ledger.filters.query': 'サマリ検索',
    'ledger.filters.minConfidence': '信頼度(最小)',
    'ledger.filters.maxConfidence': '信頼度(最大)',
    'ledger.filters.reset': 'フィルタをリセット',
    'ledger.empty': '条件に一致するエントリがありません',
    'ledger.loading': '読み込み中…',
    'ledger.export.json': 'JSONエクスポート',
    'ledger.export.csv': 'CSVエクスポート',
    'ledger.page.prev': '前へ',
    'ledger.page.next': '次へ',
    'ledger.page.of': '{total}件中 {from}–{to}',
    'ledger.lowConfidence': '低信頼度',

    'detail.title': 'エントリ詳細',
    'detail.close': '閉じる',
    'detail.summary': 'サマリ',
    'detail.confidence': '信頼度',
    'detail.decisions': '判断根拠',
    'detail.actions': '実行アクション',
    'detail.followUp': 'フォローアップ',
    'detail.outputs': '出力',
    'detail.usage': 'トークン使用量',
    'detail.criteria': '達成基準',
    'detail.none': '記録なし',

    'dashboard.title': 'ダッシュボード',
    'dashboard.total': '総エントリ数',
    'dashboard.avgConfidence': '平均信頼度',
    'dashboard.lowConfidenceRate': '低信頼度率',
    'dashboard.byStatus': 'ステータス別',
    'dashboard.byAgent': 'エージェント別',
    'dashboard.other': 'その他',
    'dashboard.empty': 'まだエントリがありません',

    'settings.title': '設定',
    'settings.version': 'バージョン',
    'settings.ingestion.title': 'Ingestionエンドポイント',
    'settings.ingestion.desc': '任意のAIシステム(comet-taskAI・ai-scheduler・独自スクリプト)からここへ実行結果をPOSTしてください。',
    'settings.keys.title': 'APIキー',
    'settings.keys.name': 'キー名',
    'settings.keys.issue': 'キーを発行',
    'settings.keys.issued': 'キーを発行しました — この場では二度と表示されないので今すぐコピーしてください',
    'settings.keys.copy': 'コピー',
    'settings.keys.revoke': '失効',
    'settings.keys.revoked': '失効済み',
    'settings.keys.empty': 'まだキーがありません。AIシステムが記録できるようキーを発行してください。',
    'settings.quit': '終了',
    'settings.quit.confirm': 'アプリを終了しますか？再度開くまでIngestionは停止します。',

    'status.success': '成功',
    'status.partial_success': '部分的成功',
    'status.failed': '失敗',
    'status.blocked': 'ブロック中',
    'status.timed_out': 'タイムアウト',
    'status.token_budget_exceeded': '予算超過',
    'status.generated': '生成済み',
  },
}

export function useI18n() {
  function t(key: string, params?: Record<string, string | number>): string {
    let msg = messages[locale.value][key] ?? messages.en[key] ?? key
    if (params) {
      for (const [k, v] of Object.entries(params)) {
        msg = msg.replace(`{${k}}`, String(v))
      }
    }
    return msg
  }

  function toggleLocale() {
    locale.value = locale.value === 'en' ? 'ja' : 'en'
    window.localStorage.setItem('locale', locale.value)
  }

  return { t, locale, toggleLocale }
}
