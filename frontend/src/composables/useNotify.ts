import { useNotification } from 'naive-ui'
import type { NotificationOptions } from 'naive-ui'

export function useNotify() {
  const notification = useNotification()
  return {
    create: (options: NotificationOptions) => notification.create(options),
  }
}
