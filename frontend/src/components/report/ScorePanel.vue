<script setup lang="ts">
defineProps<{
  scores: Array<{ name: string; score: number; max_score: number; comment: string }>
  totalScore: number
  grade: string
}>()
</script>

<template>
  <div class="rounded-lg border border-border bg-surface p-5">
    <h3 class="text-lg font-semibold text-text-primary mb-4">评分详情</h3>

    <!-- Overall Score -->
    <div class="flex items-center justify-between mb-4 pb-4 border-b border-border">
      <span class="text-sm font-medium text-text-muted">总分</span>
      <div class="flex items-center gap-2">
        <span class="text-2xl font-bold" :class="
          totalScore >= 80 ? 'text-emerald-600' :
          totalScore >= 60 ? 'text-amber-600' :
          'text-red-600'
        ">
          {{ totalScore }}
        </span>
        <span class="text-lg font-semibold px-2 py-0.5 rounded" :class="
          totalScore >= 80 ? 'bg-emerald-50 text-emerald-700 dark:bg-emerald-900/20 dark:text-emerald-300' :
          totalScore >= 60 ? 'bg-amber-50 text-amber-700 dark:bg-amber-900/20 dark:text-amber-300' :
          'bg-red-50 text-red-700 dark:bg-red-900/20 dark:text-red-300'
        ">
          {{ grade }}
        </span>
      </div>
    </div>

    <!-- Dimension Scores -->
    <div class="space-y-3">
      <div v-for="s in scores" :key="s.name" class="flex items-center justify-between gap-3">
        <div class="min-w-0 flex-1">
          <div class="flex items-center justify-between mb-1">
            <span class="text-sm font-medium text-text-primary">{{ s.name }}</span>
            <span class="text-sm font-semibold text-text-muted">{{ s.score }}/{{ s.max_score }}</span>
          </div>
          <div class="h-2 rounded-full bg-surface-muted overflow-hidden">
            <div
              class="h-full rounded-full transition-all duration-500"
              :class="
                s.score >= 80 ? 'bg-emerald-500' :
                s.score >= 60 ? 'bg-amber-500' :
                'bg-red-500'
              "
              :style="{ width: (s.score / s.max_score * 100) + '%' }"
            />
          </div>
          <p v-if="s.comment" class="mt-1 text-xs text-text-muted">{{ s.comment }}</p>
        </div>
      </div>
    </div>
  </div>
</template>
