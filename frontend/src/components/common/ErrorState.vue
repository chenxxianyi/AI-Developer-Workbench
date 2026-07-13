<template>
  <div class="flex flex-col items-center justify-center py-16 px-4 text-center">
    <div class="w-16 h-16 mb-4 rounded-full bg-red-50 flex items-center justify-center">
      <svg class="w-8 h-8 text-[var(--color-danger)]" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 9v2m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
    </div>
    <h3 class="text-base font-semibold text-[var(--color-text-primary)] mb-1">{{ title || '加载失败' }}</h3>
    <p v-if="message" class="text-sm text-[var(--color-text-muted)] max-w-sm mb-4">{{ message }}</p>
    <AppButton v-if="retryLabel" variant="secondary" size="sm" @click="emit('retry')">{{ retryLabel }}</AppButton>
  </div>
</template>

<script setup lang="ts">
import AppButton from './AppButton.vue'
defineProps<{ title?: string; message?: string; retryLabel?: string }>()
const emit = defineEmits<{ retry: [] }>()
</script>
