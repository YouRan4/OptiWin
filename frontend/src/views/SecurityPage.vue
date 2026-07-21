<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { NSwitch, NSpin, useNotification } from 'naive-ui'
import {
  GetDefenderStatus, EnableDefender, DisableDefender,
  GetUacStatus, EnableUac, DisableUac,
  GetVbsStatus, EnableVbs, DisableVbs,
  GetMemoryIntegrityStatus, EnableMemoryIntegrity, DisableMemoryIntegrity,
} from '../../wailsjs/go/main/App'

const { t: i18n } = useI18n()
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
  notify.success({ title: i18n('sec.notifyTitle'), description: i18n('sec.restartRequired'), duration: 6000 })
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
        <template #description>{{ i18n('sec.defenderOperating') }}</template>
        <div class="loading-overlay" v-if="defenderLoading"></div>
      </n-spin>
      <h2>{{ i18n('sec.title') }}</h2>
      <div class="setting-card">
        <div class="setting-card-header"><span class="header-title">{{ i18n('sec.defenderTitle') }}</span><span class="header-desc">{{ i18n('sec.defenderDesc') }}</span></div>
        <div class="setting-row"><div><div class="row-label">{{ i18n('sec.defenderLabel') }}</div><div class="row-desc">{{ i18n('sec.defenderRowDesc') }}</div></div><n-switch :value="defender" :disabled="defenderLoading" @update:value="onDefender" /></div>
      </div>
    </div>
    <div class="setting-card">
      <div class="setting-card-header"><span class="header-title">{{ i18n('sec.systemProtection') }}</span><span class="header-desc">{{ i18n('sec.systemProtectionDesc') }}</span></div>
      <div class="setting-row"><div><div class="row-label">{{ i18n('sec.uac') }}</div><div class="row-desc">{{ i18n('sec.uacDesc') }}</div></div><n-switch v-model:value="uac" @update:value="onUac" /></div>
      <div class="setting-row"><div><div class="row-label">{{ i18n('sec.vbs') }}</div><div class="row-desc">{{ i18n('sec.vbsDesc') }}</div></div><n-switch v-model:value="vbs" @update:value="onVbs" /></div>
      <div class="setting-row"><div><div class="row-label">{{ i18n('sec.memIntegrity') }}</div><div class="row-desc">{{ i18n('sec.memIntegrityDesc') }}</div></div><n-switch v-model:value="memIntegrity" @update:value="onMemIntegrity" /></div>
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
