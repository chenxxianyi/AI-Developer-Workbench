<script setup lang="ts">
/**
 * Landing Page
 * Standalone layout (no AppShell)
 */

import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { RouterLink, useRoute } from 'vue-router'
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
import LanguageSwitcher from '@/components/common/LanguageSwitcher.vue'
import { getLandingNavLinkClass } from '@/utils/landingNav'

const { t } = useI18n()
const route = useRoute()

// Stats data (used in template)
// Tools data
const tools = computed(() => [
  {
    name: t('landing.tools.items.uiReview.name'),
    subtitle: t('landing.tools.items.uiReview.subtitle'),
    icon: 'Eye',
    color: 'accent',
    route: '/tools/ui-review',
    description: t('landing.tools.items.uiReview.description'),
  },
  {
    name: t('landing.tools.items.projectDoctor.name'),
    subtitle: t('landing.tools.items.projectDoctor.subtitle'),
    icon: 'Stethoscope',
    color: 'success',
    route: '/tools/project-doctor',
    description: t('landing.tools.items.projectDoctor.description'),
  },
  {
    name: t('landing.tools.items.agentConfig.name'),
    subtitle: t('landing.tools.items.agentConfig.subtitle'),
    icon: 'Bot',
    color: 'warning',
    route: '/tools/agent-config',
    description: t('landing.tools.items.agentConfig.description'),
  },
  {
    name: t('landing.tools.items.apiDoc.name'),
    subtitle: t('landing.tools.items.apiDoc.subtitle'),
    icon: 'FileText',
    color: 'danger',
    route: '/tools/api-doc',
    description: t('landing.tools.items.apiDoc.description'),
  },
  {
    name: t('landing.tools.items.dbSchema.name'),
    subtitle: t('landing.tools.items.dbSchema.subtitle'),
    icon: 'Database',
    color: 'accent',
    route: '/tools/db-schema',
    description: t('landing.tools.items.dbSchema.description'),
  },
])

// Workflow steps
const workflowSteps = computed(() => [
  { num: 1, title: t('landing.workflow.steps.choose.title'), description: t('landing.workflow.steps.choose.description') },
  { num: 2, title: t('landing.workflow.steps.analyze.title'), description: t('landing.workflow.steps.analyze.description') },
  { num: 3, title: t('landing.workflow.steps.report.title'), description: t('landing.workflow.steps.report.description') },
  { num: 4, title: t('landing.workflow.steps.export.title'), description: t('landing.workflow.steps.export.description') },
])

