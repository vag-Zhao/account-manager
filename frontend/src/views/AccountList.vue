<template>
  <div class="account-list-container">
    <div class="card">
      <!-- 工具栏 -->
      <div class="toolbar">
        <div class="filters">
          <n-select v-model:value="filterType" :placeholder="$t('accounts.type')" :options="typeOpts" size="small" style="width:100px" clearable />
          <n-select v-model:value="filterSold" :placeholder="$t('accounts.status')" :options="soldOpts" size="small" style="width:100px" clearable />
          <n-input v-model:value="searchText" :placeholder="$t('accounts.search')" size="small" style="width:160px" clearable @keyup.enter="doSearch">
            <template #prefix><n-icon :component="SearchOutline" /></template>
          </n-input>
          <n-button size="small" type="primary" @click="doSearch">{{ $t('accounts.searchBtn') }}</n-button>
          <n-button size="small" quaternary @click="doReset">{{ $t('accounts.reset') }}</n-button>
        </div>
        <div class="actions">
          <n-button size="small" type="primary" @click="openAdd">
            <template #icon><n-icon :component="AddOutline" /></template>
            {{ $t('accounts.add') }}
          </n-button>
          <n-button size="small" @click="openBatch">
            <template #icon><n-icon :component="CloudUploadOutline" /></template>
            {{ $t('accounts.import') }}
          </n-button>
        </div>
      </div>

      <!-- 表格 -->
      <div class="table-wrapper">
        <n-data-table
          :columns="columns"
          :data="store.accounts"
          :loading="store.loading"
          :pagination="false"
          size="small"
          :row-key="(r: Account) => r.id"
          :bordered="false"
        />
      </div>

      <!-- 分页 -->
      <div class="pager">
        <n-pagination
          v-if="store.totalPages > 1"
          v-model:page="page"
          :page-count="store.totalPages"
          size="small"
          @update:page="p => store.setPage(p)"
        />
      </div>
    </div>

    <!-- 添加弹窗 -->
    <n-modal v-model:show="showAdd" preset="card" :title="$t('accounts.addAccount')" style="width:420px" :bordered="false">
      <n-form ref="addFormRef" :model="addForm" :rules="addRules" label-placement="left" label-width="70" size="small">
        <n-form-item :label="$t('accounts.account')" path="account"><n-input v-model:value="addForm.account" :placeholder="$t('accounts.accountPlaceholder')" /></n-form-item>
        <n-form-item :label="$t('accounts.password')"><n-input v-model:value="addForm.password" type="password" show-password-on="click" :placeholder="$t('accounts.passwordPlaceholder')" /></n-form-item>
        <n-form-item :label="$t('accounts.accountType')"><n-select v-model:value="addForm.accountType" :options="typeOpts" /></n-form-item>
        <n-form-item :label="$t('accounts.status')"><n-select v-model:value="addForm.isSold" :options="soldSelectOpts" /></n-form-item>
        <n-form-item :label="$t('accounts.expireAt')" v-if="addForm.accountType !== 'FREE'">
          <n-date-picker v-model:value="addForm.expireAt" type="date" clearable style="width:100%" :placeholder="$t('accounts.expirePlaceholder')" />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-space justify="end">
          <n-button size="small" @click="showAdd = false">{{ $t('accounts.cancel') }}</n-button>
          <n-button size="small" type="primary" @click="doAdd" :loading="adding">{{ $t('accounts.add') }}</n-button>
        </n-space>
      </template>
    </n-modal>

    <!-- 编辑弹窗 -->
    <n-modal v-model:show="showEdit" preset="card" :title="$t('accounts.editAccount')" style="width:420px" :bordered="false">
      <n-form :model="editForm" label-placement="left" label-width="70" size="small">
        <n-form-item :label="$t('accounts.account')"><n-input v-model:value="editForm.account" /></n-form-item>
        <n-form-item :label="$t('accounts.password')"><n-input v-model:value="editForm.password" type="password" show-password-on="click" /></n-form-item>
        <n-form-item :label="$t('accounts.accountType')"><n-select v-model:value="editForm.accountType" :options="typeOpts" /></n-form-item>
        <n-form-item :label="$t('accounts.status')"><n-select v-model:value="editForm.isSold" :options="soldSelectOpts" /></n-form-item>
        <n-form-item :label="$t('accounts.expireAt')" v-if="editForm.accountType !== 'FREE'">
          <n-date-picker v-model:value="editForm.expireAt" type="date" clearable style="width:100%" />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-space justify="end">
          <n-button size="small" @click="showEdit = false">{{ $t('accounts.cancel') }}</n-button>
          <n-button size="small" type="primary" @click="saveEdit" :loading="saving">{{ $t('accounts.save') }}</n-button>
        </n-space>
      </template>
    </n-modal>

    <!-- 批量导入弹窗 -->
    <n-modal v-model:show="showBatch" preset="card" :title="$t('accounts.batchImport')" style="width:500px" :bordered="false">
      <n-alert type="info" style="margin-bottom:16px">
        <template #header>{{ $t('accounts.importFormat') }}</template>
        {{ $t('accounts.importFormatDesc') }}
      </n-alert>
      <n-input
        v-model:value="batchText"
        type="textarea"
        :rows="8"
        placeholder="user1@example.com,password123,PLUS,2024-12-31,unsold
