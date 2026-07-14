<script setup lang="ts">
/**
 * Project Form Page (create + edit)
 */
import { ref, onMounted, computed } from 'vue'
import { onBeforeRouteLeave, useRoute, useRouter } from 'vue-router'
import { useProjectStore } from '@/stores/projectStore'
import ToolFormSection from '@/components/tool/ToolFormSection.vue'
import { ArrowLeft } from '@lucide/vue'

const route = useRoute()
const router = useRouter()
const store = useProjectStore()

const isEdit = computed(() => !!route.params.projectId)
const projectId = computed(() => (route.params.projectId as string) || '')

const name = ref('')
const description = ref('')
const repoUrl = ref('')
const frontendStack = ref('')
const backendStack = ref('')
const database = ref('')
const uiStyle = ref('')
const codingRules = ref('')

const saving = ref(false)
const error = ref<string | null>(null)
const errors = ref<Record<string, string>>({})
const initialValues = ref('')

function currentValues(): string {
  return JSON.stringify({
    name: name.value,
    description: description.value,
    repoUrl: repoUrl.value,
    frontendStack: frontendStack.value,
    backendStack: backendStack.value,
    database: database.value,
    uiStyle: uiStyle.value,
    codingRules: codingRules.value,
  })
}

const isDirty = computed(() => initialValues.value !== '' && currentValues() !== initialValues.value)

onMounted(async () => {
  if (isEdit.value) {
    await store.fetchProject(projectId.value)
    const p = store.currentProject
    if (p) {
      name.value = p.name
      description.value = p.description
      repoUrl.value = p.repo_url
      frontendStack.value = p.frontend_stack
      backendStack.value = p.backend_stack
      database.value = p.database
      uiStyle.value = p.ui_style
      codingRules.value = p.coding_rules
    }
  }
  initialValues.value = currentValues()
})

async function handleSave() {
  errors.value = {}
  if (!name.value.trim()) errors.value.name = '请输入项目名称'
  if (name.value.trim().length > 128) errors.value.name = '项目名称最多 128 个字符'
  if (description.value.trim().length > 4000) errors.value.description = '描述最多 4000 个字符'
  if (repoUrl.value.trim()) {
    try {
      const parsed = new URL(repoUrl.value.trim())
      if (!['http:', 'https:'].includes(parsed.protocol)) errors.value.repo_url = '请输入 HTTP 或 HTTPS 地址'
    } catch {
      errors.value.repo_url = '请输入有效的仓库 URL'
    }
  }
  if (repoUrl.value.length > 512) errors.value.repo_url = 'URL 最多 512 个字符'
  if ([frontendStack, backendStack, uiStyle].some((field) => field.value.length > 256)) errors.value.stack = '技术栈和 UI 风格最多 256 个字符'
  if (database.value.length > 128) errors.value.database = '数据库字段最多 128 个字符'
  if (codingRules.value.length > 12000) errors.value.coding_rules = '编码规则最多 12000 个字符'
  if (Object.keys(errors.value).length) return

  saving.value = true
  error.value = null
  try {
    const payload = {
      name: name.value.trim(),
      description: description.value.trim(),
      repo_url: repoUrl.value.trim(),
      frontend_stack: frontendStack.value.trim(),
      backend_stack: backendStack.value.trim(),
      database: database.value.trim(),
      ui_style: uiStyle.value.trim(),
      coding_rules: codingRules.value.trim(),
    }
    if (isEdit.value) {
      await store.update(projectId.value, payload)
    } else {
      const created = await store.create(payload)
      router.replace(`/projects/${created.id}`)
      return
    }
    initialValues.value = currentValues()
    router.push(`/projects/${projectId.value}`)
  } catch (err: any) {
    error.value = err.message || '保存失败'
  } finally {
    saving.value = false
  }
}

function goBack() {
  router.back()
}

onBeforeRouteLeave(() => {
  if (!saving.value && isDirty.value && !window.confirm('离开后未保存的项目修改将丢失，确定离开吗？')) {
    return false
  }
  return true
})
</script>

