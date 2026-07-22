<script setup lang="ts">
import { useRouter, useRoute } from 'vue-router'
import { ref, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { NConfigProvider, darkTheme, zhCN, dateZhCN, enUS, dateEnUS, NNotificationProvider, NDialogProvider } from 'naive-ui'

const router = useRouter()
const route = useRoute()
const { t: i18n, locale } = useI18n()

const isDark = ref(true)

onMounted(() => {
  const mq = window.matchMedia('(prefers-color-scheme: dark)')
  isDark.value = mq.matches
  applyTheme(isDark.value)
  mq.addEventListener('change', e => {
    isDark.value = e.matches
    applyTheme(e.matches)
  })
})

function applyTheme(dark: boolean) {
  document.documentElement.setAttribute('data-theme', dark ? 'dark' : 'light')
}

const theme = computed(() => isDark.value ? darkTheme : undefined)
const navItems = computed(() => [
  { path: '/', name: 'home', label: i18n('nav.home') },
  { path: '/security', name: 'security', label: i18n('nav.security') },
  { path: '/performance', name: 'performance', label: i18n('nav.performance') },
  { path: '/personalization', name: 'personalization', label: i18n('nav.personalization') },
  { path: '/utilities', name: 'utilities', label: i18n('nav.utilities') },
  { path: '/updates', name: 'updates', label: i18n('nav.updates') },
])

const naiveLocale = computed(() => locale.value === 'zh' ? zhCN : enUS)
const naiveDateLocale = computed(() => locale.value === 'zh' ? dateZhCN : dateEnUS)

const activeName = computed(() => route.name as string || 'home')

function navigate(path: string) {
  if (route.path === path) return
  router.push(path).catch(() => {})
}

function wm() { (window as any).runtime?.WindowMinimise() }
function wq() { (window as any).runtime?.Quit() }

function reExplorer() {
  import('../wailsjs/go/main/App').then(m => m.RestartExplorer()).catch(() => {})
}
</script>

<template>
  <NConfigProvider :theme="theme" :locale="naiveLocale" :date-locale="naiveDateLocale">
    <n-dialog-provider>
    <n-notification-provider>
    <div class="shell">
      <header class="titlebar">
        <div class="t-left" style="--wails-draggable: drag">
          <span class="t-title">OptiWin</span>
        </div>
        <nav class="t-nav">
          <button v-for="item in navItems" :key="item.name"
            class="nav-btn" :class="{ active: activeName === item.name }"
            @click="navigate(item.path)">{{ item.label }}</button>
        </nav>
        <div class="t-right" style="--wails-draggable: no-drag">
          <button class="tb-btn" @click="reExplorer" :title="i18n('nav.restartExplorer')">&#x21BB;</button>
          <button class="tb-btn" @click="wm">&#x2014;</button>
          <button class="tb-btn tb-close" style="--wails-draggable:no-drag" @click="wq">&#x2715;</button>
        </div>
      </header>
      <div class="body">
        <main class="content">
          <router-view v-slot="{ Component }">
            <transition name="page" mode="out-in">
              <component :is="Component" />
            </transition>
          </router-view>
        </main>
      </div>
    </div>
    </n-notification-provider>
    </n-dialog-provider>
  </NConfigProvider>
</template>

<style>
* { margin: 0; padding: 0; box-sizing: border-box; }

/* 深色主题 */
:root[data-theme="dark"] {
  --accent: #60CDFF;
  --spinner-color: #FFFFFF;
  --border: rgba(255,255,255,0.06);
  --text: #E5E5E5;
  --text2: rgba(255,255,255,0.4);
  --hover: rgba(255,255,255,0.05);
  --active: rgba(96,205,255,0.12);
  --card-bg: rgba(45,45,45,0.7);
  --section-bg: rgba(55,55,55,0.5);
}

/* 浅色主题 */
:root[data-theme="light"] {
  --accent: #1890FF;
  --spinner-color: #000000;
  --border: rgba(0,0,0,0.1);
  --text: #1A1A1A;
  --text2: rgba(0,0,0,0.45);
  --hover: rgba(0,0,0,0.04);
  --active: rgba(24,144,255,0.1);
  --card-bg: rgba(255,255,255,0.85);
  --section-bg: rgba(245,245,245,0.85);
}

body, #app {
  font-family: 'Segoe UI', system-ui, sans-serif;
  height: 100vh; overflow: hidden;
  color: var(--text); background: transparent;
}

