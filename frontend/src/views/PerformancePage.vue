<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { NSwitch, NButton } from 'naive-ui'
import { Zap, Cpu, Database, Archive, Monitor, Layout, Layers, RefreshCw, Gamepad2 } from 'lucide-vue-next'
import {
  GetUltimatePerformanceStatus, EnableUltimatePerformance, DisableUltimatePerformance,
  GetCStateStatus, EnableCState, DisableCState,
  GetSuperfetchStatus, EnableSuperfetch, DisableSuperfetch,
  GetFullscreenOptimizationStatus, EnableFullscreenOptimization, DisableFullscreenOptimization,
  GetWindowedOptimizationStatus, EnableWindowedOptimization, DisableWindowedOptimization,
  GetMpoStatus, EnableMpo, DisableMpo,
  GetMemoryCompressionStatus, EnableMemoryCompression, DisableMemoryCompression,
  ClearShaderCache,
  RemoveGameBar, RestoreGameBar,
} from '../../wailsjs/go/main/App'

import { useNotify } from '../composables/useNotify'
import { useBackendToggle } from '../composables/useBackendToggle'
import WaitModal from '../components/WaitModal.vue'

const { t: i18n } = useI18n()
const notify = useNotify()
const showModal = ref(false)
const modalText = ref('')

const { value: reviPlan, refresh: refreshReviPlan, onUpdate: onPowerPlan, busy: reviPlanBusy } = useBackendToggle(
  GetUltimatePerformanceStatus,
  v => (v ? EnableUltimatePerformance() : DisableUltimatePerformance()),
  () => i18n('perf.powerPlan'),
)
const { value: cstate, refresh: refreshCState, onUpdate: onCState, busy: cstateBusy } = useBackendToggle(
  GetCStateStatus,
  v => (v ? EnableCState() : DisableCState()),
  () => i18n('perf.cState'),
)
const { value: superfetch, refresh: refreshSuperfetch, onUpdate: onSuperfetch, busy: superfetchBusy } = useBackendToggle(
  GetSuperfetchStatus,
  v => (v ? EnableSuperfetch() : DisableSuperfetch()),
  () => i18n('perf.superfetch'),
)
const { value: memCompress, refresh: refreshMemCompress, onUpdate: onMemCompress, busy: memCompressBusy } = useBackendToggle(
  GetMemoryCompressionStatus,
  v => (v ? EnableMemoryCompression() : DisableMemoryCompression()),
  () => i18n('perf.memCompress'),
  () => notify.create({ title: i18n('perf.memCompress'), description: i18n('perf.restartRequired'), duration: 6000 }),
)
const { value: fullscreen, refresh: refreshFullscreen, onUpdate: onFullscreen, busy: fullscreenBusy } = useBackendToggle(
  GetFullscreenOptimizationStatus,
  v => (v ? EnableFullscreenOptimization() : DisableFullscreenOptimization()),
  () => i18n('perf.fullscreen'),
)
const { value: windowOpt, refresh: refreshWindowOpt, onUpdate: onWindowed, busy: windowOptBusy } = useBackendToggle(
  GetWindowedOptimizationStatus,
  v => (v ? EnableWindowedOptimization() : DisableWindowedOptimization()),
  () => i18n('perf.windowed'),
)
const { value: mpo, refresh: refreshMpo, onUpdate: onMpo, busy: mpoBusy } = useBackendToggle(
  GetMpoStatus,
  v => (v ? EnableMpo() : DisableMpo()),
  () => i18n('perf.mpo'),
)

onMounted(async () => {
  await Promise.all([
    refreshReviPlan(),
    refreshCState(),
    refreshSuperfetch(),
    refreshMemCompress(),
    refreshFullscreen(),
    refreshWindowOpt(),
    refreshMpo(),
  ])
})

function showUpdateFailed(title: string) {
  notify.create({
    title,
    description: i18n('common.updateFailedDesc'),
    type: 'error',
    duration: 6000,
  })
}

async function onClearShaderCache() {
  const msg = await ClearShaderCache()
  notify.create({ title: i18n('perf.shaderCache'), description: msg, duration: 5000 })
}

async function onRemoveGameBar() {
  modalText.value = i18n('perf.removingGameBar')
  showModal.value = true
  try {
    const ok = await RemoveGameBar()
    if (ok) {
      notify.create({ title: i18n('perf.gameBar'), description: i18n('perf.gameBarRemoved'), duration: 5000 })
    } else {
      showUpdateFailed(i18n('perf.gameBar'))
    }
  } finally {
    showModal.value = false
  }
}

