<template>
  <div class="w-full bg-[var(--color-surface-muted)] rounded-full h-2 overflow-hidden">
    <div
      class="h-full rounded-full transition-all duration-500 ease-out"
      :class="colorClass"
      :style="{ width: `${Math.min(100, Math.max(0, modelValue))}%` }"
    />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
const props = withDefaults(defineProps<{ modelValue: number; variant?: 'primary' | 'success' | 'danger' }>(), { variant: 'primary' })
const colorClass = computed(() => {
  const map: Record<string, string> = { primary: 'bg-[var(--color-accent)]', success: 'bg-[var(--color-success)]', danger: 'bg-[var(--color-danger)]' }
  return map[props.variant] || map.primary
})
</script>
