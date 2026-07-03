<script setup lang="ts">
/**
 * ConfirmDialog — accessible confirmation modal.
 */
import { AlertTriangle } from '@lucide/vue'

defineProps<{
  open: boolean
  title: string
  message: string
  confirmLabel?: string
  cancelLabel?: string
  danger?: boolean
}>()

const emit = defineEmits<{
  confirm: []
  cancel: []
}>()
</script>

<template>
  <Teleport to="body">
    <div
      v-if="open"
      class="fixed inset-0 z-50 flex items-center justify-center p-4"
      role="dialog"
      aria-modal="true"
      :aria-label="title"
    >
      <!-- Backdrop -->
      <div
        class="absolute inset-0 bg-black/40 backdrop-blur-sm"
        @click="emit('cancel')"
      />

      <!-- Dialog -->
      <div class="relative z-10 w-full max-w-md rounded-lg border border-border bg-surface p-6 shadow-lg">
        <div class="flex items-start gap-3 mb-4">
          <AlertTriangle
            :size="22"
            :class="danger ? 'text-danger' : 'text-warning'"
            class="flex-shrink-0 mt-0.5"
          />
          <div>
            <h3 class="text-lg font-semibold text-text-primary">{{ title }}</h3>
            <p class="mt-1 text-sm text-text-secondary">{{ message }}</p>
          </div>
        </div>

        <div class="flex justify-end gap-3 mt-6">
          <button
            class="rounded-md border border-border px-4 py-2 text-sm font-medium text-text-primary hover:bg-surface-hover transition-smooth"
            @click="emit('cancel')"
          >
            {{ cancelLabel || '取消' }}
          </button>
          <button
            :class="[
              'rounded-md px-4 py-2 text-sm font-semibold transition-smooth',
              danger
                ? 'bg-danger text-white hover:bg-danger/90'
                : 'bg-accent text-white hover:bg-accent/90',
            ]"
            @click="emit('confirm')"
          >
            {{ confirmLabel || '确认' }}
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>
