<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { NSelect, NSwitch } from 'naive-ui'
import { Bell, MessageCircle, Hand, List, Home, Image, ArrowUpRight, Type, ShieldOff, Activity, LayoutGrid } from 'lucide-vue-next'
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
  GetOldTaskManagerStatus, SetOldTaskManager,
  GetWidgetsStatus, DisableWidgets, EnableWidgets,
} from '../../wailsjs/go/main/App'

import { useBackendToggle } from '../composables/useBackendToggle'
import { useNotify } from '../composables/useNotify'
import { explorerRestartPrompt } from '../composables/useExplorerRestartPrompt'

const { t: i18n } = useI18n()
const notify = useNotify()

const notif = ref('0')

const { value: legacyBalloons, refresh: refreshLegacyBalloons, onUpdate: onLegacyBalloon, busy: legacyBalloonsBusy } = useBackendToggle(
  GetLegacyBalloonStatus, SetLegacyBalloon, () => i18n('pers.legacyBalloon'), explorerRestartPrompt.show,
)
const { value: edgeSwipe, refresh: refreshEdgeSwipe, onUpdate: onEdgeSwipe, busy: edgeSwipeBusy } = useBackendToggle(
  GetEdgeSwipeStatus, SetEdgeSwipe, () => i18n('pers.edgeSwipe'), explorerRestartPrompt.show,
)
const { value: newMenu, refresh: refreshNewMenu, onUpdate: onNewMenu, busy: newMenuBusy } = useBackendToggle(
  GetNewContextMenuStatus, SetNewContextMenu, () => i18n('pers.contextMenu'), explorerRestartPrompt.show,
)
const { value: explorerHome, refresh: refreshExplorerHome, onUpdate: onExplorerHome, busy: explorerHomeBusy } = useBackendToggle(
  GetExplorerHomeStatus, SetExplorerHome, () => i18n('pers.explorerHome'), explorerRestartPrompt.show,
)
const { value: explorerGallery, refresh: refreshExplorerGallery, onUpdate: onExplorerGallery, busy: explorerGalleryBusy } = useBackendToggle(
  GetExplorerGalleryStatus, SetExplorerGallery, () => i18n('pers.explorerGallery'), explorerRestartPrompt.show,
)
const { value: removeArrow, refresh: refreshRemoveArrow, onUpdate: onRemoveArrow, busy: removeArrowBusy } = useBackendToggle(
  GetRemoveShortcutArrowStatus, SetRemoveShortcutArrow, () => i18n('pers.removeArrow'), explorerRestartPrompt.show,
)
const { value: removeText, refresh: refreshRemoveText, onUpdate: onRemoveText, busy: removeTextBusy } = useBackendToggle(
  GetRemoveShortcutTextStatus, SetRemoveShortcutText, () => i18n('pers.removeText'), explorerRestartPrompt.show,
)
const { value: removeShield, refresh: refreshRemoveShield, onUpdate: onRemoveShield, busy: removeShieldBusy } = useBackendToggle(
  GetRemoveShieldStatus, SetRemoveShield, () => i18n('pers.removeShield'), explorerRestartPrompt.show,
)
const { value: win11TaskManager, refresh: refreshWin11TaskManager, onUpdate: onWin11TaskManager, busy: win11TaskManagerBusy } = useBackendToggle(
  GetOldTaskManagerStatus, SetOldTaskManager, () => i18n('pers.win11TaskManager'), explorerRestartPrompt.show,
)
const { value: widgets, refresh: refreshWidgets, onUpdate: onWidgets, busy: widgetsBusy } = useBackendToggle(
  GetWidgetsStatus, v => (v ? EnableWidgets() : DisableWidgets()), () => i18n('pers.widgets'),
)

onMounted(async () => {
  notif.value = await GetNotificationStatus()
  await Promise.all([
    refreshLegacyBalloons(),
    refreshEdgeSwipe(),
    refreshNewMenu(),
    refreshExplorerHome(),
    refreshExplorerGallery(),
    refreshRemoveArrow(),
    refreshRemoveText(),
    refreshRemoveShield(),
    refreshWin11TaskManager(),
    refreshWidgets(),
  ])
})

const notifOptions = computed(() => [
  { label: i18n('pers.notifOn'), value: '0' },
  { label: i18n('pers.notifOff'), value: '1' },
  { label: i18n('pers.notifAllOff'), value: '2' },
])

async function onNotifChange(v: string) {
  const prev = notif.value
  const ok = await SetNotificationMode(v)
  notif.value = await GetNotificationStatus()
  if (!ok) {
    notify.create({
      title: i18n('common.updateFailed'),
      description: `${i18n('pers.notification')}：${i18n('common.updateFailedDesc')}`,
      type: 'error',
      duration: 6000,
    })
    return
  }
  if (notif.value !== prev) {
    explorerRestartPrompt.show()
  }
}
</script>

