<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { NSwitch, NButton, NModal, NInput, NSelect } from 'naive-ui'
import { Shield, Lock, Server, HardDrive, Search, RefreshCw, Trash2, Plus, Bug, AlertTriangle } from 'lucide-vue-next'
import {
  GetUacStatus, EnableUac, DisableUac,
  GetCoreServicesDisabled, DisableCoreServices, EnableCoreServices,
  GetTamperProtectionStatus, GetRealtimeProtectionEnabled, OpenWindowsSecurity,
  GetVbsStatus, EnableVbs, DisableVbs,
  GetMemoryIntegrityStatus, EnableMemoryIntegrity, DisableMemoryIntegrity,
  ListIfeoEntries, AddIfeoEntry, RemoveIfeoEntry, GetRunningProcesses, SetDns, GetCurrentDns,
} from '../../wailsjs/go/main/App'
import WaitModal from '../components/WaitModal.vue'
import { useNotify } from '../composables/useNotify'
import { useBackendToggle } from '../composables/useBackendToggle'

const { t: i18n } = useI18n()
const notify = useNotify()

const showModal = ref(false)
const modalText = ref('')
const showDefenderWarn = ref(false)

const { value: uac, refresh: refreshUac, onUpdate: onUac, busy: uacBusy } = useBackendToggle(
  GetUacStatus, v => (v ? EnableUac() : DisableUac()), () => i18n('sec.uac'),
)
const { value: vbs, refresh: refreshVbs, onUpdate: onVbs, busy: vbsBusy } = useBackendToggle(
  GetVbsStatus, v => (v ? EnableVbs() : DisableVbs()), () => i18n('sec.vbs'),
)
const { value: memIntegrity, refresh: refreshMemIntegrity, onUpdate: onMemIntegrity, busy: memIntegrityBusy } = useBackendToggle(
  GetMemoryIntegrityStatus,
  v => (v ? EnableMemoryIntegrity() : DisableMemoryIntegrity()),
  () => i18n('sec.memIntegrity'),
)
const coreServicesDisabled = ref(false)

// IFEO
interface IfeoEntry {
  name: string
  debugger: string
  globalFlag: number
}

const ifeoEntries = ref<IfeoEntry[]>([])
const ifeoSearch = ref('')
const ifeoLoading = ref(false)
const showAddModal = ref(false)
const showDeleteModal = ref(false)
const newExeName = ref('')
const newDebugger = ref('')
const deleteTarget = ref('')
const expandIfeo = ref(false)
const runningProcesses = ref<string[]>([])

const blockedProcesses = [
  'csrss.exe', 'smss.exe', 'lsass.exe', 'services.exe',
  'wininit.exe', 'winlogon.exe', 'svchost.exe', 'dwm.exe',
  'conhost.exe', 'fontdrvhost.exe', 'runtimebroker.exe',
  'sihost.exe', 'shellexperiencehost.exe', 'startmenuexperiencehost.exe',
  'searchui.exe', 'searchindexer.exe', 'msmpeng.exe', 'taskhostw.exe',
  'ntoskrnl.exe', 'hal.dll', 'bootvid.dll', 'win32k.sys',
]

const telemetryProcesses = [
  'AggregatorHost.exe', 'BCILauncher.exe', 'BGAUpsell.exe',
  'BingChatInstaller.exe', 'CompatTelRunner.exe', 'DeviceCensus.exe',
  'FeatureLoader.exe'
]

const telemetryDebuggerPath = '%windir%\\system32\\taskkill.exe'

// 开关状态由实际 IFEO 条目推导，而不是在加载时算一次存进 ref；
// 7 个进程全部被劫持才算已禁用，部分成功如实显示为未启用
const telemetryBlocked = computed(() => {
  const entries = ifeoEntries.value ?? []
  return telemetryProcesses.every(p => entries.some(e => e.name.toLowerCase() === p.toLowerCase()))
})

