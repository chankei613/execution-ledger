<script setup lang="ts">
import { useRoute } from 'vue-router'
import { RouterView, RouterLink } from 'vue-router'
import { useI18n } from '@/i18n'

const route = useRoute()
const { t, toggleLocale } = useI18n()

function isActive(prefix: string) {
  return route.path.startsWith(prefix)
}

const navItems = [
  {
    to: '/ledger',
    key: 'nav.ledger',
    icon: `<path d="M4 6h16M4 12h16M4 18h10"/>`,
  },
  {
    to: '/dashboard',
    key: 'nav.dashboard',
    icon: `<rect x="3" y="12" width="4" height="8"/><rect x="10" y="8" width="4" height="12"/><rect x="17" y="4" width="4" height="16"/>`,
  },
  {
    to: '/settings',
    key: 'nav.settings',
    icon: `<circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 1 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 1 1-2.83-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 1 1 2.83-2.83l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 1 1 2.83 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z"/>`,
  },
]
</script>

<template>
  <div class="flex h-screen bg-background text-foreground overflow-hidden">
    <!-- Sidebar -->
    <aside class="w-52 border-r border-border flex flex-col shrink-0 bg-background">
      <!-- Header: TitleBarHiddenInset のトラフィックライト分の余白（pt-16） -->
      <div class="px-4 pb-4 pt-16 border-b border-border" style="-webkit-app-region: drag">
        <div class="flex items-start justify-between">
          <div class="flex items-center gap-2">
            <svg class="w-5 h-5 text-gray-700 shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75">
              <rect x="3" y="4" width="18" height="16" rx="2"/><path d="M7 8h6M7 12h10M7 16h8"/>
            </svg>
            <div>
              <h1 class="text-sm font-semibold tracking-tight text-foreground">{{ t('app.subtitle') }}</h1>
            </div>
          </div>
          <button
            @click="toggleLocale"
            style="-webkit-app-region: no-drag"
            class="text-xs text-gray-400 hover:text-gray-700 px-1.5 py-0.5 rounded border border-gray-200 hover:border-gray-400 transition-colors shrink-0 mt-0.5"
          >{{ t('lang.toggle') }}</button>
        </div>
      </div>

      <nav class="flex-1 p-2 space-y-0.5">
        <RouterLink
          v-for="item in navItems"
          :key="item.to"
          :to="item.to"
          class="flex items-center gap-2.5 px-3 py-2 rounded-md text-sm transition-colors"
          :class="isActive(item.to)
            ? 'bg-gray-900 text-white'
            : 'text-gray-500 hover:bg-gray-100 hover:text-gray-900'"
        >
          <svg class="w-4 h-4 shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" v-html="item.icon" />
          <span>{{ t(item.key) }}</span>
        </RouterLink>
      </nav>
    </aside>

    <!-- Main content -->
    <main class="flex-1 overflow-hidden bg-gray-50/50 flex flex-col">
      <RouterView class="flex-1 overflow-hidden" />
    </main>
  </div>
</template>
