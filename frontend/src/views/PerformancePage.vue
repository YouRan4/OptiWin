<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { NSwitch, useNotification } from 'naive-ui'
import {
  GetUltimatePerformanceStatus, EnableUltimatePerformance, DisableUltimatePerformance,
  GetCStateStatus, EnableCState, DisableCState,
  GetSuperfetchStatus, EnableSuperfetch, DisableSuperfetch,
  GetFullscreenOptimizationStatus, EnableFullscreenOptimization, DisableFullscreenOptimization,
  GetWindowedOptimizationStatus, EnableWindowedOptimization, DisableWindowedOptimization,
  GetMpoStatus, EnableMpo, DisableMpo,
  GetMemoryCompressionStatus, EnableMemoryCompression, DisableMemoryCompression,
} from '../../wailsjs/go/main/App'

const notify = useNotification()

const reviPlan = ref(false)
const cstate = ref(true)
const superfetch = ref(true)
const memCompress = ref(true)
const fullscreen = ref(true)
const windowOpt = ref(true)
const mpo = ref(true)

onMounted(async () => {
  reviPlan.value = await GetUltimatePerformanceStatus()
  cstate.value = await GetCStateStatus()
  superfetch.value = await GetSuperfetchStatus()
  memCompress.value = await GetMemoryCompressionStatus()
  fullscreen.value = await GetFullscreenOptimizationStatus()
  windowOpt.value = await GetWindowedOptimizationStatus()
  mpo.value = await GetMpoStatus()
})

function showRestartNotice(title: string) {
  notify.success({
    title,
    description: '需要重启系统才能生效',
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
  showRestartNotice('内存压缩')
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
</script>

<template>
  <div class="page">
    <h2>性能优化</h2>
    <div class="setting-card">
      <div class="setting-card-header"><span class="header-title">电源计划</span><span class="header-desc">控制 CPU 电源管理</span></div>
      <div class="setting-row"><div><div class="row-label">卓越性能</div><div class="row-desc">启用 Windows 内置的卓越性能电源计划</div></div><n-switch v-model:value="reviPlan" @update:value="onPowerPlan" /></div>
      <div class="setting-row"><div><div class="row-label">禁用 C State</div><div class="row-desc">禁用 CPU 深度睡眠状态以减少延迟</div></div><n-switch v-model:value="cstate" @update:value="onCState" /></div>
    </div>
    <div class="setting-card">
      <div class="setting-card-header"><span class="header-title">内存与存储</span><span class="header-desc">优化内存压缩和预缓存</span></div>
      <div class="setting-row"><div><div class="row-label">Superfetch</div><div class="row-desc">预加载常用应用到内存，加快启动速度</div></div><n-switch v-model:value="superfetch" @update:value="onSuperfetch" /></div>
      <div class="setting-row"><div><div class="row-label">内存压缩</div><div class="row-desc">压缩未用内存页面</div></div><n-switch v-model:value="memCompress" @update:value="onMemCompress" /></div>
    </div>
    <div class="setting-card">
      <div class="setting-card-header"><span class="header-title">显示</span><span class="header-desc">优化游戏和应用渲染</span></div>
      <div class="setting-row"><div><div class="row-label">全屏优化</div><div class="row-desc">为全屏游戏和应用提供更低延迟</div></div><n-switch v-model:value="fullscreen" @update:value="onFullscreen" /></div>
      <div class="setting-row"><div><div class="row-label">窗口优化</div><div class="row-desc">优化窗口化应用的渲染和输入延迟</div></div><n-switch v-model:value="windowOpt" @update:value="onWindowed" /></div>
      <div class="setting-row"><div><div class="row-label">MPO</div><div class="row-desc">多平面覆盖技术，减少窗口化游戏的输入延迟</div></div><n-switch v-model:value="mpo" @update:value="onMpo" /></div>
    </div>
  </div>
</template>
