import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { EmailConfig, EmailLog } from '../types'
import {
  GetEmailConfig,
  UpdateEmailConfig,
  TestEmailSend,
  GetEmailLogs,
  GetSMTPProviders
} from '../../wailsjs/go/main/App'

export interface SMTPProvider {
  name: string
  host: string
  port: number
  helpText: string
}

export const useEmailStore = defineStore('email', () => {
  // State
  const config = ref<EmailConfig | null>(null)
  const loading = ref(false)
  const saving = ref(false)
  const testing = ref(false)

  // Email logs
  const logs = ref<EmailLog[]>([])
  const logsTotal = ref(0)
  const logsLoading = ref(false)

  // SMTP providers
  const smtpProviders = ref<SMTPProvider[]>([])
  const providersLoaded = ref(false)

  // Computed
  const isActive = computed(() => config.value?.isActive || false)
  const smtpHost = computed(() => config.value?.smtpHost || '')
  const smtpPort = computed(() => config.value?.smtpPort || 465)
  const senderEmail = computed(() => config.value?.senderEmail || '')
  const recipientEmail = computed(() => config.value?.recipientEmail || '')

  // Actions
  async function fetchConfig() {
    loading.value = true
    try {
      const result = await GetEmailConfig()
      if (result) {
        config.value = result
      }
    } catch (error) {
      console.error('Failed to fetch email config:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  async function updateConfig(
    smtpHost: string,
    smtpPort: number,
    senderEmail: string,
    senderPassword: string,
    recipientEmail: string,
    isActive: boolean
  ) {
    saving.value = true
    try {
      await UpdateEmailConfig(
        smtpHost,
        smtpPort,
        senderEmail,
        senderPassword,
        recipientEmail,
        isActive
      )

      // Update local state
      if (config.value) {
        config.value.smtpHost = smtpHost
        config.value.smtpPort = smtpPort
        config.value.senderEmail = senderEmail
        config.value.senderPassword = senderPassword
        config.value.recipientEmail = recipientEmail
        config.value.isActive = isActive
      } else {
        config.value = {
          id: 0,
          smtpHost,
          smtpPort,
          senderEmail,
          senderPassword,
          recipientEmail,
          isActive,
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString()
        }
      }
    } catch (error) {
      console.error('Failed to update email config:', error)
      throw error
    } finally {
      saving.value = false
    }
  }

  async function testEmail() {
    testing.value = true
    try {
      await TestEmailSend()
    } catch (error) {
      console.error('Failed to send test email:', error)
      throw error
    } finally {
      testing.value = false
    }
  }

  async function fetchLogs(page: number, pageSize: number) {
    logsLoading.value = true
    try {
      const result = await GetEmailLogs(page, pageSize)
      if (result) {
        logs.value = result.logs || []
        logsTotal.value = result.total || 0
      }
    } catch (error) {
      console.error('Failed to fetch email logs:', error)
      throw error
    } finally {
      logsLoading.value = false
    }
  }

  async function loadSMTPProviders() {
    if (providersLoaded.value) return

    try {
      const providers = await GetSMTPProviders()
      smtpProviders.value = providers || []
      providersLoaded.value = true
    } catch (error) {
      console.error('Failed to load SMTP providers:', error)
      throw error
    }
  }

  function getProviderByHost(host: string): SMTPProvider | undefined {
    return smtpProviders.value.find(p => p.host === host)
  }

  return {
    // State
    config,
    loading,
    saving,
    testing,
    logs,
    logsTotal,
    logsLoading,
    smtpProviders,
    providersLoaded,

    // Computed
    isActive,
    smtpHost,
    smtpPort,
    senderEmail,
    recipientEmail,

    // Actions
    fetchConfig,
    updateConfig,
    testEmail,
    fetchLogs,
    loadSMTPProviders,
    getProviderByHost
  }
})
