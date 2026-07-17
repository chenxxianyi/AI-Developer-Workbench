<script setup lang="ts">
import { computed, reactive, ref } from 'vue'
import { RouterLink } from 'vue-router'
import {
  ArrowLeft,
  ArrowRight,
  CheckCircle2,
  FileText,
  FolderPlus,
  Gamepad2,
  Gauge,
  LayoutDashboard,
  Loader2,
  Monitor,
  Newspaper,
  Search,
  ShoppingCart,
  Wrench,
} from '@lucide/vue'
import { createProject } from '@/api/projects'
import ProjectStageShell from '@/components/project/ProjectStageShell.vue'
import ToolFormSection from '@/components/tool/ToolFormSection.vue'
import type { Project, ProjectType } from '@/types/project'

const step = ref(1)
const selectedType = ref<ProjectType>('interactive_app')
const creating = ref(false)
const error = ref('')
const createdProject = ref<Project | null>(null)
const form = reactive({
  name: '',
  description: '',
  techStacks: ['Vue 3', 'TypeScript', 'Tailwind CSS'] as string[],
})

const projectTypes = [
  {
    value: 'interactive_app' as const,
    icon: Gamepad2,
    title: '互动应用 / 游戏',
    description: '适合游戏、编辑器和强交互工具，重点生成状态与业务规则。',
  },
  {
    value: 'dashboard' as const,
    icon: LayoutDashboard,
    title: '管理后台',
    description: '适合运营和业务系统，重点生成导航、列表、表单与数据状态。',
  },
  {
    value: 'data_product' as const,
    icon: Gauge,
    title: '数据看板',
    description: '适合统计和分析产品，重点生成指标、筛选、图表和空状态。',
  },
  {
    value: 'content_site' as const,
    icon: Newspaper,
    title: '内容网站',
    description: '适合官网、资讯、博客和文档，重点生成信息架构和内容体验。',
  },
  {
    value: 'ecommerce' as const,
    icon: ShoppingCart,
    title: '电商应用',
    description: '适合商品和订单流程，重点生成商品浏览、购物车和结算状态。',
  },
  {
    value: 'utility_app' as const,
    icon: Wrench,
    title: '工具型应用',
    description: '适合计算、转换和效率工具，重点生成输入、处理和结果逻辑。',
  },
  {
    value: 'landing_page' as const,
    icon: Monitor,
    title: '营销落地页',
    description: '适合活动和产品宣传，仅此类型默认使用转化型页面结构。',
  },
  {
    value: 'analysis_existing' as const,
    icon: Search,
    title: '分析已有项目',
    description: '创建项目档案后，关联项目文件并运行质量诊断。',
  },
]

const techStacks = ['Vue 3', 'React', 'Tailwind CSS', 'TypeScript', 'Node.js', 'Go']
const frontendStacks = new Set(['Vue 3', 'React', 'Tailwind CSS', 'TypeScript'])
const backendStacks = new Set(['Node.js', 'Go'])

const stepText = computed(() => `第 ${step.value} 步，共 3 步`)
const primaryNextRoute = computed(() => {
  if (!createdProject.value) return '/projects'
  return selectedType.value === 'analysis_existing'
    ? `/tools/project-doctor?project_id=${createdProject.value.id}`
    : `/projects/${createdProject.value.id}/requirements`
})

async function submitProject() {
  if (!form.name.trim() || creating.value) return

  creating.value = true
  error.value = ''
  try {
    createdProject.value = await createProject({
      name: form.name.trim(),
      project_type: selectedType.value,
      description: form.description.trim() || undefined,
      frontend_stack: form.techStacks.filter((stack) => frontendStacks.has(stack)).join(' + ') || undefined,
      backend_stack: form.techStacks.filter((stack) => backendStacks.has(stack)).join(' + ') || undefined,
    })
    step.value = 3
  } catch (err: any) {
    error.value = err.message || '项目创建失败，请稍后重试'
  } finally {
    creating.value = false
  }
}
</script>

