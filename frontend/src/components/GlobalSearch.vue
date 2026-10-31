<script setup lang="ts">
import { ref, computed, nextTick, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { Search } from 'lucide-vue-next'
import { searchIndex, type SearchEntry } from '../searchIndex'

const { t: i18n } = useI18n()
const router = useRouter()
const route = useRoute()

const open = ref(false)
const query = ref('')
const activeIndex = ref(0)
const inputEl = ref<HTMLInputElement | null>(null)

function haystack(e: SearchEntry): string {
  return [i18n(e.label), e.desc ? i18n(e.desc) : '', i18n(e.group), i18n(e.page)]
    .join(' ')
    .toLowerCase()
}

const results = computed(() => {
  const q = query.value.trim().toLowerCase()
  if (!q) return searchIndex.slice(0, 8)
  return searchIndex.filter(e => haystack(e).includes(q)).slice(0, 8)
})

async function show() {
  open.value = true
  query.value = ''
  activeIndex.value = 0
  await nextTick()
  inputEl.value?.focus()
}

function hide() {
  open.value = false
}

function move(step: number) {
  const n = results.value.length
  if (!n) return
  activeIndex.value = (activeIndex.value + step + n) % n
}

async function select(e: SearchEntry) {
  hide()
  if (route.path !== e.path) await router.push(e.path)
  flashRow(i18n(e.label))
}

function flashRow(text: string) {
  const deadline = Date.now() + 2000
  const tick = () => {
    const hit = Array.from(document.querySelectorAll<HTMLElement>('.row-label'))
      .filter(el => (el.textContent || '').trim().startsWith(text))
      .find(el => !el.closest('.page-leave-active, .page-leave-to'))
    if (hit) {
      hit.scrollIntoView({ block: 'center', behavior: 'smooth' })
      hit.classList.add('search-hit')
      setTimeout(() => hit.classList.remove('search-hit'), 1800)
      return
    }
    if (Date.now() < deadline) requestAnimationFrame(tick)
  }
  requestAnimationFrame(tick)
}

function onKeydown(e: KeyboardEvent) {
  if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === 'k') {
    e.preventDefault()
    open.value ? hide() : show()
    return
  }
  if (!open.value) return
  if (e.key === 'Escape') { e.preventDefault(); hide() }
  else if (e.key === 'ArrowDown') { e.preventDefault(); move(1) }
  else if (e.key === 'ArrowUp') { e.preventDefault(); move(-1) }
  else if (e.key === 'Enter') {
    e.preventDefault()
    const hit = results.value[activeIndex.value]
    if (hit) select(hit)
  }
}

onMounted(() => window.addEventListener('keydown', onKeydown))
onUnmounted(() => window.removeEventListener('keydown', onKeydown))
</script>

<template>
  <button class="tb-btn" :title="i18n('common.searchPlaceholder')" @click="show">
    <Search :size="14" />
  </button>
  <Transition name="gs">
    <div v-if="open" class="gs-mask" @click.self="hide">
      <div class="gs-panel">
        <input
          ref="inputEl"
          v-model="query"
          class="gs-input"
          :placeholder="i18n('common.searchPlaceholder')"
          @input="activeIndex = 0"
        />
        <div v-if="results.length" class="gs-list">
          <button
            v-for="(r, i) in results"
            :key="r.path + r.label"
            class="gs-item"
            :class="{ 'gs-item--active': i === activeIndex }"
            @click="select(r)"
            @mouseenter="activeIndex = i"
          >
            <span class="gs-label">{{ i18n(r.label) }}</span>
            <span class="gs-meta">{{ i18n(r.page) }} · {{ i18n(r.group) }}</span>
          </button>
        </div>
        <div v-else class="gs-empty">{{ i18n('common.searchEmpty') }}</div>
      </div>
    </div>
  </Transition>
</template>

<style scoped>
.gs-mask {
  position: fixed; inset: 0; z-index: 9500;
  background: rgba(0, 0, 0, 0.5); backdrop-filter: blur(4px);
  display: flex; justify-content: center; align-items: flex-start;
  padding-top: 12vh;
}
.gs-panel {
  width: 560px; max-width: calc(100vw - 64px);
  background: var(--card-bg); border: 1px solid var(--border);
  border-radius: 16px; box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
  overflow: hidden;
}
.gs-input {
  width: 100%; box-sizing: border-box;
  padding: 14px 18px; font-size: 14px; font-family: inherit;
  color: var(--text); background: transparent; border: none; outline: none;
  border-bottom: 1px solid var(--border);
}
.gs-input::placeholder { color: var(--text2); }
.gs-list { max-height: 320px; overflow-y: auto; padding: 6px; }
.gs-item {
  width: 100%; display: flex; align-items: center; justify-content: space-between; gap: 16px;
  padding: 9px 12px; border: none; border-radius: 8px; cursor: pointer;
  background: transparent; color: var(--text); font-family: inherit; font-size: 14px;
  text-align: left;
}
.gs-item--active { background: var(--hover); }
.gs-label { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.gs-meta { flex-shrink: 0; font-size: 12px; color: var(--text2); }
.gs-empty { padding: 18px; font-size: 14px; color: var(--text2); text-align: center; }

.gs-enter-active, .gs-leave-active { transition: opacity 0.18s ease; }
.gs-enter-active .gs-panel, .gs-leave-active .gs-panel {
  transition: transform 0.18s cubic-bezier(0.2, 0.9, 0.3, 1), opacity 0.18s ease;
}
.gs-enter-from, .gs-leave-to { opacity: 0; }
.gs-enter-from .gs-panel, .gs-leave-to .gs-panel {
  transform: translateY(-10px) scale(0.97);
  opacity: 0;
}
@media (prefers-reduced-motion: reduce) {
  .gs-enter-active, .gs-leave-active,
  .gs-enter-active .gs-panel, .gs-leave-active .gs-panel { transition: none; }
  .gs-enter-from .gs-panel, .gs-leave-to .gs-panel { transform: none; }
}
</style>
