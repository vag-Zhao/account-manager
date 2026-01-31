<template>
  <div class="dashboard-container">
    <!-- 动态标题 -->
    <div class="hero-header">
      <div class="hero-bg">
        <div class="hero-orb orb-1"></div>
        <div class="hero-orb orb-2"></div>
        <div class="hero-orb orb-3"></div>
      </div>
      <div class="hero-content">
        <div class="hero-badge">
          <n-icon :component="SparklesOutline" :size="14" />
          <span>{{ $t('dashboard.welcome') || '账号管理中心' }}</span>
        </div>
        <h1 class="hero-title">
          <span class="title-char" v-for="(char, i) in titleChars" :key="i" :style="{ animationDelay: `${i * 0.08}s` }">{{ char === ' ' ? '&nbsp;' : char }}</span>
        </h1>
        <p class="hero-desc">{{ $t('dashboard.subtitle') || '轻松管理您的所有账号，掌控每一个细节' }}</p>
      </div>
      <div class="hero-decoration">
        <div class="deco-line" v-for="n in 5" :key="n" :style="{ animationDelay: `${n * 0.2}s` }"></div>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-grid">
      <div class="stat-card" v-for="stat in statItems" :key="stat.label">
        <div class="stat-left">
          <div class="stat-icon" :style="{ backgroundColor: stat.bgColor }">
            <n-icon :component="stat.icon" :size="18" color="#fff" />
          </div>
          <div class="stat-label">{{ stat.label }}</div>
        </div>
        <div class="stat-value">{{ stat.value }}</div>
      </div>
    </div>

    <!-- 即将过期 -->
    <div class="card expire-card">
      <div class="card-header-with-action">
        <div class="card-title">{{ $t('dashboard.expiringAccounts') }}</div>
        <n-button size="tiny" type="primary" :loading="checking" @click="handleCheck" class="check-btn">
          <template #icon><n-icon :component="NotificationsOutline" :size="14" /></template>
          {{ $t('dashboard.checkExpiry') }}
        </n-button>
      </div>
      <div class="expire-list" v-if="expiringList.length">
        <div class="expire-item" v-for="acc in expiringList" :key="acc.id">
          <div class="expire-info">
            <span class="expire-account">{{ acc.account }}</span>
            <n-tag size="small" :type="acc.accountType === 'PLUS' ? 'info' : 'success'" round>{{ acc.accountType }}</n-tag>
          </div>
          <div class="expire-days" :class="getDaysClass(acc.expireAt)">
            {{ getDaysText(acc.expireAt) }}
          </div>
        </div>
      </div>
      <n-empty v-else :description="$t('dashboard.noExpiring')" size="small" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, markRaw } from 'vue'
import { useI18n } from 'vue-i18n'
import { NButton, NIcon, NTag, NEmpty, useMessage } from 'naive-ui'
import { ListOutline, NotificationsOutline, PeopleOutline, StarOutline, BriefcaseOutline, GiftOutline, CartOutline, TimeOutline, SparklesOutline } from '@vicons/ionicons5'
import { useAccountStore } from '../stores/accountStore'
import { ManualCheckExpiry } from '../../wailsjs/go/main/App'

const { t: $t } = useI18n()
const message = useMessage()
const store = useAccountStore()
const checking = ref(false)

const stats = computed(() => store.stats)

const titleChars = computed(() => {
  const title = $t('dashboard.heroTitle') || '数据概览'
  return title.split('')
})

const statItems = computed(() => [
  { label: $t('dashboard.stats.total'), value: stats.value?.total || 0, icon: markRaw(PeopleOutline), bgColor: '#7C9885' },
  { label: $t('dashboard.stats.plus'), value: stats.value?.plusCount || 0, icon: markRaw(StarOutline), bgColor: '#9EB3C2' },
  { label: $t('dashboard.stats.business'), value: stats.value?.businessCount || 0, icon: markRaw(BriefcaseOutline), bgColor: '#B8A99A' },
  { label: $t('dashboard.stats.free'), value: stats.value?.freeCount || 0, icon: markRaw(GiftOutline), bgColor: '#8BAF92' },
  { label: $t('dashboard.stats.sold'), value: stats.value?.soldCount || 0, icon: markRaw(CartOutline), bgColor: '#C4A77D' },
  { label: $t('dashboard.stats.expiring'), value: stats.value?.expiringIn7Days || 0, icon: markRaw(TimeOutline), bgColor: '#C9A5A5' }
])

