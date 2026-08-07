<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { NSwitch, NButton, NModal, NInput, NSelect } from 'naive-ui'
import { Lock, Server, HardDrive, Search, RefreshCw, Trash2, Plus, Bug, AlertTriangle } from 'lucide-vue-next'
import {
  GetUacStatus, EnableUac, DisableUac,
  GetVbsStatus, EnableVbs, DisableVbs,
  GetMemoryIntegrityStatus, EnableMemoryIntegrity, DisableMemoryIntegrity,
  ListIfeoEntries, AddIfeoEntry, RemoveIfeoEntry, GetRunningProcesses, SetDns, GetCurrentDns,
} from '../../wailsjs/go/main/App'
import WaitModal from '../components/WaitModal.vue'
import { useNotify } from '../composables/useNotify'

const { t: i18n } = useI18n()
const notify = useNotify()

const showModal = ref(false)
const modalText = ref('')

const uac = ref(false)
const vbs = ref(false)
const memIntegrity = ref(false)

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
  'conhost.exe', 'fontdrvhost.exe', 'RuntimeBroker.exe',
  'sihost.exe', 'shellexperiencehost.exe', 'StartMenuExperienceHost.exe',
  'SearchUI.exe', 'SearchIndexer.exe', 'MsMpEng.exe', 'taskhostw.exe',
  'ntoskrnl.exe', 'hal.dll', 'bootvid.dll', 'win32k.sys',
  'csrss.exe', 'smss.exe', 'lsass.exe', 'svchost.exe',
]

const telemetryProcesses = [
  'AggregatorHost.exe', 'BCILauncher.exe', 'BGAUpsell.exe',
  'BingChatInstaller.exe', 'CompatTelRunner.exe', 'DeviceCensus.exe',
  'FeatureLoader.exe'
]

const telemetryBlocked = ref(false)
const telemetryDebuggerPath = '%windir%\\system32\\taskkill.exe'

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
        description: `已切换到 ${label}`,
        duration: 3000
      })
    } else {
      notify.create({ title: i18n('sec.notifyTitle'), description: 'DNS 切换失败', duration: 4000 })
    }
  } catch (e: any) {
    notify.create({ title: i18n('sec.notifyTitle'), description: e?.message || 'Error', duration: 4000 })
  }
}

function checkTelemetryStatus() {
  telemetryBlocked.value = telemetryProcesses.some(p =>
    ifeoEntries.value.some(e => e.name.toLowerCase() === p.toLowerCase())
  )
}

const filteredIfeo = computed(() => {
  if (!ifeoSearch.value) return ifeoEntries.value
  const q = ifeoSearch.value.toLowerCase()
  return ifeoEntries.value.filter(e =>
    e.name.toLowerCase().includes(q) ||
    e.debugger.toLowerCase().includes(q)
  )
})

onMounted(async () => {
  uac.value = await GetUacStatus()
  vbs.value = await GetVbsStatus()
  memIntegrity.value = await GetMemoryIntegrityStatus()
  await loadIfeoEntries()
  await loadCurrentDns()
})

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

// IFEO functions
async function loadIfeoEntries() {
  ifeoLoading.value = true
  try {
    const result = await ListIfeoEntries()
    ifeoEntries.value = JSON.parse(result)
    checkTelemetryStatus()
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
  if (blockedProcesses.some(p => exeNameLower.includes(p))) {
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
    for (const proc of telemetryProcesses) {
      if (v) {
        await AddIfeoEntry(proc, telemetryDebuggerPath)
      } else {
        await RemoveIfeoEntry(proc)
      }
    }
    await loadIfeoEntries()
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
        <span class="header-title">{{ i18n('sec.systemProtection') }}</span>
      </div>
      <div class="setting-row">
        <Lock :size="18" class="row-icon" />
        <div>
          <div class="row-label">{{ i18n('sec.uac') }}</div>
          <div class="row-desc">{{ i18n('sec.uacDesc') }}</div>
        </div>
        <n-switch v-model:value="uac" @update:value="onUac" />
      </div>
      <div class="setting-row">
        <Server :size="18" class="row-icon" />
        <div>
          <div class="row-label">{{ i18n('sec.vbs') }}</div>
          <div class="row-desc">{{ i18n('sec.vbsDesc') }}</div>
        </div>
        <n-switch v-model:value="vbs" @update:value="onVbs" />
      </div>
      <div class="setting-row">
        <HardDrive :size="18" class="row-icon" />
        <div>
          <div class="row-label">{{ i18n('sec.memIntegrity') }}</div>
          <div class="row-desc">{{ i18n('sec.memIntegrityDesc') }}</div>
        </div>
        <n-switch v-model:value="memIntegrity" @update:value="onMemIntegrity" />
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
            <span>网卡</span>
            <span>{{ currentDnsInfo.adapter || '未检测到' }}</span>
          </div>
          <div class="row-label" style="display:flex; justify-content:space-between">
            <span>主 DNS</span>
            <span>{{ currentDnsInfo.primaryDns || '未设置' }}</span>
          </div>
          <div class="row-label" style="display:flex; justify-content:space-between">
            <span>副 DNS</span>
            <span>{{ currentDnsInfo.secondaryDns || '未设置' }}</span>
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
  </div>
</template>
