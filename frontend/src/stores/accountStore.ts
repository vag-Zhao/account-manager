import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Account, AccountStats, PaginatedAccounts } from '../types'
import {
  GetAccounts,
  GetStats,
  CreateAccount,
  UpdateAccount,
  DeleteAccount,
  MarkAsSold,
  MarkAsUnsold,
  BatchImport,
  GetSystemConfig
} from '../../wailsjs/go/main/App'

export const useAccountStore = defineStore('account', () => {
  const accounts = ref<Account[]>([])
  const stats = ref<AccountStats | null>(null)
  const total = ref(0)
  const page = ref(1)
  const pageSize = ref(50) // Increased from 8 to 50 for better performance
  const totalPages = ref(0)
  const loading = ref(false)

  // 系统配置缓存
  const systemConfig = ref<any>(null)
  const configLoaded = ref(false)

  // Filters
  const accountType = ref('')
  const isSold = ref<boolean | null>(null)
  const search = ref('')

  async function fetchAccounts() {
    loading.value = true
    try {
      const result = await GetAccounts(
        accountType.value,
        isSold.value,
        search.value,
        page.value,
        pageSize.value
      ) as PaginatedAccounts
      accounts.value = result.data || []
      total.value = result.total
      totalPages.value = result.totalPages
    } catch (error) {
      console.error('Failed to fetch accounts:', error)
    } finally {
      loading.value = false
    }
  }

  async function fetchStats() {
    try {
      stats.value = await GetStats() as AccountStats
    } catch (error) {
      console.error('Failed to fetch stats:', error)
    }
  }

  async function loadSystemConfig() {
    // 总是重新加载配置，确保获取最新的标签
    try {
      const config = await GetSystemConfig()
      systemConfig.value = config
      configLoaded.value = true
      return config
    } catch (error) {
      console.error('Failed to load system config:', error)
      return null
    }
  }

  async function createAccount(
    account: string,
    password: string,
    type: string,
    expireAt: string,
    isSold: boolean
  ) {
    await CreateAccount(account, password, type, expireAt, isSold)
    await fetchAccounts()
    await fetchStats()
  }

  async function updateAccount(
    id: number,
    account: string,
    password: string,
    type: string,
    expireAt: string,
    isSold: boolean
  ) {
    await UpdateAccount(id, account, password, type, expireAt, isSold)
    await fetchAccounts()
    await fetchStats()
  }

  async function deleteAccount(id: number) {
    await DeleteAccount(id)
    await fetchAccounts()
    await fetchStats()
  }

  async function markAsSold(id: number) {
    await MarkAsSold(id)
    await fetchAccounts()
    await fetchStats()
  }

  async function markAsUnsold(id: number) {
    await MarkAsUnsold(id)
    await fetchAccounts()
    await fetchStats()
  }

  async function batchImport(data: any[]) {
    const result = await BatchImport(data)
    await fetchAccounts()
    await fetchStats()
    return result
  }

  function setPage(p: number) {
    page.value = p
    fetchAccounts()
  }

  function setPageSize(size: number) {
    pageSize.value = size
    page.value = 1
    fetchAccounts()
  }

  function setFilter(type: string, sold: boolean | null, searchText: string) {
    accountType.value = type
    isSold.value = sold
    search.value = searchText
    page.value = 1
    fetchAccounts()
  }

  function resetFilters() {
    accountType.value = ''
    isSold.value = null
    search.value = ''
    page.value = 1
    fetchAccounts()
  }

  return {
    accounts,
    stats,
    total,
    page,
    pageSize,
    totalPages,
    loading,
    accountType,
    isSold,
    search,
    systemConfig,
    fetchAccounts,
    fetchStats,
    loadSystemConfig,
    createAccount,
    updateAccount,
    deleteAccount,
    markAsSold,
    markAsUnsold,
    batchImport,
    setPage,
    setPageSize,
    setFilter,
    resetFilters
  }
})
