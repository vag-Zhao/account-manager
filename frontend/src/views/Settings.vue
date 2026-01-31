<template>
  <div class="settings-container">
    <n-spin :show="loading" style="height: 100%;">
      <n-tabs type="line" animated class="settings-tabs">
        <!-- 系统设置 -->
        <n-tab-pane name="system" :tab="$t('settings.system')">
          <div class="tab-content">
            <div class="content-card">
              <n-form :model="form" label-placement="left" label-width="140" size="medium">
                <n-form-item :label="$t('settings.defaultValidity')">
                  <n-input-number v-model:value="form.defaultValidityDays" :min="1" :max="365" style="width:200px" />
                </n-form-item>
                <n-form-item :label="$t('settings.reminderDays')">
                  <n-input-number v-model:value="form.reminderDaysBefore" :min="1" :max="30" style="width:200px" />
                </n-form-item>
              </n-form>
            </div>
          </div>
        </n-tab-pane>

        <!-- 格式模板 -->
        <n-tab-pane name="formats" :tab="$t('settings.formatTemplates')">
          <div class="tab-content">
            <div class="content-card">
              <div class="format-section">
                <div class="format-header">
                  <div class="format-title">
                    <n-icon :component="CopyOutline" :size="18" />
                    <span>{{ $t('settings.copyFormat') }}</span>
                  </div>
                  <div class="format-vars">
                    <span class="vars-label">{{ $t('settings.availableVars') }}:</span>
                    <code>{account}</code>
                    <code>{password}</code>
                    <code>{type}</code>
                    <code>{expireAt}</code>
                  </div>
                </div>
                <n-input
                  v-model:value="form.copyFormat"
                  type="textarea"
                  :rows="4"
                  :placeholder="$t('settings.copyTemplatePlaceholder')"
                />
              </div>

              <n-divider style="margin: 24px 0;" />

              <div class="format-section">
                <div class="format-header">
                  <div class="format-title">
                    <n-icon :component="MailOutline" :size="18" />
                    <span>{{ $t('settings.emailFormat') }}</span>
                  </div>
                  <div class="format-vars">
                    <span class="vars-label">{{ $t('settings.availableVars') }}:</span>
                    <code>{account}</code>
                    <code>{expireAt}</code>
                  </div>
                </div>
                <n-input
                  v-model:value="form.emailFormat"
                  type="textarea"
                  :rows="4"
                  :placeholder="$t('settings.emailTemplatePlaceholder')"
                />
              </div>
            </div>
          </div>
        </n-tab-pane>

        <!-- 标签管理 -->
        <n-tab-pane name="tags" :tab="$t('settings.tagManagement')">
          <div class="tab-content">
            <div class="content-card">
              <div class="tag-management">
                <!-- 账号类型 -->
                <div class="tag-category">
                  <div class="tag-category-header">
                    <span class="tag-category-title">{{ $t('settings.accountTypes') }}</span>
                    <n-button size="small" type="primary" @click="openAddType">
                      <template #icon><n-icon :component="AddOutline" /></template>
                      {{ $t('settings.addTag') }}
                    </n-button>
                  </div>
                  <div class="tag-list">
                    <div class="tag-item" v-for="(tag, index) in accountTypes" :key="index" @click="editType(index)">
                      <n-tag :color="{ color: tag.color, textColor: '#fff' }" size="small" round>{{ tag.label }}</n-tag>
                      <div class="tag-item-actions" @click.stop>
                        <n-popconfirm @positive-click="deleteType(index)">
                          <template #trigger>
                            <n-button size="tiny" quaternary type="error">
                              <template #icon><n-icon :component="TrashOutline" :size="14" /></template>
                            </n-button>
                          </template>
                          {{ $t('settings.confirmDeleteTag') }}
                        </n-popconfirm>
                      </div>
                    </div>
                  </div>
                </div>

                <n-divider style="margin: 24px 0;" />

                <!-- 账号状态 -->
                <div class="tag-category">
                  <div class="tag-category-header">
                    <span class="tag-category-title">{{ $t('settings.accountStatuses') }}</span>
                    <n-button size="small" type="primary" @click="openAddStatus">
                      <template #icon><n-icon :component="AddOutline" /></template>
                      {{ $t('settings.addTag') }}
                    </n-button>
                  </div>
                  <div class="tag-list">
                    <div class="tag-item" v-for="(tag, index) in accountStatuses" :key="index" @click="editStatus(index)">
                      <n-tag :color="{ color: tag.color, textColor: '#fff' }" size="small" round>{{ tag.label }}</n-tag>
                      <div class="tag-item-actions" @click.stop>
                        <n-popconfirm @positive-click="deleteStatus(index)">
                          <template #trigger>
                            <n-button size="tiny" quaternary type="error">
                              <template #icon><n-icon :component="TrashOutline" :size="14" /></template>
                            </n-button>
                          </template>
                          {{ $t('settings.confirmDeleteTag') }}
                        </n-popconfirm>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </n-tab-pane>

        <!-- 保存按钮 -->
        <template #suffix>
          <div v-if="savingSystem || savingFormats || savingTags" class="saving-indicator">
            <n-spin size="small" />
            <span style="margin-left: 8px; font-size: 13px; color: #666;">{{ $t('common.saving') }}</span>
          </div>
        </template>
      </n-tabs>
    </n-spin>

    <!-- 标签编辑弹窗 -->
    <n-modal v-model:show="showTagModal" preset="card" :title="tagModalTitle" style="width:420px" :bordered="false">
      <n-form :model="tagForm" label-placement="left" label-width="80" size="medium">
        <n-form-item :label="$t('settings.tagLabel')">
          <n-input v-model:value="tagForm.label" :placeholder="$t('settings.tagLabelPlaceholder')" />
        </n-form-item>
        <n-form-item :label="$t('settings.tagValue')">
          <n-input v-model:value="tagForm.value" :placeholder="$t('settings.tagValuePlaceholder')" />
        </n-form-item>
        <n-form-item :label="$t('settings.tagColor')">
          <n-color-picker v-model:value="tagForm.color" :modes="['hex']" :show-alpha="false" />
        </n-form-item>
        <n-form-item :label="$t('settings.preview')">
          <n-tag :color="{ color: tagForm.color, textColor: '#fff' }" size="medium" round>{{ tagForm.label || $t('settings.previewText') }}</n-tag>
        </n-form-item>
      </n-form>
      <template #footer>
        <n-space justify="end">
          <n-button size="small" @click="showTagModal = false">{{ $t('common.cancel') }}</n-button>
          <n-button size="small" type="primary" @click="saveTag">{{ $t('common.save') }}</n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { NForm, NFormItem, NInputNumber, NButton, NSpin, NIcon, NInput, useMessage, NTag, NModal, NSelect, NTabs, NTabPane, NSpace, NPopconfirm, NDivider, NColorPicker } from 'naive-ui'
