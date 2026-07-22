<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { NButton, NSelect, NSwitch, NModal, useNotification } from 'naive-ui'
import {
  RestoreDefender, DisableDefenderEngine, DisableAllServices,
  GetUacStatus, EnableUac, DisableUac,
  GetVbsStatus, EnableVbs, DisableVbs,
  GetMemoryIntegrityStatus, EnableMemoryIntegrity, DisableMemoryIntegrity,
} from '../../wailsjs/go/main/App'
import WaitModal from '../components/WaitModal.vue'

const { t: i18n } = useI18n()
const notify = useNotification()

const showModal = ref(false)
const modalText = ref('')
const showHelp = ref(false)
const selected = ref<string | null>(null)

const uac = ref(false)
const vbs = ref(false)
const memIntegrity = ref(false)

const options = [
  { label: i18n('sec.restore'), value: 'restore' },
  { label: i18n('sec.disableEngine'), value: 'disableEngine' },
  { label: i18n('sec.disableAll'), value: 'disableAll' },
]

onMounted(async () => {
  uac.value = await GetUacStatus()
  vbs.value = await GetVbsStatus()
  memIntegrity.value = await GetMemoryIntegrityStatus()
})

async function execute() {
  if (!selected.value) return

  modalText.value = i18n('sec.operating')
  showModal.value = true

  try {
    switch (selected.value) {
      case 'restore': await RestoreDefender(); break
      case 'disableEngine': await DisableDefenderEngine(); break
      case 'disableAll': await DisableAllServices(); break
    }
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
      <div class="setting-row">
        <div>
          <div class="row-label">{{ i18n('sec.defenderTitle') }}</div>
          <div class="row-desc">{{ i18n('sec.defenderDesc') }}</div>
        </div>
        <div class="action-group">
          <n-button
            quaternary
            circle
            class="help-btn"
            @click="showHelp = true"
          >?</n-button>
          <n-select
            v-model:value="selected"
            :options="options"
            :placeholder="i18n('sec.selectAction')"
            :disabled="showModal"
            style="width: 200px"
          />
          <n-button
            size="small"
            :disabled="!selected || showModal"
            @click="execute"
          >
            {{ i18n('sec.execute') }}
          </n-button>
        </div>
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

    <n-modal v-model:show="showHelp" preset="card" :title="i18n('sec.helpTitle')" style="max-width: 480px">
      <div class="help-content">
        <div class="help-item">
          <div class="help-label">{{ i18n('sec.restore') }}</div>
          <div class="help-desc">{{ i18n('sec.helpRestore') }}</div>
        </div>
        <div class="help-item">
          <div class="help-label">{{ i18n('sec.disableEngine') }}</div>
          <div class="help-desc">{{ i18n('sec.helpDisableEngine') }}</div>
        </div>
        <div class="help-item">
          <div class="help-label">{{ i18n('sec.disableAll') }}</div>
          <div class="help-desc">{{ i18n('sec.helpDisableAll') }}</div>
        </div>
      </div>
    </n-modal>
  </div>
</template>

<style scoped>
.action-group {
  display: flex;
  align-items: center;
  gap: 4px;
}
.help-btn {
  width: 32px; height: 32px;
  font-size: 14px; font-weight: 600;
  flex-shrink: 0;
}
.help-content {
  display: flex; flex-direction: column; gap: 16px;
}
.help-item {
  padding: 12px 16px;
  background: var(--section-bg);
  border-radius: 8px;
  border: 1px solid var(--border);
}
.help-label {
  font-weight: 600; margin-bottom: 4px; font-size: 14px;
}
.help-desc {
  font-size: 13px; color: var(--text2); line-height: 1.5;
}
</style>