user2@example.com,pass456,BUSINESS,,sold
user3@example.com,pass789,FREE,,"
      />
      <n-alert v-if="batchResult" :type="batchResult.errors.length ? 'warning' : 'success'" style="margin-top:12px">
        <template #header>{{ $t('accounts.importResult') }}</template>
        {{ $t('accounts.importSuccess', { n: batchResult.success }) }}
        <div v-if="batchResult.errors.length" style="margin-top:8px;font-size:12px">
          {{ $t('accounts.importFailed') }}: {{ batchResult.errors.slice(0, 3).join('; ') }}
          <span v-if="batchResult.errors.length > 3">{{ $t('common.etc') }} {{ batchResult.errors.length }} {{ $t('common.items') }}</span>
        </div>
      </n-alert>
      <template #footer>
        <n-space justify="end">
          <n-button size="small" @click="showBatch = false">{{ $t('accounts.close') }}</n-button>
          <n-button size="small" type="primary" @click="doBatch" :loading="importing">{{ $t('accounts.startImport') }}</n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, h, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { NSpace, NSelect, NInput, NButton, NDataTable, NPagination, NModal, NForm, NFormItem, NDatePicker, NTag, NPopconfirm, NIcon, NAlert, NTooltip, useMessage } from 'naive-ui'
import { SearchOutline, AddOutline, CreateOutline, TrashOutline, PricetagOutline, CloudUploadOutline, CopyOutline } from '@vicons/ionicons5'
import type { DataTableColumns, FormInst, FormRules } from 'naive-ui'
import { useAccountStore } from '../stores/accountStore'
import type { Account, Tag } from '../types'

const { t: $t } = useI18n()
const message = useMessage()
const store = useAccountStore()

const filterType = ref('')
const filterSold = ref('')
const searchText = ref('')
const page = ref(1)
const copyFormat = ref('账号：{account}\n密码：{password}')

// 标签配置
const accountTypes = ref<Tag[]>([])
const accountStatuses = ref<Tag[]>([])

// 添加相关
const showAdd = ref(false)
const adding = ref(false)
const addFormRef = ref<FormInst | null>(null)
const addForm = ref({ account: '', password: '', accountType: 'PLUS', isSold: 'unsold', expireAt: null as number | null })

// 编辑相关
const showEdit = ref(false)
const saving = ref(false)
const editForm = ref({ id: 0, account: '', password: '', accountType: 'PLUS', isSold: 'unsold', expireAt: null as number | null })

// 批量导入相关
const showBatch = ref(false)
const importing = ref(false)
const batchText = ref('')
const batchResult = ref<{ success: number; errors: string[] } | null>(null)