const expiringList = computed(() => {
  return store.accounts.filter(a => {
    if (!a.expireAt || a.isSold) return false
    const days = Math.ceil((new Date(a.expireAt).getTime() - Date.now()) / 86400000)
    return days >= 0 && days <= 7
  }).slice(0, 5)
})

onMounted(() => {
  store.fetchStats()
  // 只在账号列表为空时才加载
  if (!store.accounts.length) {
    store.fetchAccounts()
  }
})

function getDaysText(d: string | null) {
  if (!d) return '-'
  const days = Math.ceil((new Date(d).getTime() - Date.now()) / 86400000)
  if (days === 0) return $t('dashboard.today')
  if (days === 1) return $t('dashboard.tomorrow')
  return `${days}${$t('dashboard.daysLater')}`
}

function getDaysClass(d: string | null) {
  if (!d) return ''
  const days = Math.ceil((new Date(d).getTime() - Date.now()) / 86400000)
  if (days <= 1) return 'urgent'
  if (days <= 3) return 'warning'
  return 'normal'
}

async function handleCheck() {
  checking.value = true
  try {
    const n = await ManualCheckExpiry()
    message.success(n > 0 ? $t('dashboard.reminderSent', { n }) : $t('dashboard.noReminder'))
  } catch (e: any) {
    message.error(e.toString())
  }
  checking.value = false
}
</script>

<style scoped>
.dashboard-container {
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.hero-header {
  position: relative;
  padding: 24px 32px;
  margin-bottom: 16px;
  border-radius: 16px;
  background: linear-gradient(135deg, #7C9885 0%, #9EB3C2 50%, #B8A99A 100%);
  background-size: 200% 200%;
  animation: gradientShift 8s ease infinite;
  overflow: hidden;
  flex-shrink: 0;
}

@keyframes gradientShift {
  0%, 100% { background-position: 0% 50%; }
  50% { background-position: 100% 50%; }
}

.hero-bg {
  position: absolute;
  inset: 0;
  overflow: hidden;
}

.hero-orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(50px);
  opacity: 0.5;
}

.orb-1 {
  width: 150px;
  height: 150px;
  background: radial-gradient(circle, #C4A77D 0%, transparent 70%);
  top: -40px;
  right: 5%;
  animation: orbFloat1 8s ease-in-out infinite;
}

.orb-2 {
  width: 100px;
  height: 100px;
  background: radial-gradient(circle, #8BAF92 0%, transparent 70%);
  bottom: -30px;
  left: 15%;
  animation: orbFloat2 10s ease-in-out infinite;
}

.orb-3 {
  width: 120px;
  height: 120px;
  background: radial-gradient(circle, #C9A5A5 0%, transparent 70%);
  top: 40%;
  right: 25%;
  animation: orbFloat3 12s ease-in-out infinite;
}

@keyframes orbFloat1 {
  0%, 100% { transform: translate(0, 0) scale(1); }
  25% { transform: translate(-20px, 15px) scale(1.1); }
  50% { transform: translate(-10px, -20px) scale(0.9); }
  75% { transform: translate(15px, 10px) scale(1.05); }
}

@keyframes orbFloat2 {
  0%, 100% { transform: translate(0, 0) scale(1); }
  33% { transform: translate(25px, -15px) scale(1.15); }
  66% { transform: translate(-15px, 20px) scale(0.95); }
}

@keyframes orbFloat3 {
  0%, 100% { transform: translate(0, 0) scale(1) rotate(0deg); }
  50% { transform: translate(30px, -25px) scale(1.1) rotate(180deg); }
}

.hero-content {
  position: relative;
  z-index: 2;
}

.hero-badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 14px;
  background: rgba(255,255,255,0.2);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255,255,255,0.3);
  border-radius: 20px;
  font-size: 12px;
  color: #fff;
  font-weight: 500;
  margin-bottom: 14px;
  animation: badgeSlideIn 0.8s cubic-bezier(0.34, 1.56, 0.64, 1) backwards;
}

@keyframes badgeSlideIn {
  0% { opacity: 0; transform: translateX(-30px) scale(0.8); }
  100% { opacity: 1; transform: translateX(0) scale(1); }
}

.hero-title {
  font-size: 28px;
  font-weight: 700;
  color: #fff;
  margin: 0 0 8px 0;
  text-shadow: 0 3px 15px rgba(0,0,0,0.2);
  display: flex;
  flex-wrap: wrap;
  line-height: 1.2;
}

.title-char {
  display: inline-block;
  animation: charReveal 0.5s cubic-bezier(0.34, 1.56, 0.64, 1) backwards;
}

@keyframes charReveal {
  0% {
    opacity: 0;
    transform: translateY(30px) rotateX(90deg) scale(0.5);
    filter: blur(10px);
  }
  100% {
    opacity: 1;
    transform: translateY(0) rotateX(0) scale(1);
    filter: blur(0);
  }
}

.hero-desc {
  font-size: 14px;
  color: rgba(255,255,255,0.9);
  margin: 0;
  animation: descFadeIn 0.8s ease-out 0.6s backwards;
  font-weight: 400;
}

@keyframes descFadeIn {
  0% { opacity: 0; transform: translateY(10px); }
  100% { opacity: 1; transform: translateY(0); }
}

.hero-decoration {
  position: absolute;
  right: 30px;
  top: 50%;
  transform: translateY(-50%);
  display: flex;
  gap: 10px;
  align-items: flex-end;
}

.deco-line {
  width: 5px;
  background: linear-gradient(180deg, rgba(255,255,255,0.5) 0%, rgba(255,255,255,0.1) 100%);
  border-radius: 3px;
  animation: barPulse 1.5s ease-in-out infinite;
}

.deco-line:nth-child(1) { height: 25px; animation-delay: 0s; }
.deco-line:nth-child(2) { height: 45px; animation-delay: 0.1s; }
.deco-line:nth-child(3) { height: 30px; animation-delay: 0.2s; }
.deco-line:nth-child(4) { height: 55px; animation-delay: 0.3s; }
.deco-line:nth-child(5) { height: 20px; animation-delay: 0.4s; }

@keyframes barPulse {
  0%, 100% {
    opacity: 0.4;
    transform: scaleY(1);
  }
  50% {
    opacity: 0.9;
    transform: scaleY(1.3);
  }
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(6, 1fr);
  gap: 12px;
  margin-bottom: 12px;
  flex-shrink: 0;
}

@media (max-width: 1200px) {
  .stats-grid {
    grid-template-columns: repeat(3, 1fr);
  }
}

@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 576px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }
}

