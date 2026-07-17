<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  AlertCircle,
  CheckCircle2,
  Code2,
  Database,
  FileJson,
  Files,
  GitBranch,
  Layers3,
  ListChecks,
  Loader2,
  Package,
  PenTool,
  RefreshCw,
  Save,
  Workflow,
} from '@lucide/vue'
import apiClient from '@/api/client'
import ProjectStageShell from '@/components/project/ProjectStageShell.vue'

interface BlueprintPageItem { name: string; route: string; purpose?: string; key_sections?: string[] }
interface BlueprintFeature { id: string; name: string; priority: string; description?: string; acceptance_criteria?: string[] }
interface BlueprintContent {
  schema_version?: number
  app_type?: string
  product_positioning?: string
  tech_stack?: string
  ui_style?: string
  pages?: BlueprintPageItem[]
  user_flows?: Array<{ name: string; steps: string[] }>
  features?: BlueprintFeature[]
  components?: string[]
  interaction_rules?: string[]
  state_model?: Array<{ name: string; description: string }>
  domain_rules?: string[]
  api_endpoints?: Array<{ method: string; path: string; description: string }>
  data_models?: Array<{ name: string; fields: string[] }>
  visual_system?: { style?: string; colors?: string; layout?: string; accessibility?: string }
  acceptance_criteria?: string[]
  test_plan?: string[]
  implementation_notes?: string[]
  open_questions?: string[]
}
interface BlueprintRecord { id?: string; content?: string | BlueprintContent; status?: string; version?: number }

const route = useRoute()
const router = useRouter()
const projectId = route.params.projectId as string
const blueprintRecord = ref<BlueprintRecord | null>(null)
const loading = ref(true)
const generating = ref(false)
const saving = ref(false)
const confirming = ref(false)
const editing = ref(false)
const jsonDraft = ref('')
const error = ref('')
const success = ref('')
const generationElapsed = ref(0)
const requirementSchemaVersion = ref(0)
const requirementMustFeatures = ref<string[]>([])
let generationTimer: ReturnType<typeof setInterval> | undefined
let generationController: AbortController | undefined

const blueprint = computed<BlueprintContent | null>(() => {
  if (!blueprintRecord.value) return null
  const source = blueprintRecord.value.content
  if (!source) return null
  if (typeof source !== 'string') return source
  try { return JSON.parse(source) as BlueprintContent } catch { return null }
})
const features = computed(() => blueprint.value?.features ?? [])
const mustFeatures = computed(() => features.value.filter((item) => item.priority === 'must'))
const isSuperseded = computed(() => blueprintRecord.value?.status === 'superseded')
const isConfirmed = computed(() => blueprintRecord.value?.status === 'confirmed')
const confirmIssues = computed(() => {
  const issues: string[] = []
  if (isSuperseded.value) return ['需求已经更新，这份蓝图已失效，请根据最新需求重新生成']
  if (!blueprint.value?.pages?.length) issues.push('缺少页面或屏幕结构')
  if (requirementSchemaVersion.value >= 2 && blueprint.value?.schema_version !== 2) issues.push('蓝图版本过旧，不支持最新需求结构')
  if (!features.value.length) issues.push('缺少功能清单')
  else if (mustFeatures.value.length < requirementMustFeatures.value.length) issues.push(`必须功能覆盖不足：需求 ${requirementMustFeatures.value.length} 项，蓝图仅 ${mustFeatures.value.length} 项`)
  if (mustFeatures.value.some((item) => !item.acceptance_criteria?.length)) issues.push('必须功能缺少验收标准')
  if (blueprint.value?.open_questions?.length) issues.push(`仍有 ${blueprint.value.open_questions.length} 个未决问题`)
  return issues
})

async function loadBlueprint() {
  loading.value = true
  error.value = ''
  try { blueprintRecord.value = await apiClient.get(`/projects/${projectId}/blueprint`) }
  catch { blueprintRecord.value = null }
  finally { loading.value = false }
}

async function loadRequirementContext() {
  try {
    const requirement: any = await apiClient.get(`/projects/${projectId}/requirements`)
    const content = typeof requirement?.content === 'string' ? JSON.parse(requirement.content) : requirement?.content
    requirementSchemaVersion.value = Number(content?.schema_version) || 0
    requirementMustFeatures.value = Array.isArray(content?.must_have_features) ? content.must_have_features : []
  } catch {
    requirementSchemaVersion.value = 0
    requirementMustFeatures.value = []
  }
}

