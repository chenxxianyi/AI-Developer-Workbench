<script setup lang="ts">
/**
 * QualityTrend — lightweight inline-SVG sparkline for the 30-day quality trend.
 *
 * Renders an average-score line plus per-day report-count bars. When there are
 * fewer than 2 scored data points it shows an explanatory empty state instead
 * of a misleading flat line.
 */
import { computed } from 'vue'
import type { QualityTrendPoint } from '@/types/system'

const props = defineProps<{
  points: QualityTrendPoint[]
}>()

// Width/height of the SVG viewBox.
const W = 600
const H = 160
const PAD = 8

const scoredPoints = computed(() =>
  props.points.filter((p) => p.average_score !== null && p.average_score !== undefined),
)

const hasEnoughData = computed(() => scoredPoints.value.length >= 2)

const maxCount = computed(() =>
  Math.max(1, ...props.points.map((p) => p.report_count)),
)

// Map each day to an x coordinate across the full 30-day span.
const xFor = (i: number, total: number) => {
  if (total <= 1) return W / 2
  const span = W - PAD * 2
  return PAD + (span * i) / (total - 1)
}

// Map an average score (0-100) to a y coordinate.
const yFor = (score: number) => {
  const span = H - PAD * 2
  return PAD + span * (1 - score / 100)
}

const linePath = computed(() => {
  if (!hasEnoughData.value) return ''
  const pts = scoredPoints.value
  return pts
    .map((p, idx) => {
      // find the day index in the full points array for x positioning
      const dayIdx = props.points.indexOf(p)
      const x = xFor(dayIdx, props.points.length)
      const y = yFor(p.average_score!)
      return `${idx === 0 ? 'M' : 'L'}${x.toFixed(1)},${y.toFixed(1)}`
    })
    .join(' ')
})

const bars = computed(() =>
  props.points.map((p, i) => {
    const x = xFor(i, props.points.length)
    const barW = Math.max(2, (W - PAD * 2) / Math.max(1, props.points.length) - 2)
    const barH = (p.report_count / maxCount.value) * (H - PAD * 2) * 0.35
    return {
      x: x - barW / 2,
      y: H - PAD - barH,
      w: barW,
      h: barH,
      count: p.report_count,
    }
  }),
)
</script>

<template>
  <div class="rounded-lg border border-border bg-surface p-5">
    <div class="flex items-center justify-between mb-3">
      <h3 class="text-lg font-semibold text-text-primary">质量趋势（近 30 天）</h3>
      <span class="text-xs text-text-muted">按日聚合 · UTC</span>
    </div>

    <div v-if="!hasEnoughData" class="py-8 text-center">
      <p class="text-sm text-text-muted">评分数据点不足（少于 2 天），暂不绘制趋势线。</p>
      <p class="text-xs text-text-muted mt-1">运行更多评分型工具后会自动出现趋势。</p>
    </div>

    <svg
      v-else
      :viewBox="`0 0 ${W} ${H}`"
      class="w-full h-40"
      preserveAspectRatio="none"
      role="img"
      aria-label="近 30 天平均质量分趋势"
    >
      <!-- report-count bars -->
      <rect
        v-for="(b, i) in bars"
        :key="`bar-${i}`"
        :x="b.x"
        :y="b.y"
        :width="b.w"
        :height="b.h"
        class="fill-accent/20"
      />
      <!-- average-score line -->
      <path :d="linePath" fill="none" stroke="currentColor" class="text-accent" stroke-width="2" />
    </svg>

    <div class="flex items-center gap-4 mt-2 text-xs text-text-muted">
      <span class="flex items-center gap-1">
        <span class="inline-block w-3 h-0.5 bg-accent" /> 平均分
      </span>
      <span class="flex items-center gap-1">
        <span class="inline-block w-2.5 h-2.5 bg-accent/20 rounded-sm" /> 报告数
      </span>
    </div>
  </div>
</template>
