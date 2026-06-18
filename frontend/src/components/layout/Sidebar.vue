<script setup lang="ts">
/**
 * Sidebar Component
 * Desktop: Fixed left sidebar
 * Mobile: Drawer with overlay
 */

import { computed } from 'vue'
import { useRoute, RouterLink } from 'vue-router'
import { useSystemStore } from '@/stores/systemStore'
import {
  Zap,
  LayoutDashboard,
  Eye,
  Stethoscope,
  Bot,
  FileText,
  Database,
  FileStack,
  Settings,
  CheckCircle2,
  X,
} from '@lucide/vue'

const props = defineProps<{
  mobileOpen: boolean
}>()

const emit = defineEmits<{
  toggle: []
  close: []
}>()

const route = useRoute()
const systemStore = useSystemStore()

// Check if route is active
function isActive(path: string): boolean {
  return route.path === path
}

// Get tool type from route
const currentToolType = computed(() => {
  const path = route.path
  if (path.startsWith('/tools/ui-review')) return 'ui_review'
  if (path.startsWith('/tools/project-doctor')) return 'project_doctor'
  if (path.startsWith('/tools/agent-config')) return 'agent_config'
  if (path.startsWith('/tools/api-doc')) return 'api_doc'
  if (path.startsWith('/tools/db-schema')) return 'db_schema'
  return null
})

// System status
const systemHealthy = computed(() => systemStore.status?.healthy ?? false)
const providerInfo = computed(() => systemStore.providerInfo)
</script>

