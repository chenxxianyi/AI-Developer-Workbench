<template>
  <button
    :class="buttonClasses"
    :disabled="disabled || loading"
    :type="type"
  >
    <span v-if="loading" class="mr-2 inline-block h-4 w-4 animate-spin rounded-full border-2 border-current border-r-transparent" />
    <slot />
  </button>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(defineProps<{
  variant?: 'primary' | 'secondary' | 'danger' | 'ghost'
  size?: 'sm' | 'md' | 'lg'
  disabled?: boolean
  loading?: boolean
  type?: 'button' | 'submit' | 'reset'
}>(), {
  variant: 'primary',
  size: 'md',
  type: 'button',
})

const buttonClasses = computed(() => {
  const base = 'inline-flex items-center justify-center font-medium rounded-lg transition-all duration-200 focus-visible:outline-2 focus-visible:outline-offset-2 disabled:opacity-50 disabled:cursor-not-allowed'
  const sizes: Record<string, string> = {
    sm: 'px-3 py-1.5 text-sm gap-1.5',
    md: 'px-4 py-2 text-sm gap-2',
    lg: 'px-6 py-3 text-base gap-2',
  }
  const variants: Record<string, string> = {
    primary: 'bg-[var(--color-accent)] text-white hover:brightness-110 active:brightness-90 focus-visible:outline-[var(--color-accent)]',
    secondary: 'bg-[var(--color-surface)] text-[var(--color-text-primary)] border border-[var(--color-border)] hover:bg-[var(--color-surface-muted)]',
    danger: 'bg-[var(--color-danger)] text-white hover:brightness-110 active:brightness-90',
    ghost: 'text-[var(--color-text-secondary)] hover:bg-[var(--color-surface-muted)]',
  }
  return `${base} ${sizes[props.size]} ${variants[props.variant]}`
})
</script>