<template>
  <div class="max-w-3xl mx-auto">
    <button class="flex items-center gap-2 text-sm text-text-muted hover:text-accent transition-smooth mb-4" @click="goBack">
      <ArrowLeft :size="16" />
      返回
    </button>

    <h1 class="text-2xl font-bold text-text-primary mb-6">{{ isEdit ? '编辑项目' : '新建项目' }}</h1>

    <div class="rounded-lg border border-border bg-surface p-6 space-y-1">
      <ToolFormSection label="项目名称" required id-for="project-name" :error="errors.name">
        <input
          id="project-name"
          v-model="name"
          type="text"
          class="w-full px-4 py-2 bg-surface-muted border border-border/80 rounded-lg focus-visible:ring-2 focus-visible:ring-accent focus:outline-none text-text-primary"
          placeholder="如: AI Workbench"
        />
      </ToolFormSection>

      <ToolFormSection label="描述" optional id-for="project-desc" :error="errors.description">
        <textarea
          id="project-desc"
          v-model="description"
          rows="3"
          class="w-full px-4 py-2 bg-surface-muted border border-border/80 rounded-lg focus-visible:ring-2 focus-visible:ring-accent focus:outline-none text-text-primary"
          placeholder="项目简介和目标…"
        />
      </ToolFormSection>

      <ToolFormSection label="仓库 URL" optional id-for="project-repo" :error="errors.repo_url">
        <input
          id="project-repo"
          v-model="repoUrl"
          type="text"
          class="w-full px-4 py-2 bg-surface-muted border border-border/80 rounded-lg focus-visible:ring-2 focus-visible:ring-accent focus:outline-none text-text-primary"
          placeholder="https://github.com/org/repo"
        />
      </ToolFormSection>

      <ToolFormSection label="前端技术栈" optional id-for="project-fe" :error="errors.stack">
        <input id="project-fe" v-model="frontendStack" type="text" class="w-full px-4 py-2 bg-surface-muted border border-border/80 rounded-lg focus-visible:ring-2 focus-visible:ring-accent focus:outline-none text-text-primary" placeholder="Vue 3 + TypeScript + Vite" />
      </ToolFormSection>

      <ToolFormSection label="后端技术栈" optional id-for="project-be" :error="errors.stack">
        <input id="project-be" v-model="backendStack" type="text" class="w-full px-4 py-2 bg-surface-muted border border-border/80 rounded-lg focus-visible:ring-2 focus-visible:ring-accent focus:outline-none text-text-primary" placeholder="Go + Gin + GORM" />
      </ToolFormSection>

      <ToolFormSection label="数据库" optional id-for="project-db" :error="errors.database">
        <input id="project-db" v-model="database" type="text" class="w-full px-4 py-2 bg-surface-muted border border-border/80 rounded-lg focus-visible:ring-2 focus-visible:ring-accent focus:outline-none text-text-primary" placeholder="MySQL 8" />
      </ToolFormSection>

      <ToolFormSection label="UI 风格" optional id-for="project-ui" :error="errors.stack">
        <input id="project-ui" v-model="uiStyle" type="text" class="w-full px-4 py-2 bg-surface-muted border border-border/80 rounded-lg focus-visible:ring-2 focus-visible:ring-accent focus:outline-none text-text-primary" placeholder="Tailwind CSS, 简洁现代" />
      </ToolFormSection>

      <ToolFormSection label="编码规则" optional id-for="project-rules" :error="errors.coding_rules">
        <textarea
          id="project-rules"
          v-model="codingRules"
          rows="4"
          class="w-full px-4 py-2 bg-surface-muted border border-border/80 rounded-lg focus-visible:ring-2 focus-visible:ring-accent focus:outline-none text-text-primary font-mono text-sm"
          placeholder="项目编码规范、命名约定、禁止事项…"
        />
      </ToolFormSection>

      <div v-if="error" class="p-3 bg-danger/10 border border-danger/20 rounded-lg text-danger text-sm">{{ error }}</div>

      <div class="flex gap-3 pt-2">
        <button
          :disabled="saving"
          class="px-6 py-2 bg-accent text-white rounded-lg hover:bg-accent/80 transition-smooth disabled:opacity-60 focus-visible:ring-2 focus-visible:ring-accent focus:outline-none"
          @click="handleSave"
        >
          {{ saving ? '保存中…' : '保存' }}
        </button>
        <button class="px-4 py-2 bg-surface-muted text-text-secondary rounded-lg hover:bg-border transition-smooth" @click="goBack">
          取消
        </button>
      </div>
    </div>
  </div>
</template>
