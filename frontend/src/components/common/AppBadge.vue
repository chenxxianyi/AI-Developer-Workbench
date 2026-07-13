<template>
  <span :class="badgeClasses">
    <slot />
  </span>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(defineProps<{
  variant?: 'default' | 'success' | 'warning' | 'danger' | 'info'
  size?: 'sm' | 'md'
}>(), {
  variant: 'default',
  size: 'sm',
})

const badgeClasses = computed(() => {
  const base = 'inline-flex items-center font-medium rounded-full'
  const sizes: Record<string, string> = {
    sm: 'px-2 py-0.5 text-xs',
    md: 'px-3 py-1 text-sm',
  }
  const variants: Record<string, string> = {
    default: 'bg-[var(--color-surface-muted)] text-[var(--color-text-secondary)]',
    success: 'bg-green-50 text-[var(--color-success)]',
    warning: 'bg-amber-50 text-[var(--color-warning)]',
    danger: 'bg-red-50 text-[var(--color-danger)]',
    info: 'bg-[var(--color-info-soft)] text-[var(--color-info)]',
  }
  return `${base} ${sizes[props.size]} ${variants[props.variant]}`
})
</script>
