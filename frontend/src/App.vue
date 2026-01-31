<template>
  <n-config-provider :locale="naiveLocale" :date-locale="naiveDateLocale" :theme-overrides="themeOverrides">
    <n-message-provider :max="2">
      <n-dialog-provider>
        <div class="app">
          <!-- 顶部导航 -->
          <header class="header" style="--wails-draggable:drag">
            <div class="header-left">
              <span class="logo">{{ $t('settings.appNameValue') }}</span>
              <nav class="nav" style="--wails-draggable:no-drag">
                <router-link v-for="item in navItems" :key="item.path" :to="item.path" class="nav-item" :class="{ active: $route.path === item.path }">
                  <n-icon :component="item.icon" />
                  <span>{{ $t(item.label) }}</span>
                </router-link>
              </nav>
            </div>
            <div class="header-right" style="--wails-draggable:no-drag">
              <n-button quaternary circle size="small" @click="showAbout" class="about-btn">
                <template #icon>
                  <n-icon :component="SparklesOutline" />
                </template>
              </n-button>
              <n-dropdown :options="langOptions" @select="handleLangChange">
                <n-button quaternary circle size="small" class="lang-btn">
                  <template #icon>
                    <n-icon :component="LanguageOutline" />
                  </template>
                </n-button>
              </n-dropdown>
              <n-button quaternary circle size="small" @click="closeWindow" class="close-btn">
                <template #icon>
                  <n-icon :component="CloseOutline" />
                </template>
              </n-button>
            </div>
          </header>

          <!-- 主内容 -->
          <main class="main">
            <router-view v-slot="{ Component, route }">
              <transition name="fade" mode="out-in">
                <keep-alive :exclude="['About']">
                  <component :is="Component" :key="route.path" />
                </keep-alive>
              </transition>
            </router-view>
          </main>
        </div>
      </n-dialog-provider>
    </n-message-provider>
  </n-config-provider>
</template>

<script setup lang="ts">
import { ref, computed, markRaw } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { NConfigProvider, NMessageProvider, NDialogProvider, NIcon, NButton, NDropdown, zhCN, dateZhCN, enUS, dateEnUS } from 'naive-ui'
import { HomeOutline, ListOutline, MailOutline, SettingsOutline, CloseOutline, LanguageOutline, SparklesOutline } from '@vicons/ionicons5'
import { Quit } from '../wailsjs/runtime/runtime'
import { useMessage } from 'naive-ui'

const { locale } = useI18n()
const router = useRouter()

const navItems = [
  { path: '/dashboard', label: 'nav.dashboard', icon: markRaw(HomeOutline) },
  { path: '/accounts', label: 'nav.accounts', icon: markRaw(ListOutline) },
  { path: '/email/settings', label: 'nav.email', icon: markRaw(MailOutline) },
  { path: '/settings', label: 'nav.settings', icon: markRaw(SettingsOutline) }
]

const langOptions = [
  { label: '简体中文', key: 'zh-CN' },
  { label: 'English', key: 'en-US' }
]

const naiveLocale = computed(() => locale.value === 'zh-CN' ? zhCN : enUS)
const naiveDateLocale = computed(() => locale.value === 'zh-CN' ? dateZhCN : dateEnUS)

const themeOverrides = {
  common: {
    primaryColor: '#7C9885',
    primaryColorHover: '#8FAA9A',
    primaryColorPressed: '#6A8573',
    primaryColorSuppl: '#7C9885'
  },
  Button: {
    colorPrimary: '#7C9885',
    colorHoverPrimary: '#8FAA9A',
    colorPressedPrimary: '#6A8573',
    colorFocusPrimary: '#7C9885',
    borderPrimary: '1px solid #7C9885',
    borderHoverPrimary: '1px solid #8FAA9A',
    borderPressedPrimary: '1px solid #6A8573',
    borderFocusPrimary: '1px solid #7C9885'
  }
}

function handleLangChange(key: string) {
  locale.value = key
  localStorage.setItem('locale', key)
}

function showAbout() {
  router.push('/about')
}

function closeWindow() {
  Quit()
}
</script>

<style>
* { margin: 0; padding: 0; box-sizing: border-box; }

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  font-size: 13px;
  overflow: hidden;
}

/* 隐藏滚动条但保留滚动功能 */
::-webkit-scrollbar {
  display: none;
}

* {
  -ms-overflow-style: none;
  scrollbar-width: none;
}

.app {
  height: 100vh;
  background: #F8F6F4;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.header {
  height: 52px;
  background: #FFFFFF;
  border-bottom: 1px solid #E8E4E0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  flex-shrink: 0;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 32px;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 4px;
}

.about-btn:hover {
  background: rgba(124, 152, 133, 0.2) !important;
  color: #7C9885 !important;
}

.close-btn:hover {
  background: rgba(201, 165, 165, 0.2) !important;
  color: #C9A5A5 !important;
}

.lang-btn:hover {
  background: rgba(124, 152, 133, 0.2) !important;
  color: #7C9885 !important;
}

.logo {
  font-size: 16px;
  font-weight: 700;
  color: #7C9885;
}

.nav {
  display: flex;
  gap: 4px;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 14px;
  border-radius: 8px;
  color: #666;
  text-decoration: none;
  font-size: 13px;
}

.nav-item:hover {
  background: rgba(124,152,133,0.1);
  color: #7C9885;
}

.nav-item.active {
  background: #7C9885;
  color: #fff;
}

.main {
  flex: 1;
  padding: 20px 24px;
  max-width: 1400px;
  width: 100%;
  margin: 0 auto;
  overflow: hidden;
  height: calc(100vh - 52px);
}

.page-title {
  font-size: 20px;
  font-weight: 600;
  color: #333;
  margin-bottom: 16px;
}

/* 卡片样式 */
.card {
  background: #FFFFFF;
  border-radius: 12px;
  padding: 16px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.06);
  border: 1px solid #E8E4E0;
}

/* 路由过渡动画 */
.fade-enter-active {
  transition: opacity 0.12s ease;
}

.fade-leave-active {
  transition: opacity 0.08s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
