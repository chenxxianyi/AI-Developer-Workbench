<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, reactive, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  AlertCircle,
  ArrowLeft,
  ArrowRight,
  Check,
  CheckCircle2,
  ChevronDown,
  ClipboardList,
  HelpCircle,
  Laptop,
  Loader2,
  Monitor,
  Palette,
  Plus,
  RotateCcw,
  Save,
  Smartphone,
  Sparkles,
  Tablet,
  Trash2,
  Users,
  WandSparkles,
} from '@lucide/vue'
import apiClient from '@/api/client'
import { getProject } from '@/api/projects'
import ProjectStageShell from '@/components/project/ProjectStageShell.vue'
import type { ProjectType } from '@/types/project'

interface RequirementSpecV2 {
  schema_version: 2
  app_type: ProjectType
  goal: string
  target_users: string[]
  primary_scenarios: string[]
  must_have_features: string[]
  should_have_features: string[]
  screens: string[]
  interaction_rules: string[]
  data_and_storage: { persistence: string; backend_required: boolean }
  visual_preferences: { style: string; primary_color: string }
  responsive_targets: string[]
  non_functional_requirements: string[]
  acceptance_criteria: string[]
  out_of_scope: string[]
}

interface AssistQuestion {
  id: string
  question: string
  options: string[]
  blocking: boolean
}

interface AssistResult {
  spec: RequirementSpecV2
  inferred_fields: string[]
  questions: AssistQuestion[]
  ready: boolean
}

const route = useRoute()
const router = useRouter()
const projectId = route.params.projectId as string
const draftKey = `requirements-wizard:${projectId}`
const loading = ref(true)
const saving = ref(false)
const assisting = ref(false)
const error = ref('')
const success = ref('')
const saveState = ref<'idle' | 'saving' | 'saved' | 'error'>('idle')
const currentStep = ref(0)
const maxVisitedStep = ref(0)
const advancedOpen = ref(false)
const showExample = ref(false)
const questions = ref<AssistQuestion[]>([])
const previousAISnapshot = ref<RequirementSpecV2 | null>(null)
const customUser = ref('')
const scenarioInput = ref('')
const featureInput = ref('')
const featurePriority = ref<'must' | 'should' | 'later'>('must')
let saveTimer: ReturnType<typeof setTimeout> | undefined

const form = reactive<RequirementSpecV2>(createEmptySpec())

function createEmptySpec(): RequirementSpecV2 {
  return {
    schema_version: 2,
    app_type: 'utility_app',
    goal: '',
    target_users: [],
    primary_scenarios: [],
    must_have_features: [],
    should_have_features: [],
    screens: [],
    interaction_rules: [],
    data_and_storage: { persistence: '浏览器本地存储', backend_required: false },
    visual_preferences: { style: '简洁现代', primary_color: '蓝色科技' },
    responsive_targets: ['desktop', 'mobile'],
    non_functional_requirements: ['电脑和手机上都能正常使用', '操作后有清晰反馈'],
    acceptance_criteria: [],
    out_of_scope: [],
  }
}

const steps = [
  { key: 'idea', label: '你的想法', short: '想法', icon: Sparkles },
  { key: 'usage', label: '谁会使用', short: '使用方式', icon: Users },
  { key: 'features', label: '需要的功能', short: '功能', icon: ClipboardList },
  { key: 'appearance', label: '外观与设备', short: '外观', icon: Palette },
]

const projectTypes: Array<{ value: ProjectType; label: string; description: string }> = [
  { value: 'interactive_app', label: '互动应用 / 游戏', description: '有操作、状态变化或游戏规则' },
  { value: 'dashboard', label: '管理后台', description: '查看、筛选和管理业务数据' },
  { value: 'data_product', label: '数据看板', description: '用指标和图表展示数据' },
  { value: 'content_site', label: '内容网站', description: '发布、分类和阅读内容' },
  { value: 'ecommerce', label: '电商应用', description: '展示商品并完成购买流程' },
  { value: 'utility_app', label: '工具应用', description: '输入内容并得到处理结果' },
  { value: 'landing_page', label: '介绍页面', description: '集中介绍产品、活动或服务' },
  { value: 'analysis_existing', label: '分析已有项目', description: '检查和改进现有代码或产品' },
]

const userSuggestions = ['普通用户', '学生', '公司员工', '管理员', '客户', '自己使用']
const storageOptions = [
  { value: '不用保存', label: '不用保存', hint: '关闭页面后可以清空' },
  { value: '浏览器本地存储', label: '保存在当前设备', hint: '下次打开还能继续' },
  { value: '登录后保存', label: '登录后保存', hint: '每个人有自己的内容' },
  { value: '跨设备同步', label: '不同设备同步', hint: '手机和电脑看到相同内容' },
  { value: '交给 AI 推荐', label: '交给 AI 推荐', hint: '根据功能自动判断' },
]
const visualStyles = [
  { value: '简洁现代', tone: 'bg-slate-100 border-slate-300' },
  { value: '传统中式', tone: 'bg-red-50 border-red-300' },
  { value: '活泼有趣', tone: 'bg-amber-50 border-amber-300' },
  { value: '科技感', tone: 'bg-blue-50 border-blue-300' },
  { value: '稳重专业', tone: 'bg-zinc-100 border-zinc-400' },
  { value: '交给 AI 推荐', tone: 'bg-violet-50 border-violet-300' },
]
const colorOptions = [
  { value: '蓝色科技', colors: ['#2563eb', '#60a5fa'] },
  { value: '朱红与木色', colors: ['#b42318', '#7a4b2a'] },
  { value: '黑白简洁', colors: ['#111827', '#e5e7eb'] },
  { value: '绿色自然', colors: ['#15803d', '#86efac'] },
  { value: '紫色创意', colors: ['#7c3aed', '#c4b5fd'] },
]

