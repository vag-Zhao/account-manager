<template>
  <div class="about-page">
    <div class="about-content">
      <!-- 左侧 Logo 区域 -->
      <div class="logo-section">
        <div class="logo-wrapper">
          <div class="logo-circle">
            <n-icon :component="ShieldCheckmarkOutline" :size="56" class="logo-icon" />
          </div>
          <div class="logo-rings">
            <div class="ring"></div>
            <div class="ring"></div>
          </div>
        </div>
        <h1 class="app-name">{{ $t('settings.appNameValue') }}</h1>
        <p class="version">v1.0.0</p>
      </div>

      <!-- 右侧信息区域 -->
      <div class="info-section">
        <div class="info-card">
          <div class="card-icon">
            <n-icon :component="CodeSlashOutline" :size="22" />
          </div>
          <div class="card-content">
            <div class="card-title">{{ $t('about.techStack') }}</div>
            <div class="card-desc">Wails + Vue 3 + TypeScript</div>
          </div>
        </div>

        <div class="info-card">
          <div class="card-icon">
            <n-icon :component="RocketOutline" :size="22" />
          </div>
          <div class="card-content">
            <div class="card-title">{{ $t('about.features') }}</div>
            <div class="card-desc">{{ $t('about.featuresDesc') }}</div>
          </div>
        </div>

        <div class="info-card">
          <div class="card-icon">
            <n-icon :component="HeartOutline" :size="22" />
          </div>
          <div class="card-content">
            <div class="card-title">{{ $t('about.madeWith') }}</div>
            <div class="card-desc">{{ $t('about.madeWithDesc') }}</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { NIcon } from 'naive-ui'
import { ShieldCheckmarkOutline, CodeSlashOutline, RocketOutline, HeartOutline } from '@vicons/ionicons5'

// 设置组件名称用于 keep-alive 排除
defineOptions({
  name: 'About'
})
</script>

<style scoped>
.about-page {
  display: flex;
  align-items: center;
  justify-content: center;
  height: calc(100vh - 92px);
  padding: 40px;
  will-change: opacity;
}

.about-content {
  display: flex;
  flex-direction: row;
  gap: 80px;
  align-items: flex-start;
  max-width: 1000px;
}

.logo-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  animation: logoSlideIn 0.5s ease-out;
  flex-shrink: 0;
  width: 220px;
  margin-top: 66px;
}

@keyframes logoSlideIn {
  0% { opacity: 0; transform: translateX(-20px); }
  100% { opacity: 1; transform: translateX(0); }
}

.logo-wrapper {
  position: relative;
  margin-bottom: 28px;
}

.logo-circle {
  width: 120px;
  height: 120px;
  background: linear-gradient(135deg, #7C9885 0%, #9EB3C2 100%);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 12px 32px rgba(124, 152, 133, 0.25);
  position: relative;
  z-index: 2;
  animation: logoFloat 3s ease-in-out infinite;
  will-change: transform;
}

@keyframes logoFloat {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-8px); }
}

.logo-icon {
  color: #fff;
}

.logo-rings {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 120px;
  height: 120px;
  pointer-events: none;
}

.ring {
  position: absolute;
  inset: 0;
  border-radius: 50%;
  background: radial-gradient(circle, transparent 50%, rgba(124, 152, 133, 0.08) 65%, rgba(124, 152, 133, 0.06) 80%, transparent 100%);
  animation: ringPulse 3s ease-out infinite;
  will-change: transform, opacity;
  filter: blur(4px);
}

.ring:nth-child(2) {
  animation-delay: 1.5s;
}

@keyframes ringPulse {
  0% {
    transform: scale(1);
    opacity: 0.8;
  }
  100% {
    transform: scale(2.2);
    opacity: 0;
  }
}

.app-name {
  font-size: 20px;
  font-weight: 700;
  color: #333;
  margin: 0 0 10px 0;
  text-align: center;
}

.version {
  font-size: 12px;
  color: #999;
  font-weight: 500;
  padding: 4px 12px;
  background: rgba(124, 152, 133, 0.1);
  border-radius: 12px;
}

.info-section {
  display: flex;
  flex-direction: column;
  gap: 14px;
  flex: 1;
  width: 320px;
}

.info-card {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 16px 18px;
  background: #fff;
  border-radius: 10px;
  border: 1px solid #E8E4E0;
  box-shadow: 0 2px 6px rgba(0,0,0,0.04);
  transition: all 0.25s ease;
  animation: cardSlideIn 0.4s ease-out backwards;
  will-change: transform;
}

.info-card:nth-child(1) { animation-delay: 0.05s; }
.info-card:nth-child(2) { animation-delay: 0.1s; }
.info-card:nth-child(3) { animation-delay: 0.15s; }

@keyframes cardSlideIn {
  0% { opacity: 0; transform: translateX(15px); }
  100% { opacity: 1; transform: translateX(0); }
}

.info-card:hover {
  transform: translateX(4px);
  box-shadow: 0 4px 14px rgba(124, 152, 133, 0.12);
  border-color: rgba(124, 152, 133, 0.3);
}

.card-icon {
  width: 38px;
  height: 38px;
  border-radius: 9px;
  background: linear-gradient(135deg, #7C9885 0%, #9EB3C2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  flex-shrink: 0;
}

.card-content {
  flex: 1;
  min-width: 0;
}

.card-title {
  font-size: 13px;
  font-weight: 600;
  color: #333;
  margin-bottom: 4px;
}

.card-desc {
  font-size: 12px;
  color: #666;
  line-height: 1.5;
}

@media (max-width: 600px) {
  .about-content {
    flex-direction: column;
    gap: 40px;
    max-width: 500px;
  }

  .logo-section {
    width: auto;
  }

  .info-section {
    max-width: 100%;
  }
}
</style>
