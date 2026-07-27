import { defineStore } from 'pinia'
import { SearchEntries, GetEntry, ExportEntriesJSON } from '../../wailsjs/go/main/App'
import { api, db } from '../../wailsjs/go/models'

export interface EntryFiltersState {
  agentId: string
  source: string
  status: string
  subject: string
  query: string
  minConfidence: number | null
  maxConfidence: number | null
}

function emptyFilters(): EntryFiltersState {
  return {
    agentId: '',
    source: '',
    status: '',
    subject: '',
    query: '',
    minConfidence: null,
    maxConfidence: null,
  }
}

function toApiFilters(f: EntryFiltersState): api.EntryFilters {
  return api.EntryFilters.createFrom({
    AgentID: f.agentId,
    Source: f.source,
    Status: f.status,
    Subject: f.subject,
    Query: f.query,
    MinConfidence: f.minConfidence ?? undefined,
    MaxConfidence: f.maxConfidence ?? undefined,
  })
}

export const useEntriesStore = defineStore('entries', {
  state: () => ({
    filters: emptyFilters(),
    entries: [] as db.LedgerEntry[],
    total: 0,
    limit: 50,
    offset: 0,
    loading: false,
    error: null as string | null,
    selected: null as db.LedgerEntry | null,
  }),
  actions: {
    async search() {
      this.loading = true
      this.error = null
      try {
        const result = await SearchEntries(toApiFilters(this.filters), this.limit, this.offset)
        this.entries = result.entries ?? []
        this.total = result.total
      } catch (e) {
        this.error = String(e)
      } finally {
        this.loading = false
      }
    },
    async selectEntry(id: string) {
      try {
        this.selected = await GetEntry(id)
      } catch (e) {
        this.error = String(e)
      }
    },
    clearSelection() {
      this.selected = null
    },
    resetFilters() {
      this.filters = emptyFilters()
      this.offset = 0
    },
    nextPage() {
      if (this.offset + this.limit < this.total) {
        this.offset += this.limit
        this.search()
      }
    },
    prevPage() {
      if (this.offset > 0) {
        this.offset = Math.max(0, this.offset - this.limit)
        this.search()
      }
    },
    async exportJSON(): Promise<db.LedgerEntry[]> {
      return ExportEntriesJSON(toApiFilters(this.filters))
    },
  },
})
