import { defineStore } from 'pinia'
import { GetStats } from '../../wailsjs/go/main/App'
import { api } from '../../wailsjs/go/models'

export const useStatsStore = defineStore('stats', {
  state: () => ({
    total: 0,
    byStatus: {} as Record<string, number>,
    avgConfidence: 0,
    lowConfidenceRate: 0,
    byAgent: {} as Record<string, number>,
    loading: false,
    error: null as string | null,
  }),
  actions: {
    async load() {
      this.loading = true
      this.error = null
      try {
        const result = await GetStats(api.EntryFilters.createFrom({}))
        this.total = result.total
        this.byStatus = result.by_status ?? {}
        this.avgConfidence = result.avg_confidence
        this.lowConfidenceRate = result.low_confidence_rate
        this.byAgent = result.by_agent ?? {}
      } catch (e) {
        this.error = String(e)
      } finally {
        this.loading = false
      }
    },
  },
})
