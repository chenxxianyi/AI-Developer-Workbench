<script setup lang="ts">
import { Lightbulb } from '@lucide/vue'

defineProps<{
  recommendations: string[]
}>()
</script>

<template>
  <div v-if="recommendations.length" class="rounded-lg border border-border bg-surface p-5">
    <h3 class="text-lg font-semibold text-text-primary mb-4">
      <Lightbulb :size="18" class="inline-block mr-1.5 text-amber-500" />
      改进建议
    </h3>
    <ul class="space-y-2">
      <li
        v-for="(rec, idx) in recommendations"
        :key="idx"
        class="flex items-start gap-2 text-sm text-text-secondary"
      >
        <span class="mt-1.5 h-1.5 w-1.5 flex-shrink-0 rounded-full bg-accent" />
        <span>{{ rec }}</span>
      </li>
    </ul>
  </div>
</template>
