<script setup lang="ts">
/**
 * Sidebar — 统一导航
 * 分组：项目 | 生成工作室 | AI 工具 | 报告 | 管理后台 | 系统
 */
import { useRoute, RouterLink } from 'vue-router'
import {
  Zap, LayoutDashboard, FolderKanban, Eye, Stethoscope, Bot,
  FileText, Database, FileStack, Settings, X, LogOut,
  Wand2, PenTool, Play, Monitor, Files, User, Shield,
} from '@lucide/vue'
import { useAuthStore } from '@/stores/authStore'

const props = defineProps<{ mobileOpen: boolean }>()
const emit = defineEmits<{ toggle: []; close: [] }>()
const route = useRoute()
const auth = useAuthStore()

function isActive(path: string) { return route.path === path }
function isSameOrChildPath(path: string) {
  return route.path === path || route.path.startsWith(`${path}/`)
}

// Sidebar items definition
interface NavItem {
  to?: string
  icon: any
  label: string
  /**
   * Use route names for active matching to avoid broad path prefixes
   * (for example /admin or /projects/) lighting up multiple entries.
   */
  activeNames?: string[]
  admin?: boolean
}
interface NavSection { title: string; items: NavItem[]; admin?: boolean }

const sections: NavSection[] = [
  {
    title: '概览',
    items: [
      { to: '/dashboard', icon: LayoutDashboard, label: '工作台', activeNames: ['dashboard'] },
      {
        to: '/projects',
        icon: FolderKanban,
        label: '项目',
        activeNames: ['projects', 'project-overview', 'project-edit', 'project-requirements'],
      },
    ],
  },
  {
    title: '生成工作室',
    items: [
      { to: '/projects/new', icon: Wand2, label: '创建项目', activeNames: ['project-create'] },
      { to: '/projects/1/blueprint', icon: PenTool, label: '蓝图设计', activeNames: ['project-blueprint'] },
      { to: '/projects/1/generation', icon: Play, label: '代码生成', activeNames: ['project-generation'] },
      { to: '/projects/1/preview', icon: Monitor, label: '在线预览', activeNames: ['project-preview'] },
      { to: '/projects/1/files', icon: Files, label: '项目文件', activeNames: ['project-files'] },
    ],
  },
  {
    title: 'AI 开发工具',
    items: [
      { to: '/tools/ui-review', icon: Eye, label: 'UI 审查', activeNames: ['ui-review'] },
      { to: '/tools/project-doctor', icon: Stethoscope, label: '项目诊断', activeNames: ['project-doctor'] },
      { to: '/tools/agent-config', icon: Bot, label: 'Agent 配置', activeNames: ['agent-config'] },
      { to: '/tools/api-doc', icon: FileText, label: 'API 文档', activeNames: ['api-doc'] },
      { to: '/tools/db-schema', icon: Database, label: '数据库审查', activeNames: ['db-schema'] },
    ],
  },
  {
    title: '报告',
    items: [
      { to: '/reports', icon: FileStack, label: '历史报告', activeNames: ['reports', 'report-detail', 'report-compare'] },
    ],
  },
  {
    title: '管理后台',
    admin: true,
    items: [
      { to: '/admin/models', icon: Bot, label: 'AI 模型', activeNames: ['admin-models'] },
      { to: '/admin/prompts', icon: FileText, label: 'Prompt', activeNames: ['admin-prompts'] },
      { to: '/admin/users', icon: User, label: '用户管理', activeNames: ['admin-users'] },
      { to: '/admin/projects', icon: FolderKanban, label: '项目管理', activeNames: ['admin-projects'] },
    ],
  },
  {
    title: '系统',
    items: [
      { to: '/settings', icon: Settings, label: '系统状态', activeNames: ['settings'] },
    ],
  },
]

function itemActive(item: NavItem): boolean {
  const currentRouteName = typeof route.name === 'string' ? route.name : String(route.name ?? '')
  if (item.activeNames?.length) return item.activeNames.includes(currentRouteName)
  if (item.to) return isActive(item.to) || isSameOrChildPath(item.to)
  return false
}
</script>

