<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import {
  ArrowDown, ArrowRight, Bot, CheckCircle2, Code2, Database, Eye,
  FileCheck2, FileCode, FileText, FolderPlus, GitBranch, LayoutDashboard, Monitor,
  Package, Play, ShieldCheck, Sparkles, Stethoscope, WandSparkles,
  Wrench, Zap,
} from '@lucide/vue'
import LanguageSwitcher from '@/components/common/LanguageSwitcher.vue'
import { getLandingNavLinkClass } from '@/utils/landingNav'

const { t } = useI18n()
const route = useRoute()
const isLandingNavCompact = ref(false)
let landingNavScrollFrame: number | null = null

const heroLightingStyle = ref<Record<string, string>>({
  '--hero-light-x': '76%',
  '--hero-light-y': '24%',
  '--hero-light-opacity': '0.72',
})

const productTiltStyle = ref<Record<string, string>>({
  '--card-rotate-x': '1deg',
  '--card-rotate-y': '-2deg',
  '--card-translate-y': '0px',
  '--card-shadow-x': '0px',
  '--card-shadow-y': '36px',
  '--card-glare-x': '68%',
  '--card-glare-y': '18%',
  '--card-glare-opacity': '0.28',
})

const studioStages = computed(() => [
  { icon: FolderPlus, label: t('landing.studio.stages.project'), note: '01' },
  { icon: FileText, label: t('landing.studio.stages.requirements'), note: '02' },
  { icon: GitBranch, label: t('landing.studio.stages.blueprint'), note: '03' },
  { icon: Code2, label: t('landing.studio.stages.generation'), note: '04' },
  { icon: Monitor, label: t('landing.studio.stages.preview'), note: '05' },
  { icon: Package, label: t('landing.studio.stages.delivery'), note: '06' },
])

const tools = computed(() => [
  { name: t('landing.tools.items.uiReview.name'), description: t('landing.tools.items.uiReview.description'), icon: Eye, route: '/tools/ui-review', tone: 'blue' },
  { name: t('landing.tools.items.projectDoctor.name'), description: t('landing.tools.items.projectDoctor.description'), icon: Stethoscope, route: '/tools/project-doctor', tone: 'green' },
  { name: t('landing.tools.items.agentConfig.name'), description: t('landing.tools.items.agentConfig.description'), icon: Bot, route: '/tools/agent-config', tone: 'amber' },
  { name: t('landing.tools.items.apiDoc.name'), description: t('landing.tools.items.apiDoc.description'), icon: FileText, route: '/tools/api-doc', tone: 'violet' },
  { name: t('landing.tools.items.dbSchema.name'), description: t('landing.tools.items.dbSchema.description'), icon: Database, route: '/tools/db-schema', tone: 'rose' },
])

const workflowSteps = computed(() => [
  { icon: FileText, title: t('landing.workflow.steps.brief.title'), description: t('landing.workflow.steps.brief.description') },
  { icon: GitBranch, title: t('landing.workflow.steps.blueprint.title'), description: t('landing.workflow.steps.blueprint.description') },
  { icon: WandSparkles, title: t('landing.workflow.steps.generate.title'), description: t('landing.workflow.steps.generate.description') },
  { icon: FileCheck2, title: t('landing.workflow.steps.ship.title'), description: t('landing.workflow.steps.ship.description') },
])


function prefersReducedMotion() {
  return window.matchMedia('(prefers-reduced-motion: reduce)').matches
}

function handleHeroPointerMove(event: PointerEvent) {
  if (event.pointerType === 'touch' || prefersReducedMotion()) return

  const target = event.currentTarget as HTMLElement
  const rect = target.getBoundingClientRect()
  const x = ((event.clientX - rect.left) / rect.width) * 100
  const y = ((event.clientY - rect.top) / rect.height) * 100

  heroLightingStyle.value = {
    '--hero-light-x': `${Math.max(0, Math.min(100, x)).toFixed(2)}%`,
    '--hero-light-y': `${Math.max(0, Math.min(100, y)).toFixed(2)}%`,
    '--hero-light-opacity': '0.92',
  }
}

function resetHeroLighting() {
  heroLightingStyle.value = {
    '--hero-light-x': '76%',
    '--hero-light-y': '24%',
    '--hero-light-opacity': '0.72',
  }
}

function handleProductPointerMove(event: PointerEvent) {
  if (event.pointerType === 'touch' || prefersReducedMotion()) return

  const target = event.currentTarget as HTMLElement
  const rect = target.getBoundingClientRect()
  const x = (event.clientX - rect.left) / rect.width - 0.5
  const y = (event.clientY - rect.top) / rect.height - 0.5

  productTiltStyle.value = {
    '--card-rotate-x': `${(-y * 9).toFixed(2)}deg`,
    '--card-rotate-y': `${(x * 11).toFixed(2)}deg`,
    '--card-translate-y': '-8px',
    '--card-shadow-x': `${(-x * 28).toFixed(1)}px`,
    '--card-shadow-y': `${(42 + Math.abs(y) * 18).toFixed(1)}px`,
    '--card-glare-x': `${((x + 0.5) * 100).toFixed(1)}%`,
    '--card-glare-y': `${((y + 0.5) * 100).toFixed(1)}%`,
    '--card-glare-opacity': '0.48',
  }
}