import { CopyOutline, SaveOutline, MailOutline, AddOutline, CreateOutline, TrashOutline, CheckmarkOutline } from '@vicons/ionicons5'
import { GetSystemConfig, UpdateSystemConfig } from '../../wailsjs/go/main/App'

const { t: $t } = useI18n()
const message = useMessage()
const loading = ref(false)
const savingSystem = ref(false)
const savingFormats = ref(false)
const savingTags = ref(false)
const currentTab = ref('system')
const autoSaveTimer = ref<number | null>(null)
const form = ref({
  defaultValidityDays: 30,
  reminderDaysBefore: 1,
  copyFormat: '账号：{account}\n密码：{password}',
  emailFormat: '您的账号 {account} 将在 {expireAt} 过期，请及时处理。'
})
const originalForm = ref({
  defaultValidityDays: 30,
  reminderDaysBefore: 1,
  copyFormat: '账号：{account}\n密码：{password}',
  emailFormat: '您的账号 {account} 将在 {expireAt} 过期，请及时处理。'
})

// 标签管理
const accountTypes = ref([
  { label: 'PLUS', value: 'PLUS', color: '#2080f0' },
  { label: 'BUSINESS', value: 'BUSINESS', color: '#18a058' },
  { label: 'FREE', value: 'FREE', color: '#909399' }
])
const accountStatuses = ref([
  { label: '未售出', value: 'unsold', color: '#f0a020' },
  { label: '已售出', value: 'sold', color: '#18a058' }
])
const originalAccountTypes = ref<any[]>([])
const originalAccountStatuses = ref<any[]>([])
const showTagModal = ref(false)
const tagForm = ref({ label: '', value: '', color: '#2080f0' })
const tagModalTitle = ref('')
const editingType = ref<'type' | 'status' | null>(null)
const editingIndex = ref(-1)

const hasSystemChanges = computed(() => {
  return form.value.defaultValidityDays !== originalForm.value.defaultValidityDays ||
    form.value.reminderDaysBefore !== originalForm.value.reminderDaysBefore
})

const hasFormatChanges = computed(() => {
  return form.value.copyFormat !== originalForm.value.copyFormat ||
    form.value.emailFormat !== originalForm.value.emailFormat
})

const hasTagChanges = computed(() => {
  return JSON.stringify(accountTypes.value) !== JSON.stringify(originalAccountTypes.value) ||
    JSON.stringify(accountStatuses.value) !== JSON.stringify(originalAccountStatuses.value)
})

const hasAnyChanges = computed(() => {
  return hasSystemChanges.value || hasFormatChanges.value || hasTagChanges.value
})

