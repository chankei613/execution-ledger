import { describe, it, expect } from 'vitest'
import { statusRole, statusIcon, STATUS_COLORS } from '@/statusColors'

describe('statusColors', () => {
  it('maps known statuses to expected roles', () => {
    expect(statusRole('success')).toBe('good')
    expect(statusRole('partial_success')).toBe('warning')
    expect(statusRole('blocked')).toBe('serious')
    expect(statusRole('timed_out')).toBe('serious')
    expect(statusRole('failed')).toBe('critical')
    expect(statusRole('token_budget_exceeded')).toBe('critical')
    expect(statusRole('generated')).toBe('neutral')
  })

  it('falls back to neutral for unknown statuses', () => {
    expect(statusRole('something_else')).toBe('neutral')
  })

  it('every role has both light and dark colors defined', () => {
    for (const role of Object.values(STATUS_COLORS)) {
      expect(role.light).toMatch(/^#[0-9a-f]{6}$/i)
      expect(role.dark).toMatch(/^#[0-9a-f]{6}$/i)
    }
  })

  it('returns a non-empty icon for every status', () => {
    for (const s of ['success', 'partial_success', 'failed', 'blocked', 'timed_out', 'token_budget_exceeded', 'generated']) {
      expect(statusIcon(s).length).toBeGreaterThan(0)
    }
  })
})
