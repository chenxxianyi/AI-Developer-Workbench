/**
 * Vue Router Configuration
 */

import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'landing',
    component: () => import('@/pages/LandingPage.vue'),
    meta: { layout: 'landing' },
  },
  {
    path: '/dashboard',
    name: 'dashboard',
    component: () => import('@/pages/DashboardPage.vue'),
    meta: { layout: 'app' },
  },
  {
    path: '/tools/ui-review',
    name: 'ui-review',
    component: () => import('@/pages/tools/UIReviewPage.vue'),
    meta: { layout: 'app' },
  },
  {
    path: '/tools/project-doctor',
    name: 'project-doctor',
    component: () => import('@/pages/tools/ProjectDoctorPage.vue'),
    meta: { layout: 'app' },
  },
  {
    path: '/tools/agent-config',
    name: 'agent-config',
    component: () => import('@/pages/tools/AgentConfigPage.vue'),
    meta: { layout: 'app' },
  },
  {
    path: '/tools/api-doc',
    name: 'api-doc',
    component: () => import('@/pages/tools/APIDocPage.vue'),
    meta: { layout: 'app' },
  },
  {
    path: '/tools/db-schema',
    name: 'db-schema',
    component: () => import('@/pages/tools/DBSchemaPage.vue'),
    meta: { layout: 'app' },
  },
  {
    path: '/reports',
    name: 'reports',
    component: () => import('@/pages/ReportsPage.vue'),
    meta: { layout: 'app' },
  },
  {
    path: '/reports/:id',
    name: 'report-detail',
    component: () => import('@/pages/ReportDetailPage.vue'),
    meta: { layout: 'app' },
  },
  {
    path: '/settings',
    name: 'settings',
    component: () => import('@/pages/SettingsPage.vue'),
    meta: { layout: 'app' },
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/',
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior() {
    return { top: 0 } // Always scroll to top on navigation
  },
})

export default router