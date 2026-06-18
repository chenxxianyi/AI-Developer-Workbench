<script setup lang="ts">
/**
 * Landing Page
 * Standalone layout (no AppShell)
 */

import { RouterLink } from 'vue-router'
import {
  Zap,
  LayoutDashboard,
  ArrowDown,
  Eye,
  Stethoscope,
  Bot,
  FileText,
  Database,
  Wrench,
  FileCheck,
  Download,
  CheckCircle2,
  Terminal,
  ExternalLink,
  BookOpen,
} from '@lucide/vue'

// Stats data (used in template)
// Tools data
const tools = [
  {
    name: 'UI Review',
    subtitle: 'UI 质量审查',
    icon: 'Eye',
    color: 'accent',
    description: '根据截图和前端代码审查 UI 质量，识别模板化痕迹，评估设计一致性，提供改进建议。',
  },
  {
    name: 'Project Doctor',
    subtitle: '项目结构检查',
    icon: 'Stethoscope',
    color: 'success',
    description: '静态分析项目 ZIP，检查工程结构、依赖管理、代码规范和潜在风险，生成健康报告。',
  },
  {
    name: 'Agent Config Studio',
    subtitle: 'AI Agent 配置生成',
    icon: 'Bot',
    color: 'warning',
    description: '根据项目特征生成 AGENTS.md、TASK_PLAN.md 等配置文件，优化 AI Coding 效果。',
  },
  {
    name: 'API Doc Builder',
    subtitle: 'API 文档生成',
    icon: 'FileText',
    color: 'danger',
    description: '从代码或项目 ZIP 生成 Markdown/OpenAPI 文档，支持多种后端框架和输出格式。',
  },
  {
    name: 'DB Schema Review',
    subtitle: '数据库结构审查',
    icon: 'Database',
    color: 'accent',
    description: '审查 SQL、GORM、Prisma 等数据库定义，评估表结构、索引、性能和安全问题。',
  },
]

// Workflow steps
const workflowSteps = [
  { num: 1, title: '选择工具', description: '从 Dashboard 选择合适的分析工具，填写表单或上传文件' },
  { num: 2, title: '智能分析', description: '后端进行安全处理和 AI 分析，生成结构化报告和评分' },
  { num: 3, title: '查看报告', description: '查看评分、问题列表和改进建议，了解具体优化方向' },
  { num: 4, title: '复制与导出', description: '复制 Codex Prompt 到 AI Coding 工具，下载生成文件或报告' },
]

// Features
const features = [
  { title: '真实产品感', description: '避免模板化 SaaS 风格，克制、专业的开发者体验' },
  { title: '安全优先', description: '静态分析不执行代码，完整的安全处理和输入校验' },
  { title: '结构化输出', description: '评分、问题、建议清晰分离，支持 Markdown 和 OpenAPI' },
  { title: '一键集成', description: 'Codex Prompt 直接复制到 Claude Code、Cursor 等 AI 工具' },
]

function getIconComponent(iconName: string) {
  const iconMap: Record<string, any> = {
    Eye,
    Stethoscope,
    Bot,
    FileText,
    Database,
    Wrench,
    FileCheck,
    Download,
  }
  return iconMap[iconName] || Zap
}

function getColorClass(color: string) {
  return `bg-${color}-soft`
}

function getIconColorClass(color: string) {
  return `text-${color}`
}
</script>

