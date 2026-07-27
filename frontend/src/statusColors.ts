// ステータス色は dataviz スキルの検証済みパレット（status palette、固定・非テーマ）をそのまま使う。
// 色だけで意味を伝えない: 呼び出し側は必ずアイコン+ラベルとセットで使うこと。
export type StatusRole = 'good' | 'warning' | 'serious' | 'critical' | 'neutral'

export const STATUS_COLORS: Record<StatusRole, { light: string; dark: string }> = {
  good: { light: '#0ca30c', dark: '#0ca30c' },
  warning: { light: '#fab219', dark: '#fab219' },
  serious: { light: '#ec835a', dark: '#ec835a' },
  critical: { light: '#d03b3b', dark: '#d03b3b' },
  neutral: { light: '#8a8a86', dark: '#9a9a95' },
}

const ROLE_BY_STATUS: Record<string, StatusRole> = {
  success: 'good',
  partial_success: 'warning',
  blocked: 'serious',
  timed_out: 'serious',
  failed: 'critical',
  token_budget_exceeded: 'critical',
  generated: 'neutral',
}

const ICON_BY_ROLE: Record<StatusRole, string> = {
  good: '✓',
  warning: '◐',
  serious: '!',
  critical: '✕',
  neutral: '·',
}

export function statusRole(status: string): StatusRole {
  return ROLE_BY_STATUS[status] ?? 'neutral'
}

export function statusIcon(status: string): string {
  return ICON_BY_ROLE[statusRole(status)]
}

// カテゴリカル配色（エージェント別内訳など、状態ではなく識別に使う場合）。
// dataviz検証済みの先頭3スロット（all-pairsでも安全）+ それ以降は "その他" に畳む運用。
export const CATEGORICAL_COLORS = ['#2a78d6', '#eb6834', '#1baf7a']