// DNS
const dnsOptions = computed(() => [
  { label: i18n('sec.dns_isp'), value: 'isp' },
  { label: i18n('sec.dns_ali'), value: 'ali' },
  { label: i18n('sec.dns_tencent'), value: 'tencent' },
  { label: i18n('sec.dns_baidu'), value: 'baidu' },
  { label: i18n('sec.dns_google'), value: 'google' },
  { label: i18n('sec.dns_cloudflare'), value: 'cloudflare' },
  { label: i18n('sec.dns_quad9'), value: 'quad9' },
])

const selectedDns = ref('isp')
const currentDnsInfo = ref({ adapter: '', primaryDns: '', secondaryDns: '' })

async function loadCurrentDns() {
  try {
    const result = await GetCurrentDns()
    const parsed = JSON.parse(result)
    currentDnsInfo.value = {
      adapter: parsed.adapter || '',
      primaryDns: parsed.primaryDns || '',
      secondaryDns: parsed.secondaryDns || ''
    }
  } catch {
    currentDnsInfo.value = { adapter: '', primaryDns: '', secondaryDns: '' }
  }
}

async function onDnsChange(value: string) {
  try {
    const ok = await SetDns(value)
    if (ok) {
      await loadCurrentDns()
      const label = dnsOptions.value.find(o => o.value === value)?.label
      notify.create({
        title: i18n('sec.notifyTitle'),
        description: i18n('sec.dnsSwitched', { name: label }),
        duration: 3000
      })
    } else {
      notify.create({ title: i18n('sec.notifyTitle'), description: i18n('sec.dnsSwitchFailed'), duration: 4000 })
    }
  } catch (e: any) {
    notify.create({ title: i18n('sec.notifyTitle'), description: e?.message || 'Error', duration: 4000 })
  }
}

const filteredIfeo = computed(() => {
  const entries = ifeoEntries.value ?? []
  if (!ifeoSearch.value) return entries
  const q = ifeoSearch.value.toLowerCase()
  return entries.filter(e =>
    e.name.toLowerCase().includes(q) ||
    e.debugger.toLowerCase().includes(q)
  )
})

onMounted(async () => {
  await Promise.all([refreshUac(), refreshVbs(), refreshMemIntegrity()])
  coreServicesDisabled.value = await GetCoreServicesDisabled()
  await loadIfeoEntries()
  await loadCurrentDns()
})

async function onDisableDefender() {
  // 检测实时保护与篡改保护是否都已关闭
  const tpOn = await GetTamperProtectionStatus()
  const rtOn = await GetRealtimeProtectionEnabled()
  if (tpOn || rtOn) {
    showDefenderWarn.value = true // 弹窗提示，结束执行
    return
  }

  // 均已关闭：执行禁用（实时保护 + 核心服务）
  modalText.value = i18n('sec.operating')
  showModal.value = true
  try {
    const ok = await DisableCoreServices()
    coreServicesDisabled.value = await GetCoreServicesDisabled()
    notify.create({ title: i18n('sec.notifyTitle'), description: ok ? i18n('sec.defenderDisabled') : i18n('sec.defenderFail'), duration: 4000 })
  } finally {
    showModal.value = false
  }
}

async function onEnableDefender() {
  modalText.value = i18n('sec.operating')
  showModal.value = true
  try {
    const ok = await EnableCoreServices()
    coreServicesDisabled.value = await GetCoreServicesDisabled()
    notify.create({ title: i18n('sec.notifyTitle'), description: ok ? i18n('sec.defenderEnabled') : i18n('sec.defenderFail'), duration: 4000 })
  } finally {
    showModal.value = false
  }
}

// IFEO functions
async function loadIfeoEntries() {
  ifeoLoading.value = true
  try {
    const result = await ListIfeoEntries()
    const parsed = JSON.parse(result)
    ifeoEntries.value = Array.isArray(parsed) ? parsed : []
  } catch {
    ifeoEntries.value = []
  } finally {
    ifeoLoading.value = false
  }
}

