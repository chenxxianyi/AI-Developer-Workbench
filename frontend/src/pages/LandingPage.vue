<script setup lang="ts">
/**
 * Landing Page
 * Standalone layout (no AppShell)
 */

import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { RouterLink, useRoute } from 'vue-router'
import {
  Zap,
  LayoutDashboard,
  ArrowDown,
  ArrowRight,
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
  ShieldCheck,
} from '@lucide/vue'
import LanguageSwitcher from '@/components/common/LanguageSwitcher.vue'
import { getLandingNavLinkClass } from '@/utils/landingNav'

const { t } = useI18n()
const route = useRoute()

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

const workflowSteps = computed(() => [
  { num: 1, title: t('landing.workflow.steps.choose.title'), description: t('landing.workflow.steps.choose.description') },
  { num: 2, title: t('landing.workflow.steps.analyze.title'), description: t('landing.workflow.steps.analyze.description') },
  { num: 3, title: t('landing.workflow.steps.report.title'), description: t('landing.workflow.steps.report.description') },
  { num: 4, title: t('landing.workflow.steps.export.title'), description: t('landing.workflow.steps.export.description') },
])

const features = computed(() => [
  { title: t('landing.features.items.product.title'), description: t('landing.features.items.product.description'), icon: LayoutDashboard },
  { title: t('landing.features.items.safety.title'), description: t('landing.features.items.safety.description'), icon: ShieldCheck },
  { title: t('landing.features.items.output.title'), description: t('landing.features.items.output.description'), icon: FileCheck },
  { title: t('landing.features.items.integration.title'), description: t('landing.features.items.integration.description'), icon: Bot },
])

const statItems = computed(() => [
  { icon: Wrench, label: t('landing.stats.tools') },
  { icon: FileCheck, label: t('landing.stats.analysis') },
  { icon: Download, label: t('landing.stats.export') },
])

const qualityMetrics = computed(() => [
  { label: t('landing.sampleReport.score'), value: '78', suffix: '/100', tone: 'text-warning' },
  { label: t('landing.sampleReport.issues'), value: '5', suffix: '', tone: 'text-danger' },
  { label: t('landing.stats.tools'), value: '5', suffix: '', tone: 'text-accent' },
])

const sampleIssues = computed(() => [
  { level: 'H', text: t('landing.sampleReport.high'), levelClass: 'border-danger/30 bg-danger/15 text-red-200', barClass: 'w-4/5 bg-danger' },
  { level: 'M', text: t('landing.sampleReport.medium'), levelClass: 'border-warning/30 bg-warning/15 text-amber-200', barClass: 'w-3/5 bg-warning' },
  { level: 'L', text: t('landing.sampleReport.low'), levelClass: 'border-success/30 bg-success/15 text-emerald-200', barClass: 'w-2/5 bg-success' },
])

const terminalTiltStyle = ref<Record<string, string>>({
  '--terminal-rotate-x': '0deg',
  '--terminal-rotate-y': '0deg',
  '--terminal-lift': '0px',
  '--terminal-shadow-x': '0px',
  '--terminal-shadow-y': '28px',
  '--terminal-glow-x': '50%',
  '--terminal-glow-y': '50%',
})

const heroGridStyle = ref<Record<string, string>>({
  '--hero-grid-x': '50%',
  '--hero-grid-y': '42%',
  '--hero-grid-opacity': '0',
})

const isLandingNavCompact = ref(false)
let landingNavScrollFrame: number | null = null

const colorClassMap: Record<string, string> = {
  accent: 'bg-accent-soft text-accent border-accent/10',
  success: 'bg-success/10 text-success border-success/15',
  warning: 'bg-warning/10 text-warning border-warning/15',
  danger: 'bg-danger/10 text-danger border-danger/15',
}

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

function getToolIconClass(color: string) {
  return colorClassMap[color] || colorClassMap.accent
}

function handleTerminalPointerMove(event: PointerEvent) {
  if (event.pointerType === 'touch' || window.matchMedia('(prefers-reduced-motion: reduce)').matches) {
    return
  }

  const target = event.currentTarget as HTMLElement
  const rect = target.getBoundingClientRect()
  const x = (event.clientX - rect.left) / rect.width - 0.5
  const y = (event.clientY - rect.top) / rect.height - 0.5
  const rotateX = y * -8
  const rotateY = x * 10

  terminalTiltStyle.value = {
    '--terminal-rotate-x': `${rotateX.toFixed(2)}deg`,
    '--terminal-rotate-y': `${rotateY.toFixed(2)}deg`,
    '--terminal-lift': '-6px',
    '--terminal-shadow-x': `${(-x * 24).toFixed(1)}px`,
    '--terminal-shadow-y': `${(30 + Math.abs(y) * 14).toFixed(1)}px`,
    '--terminal-glow-x': `${((x + 0.5) * 100).toFixed(1)}%`,
    '--terminal-glow-y': `${((y + 0.5) * 100).toFixed(1)}%`,
  }
}

