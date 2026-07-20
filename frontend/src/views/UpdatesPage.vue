<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { NSwitch, NButton } from 'naive-ui'
import {
  GetPauseUpdatesStatus,
  EnablePauseUpdates,
  DisablePauseUpdates,
  GetVisibilityStatus,
  EnableVisibility,
  DisableVisibility,
  GetDriverUpdatesStatus,
  EnableDriverUpdates,
  DisableDriverUpdates,
  UpdateCertificates,
  UpdateKGL,
} from '../../wailsjs/go/main/App'

const pauseUpdates = ref(false)
const visibility = ref(false)
const drivers = ref(true)

onMounted(async () => {
  pauseUpdates.value = await GetPauseUpdatesStatus()
  visibility.value = await GetVisibilityStatus()
  drivers.value = await GetDriverUpdatesStatus()
})

async function doCerts() {
  const msg = await UpdateCertificates()
  alert(msg)
}

async function doKGL() {
  const msg = await UpdateKGL()
  alert(msg)
}

async function onPauseChange(v: boolean) {
  if (v) await EnablePauseUpdates() ; else await DisablePauseUpdates()
  pauseUpdates.value = await GetPauseUpdatesStatus()
}

async function onVisibilityChange(v: boolean) {
  if (v) await EnableVisibility(); else await DisableVisibility()
  visibility.value = await GetVisibilityStatus()
}

async function onDriversChange(v: boolean) {
  if (v) await EnableDriverUpdates() ; else await DisableDriverUpdates()
  drivers.value = await GetDriverUpdatesStatus()
}
</script>

<template>
  <div class="page">
    <h2>更新</h2>
    <div class="setting-card">
      <div class="setting-card-header"><span class="header-title">组件更新</span><span class="header-desc">证书和知识库更新</span></div>
      <div class="setting-row">
        <div><div class="row-label">更新证书</div><div class="row-desc">从 Windows Update 拉取最新根证书</div></div>
        <n-button size="small" @click="doCerts">更新</n-button>
      </div>
      <div class="setting-row">
        <div><div class="row-label">更新 KGL</div><div class="row-desc">获取最新知识图谱库</div></div>
        <n-button size="small" @click="doKGL">更新</n-button>
      </div>
    </div>
    <div class="setting-card">
      <div class="setting-card-header"><span class="header-title">Windows 更新</span><span class="header-desc">更新策略控制</span></div>
      <div class="setting-row"><div><div class="row-label">暂停更新</div><div class="row-desc">暂停到 2126 年</div></div><n-switch v-model:value="pauseUpdates" @update:value="onPauseChange" /></div>
      <div class="setting-row"><div><div class="row-label">可见性</div><div class="row-desc">显示或隐藏设置中的更新页面</div></div><n-switch v-model:value="visibility" @update:value="onVisibilityChange" /></div>
      <div class="setting-row"><div><div class="row-label">驱动程序</div><div class="row-desc">允许或禁止通过 Windows Update 安装驱动</div></div><n-switch v-model:value="drivers" @update:value="onDriversChange" /></div>
    </div>
  </div>
</template>