const typeOpts = computed(() => accountTypes.value.map(t => ({ label: t.label, value: t.value })))
const soldOpts = computed(() => accountStatuses.value.map(s => ({ label: s.label, value: s.value })))
const soldSelectOpts = computed(() => accountStatuses.value.map(s => ({ label: s.label, value: s.value })))

const addRules = computed<FormRules>(() => ({
  account: { required: true, message: $t('accounts.accountPlaceholder'), trigger: 'blur' }
}))

const columns = computed<DataTableColumns<Account>>(() => [
  {
    title: $t('accounts.account'), key: 'account', width: 200, titleAlign: 'center',
    render: r => h('div', { style: 'display:flex;align-items:center;gap:4px' }, [
      h(NTooltip, { trigger: 'hover', disabled: r.account.length < 20 }, {
        trigger: () => h('span', { style: 'flex:1;min-width:0;overflow:hidden;text-overflow:ellipsis;white-space:nowrap' }, r.account),
        default: () => r.account
      }),
      h(NTooltip, { trigger: 'hover' }, {
        trigger: () => h(NButton, {
          size: 'tiny',
          quaternary: true,
          style: 'flex-shrink:0',
          onClick: (e: Event) => { e.stopPropagation(); copyField(r, 'account') }
        }, { icon: () => h(NIcon, { component: CopyOutline, size: 14 }) }),
        default: () => $t('accounts.copyAccount')
      })
    ])
  },
  {
    title: $t('accounts.password'), key: 'password', width: 130, titleAlign: 'center',
    render: r => h('div', { style: 'padding-left:35px' }, [
      h(NSpace, { size: 4, align: 'center', style: 'flex-wrap:nowrap' }, () => [
        h('code', { style: 'font-size:11px;color:#666;max-width:60px;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;display:inline-block' }, r.password || '-'),
        r.password ? h(NTooltip, { trigger: 'hover' }, {
          trigger: () => h(NButton, {
            size: 'tiny',
            quaternary: true,
            onClick: (e: Event) => { e.stopPropagation(); copyField(r, 'password') }
          }, { icon: () => h(NIcon, { component: CopyOutline, size: 14 }) }),
          default: () => $t('accounts.copyPassword')
        }) : null
      ])
    ])
  },
  {
    title: $t('accounts.accountType'), key: 'accountType', width: 90, titleAlign: 'center',
    render: r => {
      const tag = accountTypes.value.find(t => t.value === r.accountType)
      return h('div', { style: 'padding-left:20px' }, [
        h(NTag, {
          size: 'small',
          round: true,
          color: tag?.color ? { color: tag.color, textColor: '#fff' } : undefined,
          type: tag?.color ? undefined : 'default'
        }, () => tag?.label || r.accountType)
      ])
    }
  },
  {
    title: $t('accounts.status'), key: 'isSold', width: 75, titleAlign: 'center',
    render: r => {
      const statusValue = r.isSold ? 'sold' : 'unsold'
      const tag = accountStatuses.value.find(s => s.value === statusValue)
      return h('div', { style: 'padding-left:15px' }, [
        h(NTag, {
          size: 'small',
          round: true,
          color: tag?.color ? { color: tag.color, textColor: '#fff' } : undefined,
          type: tag?.color ? undefined : 'default'
        }, () => tag?.label || (r.isSold ? $t('accounts.sold') : $t('accounts.unsold')))
      ])
    }
  },
  {
    title: $t('accounts.expireAt'), key: 'expireAt', width: 110, titleAlign: 'center',
    render: r => {
      if (!r.expireAt) return h('div', { style: 'padding-left:20px' }, '-')
      const d = new Date(r.expireAt)
      const isExpired = d < new Date()
      const locale = $t('common.locale') === 'zh-CN' ? 'zh-CN' : 'en-US'
      return h('span', { style: isExpired ? 'color:#cf1322;font-size:12px;padding-left:20px' : 'font-size:12px;padding-left:20px' }, d.toLocaleDateString(locale))
    }
  },
  {
    title: $t('accounts.actions'), key: 'actions', width: 140, fixed: 'right', titleAlign: 'center',
    render: r => h('div', { style: 'margin-left:-20px' }, [
      h(NSpace, { size: 4 }, () => [
        h(NTooltip, { trigger: 'hover' }, {
          trigger: () => h(NButton, {
            size: 'tiny',
            quaternary: true,
            type: 'info',
            onClick: () => copyAccount(r)
          }, { icon: () => h(NIcon, { component: CopyOutline }) }),
          default: () => $t('accounts.oneClickCopy')
        }),
        h(NButton, { size: 'tiny', quaternary: true, onClick: () => openEdit(r) }, { icon: () => h(NIcon, { component: CreateOutline }) }),
        h(NButton, {
          size: 'tiny',
          quaternary: true,
          type: r.isSold ? 'default' : 'primary',
          onClick: () => toggleSold(r)
        }, { icon: () => h(NIcon, { component: PricetagOutline }) }),
        h(NPopconfirm, { onPositiveClick: () => doDelete(r.id) }, {
          trigger: () => h(NButton, { size: 'tiny', quaternary: true, type: 'error' }, { icon: () => h(NIcon, { component: TrashOutline }) }),
          default: () => $t('accounts.confirmDelete')
        })
      ])
    ])
  }
])