function resetTerminalTilt() {
  terminalTiltStyle.value = {
    '--terminal-rotate-x': '0deg',
    '--terminal-rotate-y': '0deg',
    '--terminal-lift': '0px',
    '--terminal-shadow-x': '0px',
    '--terminal-shadow-y': '28px',
    '--terminal-glow-x': '50%',
    '--terminal-glow-y': '50%',
  }
}

function handleHeroPointerMove(event: PointerEvent) {
  if (event.pointerType === 'touch' || window.matchMedia('(prefers-reduced-motion: reduce)').matches) {
    return
  }

  const target = event.currentTarget as HTMLElement
  const rect = target.getBoundingClientRect()
  const x = ((event.clientX - rect.left) / rect.width) * 100
  const y = ((event.clientY - rect.top) / rect.height) * 100

  heroGridStyle.value = {
    '--hero-grid-x': `${Math.max(0, Math.min(100, x)).toFixed(2)}%`,
    '--hero-grid-y': `${Math.max(0, Math.min(100, y)).toFixed(2)}%`,
    '--hero-grid-opacity': '1',
  }
}

function resetHeroGridGlow() {
  heroGridStyle.value = {
    ...heroGridStyle.value,
    '--hero-grid-opacity': '0',
  }
}

function updateLandingNavState() {
  const scrollY = window.scrollY

  if (!isLandingNavCompact.value && scrollY > 72) {
    isLandingNavCompact.value = true
    return
  }

  if (isLandingNavCompact.value && scrollY < 24) {
    isLandingNavCompact.value = false
  }
}

function handleLandingNavScroll() {
  if (landingNavScrollFrame !== null) {
    return
  }

  landingNavScrollFrame = window.requestAnimationFrame(() => {
    landingNavScrollFrame = null
    updateLandingNavState()
  })
}

onMounted(() => {
  updateLandingNavState()
  window.addEventListener('scroll', handleLandingNavScroll, { passive: true })
})

onUnmounted(() => {
  window.removeEventListener('scroll', handleLandingNavScroll)

  if (landingNavScrollFrame !== null) {
    window.cancelAnimationFrame(landingNavScrollFrame)
  }
})
</script>

