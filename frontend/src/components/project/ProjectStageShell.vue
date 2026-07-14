<script setup lang="ts">
import type { Component } from 'vue'

withDefaults(
  defineProps<{
    icon: Component
    title: string
    description: string
    stepText?: string
  }>(),
  {
    stepText: '',
  },
)
</script>

<template>
  <section>
    <header class="mb-6 flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
      <div class="flex items-start gap-3">
        <div class="flex h-12 w-12 shrink-0 items-center justify-center rounded-lg bg-accent text-white shadow-sm">
          <component :is="icon" :size="24" />
        </div>
        <div class="min-w-0">
          <h1 class="text-2xl font-bold text-text-primary sm:text-3xl">{{ title }}</h1>
          <p class="mt-1 max-w-2xl text-sm text-text-secondary sm:text-base">{{ description }}</p>
        </div>
      </div>

      <div class="flex flex-col gap-3 sm:items-end">
        <div
          v-if="stepText"
          class="inline-flex w-fit items-center gap-2 rounded-full border border-accent/20 bg-accent-soft px-3 py-1.5 text-sm font-medium text-accent"
        >
          <span class="h-1.5 w-1.5 rounded-full bg-accent" />
          {{ stepText }}
        </div>
        <div v-if="$slots.actions" class="flex flex-wrap items-center gap-2">
          <slot name="actions" />
        </div>
      </div>
    </header>

    <slot />
  </section>
</template>