onMounted(async () => {
  // 加载系统配置和标签
  const config = await store.loadSystemConfig()
  if (config) {
    if (config.copyFormat) {
      copyFormat.value = config.copyFormat
    }
    // 加载账号类型标签
    if (config.accountTypes) {
      try {
        accountTypes.value = JSON.parse(config.accountTypes)
      } catch {
        // 使用默认值
        accountTypes.value = [
          { label: 'PLUS', value: 'PLUS', color: '#2080f0' },
          { label: 'BUSINESS', value: 'BUSINESS', color: '#18a058' },
          { label: 'FREE', value: 'FREE', color: '#909399' }
        ]
      }
    }
    // 加载账号状态标签
    if (config.accountStatuses) {
      try {
        accountStatuses.value = JSON.parse(config.accountStatuses)
      } catch {
        // 使用默认值
        accountStatuses.value = [
          { label: $t('accounts.unsold'), value: 'unsold', color: '#f0a020' },
          { label: $t('accounts.sold'), value: 'sold', color: '#18a058' }
        ]
      }
    }
  }
  store.fetchAccounts()
})

// 复制单个字段
function copyField(r: Account, field: 'account' | 'password') {
  const text = field === 'account' ? r.account : r.password
  navigator.clipboard.writeText(text).then(() => {
    message.success(field === 'account' ? $t('accounts.accountCopied') : $t('accounts.passwordCopied'))
  }).catch(() => {
    message.error($t('accounts.copyFailed'))
  })
}

// 一键复制（使用自定义格式）
function copyAccount(r: Account) {
  const expireDate = r.expireAt ? new Date(r.expireAt).toLocaleDateString($t('common.locale') === 'zh-CN' ? 'zh-CN' : 'en-US') : ''
  const text = copyFormat.value
    .replace(/{account}/g, r.account)
    .replace(/{password}/g, r.password || '')
    .replace(/{type}/g, r.accountType)
    .replace(/{expireAt}/g, expireDate)
  navigator.clipboard.writeText(text).then(() => {
    message.success($t('accounts.copySuccess'))
  }).catch(() => {
    message.error($t('accounts.copyFailed'))
  })
}

function doSearch() {
  const sold = filterSold.value === '' ? null : filterSold.value === 'true'
  store.setFilter(filterType.value, sold, searchText.value)
}

function doReset() {
  filterType.value = ''
  filterSold.value = ''
  searchText.value = ''
  store.resetFilters()
}

async function openAdd() {
  // 重新加载系统配置以获取最新的标签
  await loadTags()
  addForm.value = { account: '', password: '', accountType: 'PLUS', isSold: 'unsold', expireAt: null }
  showAdd.value = true
}