async function onRestoreGameBar() {
  modalText.value = i18n('perf.restoringGameBar')
  showModal.value = true
  try {
    const ok = await RestoreGameBar()
    if (ok) {
      notify.create({ title: i18n('perf.gameBar'), description: i18n('perf.gameBarRestored'), duration: 5000 })
    } else {
      showUpdateFailed(i18n('perf.gameBar'))
    }
  } finally {
    showModal.value = false
  }
}
</script>

<template>
  <div class="page">
    <WaitModal :show="showModal" :text="modalText" />
    <div class="setting-card">
      <div class="setting-card-header setting-card-header--flat"><span class="header-title">{{ i18n('perf.powerPlan') }}</span></div>
      <div class="setting-row">
        <Zap :size="18" class="row-icon" />
        <div>
          <div class="row-label">{{ i18n('perf.ultimatePerf') }}</div>
          <div class="row-desc">{{ i18n('perf.ultimatePerfDesc') }}</div>
        </div><n-switch v-model:value="reviPlan" :loading="reviPlanBusy" @update:value="onPowerPlan" />
      </div>
      <div class="setting-row">
        <Cpu :size="18" class="row-icon" />
        <div>
          <div class="row-label">{{ i18n('perf.cState') }}</div>
          <div class="row-desc">{{ i18n('perf.cStateDesc') }}</div>
        </div><n-switch v-model:value="cstate" :loading="cstateBusy" @update:value="onCState" />
      </div>
    </div>
    <div class="setting-card">
      <div class="setting-card-header setting-card-header--flat"><span class="header-title">{{ i18n('perf.memory') }}</span></div>
      <div class="setting-row">
        <Database :size="18" class="row-icon" />
        <div>
          <div class="row-label">{{ i18n('perf.superfetch') }}</div>
          <div class="row-desc">{{ i18n('perf.superfetchDesc') }}</div>
        </div><n-switch v-model:value="superfetch" :loading="superfetchBusy" @update:value="onSuperfetch" />
      </div>
      <div class="setting-row">
        <Archive :size="18" class="row-icon" />
        <div>
          <div class="row-label">{{ i18n('perf.memCompress') }}</div>
          <div class="row-desc">{{ i18n('perf.memCompressDesc') }}</div>
        </div><n-switch v-model:value="memCompress" :loading="memCompressBusy" @update:value="onMemCompress" />
      </div>
    </div>
    <div class="setting-card">
      <div class="setting-card-header setting-card-header--flat"><span class="header-title">{{ i18n('perf.display') }}</span></div>
      <div class="setting-row">
        <Monitor :size="18" class="row-icon" />
        <div>
          <div class="row-label">{{ i18n('perf.fullscreen') }}</div>
          <div class="row-desc">{{ i18n('perf.fullscreenDesc') }}</div>
        </div><n-switch v-model:value="fullscreen" :loading="fullscreenBusy" @update:value="onFullscreen" />
      </div>
      <div class="setting-row">
        <Layout :size="18" class="row-icon" />
        <div>
          <div class="row-label">{{ i18n('perf.windowed') }}</div>
          <div class="row-desc">{{ i18n('perf.windowedDesc') }}</div>
        </div><n-switch v-model:value="windowOpt" :loading="windowOptBusy" @update:value="onWindowed" />
      </div>
      <div class="setting-row">
        <Layers :size="18" class="row-icon" />
        <div>
          <div class="row-label">{{ i18n('perf.mpo') }}</div>
          <div class="row-desc">{{ i18n('perf.mpoDesc') }}</div>
        </div><n-switch v-model:value="mpo" :loading="mpoBusy" @update:value="onMpo" />
      </div>
      <div class="setting-row">
        <RefreshCw :size="18" class="row-icon" />
        <div>
          <div class="row-label">{{ i18n('perf.clearShader') }}</div>
          <div class="row-desc">{{ i18n('perf.clearShaderDesc') }}</div>
        </div>
        <n-button size="small" @click="onClearShaderCache">{{ i18n('perf.clear') }}</n-button>
      </div>
    </div>

    <div class="setting-card">
      <div class="setting-card-header setting-card-header--flat"><span class="header-title">{{ i18n('perf.xboxServices') }}</span></div>
      <div class="setting-row">
        <Gamepad2 :size="18" class="row-icon" />
        <div>
          <div class="row-label">{{ i18n('perf.gameBar') }}</div>
          <div class="row-desc">{{ i18n('perf.gameBarDesc') }}</div>
        </div>
        <div style="display: flex; gap: 8px;">
          <n-button size="small" @click="onRemoveGameBar">{{ i18n('perf.uninstall') }}</n-button>
          <n-button size="small" @click="onRestoreGameBar">{{ i18n('perf.install') }}</n-button>
        </div>
      </div>
    </div>
  </div>
</template>