async function loadRunningProcesses() {
  try {
    const result = await GetRunningProcesses()
    const procs = JSON.parse(result) as string[]
    runningProcesses.value = procs.map(p => p.toLowerCase().endsWith('.exe') ? p : p + '.exe')
  } catch {
    runningProcesses.value = []
  }
}

async function addIfeoEntry() {
  let exeName = newExeName.value.trim()
  if (exeName && !exeName.toLowerCase().endsWith('.exe')) {
    exeName = exeName + '.exe'
  }
  const debuggerPath = newDebugger.value.trim()

  if (!exeName || !debuggerPath) {
    notify.create({ title: i18n('sec.notifyTitle'), description: i18n('sec.ifeoEmptyInput'), duration: 4000 })
    return
  }

  if (!exeName.toLowerCase().endsWith('.exe')) {
    notify.create({ title: i18n('sec.notifyTitle'), description: i18n('sec.ifeoInvalidExe'), duration: 4000 })
    return
  }

  if (/[\/:;*?<>|]/.test(exeName)) {
    notify.create({ title: i18n('sec.notifyTitle'), description: i18n('sec.ifeoInvalidChars'), duration: 4000 })
    return
  }

  const exeNameLower = exeName.toLowerCase()
  if (blockedProcesses.includes(exeNameLower)) {
    notify.create({ title: i18n('sec.notifyTitle'), description: i18n('sec.ifeoBlockedProcess'), duration: 6000 })
    return
  }

  const dbgLower = debuggerPath.toLowerCase()
  if (!dbgLower.includes(':\\') && !dbgLower.includes('%')) {
    notify.create({ title: i18n('sec.notifyTitle'), description: i18n('sec.ifeoInvalidAbsPath'), duration: 4000 })
    return
  }
  if (dbgLower.startsWith('\\\\')) {
    notify.create({ title: i18n('sec.notifyTitle'), description: i18n('sec.ifeoUncPath'), duration: 4000 })
    return
  }
  if (dbgLower.includes('..')) {
    notify.create({ title: i18n('sec.notifyTitle'), description: i18n('sec.ifeoPathTraversal'), duration: 4000 })
    return
  }

  const ok = await AddIfeoEntry(exeName, debuggerPath)
  if (ok) {
    notify.create({ title: i18n('sec.notifyTitle'), description: i18n('sec.ifeoAddSuccess'), duration: 4000 })
    showAddModal.value = false
    newExeName.value = ''
    newDebugger.value = ''
    await loadIfeoEntries()
  } else {
    notify.create({ title: i18n('sec.notifyTitle'), description: i18n('sec.ifeoAddFailed'), duration: 6000 })
  }
}

function confirmDeleteIfeo(name: string) {
  deleteTarget.value = name
  showDeleteModal.value = true
}

async function doDeleteIfeo() {
  if (!deleteTarget.value) return
  showDeleteModal.value = false
  const ok = await RemoveIfeoEntry(deleteTarget.value)
  if (ok) {
    notify.create({ title: i18n('sec.notifyTitle'), description: i18n('sec.ifeoDeleteSuccess'), duration: 4000 })
    await loadIfeoEntries()
  } else {
    notify.create({ title: i18n('sec.notifyTitle'), description: i18n('sec.ifeoDeleteFailed'), duration: 6000 })
  }
}