async function generate() {
  if (generating.value) return
  generating.value = true
  generationElapsed.value = 0
  generationController = new AbortController()
  generationTimer = setInterval(() => { generationElapsed.value += 1 }, 1000)
  error.value = ''
  success.value = ''
  try {
    // Blueprint generation is bounded by the AI provider timeout on the backend.
    // Do not let the generic 100-second Axios timeout abort a still-running model call.
    blueprintRecord.value = await apiClient.post(
      `/projects/${projectId}/blueprint/generate`,
      undefined,
      { timeout: 0, signal: generationController.signal },
    )
    await loadRequirementContext()
    success.value = '蓝图已生成，请逐项评审后确认'
  } catch (err: any) { error.value = err.message || '蓝图生成失败，请稍后重试' }
  finally {
    if (generationTimer) clearInterval(generationTimer)
    generationTimer = undefined
    generationController = undefined
    generating.value = false
  }
}

function startEditing() {
  if (!blueprint.value) return
  jsonDraft.value = JSON.stringify(blueprint.value, null, 2)
  editing.value = true
  error.value = ''
}

async function saveDraft() {
  if (saving.value) return
  let parsed: BlueprintContent
  try { parsed = JSON.parse(jsonDraft.value) as BlueprintContent }
  catch { error.value = '蓝图 JSON 格式无效，请修正后再保存'; return }
  saving.value = true
  error.value = ''
  success.value = ''
  try {
    blueprintRecord.value = await apiClient.put(`/projects/${projectId}/blueprint`, { content: JSON.stringify(parsed) })
    editing.value = false
    success.value = '蓝图草稿已保存为新版本'
  } catch (err: any) { error.value = err.message || '蓝图保存失败' }
  finally { saving.value = false }
}

async function confirm() {
  if (!blueprintRecord.value || confirming.value) return
  if (confirmIssues.value.length) {
    error.value = `蓝图尚不可确认：${confirmIssues.value.join('；')}`
    return
  }
  confirming.value = true
  error.value = ''
  try {
    await apiClient.post(`/projects/${projectId}/blueprint/confirm`)
    await router.push(`/projects/${projectId}/generation`)
  } catch (err: any) { error.value = err.message || '蓝图确认失败' }
  finally { confirming.value = false }
}

function enterGeneration() {
  void router.push(`/projects/${projectId}/generation`)
}

onMounted(() => { void Promise.all([loadBlueprint(), loadRequirementContext()]) })
onBeforeUnmount(() => {
  if (generationTimer) clearInterval(generationTimer)
  generationController?.abort()
})
</script>

