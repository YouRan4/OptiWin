<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { NButton, NModal, NInput } from 'naive-ui'
import { Package, Search, RefreshCw } from 'lucide-vue-next'
import WaitModal from '../components/WaitModal.vue'
import { useNotify } from '../composables/useNotify'
import { ListAppxPackages, UninstallAppx } from '../../wailsjs/go/main/App'

const { t: i18n } = useI18n()
const notify = useNotify()

const showModal = ref(false)
const modalText = ref('')
const loading = ref(false)
const searchText = ref('')

interface AppxPackage {
  Name: string
  PackageFullName: string
  Version: string
  DisplayName: string
  Type: string
}

const packages = ref<AppxPackage[]>([])
const showUninstallModal = ref(false)
const uninstallTarget = ref<AppxPackage | null>(null)
const expandUserApps = ref(false)
const expandSystemApps = ref(false)

const filteredUserApps = computed(() => {
  const list = packages.value.filter(p => p.Type === 'user')
  if (!searchText.value) return list
  const q = searchText.value.toLowerCase()
  return list.filter(p =>
    p.Name.toLowerCase().includes(q) ||
    p.DisplayName.toLowerCase().includes(q) ||
    p.PackageFullName.toLowerCase().includes(q)
  )
})

const filteredSystemApps = computed(() => {
  const list = packages.value.filter(p => p.Type === 'system')
  if (!searchText.value) return list
  const q = searchText.value.toLowerCase()
  return list.filter(p =>
    p.Name.toLowerCase().includes(q) ||
    p.DisplayName.toLowerCase().includes(q) ||
    p.PackageFullName.toLowerCase().includes(q)
  )
})

onMounted(async () => {
  await loadPackages()
})

async function loadPackages() {
  loading.value = true
  try {
    const result = await ListAppxPackages()
    packages.value = JSON.parse(result)
  } catch {
    packages.value = []
  } finally {
    loading.value = false
  }
}

function confirmUninstall(pkg: AppxPackage) {
  uninstallTarget.value = pkg
  showUninstallModal.value = true
}

async function doUninstall() {
  if (!uninstallTarget.value) return
  showUninstallModal.value = false

  modalText.value = i18n('appx.uninstalling')
  showModal.value = true

  try {
    const ok = await UninstallAppx(uninstallTarget.value.PackageFullName)
    if (ok) {
      packages.value = packages.value.filter(p => p.PackageFullName !== uninstallTarget.value!.PackageFullName)
      notify.create({ title: i18n('appx.uninstallSuccess'), description: uninstallTarget.value.DisplayName, duration: 4000 })
    } else {
      notify.create({ title: i18n('appx.uninstallFailed'), description: uninstallTarget.value.DisplayName, duration: 6000 })
    }
  } catch (e: any) {
    notify.create({ title: i18n('appx.uninstallFailed'), description: e?.message || 'Error', duration: 6000 })
  } finally {
    showModal.value = false
    uninstallTarget.value = null
  }
}
</script>

