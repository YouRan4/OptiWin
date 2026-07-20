<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { NSwitch, NButton, useNotification } from 'naive-ui'
import {
  GetHibernateStatus, EnableHibernate, DisableHibernate,
  GetFastStartupStatus, EnableFastStartup, DisableFastStartup,
  GetPhotoViewerStatus, EnablePhotoViewer, DisablePhotoViewer,
  UninstallEdge, GetWebView2Version, InstallWebView2,
} from '../../wailsjs/go/main/App'

const notify = useNotification()

const hibernate = ref(false)
const fastStartup = ref(true)
const photoViewer = ref(false)
const webviewVer = ref('加载中...')

onMounted(async () => {
  hibernate.value = await GetHibernateStatus()
  fastStartup.value = await GetFastStartupStatus()
  photoViewer.value = await GetPhotoViewerStatus()
  webviewVer.value = await GetWebView2Version()
})

async function onHibernate(v: boolean) {
  if (v) await EnableHibernate(); else await DisableHibernate()
  hibernate.value = await GetHibernateStatus()
}

async function onFastStartup(v: boolean) {
  if (v) await EnableFastStartup(); else await DisableFastStartup()
  fastStartup.value = await GetFastStartupStatus()
}

async function onPhotoViewer(v: boolean) {
  if (v) await EnablePhotoViewer(); else await DisablePhotoViewer()
  photoViewer.value = await GetPhotoViewerStatus()
}

async function onUninstallEdge() {
  const msg = await UninstallEdge()
  notify.info({ title: '卸载 Edge', description: msg, duration: 10000 })
}

async function onInstallWebView2() {
  const n = notify.info({ title: 'WebView2', description: '正在下载安装...', duration: 0 })
  const msg = await InstallWebView2()
  n.destroy()
  notify.success({ title: 'WebView2', description: msg, duration: 8000 })
  webviewVer.value = await GetWebView2Version()
}
</script>

<template>
  <div class="page">
    <h2>实用工具</h2>
    <div class="setting-card">
      <div class="setting-card-header"><span class="header-title">电源管理</span><span class="header-desc">休眠和启动设置</span></div>
      <div class="setting-row"><div><div class="row-label">休眠</div><div class="row-desc">保存系统状态到磁盘并完全关闭电源</div></div><n-switch v-model:value="hibernate" @update:value="onHibernate" /></div>
      <div class="setting-row"><div><div class="row-label">快速启动</div><div class="row-desc">结合休眠实现更快的启动速度</div></div><n-switch v-model:value="fastStartup" @update:value="onFastStartup" /></div>
    </div>
    <div class="setting-card">
      <div class="setting-card-header"><span class="header-title">应用设置</span><span class="header-desc">默认应用和浏览器配置</span></div>
      <div class="setting-row"><div><div class="row-label">Windows 照片查看器</div><div class="row-desc">启用经典的 Windows 照片查看器来打开图片文件</div></div><n-switch v-model:value="photoViewer" @update:value="onPhotoViewer" /></div>
      <div class="setting-row">
        <div><div class="row-label">卸载 Microsoft Edge</div><div class="row-desc">安全卸载（保留 WebView2 不受影响）</div></div>
        <n-button size="small" @click="onUninstallEdge">卸载</n-button>
      </div>
      <div class="setting-row">
        <div><div class="row-label">WebView2</div><div class="row-desc">当前版本: {{ webviewVer }}</div></div>
        <n-button size="small" @click="onInstallWebView2">安装/升级</n-button>
      </div>
    </div>
  </div>
</template>
