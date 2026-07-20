<script setup lang="ts">
import { ref, onMounted, h } from 'vue'
import { NButton, useDialog, useNotification } from 'naive-ui'
import { BrowserOpenURL } from '../../wailsjs/runtime/runtime'
import { GetCurrentVersion, CheckUpdate, GetSystemInfo, GetProxyInfo } from '../../wailsjs/go/main/App'
import { marked } from 'marked'

const dialog = useDialog()
const notify = useNotification()
const version = ref('')
const checking = ref(false)
const info = ref({ os: '', build: '', cpu: '', ram: '', ipv4: '', ipv6: '' })
const proxyInfo = ref('')
let checkedAuto = false

onMounted(async () => {
  version.value = await GetCurrentVersion()
  proxyInfo.value = await GetProxyInfo()
  const raw = await GetSystemInfo()
  try { info.value = JSON.parse(raw) } catch {}

  if (checkedAuto) return
  checkedAuto = true
  const skipped = localStorage.getItem('optiwin_skip_version')
  const result = await CheckUpdate()
  if (result && !result.startsWith('err:') && result !== 'same') {
    const r = JSON.parse(result)
    if (r.version === skipped) return
    showUpdateDialog(r)
  }
})

function showUpdateDialog(r: { version: string; body: string; url: string }) {
  const htmlBody = marked.parse(r.body || '暂无更新说明')
  const d = dialog.warning({
    title: `发现新版本 ${r.version}`,
    style: 'width:70vw;max-width:800px',
    content: () => h('div', { class: 'md-content', innerHTML: htmlBody }),
    action: () =>
      h('div', { style: 'display:flex;gap:8px;justify-content:flex-end;margin-top:12px' }, [
        h(NButton, { size: 'small', type: 'primary', onClick: () => {
          BrowserOpenURL(r.url); d.destroy()
        }}, { default: () => '前往下载' }),
        h(NButton, { size: 'small', onClick: () => {
          localStorage.setItem('optiwin_skip_version', r.version); d.destroy()
        }}, { default: () => '不再提示' }),
        h(NButton, { size: 'small', quaternary: true, onClick: () => d.destroy() },
          { default: () => '取消' }),
      ])
  })
}

function copyText(text: string, label: string) {
  navigator.clipboard.writeText(text).then(() => {
    notify.success({ title: '已复制', description: label, duration: 2000 })
  })
}

async function onCheckUpdate() {
  checking.value = true
  const result = await CheckUpdate()
  checking.value = false
  if (!result || result.startsWith('err:')) {
    const msg = result && result.startsWith('err:') ? result.slice(4) : '无法连接'
    notify.error({ title: '检查更新', description: msg, duration: 5000 })
    return
  }
  if (result === 'same') {
    notify.success({ title: '更新检查', description: '已是最新版本', duration: 4000 })
    return
  }
  const r = JSON.parse(result)
  const skipped = localStorage.getItem('optiwin_skip_version')
  if (r.version === skipped) return
  showUpdateDialog(r)
}
</script>

