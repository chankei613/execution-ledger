import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useEntriesStore } from '@/stores/entries'
import { SearchEntries } from '../../wailsjs/go/main/App'
import { api } from '../../wailsjs/go/models'

describe('entries store error handling', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.mocked(SearchEntries).mockReset()
  })

  it('captures a failed search as store.error and clears loading', async () => {
    vi.mocked(SearchEntries).mockRejectedValueOnce(new Error('network down'))
    const store = useEntriesStore()

    await store.search()

    expect(store.loading).toBe(false)
    expect(store.error).toContain('network down')
  })

  it('clears the previous error on a successful retry', async () => {
    vi.mocked(SearchEntries).mockRejectedValueOnce(new Error('network down'))
    const store = useEntriesStore()
    await store.search()
    expect(store.error).not.toBeNull()

    vi.mocked(SearchEntries).mockResolvedValueOnce(api.SearchResult.createFrom({ entries: [], total: 0 }))
    await store.search()

    expect(store.error).toBeNull()
  })
})
