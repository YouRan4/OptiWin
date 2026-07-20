<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { NSwitch, NButton, NSelect, useNotification } from 'naive-ui'
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
  GetUpdateChannel,
  SetUpdateChannel,
  UpdateCertificates,
  UpdateKGL,
} from '../../wailsjs/go/main/App'

const notify = useNotification()
const pauseUpdates = ref(false)
const visibility = ref(false)
const drivers = ref(true)
const channel = ref('retail')

const channelOptions = [
  { label: '退出 Insider', value: 'retail' },
  { label: 'Release Preview', value: 'ReleasePreview' },
  { label: 'Beta 通道', value: 'Beta' },
  { label: 'Dev 通道', value: 'Dev' },
  { label: 'Canary 通道', value: 'Canary' },
]

onMounted(async () => {
  pauseUpdates.value = await GetPauseUpdatesStatus()
  visibility.value = await GetVisibilityStatus()
  drivers.value = await GetDriverUpdatesStatus()
  channel.value = await GetUpdateChannel()
})

async function doCerts() {
  const msg = await UpdateCertificates()
  notify.info({ title: '更新证书', description: msg, duration: 5000 })
}

async function doKGL() {
  const msg = await UpdateKGL()
  notify.info({ title: '更新 KGL', description: msg, duration: 5000 })
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

async function onChannelChange(v: string) {
  const ok = confirm('切换 Windows 更新通道需要重启生效，是否继续？')
  if (!ok) {
    channel.value = await GetUpdateChannel()
    return
  }
  await SetUpdateChannel(v)
  channel.value = await GetUpdateChannel()
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
      <div class="setting-row"><div><div class="row-label">暂停Windows更新</div><div class="row-desc">暂停到 2126 年</div></div><n-switch v-model:value="pauseUpdates" @update:value="onPauseChange" /></div>
      <div class="setting-row"><div><div class="row-label">隐藏Windows更新页面</div><div class="row-desc">隐藏设置中的更新页面与更新通知</div></div><n-switch v-model:value="visibility" @update:value="onVisibilityChange" /></div>
      <div class="setting-row"><div><div class="row-label">通过Windows更新安装驱动程序</div><div class="row-desc">允许 Windows Update 自动安装驱动程序</div></div><n-switch v-model:value="drivers" @update:value="onDriversChange" /></div>
      <div v-if="!pauseUpdates" class="setting-row">
        <div>
          <div class="row-label">更新通道</div>
          <div class="row-desc">切换 Windows Insider 预览通道（需要微软账号已注册 Insider 计划）</div>
        </div>
        <n-select v-model:value="channel" :options="channelOptions" style="width:160px" @update:value="onChannelChange" />
      </div>
    </div>
  </div>
</template>