<template>
  <div class="page">
    <div class="hero">
      <div class="hero-icon"><img src="../assets/logo.png" /></div>
      <h1>OptiWin</h1>
      <p class="version">{{ version }}</p>
      <button class="update-btn" :disabled="checking" @click="onCheckUpdate">
        {{ checking ? '检查中...' : '检查更新' }}
      </button>
    </div>

    <div class="bottom-row">
      <div class="bottom-left">
        <div class="info-card"><a href="https://github.com/YouRan4/OptiWin" target="_blank" class="github-link">
          <svg width="20" height="20" viewBox="0 0 16 16" fill="currentColor">
            <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0016 8c0-4.42-3.58-8-8-8z"/>
          </svg><span>YouRan4/OptiWin</span>
        </a></div>
        <div class="info-card">
          <div class="info-header">基于 meetrevision/revision-tool 二次开发</div>
          <a href="https://github.com/meetrevision/revision-tool" target="_blank" class="info-link">meetrevision/revision-tool</a>
        </div>
        <div class="info-card">
          <div class="info-header">许可证</div>
          <p class="license-text">
            遵循 <a href="https://www.gnu.org/licenses/gpl-3.0.html" target="_blank">GNU General Public License v3.0</a>
          </p>
        </div>
      </div>

      <div class="bottom-right">
        <div class="setting-card">
          <div class="setting-card-header"><span class="header-title">系统信息</span></div>
          <div class="info-row"><span class="info-label">操作系统</span><span class="info-value">{{ info.os }}</span></div>
          <div class="info-row"><span class="info-label">系统版本</span><span class="info-value">Build {{ info.build }}</span></div>
          <div class="info-row"><span class="info-label">处理器</span><span class="info-value">{{ info.cpu }}</span></div>
          <div class="info-row"><span class="info-label">内存</span><span class="info-value">{{ info.ram }}</span></div>
          <div class="info-row">
            <span class="info-label">IPv4</span>
            <span class="info-value copy-ip" @click="copyText(info.ipv4, 'IPv4')">{{ info.ipv4 }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">IPv6</span>
            <span class="info-value copy-ip" @click="copyText(info.ipv6, 'IPv6')">{{ info.ipv6 }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">代理</span>
            <span class="info-value" style="font-size:12px">{{ proxyInfo }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.hero { display: flex; flex-direction: column; align-items: center; padding: 8px 0 12px; }
.hero-icon { margin: 0 auto 10px; display: flex; align-items: center; justify-content: center; }
.hero-icon img { width: 84px; height: 84px; }
h1 { font-size: 26px; font-weight: 700; margin-bottom: 4px; }
.version { font-size: 14px; color: var(--text2); margin-bottom: 10px; }
.update-btn {
  background: transparent; border: 1px solid var(--accent); color: var(--accent);
  border-radius: 6px; padding: 5px 18px; font-size: 13px; cursor: pointer; font-family: inherit;
}
.update-btn:hover { background: rgba(96,205,255,0.1); }
.update-btn:disabled { opacity: 0.5; cursor: default; }

.bottom-row { display: flex; gap: 24px; max-width: 900px; margin: 4px auto 0; }
.bottom-left { flex: 1; display: flex; flex-direction: column; gap: 10px; }
.bottom-right { flex: 1; }

.info-card {
  background: var(--card-bg); border: 1px solid var(--border);
  border-radius: 8px; padding: 14px;
}
.github-link {
  display: flex; align-items: center; gap: 10px;
  font-size: 15px; font-weight: 600; color: var(--text);
  text-decoration: none;
}
.github-link:hover { color: var(--accent); }
.info-header { font-size: 14px; font-weight: 600; margin-bottom: 6px; }
.info-link { display: block; font-size: 13px; color: var(--accent); text-decoration: none; }
.info-link:hover { text-decoration: underline; }
.license-text { font-size: 13px; color: var(--text2); line-height: 1.6; }
.license-text a { color: var(--accent); text-decoration: none; }
.license-text a:hover { text-decoration: underline; }

.info-row {
  display: flex; align-items: center; justify-content: space-between;
  padding: 8px 20px;
}
.info-row + .info-row { border-top: 1px solid var(--border); }
.info-label { font-size: 13px; color: var(--text2); white-space: nowrap; }
.info-value { font-size: 14px; text-align: right; word-break: break-all; max-width: 70%; }
.copy-ip { cursor: pointer; }
.copy-ip:hover { color: var(--accent); text-decoration: underline; }

.md-content { font-size: 13px; line-height: 1.7; color: var(--text); max-height: 60vh; overflow-y: auto; }
.md-content h1, .md-content h2, .md-content h3 { font-size: 15px; font-weight: 600; margin: 8px 0 4px; }
.md-content ul, .md-content ol { padding-left: 20px; margin: 4px 0; }
.md-content li { margin: 2px 0; }
.md-content code { background: rgba(255,255,255,0.1); padding: 1px 4px; border-radius: 3px; font-size: 12px; }
.md-content pre { background: rgba(0,0,0,0.2); padding: 8px; border-radius: 6px; overflow-x: auto; }
.md-content pre code { background: none; padding: 0; }
.md-content a { color: var(--accent); text-decoration: none; }
.md-content a:hover { text-decoration: underline; }
</style>
