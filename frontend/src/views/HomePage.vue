<script setup lang="ts">
import { ref, onMounted, h } from 'vue'
import { NButton, useDialog, useNotification } from 'naive-ui'
import { GetCurrentVersion, CheckUpdate } from '../../wailsjs/go/main/App'
import { marked } from 'marked'

const dialog = useDialog()
const notify = useNotification()
const version = ref('')
const checking = ref(false)
let checkedAuto = false

onMounted(async () => {
  version.value = await GetCurrentVersion()
  if (checkedAuto) return
  checkedAuto = true
  // 不提示的版本
  const skipped = localStorage.getItem('optiwin_skip_version')
  const result = await CheckUpdate()
  if (result) {
    const info = JSON.parse(result)
    if (info.version === skipped) return
    showUpdateDialog(info)
  }
})

function showUpdateDialog(info: { version: string; body: string; url: string }) {
  const htmlBody = marked.parse(info.body || '暂无更新说明')
  const d = dialog.warning({
    title: `发现新版本 ${info.version}`,
    style: 'width:70vw;max-width:800px',
    content: () => h('div', {
      class: 'md-content',
      innerHTML: htmlBody
    }),
    action: () =>
      h('div', { style: 'display:flex;gap:8px;justify-content:flex-end;margin-top:12px' }, [
        h(NButton, { size: 'small', type: 'primary', onClick: () => {
          (window as any).runtime?.BrowserOpenURL(info.url)
        }}, { default: () => '前往下载' }),
        h(NButton, { size: 'small', onClick: () => {
          localStorage.setItem('optiwin_skip_version', info.version); d.destroy()
        }}, { default: () => '不再提示' }),
        h(NButton, { size: 'small', quaternary: true, onClick: () => d.destroy() },
          { default: () => '取消' }),
      ])
  })
}

async function onCheckUpdate() {
  checking.value = true
  const result = await CheckUpdate()
  checking.value = false

  if (!result) {
    notify.success({ title: '更新检查', description: '已是最新版本', duration: 4000 })
    return
  }

  const info = JSON.parse(result)
  const skipped = localStorage.getItem('optiwin_skip_version')
  if (info.version === skipped) return
  showUpdateDialog(info)
}
</script>

<template>
  <div class="page">
    <div class="hero">
      <div class="hero-icon">
        <img src="../assets/logo.png" />
      </div>
      <h1>OptiWin</h1>
      <p class="version">{{ version }}</p>
      <p class="desc">适用于任何 Windows 系统的个性化调整工具箱</p>
      <button class="update-btn" :disabled="checking" @click="onCheckUpdate">
        {{ checking ? '检查中...' : '检查更新' }}
      </button>
    </div>

    <div class="info-section">
      <div class="info-card github-card">
        <a href="https://github.com/YouRan4/OptiWin" target="_blank" class="github-link">
          <svg width="20" height="20" viewBox="0 0 16 16" fill="currentColor">
            <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0016 8c0-4.42-3.58-8-8-8z"/>
          </svg>
          <span>YouRan4/OptiWin</span>
        </a>
      </div>
      <div class="info-card">
        <div class="info-header">基于 meetrevision/revision-tool 二次开发</div>
        <a href="https://github.com/meetrevision/revision-tool" target="_blank" class="info-link">meetrevision/revision-tool</a>
      </div>
      <div class="info-card">
        <div class="info-header">许可证</div>
        <p class="license-text">
          本项目遵循
          <a href="https://www.gnu.org/licenses/gpl-3.0.html" target="_blank">GNU General Public License v3.0</a>
          开源。
        </p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.hero { text-align: center; padding: 32px 0 24px; }
.hero-icon { margin: 0 auto 12px; display: flex; align-items: center; justify-content: center; }
.hero-icon img { width: 96px; height: 96px; }
h1 { font-size: 28px; font-weight: 700; margin-bottom: 4px; }
.version { font-size: 14px; color: var(--text2); margin-bottom: 8px; }
.desc { font-size: 15px; color: var(--text2); line-height: 1.6; max-width: 480px; margin: 0 auto 16px; }
.update-btn {
  background: transparent; border: 1px solid var(--accent); color: var(--accent);
  border-radius: 6px; padding: 6px 20px; font-size: 13px; cursor: pointer; font-family: inherit;
}
.update-btn:hover { background: rgba(96,205,255,0.1); }
.update-btn:disabled { opacity: 0.5; cursor: default; }

.info-section { max-width: 600px; margin: 0 auto; display: flex; flex-direction: column; gap: 12px; }
.info-card {
  background: var(--card-bg); border: 1px solid var(--border);
  border-radius: 8px; padding: 20px;
}
.github-card { padding: 16px !important; }
.github-link {
  display: flex; align-items: center; gap: 10px;
  font-size: 15px; font-weight: 600; color: var(--text);
  text-decoration: none;
}
.github-link:hover { color: var(--accent); }
.info-header { font-size: 15px; font-weight: 600; margin-bottom: 8px; }
.info-link { display: block; font-size: 14px; color: var(--accent); text-decoration: none; }
.info-link:hover { text-decoration: underline; }
.md-content { font-size: 13px; line-height: 1.7; color: var(--text); max-height: 60vh; overflow-y: auto; }
.md-content h1, .md-content h2, .md-content h3 { font-size: 15px; font-weight: 600; margin: 8px 0 4px; }
.md-content ul, .md-content ol { padding-left: 20px; margin: 4px 0; }
.md-content li { margin: 2px 0; }
.md-content code { background: rgba(255,255,255,0.1); padding: 1px 4px; border-radius: 3px; font-size: 12px; }
.md-content pre { background: rgba(0,0,0,0.2); padding: 8px; border-radius: 6px; overflow-x: auto; }
.md-content pre code { background: none; padding: 0; }
.md-content a { color: var(--accent); text-decoration: none; }
.md-content a:hover { text-decoration: underline; }

.license-text { font-size: 13px; color: var(--text2); line-height: 1.6; }
.license-text a { color: var(--accent); text-decoration: none; }
.license-text a:hover { text-decoration: underline; }
</style>
