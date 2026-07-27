import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { createRouter, createWebHashHistory } from 'vue-router'
import 'virtual:uno.css'
import '@unocss/reset/tailwind.css'
import './assets/globals.css'
import App from './App.vue'
import LedgerView from './pages/LedgerView.vue'
import DashboardView from './pages/DashboardView.vue'
import SettingsView from './pages/SettingsView.vue'

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    { path: '/', redirect: '/ledger' },
    { path: '/ledger', component: LedgerView },
    { path: '/dashboard', component: DashboardView },
    { path: '/settings', component: SettingsView },
  ],
})

const app = createApp(App)
app.use(createPinia())
app.use(router)
app.mount('#app')
