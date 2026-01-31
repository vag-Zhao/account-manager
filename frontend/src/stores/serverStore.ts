import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import {
  GetServerConfig,
  UpdateServerConfig,
  TestServerConnection,
  DeployEmailService,
  GetServiceStatus,
  StartEmailService,
  StopEmailService,
  DetectServerInfo
} from '../../wailsjs/go/main/App'

export interface ServerConfig {
  id?: number
  host: string
  port: number
  username: string
  password: string
  privateKey: string
  deployPath: string
  isActive: boolean
  lastDeployedAt?: string
  serviceStatus?: string
  createdAt?: string
  updatedAt?: string
}

export interface ServerInfo {
  osName: string
  osPrettyName: string
  osVersion: string
  packageManager: string
  hasSystemd: boolean
  systemdVersion: string
}

export const useServerStore = defineStore('server', () => {
  // State
  const config = ref<ServerConfig>({
    host: '',
    port: 22,
    username: 'root',
    password: '',
    privateKey: '',
    deployPath: '/opt/account-manager-email',
    isActive: false
  })

  const loading = ref(false)
  const saving = ref(false)
  const testing = ref(false)
  const deploying = ref(false)
  const starting = ref(false)
  const stopping = ref(false)
  const refreshing = ref(false)

  const serviceStatus = ref<string>('unknown')
  const deployLog = ref<string>('')
  const serverInfo = ref<ServerInfo | null>(null)
  const detectingServer = ref(false)

  // Computed
  const isActive = computed(() => config.value.isActive)
  const isServiceRunning = computed(() => serviceStatus.value === 'running')
  const isServiceStopped = computed(() => serviceStatus.value === 'stopped')
  const lastDeployedAt = computed(() => {
    if (!config.value.lastDeployedAt) return null
    return new Date(config.value.lastDeployedAt).toLocaleString('zh-CN')
  })

  // Actions
  async function fetchConfig() {
    loading.value = true
    try {
      const result = await GetServerConfig()
      if (result) {
        config.value = {
          host: result.host || '',
          port: result.port || 22,
          username: result.username || 'root',
          password: result.password || '',
          privateKey: result.privateKey || '',
          deployPath: result.deployPath || '/opt/account-manager-email',
          isActive: result.isActive || false,
          lastDeployedAt: result.lastDeployedAt,
          serviceStatus: result.serviceStatus
        }
      }
    } catch (error: any) {
      // Silently handle "record not found" error on first load
      if (!error.toString().includes('record not found')) {
        console.error('Failed to fetch server config:', error)
        throw error
      }
    } finally {
      loading.value = false
    }
  }

  async function updateConfig(
    host: string,
    port: number,
    username: string,
    password: string,
    privateKey: string,
    deployPath: string,
    isActive: boolean
  ) {
    saving.value = true
    try {
      await UpdateServerConfig(
        host,
        port,
        username,
        password,
        privateKey,
        deployPath,
        isActive
      )

      // Update local state
      config.value = {
        ...config.value,
        host,
        port,
        username,
        password,
        privateKey,
        deployPath,
        isActive
      }
    } catch (error) {
      console.error('Failed to update server config:', error)
      throw error
    } finally {
      saving.value = false
    }
  }

  async function testConnection() {
    testing.value = true
    try {
      await TestServerConnection()

      // Automatically detect server info after successful connection
      await detectServer()
    } catch (error) {
      console.error('Failed to test server connection:', error)
      throw error
    } finally {
      testing.value = false
    }
  }

  async function detectServer() {
    detectingServer.value = true
    try {
      const info = await DetectServerInfo()
      serverInfo.value = info
      return info
    } catch (error) {
      console.error('Failed to detect server info:', error)
      throw error
    } finally {
      detectingServer.value = false
    }
  }

  async function deploy() {
    deploying.value = true
    deployLog.value = '开始部署...\n'

    try {
      deployLog.value += '连接服务器...\n'
      deployLog.value += '构建邮件服务...\n'
      deployLog.value += '上传文件...\n'
      deployLog.value += '配置 systemd 服务...\n'

      await DeployEmailService()

      deployLog.value += '部署成功！\n'
      deployLog.value += '服务已启动并设置为开机自启\n'

      // Refresh config and status after deployment
      await fetchConfig()
      await refreshStatus()
    } catch (error) {
      deployLog.value += '部署失败: ' + error + '\n'
      console.error('Failed to deploy service:', error)
      throw error
    } finally {
      deploying.value = false
    }
  }

  async function refreshStatus() {
    refreshing.value = true
    try {
      const status = await GetServiceStatus()
      serviceStatus.value = status
    } catch (error: any) {
      // Silently handle "record not found" error
      if (!error.toString().includes('record not found')) {
        console.error('Failed to get service status:', error)
        throw error
      }
    } finally {
      refreshing.value = false
    }
  }

  async function startService() {
    starting.value = true
    try {
      await StartEmailService()
      await refreshStatus()
    } catch (error) {
      console.error('Failed to start service:', error)
      throw error
    } finally {
      starting.value = false
    }
  }

  async function stopService() {
    stopping.value = true
    try {
      await StopEmailService()
      await refreshStatus()
    } catch (error) {
      console.error('Failed to stop service:', error)
      throw error
    } finally {
      stopping.value = false
    }
  }

  function clearDeployLog() {
    deployLog.value = ''
  }

  return {
    // State
    config,
    loading,
    saving,
    testing,
    deploying,
    starting,
    stopping,
    refreshing,
    serviceStatus,
    deployLog,
    serverInfo,
    detectingServer,

    // Computed
    isActive,
    isServiceRunning,
    isServiceStopped,
    lastDeployedAt,

    // Actions
    fetchConfig,
    updateConfig,
    testConnection,
    detectServer,
    deploy,
    refreshStatus,
    startService,
    stopService,
    clearDeployLog
  }
})