const blockingQuestions = computed(() => questions.value.filter((item) => item.blocking))
const readyItems = computed(() => [
  { label: '已经说明想做什么', done: Boolean(form.goal.trim()), step: 0 },
  { label: '已经确认谁会使用', done: form.target_users.length > 0, step: 1 },
  { label: '已经整理主要功能', done: form.must_have_features.length > 0, step: 2 },
  { label: '已经确认使用设备', done: form.responsive_targets.length > 0, step: 3 },
  { label: '已经有明确的完成效果', done: form.acceptance_criteria.length > 0, step: 3 },
])
const isReady = computed(() => readyItems.value.every((item) => item.done) && blockingQuestions.value.length === 0)
const summary = computed(() => ({
  product: form.goal || '还没有描述产品想法',
  users: form.target_users.length ? form.target_users.join('、') : '还没有选择使用人群',
  features: form.must_have_features.slice(0, 5),
  success: form.acceptance_criteria.slice(0, 4),
}))

function cleanList(value: unknown): string[] {
  return Array.isArray(value)
    ? [...new Set(value.filter((item): item is string => typeof item === 'string').map((item) => item.trim()).filter(Boolean))]
    : []
}

function applySpec(data: Partial<RequirementSpecV2>) {
  form.schema_version = 2
  form.app_type = data.app_type || form.app_type
  form.goal = typeof data.goal === 'string' ? data.goal : form.goal
  form.target_users = cleanList(data.target_users)
  form.primary_scenarios = cleanList(data.primary_scenarios)
  form.must_have_features = cleanList(data.must_have_features)
  form.should_have_features = cleanList(data.should_have_features)
  form.screens = cleanList(data.screens)
  form.interaction_rules = cleanList(data.interaction_rules)
  form.data_and_storage = {
    persistence: data.data_and_storage?.persistence || form.data_and_storage.persistence,
    backend_required: Boolean(data.data_and_storage?.backend_required),
  }
  form.visual_preferences = {
    style: data.visual_preferences?.style || form.visual_preferences.style,
    primary_color: data.visual_preferences?.primary_color || form.visual_preferences.primary_color,
  }
  form.responsive_targets = cleanList(data.responsive_targets).length ? cleanList(data.responsive_targets) : form.responsive_targets
  form.non_functional_requirements = cleanList(data.non_functional_requirements)
  form.acceptance_criteria = cleanList(data.acceptance_criteria)
  form.out_of_scope = cleanList(data.out_of_scope)
}

function snapshotSpec(): RequirementSpecV2 {
  return JSON.parse(JSON.stringify(form)) as RequirementSpecV2
}

async function loadRequirements() {
  loading.value = true
  try {
    const project = await getProject(projectId)
    form.app_type = project.project_type || 'utility_app'
    if (!form.goal && project.description) form.goal = project.description
    try {
      const result: any = await apiClient.get(`/projects/${projectId}/requirements`)
      if (result?.content) applySpec(typeof result.content === 'string' ? JSON.parse(result.content) : result.content)
    } catch {
      // New projects do not have a server requirement version yet.
    }
    const localDraft = localStorage.getItem(draftKey)
    if (localDraft) {
      const parsed = JSON.parse(localDraft)
      if (parsed?.spec) applySpec(parsed.spec)
      currentStep.value = Math.min(3, Math.max(0, Number(parsed?.step) || 0))
      maxVisitedStep.value = Math.max(maxVisitedStep.value, currentStep.value)
      saveState.value = 'saved'
    }
  } catch (err: any) {
    error.value = err.message || '加载项目信息失败'
  } finally {
    loading.value = false
  }
}

function scheduleDraftSave() {
  if (loading.value) return
  saveState.value = 'saving'
  if (saveTimer) clearTimeout(saveTimer)
  saveTimer = setTimeout(() => {
    try {
      localStorage.setItem(draftKey, JSON.stringify({ spec: snapshotSpec(), step: currentStep.value, savedAt: Date.now() }))
      saveState.value = 'saved'
    } catch {
      saveState.value = 'error'
    }
  }, 600)
}

watch(form, scheduleDraftSave, { deep: true })
watch(currentStep, scheduleDraftSave)

function validateStep(step = currentStep.value): boolean {
  error.value = ''
  if (step === 0 && !form.goal.trim()) error.value = '先用一句话描述你想做的产品，就可以继续了。'
  if (step === 1 && !form.target_users.length) error.value = '请选择至少一种使用人群。'
  if (step === 2 && !form.must_have_features.length) error.value = '请保留至少一个“必须有”的功能。'
  if (step === 3 && !form.responsive_targets.length) error.value = '请选择用户主要使用的设备。'
  return !error.value
}

async function nextStep() {
  if (!validateStep()) return
  if (currentStep.value < 3) {
    currentStep.value += 1
    maxVisitedStep.value = Math.max(maxVisitedStep.value, currentStep.value)
    await nextTick()
    document.querySelector('[data-wizard-panel]')?.scrollIntoView({ behavior: 'smooth', block: 'start' })
  }
}

function goToStep(step: number) {
  if (step <= maxVisitedStep.value) {
    currentStep.value = step
    error.value = ''
  }
}

function toggleUser(value: string) {
  const index = form.target_users.indexOf(value)
  if (index >= 0) form.target_users.splice(index, 1)
  else form.target_users.push(value)
}

