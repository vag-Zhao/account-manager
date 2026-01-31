import { ref } from 'vue'
import type { Ref } from 'vue'

export interface Tag {
  label: string
  value: string
  color: string
}

/**
 * Composable for loading and managing tags from system config
 */
export function useTagLoader() {
  const accountTypes = ref<Tag[]>([
    { label: 'PLUS', value: 'PLUS', color: '#2080f0' },
    { label: 'BUSINESS', value: 'BUSINESS', color: '#18a058' },
    { label: 'FREE', value: 'FREE', color: '#909399' }
  ])

  const accountStatuses = ref<Tag[]>([
    { label: '未售出', value: 'unsold', color: '#f0a020' },
    { label: '已售出', value: 'sold', color: '#18a058' }
  ])

  const loadAccountTypes = (configJson: string | null | undefined) => {
    if (!configJson) return
    try {
      const parsed = JSON.parse(configJson)
      if (Array.isArray(parsed) && parsed.length > 0) {
        accountTypes.value = parsed
      }
    } catch (e) {
      console.warn('Failed to parse account types:', e)
    }
  }

  const loadAccountStatuses = (configJson: string | null | undefined, t?: (key: string) => string) => {
    if (!configJson) return
    try {
      const parsed = JSON.parse(configJson)
      if (Array.isArray(parsed) && parsed.length > 0) {
        accountStatuses.value = parsed
      }
    } catch (e) {
      console.warn('Failed to parse account statuses:', e)
      // Use default with i18n if available
      if (t) {
        accountStatuses.value = [
          { label: t('accounts.unsold'), value: 'unsold', color: '#f0a020' },
          { label: t('accounts.sold'), value: 'sold', color: '#18a058' }
        ]
      }
    }
  }

  const loadAllTags = (
    accountTypesJson: string | null | undefined,
    accountStatusesJson: string | null | undefined,
    t?: (key: string) => string
  ) => {
    loadAccountTypes(accountTypesJson)
    loadAccountStatuses(accountStatusesJson, t)
  }

  return {
    accountTypes,
    accountStatuses,
    loadAccountTypes,
    loadAccountStatuses,
    loadAllTags
  }
}
