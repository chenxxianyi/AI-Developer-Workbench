<template>
  <div :class="cardClasses">
    <div v-if="$slots.header || title" class="px-5 py-4 border-b border-[var(--color-border)]">
      <h3 v-if="title" class="text-base font-semibold text-[var(--color-text-primary)]">{{ title }}</h3>
      <slot name="header" />
    </div>
    <div class="px-5 py-4" :class="padding ? '' : '!p-0'">
      <slot />
    </div>
    <div v-if="$slots.footer" class="px-5 py-3 border-t border-[var(--color-border)] bg-[var(--color-surface-muted)]/50 rounded-b-[var(--radius-card)]">
      <slot name="footer" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(defineProps<{
  title?: string
  padding?: boolean
  hover?: boolean
}>(), {
  padding: true,
})

const cardClasses = computed(() =>
  `bg-[var(--color-surface)] rounded-[var(--radius-card)] border border-[var(--color-border)] shadow-[var(--shadow-sm)] ${props.hover ? 'hover:shadow-[var(--shadow-md)] transition-shadow duration-200' : ''}`
)
</script>