function addCustomUser() {
  const value = customUser.value.trim()
  if (value && !form.target_users.includes(value)) form.target_users.push(value)
  customUser.value = ''
}

function addScenario() {
  const value = scenarioInput.value.trim()
  if (value && !form.primary_scenarios.includes(value)) form.primary_scenarios.push(value)
  scenarioInput.value = ''
}

function featureList(priority = featurePriority.value) {
  if (priority === 'should') return form.should_have_features
  if (priority === 'later') return form.out_of_scope
  return form.must_have_features
}

function addFeature() {
  const value = featureInput.value.trim()
  if (value && !featureList().includes(value)) featureList().push(value)
  featureInput.value = ''
}

function moveFeature(value: string, from: 'must' | 'should' | 'later', to: 'must' | 'should' | 'later') {
  const source = featureList(from)
  const index = source.indexOf(value)
  if (index >= 0) source.splice(index, 1)
  if (!featureList(to).includes(value)) featureList(to).push(value)
}

function setDevices(mode: string) {
  if (mode === 'all') form.responsive_targets = ['desktop', 'tablet', 'mobile']
  else form.responsive_targets = [mode]
}

function generateAcceptanceCriteria() {
  if (!form.acceptance_criteria.length) {
    form.acceptance_criteria = form.must_have_features.map((feature) => `用户可以使用“${feature}”并看到清晰、正确的结果`)
  }
}

function confirmedFields(): string[] {
  const result = ['goal', 'app_type']
  if (form.target_users.length) result.push('target_users')
  if (form.primary_scenarios.length) result.push('primary_scenarios')
  if (form.must_have_features.length) result.push('must_have_features')
  if (form.should_have_features.length) result.push('should_have_features')
  if (form.visual_preferences.style || form.visual_preferences.primary_color) result.push('visual_preferences')
  if (form.responsive_targets.length) result.push('responsive_targets')
  return result
}

async function assistRequirements() {
  if (!form.goal.trim()) {
    error.value = '先简单描述你的想法，AI 才能帮你整理。'
    return
  }
  assisting.value = true
  error.value = ''
  success.value = ''
  previousAISnapshot.value = snapshotSpec()
  try {
    const result = await apiClient.post<any, AssistResult>(`/projects/${projectId}/requirements/assist`, {
      description: form.goal,
      current_spec: snapshotSpec(),
      current_step: steps[currentStep.value].key,
      confirmed_fields: confirmedFields(),
    })
    applySpec(result.spec)
    questions.value = result.questions || []
    maxVisitedStep.value = Math.max(maxVisitedStep.value, 2)
    success.value = result.ready ? 'AI 已经整理好主要内容，你可以继续确认。' : 'AI 已整理需求，还有少量问题需要你确认。'
  } catch (err: any) {
    error.value = err.message || 'AI 暂时没有整理成功，你可以继续手动填写或稍后重试。'
  } finally {
    assisting.value = false
  }
}

function undoAssist() {
  if (!previousAISnapshot.value) return
  applySpec(previousAISnapshot.value)
  previousAISnapshot.value = null
  questions.value = []
  success.value = '已经撤销最近一次 AI 整理。'
}

function answerQuestion(question: AssistQuestion, option: string) {
  const rule = `${question.question}：${option}`
  if (!form.interaction_rules.includes(rule)) form.interaction_rules.push(rule)
  questions.value = questions.value.filter((item) => item.id !== question.id)
}

function setAdvancedList(field: 'screens' | 'interaction_rules' | 'acceptance_criteria' | 'non_functional_requirements', event: Event) {
  const value = (event.target as HTMLTextAreaElement).value
  form[field] = value.split('\n').map((item) => item.trim()).filter(Boolean)
}

function saveLocalNow() {
  if (saveTimer) clearTimeout(saveTimer)
  localStorage.setItem(draftKey, JSON.stringify({ spec: snapshotSpec(), step: currentStep.value, savedAt: Date.now() }))
  saveState.value = 'saved'
  success.value = '进度已保存在当前设备。'
}

async function saveAndContinue() {
  generateAcceptanceCriteria()
  if (!validateStep(0) || !form.target_users.length || !form.must_have_features.length || !form.acceptance_criteria.length) {
    error.value = '还需要补充使用人群、必须功能和完成效果，才能查看 AI 方案。'
    return
  }
  if (blockingQuestions.value.length) {
    error.value = `还有 ${blockingQuestions.value.length} 个关键问题需要确认。`
    return
  }
  saving.value = true
  error.value = ''
  try {
    await apiClient.put(`/projects/${projectId}/requirements`, { content: JSON.stringify(snapshotSpec()) })
    localStorage.removeItem(draftKey)
    await router.push(`/projects/${projectId}/blueprint`)
  } catch (err: any) {
    error.value = err.message || '保存需求失败，请稍后重试。'
  } finally {
    saving.value = false
  }
}

onMounted(loadRequirements)
onBeforeUnmount(() => { if (saveTimer) clearTimeout(saveTimer) })
</script>

