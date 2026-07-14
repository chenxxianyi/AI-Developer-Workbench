/**
 * Vue Router — 统一路由配置
 * 使用 createWebHistory，包含路由守卫
 */
import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  // ── 公共 ──
  {
    path: '/',
    name: 'landing',
    component: () => import('@/pages/LandingPage.vue'),
    meta: { layout: 'landing', requiresAuth: false },
  },
  {
    path: '/login',
    name: 'login',
    component: () => import('@/pages/LoginPage.vue'),
    meta: { layout: 'landing', requiresAuth: false },
  },
  {
    path: '/register',
    name: 'register',
    component: () => import('@/pages/RegisterPage.vue'),
    meta: { layout: 'landing', requiresAuth: false },
  },

  // ── 认证后 ──
  {
    path: '/dashboard',
    name: 'dashboard',
    component: () => import('@/pages/DashboardPage.vue'),
    meta: { layout: 'app', requiresAuth: true },
  },
  {
    path: '/profile',
    name: 'profile',
    component: () => import('@/pages/ProfilePage.vue'),
    meta: { layout: 'app', requiresAuth: true },
  },
  {
    path: '/settings',
    name: 'settings',
    component: () => import('@/pages/SettingsPage.vue'),
    meta: { layout: 'app', requiresAuth: true },
  },

  // ── 项目 ──
  {
    path: '/projects',
    name: 'projects',
    component: () => import('@/pages/ProjectsPage.vue'),
    meta: { layout: 'app', requiresAuth: true },
  },
  {
    path: '/projects/new',
    name: 'project-create',
    component: () => import('@/pages/ProjectCreatePage.vue'),
    meta: { layout: 'app', requiresAuth: true },
  },

  // ── 项目工作区 (嵌套路由) ──
  {
    path: '/projects/:projectId',
    component: () => import('@/components/layout/ProjectWorkspaceLayout.vue'),
    meta: { layout: 'app', requiresAuth: true },
    children: [
      { path: '', name: 'project-overview', component: () => import('@/pages/ProjectDetailPage.vue') },
      { path: 'edit', name: 'project-edit', component: () => import('@/pages/ProjectFormPage.vue') },
      { path: 'requirements', name: 'project-requirements', component: () => import('@/pages/ProjectRequirementsPage.vue') },
      { path: 'blueprint', name: 'project-blueprint', component: () => import('@/pages/BlueprintPage.vue') },
      { path: 'generation', name: 'project-generation', component: () => import('@/pages/GenerationPage.vue') },
      { path: 'preview', name: 'project-preview', component: () => import('@/pages/PreviewPage.vue') },
      { path: 'files', name: 'project-files', component: () => import('@/pages/FilesPage.vue') },
      { path: 'reports', name: 'project-reports', component: () => import('@/pages/ReportsPage.vue') },
    ],
  },

  // ── AI 工具 ──
  {
    path: '/tools/ui-review',
    name: 'ui-review',
    component: () => import('@/pages/tools/UIReviewPage.vue'),
    meta: { layout: 'app', requiresAuth: true },
  },
  {
    path: '/tools/project-doctor',
    name: 'project-doctor',
    component: () => import('@/pages/tools/ProjectDoctorPage.vue'),
    meta: { layout: 'app', requiresAuth: true },
  },
  {
    path: '/tools/agent-config',
    name: 'agent-config',
    component: () => import('@/pages/tools/AgentConfigPage.vue'),
    meta: { layout: 'app', requiresAuth: true },
  },
  {
    path: '/tools/api-doc',
    name: 'api-doc',
    component: () => import('@/pages/tools/APIDocPage.vue'),
    meta: { layout: 'app', requiresAuth: true },
  },
  {
    path: '/tools/db-schema',
    name: 'db-schema',
    component: () => import('@/pages/tools/DBSchemaPage.vue'),
    meta: { layout: 'app', requiresAuth: true },
  },

  // ── 报告 ──
  {
    path: '/reports',
    name: 'reports',
    component: () => import('@/pages/ReportsPage.vue'),
    meta: { layout: 'app', requiresAuth: true },
  },
  {
    path: '/reports/:id',
    name: 'report-detail',
    component: () => import('@/pages/ReportDetailPage.vue'),
    meta: { layout: 'app', requiresAuth: true },
  },
  {
    path: '/reports/:id/compare/:targetId',
    name: 'report-compare',
    component: () => import('@/pages/ReportComparePage.vue'),
    meta: { layout: 'app', requiresAuth: true },
  },

  // ── 管理后台 ──
  {
    path: '/admin',
    meta: { layout: 'app', requiresAuth: true, requiresAdmin: true },
    children: [
      { path: '', redirect: '/admin/models' },
      { path: 'models', name: 'admin-models', component: () => import('@/pages/admin/ModelsPage.vue') },
      { path: 'prompts', name: 'admin-prompts', component: () => import('@/pages/admin/PromptsPage.vue') },
      { path: 'users', name: 'admin-users', component: () => import('@/pages/admin/UsersPage.vue') },
      { path: 'projects', name: 'admin-projects', component: () => import('@/pages/admin/ProjectsPage.vue') },
    ],
  },

  // ── 错误 ──
  {
    path: '/403',
    name: 'forbidden',
    component: () => import('@/pages/ForbiddenPage.vue'),
    meta: { layout: 'app', requiresAuth: false },
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'not-found',
    component: () => import('@/pages/NotFoundPage.vue'),
    meta: { layout: 'app', requiresAuth: false },
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to) {
    if (to.hash) return { el: to.hash, top: 96, behavior: 'smooth' }
    return { top: 0 }
  },
})

// ── 路由守卫 ──
router.beforeEach((to, _from) => {
  void to // 路由守卫待 authStore 接入后启用
  // const authStore = useAuthStore()
  // if (to.meta.requiresAuth !== false && !authStore.isLoggedIn) {
  //   return { name: 'login', query: { redirect: to.fullPath } }
  // }
  // if (to.name === 'login' && authStore.isLoggedIn) {
  //   return { name: 'dashboard' }
  // }
  // if (to.meta.requiresAdmin && authStore.user?.role !== 'admin') {
  //   return { name: 'forbidden' }
  // }
})

export default router