<template>
  <!-- Mobile Overlay -->
  <div
    v-if="props.mobileOpen"
    class="sidebar-overlay active fixed inset-0 bg-black/20 backdrop-blur-sm z-40 md:hidden"
    @click="emit('close')"
  ></div>

  <!-- Desktop Sidebar -->
  <aside
    class="hidden md:flex fixed left-0 top-0 bottom-0 w-64 bg-surface border-r border-border flex-col z-30"
  >
    <!-- Logo -->
    <div class="px-6 py-5 border-b border-border">
      <RouterLink to="/" class="flex items-center gap-3 hover:opacity-80 transition-smooth">
        <div class="w-8 h-8 bg-accent rounded-lg flex items-center justify-center">
          <Zap :size="20" class="text-white" />
        </div>
        <span class="text-lg font-semibold text-text-primary">AI Workbench</span>
      </RouterLink>
    </div>

    <!-- Navigation -->
    <nav class="flex-1 px-4 py-6 overflow-y-auto">
      <div class="space-y-2">
        <!-- Dashboard -->
        <RouterLink
          to="/dashboard"
          :class="[
            'flex items-center gap-3 px-3 py-2 rounded-lg transition-smooth',
            isActive('/dashboard')
              ? 'bg-accent-soft text-accent font-medium'
              : 'text-text-secondary hover:bg-surface-muted hover:text-text-primary',
          ]"
        >
          <LayoutDashboard :size="20" />
          <span>工作台</span>
        </RouterLink>

        <!-- Tools Section -->
        <div class="mt-6 mb-2 px-3">
          <span class="text-sm font-semibold text-text-muted uppercase tracking-wide">分析工具</span>
        </div>

        <!-- UI Review -->
        <RouterLink
          to="/tools/ui-review"
          :class="[
            'flex items-center gap-3 px-3 py-2 rounded-lg transition-smooth',
            currentToolType === 'ui_review'
              ? 'bg-accent-soft text-accent font-medium'
              : 'text-text-secondary hover:bg-surface-muted hover:text-text-primary',
          ]"
        >
          <Eye :size="20" />
          <span>UI 审查</span>
        </RouterLink>

        <!-- Project Doctor -->
        <RouterLink
          to="/tools/project-doctor"
          :class="[
            'flex items-center gap-3 px-3 py-2 rounded-lg transition-smooth',
            currentToolType === 'project_doctor'
              ? 'bg-accent-soft text-accent font-medium'
              : 'text-text-secondary hover:bg-surface-muted hover:text-text-primary',
          ]"
        >
          <Stethoscope :size="20" />
          <span>项目诊断</span>
        </RouterLink>

        <!-- Agent Config -->
        <RouterLink
          to="/tools/agent-config"
          :class="[
            'flex items-center gap-3 px-3 py-2 rounded-lg transition-smooth',
            currentToolType === 'agent_config'
              ? 'bg-accent-soft text-accent font-medium'
              : 'text-text-secondary hover:bg-surface-muted hover:text-text-primary',
          ]"
        >
          <Bot :size="20" />
          <span>Agent 配置</span>
        </RouterLink>

        <!-- API Doc -->
        <RouterLink
          to="/tools/api-doc"
          :class="[
            'flex items-center gap-3 px-3 py-2 rounded-lg transition-smooth',
            currentToolType === 'api_doc'
              ? 'bg-accent-soft text-accent font-medium'
              : 'text-text-secondary hover:bg-surface-muted hover:text-text-primary',
          ]"
        >
          <FileText :size="20" />
          <span>API 文档</span>
        </RouterLink>

        <!-- DB Schema -->
        <RouterLink
          to="/tools/db-schema"
          :class="[
            'flex items-center gap-3 px-3 py-2 rounded-lg transition-smooth',
            currentToolType === 'db_schema'
              ? 'bg-accent-soft text-accent font-medium'
              : 'text-text-secondary hover:bg-surface-muted hover:text-text-primary',
          ]"
        >
          <Database :size="20" />
          <span>数据库审查</span>
        </RouterLink>

        <!-- Reports Section -->
        <div class="mt-6 mb-2 px-3">
          <span class="text-sm font-semibold text-text-muted uppercase tracking-wide">报告管理</span>
        </div>

        <RouterLink
          to="/reports"
          :class="[
            'flex items-center gap-3 px-3 py-2 rounded-lg transition-smooth',
            isActive('/reports')
              ? 'bg-accent-soft text-accent font-medium'
              : 'text-text-secondary hover:bg-surface-muted hover:text-text-primary',
          ]"
        >
          <FileStack :size="20" />
          <span>历史报告</span>
        </RouterLink>

        <!-- System Section -->
        <div class="mt-6 mb-2 px-3">
          <span class="text-sm font-semibold text-text-muted uppercase tracking-wide">系统设置</span>
        </div>

        <RouterLink
          to="/settings"
          :class="[
            'flex items-center gap-3 px-3 py-2 rounded-lg transition-smooth',
            isActive('/settings')
              ? 'bg-accent-soft text-accent font-medium'
              : 'text-text-secondary hover:bg-surface-muted hover:text-text-primary',
          ]"
        >
          <Settings :size="20" />
          <span>系统状态</span>
        </RouterLink>
      </div>
    </nav>

    <!-- System Status -->
    <div class="px-4 py-4 border-t border-border">
      <div
        :class="[
          'px-3 py-2 rounded-lg',
          systemHealthy ? 'bg-success/10' : 'bg-danger/10',
        ]"
      >
        <div class="flex items-center gap-2 mb-1">
          <CheckCircle2
            :size="16"
            :class="systemHealthy ? 'text-success' : 'text-danger'"
          />
          <span
            :class="[
              'text-sm font-medium',
              systemHealthy ? 'text-success' : 'text-danger',
            ]"
          >
            {{ systemHealthy ? '服务正常' : '服务异常' }}
          </span>
        </div>
        <div class="text-xs text-text-muted">
          {{ providerInfo }}
        </div>
      </div>
    </div>
  </aside>

  <!-- Mobile Sidebar (Drawer) -->
  <aside
    :class="[
      'sidebar-mobile fixed left-0 top-0 bottom-0 w-64 bg-surface border-r border-border flex-col z-50 md:hidden',
      props.mobileOpen ? 'active' : '',
    ]"
    v-if="props.mobileOpen"
  >
    <!-- Logo + Close Button -->
    <div class="px-6 py-5 border-b border-border flex items-center justify-between">
      <RouterLink to="/" class="flex items-center gap-3" @click="emit('close')">
        <div class="w-8 h-8 bg-accent rounded-lg flex items-center justify-center">
          <Zap :size="20" class="text-white" />
        </div>
        <span class="text-lg font-semibold">AI Workbench</span>
      </RouterLink>
      <button
        @click="emit('close')"
        class="p-1 hover:bg-surface-muted rounded transition-smooth"
      >
        <X :size="20" />
      </button>
    </div>

    <!-- Mobile Navigation (Simplified) -->
    <nav class="flex-1 px-4 py-6 overflow-y-auto">
      <div class="space-y-2">
        <RouterLink
          to="/dashboard"
          @click="emit('close')"
          class="flex items-center gap-3 px-3 py-2 text-text-secondary hover:bg-surface-muted rounded-lg"
        >
          <LayoutDashboard :size="20" />
          <span>工作台</span>
        </RouterLink>

        <div class="mt-6 mb-2 px-3">
          <span class="text-sm font-semibold text-text-muted uppercase tracking-wide">分析工具</span>
        </div>

        <RouterLink
          to="/tools/ui-review"
          @click="emit('close')"
          class="flex items-center gap-3 px-3 py-2 text-text-secondary hover:bg-surface-muted rounded-lg"
        >
          <Eye :size="20" />
          <span>UI 审查</span>
        </RouterLink>

        <RouterLink
          to="/tools/project-doctor"
          @click="emit('close')"
          class="flex items-center gap-3 px-3 py-2 text-text-secondary hover:bg-surface-muted rounded-lg"
        >
          <Stethoscope :size="20" />
          <span>项目诊断</span>
        </RouterLink>

        <RouterLink
          to="/tools/agent-config"
          @click="emit('close')"
          class="flex items-center gap-3 px-3 py-2 text-text-secondary hover:bg-surface-muted rounded-lg"
        >
          <Bot :size="20" />
          <span>Agent 配置</span>
        </RouterLink>

        <RouterLink
          to="/tools/api-doc"
          @click="emit('close')"
          class="flex items-center gap-3 px-3 py-2 text-text-secondary hover:bg-surface-muted rounded-lg"
        >
          <FileText :size="20" />
          <span>API 文档</span>
        </RouterLink>

        <RouterLink
          to="/tools/db-schema"
          @click="emit('close')"
          class="flex items-center gap-3 px-3 py-2 text-text-secondary hover:bg-surface-muted rounded-lg"
        >
          <Database :size="20" />
          <span>数据库审查</span>
        </RouterLink>

        <div class="mt-6 mb-2 px-3">
          <span class="text-sm font-semibold text-text-muted uppercase tracking-wide">报告管理</span>
        </div>

        <RouterLink
          to="/reports"
          @click="emit('close')"
          class="flex items-center gap-3 px-3 py-2 text-text-secondary hover:bg-surface-muted rounded-lg"
        >
          <FileStack :size="20" />
          <span>历史报告</span>
        </RouterLink>

        <div class="mt-6 mb-2 px-3">
          <span class="text-sm font-semibold text-text-muted uppercase tracking-wide">系统设置</span>
        </div>

        <RouterLink
          to="/settings"
          @click="emit('close')"
          class="flex items-center gap-3 px-3 py-2 text-text-secondary hover:bg-surface-muted rounded-lg"
        >
          <Settings :size="20" />
          <span>系统状态</span>
        </RouterLink>
      </div>
    </nav>

    <!-- Mobile Status -->
    <div class="px-4 py-4 border-t border-border">
      <div
        :class="[
          'px-3 py-2 rounded-lg',
          systemHealthy ? 'bg-success/10' : 'bg-danger/10',
        ]"
      >
        <div class="flex items-center gap-2 mb-1">
          <CheckCircle2 :size="16" :class="systemHealthy ? 'text-success' : 'text-danger'" />
          <span :class="['text-sm font-medium', systemHealthy ? 'text-success' : 'text-danger']">
            {{ systemHealthy ? '服务正常' : '服务异常' }}
          </span>
        </div>
        <div class="text-xs text-text-muted">{{ providerInfo }}</div>
      </div>
    </div>
  </aside>
</template>

<style scoped>
/* All transitions handled by global .transition-smooth class */
</style>