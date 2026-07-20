<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { NSwitch, NSpin, useNotification } from 'naive-ui'
import {
  GetDefenderStatus, EnableDefender, DisableDefender,
  GetUacStatus, EnableUac, DisableUac,
  GetVbsStatus, EnableVbs, DisableVbs,
  GetMemoryIntegrityStatus, EnableMemoryIntegrity, DisableMemoryIntegrity,
} from '../../wailsjs/go/main/App'

const notify = useNotification()

const defender = ref(false)
const defenderLoading = ref(false)
const uac = ref(false)
const vbs = ref(false)
const memIntegrity = ref(false)

onMounted(async () => {
  defender.value = await GetDefenderStatus()
  uac.value = await GetUacStatus()
  vbs.value = await GetVbsStatus()
  memIntegrity.value = await GetMemoryIntegrityStatus()
})

async function onDefender(v: boolean) {
  defenderLoading.value = true
  if (v) await EnableDefender(); else await DisableDefender()
  defender.value = await GetDefenderStatus()
  defenderLoading.value = false
  notify.success({ title: '安全中心', description: '需要重启系统才能完全生效', duration: 6000 })
}

async function onUac(v: boolean) {
  if (v) await EnableUac(); else await DisableUac()
  uac.value = await GetUacStatus()
}

async function onVbs(v: boolean) {
  if (v) await EnableVbs(); else await DisableVbs()
  vbs.value = await GetVbsStatus()
}

async function onMemIntegrity(v: boolean) {
  if (v) await EnableMemoryIntegrity(); else await DisableMemoryIntegrity()
  memIntegrity.value = await GetMemoryIntegrityStatus()
}
</script>

<template>
  <div class="page">
    <div class="defender-section" style="position:relative">
      <n-spin :show="defenderLoading">
        <template #description>正在操作安全中心...</template>
        <div class="loading-overlay" v-if="defenderLoading"></div>
      </n-spin>
      <h2>安全</h2>
      <div class="setting-card">
        <div class="setting-card-header"><span class="header-title">安全中心</span><span class="header-desc">Windows 安全中心完全启用或禁用</span></div>
        <div class="setting-row"><div><div class="row-label">Windows 安全中心</div><div class="row-desc">完全开启或关闭 Defender + 安全中心服务</div></div><n-switch :value="defender" :disabled="defenderLoading" @update:value="onDefender" /></div>
      </div>
    </div>
    <div class="setting-card">
      <div class="setting-card-header"><span class="header-title">系统防护</span><span class="header-desc">UAC、VBS、内存完整性设置</span></div>
      <div class="setting-row"><div><div class="row-label">UAC</div><div class="row-desc">用户帐户控制，限制应用对系统更改的权限</div></div><n-switch v-model:value="uac" @update:value="onUac" /></div>
      <div class="setting-row"><div><div class="row-label">VBS</div><div class="row-desc">基于虚拟化的安全，使用硬件虚拟化增强系统安全</div></div><n-switch v-model:value="vbs" @update:value="onVbs" /></div>
      <div class="setting-row"><div><div class="row-label">内存完整性</div><div class="row-desc">核心隔离功能，防止恶意代码注入高安全性进程</div></div><n-switch v-model:value="memIntegrity" @update:value="onMemIntegrity" /></div>
    </div>
  </div>
</template>

<style scoped>
.loading-overlay {
  position: fixed; inset: 0; z-index: 999;
  background: rgba(0,0,0,0.3);
  pointer-events: none;
}
</style>