<template>
  <div class="page">
    <WaitModal :show="showModal" :text="modalText" />

    <div class="setting-card">
      <div class="setting-card-header setting-card-header--flat">
        <span class="header-title">{{ i18n('appx.title') }}</span>
      </div>
      <div class="setting-row">
        <Package :size="18" class="row-icon" />
        <div style="flex:1">
          <div class="row-label">{{ i18n('appx.desc') }}</div>
        </div>
        <n-button size="small" :loading="loading" @click="loadPackages">
          <template #icon><RefreshCw :size="14" /></template>
          {{ i18n('appx.refresh') }}
        </n-button>
      </div>
    </div>

    <div class="setting-card">
      <div class="setting-row" style="padding: 8px 12px">
        <n-input
          v-model:value="searchText"
          :placeholder="i18n('appx.search')"
          clearable
          size="small"
        >
          <template #prefix><Search :size="14" /></template>
        </n-input>
      </div>
    </div>

    <div class="setting-card">
      <div class="setting-row" style="cursor:pointer" @click="expandUserApps = !expandUserApps">
        <span class="header-title">{{ i18n('appx.userApps') }}</span>
        <span style="flex:1"></span>
        <span style="font-size:12px; color:var(--text2); margin-right:8px">
          {{ filteredUserApps.length }}
        </span>
        <span style="font-size:12px; color:var(--text2); transition:transform 0.2s" :style="{ transform: expandUserApps ? 'rotate(90deg)' : 'rotate(0deg)' }">▶</span>
      </div>
      <div v-if="expandUserApps" style="max-height: 40vh; overflow-y: auto">
        <div v-for="pkg in filteredUserApps" :key="pkg.PackageFullName" class="setting-row">
          <div style="flex:1; min-width:0">
            <div class="row-label" style="font-weight:500">{{ pkg.DisplayName }}</div>
            <div class="row-desc" style="font-family:monospace; font-size:11px">{{ pkg.Name }}</div>
          </div>
          <div style="flex-shrink:0; width:80px; text-align:right">
            <div style="font-size:12px; color:var(--text2)">{{ pkg.Version }}</div>
          </div>
          <div style="flex-shrink:0">
            <n-button size="small" @click.stop="confirmUninstall(pkg)">
              {{ i18n('appx.uninstall') }}
            </n-button>
          </div>
        </div>
      </div>
      <div v-if="expandUserApps && !loading && filteredUserApps.length === 0" class="setting-row" style="justify-content:center; color:var(--text2)">
        {{ i18n('appx.noResults') }}
      </div>
    </div>

    <div class="setting-card">
      <div class="setting-row" style="cursor:pointer" @click="expandSystemApps = !expandSystemApps">
        <span class="header-title">{{ i18n('appx.systemApps') }}</span>
        <span style="flex:1"></span>
        <span style="font-size:12px; color:var(--text2); margin-right:8px">
          {{ filteredSystemApps.length }}
        </span>
        <span style="font-size:12px; color:var(--text2); transition:transform 0.2s" :style="{ transform: expandSystemApps ? 'rotate(90deg)' : 'rotate(0deg)' }">▶</span>
      </div>
      <div v-if="expandSystemApps" style="max-height: 40vh; overflow-y: auto">
        <div v-for="pkg in filteredSystemApps" :key="pkg.PackageFullName" class="setting-row">
          <div style="flex:1; min-width:0">
            <div class="row-label" style="font-weight:500">{{ pkg.DisplayName }}</div>
            <div class="row-desc" style="font-family:monospace; font-size:11px">{{ pkg.Name }}</div>
          </div>
          <div style="flex-shrink:0; width:80px; text-align:right">
            <div style="font-size:12px; color:var(--text2)">{{ pkg.Version }}</div>
          </div>
          <div style="flex-shrink:0">
            <n-button size="small" @click.stop="confirmUninstall(pkg)">
              {{ i18n('appx.uninstall') }}
            </n-button>
          </div>
        </div>
      </div>
      <div v-if="expandSystemApps && !loading && filteredSystemApps.length === 0" class="setting-row" style="justify-content:center; color:var(--text2)">
        {{ i18n('appx.noResults') }}
      </div>
    </div>

    <n-modal v-model:show="showUninstallModal" preset="card" :title="i18n('appx.uninstallConfirm')"
      style="width: 440px" :bordered="false">
      <div style="margin-bottom:12px; color:var(--text2)">
        {{ i18n('appx.uninstallDesc') }}
      </div>
      <div v-if="uninstallTarget" style="padding:8px 12px; background:var(--section-bg); border-radius:6px">
        <div style="font-weight:500">{{ uninstallTarget.DisplayName }}</div>
        <div style="font-size:12px; color:var(--text2); font-family:monospace">{{ uninstallTarget.Name }}</div>
      </div>
      <template #footer>
        <div style="display:flex; justify-content:flex-end; gap:8px">
          <n-button @click="showUninstallModal = false">{{ i18n('appx.cancel') }}</n-button>
          <n-button @click="doUninstall">{{ i18n('appx.confirm') }}</n-button>
        </div>
      </template>
    </n-modal>
  </div>
</template>
