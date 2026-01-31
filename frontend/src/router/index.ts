import { createRouter, createWebHashHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  { path: '/', redirect: '/dashboard' },
  { path: '/dashboard', name: 'Dashboard', component: () => import('../views/Dashboard.vue') },
  { path: '/accounts', name: 'AccountList', component: () => import('../views/AccountList.vue') },
  { path: '/email/settings', name: 'EmailSettings', component: () => import('../views/EmailSettings.vue') },
  { path: '/settings', name: 'Settings', component: () => import('../views/Settings.vue') },
  { path: '/about', name: 'About', component: () => import('../views/About.vue') }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router
