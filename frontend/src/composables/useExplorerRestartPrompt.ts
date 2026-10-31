import { ref } from 'vue'

const pending = ref(false)
const error = ref('')

export const explorerRestartPrompt = {
  pending,
  error,
  show: () => {
    error.value = ''
    pending.value = true
  },
  fail: (msg: string) => {
    error.value = msg
    pending.value = true
  },
  dismiss: () => {
    pending.value = false
    error.value = ''
  },
}