<template>
  <ProjectStageShell
    :icon="PenTool"
    title="产品与技术蓝图"
    description="评审功能、流程、状态、领域规则和验收标准；确认后的版本将作为代码生成的唯一输入。"
    :step-text="blueprintRecord ? `蓝图 v${blueprintRecord.version || 1}` : '蓝图阶段'"
  >
    <template #actions>
      <button v-if="blueprint && !editing" type="button" :disabled="generating" class="inline-flex min-h-10 items-center gap-2 rounded-lg border border-border bg-surface px-4 text-sm font-medium text-text-primary hover:bg-surface-muted disabled:opacity-60" @click="generate">
        <Loader2 v-if="generating" :size="16" class="animate-spin" /><RefreshCw v-else :size="16" />{{ generating ? '生成中...' : '重新生成' }}
      </button>
      <button v-if="blueprint && !editing && !isSuperseded" type="button" class="inline-flex min-h-10 items-center gap-2 rounded-lg border border-border bg-surface px-4 text-sm font-medium text-text-primary hover:bg-surface-muted" @click="startEditing">
        <FileJson :size="16" />编辑蓝图
      </button>
      <button v-if="blueprint && !editing && !isConfirmed" type="button" :disabled="confirming || confirmIssues.length > 0" class="inline-flex min-h-10 items-center gap-2 rounded-lg bg-accent px-4 text-sm font-semibold text-white hover:bg-accent/80 disabled:cursor-not-allowed disabled:opacity-50" @click="confirm">
        <Loader2 v-if="confirming" :size="16" class="animate-spin" /><CheckCircle2 v-else :size="16" />{{ confirming ? '确认中...' : '确认并进入生成' }}
      </button>
      <button v-if="blueprint && !editing && isConfirmed" type="button" class="inline-flex min-h-10 items-center gap-2 rounded-lg bg-accent px-4 text-sm font-semibold text-white hover:bg-accent/80" @click="enterGeneration"><CheckCircle2 :size="16" />进入代码生成</button>
    </template>

    <div v-if="error" role="alert" class="mb-5 flex items-center gap-2 rounded-lg border border-danger/20 bg-danger/10 px-4 py-3 text-sm text-danger"><AlertCircle :size="18" />{{ error }}</div>
    <div v-if="success" aria-live="polite" class="mb-5 flex items-center gap-2 rounded-lg border border-success/20 bg-success/10 px-4 py-3 text-sm text-success"><CheckCircle2 :size="18" />{{ success }}</div>

    <section v-if="generating" aria-live="polite" class="mb-5 flex items-start gap-3 rounded-lg border border-accent/20 bg-accent-soft/40 p-4">
      <Loader2 :size="20" class="mt-0.5 shrink-0 animate-spin text-accent" />
      <div><h2 class="font-semibold text-text-primary">AI 正在根据最新需求生成蓝图</h2><p class="mt-1 text-sm leading-6 text-text-secondary">已等待 {{ generationElapsed }} 秒。复杂蓝图通常需要 1–3 分钟，生成完成前请保持当前页面打开。</p></div>
    </section>

    <section v-if="isSuperseded && !loading" class="mb-5 flex flex-col gap-4 rounded-lg border border-warning/30 bg-warning/10 p-4 sm:flex-row sm:items-center sm:justify-between">
      <div class="flex items-start gap-3"><AlertCircle :size="20" class="mt-0.5 shrink-0 text-warning" /><div><h2 class="font-semibold text-text-primary">需求已更新，需要重新生成方案</h2><p class="mt-1 text-sm leading-6 text-text-secondary">当前展示的是旧版蓝图，不能再用于代码生成。重新生成后会使用你刚刚保存的最新需求。</p></div></div>
      <button type="button" :disabled="generating" class="inline-flex min-h-10 shrink-0 items-center justify-center gap-2 rounded-lg bg-accent px-4 text-sm font-semibold text-white hover:bg-accent/80 disabled:opacity-60" @click="generate"><Loader2 v-if="generating" :size="16" class="animate-spin" /><RefreshCw v-else :size="16" />{{ generating ? '正在生成...' : '根据最新需求重新生成' }}</button>
    </section>

    <div v-if="loading" class="grid gap-6 lg:grid-cols-[minmax(0,1fr)_300px]"><div class="h-[540px] animate-pulse rounded-lg border border-border bg-surface-muted" /><div class="h-80 animate-pulse rounded-lg border border-border bg-surface-muted" /></div>

    <section v-else-if="editing" class="rounded-lg border border-border bg-surface p-5 shadow-sm sm:p-6">
      <div class="mb-4 flex items-start justify-between gap-4">
        <div><h2 class="font-semibold text-text-primary">编辑蓝图 JSON</h2><p class="mt-1 text-sm text-text-secondary">保存会创建新草稿版本，已确认版本不会被覆盖。</p></div>
        <FileJson :size="22" class="text-accent" />
      </div>
      <textarea v-model="jsonDraft" rows="30" spellcheck="false" class="w-full rounded-lg border border-border bg-[#171717] p-4 font-mono text-xs leading-6 text-zinc-200 focus-visible:ring-2 focus-visible:ring-accent focus:outline-none" aria-label="蓝图 JSON" />
      <div class="mt-4 flex justify-end gap-3">
        <button type="button" class="min-h-10 rounded-lg border border-border px-4 text-sm font-medium text-text-secondary hover:bg-surface-muted" @click="editing = false">取消</button>
        <button type="button" :disabled="saving" class="inline-flex min-h-10 items-center gap-2 rounded-lg bg-accent px-4 text-sm font-semibold text-white disabled:opacity-60" @click="saveDraft"><Loader2 v-if="saving" :size="16" class="animate-spin" /><Save v-else :size="16" />保存新版本</button>
      </div>
    </section>

    <div v-else-if="blueprint" class="grid gap-6 lg:grid-cols-[minmax(0,1fr)_300px]">
      <div class="space-y-6">
        <section class="rounded-lg border border-border bg-surface p-5 shadow-sm sm:p-6">
          <div class="flex items-start gap-3"><div class="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-accent-soft text-accent"><GitBranch :size="20" /></div><div><p class="text-xs font-semibold uppercase tracking-wide text-text-muted">{{ blueprint.app_type || 'application' }}</p><h2 class="mt-1 text-lg font-semibold text-text-primary">{{ blueprint.product_positioning }}</h2><p class="mt-2 text-sm leading-6 text-text-secondary">{{ blueprint.ui_style }}</p></div></div>
        </section>

        <section class="rounded-lg border border-border bg-surface p-5 shadow-sm sm:p-6">
          <div class="mb-4 flex items-center justify-between"><div class="flex items-center gap-2"><ListChecks :size="18" class="text-accent" /><h2 class="font-semibold text-text-primary">功能与验收条件</h2></div><span class="text-xs text-text-muted">{{ features.length }} 项</span></div>
          <div class="space-y-3">
            <article v-for="feature in features" :key="feature.id" class="rounded-lg border border-border bg-surface-muted/50 p-4">
              <div class="flex flex-wrap items-center gap-2"><code class="text-xs text-accent">{{ feature.id }}</code><h3 class="font-semibold text-text-primary">{{ feature.name }}</h3><span class="rounded-full border px-2 py-0.5 text-xs" :class="feature.priority === 'must' ? 'border-danger/20 bg-danger/10 text-danger' : 'border-border text-text-muted'">{{ feature.priority }}</span></div>
              <p v-if="feature.description" class="mt-2 text-sm text-text-secondary">{{ feature.description }}</p>
              <ul class="mt-3 space-y-1.5 text-sm text-text-secondary"><li v-for="criterion in feature.acceptance_criteria" :key="criterion" class="flex gap-2"><CheckCircle2 :size="15" class="mt-0.5 shrink-0 text-success" />{{ criterion }}</li></ul>
            </article>
          </div>
        </section>

        <section class="rounded-lg border border-border bg-surface p-5 shadow-sm sm:p-6">
          <div class="mb-4 flex items-center gap-2"><Files :size="18" class="text-accent" /><h2 class="font-semibold text-text-primary">页面与用户流程</h2></div>
          <div class="grid gap-4 md:grid-cols-2">
            <article v-for="page in blueprint.pages" :key="page.route" class="rounded-lg border border-border p-4"><div class="flex items-center justify-between gap-3"><h3 class="font-semibold text-text-primary">{{ page.name }}</h3><code class="rounded bg-surface-muted px-2 py-1 text-xs text-text-muted">{{ page.route }}</code></div><p class="mt-2 text-sm text-text-secondary">{{ page.purpose }}</p><div class="mt-3 flex flex-wrap gap-1.5"><span v-for="section in page.key_sections" :key="section" class="rounded-md bg-accent-soft px-2 py-1 text-xs text-accent">{{ section }}</span></div></article>
          </div>
          <div v-if="blueprint.user_flows?.length" class="mt-5 space-y-3 border-t border-border pt-5"><article v-for="flow in blueprint.user_flows" :key="flow.name"><h3 class="text-sm font-semibold text-text-primary">{{ flow.name }}</h3><ol class="mt-2 flex flex-wrap items-center gap-2 text-xs text-text-secondary"><li v-for="(step, index) in flow.steps" :key="step" class="flex items-center gap-2"><span class="rounded-md border border-border bg-surface-muted px-2 py-1">{{ index + 1 }}. {{ step }}</span><span v-if="index < flow.steps.length - 1">→</span></li></ol></article></div>
        </section>

        <section class="grid gap-6 md:grid-cols-2">
          <div class="rounded-lg border border-border bg-surface p-5 shadow-sm"><div class="mb-3 flex items-center gap-2"><Workflow :size="18" class="text-accent" /><h2 class="font-semibold text-text-primary">状态与交互</h2></div><dl class="space-y-3"><div v-for="state in blueprint.state_model" :key="state.name"><dt class="text-sm font-medium text-text-primary">{{ state.name }}</dt><dd class="mt-0.5 text-sm text-text-secondary">{{ state.description }}</dd></div></dl><ul class="mt-4 space-y-2 border-t border-border pt-4 text-sm text-text-secondary"><li v-for="rule in blueprint.interaction_rules" :key="rule">• {{ rule }}</li></ul></div>
          <div class="rounded-lg border border-border bg-surface p-5 shadow-sm"><div class="mb-3 flex items-center gap-2"><Layers3 :size="18" class="text-accent" /><h2 class="font-semibold text-text-primary">领域规则</h2></div><ul class="space-y-2 text-sm leading-6 text-text-secondary"><li v-for="rule in blueprint.domain_rules" :key="rule">• {{ rule }}</li></ul></div>
        </section>

        <section class="rounded-lg border border-border bg-surface p-5 shadow-sm sm:p-6">
          <div class="mb-4 flex items-center gap-2"><ListChecks :size="18" class="text-accent" /><h2 class="font-semibold text-text-primary">全局验收与测试计划</h2></div>
          <div class="grid gap-5 md:grid-cols-2"><div><h3 class="text-sm font-semibold text-text-primary">验收标准</h3><ul class="mt-2 space-y-2 text-sm text-text-secondary"><li v-for="item in blueprint.acceptance_criteria" :key="item" class="flex gap-2"><CheckCircle2 :size="15" class="mt-0.5 shrink-0 text-success" />{{ item }}</li></ul></div><div><h3 class="text-sm font-semibold text-text-primary">测试计划</h3><ul class="mt-2 space-y-2 text-sm text-text-secondary"><li v-for="item in blueprint.test_plan" :key="item">• {{ item }}</li></ul></div></div>
        </section>
      </div>

      <aside class="h-fit space-y-5 lg:sticky lg:top-6">
        <section class="rounded-lg border border-border bg-surface p-5 shadow-sm"><div class="mb-3 flex items-center gap-2"><Package :size="18" class="text-accent" /><h2 class="font-semibold text-text-primary">技术方案</h2></div><p class="text-sm leading-6 text-text-secondary">{{ blueprint.tech_stack }}</p><div class="mt-4 flex flex-wrap gap-1.5"><span v-for="component in blueprint.components" :key="component" class="rounded-md border border-border bg-surface-muted px-2 py-1 text-xs text-text-secondary">{{ component }}</span></div></section>
        <section v-if="blueprint.api_endpoints?.length || blueprint.data_models?.length" class="rounded-lg border border-border bg-surface p-5 shadow-sm"><div class="mb-3 flex items-center gap-2"><Database :size="18" class="text-accent" /><h2 class="font-semibold text-text-primary">数据与 API</h2></div><div class="space-y-3"><div v-for="api in blueprint.api_endpoints" :key="`${api.method}-${api.path}`" class="text-xs"><code class="text-accent">{{ api.method }} {{ api.path }}</code><p class="mt-1 text-text-secondary">{{ api.description }}</p></div><div v-for="model in blueprint.data_models" :key="model.name" class="border-t border-border pt-3"><p class="text-sm font-medium text-text-primary">{{ model.name }}</p><p class="mt-1 text-xs text-text-muted">{{ model.fields.join(' · ') }}</p></div></div></section>
        <section class="rounded-lg border border-border bg-surface p-5 shadow-sm"><h2 class="font-semibold text-text-primary">确认门禁</h2><div v-if="confirmIssues.length" class="mt-3 space-y-2"><p v-for="issue in confirmIssues" :key="issue" class="flex gap-2 text-sm text-danger"><AlertCircle :size="16" class="mt-0.5 shrink-0" />{{ issue }}</p></div><p v-else class="mt-3 flex gap-2 text-sm text-success"><CheckCircle2 :size="16" />蓝图已满足确认条件</p><dl class="mt-4 space-y-2 border-t border-border pt-4 text-sm"><div class="flex justify-between"><dt class="text-text-muted">状态</dt><dd class="font-medium" :class="isSuperseded ? 'text-warning' : 'text-text-primary'">{{ isSuperseded ? '需求更新后已失效' : isConfirmed ? '已确认' : blueprintRecord?.status }}</dd></div><div class="flex justify-between"><dt class="text-text-muted">版本</dt><dd class="font-medium text-text-primary">v{{ blueprintRecord?.version }}</dd></div><div class="flex justify-between"><dt class="text-text-muted">必须功能</dt><dd class="font-medium text-text-primary">{{ mustFeatures.length }} / {{ requirementMustFeatures.length }}</dd></div></dl></section>
        <section v-if="blueprint.open_questions?.length" class="rounded-lg border border-warning/30 bg-warning/10 p-5"><h2 class="font-semibold text-text-primary">未决问题</h2><ul class="mt-3 space-y-2 text-sm text-text-secondary"><li v-for="item in blueprint.open_questions" :key="item">• {{ item }}</li></ul></section>
      </aside>
    </div>

    <section v-else class="flex min-h-96 flex-col items-center justify-center rounded-lg border border-border bg-surface px-6 py-12 text-center shadow-sm"><div class="mb-4 flex h-14 w-14 items-center justify-center rounded-lg bg-accent-soft text-accent"><PenTool :size="26" /></div><h2 class="text-lg font-semibold text-text-primary">尚未生成蓝图</h2><p class="mt-2 max-w-md text-sm leading-6 text-text-secondary">系统将根据结构化需求生成页面、功能、状态、领域规则和测试计划。</p><button type="button" :disabled="generating" class="mt-6 inline-flex min-h-10 items-center gap-2 rounded-lg bg-accent px-5 text-sm font-semibold text-white disabled:opacity-60" @click="generate"><Loader2 v-if="generating" :size="16" class="animate-spin" /><Code2 v-else :size="16" />{{ generating ? '生成中...' : '生成完整蓝图' }}</button></section>
  </ProjectStageShell>
</template>