async function doAdd() {
  try { await addFormRef.value?.validate() } catch { return }
  adding.value = true
  try {
    const exp = addForm.value.expireAt ? new Date(addForm.value.expireAt).toISOString().split('T')[0] : ''
    await store.createAccount(addForm.value.account, addForm.value.password, addForm.value.accountType, exp, addForm.value.isSold === 'sold')
    message.success($t('accounts.addSuccess'))
    showAdd.value = false
  } catch (e: any) { message.error(e.toString()) }
  adding.value = false
}

async function openEdit(r: Account) {
  // 重新加载系统配置以获取最新的标签
  await loadTags()
  editForm.value = { id: r.id, account: r.account, password: r.password, accountType: r.accountType, isSold: r.isSold ? 'sold' : 'unsold', expireAt: r.expireAt ? new Date(r.expireAt).getTime() : null }
  showEdit.value = true
}

async function loadTags() {
  const config = await store.loadSystemConfig()
  if (config) {
    if (config.accountTypes) {
      try {
        accountTypes.value = JSON.parse(config.accountTypes)
      } catch {}
    }
    if (config.accountStatuses) {
      try {
        accountStatuses.value = JSON.parse(config.accountStatuses)
      } catch {}
    }
  }
}

async function saveEdit() {
  saving.value = true
  try {
    const exp = editForm.value.expireAt ? new Date(editForm.value.expireAt).toISOString().split('T')[0] : ''
    await store.updateAccount(editForm.value.id, editForm.value.account, editForm.value.password, editForm.value.accountType, exp, editForm.value.isSold === 'sold')
    message.success($t('accounts.saveSuccess'))
    showEdit.value = false
  } catch (e: any) { message.error(e.toString()) }
  saving.value = false
}

async function toggleSold(r: Account) {
  try {
    r.isSold ? await store.markAsUnsold(r.id) : await store.markAsSold(r.id)
    message.success($t('accounts.updateSuccess'))
  } catch (e: any) { message.error(e.toString()) }
}

async function doDelete(id: number) {
  try {
    await store.deleteAccount(id)
    message.success($t('accounts.deleteSuccess'))
  } catch (e: any) { message.error(e.toString()) }
}

function openBatch() {
  batchText.value = ''
  batchResult.value = null
  showBatch.value = true
}

async function doBatch() {
  if (!batchText.value.trim()) { message.warning($t('accounts.inputData')); return }
  const lines = batchText.value.trim().split('\n')
  const list = lines.map(l => {
    const p = l.split(',').map(s => s.trim())
    const isSold = p[4]?.toLowerCase() === 'sold'
    return { account: p[0], password: p[1] || '', accountType: p[2] || 'PLUS', expireAt: p[3] || '', isSold }
  }).filter(x => x.account)
  if (!list.length) { message.warning($t('accounts.noData')); return }
  importing.value = true
  try {
    const r = await store.batchImport(list)
    batchResult.value = { success: r.success, errors: r.errors || [] }
    if (r.success > 0) message.success($t('accounts.importSuccess', { n: r.success }))
  } catch (e: any) { message.error(e.toString()) }
  importing.value = false
}
</script>

<style scoped>
.account-list-container {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.card {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
  flex-shrink: 0;
}

.filters {
  display: flex;
  gap: 8px;
  align-items: center;
  flex-wrap: wrap;
}

.actions {
  display: flex;
  gap: 8px;
}

.table-wrapper {
  flex: 1;
  margin-top: 16px;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.table-wrapper :deep(.n-data-table) {
  height: 100%;
}

.table-wrapper :deep(.n-data-table-wrapper) {
  height: 100%;
  overflow-y: auto;
}

@media (max-width: 768px) {
  .toolbar {
    flex-direction: column;
    align-items: stretch;
  }

  .filters {
    width: 100%;
  }

  .filters > * {
    flex: 1;
    min-width: 80px;
  }

  .actions {
    justify-content: flex-end;
  }
}

.pager {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
  height: 32px;
  flex-shrink: 0;
  align-items: center;
}
</style>