async function onTelemetryToggle(v: boolean) {
  modalText.value = i18n('sec.operating')
  showModal.value = true
  try {
    const results = await Promise.all(telemetryProcesses.map(proc =>
      v ? AddIfeoEntry(proc, telemetryDebuggerPath) : RemoveIfeoEntry(proc)
    ))
    const failed = results.filter(r => !r).length
    await loadIfeoEntries()
    if (failed > 0) {
      notify.create({
        title: i18n('common.updateFailed'),
        description: `${i18n('sec.telemetryTitle')}：${failed}/${telemetryProcesses.length}`,
        type: 'error',
        duration: 6000,
      })
      return
    }
    notify.create({
      title: i18n('sec.notifyTitle'),
      description: v ? i18n('sec.telemetryEnabled') : i18n('sec.telemetryDisabled'),
      duration: 4000
    })
  } catch (e: any) {
    notify.create({ title: i18n('sec.notifyTitle'), description: e?.message || 'Error', duration: 6000 })
  } finally {
    showModal.value = false
  }
}
</script>

<template>
  <div class="page">
    <WaitModal :show="showModal" :text="modalText" />
    <div class="setting-card">
      <div class="setting-card-header setting-card-header--flat">
        <span class="header-title">{{ i18n('sec.defenderModule') }}</span>
      </div>
      <div class="setting-row">
        <Shield :size="18" class="row-icon" />
        <div>
          <div class="row-label">{{ i18n('sec.defenderTitle') }}</div>
          <div class="row-desc">{{ coreServicesDisabled ? i18n('sec.defenderStatusDisabled') : i18n('sec.defenderStatusEnabled') }}</div>
        </div>
        <div style="display:flex; gap:8px; flex-shrink:0">
          <n-button size="small" strong :disabled="coreServicesDisabled" @click="onDisableDefender">
            {{ i18n('sec.defenderDisableBtn') }}
          </n-button>
          <n-button size="small" strong :disabled="!coreServicesDisabled" @click="onEnableDefender">
            {{ i18n('sec.defenderEnableBtn') }}
          </n-button>
        </div>
      </div>
      <div class="setting-row">
        <Shield :size="18" class="row-icon" />
        <div>
          <div class="row-label">{{ i18n('sec.defenderTpTitle') }}</div>
        </div>
        <n-button size="small" strong @click="OpenWindowsSecurity()">{{ i18n('sec.defenderTpBtn') }}</n-button>
      </div>
      <div style="padding:8px 12px 12px; font-size:12px; color:var(--text2); line-height:1.8">
        <div>{{ i18n('sec.defenderNote') }}</div>
        <div>1. {{ i18n('sec.defenderStep1') }}</div>
        <div>2. {{ i18n('sec.defenderStep2') }}</div>
        <div>3. {{ i18n('sec.defenderStep3') }}</div>
      </div>
    </div>

    <div class="setting-card">
      <div class="setting-card-header setting-card-header--flat">
        <span class="header-title">{{ i18n('sec.systemProtection') }}</span>
      </div>
      <div class="setting-row">
        <Lock :size="18" class="row-icon" />
        <div>
          <div class="row-label">{{ i18n('sec.uac') }}</div>
          <div class="row-desc">{{ i18n('sec.uacDesc') }}</div>
        </div>
        <n-switch v-model:value="uac" :loading="uacBusy" @update:value="onUac" />
      </div>
      <div class="setting-row">
        <Server :size="18" class="row-icon" />
        <div>
          <div class="row-label">{{ i18n('sec.vbs') }}</div>
          <div class="row-desc">{{ i18n('sec.vbsDesc') }}</div>
        </div>
        <n-switch v-model:value="vbs" :loading="vbsBusy" @update:value="onVbs" />
      </div>
      <div class="setting-row">
        <HardDrive :size="18" class="row-icon" />
        <div>
          <div class="row-label">{{ i18n('sec.memIntegrity') }}</div>
          <div class="row-desc">{{ i18n('sec.memIntegrityDesc') }}</div>
        </div>
        <n-switch v-model:value="memIntegrity" :loading="memIntegrityBusy" @update:value="onMemIntegrity" />
      </div>
    </div>

    <div class="setting-card">
      <div class="setting-card-header setting-card-header--flat">
        <span class="header-title">{{ i18n('sec.privacy') }}</span>
      </div>
      <div class="setting-row">
        <AlertTriangle :size="18" class="row-icon" />
        <div>
          <div class="row-label">{{ i18n('sec.telemetryTitle') }}</div>
          <div class="row-desc">{{ i18n('sec.telemetryDesc') }}</div>
        </div>
        <n-switch :value="telemetryBlocked" :disabled="showModal" @update:value="onTelemetryToggle" />
      </div>
    </div>

    <!-- NetWork -->
    <div class="setting-card">
      <div class="setting-card-header setting-card-header--flat">
        <span class="header-title">{{ i18n('sec.network') }}</span>
      </div>
      <!-- DNS Settings -->
      <div class="setting-row">
        <AlertTriangle :size="18" class="row-icon" />
        <div>
          <div class="row-label">{{ i18n('sec.dns') }}</div>
          <div class="row-desc">{{ i18n('sec.dnsDesc') }}</div>
        </div>
        <n-select
          v-model:value="selectedDns"
          :options="dnsOptions"
          size="small"
          style="width: 200px"
          @update:value="onDnsChange"
        />
      </div>
      <div class="setting-row">
        <div style="flex:2">
          <div class="row-label" style="display:flex; justify-content:space-between">
            <span>{{ i18n('sec.dnsAdapter') }}</span>
            <span>{{ currentDnsInfo.adapter || i18n('sec.dnsNotDetected') }}</span>
          </div>
          <div class="row-label" style="display:flex; justify-content:space-between">
            <span>{{ i18n('sec.dnsPrimary') }}</span>
            <span>{{ currentDnsInfo.primaryDns || i18n('sec.dnsNotSet') }}</span>
          </div>
          <div class="row-label" style="display:flex; justify-content:space-between">
            <span>{{ i18n('sec.dnsSecondary') }}</span>
            <span>{{ currentDnsInfo.secondaryDns || i18n('sec.dnsNotSet') }}</span>
          </div>
        </div>
      </div>
    </div>

    <div class="setting-card">
      <div class="setting-card-header setting-card-header--flat">
        <span class="header-title">{{ i18n('sec.ifeoTitle') }}</span>
      </div>
      <div class="setting-row" style="cursor:pointer; align-items:center" @click="expandIfeo = !expandIfeo">
        <Bug :size="18" class="row-icon" />
        <span class="header-title">{{ i18n('sec.ifeoTitle') }}</span>
        <span style="flex:1"></span>
        <span style="font-size:12px; color:var(--text2); margin-right:8px">
          {{ ifeoEntries.length }}
        </span>
        <span style="font-size:12px; color:var(--text2); transition:transform 0.2s; display:inline-flex; align-items:center" :style="{ transform: expandIfeo ? 'rotate(90deg)' : 'rotate(0deg)' }">▶</span>
      </div>
      <template v-if="expandIfeo">
        <div class="setting-row">
          <div style="flex:1">
            <div class="row-label">{{ i18n('sec.ifeoDesc') }}</div>
          </div>
          <n-button size="small" :loading="ifeoLoading" @click.stop="loadIfeoEntries">
            <template #icon><RefreshCw :size="14" /></template>
          </n-button>
          <n-button size="small" @click.stop="newExeName = ''; newDebugger = ''; loadRunningProcesses(); showAddModal = true">
            <template #icon><Plus :size="14" /></template>
            {{ i18n('sec.ifeoAdd') }}
          </n-button>
        </div>
        <div class="setting-row" style="padding: 8px 12px">
          <n-input
            v-model:value="ifeoSearch"
            :placeholder="i18n('sec.ifeoSearch')"
            clearable
            size="small"
          >
            <template #prefix><Search :size="14" /></template>
          </n-input>
        </div>
        <div style="max-height: 40vh; overflow-y: auto">
          <div v-for="entry in filteredIfeo" :key="entry.name" class="setting-row">
            <div style="flex:1; min-width:0">
              <div class="row-label" style="font-weight:500">{{ entry.name }}</div>
              <div class="row-desc" style="font-family:monospace; font-size:11px">
                {{ entry.debugger || 'GlobalFlag: ' + entry.globalFlag }}
              </div>
            </div>
            <div style="flex-shrink:0; text-align:right">
              <n-button size="small" type="error" @click.stop="confirmDeleteIfeo(entry.name)">
                <template #icon><Trash2 :size="14" /></template>
              </n-button>
            </div>
          </div>
        </div>
        <div v-if="!ifeoLoading && filteredIfeo.length === 0" class="setting-row" style="justify-content:center; color:var(--text2)">
          {{ i18n('sec.ifeoNoResults') }}
        </div>
      </template>
    </div>

    <n-modal v-model:show="showAddModal" preset="card" :title="i18n('sec.ifeoAddTitle')"
      style="width: 480px" :bordered="false">
      <div style="margin-bottom:12px; padding:10px 12px; background:rgba(255,200,0,0.08); border:1px solid rgba(255,200,0,0.3); border-radius:6px; color:#e6a817; font-size:13px; line-height:1.5">
        ⚠️ {{ i18n('sec.ifeoWarning') }}
      </div>
      <div style="display:flex; flex-direction:column; gap:12px">
        <div>
          <div style="font-size:12px; color:var(--text2); margin-bottom:4px">{{ i18n('sec.ifeoExeName') }}</div>
          <n-select
            v-model:value="newExeName"
            :options="runningProcesses.map(p => ({ label: p, value: p }))"
            :placeholder="i18n('sec.ifeoExePlaceholder')"
            filterable
            tag
            clearable
            size="small"
          />
        </div>
        <div>
          <div style="font-size:12px; color:var(--text2); margin-bottom:4px">{{ i18n('sec.ifeoDebugger') }}</div>
          <n-input v-model:value="newDebugger" :placeholder="i18n('sec.ifeoDebuggerPlaceholder')" size="small" />
        </div>
      </div>
      <template #footer>
        <div style="display:flex; justify-content:flex-end; gap:8px">
          <n-button @click="showAddModal = false">{{ i18n('sec.ifeoCancel') }}</n-button>
          <n-button @click="addIfeoEntry">{{ i18n('sec.ifeoConfirm') }}</n-button>
        </div>
      </template>
    </n-modal>

    <n-modal v-model:show="showDeleteModal" preset="card" :title="i18n('sec.ifeoDeleteConfirm')"
      style="width: 440px" :bordered="false">
      <div style="margin-bottom:12px; color:var(--text2)">
        {{ i18n('sec.ifeoDeleteDesc') }}
      </div>
      <div v-if="deleteTarget" style="padding:8px 12px; background:var(--section-bg); border-radius:6px">
        <div style="font-weight:500; font-family:monospace">{{ deleteTarget }}</div>
      </div>
      <template #footer>
        <div style="display:flex; justify-content:flex-end; gap:8px">
          <n-button @click="showDeleteModal = false">{{ i18n('sec.ifeoCancel') }}</n-button>
          <n-button @click="doDeleteIfeo">{{ i18n('sec.ifeoConfirm') }}</n-button>
        </div>
      </template>
    </n-modal>

    <n-modal
      v-model:show="showDefenderWarn"
      preset="card"
      :title="i18n('sec.defenderWarnTitle')"
      style="width: 440px; max-width: 90vw;"
      :mask-closable="false"
      closable
      @close="showDefenderWarn = false"
    >
      <p style="color: var(--text2); line-height: 1.8;">{{ i18n('sec.defenderWarnDesc') }}</p>
      <template #footer>
        <div style="display: flex; justify-content: flex-end; gap: 8px;">
          <n-button @click="showDefenderWarn = false">{{ i18n('sec.defenderOkBtn') }}</n-button>
        </div>
      </template>
    </n-modal>
  </div>
</template>
