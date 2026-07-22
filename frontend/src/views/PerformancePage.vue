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
  GetGameBarStatus, EnableGameBar, DisableGameBar,
} from '../../wailsjs/go/main/App'

import { useNotify } from '../composables/useNotify'

const { t: i18n } = useI18n()
const notify = useNotify()

const reviPlan = ref(false)
const cstate = ref(true)
const superfetch = ref(true)
const memCompress = ref(true)
const fullscreen = ref(true)
const windowOpt = ref(true)
const mpo = ref(true)
const gameBar = ref(true)

onMounted(async () => {
  reviPlan.value = await GetUltimatePerformanceStatus()
  cstate.value = await GetCStateStatus()
  superfetch.value = await GetSuperfetchStatus()
  memCompress.value = await GetMemoryCompressionStatus()
  fullscreen.value = await GetFullscreenOptimizationStatus()
  windowOpt.value = await GetWindowedOptimizationStatus()
  mpo.value = await GetMpoStatus()
  gameBar.value = await GetGameBarStatus()
})

function showRestartNotice(title: string) {
  notify.create({
    title,
    description: i18n('perf.restartRequired'),
    duration: 6000,
  })
}

async function onPowerPlan(v: boolean) {
  if (v) await EnableUltimatePerformance(); else await DisableUltimatePerformance()
  reviPlan.value = await GetUltimatePerformanceStatus()
}
async function onCState(v: boolean) {
  if (v) await EnableCState(); else await DisableCState()
  cstate.value = await GetCStateStatus()
}
async function onSuperfetch(v: boolean) {
  if (v) await EnableSuperfetch(); else await DisableSuperfetch()
  superfetch.value = await GetSuperfetchStatus()
}
async function onMemCompress(v: boolean) {
  if (v) await EnableMemoryCompression(); else await DisableMemoryCompression()
  memCompress.value = await GetMemoryCompressionStatus()
  showRestartNotice(i18n('perf.memCompress'))
}
async function onFullscreen(v: boolean) {
  if (v) await EnableFullscreenOptimization(); else await DisableFullscreenOptimization()
  fullscreen.value = await GetFullscreenOptimizationStatus()
}
async function onWindowed(v: boolean) {
  if (v) await EnableWindowedOptimization(); else await DisableWindowedOptimization()
  windowOpt.value = await GetWindowedOptimizationStatus()
}
async function onMpo(v: boolean) {
  if (v) await EnableMpo(); else await DisableMpo()
  mpo.value = await GetMpoStatus()
}

async function onClearShaderCache() {
  const msg = await ClearShaderCache()
  notify.create({ title: i18n('perf.shaderCache'), description: msg, duration: 5000 })
}

async function onGameBar(v: boolean) {
  if (v) await EnableGameBar(); else await DisableGameBar()
  gameBar.value = await GetGameBarStatus()
}
</script>

<template>
  <div class="page">
    <div class="setting-card">
      <div class="setting-card-header setting-card-header--flat"><span class="header-title">{{ i18n('perf.powerPlan') }}</span></div>
      <div class="setting-row">
        <Zap :size="18" class="row-icon" />
        <div>
          <div class="row-label">{{ i18n('perf.ultimatePerf') }}</div>
          <div class="row-desc">{{ i18n('perf.ultimatePerfDesc') }}</div>
        </div><n-switch v-model:value="reviPlan" @update:value="onPowerPlan" />
      </div>
      <div class="setting-row">
        <Cpu :size="18" class="row-icon" />
        <div>
          <div class="row-label">{{ i18n('perf.cState') }}</div>
          <div class="row-desc">{{ i18n('perf.cStateDesc') }}</div>
        </div><n-switch v-model:value="cstate" @update:value="onCState" />
      </div>
    </div>
    <div class="setting-card">
      <div class="setting-card-header setting-card-header--flat"><span class="header-title">{{ i18n('perf.memory') }}</span></div>
      <div class="setting-row">
        <Database :size="18" class="row-icon" />
        <div>
          <div class="row-label">{{ i18n('perf.superfetch') }}</div>
          <div class="row-desc">{{ i18n('perf.superfetchDesc') }}</div>
        </div><n-switch v-model:value="superfetch" @update:value="onSuperfetch" />
      </div>
      <div class="setting-row">
        <Archive :size="18" class="row-icon" />
        <div>
          <div class="row-label">{{ i18n('perf.memCompress') }}</div>
          <div class="row-desc">{{ i18n('perf.memCompressDesc') }}</div>
        </div><n-switch v-model:value="memCompress" @update:value="onMemCompress" />
      </div>
    </div>
    <div class="setting-card">
      <div class="setting-card-header setting-card-header--flat"><span class="header-title">{{ i18n('perf.display') }}</span></div>
      <div class="setting-row">
        <Monitor :size="18" class="row-icon" />
        <div>
          <div class="row-label">{{ i18n('perf.fullscreen') }}</div>
          <div class="row-desc">{{ i18n('perf.fullscreenDesc') }}</div>
        </div><n-switch v-model:value="fullscreen" @update:value="onFullscreen" />
      </div>
      <div class="setting-row">
        <Layout :size="18" class="row-icon" />
        <div>
          <div class="row-label">{{ i18n('perf.windowed') }}</div>
          <div class="row-desc">{{ i18n('perf.windowedDesc') }}</div>
        </div><n-switch v-model:value="windowOpt" @update:value="onWindowed" />
      </div>
      <div class="setting-row">
        <Layers :size="18" class="row-icon" />
        <div>
          <div class="row-label">{{ i18n('perf.mpo') }}</div>
          <div class="row-desc">{{ i18n('perf.mpoDesc') }}</div>
        </div><n-switch v-model:value="mpo" @update:value="onMpo" />
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
        <n-switch v-model:value="gameBar" @update:value="onGameBar" />
      </div>
    </div>
  </div>
</template>
