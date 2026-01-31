import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { SystemConfig, Tag } from '../types'
import { GetSystemConfig, UpdateSystemConfig } from '../../wailsjs/go/main/App'

export const useSettingsStore = defineStore('settings', () => {
  // State
  const config = ref<SystemConfig | null>(null)
  const loading = ref(false)
  const saving = ref(false)

  // Tag state
  const accountTypes = ref<Tag[]>([
    { label: 'PLUS', value: 'PLUS', color: '#2080f0' },
    { label: 'BUSINESS', value: 'BUSINESS', color: '#18a058' },
    { label: 'FREE', value: 'FREE', color: '#909399' }
  ])

  const accountStatuses = ref<Tag[]>([
    { label: '未售出', value: 'unsold', color: '#f0a020' },
    { label: '已售出', value: 'sold', color: '#18a058' }
  ])

  // Computed
  const defaultValidityDays = computed(() => config.value?.defaultValidityDays || 30)
  const reminderDaysBefore = computed(() => config.value?.reminderDaysBefore || 1)
  const copyFormat = computed(() => config.value?.copyFormat || '账号：{account}\n密码：{password}')
  const emailFormat = computed(() => config.value?.emailFormat || '您的账号 {account} 将在 {expireAt} 过期，请及时处理。')

  // Actions
  async function fetchConfig() {
    loading.value = true
    try {
      const result = await GetSystemConfig()
      if (result) {
        config.value = result

        // Load tags from config
        if (result.accountTypes) {
          try {
            const parsed = JSON.parse(result.accountTypes)
            if (Array.isArray(parsed) && parsed.length > 0) {
              accountTypes.value = parsed
            }
          } catch (e) {
            console.warn('Failed to parse account types:', e)
          }
        }

        if (result.accountStatuses) {
          try {
            const parsed = JSON.parse(result.accountStatuses)
            if (Array.isArray(parsed) && parsed.length > 0) {
              accountStatuses.value = parsed
            }
          } catch (e) {
            console.warn('Failed to parse account statuses:', e)
          }
        }
      }
    } catch (error) {
      console.error('Failed to fetch system config:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  async function updateConfig(
    defaultValidityDays: number,
    reminderDaysBefore: number,
    copyFormat: string,
    emailFormat: string
  ) {
    saving.value = true
    try {
      await UpdateSystemConfig(
        defaultValidityDays,
        reminderDaysBefore,
        copyFormat,
        emailFormat,
        JSON.stringify(accountTypes.value),
        JSON.stringify(accountStatuses.value)
      )

      // Update local state
      if (config.value) {
        config.value.defaultValidityDays = defaultValidityDays
        config.value.reminderDaysBefore = reminderDaysBefore
        config.value.copyFormat = copyFormat
        config.value.emailFormat = emailFormat
      }
    } catch (error) {
      console.error('Failed to update system config:', error)
      throw error
    } finally {
      saving.value = false
    }
  }

  async function updateTags() {
    if (!config.value) return

    saving.value = true
    try {
      await UpdateSystemConfig(
        config.value.defaultValidityDays,
        config.value.reminderDaysBefore,
        config.value.copyFormat,
        config.value.emailFormat,
        JSON.stringify(accountTypes.value),
        JSON.stringify(accountStatuses.value)
      )

      // Update local state
      config.value.accountTypes = JSON.stringify(accountTypes.value)
      config.value.accountStatuses = JSON.stringify(accountStatuses.value)
    } catch (error) {
      console.error('Failed to update tags:', error)
      throw error
    } finally {
      saving.value = false
    }
  }

  function addAccountType(tag: Tag) {
    accountTypes.value.push(tag)
  }

  function updateAccountType(index: number, tag: Tag) {
    if (index >= 0 && index < accountTypes.value.length) {
      accountTypes.value[index] = tag
    }
  }

  function deleteAccountType(index: number) {
    if (index >= 0 && index < accountTypes.value.length) {
      accountTypes.value.splice(index, 1)
    }
  }

  function addAccountStatus(tag: Tag) {
    accountStatuses.value.push(tag)
  }

  function updateAccountStatus(index: number, tag: Tag) {
    if (index >= 0 && index < accountStatuses.value.length) {
      accountStatuses.value[index] = tag
    }
  }

  function deleteAccountStatus(index: number) {
    if (index >= 0 && index < accountStatuses.value.length) {
      accountStatuses.value.splice(index, 1)
    }
  }

  return {
    // State
    config,
    loading,
    saving,
    accountTypes,
    accountStatuses,

    // Computed
    defaultValidityDays,
    reminderDaysBefore,
    copyFormat,
    emailFormat,

    // Actions
    fetchConfig,
    updateConfig,
    updateTags,
    addAccountType,
    updateAccountType,
    deleteAccountType,
    addAccountStatus,
    updateAccountStatus,
    deleteAccountStatus
  }
})