function resetProductTilt() {
  productTiltStyle.value = {
    '--card-rotate-x': '1deg',
    '--card-rotate-y': '-2deg',
    '--card-translate-y': '0px',
    '--card-shadow-x': '0px',
    '--card-shadow-y': '36px',
    '--card-glare-x': '68%',
    '--card-glare-y': '18%',
    '--card-glare-opacity': '0.28',
  }
}

function updateLandingNavState() { isLandingNavCompact.value = window.scrollY > 40 }
function handleLandingNavScroll() {
  if (landingNavScrollFrame !== null) return
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
  if (landingNavScrollFrame !== null) window.cancelAnimationFrame(landingNavScrollFrame)
})
</script>

<template>
  <div class="landing-page min-h-screen overflow-x-hidden bg-background text-text-primary">
    <nav
      class="landing-nav fixed z-50 rounded-2xl border border-border/80 bg-surface/90 shadow-[0_16px_48px_rgba(15,23,42,0.10)] backdrop-blur-xl"
      :class="{ 'landing-nav--compact': isLandingNavCompact }"
      aria-label="Primary navigation"
    >
      <div class="landing-nav__inner mx-auto flex max-w-content items-center justify-between gap-4 px-4 py-3 md:px-6">
        <RouterLink to="/" class="group flex min-w-0 items-center gap-3 rounded-lg" aria-label="AI Developer Workbench">
          <span class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-accent text-white shadow-[0_10px_24px_rgba(37,99,235,0.24)] transition-transform duration-200 group-hover:-translate-y-0.5"><Zap :size="21" /></span>
          <span class="hidden truncate text-lg font-bold tracking-[-0.02em] text-text-primary sm:block">AI Workbench</span>
        </RouterLink>

        <div class="hidden items-center gap-1 rounded-xl border border-border/70 bg-background/70 p-1 md:flex">
          <RouterLink to="/" :class="getLandingNavLinkClass(route.hash, '')">{{ t('landing.nav.home') }}</RouterLink>
          <RouterLink :to="{ path: '/', hash: '#studio' }" :class="getLandingNavLinkClass(route.hash, '#studio')">{{ t('landing.nav.studio') }}</RouterLink>
          <RouterLink :to="{ path: '/', hash: '#tools' }" :class="getLandingNavLinkClass(route.hash, '#tools')">{{ t('landing.nav.tools') }}</RouterLink>
          <RouterLink :to="{ path: '/', hash: '#workflow' }" :class="getLandingNavLinkClass(route.hash, '#workflow')">{{ t('landing.nav.workflow') }}</RouterLink>
        </div>

        <div class="flex shrink-0 items-center gap-2 sm:gap-3">
          <LanguageSwitcher />
          <RouterLink to="/projects/new" class="inline-flex min-h-10 items-center justify-center gap-2 rounded-xl bg-text-primary px-4 py-2 text-sm font-bold text-white transition-colors duration-200 hover:bg-accent sm:px-5">
            <Sparkles :size="16" /><span class="hidden sm:inline">{{ t('landing.nav.start') }}</span>
          </RouterLink>
        </div>
      </div>
    </nav>

    <main>
      <section
        class="hero-section relative px-4 pb-20 pt-32 md:px-8 md:pb-24 md:pt-40"
        :style="heroLightingStyle"
        @pointermove="handleHeroPointerMove"
        @pointerleave="resetHeroLighting"
        @pointercancel="resetHeroLighting"
      >
        <div class="hero-grid absolute inset-0" aria-hidden="true"></div>
        <div class="hero-interactive-light absolute inset-0" aria-hidden="true"></div>
        <div class="hero-orb hero-orb--one" aria-hidden="true"></div>
        <div class="hero-orb hero-orb--two" aria-hidden="true"></div>
        <div class="hero-orb hero-orb--three" aria-hidden="true"></div>

        <div class="relative mx-auto grid max-w-content items-center gap-14 lg:grid-cols-[minmax(0,0.92fr)_minmax(520px,1.08fr)] lg:gap-16">
          <div class="mx-auto max-w-3xl text-center lg:mx-0 lg:text-left">
            <div class="mb-6 inline-flex items-center gap-2 rounded-full border border-accent/20 bg-surface/80 px-3.5 py-1.5 text-sm font-semibold text-accent shadow-sm backdrop-blur">
              <Sparkles :size="15" /><span>{{ t('landing.hero.eyebrow') }}</span>
            </div>
            <h1 class="text-balance text-4xl font-bold leading-[1.08] tracking-[-0.045em] text-text-primary sm:text-5xl md:text-6xl lg:text-[68px]">
              {{ t('landing.hero.titleLead') }} <span class="hero-title-accent">{{ t('landing.hero.titleAccent') }}</span> {{ t('landing.hero.titleTail') }}
            </h1>
            <p class="mx-auto mt-7 max-w-2xl text-lg leading-8 text-text-secondary md:text-xl lg:mx-0">{{ t('landing.hero.description') }}</p>

            <div class="mt-9 flex flex-col justify-center gap-3 sm:flex-row lg:justify-start">
              <RouterLink to="/projects/new" class="inline-flex min-h-12 items-center justify-center gap-2 rounded-xl bg-accent px-6 py-3 font-bold text-white shadow-[0_14px_30px_rgba(37,99,235,0.22)] transition-all duration-200 hover:-translate-y-0.5 hover:bg-blue-700 hover:shadow-[0_18px_36px_rgba(37,99,235,0.28)]">
                <WandSparkles :size="19" />{{ t('landing.hero.studioCta') }}<ArrowRight :size="18" />
              </RouterLink>
              <a href="#tools" class="inline-flex min-h-12 items-center justify-center gap-2 rounded-xl border border-border bg-surface/80 px-6 py-3 font-bold text-text-primary shadow-sm backdrop-blur transition-colors duration-200 hover:border-accent/30 hover:bg-accent-soft">
                <Wrench :size="18" />{{ t('landing.hero.toolsCta') }}
              </a>
            </div>

            <div class="mt-9 flex flex-wrap justify-center gap-x-6 gap-y-3 text-sm font-medium text-text-muted lg:justify-start">
              <span class="inline-flex items-center gap-2"><CheckCircle2 :size="16" class="text-success" />{{ t('landing.hero.proof.blueprint') }}</span>
              <span class="inline-flex items-center gap-2"><CheckCircle2 :size="16" class="text-success" />{{ t('landing.hero.proof.preview') }}</span>
              <span class="inline-flex items-center gap-2"><CheckCircle2 :size="16" class="text-success" />{{ t('landing.hero.proof.quality') }}</span>
            </div>
          </div>

          <div
            class="hero-product-shell relative mx-auto w-full max-w-2xl lg:max-w-none"
            :style="productTiltStyle"
            @pointermove="handleProductPointerMove"
            @pointerleave="resetProductTilt"
            @pointercancel="resetProductTilt"
          >
            <div class="hero-product-glow" aria-hidden="true"></div>
            <div class="hero-floating-badge hero-floating-badge--top hidden sm:flex" aria-hidden="true">
              <CheckCircle2 :size="14" />
              <span>{{ t('landing.hero.proof.blueprint') }}</span>
            </div>
            <div class="hero-floating-badge hero-floating-badge--bottom hidden sm:flex" aria-hidden="true">
              <ShieldCheck :size="14" />
              <span>{{ t('landing.hero.proof.quality') }}</span>
            </div>
            <div class="hero-product-card relative overflow-hidden rounded-[26px] border border-slate-700/80 bg-[#111827]">
              <div class="hero-card-glare" aria-hidden="true"></div>
              <div class="flex items-center justify-between border-b border-white/10 px-5 py-4 sm:px-6">
                <div class="flex min-w-0 items-center gap-3">
                  <span class="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-blue-500/15 text-blue-300"><WandSparkles :size="18" /></span>
                  <div class="min-w-0">
                    <p class="truncate text-sm font-semibold text-white">{{ t('landing.hero.preview.title') }}</p>
                    <p class="truncate text-xs text-slate-400">{{ t('landing.hero.preview.subtitle') }}</p>
                  </div>
                </div>
                <span class="inline-flex items-center gap-2 rounded-full border border-emerald-400/20 bg-emerald-400/10 px-3 py-1 text-xs font-semibold text-emerald-300">
                  <span class="h-1.5 w-1.5 rounded-full bg-emerald-300"></span>{{ t('landing.hero.preview.running') }}
                </span>
              </div>

              <div class="grid sm:grid-cols-[190px_1fr]">
                <aside class="border-b border-white/10 bg-white/[0.025] p-4 sm:border-b-0 sm:border-r sm:p-5">
                  <p class="mb-3 text-xs font-semibold uppercase tracking-[0.16em] text-slate-500">{{ t('landing.hero.preview.pipeline') }}</p>
                  <div class="grid grid-cols-3 gap-2 sm:grid-cols-1">
                    <div v-for="(stage, index) in studioStages" :key="stage.note" class="flex min-w-0 items-center gap-2.5 rounded-lg px-2.5 py-2 text-xs sm:text-sm" :class="index === 3 ? 'bg-blue-500/15 text-blue-200' : index < 3 ? 'text-slate-300' : 'text-slate-500'">
                      <span class="flex h-6 w-6 shrink-0 items-center justify-center rounded-md border" :class="index < 3 ? 'border-emerald-400/20 bg-emerald-400/10 text-emerald-300' : index === 3 ? 'border-blue-400/20 bg-blue-400/10 text-blue-300' : 'border-white/10 bg-white/[0.03]'">
                        <CheckCircle2 v-if="index < 3" :size="13" /><Play v-else-if="index === 3" :size="12" fill="currentColor" /><span v-else class="text-[10px] font-bold">{{ stage.note }}</span>
                      </span>
                      <span class="truncate">{{ stage.label }}</span>
                    </div>
                  </div>
                </aside>

                <div class="min-w-0 p-5 sm:p-6">
                  <div class="mb-5 flex flex-wrap items-center justify-between gap-3">
                    <div><p class="text-xs font-semibold uppercase tracking-[0.16em] text-blue-300">{{ t('landing.hero.preview.current') }}</p><h2 class="mt-1 text-lg font-bold text-white">{{ t('landing.hero.preview.generating') }}</h2></div>
                    <span class="font-mono text-sm font-semibold text-blue-300">72%</span>
                  </div>
                  <div class="mb-6 h-2 overflow-hidden rounded-full bg-white/10"><div class="generation-progress h-full w-[72%] rounded-full bg-gradient-to-r from-blue-500 to-cyan-300"></div></div>
                  <div class="code-panel rounded-xl border border-white/10 bg-black/25 p-4 font-mono text-xs leading-6 text-slate-400 sm:text-sm">
                    <p><span class="text-violet-300">const</span> <span class="text-blue-200">project</span> = <span class="text-amber-200">await</span> studio.generate&#40;&#41;</p>
                    <p><span class="text-emerald-300">✓</span> src/pages/Dashboard.vue</p>
                    <p><span class="text-emerald-300">✓</span> src/components/AppShell.vue</p>
                    <p><span class="text-blue-300">●</span> {{ t('landing.hero.preview.log') }}</p>
                  </div>
                  <div class="mt-5 grid grid-cols-3 gap-2">
                    <div class="hero-metric-card rounded-lg border border-white/10 bg-white/[0.035] px-3 py-3"><FileCode :size="16" class="mb-2 text-blue-300" /><p class="text-xs text-slate-500">{{ t('landing.hero.preview.files') }}</p><p class="mt-1 text-sm font-bold text-white">24</p></div>
                    <div class="hero-metric-card rounded-lg border border-white/10 bg-white/[0.035] px-3 py-3"><Monitor :size="16" class="mb-2 text-cyan-300" /><p class="text-xs text-slate-500">{{ t('landing.hero.preview.pages') }}</p><p class="mt-1 text-sm font-bold text-white">8</p></div>
                    <div class="hero-metric-card rounded-lg border border-white/10 bg-white/[0.035] px-3 py-3"><ShieldCheck :size="16" class="mb-2 text-emerald-300" /><p class="text-xs text-slate-500">{{ t('landing.hero.preview.checks') }}</p><p class="mt-1 text-sm font-bold text-white">12/12</p></div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      <section id="studio" class="scroll-mt-28 border-y border-border bg-surface px-4 py-16 md:px-8 md:py-20">
        <div class="mx-auto max-w-content">
          <div class="mb-9 grid gap-4 lg:grid-cols-[0.85fr_1fr] lg:items-end">
            <div>
              <p class="section-kicker"><Sparkles :size="15" />{{ t('landing.overview.kicker') }}</p>
              <h2 class="mt-3 text-3xl font-bold tracking-[-0.035em] text-text-primary md:text-4xl">{{ t('landing.overview.title') }}</h2>
            </div>
            <p class="max-w-xl text-base leading-7 text-text-secondary lg:justify-self-end">{{ t('landing.overview.description') }}</p>
          </div>

          <div class="grid gap-5 lg:grid-cols-[1.3fr_0.7fr]">
            <article class="studio-card relative overflow-hidden rounded-[24px] bg-[#111827] p-6 text-white shadow-[0_24px_60px_rgba(15,23,42,0.16)] md:p-8">
              <div class="studio-card-grid absolute inset-0" aria-hidden="true"></div>
              <div class="relative">
                <div class="flex flex-col gap-6 sm:flex-row sm:items-start sm:justify-between">
                  <div class="max-w-xl">
                    <span class="inline-flex items-center gap-2 rounded-full border border-blue-300/20 bg-blue-300/10 px-3 py-1 text-sm font-semibold text-blue-200"><WandSparkles :size="15" />{{ t('landing.studio.badge') }}</span>
                    <h3 class="mt-4 text-2xl font-bold tracking-[-0.03em] md:text-3xl">{{ t('landing.studio.title') }}</h3>
                    <p class="mt-3 max-w-xl text-sm leading-6 text-slate-300">{{ t('landing.studio.description') }}</p>
                  </div>
                  <RouterLink to="/projects/new" class="inline-flex min-h-11 shrink-0 items-center justify-center gap-2 rounded-xl bg-white px-5 py-2.5 text-sm font-bold text-slate-950 transition-colors duration-200 hover:bg-blue-100">
                    {{ t('landing.studio.cta') }}<ArrowRight :size="16" />
                  </RouterLink>
                </div>

                <div class="mt-7 grid grid-cols-2 gap-3 md:grid-cols-3">
                  <div v-for="stage in studioStages" :key="stage.note" class="group rounded-xl border border-white/10 bg-white/[0.045] p-3.5 transition-colors duration-200 hover:border-blue-300/30 hover:bg-blue-300/[0.08]">
                    <div class="mb-4 flex items-center justify-between gap-3">
                      <span class="flex h-9 w-9 items-center justify-center rounded-lg bg-white/10 text-blue-200"><component :is="stage.icon" :size="18" /></span>
                      <span class="font-mono text-xs font-bold text-slate-500">{{ stage.note }}</span>
                    </div>
                    <p class="text-sm font-semibold text-slate-100">{{ stage.label }}</p>
                  </div>
                </div>

                <div class="mt-5 flex flex-wrap gap-2 border-t border-white/10 pt-5 text-xs font-medium text-slate-300">
                  <span class="rounded-full border border-white/10 bg-white/[0.04] px-3 py-1.5">{{ t('landing.studio.outputs.source') }}</span>
                  <span class="rounded-full border border-white/10 bg-white/[0.04] px-3 py-1.5">{{ t('landing.studio.outputs.preview') }}</span>
                  <span class="rounded-full border border-white/10 bg-white/[0.04] px-3 py-1.5">{{ t('landing.studio.outputs.blueprint') }}</span>
                  <span class="rounded-full border border-white/10 bg-white/[0.04] px-3 py-1.5">{{ t('landing.studio.outputs.package') }}</span>
                </div>
              </div>
            </article>

            <article class="flex flex-col rounded-[24px] border border-border bg-background p-6 md:p-8">
              <div class="flex h-12 w-12 items-center justify-center rounded-xl border border-accent/15 bg-accent-soft text-accent"><Wrench :size="23" /></div>
              <span class="mt-6 text-sm font-bold uppercase tracking-[0.14em] text-accent">{{ t('landing.toolbox.badge') }}</span>
              <h3 class="mt-3 text-2xl font-bold tracking-[-0.025em] text-text-primary md:text-3xl">{{ t('landing.toolbox.title') }}</h3>
              <p class="mt-4 leading-7 text-text-secondary">{{ t('landing.toolbox.description') }}</p>
              <div class="mt-7 space-y-3">
                <div class="flex items-center gap-3 rounded-xl border border-border bg-surface px-4 py-3"><Eye :size="17" class="text-accent" /><span class="text-sm font-semibold text-text-primary">{{ t('landing.tools.items.uiReview.name') }}</span></div>
                <div class="flex items-center gap-3 rounded-xl border border-border bg-surface px-4 py-3"><Stethoscope :size="17" class="text-success" /><span class="text-sm font-semibold text-text-primary">{{ t('landing.tools.items.projectDoctor.name') }}</span></div>
                <div class="flex items-center gap-3 rounded-xl border border-border bg-surface px-4 py-3"><Bot :size="17" class="text-warning" /><span class="text-sm font-semibold text-text-primary">{{ t('landing.tools.items.agentConfig.name') }}</span></div>
              </div>
              <a href="#tools" class="mt-auto inline-flex min-h-11 items-center gap-2 pt-7 text-sm font-bold text-accent hover:text-blue-700">{{ t('landing.toolbox.cta') }}<ArrowDown :size="16" /></a>
            </article>
          </div>
        </div>
      </section>

      <section id="tools" class="scroll-mt-28 px-4 py-16 md:px-8 md:py-20">
        <div class="mx-auto max-w-content">
          <div class="mx-auto mb-9 max-w-2xl text-center">
            <p class="section-kicker justify-center"><Wrench :size="15" />{{ t('landing.tools.kicker') }}</p>
            <h2 class="mt-3 text-3xl font-bold tracking-[-0.035em] text-text-primary md:text-4xl">{{ t('landing.tools.title') }}</h2>
            <p class="mt-4 text-base leading-7 text-text-secondary">{{ t('landing.tools.description') }}</p>
          </div>

          <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-6">
            <RouterLink v-for="(tool, index) in tools" :key="tool.name" :to="tool.route" class="tool-card group flex min-h-[220px] flex-col rounded-2xl border border-border bg-surface p-5 shadow-sm transition-all duration-200 hover:-translate-y-1 hover:border-accent/30 hover:shadow-[0_18px_44px_rgba(15,23,42,0.10)]" :class="index < 3 ? 'lg:col-span-2' : 'lg:col-span-3'">
              <div class="flex items-start justify-between gap-4">
                <span class="tool-icon flex h-12 w-12 items-center justify-center rounded-xl border" :data-tone="tool.tone"><component :is="tool.icon" :size="22" /></span>
                <ArrowRight :size="18" class="text-text-muted transition-transform duration-200 group-hover:translate-x-1 group-hover:text-accent" />
              </div>
              <h3 class="mt-6 text-lg font-bold tracking-[-0.02em] text-text-primary">{{ tool.name }}</h3>
              <p class="mt-2 flex-1 text-sm leading-6 text-text-secondary">{{ tool.description }}</p>
              <span class="mt-5 text-sm font-bold text-accent">{{ t('landing.tools.open') }}</span>
            </RouterLink>
          </div>
        </div>
      </section>

      <section id="workflow" class="workflow-section scroll-mt-28 border-y border-border px-4 py-16 md:px-8 md:py-20">
        <div class="mx-auto max-w-content">
          <div class="mx-auto mb-9 max-w-2xl text-center">
            <p class="section-kicker justify-center"><GitBranch :size="15" />{{ t('landing.workflow.kicker') }}</p>
            <h2 class="mt-3 text-3xl font-bold tracking-[-0.035em] text-text-primary md:text-4xl">{{ t('landing.workflow.title') }}</h2>
            <p class="mt-4 text-base leading-7 text-text-secondary">{{ t('landing.workflow.description') }}</p>
          </div>

          <ol class="grid gap-3 sm:grid-cols-2 lg:grid-cols-4">
            <li v-for="(step, index) in workflowSteps" :key="step.title" class="workflow-card rounded-2xl border border-border bg-surface p-5 shadow-sm">
              <div class="flex items-center justify-between gap-3">
                <span class="flex h-10 w-10 items-center justify-center rounded-xl bg-accent-soft text-accent"><component :is="step.icon" :size="19" /></span>
                <span class="font-mono text-xs font-bold text-text-muted">0{{ index + 1 }}</span>
              </div>
              <h3 class="mt-5 text-lg font-bold text-text-primary">{{ step.title }}</h3>
              <p class="mt-2 text-sm leading-6 text-text-secondary">{{ step.description }}</p>
            </li>
          </ol>

          <div class="mt-7 text-center">
            <RouterLink to="/projects/new" class="inline-flex min-h-11 items-center gap-2 rounded-xl bg-text-primary px-5 py-2.5 text-sm font-bold text-white transition-colors duration-200 hover:bg-accent">{{ t('landing.workflow.cta') }}<ArrowRight :size="16" /></RouterLink>
          </div>
        </div>
      </section>

      <section class="px-4 pb-8 md:px-8 md:pb-12">
        <div class="cta-panel relative mx-auto max-w-content overflow-hidden rounded-[28px] bg-accent px-6 py-10 text-white shadow-[0_28px_70px_rgba(37,99,235,0.24)] md:px-10 md:py-12 lg:px-14">
          <div class="cta-grid absolute inset-0" aria-hidden="true"></div>
          <div class="relative grid gap-6 lg:grid-cols-[1fr_auto] lg:items-center">
            <div>
              <p class="inline-flex items-center gap-2 text-sm font-bold uppercase tracking-[0.14em] text-blue-100"><Zap :size="16" />{{ t('landing.cta.kicker') }}</p>
              <h2 class="mt-3 max-w-2xl text-3xl font-bold tracking-[-0.035em] md:text-4xl">{{ t('landing.cta.title') }}</h2>
              <p class="mt-3 max-w-xl text-base leading-7 text-blue-100">{{ t('landing.cta.description') }}</p>
            </div>
            <div class="flex flex-col gap-3 sm:flex-row lg:flex-col">
              <RouterLink to="/projects/new" class="inline-flex min-h-12 items-center justify-center gap-2 rounded-xl bg-white px-7 py-3 font-bold text-accent transition-colors duration-200 hover:bg-blue-50"><WandSparkles :size="18" />{{ t('landing.cta.primary') }}<ArrowRight :size="17" /></RouterLink>
              <RouterLink to="/dashboard" class="inline-flex min-h-12 items-center justify-center gap-2 rounded-xl border border-white/25 bg-white/10 px-7 py-3 font-bold text-white transition-colors duration-200 hover:bg-white/20"><LayoutDashboard :size="18" />{{ t('landing.cta.secondary') }}</RouterLink>
            </div>
          </div>
        </div>
      </section>
    </main>

    <footer class="px-4 py-8 md:px-8">
      <div class="mx-auto flex max-w-content flex-col items-center justify-between gap-5 border-t border-border pt-8 text-center sm:flex-row sm:text-left">
        <div class="flex items-center gap-3">
          <span class="flex h-9 w-9 items-center justify-center rounded-xl bg-text-primary text-white"><Zap :size="18" /></span>
          <div><p class="font-bold text-text-primary">AI Developer Workbench</p><p class="mt-0.5 text-xs text-text-muted">{{ t('landing.footer.tagline') }}</p></div>
        </div>
        <div class="flex items-center gap-5 text-sm font-semibold text-text-secondary">
          <a href="#studio" class="hover:text-accent">{{ t('landing.nav.studio') }}</a>
          <a href="#tools" class="hover:text-accent">{{ t('landing.nav.tools') }}</a>
          <RouterLink to="/dashboard" class="hover:text-accent">Dashboard</RouterLink>
        </div>
      </div>
    </footer>
  </div>