<template>
  <!-- Mobile Overlay -->
  <div v-if="props.mobileOpen" class="sidebar-overlay active fixed inset-0 bg-black/20 backdrop-blur-sm z-40 md:hidden" @click="emit('close')" />

  <!-- Desktop -->
  <aside class="hidden md:flex fixed left-0 top-0 bottom-0 w-64 bg-surface border-r border-border flex-col z-30">
    <div class="px-6 py-5 border-b border-border">
      <RouterLink to="/" class="flex items-center gap-3 hover:opacity-80 transition-smooth">
        <div class="w-8 h-8 bg-accent rounded-lg flex items-center justify-center"><Zap :size="20" class="text-white" /></div>
        <span class="text-lg font-semibold text-text-primary">AI Workbench</span>
      </RouterLink>
    </div>

    <nav class="flex-1 px-4 py-4 overflow-y-auto">
      <template v-for="section in sections" :key="section.title">
        <div class="mt-4 mb-2 px-3">
          <span class="text-xs font-semibold text-text-muted uppercase tracking-wide flex items-center gap-1.5">
            {{ section.title }}
            <Shield v-if="section.admin" :size="12" class="text-warning" />
          </span>
        </div>
        <RouterLink
          v-for="item in section.items" :key="item.label"
          :to="item.to!"
          :class="[
            'flex items-center gap-3 px-3 py-2 rounded-lg transition-smooth mb-0.5',
            itemActive(item)
              ? 'bg-accent-soft text-accent font-medium'
              : 'text-text-secondary hover:bg-surface-muted hover:text-text-primary',
          ]"
        >
          <component :is="item.icon" :size="20" />
          <span>{{ item.label }}</span>
        </RouterLink>
      </template>
    </nav>

    <!-- User area -->
    <div class="px-4 py-3 border-t border-border">
      <template v-if="auth.isLoggedIn">
        <div class="flex items-center justify-between px-3 py-2">
          <div class="flex items-center gap-3 min-w-0">
            <div class="w-8 h-8 rounded-full bg-accent-soft flex items-center justify-center shrink-0">
              <span class="text-sm font-semibold text-accent">{{ auth.user?.username?.charAt(0).toUpperCase() }}</span>
            </div>
            <div class="min-w-0">
              <p class="text-sm font-medium text-text-primary truncate">{{ auth.user?.username }}</p>
              <p class="text-xs text-text-muted truncate">{{ auth.user?.role === 'admin' ? '管理员' : '用户' }}</p>
            </div>
          </div>
          <button @click="auth.logout(); window.location.href = '/login'" class="p-1.5 rounded-lg text-text-muted hover:text-danger hover:bg-danger/10 transition-smooth" title="退出登录">
            <LogOut :size="16" />
          </button>
        </div>
      </template>
      <RouterLink v-else to="/login" class="flex items-center gap-3 px-3 py-2 rounded-lg text-text-secondary hover:bg-surface-muted transition-smooth">
        <User :size="18" />
        <span class="text-sm">登录</span>
      </RouterLink>
    </div>
  </aside>

  <!-- Mobile -->
  <aside v-if="props.mobileOpen" class="sidebar-mobile active fixed left-0 top-0 bottom-0 w-64 bg-surface border-r border-border flex-col z-50 md:hidden">
    <div class="px-6 py-5 border-b border-border flex items-center justify-between">
      <RouterLink to="/" class="flex items-center gap-3" @click="emit('close')">
        <div class="w-8 h-8 bg-accent rounded-lg flex items-center justify-center"><Zap :size="20" class="text-white" /></div>
        <span class="text-lg font-semibold">AI Workbench</span>
      </RouterLink>
      <button @click="emit('close')" class="p-1 hover:bg-surface-muted rounded"><X :size="20" /></button>
    </div>
    <nav class="flex-1 px-4 py-4 overflow-y-auto">
      <template v-for="section in sections" :key="section.title">
        <div class="mt-4 mb-2 px-3"><span class="text-xs font-semibold text-text-muted uppercase">{{ section.title }}</span></div>
        <RouterLink
          v-for="item in section.items" :key="item.label" :to="item.to!" @click="emit('close')"
          :class="[
            'flex items-center gap-3 px-3 py-2 rounded-lg transition-smooth mb-0.5',
            itemActive(item)
              ? 'bg-accent-soft text-accent font-medium'
              : 'text-text-secondary hover:bg-surface-muted hover:text-text-primary',
          ]"
        >
          <component :is="item.icon" :size="20" /><span>{{ item.label }}</span>
        </RouterLink>
      </template>
    </nav>
    <div class="px-4 py-3 border-t border-border">
      <template v-if="auth.isLoggedIn">
        <div class="flex items-center justify-between px-3 py-2">
          <div class="flex items-center gap-3 min-w-0">
            <div class="w-8 h-8 rounded-full bg-accent-soft flex items-center justify-center shrink-0">
              <span class="text-sm font-semibold text-accent">{{ auth.user?.username?.charAt(0).toUpperCase() }}</span>
            </div>
            <div class="min-w-0">
              <p class="text-sm font-medium text-text-primary truncate">{{ auth.user?.username }}</p>
              <p class="text-xs text-text-muted truncate">{{ auth.user?.role === 'admin' ? '管理员' : '用户' }}</p>
            </div>
          </div>
          <button @click="auth.logout(); window.location.href = '/login'" class="p-1.5 rounded-lg text-text-muted hover:text-danger hover:bg-danger/10 transition-smooth" title="退出登录">
            <LogOut :size="16" />
          </button>
        </div>
      </template>
      <RouterLink v-else to="/login" @click="emit('close')" class="flex items-center gap-3 px-3 py-2 rounded-lg text-text-secondary hover:bg-surface-muted">
        <User :size="18" /><span class="text-sm">登录</span>
      </RouterLink>
    </div>
  </aside>
</template>