.stat-card {
  background: #FFFFFF;
  border-radius: 12px;
  padding: 16px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border: 1px solid #E8E4E0;
  box-shadow: 0 2px 8px rgba(0,0,0,0.06);
  transition: all 0.2s;
}

.stat-card:hover {
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
  transform: translateY(-2px);
}

.stat-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.stat-icon {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.stat-label {
  font-size: 13px;
  color: #666;
  font-weight: 500;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: #333;
  line-height: 1;
}

.expire-card {
  flex: 1;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.expire-card :deep(.n-empty) {
  padding: 20px 0;
}

.card-header-with-action {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
  flex-shrink: 0;
}

.card-title {
  font-size: 14px;
  font-weight: 600;
  color: #333;
}

.check-btn {
  border-radius: 6px;
  font-size: 12px;
  padding: 0 12px;
  height: 28px;
  box-shadow: 0 2px 4px rgba(124, 152, 133, 0.2);
  transition: all 0.3s ease;
}

.check-btn:hover {
  box-shadow: 0 4px 8px rgba(124, 152, 133, 0.3);
  transform: translateY(-1px);
}

.check-btn:active {
  transform: translateY(0);
}

.expire-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  overflow-y: auto;
  flex: 1;
}

.expire-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 12px;
  background: rgba(0,0,0,0.02);
  border-radius: 8px;
  flex-shrink: 0;
}

.expire-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.expire-account {
  font-size: 13px;
  color: #333;
}

.expire-days {
  font-size: 12px;
  font-weight: 500;
  padding: 4px 10px;
  border-radius: 12px;
}

.expire-days.urgent {
  background: #fff1f0;
  color: #cf1322;
}

.expire-days.warning {
  background: #fffbe6;
  color: #d48806;
}

.expire-days.normal {
  background: #f6ffed;
  color: #389e0d;
}
</style>