</template>

<style scoped>
:global(html) { scroll-behavior: smooth; }

.landing-nav {
  top: 16px;
  left: 16px;
  right: 16px;
  transition: top 220ms ease, box-shadow 220ms ease;
}
.landing-nav--compact { top: 8px; box-shadow: 0 12px 34px rgba(15, 23, 42, 0.13); }

.hero-section {
  --hero-light-x: 76%;
  --hero-light-y: 24%;
  --hero-light-opacity: 0.72;
  background:
    radial-gradient(circle at 16% 18%, rgba(37, 99, 235, 0.13), transparent 31%),
    radial-gradient(circle at 82% 32%, rgba(14, 165, 233, 0.11), transparent 30%),
    linear-gradient(180deg, #fbfdff 0%, var(--color-background) 78%);
}
.hero-grid {
  opacity: 0.5;
  background-image:
    linear-gradient(rgba(148, 163, 184, 0.16) 1px, transparent 1px),
    linear-gradient(90deg, rgba(148, 163, 184, 0.16) 1px, transparent 1px);
  background-position: 0 0, 0 0;
  background-size: 46px 46px;
  mask-image: linear-gradient(to bottom, black 0%, transparent 82%);
}
.hero-interactive-light {
  opacity: var(--hero-light-opacity);
  background:
    radial-gradient(circle at var(--hero-light-x) var(--hero-light-y), rgba(255, 255, 255, 0.95) 0%, rgba(125, 211, 252, 0.32) 12%, rgba(37, 99, 235, 0.14) 26%, transparent 48%);
  mix-blend-mode: screen;
  pointer-events: none;
  transition: opacity 260ms ease;
}
.hero-orb {
  position: absolute;
  width: 280px;
  height: 280px;
  border-radius: 9999px;
  filter: blur(80px);
  opacity: 0.2;
  pointer-events: none;
  will-change: transform;
}
.hero-orb--one { top: 90px; left: -100px; background: #60a5fa; animation: orb-drift-a 12s ease-in-out infinite; }
.hero-orb--two { right: -90px; bottom: 40px; background: #22d3ee; animation: orb-drift-b 14s ease-in-out infinite; }
.hero-orb--three { top: 42%; left: 42%; width: 190px; height: 190px; background: #a78bfa; opacity: 0.12; animation: orb-drift-c 16s ease-in-out infinite; }

.hero-title-accent { color: var(--color-accent); position: relative; white-space: nowrap; }
.hero-title-accent::after {
  content: '';
  position: absolute;
  left: 2%;
  right: 2%;
  bottom: -5px;
  height: 8px;
  border-radius: 9999px;
  background: linear-gradient(90deg, rgba(37, 99, 235, 0.2), rgba(14, 165, 233, 0.5));
  transform: skewX(-16deg);
  z-index: -1;
}

.hero-product-shell {
  --card-rotate-x: 1deg;
  --card-rotate-y: -2deg;
  --card-translate-y: 0px;
  --card-shadow-x: 0px;
  --card-shadow-y: 36px;
  --card-glare-x: 68%;
  --card-glare-y: 18%;
  --card-glare-opacity: 0.28;
  perspective: 1200px;
  transform-style: preserve-3d;
}
.hero-product-card {
  transform: translateY(var(--card-translate-y)) rotateX(var(--card-rotate-x)) rotateY(var(--card-rotate-y));
  transform-origin: center;
  transform-style: preserve-3d;
  box-shadow:
    var(--card-shadow-x) var(--card-shadow-y) 92px rgba(15, 23, 42, 0.28),
    0 1px 0 rgba(255, 255, 255, 0.06) inset;
  transition: transform 220ms ease, box-shadow 220ms ease;
  will-change: transform;
}
.hero-card-glare {
  position: absolute;
  inset: 0;
  z-index: 2;
  pointer-events: none;
  opacity: var(--card-glare-opacity);
  background:
    radial-gradient(circle at var(--card-glare-x) var(--card-glare-y), rgba(255, 255, 255, 0.34), rgba(255, 255, 255, 0.08) 18%, transparent 42%);
  mix-blend-mode: screen;
  transition: opacity 220ms ease;
}
.hero-product-card > :not(.hero-card-glare) { position: relative; z-index: 3; }
.hero-product-glow {
  position: absolute;
  inset: 12% 6% -4%;
  border-radius: 32px;
  background: linear-gradient(135deg, rgba(37, 99, 235, 0.24), rgba(34, 211, 238, 0.18));
  filter: blur(44px);
  transform: translateZ(-60px);
}
.hero-floating-badge {
  position: absolute;
  z-index: 8;
  align-items: center;
  gap: 7px;
  border: 1px solid rgba(255, 255, 255, 0.64);
  border-radius: 9999px;
  background: rgba(255, 255, 255, 0.78);
  padding: 9px 12px;
  color: #0f172a;
  font-size: 12px;
  font-weight: 800;
  box-shadow: 0 18px 42px rgba(15, 23, 42, 0.14);
  backdrop-filter: blur(16px) saturate(160%);
  -webkit-backdrop-filter: blur(16px) saturate(160%);
}
.hero-floating-badge--top { top: 9%; right: -18px; transform: translateZ(64px) rotate(2deg); }
.hero-floating-badge--bottom { left: -18px; bottom: 12%; transform: translateZ(72px) rotate(-2deg); }
.hero-metric-card {
  transform: translateZ(22px);
  transition: transform 200ms ease, border-color 200ms ease, background 200ms ease;
}
.hero-metric-card:hover {
  transform: translateZ(30px) translateY(-2px);
  border-color: rgba(125, 211, 252, 0.3);
  background: rgba(255, 255, 255, 0.065);
}
.generation-progress {
  box-shadow: 0 0 18px rgba(56, 189, 248, 0.52);
  animation: progress-pulse 2.4s ease-in-out infinite;
}

.section-kicker {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  color: var(--color-accent);
  font-size: 0.875rem;
  font-weight: 800;
  letter-spacing: 0.12em;
  text-transform: uppercase;
}
.studio-card-grid,
.cta-grid {
  background-image:
    linear-gradient(rgba(255, 255, 255, 0.045) 1px, transparent 1px),
    linear-gradient(90deg, rgba(255, 255, 255, 0.045) 1px, transparent 1px);
  background-size: 34px 34px;
  mask-image: linear-gradient(135deg, black, transparent 80%);
}

.tool-icon[data-tone='blue'] { border-color: rgba(37, 99, 235, 0.14); background: #dbeafe; color: #2563eb; }
.tool-icon[data-tone='green'] { border-color: rgba(22, 163, 74, 0.14); background: #dcfce7; color: #15803d; }
.tool-icon[data-tone='amber'] { border-color: rgba(217, 119, 6, 0.14); background: #fef3c7; color: #b45309; }
.tool-icon[data-tone='violet'] { border-color: rgba(124, 58, 237, 0.14); background: #ede9fe; color: #7c3aed; }
.tool-icon[data-tone='rose'] { border-color: rgba(225, 29, 72, 0.14); background: #ffe4e6; color: #e11d48; }

.workflow-section {
  background:
    radial-gradient(circle at 12% 18%, rgba(37, 99, 235, 0.08), transparent 24%),
    #f3f6fb;
}
.workflow-card { transition: border-color 200ms ease, transform 200ms ease, box-shadow 200ms ease; }
.workflow-card:hover {
  transform: translateY(-3px);
  border-color: rgba(37, 99, 235, 0.28);
  box-shadow: 0 16px 40px rgba(15, 23, 42, 0.08);
}

.cta-panel { isolation: isolate; }
.cta-panel::after {
  content: '';
  position: absolute;
  right: -110px;
  bottom: -180px;
  width: 420px;
  height: 420px;
  border-radius: 9999px;
  background: rgba(255, 255, 255, 0.16);
  filter: blur(12px);
}

@keyframes progress-pulse {
  0%, 100% { opacity: 0.86; }
  50% { opacity: 1; }
}
@keyframes orb-drift-a {
  0%, 100% { transform: translate3d(0, 0, 0) scale(1); }
  50% { transform: translate3d(34px, -18px, 0) scale(1.08); }
}
@keyframes orb-drift-b {
  0%, 100% { transform: translate3d(0, 0, 0) scale(1); }
  50% { transform: translate3d(-28px, 22px, 0) scale(1.06); }
}
@keyframes orb-drift-c {
  0%, 100% { transform: translate3d(0, 0, 0) scale(1); }
  50% { transform: translate3d(18px, 28px, 0) scale(1.12); }
}

@media (max-width: 1023px) {
  .hero-product-card { transform: none; }
  .hero-floating-badge { transform: none; }
}
@media (max-width: 639px) {
  .landing-nav { left: 10px; right: 10px; top: 10px; }
  .hero-title-accent { white-space: normal; }
  .code-panel { overflow-x: auto; }
}
@media (prefers-reduced-motion: reduce) {
  :global(html) { scroll-behavior: auto; }
  .generation-progress, .hero-orb { animation: none; }
  .workflow-card, .hero-product-card, .hero-floating-badge { transition: none; transform: none; }
}
</style>