// 自动保存函数
function triggerAutoSave() {
  if (autoSaveTimer.value) {
    clearTimeout(autoSaveTimer.value)
  }
  autoSaveTimer.value = window.setTimeout(async () => {
    if (hasSystemChanges.value) {
      await doSaveSystem()
    }
    if (hasFormatChanges.value) {
      await doSaveFormats()
    }
    if (hasTagChanges.value) {
      await doSaveTags()
    }
  }, 300)
}

// 监听表单变化
watch(() => form.value.defaultValidityDays, () => {
  if (!loading.value) triggerAutoSave()
})

watch(() => form.value.reminderDaysBefore, () => {
  if (!loading.value) triggerAutoSave()
})

watch(() => form.value.copyFormat, () => {
  if (!loading.value) triggerAutoSave()
})

watch(() => form.value.emailFormat, () => {
  if (!loading.value) triggerAutoSave()
})

watch(() => accountTypes.value, () => {
  if (!loading.value) triggerAutoSave()
}, { deep: true })

watch(() => accountStatuses.value, () => {
  if (!loading.value) triggerAutoSave()
}, { deep: true })

onMounted(async () => {
  loading.value = true
  try {
    const c = await GetSystemConfig()
    if (c) {
      form.value = {
        defaultValidityDays: c.defaultValidityDays,
        reminderDaysBefore: c.reminderDaysBefore,
        copyFormat: c.copyFormat || '账号：{account}\n密码：{password}',
        emailFormat: c.emailFormat || '您的账号 {account} 将在 {expireAt} 过期，请及时处理。'
      }
      originalForm.value = { ...form.value }

      // 加载标签配置
      if (c.accountTypes) {
        try {
          accountTypes.value = JSON.parse(c.accountTypes)
          originalAccountTypes.value = JSON.parse(c.accountTypes)
        } catch {}
      }
      if (c.accountStatuses) {
        try {
          accountStatuses.value = JSON.parse(c.accountStatuses)
          originalAccountStatuses.value = JSON.parse(c.accountStatuses)
        } catch {}
      }
    }
  } catch {}
  loading.value = false
})

async function saveCurrentTab() {
  if (hasSystemChanges.value) {
    await doSaveSystem()
  }
  if (hasFormatChanges.value) {
    await doSaveFormats()
  }
  if (hasTagChanges.value) {
    await doSaveTags()
  }
}

async function doSaveSystem() {
  savingSystem.value = true
  try {
    await UpdateSystemConfig(
      form.value.defaultValidityDays,
      form.value.reminderDaysBefore,
      form.value.copyFormat,
      form.value.emailFormat,
      JSON.stringify(accountTypes.value),
      JSON.stringify(accountStatuses.value)
    )
    message.success($t('settings.settingsSaved'))
    originalForm.value.defaultValidityDays = form.value.defaultValidityDays
    originalForm.value.reminderDaysBefore = form.value.reminderDaysBefore
  }
  catch (e: any) { message.error(e.toString()) }
  savingSystem.value = false
}

async function doSaveFormats() {
  savingFormats.value = true
  try {
    await UpdateSystemConfig(
      form.value.defaultValidityDays,
      form.value.reminderDaysBefore,
      form.value.copyFormat,
      form.value.emailFormat,
      JSON.stringify(accountTypes.value),
      JSON.stringify(accountStatuses.value)
    )
    message.success($t('settings.settingsSaved'))
    originalForm.value.copyFormat = form.value.copyFormat
    originalForm.value.emailFormat = form.value.emailFormat
  }
  catch (e: any) { message.error(e.toString()) }
  savingFormats.value = false
}

async function doSaveTags() {
  savingTags.value = true
  try {
    await UpdateSystemConfig(
      form.value.defaultValidityDays,
      form.value.reminderDaysBefore,
      form.value.copyFormat,
      form.value.emailFormat,
      JSON.stringify(accountTypes.value),
      JSON.stringify(accountStatuses.value)
    )
    message.success($t('settings.settingsSaved'))
    originalAccountTypes.value = JSON.parse(JSON.stringify(accountTypes.value))
    originalAccountStatuses.value = JSON.parse(JSON.stringify(accountStatuses.value))
  }
  catch (e: any) { message.error(e.toString()) }
  savingTags.value = false
}

// 标签管理函数
function openAddType() {
  tagForm.value = { label: '', value: '', color: 'default' }
  tagModalTitle.value = $t('settings.addAccountType')
  editingType.value = 'type'
  editingIndex.value = -1
  showTagModal.value = true
}

function openAddStatus() {
  tagForm.value = { label: '', value: '', color: 'default' }
  tagModalTitle.value = $t('settings.addAccountStatus')
  editingType.value = 'status'
  editingIndex.value = -1
  showTagModal.value = true
}