<template>
  <div class="page">
    <div class="setting-card">
      <div class="setting-card-header setting-card-header--flat"><span class="header-title">{{ i18n('pers.desktop') }}</span></div>
      <div class="setting-row">
        <Bell :size="18" class="row-icon" />
        <div><div class="row-label">{{ i18n('pers.notification') }}</div><div class="row-desc">{{ i18n('pers.notificationDesc') }}</div></div>
        <n-select v-model:value="notif" :options="notifOptions" style="width:140px" @update:value="onNotifChange" />
      </div>
      <div class="setting-row">
        <MessageCircle :size="18" class="row-icon" />
        <div><div class="row-label">{{ i18n('pers.legacyBalloon') }}</div><div class="row-desc">{{ i18n('pers.legacyBalloonDesc') }}</div></div><n-switch v-model:value="legacyBalloons" :loading="legacyBalloonsBusy" @update:value="onLegacyBalloon" />
      </div>
      <div class="setting-row">
        <Hand :size="18" class="row-icon" />
        <div><div class="row-label">{{ i18n('pers.edgeSwipe') }}</div><div class="row-desc">{{ i18n('pers.edgeSwipeDesc') }}</div></div><n-switch v-model:value="edgeSwipe" :loading="edgeSwipeBusy" @update:value="onEdgeSwipe" />
      </div>
      <div class="setting-row">
        <List :size="18" class="row-icon" />
        <div><div class="row-label">{{ i18n('pers.contextMenu') }}</div><div class="row-desc">{{ i18n('pers.contextMenuDesc') }}</div></div><n-switch v-model:value="newMenu" :loading="newMenuBusy" @update:value="onNewMenu" />
      </div>
      <div class="setting-row">
        <Activity :size="18" class="row-icon" />
        <div><div class="row-label">{{ i18n('pers.win11TaskManager') }}</div><div class="row-desc">{{ i18n('pers.win11TaskManagerDesc') }}</div></div><n-switch v-model:value="win11TaskManager" :loading="win11TaskManagerBusy" @update:value="onWin11TaskManager" />
      </div>
      <div class="setting-row">
        <LayoutGrid :size="18" class="row-icon" />
        <div><div class="row-label">{{ i18n('pers.widgets') }}</div><div class="row-desc">{{ i18n('pers.widgetsDesc') }}</div></div><n-switch v-model:value="widgets" :loading="widgetsBusy" @update:value="onWidgets" />
      </div>
      <div class="setting-row">
        <Home :size="18" class="row-icon" />
        <div><div class="row-label">{{ i18n('pers.explorerHome') }}</div><div class="row-desc">{{ i18n('pers.explorerHomeDesc') }}</div></div><n-switch v-model:value="explorerHome" :loading="explorerHomeBusy" @update:value="onExplorerHome" />
      </div>
      <div class="setting-row">
        <Image :size="18" class="row-icon" />
        <div><div class="row-label">{{ i18n('pers.explorerGallery') }}</div><div class="row-desc">{{ i18n('pers.explorerGalleryDesc') }}</div></div><n-switch v-model:value="explorerGallery" :loading="explorerGalleryBusy" @update:value="onExplorerGallery" />
      </div>
    </div>
    <div class="setting-card">
      <div class="setting-card-header setting-card-header--flat"><span class="header-title">{{ i18n('pers.appearance') }}</span></div>
      <div class="setting-row">
        <ArrowUpRight :size="18" class="row-icon" />
        <div><div class="row-label">{{ i18n('pers.removeArrow') }}</div><div class="row-desc">{{ i18n('pers.removeArrowDesc') }}</div></div><n-switch v-model:value="removeArrow" :loading="removeArrowBusy" @update:value="onRemoveArrow" />
      </div>
      <div class="setting-row">
        <Type :size="18" class="row-icon" />
        <div><div class="row-label">{{ i18n('pers.removeText') }}</div><div class="row-desc">{{ i18n('pers.removeTextDesc') }}</div></div><n-switch v-model:value="removeText" :loading="removeTextBusy" @update:value="onRemoveText" />
      </div>
      <div class="setting-row">
        <ShieldOff :size="18" class="row-icon" />
        <div><div class="row-label">{{ i18n('pers.removeShield') }}</div><div class="row-desc">{{ i18n('pers.removeShieldDesc') }}</div></div><n-switch v-model:value="removeShield" :loading="removeShieldBusy" @update:value="onRemoveShield" />
      </div>
    </div>
  </div>
</template>
