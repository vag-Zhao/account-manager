<template>
  <div class="email-settings-container">
    <!-- 功能说明 / 服务状态 -->
    <div class="card intro-card">
      <!-- 邮件配置标签页的介绍 -->
      <div v-if="activeTab === 'email'" class="intro-section">
        <div class="card-header-simple">
          <n-icon :component="NotificationsOutline" :size="20" />
          <span>{{ $t('email.title') }}</span>
        </div>
        <div class="intro-content">
          <p class="intro-desc">{{ $t('email.intro') }}</p>
          <div class="intro-features">
            <div class="feature-item">
              <n-icon :component="TimeOutline" :size="16" />
              <span>{{ $t('email.features.autoDetect') }}</span>
            </div>
            <div class="feature-item">
              <n-icon :component="MailOutline" :size="16" />
              <span>{{ $t('email.features.sendEmail') }}</span>
            </div>
            <div class="feature-item">
              <n-icon :component="SettingsOutline" :size="16" />
              <span>{{ $t('email.features.configReminder') }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 服务器部署标签页的状态 -->
      <div v-if="activeTab === 'server'" class="status-section">
        <div class="card-header-simple">
          <n-icon :component="SettingsOutline" :size="20" />
          <span>{{ $t('email.serviceStatus') }}</span>
          <n-button size="small" @click="refreshStatus" :loading="refreshing" style="margin-left: auto;">
            {{ $t('email.refreshStatus') }}
          </n-button>
        </div>
        <div class="status-grid">
          <div class="status-item">
            <span class="status-item-label">{{ $t('email.serviceStatus') }}</span>
            <span class="status-item-value" :class="'status-' + serviceStatus">{{ serviceStatusText }}</span>
          </div>
          <div class="status-item">
            <span class="status-item-label">{{ $t('email.serverAddress') }}</span>
            <span class="status-item-value">{{ serverConfig.host || $t('email.notConfigured') }}</span>
          </div>
          <div class="status-item">
            <span class="status-item-label">{{ $t('email.deployPath') }}</span>
            <span class="status-item-value">{{ serverConfig.deployPath || $t('email.notConfigured') }}</span>
          </div>
          <div class="status-item">
            <span class="status-item-label">{{ $t('email.lastDeployed') }}</span>
            <span class="status-item-value">{{ serverConfig.lastDeployedAt || $t('email.neverDeployed') }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 标签页主体 -->
    <div class="card tabs-card">
      <n-tabs type="line" animated v-model:value="activeTab">
        <template #suffix>
          <div v-if="activeTab === 'server'" class="service-control-buttons">
            <n-button
              size="small"
              @click.stop="startService"
              :loading="starting"
              :disabled="serviceStatus === 'running'"
              class="service-btn start-btn"
            >
              <template #icon>
                <n-icon :component="PlayOutline" :size="14" />
              </template>
              启动
            </n-button>
            <n-button
              size="small"
              @click.stop="stopService"
              :loading="stopping"
              :disabled="serviceStatus !== 'running'"
              class="service-btn stop-btn"
            >
              <template #icon>
                <n-icon :component="StopOutline" :size="14" />
              </template>
              停止
            </n-button>
          </div>
        </template>
        <!-- 邮件配置标签页 -->
        <n-tab-pane name="email" :tab="$t('email.config')">
          <div class="tab-content">
            <div class="config-section">
              <div class="section-header">
                <div class="header-left">
                  <n-icon :component="MailOutline" :size="20" />
                  <span>{{ $t('email.smtpConfig') }}</span>
                  <div class="service-toggle">
                    <n-switch v-model:value="form.isActive" size="small" @update:value="debouncedSaveEmail" />
                  </div>
                </div>
              </div>
              <n-spin :show="loading">
                <n-form :model="form" label-placement="top" size="small">
                  <n-grid :cols="2" :x-gap="16">
                    <n-form-item-gi :label="$t('email.smtpHost')">
                      <n-select
                        v-model:value="form.smtpHost"
                        :options="smtpHostOptions"
                        filterable
                        tag
                        :placeholder="$t('email.smtpHostPlaceholder')"
                        :disabled="!form.isActive"
                        @update:value="handleSmtpHostChange"
                      />
                    </n-form-item-gi>
                    <n-form-item-gi :label="$t('email.port')">
                      <n-input-number v-model:value="form.smtpPort" :min="1" style="width:100%" :placeholder="$t('email.portPlaceholder')" :disabled="!form.isActive || isKnownProvider" @update:value="debouncedSaveEmail" />
                    </n-form-item-gi>
                    <n-form-item-gi :label="$t('email.senderEmail')">
                      <n-input v-model:value="form.senderEmail" :placeholder="$t('email.senderEmailPlaceholder')" :disabled="!form.isActive" @update:value="debouncedSaveEmail" />
                    </n-form-item-gi>
                    <n-form-item-gi :label="$t('email.authCode')">
                      <n-input v-model:value="form.senderPassword" type="password" show-password-on="click" :placeholder="$t('email.authCodePlaceholder')" :disabled="!form.isActive" @update:value="debouncedSaveEmail" />
                    </n-form-item-gi>
                  </n-grid>
                  <div class="recipient-row">
                    <div class="recipient-input">
                      <div class="form-label">{{ $t('email.recipientEmail') }}</div>
                      <n-input v-model:value="form.recipientEmail" :placeholder="$t('email.recipientEmailPlaceholder')" :disabled="!form.isActive" size="small" @update:value="debouncedSaveEmail" />
                    </div>
                    <div class="test-button">
                      <div class="form-label">&nbsp;</div>
                      <n-button secondary @click="doTest" :loading="testing" :disabled="!form.isActive" class="test-btn" size="small">
                        <template #icon>
                          <n-icon :component="MailOutline" :size="16" />
                        </template>
                        {{ $t('email.sendTest') }}
                      </n-button>
                    </div>
                  </div>
                  <!-- 授权码帮助提示 -->
                  <n-alert v-if="providerHelpText" type="warning" style="margin-top: 16px;">
                    <strong>{{ $t('email.howToGetAuthCode') }}</strong><br />
                    {{ providerHelpText }}
                  </n-alert>
                </n-form>
              </n-spin>
            </div>
          </div>
        </n-tab-pane>

        <!-- 服务器部署标签页 -->
        <n-tab-pane name="server" :tab="$t('email.serverDeploy')">
          <div class="tab-content">
            <n-alert type="info" style="margin-bottom: 6px;">
              {{ $t('email.deployIntro') }}
            </n-alert>

            <div class="deploy-grid">
              <!-- SSH配置 -->
              <div class="config-section">
                <div class="section-header">
                  <div class="header-left">
                    <n-icon :component="SettingsOutline" :size="20" />
                    <span>{{ $t('email.sshConfig') }}</span>
                    <div class="service-toggle">
                      <n-switch v-model:value="serverConfig.isActive" @update:value="debouncedSaveServerConfig" size="small" />
                      <span style="margin-left: 8px; font-size: 12px; color: #999;">
                        {{ serverConfig.isActive ? $t('email.enabled') : $t('email.disabled') }}
                      </span>
                    </div>
                  </div>
                </div>
                <n-form :model="serverConfig" label-placement="top" size="small">
                  <n-grid :cols="2" :x-gap="16">
                    <n-form-item-gi :label="$t('email.serverAddress')">
                      <n-input v-model:value="serverConfig.host" :placeholder="$t('email.serverAddressPlaceholder')" :disabled="!serverConfig.isActive" @update:value="debouncedSaveServerConfig" />
                    </n-form-item-gi>
                    <n-form-item-gi :label="$t('email.sshPort')">
                      <n-input-number v-model:value="serverConfig.port" :min="1" :max="65535" style="width: 100%;" :disabled="!serverConfig.isActive" @update:value="debouncedSaveServerConfig" />
                    </n-form-item-gi>
                    <n-form-item-gi :label="$t('email.username')">
                      <n-input v-model:value="serverConfig.username" :placeholder="$t('email.usernamePlaceholder')" :disabled="!serverConfig.isActive" @update:value="debouncedSaveServerConfig" />
                    </n-form-item-gi>
                    <n-form-item-gi :label="$t('email.deployPath')">
                      <n-input v-model:value="serverConfig.deployPath" :placeholder="$t('email.deployPathPlaceholder')" :disabled="!serverConfig.isActive" @update:value="debouncedSaveServerConfig" />
                    </n-form-item-gi>
                  </n-grid>
                  <n-form-item :label="$t('email.authMethod')">
                    <n-radio-group v-model:value="authMethod" :disabled="!serverConfig.isActive" @update:value="debouncedSaveServerConfig">
                      <n-radio value="password">{{ $t('email.authPassword') }}</n-radio>
                      <n-radio value="privateKey">{{ $t('email.authPrivateKey') }}</n-radio>
                    </n-radio-group>
                  </n-form-item>
                  <div v-if="authMethod === 'password'" class="password-row">
                    <div class="password-input">
                      <n-input
                        v-model:value="serverConfig.password"
                        type="password"
                        show-password-on="click"
                        :placeholder="$t('email.sshPasswordPlaceholder')"
                        :disabled="!serverConfig.isActive"
                        size="small"
                        @update:value="debouncedSaveServerConfig"
                      />
                    </div>
                    <div class="test-connection-button">
                      <n-button secondary @click="testConnection" :loading="testing" :disabled="!serverConfig.isActive" size="small">
                        <template #icon>
                          <n-icon :component="SettingsOutline" :size="16" />
                        </template>
                        {{ $t('email.testConnection') }}
                      </n-button>
                    </div>
                  </div>
                  <div v-if="authMethod === 'privateKey'">
                    <n-input
                      v-model:value="serverConfig.privateKey"
                      type="textarea"
                      :rows="6"
                      :placeholder="$t('email.privateKeyPlaceholder')"
                      :disabled="!serverConfig.isActive"
                      @update:value="debouncedSaveServerConfig"
                    />
                    <n-button secondary @click="testConnection" :loading="testing" :disabled="!serverConfig.isActive" size="small" style="margin-top: 12px;">
                      <template #icon>
                        <n-icon :component="SettingsOutline" :size="16" />
                      </template>
                      {{ $t('email.testConnection') }}
                    </n-button>
                  </div>
                </n-form>

                <!-- Server Detection Results -->
                <n-spin :show="detectingServer">
                  <n-alert v-if="serverInfo" :type="serverInfo.hasSystemd ? 'success' : 'warning'" style="margin-top: 16px;">
                    <template #header>
                      <strong>{{ $t('email.serverDetectionResult') }}</strong>
                    </template>
                    <div style="font-size: 13px; line-height: 1.8;">
                      <div><strong>{{ $t('email.osInfo') }}:</strong> {{ serverInfo.osPrettyName || serverInfo.osName }}</div>
                      <div><strong>{{ $t('email.osVersion') }}:</strong> {{ serverInfo.osVersion }}</div>
                      <div><strong>{{ $t('email.packageManager') }}:</strong> {{ serverInfo.packageManager }}</div>
                      <div><strong>{{ $t('email.systemdStatus') }}:</strong>
                        <n-tag :type="serverInfo.hasSystemd ? 'success' : 'error'" size="small" style="margin-left: 8px;">
                          {{ serverInfo.hasSystemd ? $t('email.systemdInstalled') : $t('email.systemdNotInstalled') }}
                        </n-tag>
                        <span v-if="serverInfo.hasSystemd" style="margin-left: 8px; color: #666;">
                          ({{ $t('email.systemdVersion') }}: {{ serverInfo.systemdVersion }})
                        </span>
                      </div>
                      <div v-if="!serverInfo.hasSystemd" style="margin-top: 8px; color: #d03050;">
                        <strong>{{ $t('common.warning') }}:</strong> {{ $t('email.systemdWarning') }}
                      </div>
                    </div>
                  </n-alert>
                </n-spin>
              </div>

              <!-- 部署操作 -->
              <div class="config-section">
                <div class="section-header-simple">
                  <n-icon :component="SettingsOutline" :size="20" />
                  <span>{{ $t('email.deployManagement') }}</span>
                </div>

                <!-- 部署操作 -->
                <div class="deploy-actions">
                  <n-space vertical>
                    <n-button
                      type="primary"
                      @click="deployService"
                      :loading="deploying"
                      :disabled="!serverConfig.isActive"
                      style="width: 100%;"
                    >
                      {{ $t('email.deployService') }}
                    </n-button>
                  </n-space>
                </div>

                <!-- 部署日志 -->
                <div class="deploy-log">
                  <div class="log-header">{{ $t('email.deployLog') }}</div>
                  <n-log :log="deployLog" :rows="10" language="log" />
                </div>
              </div>
            </div>
          </div>
        </n-tab-pane>
      </n-tabs>
    </div>

    <!-- SSH主机密钥确认对话框 -->
    <n-modal v-model:show="showHostKeyDialog" preset="dialog" :title="$t('email.hostKeyVerification')" style="width: 600px;">
      <n-alert type="warning" style="margin-bottom: 16px;">
        首次连接到此SSH服务器，需要验证主机密钥以确保连接安全。
      </n-alert>
      <n-descriptions :column="1" bordered size="small">
        <n-descriptions-item label="主机">{{ hostKeyInfo.host }}</n-descriptions-item>
        <n-descriptions-item label="密钥类型">{{ hostKeyInfo.keyType }}</n-descriptions-item>
        <n-descriptions-item label="指纹">
          <code style="font-size: 11px; word-break: break-all;">{{ hostKeyInfo.fingerprint }}</code>
        </n-descriptions-item>
      </n-descriptions>
      <template #action>
        <n-space>
          <n-button @click="showHostKeyDialog = false">取消</n-button>
          <n-button type="primary" @click="trustAndRetry">信任并继续</n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { NForm, NFormItem, NFormItemGi, NGrid, NInput, NInputNumber, NSwitch, NButton, NSpace, NSpin, NIcon, useMessage, NTabs, NTabPane, NAlert, NRadioGroup, NRadio, NTag, NLog, NSelect, NModal, NDescriptions, NDescriptionsItem } from 'naive-ui'
import { MailOutline, InformationCircleOutline, NotificationsOutline, TimeOutline, SettingsOutline, PlayOutline, StopOutline } from '@vicons/ionicons5'
import { GetEmailConfig, UpdateEmailConfig, TestEmailSend, GetServerConfig, UpdateServerConfig, TestServerConnection, DeployEmailService, GetServiceStatus, StartEmailService, StopEmailService, GetSMTPProviders, DetectServerInfo, TrustHostKey } from '../../wailsjs/go/main/App'

const { t: $t } = useI18n()
const message = useMessage()

// Tab state
const activeTab = ref('email')

// Email config state
const loading = ref(false)
const testing = ref(false)
const form = ref({ smtpHost: '', smtpPort: 465, senderEmail: '', senderPassword: '', recipientEmail: '', isActive: false })

// SMTP Providers
const smtpProviders = ref<any[]>([])
const smtpHostOptions = ref<any[]>([])
const providerHelpText = ref<string>('')

// Server config state
const serverConfig = ref({
  host: '',
  port: 22,
  username: 'root',
  password: '',
  privateKey: '',
  deployPath: '/opt/account-manager-email',
  isActive: false,
  lastDeployedAt: '',
  serviceStatus: 'unknown'
})
const authMethod = ref<'password' | 'privateKey'>('password')
const deploying = ref(false)
const starting = ref(false)
const stopping = ref(false)
const refreshing = ref(false)
const serviceStatus = ref('unknown')
const deployLog = ref('')
const serverInfo = ref<any>(null)
const detectingServer = ref(false)

// SSH主机密钥确认对话框
const showHostKeyDialog = ref(false)
const hostKeyInfo = ref({
  id: 0,
  host: '',
  keyType: '',
  fingerprint: ''
})
const pendingAction = ref<(() => Promise<void>) | null>(null)

const serviceStatusType = computed(() => {
  switch (serviceStatus.value) {
    case 'running': return 'success'
    case 'stopped': return 'warning'
    case 'disabled': return 'default'
    default: return 'info'
  }
})

const serviceStatusText = computed(() => {
  switch (serviceStatus.value) {
    case 'running': return $t('email.statusRunning')
    case 'stopped': return $t('email.statusStopped')
    case 'disabled': return $t('email.statusDisabled')
    default: return $t('email.statusUnknown')
  }
})

// 检查是否为已知的SMTP服务商
const isKnownProvider = computed(() => {
  return smtpHostOptions.value.some(p => p.value === form.value.smtpHost)
})

onMounted(async () => {
  loading.value = true
  try {
    await loadSMTPProviders()
    const c = await GetEmailConfig()
    if (c && c.smtpHost) {
      form.value = { smtpHost: c.smtpHost, smtpPort: c.smtpPort, senderEmail: c.senderEmail, senderPassword: c.senderPassword, recipientEmail: c.recipientEmail, isActive: c.isActive }
    }
    await loadServerConfig()
    await refreshStatus()
  } catch {}
  loading.value = false
})

async function doSave() {
  try {
    await UpdateEmailConfig(form.value.smtpHost, form.value.smtpPort, form.value.senderEmail, form.value.senderPassword, form.value.recipientEmail, form.value.isActive)
  } catch (e: any) {
    console.error('保存邮件配置失败:', e)
  }
}

// Debounced save function for auto-save email config
let saveEmailTimer: number | null = null
function debouncedSaveEmail() {
  if (saveEmailTimer) {
    clearTimeout(saveEmailTimer)
  }
  saveEmailTimer = setTimeout(async () => {
    await doSave()
  }, 300) as unknown as number
}

async function doTest() {
  testing.value = true
  try { await TestEmailSend(); message.success($t('email.testSent')) }
  catch (e: any) { message.error(e.toString()) }
  testing.value = false
}

// SMTP Provider functions
async function loadSMTPProviders() {
  try {
    const providers = await GetSMTPProviders()
    smtpProviders.value = providers
    smtpHostOptions.value = providers.map((p: any) => ({
      label: `${p.name} (${p.host})`,
      value: p.host,
      host: p.host,
      port: p.port,
      helpText: p.helpText
    }))
  } catch (error: any) {
    console.error('加载SMTP服务商失败:', error)
  }
}

function handleSmtpHostChange(value: string) {
  // Extract hostname if user entered "host:port" format
  let hostname = value
  if (value && value.includes(':')) {
    const parts = value.split(':')
    hostname = parts[0]
    // If port is provided, update it
    if (parts[1] && !isNaN(parseInt(parts[1]))) {
      form.value.smtpPort = parseInt(parts[1])
    }
  }

  // Update the form value to hostname only
  if (hostname !== value) {
    form.value.smtpHost = hostname
  }

  const provider = smtpHostOptions.value.find(p => p.value === hostname)
  if (provider) {
    form.value.smtpPort = provider.port
    providerHelpText.value = provider.helpText
  } else {
    providerHelpText.value = ''
  }
  // Auto-save after host change
  debouncedSaveEmail()
}

// Server config functions
async function loadServerConfig() {
  try {
    const config = await GetServerConfig()
    if (config) {
      serverConfig.value = {
        host: config.host || '',
        port: config.port || 22,
        username: config.username || 'root',
        password: config.password || '',
        privateKey: config.privateKey || '',
        deployPath: config.deployPath || '/opt/account-manager-email',
        isActive: config.isActive || false,
        lastDeployedAt: config.lastDeployedAt ? new Date(config.lastDeployedAt).toLocaleString('zh-CN') : '',
        serviceStatus: config.serviceStatus || 'unknown'
      }
    }
  } catch (error: any) {
    // 静默处理首次加载时的 record not found 错误
    if (!error.toString().includes('record not found')) {
      console.error('加载服务器配置失败:', error)
    }
  }
}

async function saveServerConfig() {
  try {
    await UpdateServerConfig(
      serverConfig.value.host,
      serverConfig.value.port,
      serverConfig.value.username,
      authMethod.value === 'password' ? serverConfig.value.password : '',
      authMethod.value === 'privateKey' ? serverConfig.value.privateKey : '',
      serverConfig.value.deployPath,
      serverConfig.value.isActive
    )
  } catch (error: any) {
    console.error('保存服务器配置失败:', error)
  }
}

// Debounced save function for auto-save
let saveServerConfigTimer: number | null = null
function debouncedSaveServerConfig() {
  if (saveServerConfigTimer) {
    clearTimeout(saveServerConfigTimer)
  }
  saveServerConfigTimer = setTimeout(async () => {
    await saveServerConfig()
  }, 300) as unknown as number
}

async function testConnection() {
  testing.value = true
  try {
    await TestServerConnection()
    message.success($t('email.testConnectionSuccess'))

    // Automatically detect server info after successful connection
    detectingServer.value = true
    try {
      serverInfo.value = await DetectServerInfo()

      // Show warning if systemd is not available
      if (!serverInfo.value.hasSystemd) {
        message.warning($t('email.systemdWarning'))
      }
    } catch (error: any) {
      console.error('检测服务器信息失败:', error)
      message.warning($t('email.serverDetectionFailed') + ': ' + error)
    } finally {
      detectingServer.value = false
    }
  } catch (error: any) {
    // 检查是否是未信任的主机密钥错误
    if (handleHostKeyError(error.toString(), testConnection)) {
      return
    }
    message.error($t('email.testConnectionFailed') + ': ' + error)
  } finally {
    testing.value = false
  }
}

async function deployService() {
  deploying.value = true
  deployLog.value = $t('email.deploying') + '\n'
  try {
    deployLog.value += 'Connecting to server...\n'
    deployLog.value += 'Building email service...\n'
    deployLog.value += 'Uploading files...\n'
    deployLog.value += 'Configuring systemd service...\n'

    await DeployEmailService()

    deployLog.value += 'Deployment successful!\n'
    deployLog.value += 'Service started and enabled on boot\n'
    message.success($t('email.deploySuccess'))

    await loadServerConfig()
    await refreshStatus()
  } catch (error: any) {
    // 检查是否是未信任的主机密钥错误
    if (handleHostKeyError(error.toString(), deployService)) {
      return
    }
    deployLog.value += $t('email.deployFailed') + ': ' + error + '\n'
    message.error($t('email.deployFailed') + ': ' + error)
  } finally {
    deploying.value = false
  }
}

async function refreshStatus() {
  refreshing.value = true
  try {
    const status = await GetServiceStatus()
    serviceStatus.value = status
    message.success($t('email.statusRefreshed'))
  } catch (error: any) {
    // 检查是否是未信任的主机密钥错误
    if (handleHostKeyError(error.toString(), refreshStatus)) {
      return
    }
    // 静默处理首次加载时的 record not found 错误
    if (!error.toString().includes('record not found')) {
      console.error('获取服务状态失败:', error)
      message.error($t('email.statusRefreshFailed') + ': ' + error)
    }
  } finally {
    refreshing.value = false
  }
}

async function startService() {
  starting.value = true
  try {
    await StartEmailService()
    message.success($t('email.serviceStarted'))
    await refreshStatus()
  } catch (error: any) {
    message.error($t('common.error') + ': ' + error)
  } finally {
    starting.value = false
  }
}

async function stopService() {
  stopping.value = true
  try {
    await StopEmailService()
    message.success($t('email.serviceStopped'))
    await refreshStatus()
  } catch (error: any) {
    message.error($t('common.error') + ': ' + error)
  } finally {
    stopping.value = false
  }
}

// 处理SSH主机密钥错误
function handleHostKeyError(errorMsg: string, retryAction: () => Promise<void>): boolean {
  // 检查错误消息格式: UNTRUSTED_HOST_KEY|id|host|keyType|fingerprint
  if (errorMsg.includes('UNTRUSTED_HOST_KEY')) {
    const parts = errorMsg.split('|')
    if (parts.length >= 5) {
      hostKeyInfo.value = {
        id: parseInt(parts[1]),
        host: parts[2],
        keyType: parts[3],
        fingerprint: parts[4]
      }
      pendingAction.value = retryAction
      showHostKeyDialog.value = true
      return true
    }
  }
  return false
}

// 信任主机密钥并重试
async function trustAndRetry() {
  try {
    await TrustHostKey(hostKeyInfo.value.id)
    message.success('主机密钥已信任')
    showHostKeyDialog.value = false

    // 重试之前的操作
    if (pendingAction.value) {
      await pendingAction.value()
      pendingAction.value = null
    }
  } catch (error: any) {
    message.error('信任失败: ' + error.toString())
  }
}
</script>

<style scoped>
.email-settings-container {
  height: 100%;
  overflow-y: auto;
}

.intro-card {
  margin-bottom: 16px;
}

.intro-section,
.status-section {
  height: 120px;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.intro-content {
  display: flex;
  flex-direction: column;
  gap: 12px;
  flex: 1;
}

.intro-desc {
  font-size: 13px;
  color: #666;
  line-height: 1.5;
  margin: 0;
}

.intro-features {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.feature-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: #7C9885;
  background: rgba(124, 152, 133, 0.1);
  padding: 5px 10px;
  border-radius: 6px;
}

.status-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 10px;
  margin-top: 0;
  flex: 1;
  align-content: flex-start;
}

.status-grid .status-item {
  min-width: 0;
}

@media (max-width: 768px) {
  .status-grid {
    gap: 8px;
  }

  .status-item {
    padding: 6px;
  }
}

@media (max-width: 600px) {
  .intro-section,
  .status-section {
    height: auto;
    min-height: 120px;
  }

  .status-item {
    padding: 5px 6px;
  }

  .status-item-label {
    font-size: 10px;
  }

  .status-item-value {
    font-size: 11px;
  }
}

.status-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: 6px 8px;
  background: #FFFFFF;
  border: 1px solid #E8E4E0;
  border-radius: 6px;
  transition: all 0.2s ease;
  box-shadow: 0 1px 3px rgba(0,0,0,0.04);
}

.status-item:hover {
  border-color: rgba(124, 152, 133, 0.3);
  box-shadow: 0 2px 6px rgba(0,0,0,0.08);
  transform: translateY(-1px);
}

.status-item-label {
  font-size: 11px;
  color: #999;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.3px;
  line-height: 1.2;
}

.status-item-value {
  font-size: 12px;
  color: #333;
  font-family: monospace;
  font-weight: 500;
  line-height: 1.3;
  word-break: break-all;
}

.status-item-value.status-running {
  color: #18a058;
  font-weight: 600;
}

.status-item-value.status-stopped {
  color: #f0a020;
  font-weight: 600;
}

.status-item-value.status-disabled {
  color: #999;
}

.status-item-value.status-unknown {
  color: #666;
}

.tabs-card {
  margin-bottom: 0;
}

.service-control-buttons {
  display: flex;
  gap: 8px;
  align-items: center;
}

.service-btn {
  font-size: 12px;
  height: 26px;
  padding: 0 10px;
  border-radius: 6px;
  transition: all 0.2s ease;
  white-space: nowrap;
}

.start-btn {
  background: linear-gradient(135deg, #7C9885 0%, #6B8574 100%);
  color: white;
  border: none;
}

.start-btn:hover:not(:disabled) {
  background: linear-gradient(135deg, #6B8574 0%, #5A7463 100%);
  box-shadow: 0 2px 8px rgba(124, 152, 133, 0.3);
  transform: translateY(-1px);
}

.start-btn:disabled {
  background: #E8E4E0;
  color: #999;
  cursor: not-allowed;
}

.stop-btn {
  background: linear-gradient(135deg, #D4A574 0%, #C39463 100%);
  color: white;
  border: none;
}

.stop-btn:hover:not(:disabled) {
  background: linear-gradient(135deg, #C39463 0%, #B28352 100%);
  box-shadow: 0 2px 8px rgba(212, 165, 116, 0.3);
  transform: translateY(-1px);
}

.stop-btn:disabled {
  background: #E8E4E0;
  color: #999;
  cursor: not-allowed;
}

.tab-content {
  padding: 16px 0;
}

.config-section {
  background: rgba(0,0,0,0.02);
  border-radius: 8px;
  padding: 16px;
}

.deploy-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}

@media (max-width: 1100px) {
  .deploy-grid {
    grid-template-columns: 1fr;
  }
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  font-size: 14px;
  font-weight: 600;
  color: #333;
  margin-bottom: 16px;
}

.section-header-simple {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  font-weight: 600;
  color: #333;
  margin-bottom: 16px;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  font-size: 14px;
  font-weight: 600;
  color: #333;
  margin-bottom: 16px;
}

.card-header-simple {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  font-weight: 600;
  color: #333;
  margin-bottom: 12px;
  height: 32px;
  flex-shrink: 0;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.service-toggle {
  display: flex;
  align-items: center;
  padding-left: 16px;
  border-left: 1px solid rgba(0,0,0,0.08);
}

.save-btn {
  font-size: 13px;
  padding: 4px 12px;
  height: 28px;
  transition: opacity 0.2s ease;
}

.btn-hidden {
  opacity: 0;
  pointer-events: none;
}

.test-btn {
  font-size: 13px;
}

.recipient-row {
  display: flex;
  gap: 16px;
  margin-top: 0;
  align-items: flex-end;
}

.recipient-input {
  flex: 1;
}

.test-button {
  flex-shrink: 0;
}

.password-row {
  display: flex;
  gap: 16px;
  margin-top: 0;
  align-items: flex-end;
}

.password-input {
  flex: 1;
}

.test-connection-button {
  flex-shrink: 0;
}

.form-label {
  font-size: 14px;
  font-weight: 500;
  color: #333;
  margin-bottom: 8px;
  line-height: 1.5;
}

.deploy-actions {
  background: #FFFFFF;
  border: 1px solid #E8E4E0;
  border-radius: 8px;
  padding: 16px;
  margin-bottom: 16px;
}

.deploy-log {
  background: #FFFFFF;
  border: 1px solid #E8E4E0;
  border-radius: 8px;
  padding: 16px;
}

.log-header {
  font-size: 13px;
  font-weight: 600;
  color: #333;
  margin-bottom: 12px;
}
</style>
