<script setup lang="ts">
import { onBeforeUnmount, ref, watch } from 'vue'
import type { Issue } from '@/types/report'
const props = defineProps<{ file: File | null; viewport: 'desktop' | 'mobile'; viewportLabel: string; issues: Issue[] }>()
const imageUrl = ref('')
watch(() => props.file, (file) => { if (imageUrl.value) URL.revokeObjectURL(imageUrl.value); imageUrl.value = file ? URL.createObjectURL(file) : '' }, { immediate: true })
onBeforeUnmount(() => { if (imageUrl.value) URL.revokeObjectURL(imageUrl.value) })
</script>
<template>
  <figure v-if="imageUrl" class="rounded-xl border border-border bg-surface p-3" :data-testid="`annotation-${viewport}`">
    <figcaption class="mb-2 flex items-center justify-between text-sm"><strong class="text-text-primary">{{ viewport === 'desktop' ? 'Desktop' : 'Mobile' }}</strong><span class="text-text-muted">{{ viewportLabel }}</span></figcaption>
    <div class="relative overflow-hidden rounded-lg bg-surface-muted">
      <img :src="imageUrl" :alt="`${viewportLabel} ????`" class="block h-auto w-full" />
      <template v-for="(issue, index) in issues.filter(i => i.viewport === viewport && i.region)" :key="`${issue.title}-${index}`">
        <span class="absolute flex min-h-6 min-w-6 items-center justify-center rounded-full bg-red-600 px-1 text-xs font-bold text-white shadow" :style="{ left: `${issue.region!.x}%`, top: `${issue.region!.y}%`, width: `${issue.region!.width}%`, height: `${issue.region!.height}%` }" :title="issue.title" :aria-label="`${index + 1}. ${issue.title}`">{{ index + 1 }}<span class="sr-only">{{ issue.title }}</span></span>
      </template>
    </div>
  </figure>
</template>