<template>
  <ProjectStageShell
    :icon="WandSparkles"
    title="告诉 AI 你想做什么"
    description="不需要懂产品或技术。用自己的话描述想法，AI 会帮你整理成功能方案。"
    :step-text="`第 ${currentStep + 1} 步，共 4 步`"
  >
    <div v-if="loading" class="grid gap-6 lg:grid-cols-[minmax(0,1fr)_280px]">
      <div class="h-[560px] animate-pulse rounded-xl border border-border bg-surface-muted" />
      <div class="h-72 animate-pulse rounded-xl border border-border bg-surface-muted" />
    </div>

    <div v-else class="space-y-5">
      <div class="flex flex-wrap items-center justify-between gap-3 rounded-xl border border-border bg-surface px-4 py-3">
        <nav aria-label="需求填写步骤" class="min-w-0 flex-1 overflow-x-auto">
          <ol class="flex min-w-max items-center gap-1 sm:gap-2">
            <li v-for="(step, index) in steps" :key="step.key" class="flex items-center">
              <button
                type="button"
                class="group flex min-h-11 cursor-pointer items-center gap-2 rounded-lg px-2.5 text-sm transition-colors duration-200 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent disabled:cursor-not-allowed disabled:opacity-45 sm:px-3"
                :class="index === currentStep ? 'bg-accent-soft text-accent' : index < currentStep ? 'text-success hover:bg-success/10' : 'text-text-muted hover:bg-surface-muted'"
                :disabled="index > maxVisitedStep"
                :aria-current="index === currentStep ? 'step' : undefined"
                @click="goToStep(index)"
              >
                <span class="grid size-7 place-items-center rounded-full border" :class="index === currentStep ? 'border-accent bg-accent text-white' : index < currentStep ? 'border-success bg-success text-white' : 'border-border bg-surface'">
                  <Check v-if="index < currentStep" :size="15" aria-hidden="true" />
                  <span v-else>{{ index + 1 }}</span>
                </span>
                <span class="hidden sm:inline">{{ step.short }}</span>
              </button>
              <span v-if="index < steps.length - 1" class="h-px w-3 bg-border sm:w-5" aria-hidden="true" />
            </li>
          </ol>
        </nav>
        <div class="flex items-center gap-2 text-xs" aria-live="polite">
          <Loader2 v-if="saveState === 'saving'" :size="14" class="animate-spin text-text-muted" />
          <CheckCircle2 v-else-if="saveState === 'saved'" :size="14" class="text-success" />
          <AlertCircle v-else-if="saveState === 'error'" :size="14" class="text-danger" />
          <span :class="saveState === 'error' ? 'text-danger' : 'text-text-muted'">
            {{ saveState === 'saving' ? '正在保存' : saveState === 'saved' ? '已自动保存' : saveState === 'error' ? '自动保存失败' : '修改会自动保存' }}
          </span>
        </div>
      </div>

      <div class="grid items-start gap-6 lg:grid-cols-[minmax(0,1fr)_300px]">
        <main data-wizard-panel class="overflow-hidden rounded-xl border border-border bg-surface shadow-sm">
          <section v-if="currentStep === 0" aria-labelledby="step-idea-title" class="space-y-6 p-5 sm:p-7">
            <header class="max-w-2xl">
              <p class="text-sm font-semibold text-accent">第一步</p>
              <h2 id="step-idea-title" class="mt-1 text-xl font-bold text-text-primary">简单说说你想做什么</h2>
              <p class="mt-2 text-sm leading-6 text-text-secondary">像和朋友聊天一样描述即可，不需要整理成专业文档。</p>
            </header>

            <label class="block space-y-2">
              <span class="text-sm font-semibold text-text-primary">你的产品想法</span>
              <textarea
                v-model.trim="form.goal"
                rows="7"
                class="w-full resize-y rounded-xl border border-border bg-surface-muted px-4 py-3 text-base leading-7 text-text-primary outline-none transition-colors focus:border-accent focus:ring-2 focus:ring-accent/20"
                placeholder="例如：我想做一个可以和电脑下中国象棋的网站，打开后就能开始下棋，可以选择难度，也可以重新开始。"
              />
            </label>

            <div class="flex flex-wrap gap-2">
              <button type="button" class="inline-flex min-h-11 cursor-pointer items-center gap-2 rounded-lg bg-accent px-4 text-sm font-semibold text-white transition-colors hover:bg-accent/85 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-60" :disabled="assisting" @click="assistRequirements">
                <Loader2 v-if="assisting" :size="17" class="animate-spin" /><Sparkles v-else :size="17" />
                {{ assisting ? 'AI 正在理解你的想法' : '让 AI 帮我整理' }}
              </button>
              <button type="button" class="inline-flex min-h-11 cursor-pointer items-center gap-2 rounded-lg border border-border px-4 text-sm font-medium text-text-secondary transition-colors hover:bg-surface-muted focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent" @click="showExample = !showExample">
                <HelpCircle :size="17" />{{ showExample ? '收起示例' : '不知道怎么写？' }}
              </button>
              <button v-if="previousAISnapshot" type="button" class="inline-flex min-h-11 cursor-pointer items-center gap-2 rounded-lg px-3 text-sm font-medium text-text-secondary transition-colors hover:bg-surface-muted focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent" @click="undoAssist">
                <RotateCcw :size="16" />撤销 AI 整理
              </button>
            </div>

            <div v-if="showExample" class="rounded-xl border border-accent/20 bg-accent-soft/40 p-4 text-sm leading-6 text-text-secondary">
              <p class="font-semibold text-text-primary">你可以这样写</p>
              <p class="mt-1">“我想做一个给销售团队使用的客户管理页面，可以查看客户列表、搜索客户、记录跟进情况，也希望手机上能用。”</p>
            </div>

            <fieldset class="space-y-3">
              <legend class="text-sm font-semibold text-text-primary">这更像哪一种产品？</legend>
              <p class="text-sm text-text-secondary">如果不确定，可以先使用 AI 的判断，之后仍然可以修改。</p>
              <div class="grid gap-2 sm:grid-cols-2">
                <label v-for="type in projectTypes" :key="type.value" class="flex min-h-20 cursor-pointer items-start gap-3 rounded-xl border p-3 transition-colors duration-200 focus-within:ring-2 focus-within:ring-accent" :class="form.app_type === type.value ? 'border-accent bg-accent-soft/50' : 'border-border hover:bg-surface-muted'">
                  <input v-model="form.app_type" type="radio" :value="type.value" class="mt-1 accent-[var(--color-accent)]" />
                  <span><strong class="block text-sm text-text-primary">{{ type.label }}</strong><span class="mt-1 block text-xs leading-5 text-text-secondary">{{ type.description }}</span></span>
                </label>
              </div>
            </fieldset>
          </section>

          <section v-else-if="currentStep === 1" aria-labelledby="step-usage-title" class="space-y-7 p-5 sm:p-7">
            <header>
              <p class="text-sm font-semibold text-accent">第二步</p>
              <h2 id="step-usage-title" class="mt-1 text-xl font-bold text-text-primary">谁会使用？通常怎么使用？</h2>
              <p class="mt-2 text-sm leading-6 text-text-secondary">选择最接近的答案即可，也可以添加自己的描述。</p>
            </header>

            <fieldset class="space-y-3">
              <legend class="text-sm font-semibold text-text-primary">谁会使用这个产品？</legend>
              <div class="flex flex-wrap gap-2">
                <button v-for="user in userSuggestions" :key="user" type="button" class="inline-flex min-h-11 cursor-pointer items-center gap-2 rounded-full border px-4 text-sm transition-colors duration-200 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent" :class="form.target_users.includes(user) ? 'border-accent bg-accent text-white' : 'border-border bg-surface hover:bg-surface-muted text-text-secondary'" :aria-pressed="form.target_users.includes(user)" @click="toggleUser(user)">
                  <Check v-if="form.target_users.includes(user)" :size="15" />{{ user }}
                </button>
              </div>
              <div class="flex max-w-lg gap-2">
                <label class="sr-only" for="custom-user">添加其他使用人群</label>
                <input id="custom-user" v-model="customUser" class="min-h-11 min-w-0 flex-1 rounded-lg border border-border bg-surface-muted px-3 text-sm outline-none focus:border-accent focus:ring-2 focus:ring-accent/20" placeholder="添加其他使用人群" @keydown.enter.prevent="addCustomUser" />
                <button type="button" class="inline-flex min-h-11 cursor-pointer items-center gap-1 rounded-lg border border-border px-3 text-sm font-medium hover:bg-surface-muted focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent" @click="addCustomUser"><Plus :size="16" />添加</button>
              </div>
            </fieldset>

            <fieldset class="space-y-3">
              <legend class="text-sm font-semibold text-text-primary">用户打开后最想完成什么？</legend>
              <ul v-if="form.primary_scenarios.length" class="space-y-2">
                <li v-for="(scenario, index) in form.primary_scenarios" :key="`${scenario}-${index}`" class="flex items-center gap-2 rounded-lg border border-border bg-surface-muted p-2">
                  <CheckCircle2 :size="17" class="shrink-0 text-success" />
                  <input v-model="form.primary_scenarios[index]" :aria-label="`使用方式 ${index + 1}`" class="min-h-9 min-w-0 flex-1 bg-transparent text-sm text-text-primary outline-none focus:ring-2 focus:ring-accent/20" />
                  <button type="button" class="grid size-9 shrink-0 cursor-pointer place-items-center rounded-md text-text-muted hover:bg-danger/10 hover:text-danger focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent" :aria-label="`删除使用方式：${scenario}`" @click="form.primary_scenarios.splice(index, 1)"><Trash2 :size="16" /></button>
                </li>
              </ul>
              <div class="flex max-w-2xl gap-2">
                <label class="sr-only" for="scenario-input">添加一种使用方式</label>
                <input id="scenario-input" v-model="scenarioInput" class="min-h-11 min-w-0 flex-1 rounded-lg border border-border bg-surface-muted px-3 text-sm outline-none focus:border-accent focus:ring-2 focus:ring-accent/20" placeholder="例如：打开页面后立即开始人机对战" @keydown.enter.prevent="addScenario" />
                <button type="button" class="inline-flex min-h-11 cursor-pointer items-center gap-1 rounded-lg border border-border px-3 text-sm font-medium hover:bg-surface-muted focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent" @click="addScenario"><Plus :size="16" />添加</button>
              </div>
            </fieldset>

            <fieldset class="space-y-3">
              <legend class="text-sm font-semibold text-text-primary">使用过程中需要保存内容吗？</legend>
              <div class="grid gap-2 sm:grid-cols-2">
                <label v-for="option in storageOptions" :key="option.value" class="flex min-h-16 cursor-pointer gap-3 rounded-xl border p-3 transition-colors focus-within:ring-2 focus-within:ring-accent" :class="form.data_and_storage.persistence === option.value ? 'border-accent bg-accent-soft/50' : 'border-border hover:bg-surface-muted'">
                  <input v-model="form.data_and_storage.persistence" type="radio" :value="option.value" class="mt-1 accent-[var(--color-accent)]" />
                  <span><strong class="block text-sm text-text-primary">{{ option.label }}</strong><span class="mt-0.5 block text-xs text-text-secondary">{{ option.hint }}</span></span>
                </label>
              </div>
            </fieldset>
          </section>

          <section v-else-if="currentStep === 2" aria-labelledby="step-features-title" class="space-y-7 p-5 sm:p-7">
            <header>
              <p class="text-sm font-semibold text-accent">第三步</p>
              <h2 id="step-features-title" class="mt-1 text-xl font-bold text-text-primary">确认这次需要完成的功能</h2>
              <p class="mt-2 text-sm leading-6 text-text-secondary">可以直接修改功能名称，或把功能移动到更合适的分组。</p>
            </header>

            <div class="grid gap-4 xl:grid-cols-3">
              <section v-for="group in [{ key: 'must', title: '必须有', note: '没有这些就不能使用', list: form.must_have_features }, { key: 'should', title: '最好有', note: '有了体验会更完整', list: form.should_have_features }, { key: 'later', title: '以后再做', note: '本次明确不生成', list: form.out_of_scope }]" :key="group.key" class="rounded-xl border border-border bg-surface-muted/50 p-3">
                <div class="mb-3">
                  <h3 class="font-semibold text-text-primary">{{ group.title }} <span class="text-xs font-normal text-text-muted">{{ group.list.length }}</span></h3>
                  <p class="mt-0.5 text-xs text-text-secondary">{{ group.note }}</p>
                </div>
                <div class="space-y-2">
                  <article v-for="(feature, index) in group.list" :key="`${group.key}-${index}`" class="rounded-lg border border-border bg-surface p-2.5">
                    <div class="flex items-start gap-2">
                      <input v-model="group.list[index]" :aria-label="`${group.title}功能 ${index + 1}`" class="min-h-9 min-w-0 flex-1 rounded-md bg-transparent px-1 text-sm font-medium text-text-primary outline-none focus:ring-2 focus:ring-accent/20" />
                      <button type="button" class="grid size-9 shrink-0 cursor-pointer place-items-center rounded-md text-text-muted hover:bg-danger/10 hover:text-danger focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent" :aria-label="`删除功能：${feature}`" @click="group.list.splice(index, 1)"><Trash2 :size="15" /></button>
                    </div>
                    <label class="mt-2 block text-xs text-text-muted">
                      优先级
                      <select :value="group.key" class="ml-1 rounded-md border border-border bg-surface px-2 py-1 text-xs text-text-secondary outline-none focus:ring-2 focus:ring-accent" @change="moveFeature(feature, group.key as any, ($event.target as HTMLSelectElement).value as any)">
                        <option value="must">必须有</option><option value="should">最好有</option><option value="later">以后再做</option>
                      </select>
                    </label>
                  </article>
                  <p v-if="!group.list.length" class="rounded-lg border border-dashed border-border p-3 text-center text-xs text-text-muted">这里还没有功能</p>
                </div>
              </section>
            </div>

            <div class="rounded-xl border border-accent/20 bg-accent-soft/30 p-4">
              <label for="feature-input" class="text-sm font-semibold text-text-primary">添加一个功能</label>
              <div class="mt-2 flex flex-col gap-2 sm:flex-row">
                <input id="feature-input" v-model="featureInput" class="min-h-11 min-w-0 flex-1 rounded-lg border border-border bg-surface px-3 text-sm outline-none focus:border-accent focus:ring-2 focus:ring-accent/20" placeholder="例如：重新开始游戏" @keydown.enter.prevent="addFeature" />
                <select v-model="featurePriority" aria-label="新功能优先级" class="min-h-11 rounded-lg border border-border bg-surface px-3 text-sm text-text-secondary outline-none focus:ring-2 focus:ring-accent"><option value="must">必须有</option><option value="should">最好有</option><option value="later">以后再做</option></select>
                <button type="button" class="inline-flex min-h-11 cursor-pointer items-center justify-center gap-1 rounded-lg bg-accent px-4 text-sm font-semibold text-white hover:bg-accent/85 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent" @click="addFeature"><Plus :size="16" />添加</button>
              </div>
            </div>
          </section>

          <section v-else aria-labelledby="step-appearance-title" class="space-y-7 p-5 sm:p-7">
            <header>
              <p class="text-sm font-semibold text-accent">第四步</p>
              <h2 id="step-appearance-title" class="mt-1 text-xl font-bold text-text-primary">选择喜欢的外观和使用设备</h2>
              <p class="mt-2 text-sm leading-6 text-text-secondary">只需选择大致方向，具体设计会在方案中继续确认。</p>
            </header>

            <fieldset class="space-y-3">
              <legend class="text-sm font-semibold text-text-primary">你喜欢哪种感觉？</legend>
              <div class="grid grid-cols-2 gap-2 sm:grid-cols-3">
                <label v-for="style in visualStyles" :key="style.value" class="min-h-24 cursor-pointer rounded-xl border p-3 transition-colors focus-within:ring-2 focus-within:ring-accent" :class="form.visual_preferences.style === style.value ? 'border-accent ring-1 ring-accent' : `${style.tone} hover:border-accent/60`">
                  <input v-model="form.visual_preferences.style" type="radio" :value="style.value" class="sr-only" />
                  <span class="flex items-center justify-between"><Palette :size="19" class="text-text-secondary" /><CheckCircle2 v-if="form.visual_preferences.style === style.value" :size="18" class="text-accent" /></span>
                  <strong class="mt-4 block text-sm text-text-primary">{{ style.value }}</strong>
                </label>
              </div>
            </fieldset>

            <fieldset class="space-y-3">
              <legend class="text-sm font-semibold text-text-primary">喜欢什么颜色？</legend>
              <div class="flex flex-wrap gap-2">
                <label v-for="color in colorOptions" :key="color.value" class="flex min-h-12 cursor-pointer items-center gap-2 rounded-full border px-3 transition-colors focus-within:ring-2 focus-within:ring-accent" :class="form.visual_preferences.primary_color === color.value ? 'border-accent bg-accent-soft' : 'border-border hover:bg-surface-muted'">
                  <input v-model="form.visual_preferences.primary_color" type="radio" :value="color.value" class="sr-only" />
                  <span class="flex -space-x-1" aria-hidden="true"><span v-for="tone in color.colors" :key="tone" class="size-5 rounded-full border-2 border-white" :style="{ backgroundColor: tone }" /></span>
                  <span class="text-sm text-text-primary">{{ color.value }}</span>
                </label>
              </div>
            </fieldset>

            <fieldset class="space-y-3">
              <legend class="text-sm font-semibold text-text-primary">用户主要会在哪里使用？</legend>
              <div class="grid grid-cols-2 gap-2 sm:grid-cols-4">
                <button v-for="device in [{ value: 'desktop', label: '电脑', icon: Monitor }, { value: 'mobile', label: '手机', icon: Smartphone }, { value: 'tablet', label: '平板', icon: Tablet }, { value: 'all', label: '都需要', icon: Laptop }]" :key="device.value" type="button" class="flex min-h-24 cursor-pointer flex-col items-center justify-center gap-2 rounded-xl border text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent" :class="device.value === 'all' ? (form.responsive_targets.length === 3 ? 'border-accent bg-accent-soft text-accent' : 'border-border hover:bg-surface-muted text-text-secondary') : (form.responsive_targets.length === 1 && form.responsive_targets[0] === device.value ? 'border-accent bg-accent-soft text-accent' : 'border-border hover:bg-surface-muted text-text-secondary')" @click="setDevices(device.value)"><component :is="device.icon" :size="23" />{{ device.label }}</button>
              </div>
            </fieldset>

            <section class="rounded-xl border border-success/25 bg-success/10 p-4">
              <div class="flex items-start gap-3">
                <CheckCircle2 :size="20" class="mt-0.5 shrink-0 text-success" />
                <div><h3 class="font-semibold text-text-primary">AI 会自动整理完成效果</h3><p class="mt-1 text-sm leading-6 text-text-secondary">系统会根据“必须有”的功能生成可以验证的完成标准，你也可以在高级设置中修改。</p><button type="button" class="mt-2 cursor-pointer text-sm font-semibold text-accent hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent" @click="generateAcceptanceCriteria">现在生成完成效果</button></div>
              </div>
            </section>
          </section>

          <div v-if="error" role="alert" class="mx-5 mb-4 flex items-start gap-2 rounded-lg border border-danger/25 bg-danger/10 px-4 py-3 text-sm text-danger sm:mx-7"><AlertCircle :size="18" class="mt-0.5 shrink-0" />{{ error }}</div>
          <div v-if="success" aria-live="polite" class="mx-5 mb-4 flex items-start gap-2 rounded-lg border border-success/25 bg-success/10 px-4 py-3 text-sm text-success sm:mx-7"><CheckCircle2 :size="18" class="mt-0.5 shrink-0" />{{ success }}</div>

          <footer class="flex flex-col-reverse gap-3 border-t border-border bg-surface-muted/40 p-4 sm:flex-row sm:items-center sm:justify-between sm:px-7">
            <button v-if="currentStep > 0" type="button" class="inline-flex min-h-11 cursor-pointer items-center justify-center gap-2 rounded-lg border border-border bg-surface px-4 text-sm font-medium text-text-primary hover:bg-surface-muted focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent" @click="currentStep--"><ArrowLeft :size="17" />上一步</button><span v-else />
            <div class="flex flex-col gap-2 sm:flex-row">
              <button type="button" class="inline-flex min-h-11 cursor-pointer items-center justify-center gap-2 rounded-lg border border-border bg-surface px-4 text-sm font-medium text-text-primary hover:bg-surface-muted focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent" @click="saveLocalNow"><Save :size="16" />暂时保存</button>
              <button v-if="currentStep < 3" type="button" class="inline-flex min-h-11 cursor-pointer items-center justify-center gap-2 rounded-lg bg-accent px-5 text-sm font-semibold text-white hover:bg-accent/85 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent focus-visible:ring-offset-2" @click="nextStep">继续<ArrowRight :size="17" /></button>
              <button v-else type="button" :disabled="saving" class="inline-flex min-h-11 cursor-pointer items-center justify-center gap-2 rounded-lg bg-accent px-5 text-sm font-semibold text-white hover:bg-accent/85 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-60" @click="saveAndContinue"><Loader2 v-if="saving" :size="17" class="animate-spin" /><Sparkles v-else :size="17" />整理需求并查看方案</button>
            </div>
          </footer>
        </main>

        <aside class="space-y-4 lg:sticky lg:top-6">
          <section class="rounded-xl border border-border bg-surface p-5 shadow-sm">
            <div class="flex items-center justify-between gap-3"><h2 class="font-semibold text-text-primary">需求准备情况</h2><span class="rounded-full px-2.5 py-1 text-xs font-semibold" :class="isReady ? 'bg-success/10 text-success' : 'bg-warning/10 text-warning'">{{ isReady ? '可以查看方案' : '继续完善' }}</span></div>
            <ul class="mt-4 space-y-3">
              <li v-for="item in readyItems" :key="item.label" class="flex items-center gap-2 text-sm" :class="item.done ? 'text-text-primary' : 'text-text-muted'"><span class="grid size-5 shrink-0 place-items-center rounded-full" :class="item.done ? 'bg-success text-white' : 'border border-border bg-surface-muted'"><Check v-if="item.done" :size="13" /></span><span>{{ item.label }}</span><button v-if="!item.done && item.step <= maxVisitedStep" type="button" class="ml-auto cursor-pointer text-xs font-medium text-accent hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent" @click="goToStep(item.step)">去完成</button></li>
            </ul>
            <p class="mt-4 rounded-lg bg-surface-muted p-3 text-xs leading-5 text-text-secondary">{{ isReady ? '已经可以整理方案，其他细节之后仍然可以修改。' : '只需要完成关键内容，不必把所有高级设置都填满。' }}</p>
          </section>

          <section v-if="questions.length" class="rounded-xl border border-warning/30 bg-warning/10 p-4">
            <div class="flex items-center gap-2 text-text-primary"><HelpCircle :size="18" class="text-warning" /><h2 class="font-semibold">AI 想再确认</h2></div>
            <article v-for="question in questions" :key="question.id" class="mt-3 border-t border-warning/20 pt-3 first:border-0 first:pt-0">
              <p class="text-sm font-medium leading-6 text-text-primary">{{ question.question }}</p>
              <div class="mt-2 space-y-1.5"><button v-for="option in question.options" :key="option" type="button" class="w-full cursor-pointer rounded-lg border border-border bg-surface px-3 py-2 text-left text-xs leading-5 text-text-secondary transition-colors hover:border-accent hover:bg-accent-soft focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent" @click="answerQuestion(question, option)">{{ option }}</button></div>
            </article>
          </section>

          <section class="rounded-xl border border-accent/20 bg-accent-soft/35 p-4">
            <div class="flex items-center gap-2"><Sparkles :size="18" class="text-accent" /><h2 class="font-semibold text-text-primary">AI 整理摘要</h2></div>
            <dl class="mt-3 space-y-3 text-sm"><div><dt class="text-xs font-medium text-text-muted">你要做的是</dt><dd class="mt-1 line-clamp-3 leading-5 text-text-primary">{{ summary.product }}</dd></div><div><dt class="text-xs font-medium text-text-muted">主要用户</dt><dd class="mt-1 text-text-primary">{{ summary.users }}</dd></div><div v-if="summary.features.length"><dt class="text-xs font-medium text-text-muted">核心功能</dt><dd class="mt-1 text-text-primary">{{ summary.features.join('、') }}</dd></div><div v-if="summary.success.length"><dt class="text-xs font-medium text-text-muted">完成效果</dt><dd class="mt-1 line-clamp-3 text-text-primary">{{ summary.success.join('；') }}</dd></div></dl>
          </section>
        </aside>
      </div>

      <details class="group rounded-xl border border-border bg-surface" :open="advancedOpen" @toggle="advancedOpen = ($event.target as HTMLDetailsElement).open">
        <summary class="flex min-h-14 cursor-pointer list-none items-center justify-between gap-3 px-5 font-semibold text-text-primary focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-accent"><span>高级设置：规则、数据和完成要求</span><span class="flex items-center gap-2 text-xs font-normal text-text-muted">懂产品设计时再填写<ChevronDown :size="18" class="transition-transform duration-200 group-open:rotate-180" /></span></summary>
        <div class="grid gap-4 border-t border-border p-5 md:grid-cols-2">
          <label class="space-y-1.5 text-sm font-medium text-text-primary">页面或主要区域<textarea :value="form.screens.join('\n')" rows="4" class="w-full rounded-lg border border-border bg-surface-muted px-3 py-2 font-normal outline-none focus:ring-2 focus:ring-accent" placeholder="例如：游戏主界面" @input="setAdvancedList('screens', $event)" /></label>
          <label class="space-y-1.5 text-sm font-medium text-text-primary">操作和业务规则<textarea :value="form.interaction_rules.join('\n')" rows="4" class="w-full rounded-lg border border-border bg-surface-muted px-3 py-2 font-normal outline-none focus:ring-2 focus:ring-accent" placeholder="例如：红方先行" @input="setAdvancedList('interaction_rules', $event)" /></label>
          <label class="space-y-1.5 text-sm font-medium text-text-primary">完成后应该能做到什么<textarea :value="form.acceptance_criteria.join('\n')" rows="5" class="w-full rounded-lg border border-border bg-surface-muted px-3 py-2 font-normal outline-none focus:ring-2 focus:ring-accent" placeholder="系统会根据必须功能自动生成" @input="setAdvancedList('acceptance_criteria', $event)" /></label>
          <label class="space-y-1.5 text-sm font-medium text-text-primary">质量和体验要求<textarea :value="form.non_functional_requirements.join('\n')" rows="5" class="w-full rounded-lg border border-border bg-surface-muted px-3 py-2 font-normal outline-none focus:ring-2 focus:ring-accent" placeholder="例如：手机上不能横向滚动" @input="setAdvancedList('non_functional_requirements', $event)" /></label>
          <label class="flex min-h-12 cursor-pointer items-center gap-3 rounded-lg border border-border bg-surface-muted px-4 text-sm text-text-primary"><input v-model="form.data_and_storage.backend_required" type="checkbox" class="accent-[var(--color-accent)]" />需要账号、在线保存或跨设备同步</label>
        </div>
      </details>
    </div>
  </ProjectStageShell>
</template>

<style scoped>
@media (prefers-reduced-motion: reduce) {
  *, *::before, *::after { scroll-behavior: auto !important; transition-duration: 0.01ms !important; animation-duration: 0.01ms !important; }
}
</style>
