<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { NSelect, NSwitch } from 'naive-ui'
import {
  GetNotificationStatus, SetNotificationMode,
  GetLegacyBalloonStatus, SetLegacyBalloon,
  GetEdgeSwipeStatus, SetEdgeSwipe,
  GetNewContextMenuStatus, SetNewContextMenu,
  GetExplorerHomeStatus, SetExplorerHome,
  GetExplorerGalleryStatus, SetExplorerGallery,
  GetRemoveShortcutArrowStatus, SetRemoveShortcutArrow,
  GetRemoveShortcutTextStatus, SetRemoveShortcutText,
  GetRemoveShieldStatus, SetRemoveShield,
} from '../../wailsjs/go/main/App'

const { t: i18n } = useI18n()

const notif = ref('开启')
const legacyBalloons = ref(false)
const edgeSwipe = ref(true)
const newMenu = ref(true)
const explorerHome = ref(true)
const explorerGallery = ref(true)
const removeArrow = ref(false)
const removeText = ref(false)
const removeShield = ref(false)

onMounted(async () => {
  notif.value = await GetNotificationStatus()
  legacyBalloons.value = await GetLegacyBalloonStatus()
  edgeSwipe.value = await GetEdgeSwipeStatus()
  newMenu.value = await GetNewContextMenuStatus()
  explorerHome.value = await GetExplorerHomeStatus()
  explorerGallery.value = await GetExplorerGalleryStatus()
  removeArrow.value = await GetRemoveShortcutArrowStatus()
  removeText.value = await GetRemoveShortcutTextStatus()
  removeShield.value = await GetRemoveShieldStatus()
})

const notifOptions = computed(() => [
  { label: i18n('pers.notifOn'), value: '0' },
  { label: i18n('pers.notifOff'), value: '1' },
  { label: i18n('pers.notifAllOff'), value: '2' },
])

async function onNotifChange(v: string) {
  await SetNotificationMode(v)
  notif.value = await GetNotificationStatus()
}

async function onLegacyBalloon(v: boolean) {
  await SetLegacyBalloon(v)
  legacyBalloons.value = await GetLegacyBalloonStatus()
}

async function onEdgeSwipe(v: boolean) {
  await SetEdgeSwipe(v)
  edgeSwipe.value = await GetEdgeSwipeStatus()
}

async function onNewMenu(v: boolean) {
  await SetNewContextMenu(v)
  newMenu.value = await GetNewContextMenuStatus()
}

async function onExplorerHome(v: boolean) {
  await SetExplorerHome(v)
  explorerHome.value = await GetExplorerHomeStatus()
}

async function onExplorerGallery(v: boolean) {
  await SetExplorerGallery(v)
  explorerGallery.value = await GetExplorerGalleryStatus()
}

async function onRemoveArrow(v: boolean) {
  await SetRemoveShortcutArrow(v)
  removeArrow.value = await GetRemoveShortcutArrowStatus()
}

async function onRemoveText(v: boolean) {
  await SetRemoveShortcutText(v)
  removeText.value = await GetRemoveShortcutTextStatus()
}

async function onRemoveShield(v: boolean) {
  await SetRemoveShield(v)
  removeShield.value = await GetRemoveShieldStatus()
}
</script>

<template>
  <div class="page">
    <h2>{{ i18n('pers.title') }}</h2>
    <div class="setting-card">
      <div class="setting-card-header"><span class="header-title">{{ i18n('pers.desktop') }}</span><span class="header-desc">{{ i18n('pers.desktopDesc') }}</span></div>
      <div class="setting-row">
        <div><div class="row-label">{{ i18n('pers.notification') }}</div><div class="row-desc">{{ i18n('pers.notificationDesc') }}</div></div>
        <n-select v-model:value="notif" :options="notifOptions" style="width:140px" @update:value="onNotifChange" />
      </div>
      <div class="setting-row"><div><div class="row-label">{{ i18n('pers.legacyBalloon') }}</div><div class="row-desc">{{ i18n('pers.legacyBalloonDesc') }}</div></div><n-switch v-model:value="legacyBalloons" @update:value="onLegacyBalloon" /></div>
      <div class="setting-row"><div><div class="row-label">{{ i18n('pers.edgeSwipe') }}</div><div class="row-desc">{{ i18n('pers.edgeSwipeDesc') }}</div></div><n-switch v-model:value="edgeSwipe" @update:value="onEdgeSwipe" /></div>
      <div class="setting-row"><div><div class="row-label">{{ i18n('pers.contextMenu') }}</div><div class="row-desc">{{ i18n('pers.contextMenuDesc') }}</div></div><n-switch v-model:value="newMenu" @update:value="onNewMenu" /></div>
      <div class="setting-row"><div><div class="row-label">{{ i18n('pers.explorerHome') }}</div><div class="row-desc">{{ i18n('pers.explorerHomeDesc') }}</div></div><n-switch v-model:value="explorerHome" @update:value="onExplorerHome" /></div>
      <div class="setting-row"><div><div class="row-label">{{ i18n('pers.explorerGallery') }}</div><div class="row-desc">{{ i18n('pers.explorerGalleryDesc') }}</div></div><n-switch v-model:value="explorerGallery" @update:value="onExplorerGallery" /></div>
    </div>
    <div class="setting-card">
      <div class="setting-card-header"><span class="header-title">{{ i18n('pers.appearance') }}</span><span class="header-desc">{{ i18n('pers.appearanceDesc') }}</span></div>
      <div class="setting-row"><div><div class="row-label">{{ i18n('pers.removeArrow') }}</div><div class="row-desc">{{ i18n('pers.removeArrowDesc') }}</div></div><n-switch v-model:value="removeArrow" @update:value="onRemoveArrow" /></div>
      <div class="setting-row"><div><div class="row-label">{{ i18n('pers.removeText') }}</div><div class="row-desc">{{ i18n('pers.removeTextDesc') }}</div></div><n-switch v-model:value="removeText" @update:value="onRemoveText" /></div>
      <div class="setting-row"><div><div class="row-label">{{ i18n('pers.removeShield') }}</div><div class="row-desc">{{ i18n('pers.removeShieldDesc') }}</div></div><n-switch v-model:value="removeShield" @update:value="onRemoveShield" /></div>
    </div>
  </div>
</template>
