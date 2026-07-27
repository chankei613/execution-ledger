<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useEntriesStore } from '@/stores/entries'
import { useI18n } from '@/i18n'
import { statusRole, statusIcon, STATUS_COLORS } from '@/statusColors'
import EntryDetailDrawer from '@/components/EntryDetailDrawer.vue'

const { t } = useI18n()
const store = useEntriesStore()

let debounceTimer: ReturnType<typeof setTimeout> | null = null
function onFilterChange() {
  store.offset = 0
  if (debounceTimer) clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => store.search(), 250)
}

const statusOptions = [
  'success', 'partial_success', 'failed', 'blocked', 'timed_out', 'token_budget_exceeded', 'generated',
]

function formatTime(v: any): string {
  const d = new Date(v)
  if (isNaN(d.getTime())) return String(v)
  return d.toLocaleString(undefined, { month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit', second: '2-digit' })
}

function summarizeAction(entry: any): string {
  const actions = entry.actions_taken ?? []
  if (actions.length === 0) return entry.status === 'blocked' ? 'BLOCK' : 'RUN'
  const a = actions[actions.length - 1]
  return `${a.tool}: ${a.input_summary}`.slice(0, 80)
}

async function exportCSV() {
  const params = new URLSearchParams()
  const f = store.filters
  if (f.agentId) params.set('agent_id', f.agentId)
  if (f.source) params.set('source', f.source)
  if (f.status) params.set('status', f.status)
  if (f.subject) params.set('subject', f.subject)
  if (f.query) params.set('q', f.query)
  if (f.minConfidence != null) params.set('min_confidence', String(f.minConfidence))
  if (f.maxConfidence != null) params.set('max_confidence', String(f.maxConfidence))
  params.set('format', 'csv')
  // CSVはローカルAPIキーが要るためJSON経由で取得しBlobダウンロードする
  const entries = await store.exportJSON()
  const rows = [['id', 'received_at', 'source', 'agent_id', 'subject', 'status', 'confidence_overall', 'summary']]
  for (const e of entries) {
    rows.push([e.id, String(e.received_at), e.source, e.agent_id, e.subject, e.status, String(e.confidence_overall), e.summary])
  }
  const csv = rows.map(r => r.map(v => `"${String(v).replace(/"/g, '""')}"`).join(',')).join('\n')
  downloadBlob(csv, 'text/csv', 'ledger-export.csv')
}

async function exportJSON() {
  const entries = await store.exportJSON()
  downloadBlob(JSON.stringify(entries, null, 2), 'application/json', 'ledger-export.json')
}

function downloadBlob(content: string, mime: string, filename: string) {
  const blob = new Blob([content], { type: mime })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = filename
  a.click()
  URL.revokeObjectURL(url)
}

onMounted(() => store.search())
</script>

<template>
  <div class="flex h-full">
    <!-- フィルタサイドバー -->
    <aside class="w-64 border-r border-border p-4 space-y-3 overflow-y-auto shrink-0">
      <h2 class="text-sm font-semibold">{{ t('ledger.title') }}</h2>

      <div class="space-y-1">
        <label class="text-xs text-muted-foreground">{{ t('ledger.filters.agent') }}</label>
        <input v-model="store.filters.agentId" @input="onFilterChange" class="w-full text-sm border border-border rounded px-2 py-1" />
      </div>
      <div class="space-y-1">
        <label class="text-xs text-muted-foreground">{{ t('ledger.filters.source') }}</label>
        <input v-model="store.filters.source" @input="onFilterChange" class="w-full text-sm border border-border rounded px-2 py-1" />
      </div>
      <div class="space-y-1">
        <label class="text-xs text-muted-foreground">{{ t('ledger.filters.status') }}</label>
        <select v-model="store.filters.status" @change="onFilterChange" class="w-full text-sm border border-border rounded px-2 py-1">
          <option value="">—</option>
          <option v-for="s in statusOptions" :key="s" :value="s">{{ t('status.' + s) }}</option>
        </select>
      </div>
      <div class="space-y-1">
        <label class="text-xs text-muted-foreground">{{ t('ledger.filters.subject') }}</label>
        <input v-model="store.filters.subject" @input="onFilterChange" class="w-full text-sm border border-border rounded px-2 py-1" />
      </div>
      <div class="space-y-1">
        <label class="text-xs text-muted-foreground">{{ t('ledger.filters.query') }}</label>
        <input v-model="store.filters.query" @input="onFilterChange" class="w-full text-sm border border-border rounded px-2 py-1" />
      </div>
      <div class="grid grid-cols-2 gap-2">
        <div class="space-y-1">
          <label class="text-xs text-muted-foreground">{{ t('ledger.filters.minConfidence') }}</label>
          <input type="number" min="0" max="1" step="0.1" v-model.number="store.filters.minConfidence" @input="onFilterChange" class="w-full text-sm border border-border rounded px-2 py-1" />
        </div>
        <div class="space-y-1">
          <label class="text-xs text-muted-foreground">{{ t('ledger.filters.maxConfidence') }}</label>
          <input type="number" min="0" max="1" step="0.1" v-model.number="store.filters.maxConfidence" @input="onFilterChange" class="w-full text-sm border border-border rounded px-2 py-1" />
        </div>
      </div>
      <button @click="store.resetFilters(); store.search()" class="text-xs text-muted-foreground hover:text-foreground underline">
        {{ t('ledger.filters.reset') }}
      </button>

      <div class="pt-3 border-t border-border space-y-1">
        <button @click="exportJSON" class="w-full text-xs text-left px-2 py-1.5 rounded border border-border hover:bg-gray-50">{{ t('ledger.export.json') }}</button>
        <button @click="exportCSV" class="w-full text-xs text-left px-2 py-1.5 rounded border border-border hover:bg-gray-50">{{ t('ledger.export.csv') }}</button>
      </div>
    </aside>

    <!-- タイムライン -->
    <main class="flex-1 overflow-y-auto p-4">
      <div v-if="store.loading" class="text-sm text-muted-foreground">{{ t('ledger.loading') }}</div>
      <div v-else-if="store.entries.length === 0" class="text-sm text-muted-foreground">{{ t('ledger.empty') }}</div>

      <div v-else class="font-mono text-[13px] leading-relaxed space-y-0.5">
        <div
          v-for="entry in store.entries"
          :key="entry.id"
          @click="store.selectEntry(entry.id)"
          class="flex items-center gap-2 px-2 py-1.5 rounded cursor-pointer hover:bg-gray-50 border-l-2"
          :class="entry.confidence_overall < 0.6 ? 'border-l-[--critical]' : 'border-l-transparent'"
          :style="{ '--critical': STATUS_COLORS.critical.light }"
        >
          <span class="text-muted-foreground shrink-0">[{{ formatTime(entry.received_at) }}]</span>
          <span class="shrink-0 font-semibold" :title="t('status.' + entry.status)">
            {{ entry.agent_id }}
          </span>
          <span
            class="shrink-0 inline-flex items-center gap-1 px-1.5 rounded text-white text-[11px]"
            :style="{ background: STATUS_COLORS[statusRole(entry.status)].light }"
          >
            <span aria-hidden="true">{{ statusIcon(entry.status) }}</span>
            {{ t('status.' + entry.status) }}
          </span>
          <span class="truncate flex-1">{{ entry.subject }} — {{ summarizeAction(entry) }}</span>
          <span
            class="shrink-0 tabular-nums"
            :class="entry.confidence_overall < 0.6 ? 'font-semibold' : 'text-muted-foreground'"
            :style="entry.confidence_overall < 0.6 ? { color: STATUS_COLORS.critical.light } : {}"
          >
            confidence: {{ entry.confidence_overall.toFixed(2) }}
          </span>
        </div>
      </div>

      <div v-if="store.total > 0" class="flex items-center justify-between mt-4 text-xs text-muted-foreground">
        <button :disabled="store.offset === 0" @click="store.prevPage" class="disabled:opacity-30 underline">{{ t('ledger.page.prev') }}</button>
        <span>{{ t('ledger.page.of', { from: store.offset + 1, to: Math.min(store.offset + store.limit, store.total), total: store.total }) }}</span>
        <button :disabled="store.offset + store.limit >= store.total" @click="store.nextPage" class="disabled:opacity-30 underline">{{ t('ledger.page.next') }}</button>
      </div>
    </main>

    <EntryDetailDrawer v-if="store.selected" :entry="store.selected" @close="store.clearSelection" />
  </div>
</template>
