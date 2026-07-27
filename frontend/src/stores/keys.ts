import { defineStore } from 'pinia'
import { ListKeys, IssueKey, RevokeKey, GetIngestionURL, GetAppVersion } from '../../wailsjs/go/main/App'
import { db } from '../../wailsjs/go/models'

export const useKeysStore = defineStore('keys', {
  state: () => ({
    keys: [] as db.AgentKey[],
    ingestionURL: '',
    appVersion: '',
    lastIssuedKey: null as { name: string; apiKey: string } | null,
    loading: false,
    error: null as string | null,
  }),
  actions: {
    async load() {
      this.loading = true
      this.error = null
      try {
        const [keys, url, version] = await Promise.all([
          ListKeys(),
          GetIngestionURL(),
          GetAppVersion(),
        ])
        this.keys = keys ?? []
        this.ingestionURL = url
        this.appVersion = version
      } catch (e) {
        this.error = String(e)
      } finally {
        this.loading = false
      }
    },
    async issue(name: string) {
      const result = await IssueKey(name)
      this.lastIssuedKey = { name: result.name, apiKey: result.api_key }
      await this.load()
    },
    async revoke(id: string) {
      await RevokeKey(id)
      await this.load()
    },
    dismissIssuedKey() {
      this.lastIssuedKey = null
    },
  },
})
