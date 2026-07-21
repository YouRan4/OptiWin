<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
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

const { t: i18n } = useI18n()
const notify = useNotification()
const pauseUpdates = ref(false)
const visibility = ref(false)
const drivers = ref(true)
const channel = ref('retail')

const channelOptions = computed(() => [
  { label: i18n('upd.exitInsider'), value: 'retail' },
  { label: 'Release Preview', value: 'ReleasePreview' },
  { label: i18n('upd.beta'), value: 'Beta' },
  { label: i18n('upd.dev'), value: 'Dev' },
  { label: i18n('upd.canary'), value: 'Canary' },
])

onMounted(async () => {
  pauseUpdates.value = await GetPauseUpdatesStatus()
  visibility.value = await GetVisibilityStatus()
  drivers.value = await GetDriverUpdatesStatus()
  channel.value = await GetUpdateChannel()
})

async function doCerts() {
  const msg = await UpdateCertificates()
  notify.info({ title: i18n('upd.certsTitle'), description: msg, duration: 5000 })
}

async function doKGL() {
  const msg = await UpdateKGL()
  notify.info({ title: i18n('upd.kglTitle'), description: msg, duration: 5000 })
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
  const ok = confirm(i18n('upd.channelConfirm'))
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
    <h2>{{ i18n('upd.title') }}</h2>
    <div class="setting-card">
      <div class="setting-card-header"><span class="header-title">{{ i18n('upd.components') }}</span><span class="header-desc">{{ i18n('upd.componentsDesc') }}</span></div>
      <div class="setting-row">
        <div><div class="row-label">{{ i18n('upd.certs') }}</div><div class="row-desc">{{ i18n('upd.certsDesc') }}</div></div>
        <n-button size="small" @click="doCerts">{{ i18n('upd.update') }}</n-button>
      </div>
      <div class="setting-row">
        <div><div class="row-label">{{ i18n('upd.kgl') }}</div><div class="row-desc">{{ i18n('upd.kglDesc') }}</div></div>
        <n-button size="small" @click="doKGL">{{ i18n('upd.update') }}</n-button>
      </div>
    </div>
    <div class="setting-card">
      <div class="setting-card-header"><span class="header-title">{{ i18n('upd.windowsUpdate') }}</span><span class="header-desc">{{ i18n('upd.windowsUpdateDesc') }}</span></div>
      <div class="setting-row"><div><div class="row-label">{{ i18n('upd.pause') }}</div><div class="row-desc">{{ i18n('upd.pauseDesc') }}</div></div><n-switch v-model:value="pauseUpdates" @update:value="onPauseChange" /></div>
      <div class="setting-row"><div><div class="row-label">{{ i18n('upd.hidePage') }}</div><div class="row-desc">{{ i18n('upd.hidePageDesc') }}</div></div><n-switch v-model:value="visibility" @update:value="onVisibilityChange" /></div>
      <div class="setting-row"><div><div class="row-label">{{ i18n('upd.drivers') }}</div><div class="row-desc">{{ i18n('upd.driversDesc') }}</div></div><n-switch v-model:value="drivers" @update:value="onDriversChange" /></div>
      <div v-if="!pauseUpdates" class="setting-row">
        <div>
          <div class="row-label">{{ i18n('upd.channel') }}</div>
          <div class="row-desc">{{ i18n('upd.channelDesc') }}</div>
        </div>
        <n-select v-model:value="channel" :options="channelOptions" style="width:160px" @update:value="onChannelChange" />
      </div>
    </div>
  </div>
</template>
