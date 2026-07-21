<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { NSwitch, NButton, useNotification } from 'naive-ui'
import {
  GetHibernateStatus, EnableHibernate, DisableHibernate,
  GetFastStartupStatus, EnableFastStartup, DisableFastStartup,
  GetPhotoViewerStatus, EnablePhotoViewer, DisablePhotoViewer,
  UninstallEdge, GetWebView2Version, InstallWebView2,
  SetSafeBoot, RebootSystem, RebootToBios,
} from '../../wailsjs/go/main/App'

const { t: i18n } = useI18n()
const notify = useNotification()

const hibernate = ref(false)
const fastStartup = ref(true)
const photoViewer = ref(false)
const webviewVer = ref('')
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
  notify.info({ title: i18n('util.uninstallEdge'), description: msg, duration: 10000 })
}

async function onInstallWebView2() {
  const n = notify.info({ title: 'WebView2', description: i18n('util.downloading'), duration: 0 })
  const msg = await InstallWebView2()
  n.destroy()
  notify.success({ title: 'WebView2', description: msg, duration: 8000 })
  webviewVer.value = await GetWebView2Version()
}

async function onSafeBoot(mode: string) {
  const tips: Record<string,string> = {
    minimal: i18n('util.rebootSafeMode'),
    network: i18n('util.rebootSafeModeNet'),
    normal: i18n('util.rebootNormal'),
  }
  await SetSafeBoot(mode)
  notify.warning({ title: tips[mode] || i18n('util.rebooting'), description: i18n('util.saveWork'), duration: 5000 })
  await RebootSystem()
}

async function onBios() {
  notify.warning({ title: i18n('util.enterBios'), description: i18n('util.biosDesc'), duration: 5000 })
  await RebootToBios()
}
</script>

<template>
  <div class="page">
    <h2>{{ i18n('util.title') }}</h2>
    <div class="setting-card">
      <div class="setting-card-header"><span class="header-title">{{ i18n('util.powerManagement') }}</span><span class="header-desc">{{ i18n('util.powerManagementDesc') }}</span></div>
      <div class="setting-row"><div><div class="row-label">{{ i18n('util.hibernate') }}</div><div class="row-desc">{{ i18n('util.hibernateDesc') }}</div></div><n-switch v-model:value="hibernate" @update:value="onHibernate" /></div>
      <div class="setting-row"><div><div class="row-label">{{ i18n('util.fastStartup') }}</div><div class="row-desc">{{ i18n('util.fastStartupDesc') }}</div></div><n-switch v-model:value="fastStartup" @update:value="onFastStartup" /></div>
    </div>
    <div class="setting-card">
      <div class="setting-card-header"><span class="header-title">{{ i18n('util.apps') }}</span><span class="header-desc">{{ i18n('util.appsDesc') }}</span></div>
      <div class="setting-row"><div><div class="row-label">{{ i18n('util.photoViewer') }}</div><div class="row-desc">{{ i18n('util.photoViewerDesc') }}</div></div><n-switch v-model:value="photoViewer" @update:value="onPhotoViewer" /></div>
      <div class="setting-row">
        <div><div class="row-label">{{ i18n('util.edgeLabel') }}</div><div class="row-desc">{{ i18n('util.edgeDesc') }}</div></div>
        <n-button size="small" @click="onUninstallEdge">{{ i18n('util.uninstall') }}</n-button>
      </div>
      <div class="setting-row">
        <div><div class="row-label">WebView2</div><div class="row-desc">{{ i18n('util.currentVersion') }} {{ webviewVer }}</div></div>
        <n-button size="small" @click="onInstallWebView2">{{ i18n('util.install') }}</n-button>
      </div>
    </div>
    <div class="setting-card">
      <div class="setting-card-header">
        <span class="header-title">{{ i18n('util.bootOptions') }}</span>
        <span class="header-desc">{{ i18n('util.bootDesc') }}</span>
        <span class="warning-text">{{ i18n('util.bootWarning') }}</span>
      </div>
      <div class="boot-buttons">
        <n-button size="small" strong style="flex:1" @click="onSafeBoot('minimal')">{{ i18n('util.bootSafe') }}</n-button>
        <n-button size="small" strong style="flex:1" @click="onSafeBoot('network')">{{ i18n('util.bootSafeNet') }}</n-button>
        <n-button size="small" strong style="flex:1" @click="onSafeBoot('normal')">{{ i18n('util.bootNormal') }}</n-button>
        <n-button size="small" strong style="flex:1" type="warning" @click="onBios">{{ i18n('util.bootBios') }}</n-button>
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
