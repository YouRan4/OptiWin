<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { NSwitch, NButton, NModal } from 'naive-ui'
import { Moon, Zap, Image, Globe, Download } from 'lucide-vue-next'
import {
  GetHibernateStatus, EnableHibernate, DisableHibernate,
  GetFastStartupStatus, EnableFastStartup, DisableFastStartup,
  GetPhotoViewerStatus, EnablePhotoViewer, DisablePhotoViewer,
  UninstallEdge, InstallEdge, GetWebView2Version, InstallWebView2,
  SetSafeBoot, RebootSystem, RebootToBios,
} from '../../wailsjs/go/main/App'

import { useNotify } from '../composables/useNotify'
import { useBackendToggle } from '../composables/useBackendToggle'

const { t: i18n } = useI18n()
const notify = useNotify()

const { value: hibernate, refresh: refreshHibernate, onUpdate: onHibernate, busy: hibernateBusy } = useBackendToggle(
  GetHibernateStatus, v => (v ? EnableHibernate() : DisableHibernate()), () => i18n('util.hibernate'),
)
const { value: fastStartup, refresh: refreshFastStartup, onUpdate: onFastStartup, busy: fastStartupBusy } = useBackendToggle(
  GetFastStartupStatus, v => (v ? EnableFastStartup() : DisableFastStartup()), () => i18n('util.fastStartup'),
)
const { value: photoViewer, refresh: refreshPhotoViewer, onUpdate: onPhotoViewer, busy: photoViewerBusy } = useBackendToggle(
  GetPhotoViewerStatus, v => (v ? EnablePhotoViewer() : DisablePhotoViewer()), () => i18n('util.photoViewer'),
)

const webviewVer = ref('')
const showBootConfirm = ref(false)
const pendingBootMode = ref('')
const pendingBootTitle = ref('')
const pendingBootDesc = ref('')

onMounted(async () => {
  await Promise.all([refreshHibernate(), refreshFastStartup(), refreshPhotoViewer()])
  webviewVer.value = await GetWebView2Version()
})

async function onUninstallEdge() {
  const msg = await UninstallEdge()
  notify.create({ title: i18n('util.uninstallEdge'), description: msg, duration: 10000 })
}

async function onInstallEdge() {
  const n = notify.create({ title: i18n('util.installEdge'), description: i18n('util.downloading'), duration: 0 })
  const msg = await InstallEdge()
  n.destroy()
  notify.create({ title: i18n('util.installEdge'), description: msg, duration: 10000 })
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
    if (!(await RebootToBios())) {
      notify.create({ title: i18n('common.updateFailed'), description: i18n('util.bootBios'), type: 'error', duration: 6000 })
    }
    return
  }
  const ok = await SetSafeBoot(pendingBootMode.value)
  if (!ok) {
    notify.create({
      title: i18n('common.updateFailed'),
      description: `${i18n('util.bootOptions')}：${i18n('common.updateFailedDesc')}`,
      type: 'error',
      duration: 8000,
    })
    return
  }
  await RebootSystem()
}
</script>

<template>
  <div class="page">
    <div class="setting-card">
      <div class="setting-card-header setting-card-header--flat"><span class="header-title">{{ i18n('util.powerManagement') }}</span></div>
      <div class="setting-row">
        <Moon :size="18" class="row-icon" />
        <div><div class="row-label">{{ i18n('util.hibernate') }}</div><div class="row-desc">{{ i18n('util.hibernateDesc') }}</div></div><n-switch v-model:value="hibernate" :loading="hibernateBusy" @update:value="onHibernate" />
      </div>
      <div class="setting-row">
        <Zap :size="18" class="row-icon" />
        <div><div class="row-label">{{ i18n('util.fastStartup') }}</div><div class="row-desc">{{ i18n('util.fastStartupDesc') }}</div></div><n-switch v-model:value="fastStartup" :loading="fastStartupBusy" @update:value="onFastStartup" />
      </div>
    </div>
    <div class="setting-card">
      <div class="setting-card-header setting-card-header--flat"><span class="header-title">{{ i18n('util.apps') }}</span></div>
      <div class="setting-row">
        <Image :size="18" class="row-icon" />
        <div><div class="row-label">{{ i18n('util.photoViewer') }}</div><div class="row-desc">{{ i18n('util.photoViewerDesc') }}</div></div><n-switch v-model:value="photoViewer" :loading="photoViewerBusy" @update:value="onPhotoViewer" />
      </div>
      <div class="setting-row">
        <Globe :size="18" class="row-icon" />
        <div><div class="row-label">{{ i18n('util.edgeLabel') }}</div><div class="row-desc">{{ i18n('util.edgeDesc') }}</div></div>
        <div style="display:flex; gap:8px; flex-shrink:0">
          <n-button size="small" @click="onUninstallEdge">{{ i18n('util.uninstall') }}</n-button>
          <n-button size="small" @click="onInstallEdge">{{ i18n('util.installEdge') }}</n-button>
        </div>
      </div>
      <div class="setting-row">
        <Download :size="18" class="row-icon" />
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
