<script setup lang="ts">
/**
 * ProjectPicker - optional, searchable project ownership selector for tool runs.
 * It loads a small page of project summaries from the existing project API and
 * leaves standalone runs available through the empty option.
 */
import { onMounted, ref, watch } from 'vue'
import { Search } from '@lucide/vue'
import { listProjects } from '@/api/projects'
import type { ProjectSummary } from '@/types/project'

const props = defineProps<{
  modelValue: string
  inputId: string
  helpId?: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const query = ref('')
const projects = ref<ProjectSummary[]>([])
const loading = ref(false)
const loadError = ref<string | null>(null)

async function loadProjects() {
  loading.value = true
  loadError.value = null
  try {
    const result = await listProjects({
      search: query.value.trim() || undefined,
      page: 1,
      page_size: 100,
    })
    projects.value = result.items
  } catch (error: any) {
    projects.value = []
    loadError.value = error.message || '无法加载项目列表'
  } finally {
    loading.value = false
  }
}

function onSelect(event: Event) {
  emit('update:modelValue', (event.target as HTMLSelectElement).value)
}

watch(query, () => {
  void loadProjects()
})

onMounted(() => {
  void loadProjects()
})
</script>

<template>
  <div class="space-y-2">
    <div class="relative">
      <Search :size="15" class="absolute left-3 top-1/2 -translate-y-1/2 text-text-muted pointer-events-none" />
      <input
        :id="`${inputId}-search`"
        v-model="query"
        type="search"
        class="w-full rounded-lg border border-border/80 bg-surface-muted py-2 pl-9 pr-3 text-sm text-text-primary placeholder:text-text-muted focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent"
        placeholder="搜索项目名称或描述"
        aria-label="搜索项目"
      />
    </div>

    <select
      :id="inputId"
      :value="props.modelValue"
      :aria-describedby="helpId"
      class="w-full rounded-lg border border-border/80 bg-surface-muted px-3 py-2 text-text-primary focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent"
      @change="onSelect"
    >
      <option value="">不关联项目</option>
      <option v-for="project in projects" :key="project.id" :value="project.id">
        {{ project.name }}
      </option>
    </select>

    <p v-if="loading" class="text-xs text-text-muted">正在加载项目…</p>
    <p v-else-if="loadError" class="text-xs text-danger">{{ loadError }}</p>
    <p v-else-if="query && !projects.length" class="text-xs text-text-muted">没有匹配的项目。</p>
  </div>
</template>
