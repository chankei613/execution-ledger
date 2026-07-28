<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useStatsStore } from '@/stores/stats'
import { useI18n } from '@/i18n'
import { statusRole, STATUS_COLORS, CATEGORICAL_COLORS } from '@/statusColors'

const { t } = useI18n()
const store = useStatsStore()

onMounted(() => store.load())

const statusBars = computed(() => {
  const entries = Object.entries(store.byStatus)
  const max = Math.max(1, ...entries.map(([, v]) => v))
  return entries
    .sort((a, b) => b[1] - a[1])
    .map(([status, count]) => ({
      status,
      count,
      pct: (count / max) * 100,
      color: STATUS_COLORS[statusRole(status)].light,
    }))
})

// カテゴリカル配色は上位3件のみ割り当て、残りは「その他」に畳む（all-pairs安全な範囲を超えないため）
const agentBars = computed(() => {
  const entries = Object.entries(store.byAgent).sort((a, b) => b[1] - a[1])
  const top = entries.slice(0, 3)
  const restTotal = entries.slice(3).reduce((sum, [, v]) => sum + v, 0)
  const rows = top.map(([agent, count], i) => ({ agent, count, color: CATEGORICAL_COLORS[i] }))
  if (restTotal > 0) rows.push({ agent: t('dashboard.other'), count: restTotal, color: '#8a8a86' })
  const max = Math.max(1, ...rows.map(r => r.count))
  return rows.map(r => ({ ...r, pct: (r.count / max) * 100 }))
})
</script>

<template>
  <div class="p-6 space-y-6 overflow-y-auto h-full">
    <h2 class="text-sm font-semibold">{{ t('dashboard.title') }}</h2>

    <div v-if="store.error" class="text-sm border rounded px-3 py-2" :style="{ borderColor: STATUS_COLORS.critical.light, color: STATUS_COLORS.critical.light }">
      {{ t('error.prefix') }}{{ store.error }}
      <button @click="store.load" class="ml-2 underline">{{ t('error.retry') }}</button>
    </div>
    <div v-else-if="store.loading" class="text-sm text-muted-foreground">{{ t('ledger.loading') }}</div>
    <div v-else-if="store.total === 0" class="text-sm text-muted-foreground">{{ t('dashboard.empty') }}</div>

    <template v-else>
      <!-- 統計カード（単一の見出し数値には図表ではなくstat tileが適切） -->
      <div class="grid grid-cols-3 gap-3">
        <div class="border border-border rounded-lg p-4">
          <div class="text-xs text-muted-foreground">{{ t('dashboard.total') }}</div>
          <div class="text-2xl font-semibold tabular-nums mt-1">{{ store.total }}</div>
        </div>
        <div class="border border-border rounded-lg p-4">
          <div class="text-xs text-muted-foreground">{{ t('dashboard.avgConfidence') }}</div>
          <div class="text-2xl font-semibold tabular-nums mt-1">{{ store.avgConfidence.toFixed(2) }}</div>
        </div>
        <div class="border border-border rounded-lg p-4">
          <div class="text-xs text-muted-foreground">{{ t('dashboard.lowConfidenceRate') }}</div>
          <div
            class="text-2xl font-semibold tabular-nums mt-1"
            :style="store.lowConfidenceRate > 0.2 ? { color: STATUS_COLORS.critical.light } : {}"
          >
            {{ (store.lowConfidenceRate * 100).toFixed(0) }}%
          </div>
        </div>
      </div>

      <!-- ステータス別（状態なのでstatus paletteをアイコン+ラベル付きで使う） -->
      <div>
        <h3 class="text-xs font-semibold text-muted-foreground mb-2">{{ t('dashboard.byStatus') }}</h3>
        <div class="space-y-1.5">
          <div v-for="row in statusBars" :key="row.status" class="flex items-center gap-2 text-sm">
            <span class="w-32 shrink-0 text-xs">{{ t('status.' + row.status) }}</span>
            <div class="flex-1 h-2 bg-gray-100 rounded-full overflow-hidden">
              <div class="h-full rounded-full" :style="{ width: row.pct + '%', background: row.color }" />
            </div>
            <span class="w-8 text-right text-xs tabular-nums text-muted-foreground">{{ row.count }}</span>
          </div>
        </div>
      </div>

      <!-- エージェント別（識別なのでカテゴリカル配色。凡例として色+ラベルを併記） -->
      <div v-if="agentBars.length">
        <h3 class="text-xs font-semibold text-muted-foreground mb-2">{{ t('dashboard.byAgent') }}</h3>
        <div class="space-y-1.5">
          <div v-for="row in agentBars" :key="row.agent" class="flex items-center gap-2 text-sm">
            <span class="w-32 shrink-0 text-xs truncate flex items-center gap-1.5">
              <span class="w-2 h-2 rounded-full shrink-0" :style="{ background: row.color }" />
              {{ row.agent }}
            </span>
            <div class="flex-1 h-2 bg-gray-100 rounded-full overflow-hidden">
              <div class="h-full rounded-full" :style="{ width: row.pct + '%', background: row.color }" />
            </div>
            <span class="w-8 text-right text-xs tabular-nums text-muted-foreground">{{ row.count }}</span>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>
