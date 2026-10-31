import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { LastError } from '../../wailsjs/go/main/App'
import { useNotify } from './useNotify'

/** 统一「调后端开关 → 重读状态 → 失败提示」；apply 返回后端是否真正写入成功 */
export function useBackendToggle(
  getStatus: () => Promise<boolean>,
  apply: (next: boolean) => Promise<boolean>,
  label: () => string,
  after?: () => void,
) {
  const { t: i18n } = useI18n()
  const notify = useNotify()
  const value = ref(false)
  const busy = ref(false)

  async function refresh() {
    try {
      value.value = await getStatus()
    } catch (e) {
      console.error('[OptiWin] 读取开关状态失败', e)
    }
  }

  async function onUpdate(next: boolean) {
    if (busy.value) return
    busy.value = true
    let ok = false
    try {
      ok = await apply(next)
    } catch (e) {
      console.error('[OptiWin] 开关调用异常', e)
    }
    await refresh()
    busy.value = false
    if (!ok) {
      const reason = await LastError().catch(() => '')
      notify.create({
        title: i18n('common.updateFailed'),
        description: reason ? `${label()}：${reason}` : `${label()}：${i18n('common.updateFailedDesc')}`,
        type: 'error',
        duration: 6000,
      })
      return
    }
    after?.()
  }

  return { value, busy, refresh, onUpdate }
}
