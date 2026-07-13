<template>
  <div :class="skeletonClass" class="animate-pulse bg-[var(--color-surface-muted)] rounded" />
</template>

<script setup lang="ts">
import { computed } from 'vue'
const props = withDefaults(defineProps<{ width?: string; height?: string; circle?: boolean }>(), { width: '100%', height: '16px' })
const skeletonClass = computed(() => props.circle ? 'rounded-full' : 'rounded-[var(--radius-sm)]')
</script>
