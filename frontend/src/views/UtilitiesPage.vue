<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { NSwitch, NButton, useNotification } from 'naive-ui'
import {
  GetHibernateStatus, EnableHibernate, DisableHibernate,
  GetFastStartupStatus, EnableFastStartup, DisableFastStartup,
  GetPhotoViewerStatus, EnablePhotoViewer, DisablePhotoViewer,
  UninstallEdge, GetWebView2Version, InstallWebView2,
  SetSafeBoot, RebootSystem, RebootToBios,
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

async function onSafeBoot(mode: string) {
  const tips: Record<string,string> = {
    minimal: '即将以安全模式重启',
    network: '即将以带网络的安全模式重启',
    normal: '即将以正常模式重启',
  }
  await SetSafeBoot(mode)
  notify.warning({ title: tips[mode] || '即将重启', description: '请保存好你的工作', duration: 5000 })
  await RebootSystem()
}

async function onBios() {
  notify.warning({ title: '即将进入 BIOS', description: '系统将在重启后进入 UEFI 固件设置，下次开机恢复正常', duration: 5000 })
  await RebootToBios()
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
    <div class="setting-card">
      <div class="setting-card-header">
        <span class="header-title">启动选项</span>
        <span class="header-desc">安全模式与系统重启</span>
        <span class="warning-text">— 请确定你知道你在干什么</span>
      </div>
      <div class="boot-buttons">
        <n-button size="small" strong style="flex:1" @click="onSafeBoot('minimal')">重启至安全模式</n-button>
        <n-button size="small" strong style="flex:1" @click="onSafeBoot('network')">重启至带网络连接的安全模式</n-button>
        <n-button size="small" strong style="flex:1" @click="onSafeBoot('normal')">重启至正常模式</n-button>
        <n-button size="small" strong style="flex:1" type="warning" @click="onBios">进入 BIOS（重启）</n-button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.warning-text { font-size: 12px; color: #ffb347; margin-left: auto; }
.boot-buttons {
  padding: 12px 20px 16px;
  display: flex;
  gap: 8px;
}
</style>
