<script setup lang="ts">
/**
 * App Header Component
 * Sticky header with title and mobile menu toggle
 */

import { computed } from 'vue'
import { useRoute, RouterLink } from 'vue-router'
import { Menu, ArrowLeft, Trash2, ExternalLink } from '@lucide/vue'
import LanguageSwitcher from '@/components/common/LanguageSwitcher.vue'

const emit = defineEmits(['toggle-sidebar'])

const route = useRoute()

// Dynamic page title based on route
const pageTitle = computed(() => {
  const path = route.path

  if (path === '/dashboard') return '工作台'
  if (path === '/projects/new') return '创建项目'
  if (path.endsWith('/blueprint')) return '蓝图设计'
  if (path.endsWith('/generation')) return '代码生成'
  if (path.endsWith('/preview')) return '在线预览'
  if (path.endsWith('/files')) return '项目文件'
  if (path.startsWith('/projects/')) return '项目工作区'
  if (path === '/projects') return '项目'
  if (path.startsWith('/tools/ui-review')) return 'UI 审查'
  if (path.startsWith('/tools/project-doctor')) return '项目诊断'
  if (path.startsWith('/tools/agent-config')) return 'Agent 配置'
  if (path.startsWith('/tools/api-doc')) return 'API 文档'
  if (path.startsWith('/tools/db-schema')) return '数据库审查'
  if (path === '/reports') return '历史报告'
  if (path.startsWith('/reports/')) return '报告详情'
  if (path === '/admin/models') return 'AI 模型管理'
  if (path === '/admin/prompts') return 'Prompt 模板管理'
  if (path === '/admin/users') return '用户管理'
  if (path === '/admin/projects') return '项目管理'
  if (path === '/settings') return '系统状态'

  return 'AI Workbench'
})

// Show back button on detail pages
const showBackButton = computed(() => {
  return route.path.startsWith('/reports/') && route.path !== '/reports'
})

// Show delete button only on report detail
const showDeleteButton = computed(() => {
  return route.path.startsWith('/reports/') && route.params.id && route.path !== '/reports'
})
</script>

<template>
  <header class="sticky top-0 bg-surface/95 backdrop-blur-sm border-b border-border z-20">
    <div class="max-w-content mx-auto px-4 md:px-8 py-4 flex items-center justify-between">
      <!-- Mobile Menu Button -->
      <button
        @click="emit('toggle-sidebar')"
        class="md:hidden p-2 hover:bg-surface-muted rounded-lg transition-smooth"
      >
        <Menu :size="20" />
      </button>

      <!-- Page Title -->
      <div class="flex-1 md:flex-none">
        <h1 class="text-xl md:text-2xl font-bold text-text-primary">{{ pageTitle }}</h1>
      </div>

      <!-- Actions -->
      <div class="flex items-center gap-3">
        <LanguageSwitcher />

        <!-- Back Button -->
        <RouterLink
          v-if="showBackButton"
          to="/reports"
          class="hidden md:flex items-center gap-2 px-4 py-2 bg-surface-muted border border-border text-text-primary rounded-lg hover:bg-border transition-smooth"
        >
          <ArrowLeft :size="16" />
          <span>返回列表</span>
        </RouterLink>

        <!-- Delete Button (placeholder for now) -->
        <button
          v-if="showDeleteButton"
          class="hidden md:flex items-center gap-2 px-4 py-2 bg-danger/10 border border-danger/20 text-danger rounded-lg hover:bg-danger/20 transition-smooth"
        >
          <Trash2 :size="16" />
          <span>删除</span>
        </button>

        <!-- Dashboard Link -->
        <RouterLink
          v-if="route.path !== '/dashboard' && !route.path.startsWith('/reports/')"
          to="/dashboard"
          class="hidden md:flex items-center gap-2 px-4 py-2 bg-surface-muted border border-border text-text-primary rounded-lg hover:bg-border transition-smooth"
        >
          <ExternalLink :size="16" />
          <span>返回工作台</span>
        </RouterLink>
      </div>
    </div>
  </header>
</template>
