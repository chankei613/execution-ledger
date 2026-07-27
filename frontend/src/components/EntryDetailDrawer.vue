<script setup lang="ts">
import { useI18n } from '@/i18n'
import { statusRole, statusIcon, STATUS_COLORS } from '@/statusColors'

const props = defineProps<{ entry: any }>()
defineEmits<{ close: [] }>()
const { t } = useI18n()

function fmt(v: any): string {
  const d = new Date(v)
  return isNaN(d.getTime()) ? String(v) : d.toLocaleString()
}
</script>

<template>
  <div class="w-96 border-l border-border overflow-y-auto shrink-0 bg-background">
    <div class="sticky top-0 bg-background border-b border-border px-4 py-3 flex items-center justify-between">
      <h3 class="text-sm font-semibold">{{ t('detail.title') }}</h3>
      <button @click="$emit('close')" class="text-muted-foreground hover:text-foreground text-lg leading-none">×</button>
    </div>

    <div class="p-4 space-y-4 text-sm">
      <div>
        <div class="flex items-center gap-2 mb-1">
          <span
            class="inline-flex items-center gap-1 px-1.5 rounded text-white text-[11px]"
            :style="{ background: STATUS_COLORS[statusRole(entry.status)].light }"
          >
            <span aria-hidden="true">{{ statusIcon(entry.status) }}</span>
            {{ t('status.' + entry.status) }}
          </span>
          <span class="text-xs text-muted-foreground">{{ fmt(entry.received_at) }}</span>
        </div>
        <div class="text-xs text-muted-foreground">{{ entry.source }} · {{ entry.agent_id }} · {{ entry.subject }}</div>
      </div>

      <section>
        <h4 class="text-xs font-semibold text-muted-foreground mb-1">{{ t('detail.summary') }}</h4>
        <p class="text-sm">{{ entry.summary || t('detail.none') }}</p>
      </section>

      <section>
        <h4 class="text-xs font-semibold text-muted-foreground mb-1">{{ t('detail.confidence') }}</h4>
        <div class="text-sm mb-1">overall: <strong>{{ entry.confidence_overall.toFixed(2) }}</strong></div>
        <div class="grid grid-cols-2 gap-1 text-xs text-muted-foreground" v-if="entry.confidence_breakdown">
          <div>task_understood: {{ entry.confidence_breakdown.task_understood }}</div>
          <div>execution_complete: {{ entry.confidence_breakdown.execution_complete }}</div>
          <div>correctness: {{ entry.confidence_breakdown.correctness }}</div>
          <div>side_effects_clean: {{ entry.confidence_breakdown.side_effects_clean }}</div>
        </div>
        <div v-if="entry.low_confidence_areas?.length" class="mt-1 flex flex-wrap gap-1">
          <span v-for="a in entry.low_confidence_areas" :key="a" class="text-[11px] px-1.5 py-0.5 rounded border" :style="{ borderColor: STATUS_COLORS.warning.light, color: STATUS_COLORS.warning.light }">
            {{ a }}
          </span>
        </div>
      </section>

      <section v-if="entry.criteria_results?.length">
        <h4 class="text-xs font-semibold text-muted-foreground mb-1">{{ t('detail.criteria') }}</h4>
        <ul class="space-y-1">
          <li v-for="(c, i) in entry.criteria_results" :key="i" class="flex items-start gap-1.5">
            <span :style="{ color: c.met ? STATUS_COLORS.good.light : STATUS_COLORS.critical.light }">{{ c.met ? '✓' : '✕' }}</span>
            <span>{{ c.description }}</span>
          </li>
        </ul>
      </section>

      <section>
        <h4 class="text-xs font-semibold text-muted-foreground mb-1">{{ t('detail.decisions') }}</h4>
        <div v-if="!entry.decisions?.length" class="text-xs text-muted-foreground">{{ t('detail.none') }}</div>
        <div v-else class="space-y-2">
          <div v-for="(d, i) in entry.decisions" :key="i" class="border border-border rounded p-2">
            <div class="font-medium">{{ d.description }}</div>
            <div class="text-xs text-muted-foreground mt-0.5">{{ d.rationale }}</div>
            <div v-if="d.alternatives_considered?.length" class="text-xs text-muted-foreground mt-1">
              alternatives: {{ d.alternatives_considered.join(', ') }}
            </div>
          </div>
        </div>
      </section>

      <section>
        <h4 class="text-xs font-semibold text-muted-foreground mb-1">{{ t('detail.actions') }}</h4>
        <div v-if="!entry.actions_taken?.length" class="text-xs text-muted-foreground">{{ t('detail.none') }}</div>
        <ul v-else class="space-y-1 font-mono text-xs">
          <li v-for="(a, i) in entry.actions_taken" :key="i">
            <span class="text-muted-foreground">{{ fmt(a.timestamp) }}</span> {{ a.tool }}: {{ a.input_summary }}
          </li>
        </ul>
      </section>

      <section v-if="entry.follow_up?.length">
        <h4 class="text-xs font-semibold text-muted-foreground mb-1">{{ t('detail.followUp') }}</h4>
        <ul class="space-y-1 text-xs">
          <li v-for="(f, i) in entry.follow_up" :key="i">{{ f.description }}</li>
        </ul>
      </section>

      <section v-if="entry.outputs && Object.keys(entry.outputs).length">
        <h4 class="text-xs font-semibold text-muted-foreground mb-1">{{ t('detail.outputs') }}</h4>
        <pre class="text-xs bg-gray-50 rounded p-2 overflow-x-auto">{{ JSON.stringify(entry.outputs, null, 2) }}</pre>
      </section>

      <section v-if="entry.usage">
        <h4 class="text-xs font-semibold text-muted-foreground mb-1">{{ t('detail.usage') }}</h4>
        <div class="text-xs text-muted-foreground">
          input: {{ entry.usage.input_tokens }} / output: {{ entry.usage.output_tokens }}
        </div>
      </section>
    </div>
  </div>
</template>
