<script setup lang="ts">
import { ref, onMounted } from 'vue'
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

const notifOptions = [
  { label: '开启', value: '开启' },
  { label: '仅关闭通知', value: '仅关闭通知' },
  { label: '完全关闭', value: '完全关闭' },
]

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
    <h2>个性化</h2>
    <div class="setting-card">
      <div class="setting-card-header"><span class="header-title">桌面个性化</span><span class="header-desc">通知、菜单和资源管理器设置</span></div>
      <div class="setting-row">
        <div><div class="row-label">通知</div><div class="row-desc">控制系统通知显示方式</div></div>
        <n-select v-model:value="notif" :options="notifOptions" style="width:140px" @update:value="onNotifChange" />
      </div>
      <div class="setting-row"><div><div class="row-label">旧版气球通知</div><div class="row-desc">使用 Win7 风格通知</div></div><n-switch v-model:value="legacyBalloons" @update:value="onLegacyBalloon" /></div>
      <div class="setting-row"><div><div class="row-label">屏幕边缘滑动</div><div class="row-desc">触控屏边缘手势</div></div><n-switch v-model:value="edgeSwipe" @update:value="onEdgeSwipe" /></div>
      <div class="setting-row"><div><div class="row-label">新版上下文菜单</div><div class="row-desc">Win11 新版菜单与旧版菜单切换</div></div><n-switch v-model:value="newMenu" @update:value="onNewMenu" /></div>
      <div class="setting-row"><div><div class="row-label">Explorer 主页</div><div class="row-desc">在资源管理器中显示/隐藏主页入口</div></div><n-switch v-model:value="explorerHome" @update:value="onExplorerHome" /></div>
      <div class="setting-row"><div><div class="row-label">Explorer 图库</div><div class="row-desc">在资源管理器中显示/隐藏图库入口</div></div><n-switch v-model:value="explorerGallery" @update:value="onExplorerGallery" /></div>
    </div>
    <div class="setting-card">
      <div class="setting-card-header"><span class="header-title">外观调整</span><span class="header-desc">桌面图标和显示设置</span></div>
      <div class="setting-row"><div><div class="row-label">移除快捷方式小箭头</div><div class="row-desc">去除桌面快捷方式图标上的小箭头覆盖</div></div><n-switch v-model:value="removeArrow" @update:value="onRemoveArrow" /></div>
      <div class="setting-row"><div><div class="row-label">移除"快捷方式"文字</div><div class="row-desc">创建快捷方式时不添加"-快捷方式"后缀</div></div><n-switch v-model:value="removeText" @update:value="onRemoveText" /></div>
      <div class="setting-row"><div><div class="row-label">移除小盾牌</div><div class="row-desc">去除可执行文件图标上的 UAC 盾牌覆盖</div></div><n-switch v-model:value="removeShield" @update:value="onRemoveShield" /></div>
    </div>
  </div>
</template>