// Features
const features = computed(() => [
  { title: t('landing.features.items.product.title'), description: t('landing.features.items.product.description') },
  { title: t('landing.features.items.safety.title'), description: t('landing.features.items.safety.description') },
  { title: t('landing.features.items.output.title'), description: t('landing.features.items.output.description') },
  { title: t('landing.features.items.integration.title'), description: t('landing.features.items.integration.description') },
])

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
          <RouterLink to="/" :class="getLandingNavLinkClass(route.hash, '')">{{ t('landing.nav.home') }}</RouterLink>
          <RouterLink :to="{ path: '/', hash: '#tools' }" :class="getLandingNavLinkClass(route.hash, '#tools')">{{ t('landing.nav.tools') }}</RouterLink>
          <RouterLink :to="{ path: '/', hash: '#workflow' }" :class="getLandingNavLinkClass(route.hash, '#workflow')">{{ t('landing.nav.workflow') }}</RouterLink>
          <RouterLink :to="{ path: '/', hash: '#features' }" :class="getLandingNavLinkClass(route.hash, '#features')">{{ t('landing.nav.features') }}</RouterLink>
        </div>

        <div class="flex items-center gap-3">
          <LanguageSwitcher />
          <RouterLink
            to="/dashboard"
            class="px-4 py-2 bg-accent text-white rounded-lg font-medium hover:bg-accent/90 transition-smooth"
          >
            {{ t('landing.nav.start') }}
          </RouterLink>
        </div>
      </div>
    </nav>

    <!-- Hero Section -->
    <section class="min-h-screen pt-32 pb-20 px-4 md:px-8 flex items-center">
      <div class="max-w-content mx-auto">
        <div class="text-center max-w-3xl mx-auto">
          <h1 class="text-4xl md:text-5xl lg:text-6xl font-bold mb-6 leading-tight text-text-primary">
            {{ t('landing.hero.title') }}
          </h1>

          <p class="text-lg md:text-xl text-text-secondary mb-8 leading-relaxed">
            {{ t('landing.hero.description') }}
          </p>

          <div class="flex flex-col sm:flex-row items-center justify-center gap-4 mb-12">
            <RouterLink
              to="/dashboard"
              class="w-full sm:w-auto px-8 py-3 bg-accent text-white rounded-lg font-semibold hover:bg-accent/90 transition-smooth flex items-center justify-center gap-2"
            >
              <LayoutDashboard :size="20" />
              {{ t('landing.hero.enter') }}
            </RouterLink>
            <a
              href="#tools"
              class="w-full sm:w-auto px-8 py-3 bg-surface border border-border text-text-primary rounded-lg font-semibold hover:bg-surface-muted transition-smooth flex items-center justify-center gap-2"
            >
              <ArrowDown :size="20" />
              {{ t('landing.hero.learnMore') }}
            </a>
          </div>

          <!-- Stats -->
          <div class="flex items-center justify-center gap-8 md:gap-12 text-sm text-text-muted">
            <div class="flex items-center gap-2">
              <Wrench :size="16" />
              <span>{{ t('landing.stats.tools') }}</span>
            </div>
            <div class="flex items-center gap-2">
              <FileCheck :size="16" />
              <span>{{ t('landing.stats.analysis') }}</span>
            </div>
            <div class="flex items-center gap-2">
              <Download :size="16" />
              <span>{{ t('landing.stats.export') }}</span>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Tools Section -->
    <section id="tools" class="py-20 px-4 md:px-8 bg-surface">
      <div class="max-w-content mx-auto">
        <div class="text-center mb-12">
          <h2 class="text-3xl md:text-4xl font-bold mb-4 text-text-primary">{{ t('landing.tools.title') }}</h2>
          <p class="text-lg text-text-secondary max-w-2xl mx-auto">
            {{ t('landing.tools.description') }}
          </p>
        </div>

        <div class="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
          <RouterLink
            v-for="tool in tools"
            :key="tool.name"
            :to="tool.route"
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
          <h2 class="text-3xl md:text-4xl font-bold mb-4 text-text-primary">{{ t('landing.workflow.title') }}</h2>
          <p class="text-lg text-text-secondary max-w-2xl mx-auto">
            {{ t('landing.workflow.description') }}
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
            <h2 class="text-3xl md:text-4xl font-bold mb-6 text-text-primary">{{ t('landing.features.title') }}</h2>
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
              <span class="font-semibold text-text-primary">{{ t('landing.sampleReport.title') }}</span>
            </div>
            <div class="bg-surface-muted rounded p-4 text-sm font-mono leading-relaxed">
              <div class="text-text-muted mb-2"># UI 审查报告</div>
              <div class="mb-3">
                <span class="text-warning font-semibold">{{ t('landing.sampleReport.score') }}:</span>
                <span class="text-text-primary"> 78/100</span>
              </div>
              <div class="mb-3">
                <span class="text-danger font-semibold">{{ t('landing.sampleReport.issues') }}:</span>
                <span class="text-text-primary"> 5</span>
              </div>
              <div class="text-text-secondary">
                <div class="mb-2">• {{ t('landing.sampleReport.high') }}</div>
                <div class="mb-2">• {{ t('landing.sampleReport.medium') }}</div>
                <div>• {{ t('landing.sampleReport.low') }}</div>
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
          <h2 class="text-3xl md:text-4xl font-bold mb-4 text-text-primary">{{ t('landing.cta.title') }}</h2>
          <p class="text-lg text-text-secondary mb-8 max-w-2xl mx-auto">
            {{ t('landing.cta.description') }}
          </p>
          <RouterLink
            to="/dashboard"
            class="inline-flex items-center gap-2 px-8 py-3 bg-accent text-white rounded-lg font-semibold hover:bg-accent/90 transition-smooth"
          >
            <LayoutDashboard :size="20" />
            {{ t('landing.cta.button') }}
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
            {{ t('landing.footer.tagline') }}
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
