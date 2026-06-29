<script setup lang="ts">
/**
 * Landing Page
 * Standalone layout (no AppShell)
 */

import { computed, ref } from 'vue'
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
  { title: t('landing.features.items.product.title'), description: t('landing.features.items.product.description') },
  { title: t('landing.features.items.safety.title'), description: t('landing.features.items.safety.description') },
  { title: t('landing.features.items.output.title'), description: t('landing.features.items.output.description') },
  { title: t('landing.features.items.integration.title'), description: t('landing.features.items.integration.description') },
])

const statItems = computed(() => [
  { icon: Wrench, label: t('landing.stats.tools') },
  { icon: FileCheck, label: t('landing.stats.analysis') },
  { icon: Download, label: t('landing.stats.export') },
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
</script>

<template>
  <div class="min-h-screen overflow-x-hidden bg-background text-text-primary">
    <!-- Navigation -->
    <nav class="fixed left-3 right-3 top-3 z-50 overflow-hidden rounded-xl border border-border bg-surface/90 shadow-sm backdrop-blur-md transition-smooth md:left-4 md:right-4 md:top-4">
      <div class="mx-auto flex max-w-content items-center justify-between gap-3 px-3 py-2.5 md:gap-4 md:px-5">
        <RouterLink to="/" class="flex min-w-0 items-center gap-3 transition-smooth hover:opacity-80">
          <div class="flex h-9 w-9 flex-shrink-0 items-center justify-center rounded-lg bg-accent">
            <Zap :size="20" class="text-white" />
          </div>
          <span class="hidden truncate text-lg font-semibold text-text-primary sm:inline">AI Workbench</span>
        </RouterLink>

        <div class="hidden items-center gap-1 md:flex">
          <RouterLink to="/" :class="getLandingNavLinkClass(route.hash, '')">{{ t('landing.nav.home') }}</RouterLink>
          <RouterLink :to="{ path: '/', hash: '#tools' }" :class="getLandingNavLinkClass(route.hash, '#tools')">{{ t('landing.nav.tools') }}</RouterLink>
          <RouterLink :to="{ path: '/', hash: '#workflow' }" :class="getLandingNavLinkClass(route.hash, '#workflow')">{{ t('landing.nav.workflow') }}</RouterLink>
          <RouterLink :to="{ path: '/', hash: '#features' }" :class="getLandingNavLinkClass(route.hash, '#features')">{{ t('landing.nav.features') }}</RouterLink>
        </div>

        <div class="flex flex-shrink-0 items-center gap-2">
          <div class="landing-language-switcher">
            <LanguageSwitcher />
          </div>
          <RouterLink
            to="/dashboard"
            class="hidden items-center gap-2 rounded-lg bg-accent px-4 py-2 text-sm font-semibold text-white shadow-sm transition-smooth hover:bg-accent/90 sm:inline-flex"
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
    <section id="features" class="border-y border-border bg-surface px-4 py-16 md:px-8 md:py-20">
      <div class="mx-auto max-w-content">
        <div class="grid gap-10 md:grid-cols-[0.9fr_1.1fr] md:items-center">
          <div>
            <h2 class="mb-6 text-3xl font-bold text-text-primary md:text-4xl">{{ t('landing.features.title') }}</h2>
            <div class="space-y-5">
              <div v-for="feature in features" :key="feature.title" class="flex items-start gap-3">
                <CheckCircle2 :size="24" class="mt-0.5 flex-shrink-0 text-success" />
                <div>
                  <h3 class="mb-1 font-semibold text-text-primary">{{ feature.title }}</h3>
                  <p class="leading-relaxed text-text-secondary">{{ feature.description }}</p>
                </div>
              </div>
            </div>
          </div>

          <div class="rounded-lg border border-border bg-background p-4 shadow-sm md:p-6">
            <div class="mb-4 flex items-center gap-2">
              <Terminal :size="20" class="text-accent" />
              <span class="font-semibold text-text-primary">{{ t('landing.sampleReport.title') }}</span>
            </div>
            <div class="rounded-lg bg-text-primary p-4 font-mono text-sm leading-relaxed text-white/85">
              <div class="mb-2 text-white/50"># UI Review Report</div>
              <div class="mb-3">
                <span class="font-semibold text-warning">{{ t('landing.sampleReport.score') }}:</span>
                <span> 78/100</span>
              </div>
              <div class="mb-3">
                <span class="font-semibold text-danger">{{ t('landing.sampleReport.issues') }}:</span>
                <span> 5</span>
              </div>
              <div class="text-white/70">
                <div class="mb-2">- {{ t('landing.sampleReport.high') }}</div>
                <div class="mb-2">- {{ t('landing.sampleReport.medium') }}</div>
                <div>- {{ t('landing.sampleReport.low') }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- CTA Section -->
    <section class="bg-text-primary px-4 py-16 text-white md:px-8 md:py-20">
      <div class="mx-auto max-w-content">
        <div class="grid gap-8 md:grid-cols-[1fr_auto] md:items-center">
          <div>
            <h2 class="mb-4 text-3xl font-bold md:text-4xl">{{ t('landing.cta.title') }}</h2>
            <p class="max-w-2xl text-lg leading-relaxed text-white/70">
              {{ t('landing.cta.description') }}
            </p>
          </div>
          <RouterLink
            to="/dashboard"
            class="inline-flex min-h-12 items-center justify-center gap-2 rounded-lg bg-white px-8 py-3 font-semibold text-text-primary shadow-sm transition-smooth hover:bg-background"
          >
            <LayoutDashboard :size="20" />
            {{ t('landing.cta.button') }}
            <ArrowRight :size="18" />
          </RouterLink>
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
}

@media (prefers-reduced-motion: reduce) {
  .landing-grid-glow {
    display: none;
  }

  .terminal-window,
  .terminal-window::before {
    transition: none;
  }

  .terminal-window {
    transform: none;
  }
}
</style>