<template>
  <div class="min-h-screen overflow-x-hidden bg-background text-text-primary">
    <!-- Navigation -->
    <nav
      class="landing-nav fixed z-50 overflow-hidden rounded-xl border border-border/80 bg-surface/92 shadow-[0_14px_40px_rgba(15,23,42,0.10)] backdrop-blur-md"
      :class="{ 'landing-nav--compact': isLandingNavCompact }"
    >
      <div class="landing-nav__inner mx-auto flex max-w-content items-center justify-between gap-4 px-4 py-3 md:gap-6 md:px-7">
        <RouterLink to="/" class="flex min-w-0 items-center gap-3.5 transition-smooth hover:opacity-85">
          <div class="flex h-11 w-11 flex-shrink-0 items-center justify-center rounded-lg bg-accent shadow-[0_10px_24px_rgba(37,99,235,0.28)] md:h-12 md:w-12">
            <Zap :size="24" class="text-white" />
          </div>
          <span class="hidden truncate text-xl font-bold text-text-primary sm:inline md:text-2xl">AI Workbench</span>
        </RouterLink>

        <div class="hidden items-center gap-1.5 rounded-lg border border-border/70 bg-background/70 p-1 md:flex">
          <RouterLink to="/" :class="getLandingNavLinkClass(route.hash, '')">{{ t('landing.nav.home') }}</RouterLink>
          <RouterLink :to="{ path: '/', hash: '#tools' }" :class="getLandingNavLinkClass(route.hash, '#tools')">{{ t('landing.nav.tools') }}</RouterLink>
          <RouterLink :to="{ path: '/', hash: '#workflow' }" :class="getLandingNavLinkClass(route.hash, '#workflow')">{{ t('landing.nav.workflow') }}</RouterLink>
          <RouterLink :to="{ path: '/', hash: '#features' }" :class="getLandingNavLinkClass(route.hash, '#features')">{{ t('landing.nav.features') }}</RouterLink>
        </div>

        <div class="flex flex-shrink-0 items-center gap-3">
          <div class="landing-language-switcher">
            <LanguageSwitcher />
          </div>
          <RouterLink
            to="/dashboard"
            class="hidden min-h-11 items-center gap-2 rounded-lg bg-accent px-5 py-2.5 text-base font-bold text-white shadow-[0_10px_24px_rgba(37,99,235,0.24)] transition-smooth hover:-translate-y-0.5 hover:bg-accent/90 hover:shadow-[0_14px_28px_rgba(37,99,235,0.30)] sm:inline-flex"
          >
            {{ t('landing.nav.start') }}
          </RouterLink>
        </div>
      </div>
    </nav>

    <!-- Hero Section -->
    <section
      class="relative overflow-hidden px-4 pb-12 pt-20 md:px-8 md:pb-14 md:pt-24"
      :style="heroGridStyle"
      @pointermove="handleHeroPointerMove"
      @pointerleave="resetHeroGridGlow"
      @pointercancel="resetHeroGridGlow"
    >
      <div class="landing-grid-pattern absolute inset-0" aria-hidden="true"></div>
      <div class="landing-grid-glow absolute inset-0" aria-hidden="true"></div>

      <div class="relative mx-auto max-w-content">
        <div class="grid min-h-[auto] items-center gap-10 md:min-h-[620px] lg:grid-cols-[minmax(0,0.82fr)_minmax(520px,1.08fr)] lg:gap-14">
          <div class="mx-auto max-w-3xl text-center lg:mx-0 lg:text-left">
            <div class="mb-5 inline-flex items-center gap-2 rounded-full border border-border bg-surface/80 px-3 py-1 text-sm font-semibold text-text-secondary shadow-sm">
              <ShieldCheck :size="16" class="text-success" />
              <span>AI Developer Workbench</span>
            </div>

            <h1 class="mx-auto max-w-[10em] text-2xl font-bold leading-tight text-text-primary sm:max-w-[12em] sm:text-4xl md:max-w-none md:text-5xl lg:mx-0 lg:text-6xl">
              {{ t('landing.hero.title') }}
            </h1>

            <p class="mx-auto mt-6 max-w-2xl text-lg leading-relaxed text-text-secondary md:text-xl lg:mx-0">
              {{ t('landing.hero.description') }}
            </p>

            <div class="mt-8 flex flex-col items-stretch justify-center gap-3 sm:flex-row lg:justify-start">
              <RouterLink
                to="/dashboard"
                class="inline-flex min-h-12 items-center justify-center gap-2 rounded-lg bg-accent px-7 py-3 font-semibold text-white shadow-sm transition-smooth hover:bg-accent/90 hover:shadow-md"
              >
                <LayoutDashboard :size="20" />
                {{ t('landing.hero.enter') }}
              </RouterLink>
              <a
                href="#tools"
                class="inline-flex min-h-12 items-center justify-center gap-2 rounded-lg border border-border bg-surface px-7 py-3 font-semibold text-text-primary shadow-sm transition-smooth hover:border-accent/30 hover:bg-surface-muted"
              >
                <ArrowDown :size="20" />
                {{ t('landing.hero.learnMore') }}
              </a>
            </div>

            <div class="mt-9 grid gap-3 text-sm text-text-muted sm:grid-cols-3">
              <div
                v-for="item in statItems"
                :key="item.label"
                class="flex items-center justify-center gap-2 rounded-lg border border-border bg-surface/70 px-3 py-3 lg:justify-start"
              >
                <component :is="item.icon" :size="16" class="text-accent" />
                <span>{{ item.label }}</span>
              </div>
            </div>
          </div>

          <div
            class="terminal-showcase relative mx-auto w-full max-w-xl min-w-0 lg:max-w-none"
            :style="terminalTiltStyle"
            @pointermove="handleTerminalPointerMove"
            @pointerleave="resetTerminalTilt"
            @pointercancel="resetTerminalTilt"
          >
            <p class="mb-5 text-left text-sm font-medium tracking-wide text-text-muted md:text-base">
              图 01 — 一场 AI 项目交付会话，缓存依然温热
            </p>

            <div class="terminal-window overflow-hidden rounded-[28px] bg-[#202326] shadow-[0_28px_80px_rgba(24,24,27,0.22)] ring-1 ring-black/10">
              <div class="relative flex h-14 items-center justify-center border-b border-white/8 px-5 md:h-16">
                <div class="absolute left-5 flex items-center gap-2.5">
                  <span class="h-3 w-3 rounded-full bg-[#52575d]"></span>
                  <span class="h-3 w-3 rounded-full bg-[#52575d]"></span>
                  <span class="h-3 w-3 rounded-full bg-[#52575d]"></span>
                </div>
                <div class="flex items-center gap-2 font-mono text-sm font-semibold text-white/38 md:text-lg">
                  <Terminal :size="18" />
                  <span>~/workspace — ai-workbench</span>
                </div>
              </div>

              <div class="min-h-[320px] px-5 py-7 font-mono text-sm leading-relaxed text-white/76 md:min-h-[390px] md:px-9 md:py-10 md:text-base lg:text-lg">
                <div class="mb-4 flex flex-wrap items-center gap-x-3 gap-y-1">
                  <span class="text-accent">◆</span>
                  <span class="font-bold text-white/86">ai-workbench</span>
                  <span class="text-white/35">v0.1.0 · codex · ~/project</span>
                </div>

                <div class="mb-4 flex items-start gap-3">
                  <span class="text-[#8b5cf6]">›</span>
                  <span class="font-semibold text-white/90">review UI quality and project structure</span>
                </div>

                <div class="mb-3 flex items-start gap-3">
                  <span class="text-[#5fd18a]">✓</span>
                  <span>
                    scan
                    <span class="text-white/38">frontend/src/pages/LandingPage.vue</span>
                    <span class="font-semibold text-[#5fd18a]"> ok</span>
                  </span>
                </div>

                <div class="mb-3 flex items-start gap-3">
                  <span class="text-[#5fd18a]">✓</span>
                  <span>
                    generate
                    <span class="text-white/38">AGENTS.md prompt</span>
                    <span class="font-semibold text-[#5fd18a]"> ready</span>
                  </span>
                </div>

                <div class="mb-3 flex items-start gap-3">
                  <span class="text-[#5fd18a]">✓</span>
                  <span>
                    run
                    <span class="text-white/38">npm run build</span>
                    <span class="font-semibold text-[#5fd18a]"> passed</span>
                    <span class="text-white/36"> (0.74s)</span>
                  </span>
                </div>

                <div class="flex items-start gap-3 text-white/10">
                  <span>●</span>
                  <span>5 tools · cache 92.4% → 96.8%</span>
                </div>
              </div>

              <div class="grid gap-3 border-t border-white/8 px-5 py-4 font-mono text-sm text-white/42 sm:grid-cols-4 md:px-9 md:text-base">
                <div>cache <span class="font-bold text-white/82">96.8%</span></div>
                <div>session <span class="font-bold text-white/82">4m</span></div>
                <div>model <span class="font-bold text-white/82">codex</span></div>
                <div class="sm:text-right">cost <span class="font-bold text-white/82">$0.002</span></div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Tools Section -->
    <section id="tools" class="border-y border-border bg-surface px-4 py-16 md:px-8 md:py-20">
      <div class="mx-auto max-w-content">
        <div class="mb-10 grid gap-4 md:grid-cols-[0.8fr_1fr] md:items-end">
          <h2 class="text-3xl font-bold text-text-primary md:text-4xl">{{ t('landing.tools.title') }}</h2>
          <p class="text-lg leading-relaxed text-text-secondary">
            {{ t('landing.tools.description') }}
          </p>
        </div>

        <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-6">
          <RouterLink
            v-for="(tool, index) in tools"
            :key="tool.name"
            :to="tool.route"
            class="group rounded-lg border border-border bg-background p-5 transition-smooth hover:border-accent/40 hover:bg-surface hover:shadow-md md:p-6"
            :class="index < 3 ? 'lg:col-span-2' : 'lg:col-span-3'"
          >
            <div class="mb-5 flex items-start justify-between gap-4">
              <div :class="['flex h-12 w-12 flex-shrink-0 items-center justify-center rounded-lg border', getToolIconClass(tool.color)]">
                <component
                  :is="getIconComponent(tool.icon)"
                  :size="24"
                />
              </div>
              <ArrowRight :size="18" class="mt-1 text-text-muted transition-smooth group-hover:translate-x-1 group-hover:text-accent" />
            </div>

            <h3 class="mb-2 text-xl font-semibold text-text-primary">{{ tool.name }}</h3>
            <p class="mb-4 text-sm font-semibold text-text-muted">{{ tool.subtitle }}</p>
            <p class="leading-relaxed text-text-secondary">{{ tool.description }}</p>
          </RouterLink>
        </div>
      </div>
    </section>

    <!-- Workflow Section -->
    <section id="workflow" class="px-4 py-16 md:px-8 md:py-20">
      <div class="mx-auto max-w-content">
        <div class="mx-auto mb-12 max-w-3xl text-center">
          <h2 class="mb-4 text-3xl font-bold text-text-primary md:text-4xl">{{ t('landing.workflow.title') }}</h2>
          <p class="text-lg leading-relaxed text-text-secondary">
            {{ t('landing.workflow.description') }}
          </p>
        </div>

        <div class="grid gap-4 md:grid-cols-4">
          <div
            v-for="step in workflowSteps"
            :key="step.num"
            class="rounded-lg border border-border bg-surface p-5 shadow-sm"
          >
            <div class="mb-5 flex items-center justify-between">
              <div class="flex h-10 w-10 items-center justify-center rounded-full bg-accent text-sm font-bold text-white">
                {{ step.num }}
              </div>
              <div class="hidden h-px flex-1 bg-border md:ml-4 md:block"></div>
            </div>
            <h3 class="mb-2 text-lg font-semibold text-text-primary">{{ step.title }}</h3>
            <p class="leading-relaxed text-text-secondary">{{ step.description }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- Features Section -->
    <section id="features" class="features-lab relative overflow-hidden border-y border-border px-4 py-16 md:px-8 md:py-24">
      <div class="relative mx-auto max-w-content">
        <div class="mb-10 grid gap-5 lg:grid-cols-[0.82fr_1fr] lg:items-end">
          <div>
            <div class="mb-4 inline-flex items-center gap-2 rounded-full border border-accent/15 bg-surface/80 px-3 py-1 text-sm font-semibold text-accent shadow-sm">
              <Terminal :size="15" />
              <span>{{ t('landing.features.kicker') }}</span>
            </div>
            <h2 class="max-w-2xl text-3xl font-bold leading-tight text-text-primary md:text-4xl">
              {{ t('landing.features.title') }}
            </h2>
          </div>
          <p class="max-w-2xl text-lg leading-relaxed text-text-secondary lg:justify-self-end">
            {{ t('landing.tools.description') }}
          </p>
        </div>

        <div class="grid gap-5 lg:grid-cols-[minmax(0,0.9fr)_minmax(480px,1.1fr)] lg:items-stretch">
          <div class="grid gap-4 sm:grid-cols-2">
            <div
              v-for="feature in features"
              :key="feature.title"
              class="feature-card group relative overflow-hidden rounded-lg border border-border bg-surface/85 p-5 shadow-sm backdrop-blur transition-smooth hover:-translate-y-0.5 hover:border-accent/30 hover:shadow-md"
            >
              <div class="mb-5 flex items-center justify-between gap-4">
                <div class="flex h-11 w-11 items-center justify-center rounded-lg border border-accent/10 bg-accent-soft text-accent transition-smooth group-hover:border-accent/25 group-hover:bg-accent group-hover:text-white">
                  <component :is="feature.icon" :size="22" />
                </div>
                <CheckCircle2 :size="18" class="text-success/80" />
              </div>
              <h3 class="mb-2 text-lg font-semibold text-text-primary">{{ feature.title }}</h3>
              <p class="leading-relaxed text-text-secondary">{{ feature.description }}</p>
            </div>
          </div>

          <div class="quality-panel relative overflow-hidden rounded-lg border border-slate-800 bg-[#151a21] p-4 text-white shadow-[0_24px_70px_rgba(15,23,42,0.28)] md:p-6">
            <div class="relative z-10">
              <div class="mb-6 flex flex-wrap items-center justify-between gap-3 border-b border-white/10 pb-4">
                <div class="flex items-center gap-3">
                  <div class="flex items-center gap-1.5">
                    <span class="h-2.5 w-2.5 rounded-full bg-white/25"></span>
                    <span class="h-2.5 w-2.5 rounded-full bg-white/25"></span>
                    <span class="h-2.5 w-2.5 rounded-full bg-white/25"></span>
                  </div>
                  <div class="font-mono text-sm font-semibold text-white/72">{{ t('landing.sampleReport.title') }}</div>
                </div>
                <div class="inline-flex items-center gap-2 rounded-full border border-success/25 bg-success/10 px-3 py-1 text-xs font-semibold text-emerald-200">
                  <CheckCircle2 :size="14" />
                  <span>{{ t('landing.sampleReport.promptReady') }}</span>
                </div>
              </div>

              <div class="grid gap-5 md:grid-cols-[170px_1fr]">
                <div class="quality-score-orb flex h-36 items-center justify-center rounded-lg border border-white/10 bg-white/[0.03] md:aspect-square md:h-auto">
                  <div class="text-center">
                    <div class="font-mono text-5xl font-bold leading-none text-white">78</div>
                    <div class="mt-2 text-xs font-semibold text-white/42">{{ t('landing.sampleReport.qualityLabel') }}</div>
                  </div>
                </div>

                <div class="min-w-0">
                  <div class="mb-4 grid grid-cols-3 gap-2">
                    <div
                      v-for="metric in qualityMetrics"
                      :key="metric.label"
                      class="rounded-md border border-white/10 bg-white/[0.035] px-3 py-3"
                    >
                      <div class="mb-1 truncate text-xs font-medium text-white/42">{{ metric.label }}</div>
                      <div class="font-mono text-xl font-bold text-white">
                        <span :class="metric.tone">{{ metric.value }}</span><span class="text-sm text-white/42">{{ metric.suffix }}</span>
                      </div>
                    </div>
                  </div>

                  <div class="space-y-3">
                    <div
                      v-for="issue in sampleIssues"
                      :key="issue.text"
                      class="grid grid-cols-[auto_1fr] items-center gap-3 rounded-md border border-white/10 bg-black/12 px-3 py-3 sm:grid-cols-[auto_1fr_84px]"
                    >
                      <span :class="['flex h-7 w-7 items-center justify-center rounded-md border font-mono text-xs font-bold', issue.levelClass]">
                        {{ issue.level }}
                      </span>
                      <span class="min-w-0 text-sm leading-relaxed text-white/72">{{ issue.text }}</span>
                      <span class="col-span-2 h-1.5 overflow-hidden rounded-full bg-white/10 sm:col-span-1">
                        <span :class="['block h-full rounded-full', issue.barClass]"></span>
                      </span>
                    </div>
                  </div>
                </div>
              </div>

              <div class="mt-5 grid grid-cols-2 gap-3 border-t border-white/10 pt-5 sm:grid-cols-4">
                <div
                  v-for="step in workflowSteps"
                  :key="step.num"
                  class="flex items-center gap-3 rounded-md border border-white/10 bg-white/[0.03] px-3 py-3"
                >
                  <span class="flex h-7 w-7 flex-shrink-0 items-center justify-center rounded-md bg-white/10 font-mono text-xs font-bold text-white">
                    0{{ step.num }}
                  </span>
                  <span class="truncate text-sm font-semibold text-white/70">{{ step.title }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- CTA Section -->
    <section class="cta-command relative overflow-hidden px-4 py-14 text-white md:px-8 md:py-16">
      <div class="relative mx-auto max-w-content">
        <div class="grid gap-8 lg:grid-cols-[1fr_auto] lg:items-center">
          <div>
            <div class="mb-4 inline-flex items-center gap-2 rounded-full border border-white/10 bg-white/[0.04] px-3 py-1 text-sm font-semibold text-white/62">
              <Terminal :size="15" />
              <span>{{ t('landing.cta.command') }}</span>
            </div>
            <h2 class="max-w-2xl text-3xl font-bold leading-tight md:text-4xl">{{ t('landing.cta.title') }}</h2>
            <p class="mt-4 max-w-2xl text-lg leading-relaxed text-white/68">
              {{ t('landing.cta.description') }}
            </p>
          </div>
          <div class="flex flex-col gap-3 sm:flex-row lg:flex-col lg:items-stretch">
            <RouterLink
              to="/dashboard"
              class="inline-flex min-h-12 items-center justify-center gap-2 rounded-lg bg-white px-8 py-3 font-semibold text-text-primary shadow-sm transition-smooth hover:bg-accent hover:text-white"
            >
              <LayoutDashboard :size="20" />
              {{ t('landing.cta.button') }}
              <ArrowRight :size="18" />
            </RouterLink>
            <div class="inline-flex min-h-12 items-center justify-center gap-2 rounded-lg border border-white/10 bg-white/[0.04] px-5 py-3 font-mono text-sm text-white/64">
              <CheckCircle2 :size="16" class="text-success" />
              <span>{{ t('landing.stats.analysis') }} · {{ t('landing.stats.export') }}</span>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Footer -->
    <footer class="border-t border-border bg-surface px-4 py-10 md:px-8">
      <div class="mx-auto max-w-content">
        <div class="flex flex-col items-center justify-between gap-4 md:flex-row">
          <div class="flex items-center gap-3">
            <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-accent">
              <Zap :size="20" class="text-white" />
            </div>
            <span class="text-lg font-semibold text-text-primary">AI Developer Workbench</span>
          </div>

          <div class="text-sm text-text-muted">
            {{ t('landing.footer.tagline') }}
          </div>

          <div class="flex items-center gap-4">
            <a href="#" class="text-text-secondary transition-smooth hover:text-accent">
              <ExternalLink :size="20" />
            </a>
            <a href="#" class="text-text-secondary transition-smooth hover:text-accent">
              <BookOpen :size="20" />
            </a>
          </div>
        </div>
      </div>
    </footer>
  </div>
</template>

<style scoped>
.landing-nav {
  top: 0.75rem;
  left: 50%;
  width: calc(100vw - 1.5rem);
  transform: translateX(-50%) translateZ(0);
  transition:
    width 900ms cubic-bezier(0.16, 1, 0.3, 1),
    transform 900ms cubic-bezier(0.16, 1, 0.3, 1),
    box-shadow 720ms ease-out,
    background-color 720ms ease-out;
}

.landing-nav__inner {
  min-height: 64px;
  transition:
    min-height 900ms cubic-bezier(0.16, 1, 0.3, 1),
    padding 900ms cubic-bezier(0.16, 1, 0.3, 1);
}

.landing-nav--compact {
  width: min(1120px, calc(100vw - 4rem));
  transform: translateX(-50%) translateZ(0);
  background-color: rgba(255, 255, 255, 0.96);
  box-shadow: 0 18px 46px rgba(15, 23, 42, 0.14);
}

.landing-nav--compact .landing-nav__inner {
  min-height: 60px;
  padding-block: 0.625rem;
}

@media (min-width: 768px) {
  .landing-nav {
    top: 1.25rem;
    width: calc(100vw - 2.5rem);
  }

  .landing-nav__inner {
    min-height: 72px;
  }

  .landing-nav--compact {
    width: min(1040px, calc(100vw - 7.5rem));
    transform: translateX(-50%) translateZ(0);
  }

  .landing-nav--compact .landing-nav__inner {
    min-height: 62px;
    padding-inline: 1.25rem;
  }
}

@media (max-width: 639px) {
  .landing-nav--compact {
    width: calc(100vw - 1.5rem);
  }

  .landing-nav--compact .landing-nav__inner {
    min-height: 64px;
    padding-block: 0.75rem;
  }
}

.landing-grid-pattern {
  background-image:
    linear-gradient(rgba(37, 99, 235, 0.07) 1px, transparent 1px),
    linear-gradient(90deg, rgba(37, 99, 235, 0.07) 1px, transparent 1px);
  background-size: 44px 44px;
  mask-image: linear-gradient(to bottom, black 0%, black 70%, transparent 100%);
}

.landing-grid-glow {
  pointer-events: none;
  opacity: var(--hero-grid-opacity);
  background-image:
    linear-gradient(rgba(37, 99, 235, 0.38) 1px, transparent 1px),
    linear-gradient(90deg, rgba(37, 99, 235, 0.38) 1px, transparent 1px),
    radial-gradient(
      circle at var(--hero-grid-x) var(--hero-grid-y),
      rgba(37, 99, 235, 0.24),
      rgba(37, 99, 235, 0.1) 28%,
      transparent 56%
    );
  background-position:
    0 0,
    0 0,
    0 0;
  background-size:
    44px 44px,
    44px 44px,
    100% 100%;
  mask-image:
    radial-gradient(
      circle 280px at var(--hero-grid-x) var(--hero-grid-y),
      black 0%,
      rgba(0, 0, 0, 0.9) 36%,
      transparent 76%
    ),
    linear-gradient(to bottom, black 0%, black 72%, transparent 100%);
  mask-composite: intersect;
  transition: opacity 180ms ease-out;
}

.terminal-showcase {
  --terminal-rotate-x: 0deg;
  --terminal-rotate-y: 0deg;
  --terminal-lift: 0px;
  --terminal-shadow-x: 0px;
  --terminal-shadow-y: 28px;
  --terminal-glow-x: 50%;
  --terminal-glow-y: 50%;
  perspective: 1200px;
}

.terminal-window {
  color-scheme: dark;
  position: relative;
  transform:
    translateY(var(--terminal-lift))
    rotateX(var(--terminal-rotate-x))
    rotateY(var(--terminal-rotate-y));
  transform-style: preserve-3d;
  transition:
    transform 180ms cubic-bezier(0.2, 0.8, 0.2, 1),
    box-shadow 180ms cubic-bezier(0.2, 0.8, 0.2, 1);
  box-shadow:
    var(--terminal-shadow-x) var(--terminal-shadow-y) 70px rgba(15, 23, 42, 0.24),
    0 10px 24px rgba(15, 23, 42, 0.14);
  will-change: transform, box-shadow;
}

.terminal-window::before {
  content: '';
  position: absolute;
  inset: 0;
  z-index: 1;
  pointer-events: none;
  background:
    radial-gradient(
      circle at var(--terminal-glow-x) var(--terminal-glow-y),
      rgba(255, 255, 255, 0.16),
      rgba(255, 255, 255, 0.04) 28%,
      transparent 55%
    );
  opacity: 0;
  transition: opacity 180ms ease-out;
}

.terminal-showcase:hover .terminal-window::before {
  opacity: 1;
}

.terminal-window > * {
  position: relative;
  z-index: 2;
  transform: translateZ(18px);
}

.features-lab {
  background:
    linear-gradient(180deg, rgba(255, 255, 255, 0.86), rgba(247, 247, 245, 0.92)),
    linear-gradient(rgba(37, 99, 235, 0.06) 1px, transparent 1px),
    linear-gradient(90deg, rgba(37, 99, 235, 0.06) 1px, transparent 1px);
  background-size:
    100% 100%,
    36px 36px,
    36px 36px;
}

.features-lab::before {
  content: '';
  position: absolute;
  inset: 0;
  pointer-events: none;
  background:
    radial-gradient(circle at 18% 16%, rgba(37, 99, 235, 0.12), transparent 26%),
    radial-gradient(circle at 86% 54%, rgba(22, 163, 74, 0.1), transparent 24%);
}

.feature-card::before {
  content: '';
  position: absolute;
  inset: 0;
  pointer-events: none;
  background: linear-gradient(135deg, rgba(37, 99, 235, 0.11), transparent 38%);
  opacity: 0;
  transition: opacity 200ms ease-out;
}

.feature-card:hover::before {
  opacity: 1;
}

.quality-panel::before {
  content: '';
  position: absolute;
  inset: 0;
  pointer-events: none;
  background:
    linear-gradient(rgba(148, 163, 184, 0.08) 1px, transparent 1px),
    linear-gradient(90deg, rgba(148, 163, 184, 0.08) 1px, transparent 1px);
  background-size: 34px 34px;
  mask-image: linear-gradient(to bottom, black 0%, transparent 86%);
}

.quality-panel::after {
  content: '';
  position: absolute;
  inset: -30%;
  pointer-events: none;
  background:
    radial-gradient(circle at 76% 18%, rgba(37, 99, 235, 0.26), transparent 24%),
    radial-gradient(circle at 28% 78%, rgba(22, 163, 74, 0.18), transparent 24%);
  opacity: 0.86;
}

.quality-score-orb {
  background:
    radial-gradient(circle at 50% 42%, rgba(37, 99, 235, 0.26), transparent 46%),
    linear-gradient(135deg, rgba(255, 255, 255, 0.08), rgba(255, 255, 255, 0.02));
}

.cta-command {
  background:
    radial-gradient(circle at 18% 0%, rgba(37, 99, 235, 0.3), transparent 26%),
    linear-gradient(rgba(255, 255, 255, 0.055) 1px, transparent 1px),
    linear-gradient(90deg, rgba(255, 255, 255, 0.055) 1px, transparent 1px),
    #111827;
  background-size:
    100% 100%,
    40px 40px,
    40px 40px,
    100% 100%;
}

.cta-command::before {
  content: '';
  position: absolute;
  inset: 0;
  pointer-events: none;
  background: linear-gradient(180deg, rgba(21, 26, 33, 0.08), rgba(17, 24, 39, 0.88));
}

@media (max-width: 479px) {
  :global(html),
  :global(body) {
    overflow-x: hidden;
  }

  .terminal-window {
    border-radius: 1.25rem;
  }

  .landing-language-switcher :deep(svg) {
    display: none;
  }

  .landing-language-switcher :deep(button) {
    min-width: 2rem;
    padding-inline: 0.4rem;
  }

  .quality-panel {
    padding: 1rem;
  }
}

@media (prefers-reduced-motion: reduce) {
  .landing-grid-glow {
    display: none;
  }

  .landing-nav,
  .landing-nav__inner,
  .terminal-window,
  .terminal-window::before {
    transition: none;
  }

  .terminal-window {
    transform: none;
  }

  .feature-card {
    transform: none;
  }
}
</style>