<template>
  <div class="min-h-screen bg-background">
    <!-- Navigation -->
    <nav class="fixed top-4 left-4 right-4 z-50 bg-surface/95 backdrop-blur-sm border border-border rounded-lg shadow-sm transition-smooth">
      <div class="max-w-content mx-auto px-6 md:px-8 py-3 flex items-center justify-between">
        <RouterLink to="/" class="flex items-center gap-3 hover:opacity-80 transition-smooth">
          <div class="w-8 h-8 bg-accent rounded-lg flex items-center justify-center">
            <Zap :size="20" class="text-white" />
          </div>
          <span class="text-lg font-semibold text-text-primary">AI Workbench</span>
        </RouterLink>

        <div class="hidden md:flex items-center gap-6">
          <a href="#tools" class="text-text-secondary hover:text-accent transition-smooth">工具</a>
          <a href="#workflow" class="text-text-secondary hover:text-accent transition-smooth">工作流</a>
          <a href="#features" class="text-text-secondary hover:text-accent transition-smooth">特性</a>
        </div>

        <RouterLink
          to="/dashboard"
          class="px-4 py-2 bg-accent text-white rounded-lg font-medium hover:bg-accent/90 transition-smooth"
        >
          开始使用
        </RouterLink>
      </div>
    </nav>

    <!-- Hero Section -->
    <section class="min-h-screen pt-32 pb-20 px-4 md:px-8 flex items-center">
      <div class="max-w-content mx-auto">
        <div class="text-center max-w-3xl mx-auto">
          <h1 class="text-4xl md:text-5xl lg:text-6xl font-bold mb-6 leading-tight text-text-primary">
            Build better AI-generated projects
          </h1>

          <p class="text-lg md:text-xl text-text-secondary mb-8 leading-relaxed">
            Review UI quality, inspect project structure, generate AGENTS.md,
            build API docs, and improve database schemas in one developer workbench.
          </p>

          <div class="flex flex-col sm:flex-row items-center justify-center gap-4 mb-12">
            <RouterLink
              to="/dashboard"
              class="w-full sm:w-auto px-8 py-3 bg-accent text-white rounded-lg font-semibold hover:bg-accent/90 transition-smooth flex items-center justify-center gap-2"
            >
              <LayoutDashboard :size="20" />
              进入工作台
            </RouterLink>
            <a
              href="#tools"
              class="w-full sm:w-auto px-8 py-3 bg-surface border border-border text-text-primary rounded-lg font-semibold hover:bg-surface-muted transition-smooth flex items-center justify-center gap-2"
            >
              <ArrowDown :size="20" />
              了解更多
            </a>
          </div>

          <!-- Stats -->
          <div class="flex items-center justify-center gap-8 md:gap-12 text-sm text-text-muted">
            <div class="flex items-center gap-2">
              <Wrench :size="16" />
              <span>5 核心工具</span>
            </div>
            <div class="flex items-center gap-2">
              <FileCheck :size="16" />
              <span>智能分析</span>
            </div>
            <div class="flex items-center gap-2">
              <Download :size="16" />
              <span>一键导出</span>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Tools Section -->
    <section id="tools" class="py-20 px-4 md:px-8 bg-surface">
      <div class="max-w-content mx-auto">
        <div class="text-center mb-12">
          <h2 class="text-3xl md:text-4xl font-bold mb-4 text-text-primary">五大核心工具</h2>
          <p class="text-lg text-text-secondary max-w-2xl mx-auto">
            专为 AI Coding 开发者设计，覆盖 UI、项目结构、Agent 配置、API 文档和数据库设计全流程
          </p>
        </div>

        <div class="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
          <RouterLink
            v-for="tool in tools"
            :key="tool.name"
            :to="`/tools/${tool.name.toLowerCase().replace(' ', '-')}`"
            class="group p-6 bg-background border border-border rounded-lg hover:border-accent hover:shadow-md transition-smooth"
          >
            <div class="flex items-start gap-4 mb-4">
              <div
                :class="[
                  'w-12 h-12 rounded-lg flex items-center justify-center flex-shrink-0',
                  getColorClass(tool.color),
                ]"
              >
                <component
                  :is="getIconComponent(tool.icon)"
                  :size="24"
                  :class="getIconColorClass(tool.color)"
                />
              </div>
              <div>
                <h3 class="text-xl font-semibold mb-2 text-text-primary">{{ tool.name }}</h3>
                <p class="text-text-secondary">{{ tool.subtitle }}</p>
              </div>
            </div>
            <p class="text-text-secondary leading-relaxed">{{ tool.description }}</p>
          </RouterLink>

          <!-- Empty slot for grid alignment -->
          <div class="hidden lg:block"></div>
        </div>
      </div>
    </section>

    <!-- Workflow Section -->
    <section id="workflow" class="py-20 px-4 md:px-8">
      <div class="max-w-content mx-auto">
        <div class="text-center mb-12">
          <h2 class="text-3xl md:text-4xl font-bold mb-4 text-text-primary">统一工作流</h2>
          <p class="text-lg text-text-secondary max-w-2xl mx-auto">
            从输入到报告，每个工具遵循一致的交互闭环
          </p>
        </div>

        <div class="max-w-3xl mx-auto">
          <div class="space-y-6">
            <div
              v-for="step in workflowSteps"
              :key="step.num"
              class="flex items-start gap-4 p-6 bg-surface border border-border rounded-lg"
            >
              <div class="w-10 h-10 bg-accent text-white rounded-full flex items-center justify-center font-bold flex-shrink-0">
                {{ step.num }}
              </div>
              <div>
                <h3 class="text-xl font-semibold mb-2 text-text-primary">{{ step.title }}</h3>
                <p class="text-text-secondary">{{ step.description }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Features Section -->
    <section id="features" class="py-20 px-4 md:px-8 bg-surface">
      <div class="max-w-content mx-auto">
        <div class="grid md:grid-cols-2 gap-8 items-center">
          <div>
            <h2 class="text-3xl md:text-4xl font-bold mb-6 text-text-primary">专为开发者设计</h2>
            <div class="space-y-4">
              <div v-for="feature in features" :key="feature.title" class="flex items-start gap-3">
                <CheckCircle2 :size="24" class="text-success flex-shrink-0" />
                <div>
                  <h3 class="font-semibold mb-1 text-text-primary">{{ feature.title }}</h3>
                  <p class="text-text-secondary">{{ feature.description }}</p>
                </div>
              </div>
            </div>
          </div>

          <div class="bg-background border border-border rounded-lg p-6">
            <div class="flex items-center gap-2 mb-4">
              <Terminal :size="20" class="text-accent" />
              <span class="font-semibold text-text-primary">示例报告</span>
            </div>
            <div class="bg-surface-muted rounded p-4 text-sm font-mono leading-relaxed">
              <div class="text-text-muted mb-2"># UI Review Report</div>
              <div class="mb-3">
                <span class="text-warning font-semibold">Total Score:</span>
                <span class="text-text-primary"> 78/100</span>
              </div>
              <div class="mb-3">
                <span class="text-danger font-semibold">Issues Found:</span>
                <span class="text-text-primary"> 5</span>
              </div>
              <div class="text-text-secondary">
                <div class="mb-2">• High: AI template risk detected</div>
                <div class="mb-2">• Medium: Inconsistent button spacing</div>
                <div>• Low: Missing hover states</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- CTA Section -->
    <section class="py-20 px-4 md:px-8">
      <div class="max-w-content mx-auto">
        <div class="bg-surface border border-border rounded-lg p-8 md:p-12 text-center">
          <h2 class="text-3xl md:text-4xl font-bold mb-4 text-text-primary">开始提升你的 AI 项目质量</h2>
          <p class="text-lg text-text-secondary mb-8 max-w-2xl mx-auto">
            立即进入工作台，体验智能分析工具带来的效率提升
          </p>
          <RouterLink
            to="/dashboard"
            class="inline-flex items-center gap-2 px-8 py-3 bg-accent text-white rounded-lg font-semibold hover:bg-accent/90 transition-smooth"
          >
            <LayoutDashboard :size="20" />
            进入 Dashboard
          </RouterLink>
        </div>
      </div>
    </section>

    <!-- Footer -->
    <footer class="py-12 px-4 md:px-8 bg-surface border-t border-border">
      <div class="max-w-content mx-auto">
        <div class="flex flex-col md:flex-row items-center justify-between gap-4">
          <div class="flex items-center gap-3">
            <div class="w-8 h-8 bg-accent rounded-lg flex items-center justify-center">
              <Zap :size="20" class="text-white" />
            </div>
            <span class="text-lg font-semibold text-text-primary">AI Developer Workbench</span>
          </div>

          <div class="text-sm text-text-muted">
            MVP 0.1.0 · 专为 AI Coding 开发者设计
          </div>

          <div class="flex items-center gap-4">
            <a href="#" class="text-text-secondary hover:text-accent transition-smooth">
              <ExternalLink :size="20" />
            </a>
            <a href="#" class="text-text-secondary hover:text-accent transition-smooth">
              <BookOpen :size="20" />
            </a>
          </div>
        </div>
      </div>
    </footer>
  </div>
</template>