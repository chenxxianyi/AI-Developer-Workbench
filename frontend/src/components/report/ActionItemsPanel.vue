<script setup lang="ts">
/**
 * ActionItemsPanel — grouped report action plan with session completion state.
 */
import { computed, onMounted, ref } from 'vue'
import type { ActionItem } from '@/types/report'
import {
  buildActionItemsMarkdown,
  groupActionItems,
  normalizeActionItems,
  type DisplayActionItem,
} from '@/utils/actionItems'
import { copyToClipboard } from '@/utils/clipboard'
import { downloadBlob } from '@/utils/download'
import { downloadGitHubIssues } from '@/api/reports'
import { Check, ClipboardList, Download } from '@lucide/vue'
import ActionItemCard from './ActionItemCard.vue'

const props = defineProps<{
  reportId: string
  reportTitle: string
  actionItems?: ActionItem[]
  recommendations?: string[]
}>()

const completedIds = ref<Set<string>>(new Set())
const copiedItemId = ref<string | null>(null)
const copiedAll = ref(false)
const copyError = ref('')
const issueExportError = ref('')
const issueExportSuccess = ref(false)
const issueExporting = ref(false)

const storageKey = computed(() => `ai-workbench-action-items:${props.reportId}`)
const items = computed(() => normalizeActionItems(props.actionItems, props.recommendations ?? []))
const visibleTopItems = computed(() => items.value.slice(0, 3))
const remainingItems = computed(() => items.value.slice(3))
const grouped = computed(() => groupActionItems(remainingItems.value))
const hasItems = computed(() => items.value.length > 0)

onMounted(() => {
  try {
    const raw = sessionStorage.getItem(storageKey.value)
    if (raw) {
      completedIds.value = new Set(JSON.parse(raw))
    }
  } catch {
    completedIds.value = new Set()
  }
})

function toggle(id: string) {
  const next = new Set(completedIds.value)
  if (next.has(id)) {
    next.delete(id)
  } else {
    next.add(id)
  }
  completedIds.value = next
  sessionStorage.setItem(storageKey.value, JSON.stringify([...next]))
}

async function copyItem(item: DisplayActionItem) {
  try {
    copyError.value = ''
    await copyToClipboard(item.suggested_prompt || item.issue_body || item.title)
    copiedItemId.value = item.id
    window.setTimeout(() => {
      if (copiedItemId.value === item.id) copiedItemId.value = null
    }, 1800)
  } catch (error) {
    copyError.value = error instanceof Error ? error.message : '复制失败，请手动复制'
  }
}

async function copyAll() {
  try {
    copyError.value = ''
    await copyToClipboard(markdown.value)
    copiedAll.value = true
    window.setTimeout(() => { copiedAll.value = false }, 1800)
  } catch (error) {
    copyError.value = error instanceof Error ? error.message : '复制失败，请手动复制'
  }
}

async function downloadChecklist() {
  await downloadBlob(
    new Blob([markdown.value], { type: 'text/markdown;charset=utf-8' }),
    `${props.reportId}-action-items.md`,
  )
}

async function downloadIssueDrafts() {
  try {
    issueExportError.value = ''
    issueExportSuccess.value = false
    issueExporting.value = true
    await downloadGitHubIssues(props.reportId)
    issueExportSuccess.value = true
    window.setTimeout(() => { issueExportSuccess.value = false }, 1800)
  } catch (error) {
    issueExportError.value = error instanceof Error ? error.message : '导出 Issue 草稿失败'
  } finally {
    issueExporting.value = false
  }
}

const markdown = computed(() => buildActionItemsMarkdown(props.reportId, props.reportTitle, items.value))

const groupMeta = [
  { key: 'high', label: '高优先级' },
  { key: 'medium', label: '中优先级' },
  { key: 'low', label: '低优先级' },
] as const
</script>

<template>
  <section class="rounded-lg border border-border bg-surface p-5">
    <div class="mb-4 flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
      <div>
        <div class="mb-1 flex items-center gap-2">
          <ClipboardList :size="18" class="text-accent" />
          <h3 class="text-lg font-semibold text-text-primary">行动计划</h3>
        </div>
        <p class="text-sm text-text-muted">
          {{ hasItems ? `优先处理前 ${visibleTopItems.length} 项，勾选状态保存在当前浏览器会话。` : '暂无可执行行动项。' }}
        </p>
      </div>

      <div v-if="hasItems" class="flex flex-wrap gap-2">
        <button
          class="inline-flex items-center gap-1.5 rounded-md border border-border bg-surface px-3 py-1.5 text-sm font-medium text-text-primary transition-smooth hover:bg-surface-hover"
          @click="copyAll"
        >
          <Check v-if="copiedAll" :size="15" class="text-success" />
          <ClipboardList v-else :size="15" />
          {{ copiedAll ? '已复制' : '复制全部' }}
        </button>
        <button
          class="inline-flex items-center gap-1.5 rounded-md border border-border bg-surface px-3 py-1.5 text-sm font-medium text-text-primary transition-smooth hover:bg-surface-hover"
          @click="downloadChecklist"
        >
          <Download :size="15" />
          下载清单
        </button>
        <button
          class="inline-flex items-center gap-1.5 rounded-md border border-border bg-surface px-3 py-1.5 text-sm font-medium text-text-primary transition-smooth hover:bg-surface-hover"
          :disabled="issueExporting"
          @click="downloadIssueDrafts"
        >
          <Check v-if="issueExportSuccess" :size="15" class="text-success" />
          <Download v-else :size="15" />
          {{ issueExporting ? '导出中' : issueExportSuccess ? '已导出' : 'Issue 草稿' }}
        </button>
      </div>
    </div>

    <div v-if="!hasItems" class="rounded-md border border-dashed border-border py-8 text-center text-sm text-text-muted">
      旧报告中也没有可降级展示的建议。
    </div>

    <div v-else class="space-y-5">
      <div v-if="copyError" class="rounded-md border border-danger/30 bg-danger/5 px-3 py-2 text-sm text-danger" role="alert">
        {{ copyError }}
      </div>
      <div v-if="issueExportError" class="rounded-md border border-danger/30 bg-danger/5 px-3 py-2 text-sm text-danger" role="alert">
        {{ issueExportError }}
      </div>

      <div class="rounded-md border border-accent/20 bg-accent-soft/40 p-4">
        <div class="mb-3 text-sm font-semibold text-text-primary">首要行动</div>
        <div class="space-y-3">
          <ActionItemCard
            v-for="item in visibleTopItems"
            :key="`top-${item.id}`"
            :item="item"
            :completed="completedIds.has(item.id)"
            :copied="copiedItemId === item.id"
            @toggle="toggle"
            @copy="copyItem"
          />
        </div>
      </div>

      <div v-for="group in groupMeta" :key="group.key">
        <div v-if="grouped[group.key].length" class="mb-2 text-sm font-semibold text-text-muted">
          {{ group.label }}
        </div>
        <div v-if="grouped[group.key].length" class="space-y-3">
          <ActionItemCard
            v-for="item in grouped[group.key]"
            :key="item.id"
            :item="item"
            :completed="completedIds.has(item.id)"
            :copied="copiedItemId === item.id"
            @toggle="toggle"
            @copy="copyItem"
          />
        </div>
      </div>
    </div>
  </section>
</template>