.shell { height: 100vh; display: flex; flex-direction: column; background: transparent; }

/* 标题栏 */
.titlebar {
  display: flex; align-items: center;
  padding: 0 8px 0 16px; height: 48px; flex-shrink: 0; user-select: none;
  background: transparent; border-bottom: none;
  box-shadow: 0 1px 8px rgba(0,0,0,0.15);
  gap: 24px; z-index: 10;
}
.t-left { display: flex; align-items: center; }
.t-title { font-size: 15px; font-weight: 600; }

/* 导航按钮 */
.t-nav { display: flex; align-items: center; gap: 4px; flex: 1; overflow-x: auto; }
.nav-btn {
  background: none; border: none; cursor: pointer; color: var(--text);
  font-size: 14px; padding: 0 16px; height: 48px; white-space: nowrap;
  display: flex; align-items: center; position: relative; opacity: 0.6;
}
.nav-btn:hover { opacity: 1; background: var(--hover); }
.nav-btn.active { opacity: 1; font-weight: 600; }
.nav-btn.active::after {
  content: ''; position: absolute; bottom: 0; left: 8px; right: 8px;
  height: 2px; background: var(--accent); border-radius: 1px;
}

/* 窗口控制按钮 */
.t-right { display: flex; }
.tb-btn {
  background: none; border: none; cursor: pointer;
  font-size: 12px; padding: 0 16px; height: 48px;
  color: var(--text); display: flex; align-items: center;
}
.tb-btn:hover { background: var(--hover); }
.tb-close:hover { background: #e81123; color: #fff; }

/* 主内容区 */
.body { display: flex; flex: 1; overflow: hidden; background: transparent; }
.content { flex: 1; overflow-y: auto; padding: 32px 40px; }

h2 { font-size: 24px; font-weight: 600; margin-bottom: 4px; }
.setting-card {
  background: var(--card-bg); border: 1px solid var(--border);
  border-radius: 8px; margin-bottom: 12px; overflow: hidden;
}
.setting-card-header {
  padding: 16px 20px 12px;
  background: var(--section-bg);
  border-bottom: 1px solid var(--border);
  display: flex; align-items: baseline; gap: 8px;
}
.setting-card-header::before {
  content: ''; width: 3px; height: 18px;
  background: var(--accent); border-radius: 2px; flex-shrink: 0;
  align-self: center;
}
.setting-card-header .header-title { font-size: 16px; font-weight: 600; }
.setting-card-header .header-desc { font-size: 12px; color: var(--text2); }

.setting-row {
  display: flex; align-items: center; justify-content: space-between;
  padding: 12px 20px 12px 28px;
}
.setting-row + .setting-row { border-top: 1px solid var(--border); }
.setting-row .row-label { font-size: 14px; }
.setting-row .row-desc { font-size: 12px; color: var(--text2); margin-top: 1px; }

.page-enter-active, .page-leave-active { transition: opacity 0.3s cubic-bezier(0.4, 0, 0.2, 1), transform 0.3s cubic-bezier(0.4, 0, 0.2, 1); }
.page-enter-from { opacity: 0; transform: translateY(12px); }
.page-leave-to { opacity: 0; transform: translateY(-12px); }

::-webkit-scrollbar { width: 8px; height: 8px; }
::-webkit-scrollbar-track { background: transparent; }
[data-theme="dark"] ::-webkit-scrollbar-thumb { background: rgba(255,255,255,0.12); border-radius: 4px; }
[data-theme="dark"] ::-webkit-scrollbar-thumb:hover { background: rgba(255,255,255,0.2); }
[data-theme="light"] ::-webkit-scrollbar-thumb { background: rgba(0,0,0,0.12); border-radius: 4px; }
[data-theme="light"] ::-webkit-scrollbar-thumb:hover { background: rgba(0,0,0,0.2); }
</style>
