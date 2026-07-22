<script setup lang="ts">
import { NSpin } from 'naive-ui'
defineProps<{ show: boolean; text: string }>()
</script>

<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="show" class="wait-overlay">
        <div class="wait-modal">
          <n-spin size="large" />
          <div class="wait-text">{{ text }}</div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.wait-overlay {
  position: fixed; inset: 0; z-index: 9999;
  background: rgba(0,0,0,0.5);
  display: flex; align-items: center; justify-content: center;
  backdrop-filter: blur(4px);
}
.wait-modal {
  background: var(--card-bg);
  border: 1px solid var(--border);
  border-radius: 16px;
  padding: 48px 56px;
  text-align: center;
  box-shadow: 0 8px 32px rgba(0,0,0,0.3);
  min-width: 320px;
  display: flex; flex-direction: column; align-items: center; gap: 24px;
}
.wait-text { font-size: 16px; font-weight: 600; }

.modal-enter-active, .modal-leave-active { transition: opacity 0.25s ease; }
.modal-enter-from, .modal-leave-to { opacity: 0; }

.n-spin .n-spin-body { --n-color: var(--spinner-color); }
</style>
