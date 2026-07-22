<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { NSwitch, NButton, NModal } from 'naive-ui'
import {
  GetHibernateStatus, EnableHibernate, DisableHibernate,
  GetFastStartupStatus, EnableFastStartup, DisableFastStartup,
  GetPhotoViewerStatus, EnablePhotoViewer, DisablePhotoViewer,
  UninstallEdge, GetWebView2Version, InstallWebView2,
  SetSafeBoot, RebootSystem, RebootToBios,
} from '../../wailsjs/go/main/App'

import { useNotify } from '../composables/useNotify'

const { t: i18n } = useI18n()
const notify = useNotify()

const hibernate = ref(false)
const fastStartup = ref(true)
const photoViewer = ref(false)
const webviewVer = ref('')
const showBootConfirm = ref(false)
const pendingBootMode = ref('')
const pendingBootTitle = ref('')
const pendingBootDesc = ref('')

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
  notify.create({ title: i18n('util.uninstallEdge'), description: msg, duration: 10000 })
}

async function onInstallWebView2() {
  const n = notify.create({ title: 'WebView2', description: i18n('util.downloading'), duration: 0 })
  const msg = await InstallWebView2()
  n.destroy()
  notify.create({ title: 'WebView2', description: msg, duration: 8000 })
  webviewVer.value = await GetWebView2Version()
}

function onBootClick(mode: string) {
  pendingBootMode.value = mode
  const map: Record<string, [string, string]> = {
    minimal: [i18n('util.bootSafe'), i18n('util.saveWork')],
    network: [i18n('util.bootSafeNet'), i18n('util.saveWork')],
    normal: [i18n('util.bootNormal'), i18n('util.saveWork')],
    bios: [i18n('util.bootBios'), i18n('util.biosDesc')],
  }
  ;[pendingBootTitle.value, pendingBootDesc.value] = map[mode]
  showBootConfirm.value = true
}

async function confirmBoot() {
  showBootConfirm.value = false
  if (pendingBootMode.value === 'bios') {
    await RebootToBios()
  } else {
    await SetSafeBoot(pendingBootMode.value)
    await RebootSystem()
  }
}
</script>

<template>
  <div class="page">
    <div class="setting-card">
      <div class="setting-card-header setting-card-header--flat"><span class="header-title">{{ i18n('util.powerManagement') }}</span></div>
      <div class="setting-row"><div><div class="row-label">{{ i18n('util.hibernate') }}</div><div class="row-desc">{{ i18n('util.hibernateDesc') }}</div></div><n-switch v-model:value="hibernate" @update:value="onHibernate" /></div>
      <div class="setting-row"><div><div class="row-label">{{ i18n('util.fastStartup') }}</div><div class="row-desc">{{ i18n('util.fastStartupDesc') }}</div></div><n-switch v-model:value="fastStartup" @update:value="onFastStartup" /></div>
    </div>
    <div class="setting-card">
      <div class="setting-card-header setting-card-header--flat"><span class="header-title">{{ i18n('util.apps') }}</span></div>
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
      <div class="setting-card-header setting-card-header--flat">
        <span class="header-title">{{ i18n('util.bootOptions') }}</span>
      </div>
      <div class="boot-buttons setting-row">
        <n-button size="small" strong style="flex:1" @click="onBootClick('minimal')">{{ i18n('util.bootSafe') }}</n-button>
        <n-button size="small" strong style="flex:1" @click="onBootClick('network')">{{ i18n('util.bootSafeNet') }}</n-button>
        <n-button size="small" strong style="flex:1" @click="onBootClick('normal')">{{ i18n('util.bootNormal') }}</n-button>
        <n-button size="small" strong style="flex:1" @click="onBootClick('bios')">{{ i18n('util.bootBios') }}</n-button>
      </div>
    </div>

    <n-modal
      v-model:show="showBootConfirm"
      preset="card"
      :title="pendingBootTitle"
      style="width: 480px; max-width: 90vw;"
      :mask-closable="false"
      closable
      @close="showBootConfirm = false"
    >
      <p>{{ pendingBootDesc }}</p>
      <p v-if="pendingBootMode !== 'normal' && pendingBootMode !== 'bios'" style="margin-top: 8px; color: var(--text2); font-size: 13px;">
        {{ i18n('util.safeModeExitNote') }}
      </p>
      <template #footer>
        <div style="display: flex; justify-content: flex-end; gap: 8px;">
          <n-button @click="showBootConfirm = false">{{ i18n('home.cancel') }}</n-button>
          <n-button @click="confirmBoot">{{ i18n('util.confirmReboot') }}</n-button>
        </div>
      </template>
    </n-modal>
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
