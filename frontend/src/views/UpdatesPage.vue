<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { NSwitch, NButton, NSelect, NModal } from 'naive-ui'
import { FileCheck, Book, PauseCircle, EyeOff, Truck, Radio } from 'lucide-vue-next'
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

import { useNotify } from '../composables/useNotify'

const { t: i18n } = useI18n()
const notify = useNotify()
const pauseUpdates = ref(false)
const visibility = ref(false)
const drivers = ref(true)
const channel = ref('retail')
const showChannelConfirm = ref(false)
const pendingChannel = ref('')

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
  notify.create({ title: i18n('upd.certsTitle'), description: msg, duration: 5000 })
}

async function doKGL() {
  const msg = await UpdateKGL()
  notify.create({ title: i18n('upd.kglTitle'), description: msg, duration: 5000 })
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
  pendingChannel.value = v
  showChannelConfirm.value = true
}

async function confirmChannel() {
  showChannelConfirm.value = false
  await SetUpdateChannel(pendingChannel.value)
  channel.value = await GetUpdateChannel()
}

async function cancelChannel() {
  showChannelConfirm.value = false
  channel.value = await GetUpdateChannel()
}
</script>

<template>
  <div class="page">
    <div class="setting-card">
      <div class="setting-card-header setting-card-header--flat"><span class="header-title">{{ i18n('upd.components') }}</span></div>
      <div class="setting-row">
        <FileCheck :size="18" class="row-icon" />
        <div><div class="row-label">{{ i18n('upd.certs') }}</div><div class="row-desc">{{ i18n('upd.certsDesc') }}</div></div>
        <n-button size="small" @click="doCerts">{{ i18n('upd.update') }}</n-button>
      </div>
      <div class="setting-row">
        <Book :size="18" class="row-icon" />
        <div><div class="row-label">{{ i18n('upd.kgl') }}</div><div class="row-desc">{{ i18n('upd.kglDesc') }}</div></div>
        <n-button size="small" @click="doKGL">{{ i18n('upd.update') }}</n-button>
      </div>
    </div>
    <div class="setting-card">
      <div class="setting-card-header setting-card-header--flat"><span class="header-title">{{ i18n('upd.windowsUpdate') }}</span></div>
      <div class="setting-row">
        <PauseCircle :size="18" class="row-icon" />
        <div><div class="row-label">{{ i18n('upd.pause') }}</div><div class="row-desc">{{ i18n('upd.pauseDesc') }}</div></div><n-switch v-model:value="pauseUpdates" @update:value="onPauseChange" />
      </div>
      <div class="setting-row">
        <EyeOff :size="18" class="row-icon" />
        <div><div class="row-label">{{ i18n('upd.hidePage') }}</div><div class="row-desc">{{ i18n('upd.hidePageDesc') }}</div></div><n-switch v-model:value="visibility" @update:value="onVisibilityChange" />
      </div>
      <div class="setting-row">
        <Truck :size="18" class="row-icon" />
        <div><div class="row-label">{{ i18n('upd.drivers') }}</div><div class="row-desc">{{ i18n('upd.driversDesc') }}</div></div><n-switch v-model:value="drivers" @update:value="onDriversChange" />
      </div>
      <div v-if="!pauseUpdates" class="setting-row">
        <Radio :size="18" class="row-icon" />
        <div>
          <div class="row-label">{{ i18n('upd.channel') }}</div>
          <div class="row-desc">{{ i18n('upd.channelDesc') }}</div>
        </div>
        <n-select v-model:value="channel" :options="channelOptions" style="width:160px" @update:value="onChannelChange" />
      </div>
    </div>

    <n-modal
      v-model:show="showChannelConfirm"
      preset="card"
      :title="i18n('upd.channel')"
      style="width: 600px; max-width: 90vw;"
      :mask-closable="false"
      closable
      @close="cancelChannel"
    >
      <p>{{ i18n('upd.channelConfirm') }}</p>
      <template #footer>
        <div style="display: flex; justify-content: flex-end; gap: 8px;">
          <n-button @click="cancelChannel">{{ i18n('home.cancel') }}</n-button>
          <n-button @click="confirmChannel">{{ i18n('upd.update') }}</n-button>
        </div>
      </template>
    </n-modal>
  </div>
</template>