function editType(index: number) {
  tagForm.value = { ...accountTypes.value[index] }
  tagModalTitle.value = $t('settings.editAccountType')
  editingType.value = 'type'
  editingIndex.value = index
  showTagModal.value = true
}

function editStatus(index: number) {
  tagForm.value = { ...accountStatuses.value[index] }
  tagModalTitle.value = $t('settings.editAccountStatus')
  editingType.value = 'status'
  editingIndex.value = index
  showTagModal.value = true
}

function deleteType(index: number) {
  accountTypes.value.splice(index, 1)
  message.success($t('settings.tagDeleted'))
}

function deleteStatus(index: number) {
  accountStatuses.value.splice(index, 1)
  message.success($t('settings.tagDeleted'))
}

function saveTag() {
  if (!tagForm.value.label || !tagForm.value.value) {
    message.warning($t('settings.tagRequired'))
    return
  }

  if (editingType.value === 'type') {
    if (editingIndex.value >= 0) {
      accountTypes.value[editingIndex.value] = { ...tagForm.value }
    } else {
      accountTypes.value.push({ ...tagForm.value })
    }
  } else if (editingType.value === 'status') {
    if (editingIndex.value >= 0) {
      accountStatuses.value[editingIndex.value] = { ...tagForm.value }
    } else {
      accountStatuses.value.push({ ...tagForm.value })
    }
  }

  message.success($t('settings.tagSaved'))
  showTagModal.value = false
}
</script>

<style scoped>
.settings-container {
  /* Counteract the main container's padding to make tabbar full width */
  margin: -20px -24px;
  height: calc(100vh - 52px);
  display: flex;
  flex-direction: column;
  background: #F8F6F4;
  overflow: hidden;
}

.settings-tabs {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.settings-tabs :deep(.n-tabs-nav-scroll-wrapper) {
  flex-shrink: 0;
}

.settings-tabs :deep(.n-tabs-nav) {
  background: #FFFFFF;
  padding: 0;
  margin: 0 !important;
  box-shadow: none;
  border-bottom: none;
  height: 44px;
  position: relative;
}

.settings-tabs :deep(.n-tabs-nav-scroll-content) {
  padding: 0 24px;
  height: 44px;
  display: flex;
  align-items: center;
}

.settings-tabs :deep(.n-tabs-tab) {
  font-size: 13px;
  font-weight: 400;
  padding: 0 16px;
  height: 44px;
  color: #666;
}

.settings-tabs :deep(.n-tabs-tab--active) {
  color: #7C9885;
  font-weight: 500;
}

.settings-tabs :deep(.n-tabs-bar) {
  display: none;
}

.settings-tabs :deep(.n-tabs-pane-wrapper) {
  flex: 1;
  overflow-y: auto;
  background: #F8F6F4;
  min-height: 0;
}

.settings-tabs :deep(.n-tabs-nav__suffix) {
  display: flex;
  align-items: center;
  padding-right: 24px;
}

.saving-indicator {
  display: flex;
  align-items: center;
  padding: 0 16px;
}

.tab-content {
  max-width: 900px;
  margin: 0 auto;
  padding: 24px;
}

.content-card {
  background: white;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.06);
  border: 1px solid #E8E4E0;
}

/* 格式模板样式 */
.format-section {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.format-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 12px;
  margin-bottom: 4px;
}

.format-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  font-weight: 600;
  color: #333;
}

.format-vars {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.vars-label {
  font-size: 12px;
  color: #999;
}

.format-vars code {
  background: #f5f5f5;
  padding: 2px 8px;
  border-radius: 4px;
  color: #7C9885;
  font-family: 'Consolas', 'Monaco', monospace;
  font-size: 12px;
  font-weight: 500;
}

/* 标签管理样式 */
.tag-management {
  display: flex;
  flex-direction: column;
}

.tag-category {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.tag-category-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.tag-category-title {
  font-size: 14px;
  font-weight: 600;
  color: #333;
}

.tag-list {
  display: grid;
  grid-template-columns: repeat(6, 1fr);
  gap: 10px;
}

.tag-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 12px;
  background: #fafafa;
  border: 1px solid #e8e8e8;
  border-radius: 6px;
  transition: all 0.2s;
  cursor: pointer;
}

.tag-item:hover {
  background: #f0f0f0;
  border-color: #d9d9d9;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.08);
  transform: translateY(-1px);
}

.tag-item-actions {
  display: flex;
  gap: 4px;
  opacity: 0;
  transition: opacity 0.2s;
}

.tag-item:hover .tag-item-actions {
  opacity: 1;
}

@media (max-width: 768px) {
  .tab-content {
    padding: 16px;
  }

  .format-header {
    flex-direction: column;
    align-items: flex-start;
  }

  .tag-list {
    grid-template-columns: repeat(3, 1fr);
  }
}
</style>
