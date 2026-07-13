<template>
  <Teleport to="body">
    <TransitionGroup
      tag="div"
      name="toast"
      class="fixed bottom-4 right-4 z-[100] flex flex-col gap-2 max-w-sm"
    >
      <div
        v-for="t in toasts"
        :key="t.id"
        :class="toastClass(t.type)"
        class="px-4 py-3 rounded-[var(--radius-card)] shadow-[var(--shadow-md)] text-sm flex items-center gap-2"
      >
        <span class="flex-1">{{ t.message }}</span>
        <button class="text-current opacity-60 hover:opacity-100 ml-2" @click="remove(t.id)">✕</button>
      </div>
    </TransitionGroup>
  </Teleport>
</template>

<script setup lang="ts">
import { ref } from 'vue'

interface Toast { id: number; message: string; type: 'success' | 'error' | 'warning' | 'info' }
const toasts = ref<Toast[]>([])
let nextId = 0

function show(message: string, type: Toast['type'] = 'info', duration = 4000) {
  const id = nextId++
  toasts.value.push({ id, message, type })
  if (duration > 0) setTimeout(() => remove(id), duration)
}
function remove(id: number) { toasts.value = toasts.value.filter(t => t.id !== id) }

function toastClass(type: string) {
  const map: Record<string, string> = {
    success: 'bg-green-50 text-[var(--color-success)] border border-green-200',
    error: 'bg-red-50 text-[var(--color-danger)] border border-red-200',
    warning: 'bg-amber-50 text-[var(--color-warning)] border border-amber-200',
    info: 'bg-[var(--color-surface)] text-[var(--color-text-primary)] border border-[var(--color-border)]',
  }
  return map[type] || map.info
}

defineExpose({ show, remove })
</script>

<style scoped>
.toast-enter-active { transition: all 0.3s ease; }
.toast-leave-active { transition: all 0.2s ease; }
.toast-enter-from { opacity: 0; transform: translateX(50px); }
.toast-leave-to { opacity: 0; transform: translateX(50px); }
</style>
