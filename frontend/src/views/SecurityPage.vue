<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { NSwitch, useNotification } from 'naive-ui'
import {
  GetSecurityHealthServiceStatus, RestoreDefender, DisableAllServices,
  GetUacStatus, EnableUac, DisableUac,
  GetVbsStatus, EnableVbs, DisableVbs,
  GetMemoryIntegrityStatus, EnableMemoryIntegrity, DisableMemoryIntegrity,
} from '../../wailsjs/go/main/App'
import WaitModal from '../components/WaitModal.vue'

const { t: i18n } = useI18n()
const notify = useNotification()

const showModal = ref(false)
const modalText = ref('')
const serviceDisabled = ref(false)

const uac = ref(false)
const vbs = ref(false)
const memIntegrity = ref(false)

onMounted(async () => {
  serviceDisabled.value = await GetSecurityHealthServiceStatus()
  uac.value = await GetUacStatus()
  vbs.value = await GetVbsStatus()
  memIntegrity.value = await GetMemoryIntegrityStatus()
})

async function onServiceToggle(v: boolean) {
  modalText.value = i18n('sec.operating')
  showModal.value = true

  try {
    if (v) await RestoreDefender(); else await DisableAllServices()
    notify.success({ title: i18n('sec.notifyTitle'), description: i18n('sec.restartRequired'), duration: 6000 })
  } catch (e: any) {
    notify.error({ title: i18n('sec.notifyTitle'), description: e?.message || 'Error', duration: 6000 })
  } finally {
    showModal.value = false
  }
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
    <WaitModal :show="showModal" :text="modalText" />
    <h2>{{ i18n('sec.title') }}</h2>

    <div class="setting-card">
      <div class="setting-card-header">
        <span class="header-title">{{ i18n('sec.defenderTitle') }}</span>
      </div>
      <div class="setting-row">
        <div>
          <div class="row-label">{{ i18n('sec.disableAll') }}</div>
          <div class="row-desc">{{ i18n('sec.disableAllDesc') }}</div>
        </div>
        <n-switch v-model:value="serviceDisabled" :disabled="showModal" @update:value="onServiceToggle" />
      </div>
    </div>

    <div class="setting-card">
      <div class="setting-card-header">
        <span class="header-title">{{ i18n('sec.systemProtection') }}</span>
        <span class="header-desc">{{ i18n('sec.systemProtectionDesc') }}</span>
      </div>
      <div class="setting-row">
        <div>
          <div class="row-label">{{ i18n('sec.uac') }}</div>
          <div class="row-desc">{{ i18n('sec.uacDesc') }}</div>
        </div>
        <n-switch v-model:value="uac" @update:value="onUac" />
      </div>
      <div class="setting-row">
        <div>
          <div class="row-label">{{ i18n('sec.vbs') }}</div>
          <div class="row-desc">{{ i18n('sec.vbsDesc') }}</div>
        </div>
        <n-switch v-model:value="vbs" @update:value="onVbs" />
      </div>
      <div class="setting-row">
        <div>
          <div class="row-label">{{ i18n('sec.memIntegrity') }}</div>
          <div class="row-desc">{{ i18n('sec.memIntegrityDesc') }}</div>
        </div>
        <n-switch v-model:value="memIntegrity" @update:value="onMemIntegrity" />
      </div>
    </div>
  </div>
</template>