<template>
  <div class="mx-auto max-w-6xl">
    <ProjectStageShell
      :icon="FolderPlus"
      title="创建新项目"
      description="建立项目档案，并选择接下来要进入的开发流程。"
      :step-text="stepText"
    >
      <div class="grid grid-cols-1 gap-6 lg:grid-cols-[minmax(0,1fr)_320px]">
        <section class="rounded-lg border border-border bg-surface p-5 shadow-sm sm:p-6">
          <template v-if="step === 1">
            <div class="mb-5">
              <h2 class="text-lg font-semibold text-text-primary">选择项目类型</h2>
              <p class="mt-1 text-sm text-text-secondary">项目类型决定创建完成后的下一步工作区。</p>
            </div>

            <div class="grid grid-cols-1 gap-3 sm:grid-cols-2">
              <button
                v-for="option in projectTypes"
                :key="option.value"
                type="button"
                :aria-pressed="selectedType === option.value"
                :class="[
                  'group relative min-h-40 rounded-lg border p-5 text-left transition-smooth focus-visible:ring-2 focus-visible:ring-accent focus:outline-none',
                  selectedType === option.value
                    ? 'border-accent bg-accent-soft/60'
                    : 'border-border bg-surface hover:border-accent/40 hover:bg-surface-muted',
                ]"
                @click="selectedType = option.value"
              >
                <div
                  :class="[
                    'mb-4 flex h-10 w-10 items-center justify-center rounded-lg',
                    selectedType === option.value ? 'bg-accent text-white' : 'bg-surface-muted text-text-secondary',
                  ]"
                >
                  <component :is="option.icon" :size="20" />
                </div>
                <h3 class="font-semibold text-text-primary">{{ option.title }}</h3>
                <p class="mt-1 text-sm leading-6 text-text-secondary">{{ option.description }}</p>
                <CheckCircle2
                  v-if="selectedType === option.value"
                  :size="18"
                  class="absolute right-4 top-4 text-accent"
                />
              </button>
            </div>

            <div class="mt-6 flex flex-col-reverse gap-3 border-t border-border pt-4 sm:flex-row sm:items-center sm:justify-between">
              <RouterLink
                to="/projects"
                class="inline-flex min-h-10 items-center justify-center gap-2 rounded-lg border border-border bg-surface px-4 text-sm font-medium text-text-secondary transition-smooth hover:bg-surface-muted"
              >
                <ArrowLeft :size="16" />
                返回项目列表
              </RouterLink>
              <button
                type="button"
                class="inline-flex min-h-10 items-center justify-center gap-2 rounded-lg bg-accent px-5 text-sm font-semibold text-white transition-smooth hover:bg-accent/80 focus-visible:ring-2 focus-visible:ring-accent focus:outline-none"
                @click="step = 2"
              >
                继续
                <ArrowRight :size="16" />
              </button>
            </div>
          </template>

          <template v-else-if="step === 2">
            <div class="mb-5">
              <h2 class="text-lg font-semibold text-text-primary">填写项目资料</h2>
              <p class="mt-1 text-sm text-text-secondary">
                当前类型：{{ projectTypes.find((option) => option.value === selectedType)?.title }}
              </p>
            </div>

            <ToolFormSection label="项目名称" required id-for="project-name">
              <input
                id="project-name"
                v-model="form.name"
                type="text"
                class="w-full rounded-lg border border-border/80 bg-surface-muted px-4 py-2.5 text-text-primary placeholder:text-text-muted focus-visible:border-accent focus-visible:ring-2 focus-visible:ring-accent focus:outline-none"
                placeholder="例如：企业官网重构"
              />
            </ToolFormSection>

            <ToolFormSection label="项目描述" optional id-for="project-description">
              <textarea
                id="project-description"
                v-model="form.description"
                rows="4"
                class="w-full resize-y rounded-lg border border-border/80 bg-surface-muted px-4 py-2.5 text-text-primary placeholder:text-text-muted focus-visible:border-accent focus-visible:ring-2 focus-visible:ring-accent focus:outline-none"
                placeholder="概括项目目标、目标用户和主要交付内容"
              />
            </ToolFormSection>

            <ToolFormSection
              v-if="selectedType !== 'analysis_existing'"
              label="技术栈偏好"
              optional
              as-label
            >
              <div class="grid grid-cols-1 gap-2 sm:grid-cols-2">
                <label
                  v-for="stack in techStacks"
                  :key="stack"
                  class="flex min-h-11 cursor-pointer items-center gap-3 rounded-lg border border-border bg-surface px-3 text-sm text-text-primary transition-smooth hover:bg-surface-muted"
                >
                  <input
                    v-model="form.techStacks"
                    type="checkbox"
                    :value="stack"
                    class="h-4 w-4 rounded border-border accent-[var(--color-accent)]"
                  />
                  <span>{{ stack }}</span>
                </label>
              </div>
            </ToolFormSection>

            <p v-if="error" role="alert" class="mb-4 rounded-lg border border-danger/20 bg-danger/10 px-4 py-3 text-sm text-danger">
              {{ error }}
            </p>

            <div class="flex flex-col-reverse gap-3 border-t border-border pt-4 sm:flex-row sm:items-center sm:justify-between">
              <button
                type="button"
                class="inline-flex min-h-10 items-center justify-center gap-2 rounded-lg border border-border bg-surface px-4 text-sm font-medium text-text-secondary transition-smooth hover:bg-surface-muted"
                @click="step = 1"
              >
                <ArrowLeft :size="16" />
                返回
              </button>
              <button
                type="button"
                :disabled="!form.name.trim() || creating"
                class="inline-flex min-h-10 items-center justify-center gap-2 rounded-lg bg-accent px-5 text-sm font-semibold text-white transition-smooth hover:bg-accent/80 disabled:cursor-not-allowed disabled:bg-surface-muted disabled:text-text-muted focus-visible:ring-2 focus-visible:ring-accent focus:outline-none"
                @click="submitProject"
              >
                <Loader2 v-if="creating" :size="16" class="animate-spin" />
                <FolderPlus v-else :size="16" />
                {{ creating ? '创建中...' : '创建项目' }}
              </button>
            </div>
          </template>

          <template v-else>
            <div class="flex min-h-96 flex-col items-center justify-center px-4 py-10 text-center">
              <div class="mb-4 flex h-14 w-14 items-center justify-center rounded-lg bg-success/10 text-success">
                <CheckCircle2 :size="28" />
              </div>
              <h2 class="text-xl font-bold text-text-primary">项目创建成功</h2>
              <p class="mt-2 max-w-md text-sm leading-6 text-text-secondary">
                <strong class="font-semibold text-text-primary">{{ createdProject?.name }}</strong>
                已加入工作区，可以继续完善项目资料。
              </p>
              <div class="mt-6 flex w-full max-w-sm flex-col gap-3 sm:flex-row sm:justify-center">
                <RouterLink
                  :to="primaryNextRoute"
                  class="inline-flex min-h-10 items-center justify-center gap-2 rounded-lg bg-accent px-4 text-sm font-semibold text-white transition-smooth hover:bg-accent/80"
                >
                  {{ selectedType === 'analysis_existing' ? '开始项目诊断' : '填写项目需求' }}
                  <ArrowRight :size="16" />
                </RouterLink>
                <RouterLink
                  :to="`/projects/${createdProject?.id}`"
                  class="inline-flex min-h-10 items-center justify-center rounded-lg border border-border bg-surface px-4 text-sm font-medium text-text-secondary transition-smooth hover:bg-surface-muted"
                >
                  进入项目
                </RouterLink>
              </div>
            </div>
          </template>
        </section>

        <aside class="rounded-lg border border-border bg-surface p-5 shadow-sm">
          <div class="mb-5 flex items-center gap-3">
            <div class="flex h-9 w-9 items-center justify-center rounded-lg bg-accent-soft text-accent">
              <FileText :size="18" />
            </div>
            <div>
              <h2 class="font-semibold text-text-primary">创建流程</h2>
              <p class="text-xs text-text-muted">三步完成项目初始化</p>
            </div>
          </div>

          <ol class="space-y-1">
            <li
              v-for="(label, index) in ['选择类型', '填写资料', '完成创建']"
              :key="label"
              :class="[
                'flex items-center gap-3 rounded-lg px-3 py-3 text-sm',
                step === index + 1 ? 'bg-accent-soft text-accent' : 'text-text-secondary',
              ]"
            >
              <span
                :class="[
                  'flex h-7 w-7 shrink-0 items-center justify-center rounded-full border text-xs font-semibold',
                  step > index + 1
                    ? 'border-success bg-success text-white'
                    : step === index + 1
                      ? 'border-accent bg-accent text-white'
                      : 'border-border bg-surface text-text-muted',
                ]"
              >
                <CheckCircle2 v-if="step > index + 1" :size="15" />
                <span v-else>{{ index + 1 }}</span>
              </span>
              <span class="font-medium">{{ label }}</span>
            </li>
          </ol>

          <div class="mt-6 border-t border-border pt-5">
            <p class="text-xs font-medium text-text-muted">当前选择</p>
            <div class="mt-3 flex items-center gap-3">
              <div class="flex h-9 w-9 items-center justify-center rounded-lg bg-surface-muted text-text-secondary">
                <component
                  :is="projectTypes.find((option) => option.value === selectedType)?.icon"
                  :size="18"
                />
              </div>
              <div>
                <p class="text-sm font-semibold text-text-primary">
                  {{ projectTypes.find((option) => option.value === selectedType)?.title }}
                </p>
                <p class="text-xs text-text-muted">{{ form.techStacks.length }} 项技术偏好</p>
              </div>
            </div>
          </div>
        </aside>
      </div>
    </ProjectStageShell>
  </div>
</template>
