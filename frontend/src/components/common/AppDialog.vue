<template>
  <Teleport to="body">
    <Transition name="dialog">
      <div v-if="open" class="fixed inset-0 z-50 flex items-center justify-center p-4" @click.self="closeOnOverlay && emit('close')">
        <div class="fixed inset-0 bg-black/40" />
        <div class="relative bg-[var(--color-surface)] rounded-[var(--radius-panel)] shadow-[var(--shadow-lg)] max-w-lg w-full max-h-[90vh] overflow-auto" role="dialog" aria-modal="true">
          <div class="flex items-center justify-between px-5 py-4 border-b border-[var(--color-border)]">
            <h2 class="text-lg font-semibold text-[var(--color-text-primary)]">{{ title }}</h2>
            <button class="p-1 rounded-md hover:bg-[var(--color-surface-muted)] text-[var(--color-text-muted)]" @click="emit('close')" aria-label="关闭">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/></svg>
            </button>
          </div>
          <div class="p-5"><slot /></div>
          <div v-if="$slots.footer" class="px-5 py-3 border-t border-[var(--color-border)] flex justify-end gap-3"><slot name="footer" /></div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
withDefaults(defineProps<{ open: boolean; title: string; closeOnOverlay?: boolean }>(), { closeOnOverlay: true })
const emit = defineEmits<{ close: [] }>()
</script>

<style scoped>
.dialog-enter-active, .dialog-leave-active { transition: opacity 0.2s ease; }
.dialog-enter-from, .dialog-leave-to { opacity: 0; }
</style>
