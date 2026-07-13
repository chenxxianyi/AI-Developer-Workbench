/**
 * Project Store
 * List, detail, create, update, delete, stats
 */

import { defineStore } from 'pinia'
import { ref } from 'vue'
import type {
  Project,
  ProjectSummary,
  ProjectStats,
  ProjectCreateInput,
  ProjectUpdateInput,
  ProjectListParams,
  ProjectDeleteResult,
} from '@/types/project'
import type { Report } from '@/types/report'
import type { PaginatedData } from '@/types/api'
import {
  listProjects,
  getProject,
  createProject,
  updateProject,
  deleteProject,
  getProjectStats,
  listProjectReports,
} from '@/api/projects'

export const useProjectStore = defineStore('project', () => {
  const projects = ref<ProjectSummary[]>([])
  const total = ref(0)
  const currentPage = ref(1)
  const pageSize = ref(10)
  const search = ref('')

  const loading = ref(false)
  const error = ref<string | null>(null)

  const currentProject = ref<Project | null>(null)
  const currentStats = ref<ProjectStats | null>(null)
  const detailLoading = ref(false)
  const detailError = ref<string | null>(null)
  const projectReports = ref<Report[]>([])
  const projectReportsTotal = ref(0)
  const projectReportPage = ref(1)
  const projectReportsLoading = ref(false)
  const projectReportsError = ref<string | null>(null)

  async function fetchProjects() {
    loading.value = true
    error.value = null
    try {
      const params: ProjectListParams = {
        page: currentPage.value,
        page_size: pageSize.value,
      }
      if (search.value) params.search = search.value
      const res: PaginatedData<ProjectSummary> = await listProjects(params)
      projects.value = res.items
      total.value = res.total
    } catch (err: any) {
      error.value = err.message || '获取项目列表失败'
      projects.value = []
      total.value = 0
    } finally {
      loading.value = false
    }
  }

  function setSearch(s: string) {
    search.value = s
    currentPage.value = 1
    fetchProjects()
  }

  function setPage(p: number) {
    currentPage.value = p
    fetchProjects()
  }

  async function fetchProject(id: string) {
    detailLoading.value = true
    detailError.value = null
    try {
      currentProject.value = await getProject(id)
    } catch (err: any) {
      currentProject.value = null
      detailError.value = err.message || '获取项目详情失败'
    } finally {
      detailLoading.value = false
    }
  }

  async function fetchStats(id: string) {
    try {
      currentStats.value = await getProjectStats(id)
    } catch (err: any) {
      currentStats.value = null
    }
  }

  async function fetchProjectReports(id: string, page = projectReportPage.value) {
    projectReportsLoading.value = true
    projectReportsError.value = null
    projectReportPage.value = page
    try {
      const result = await listProjectReports(id, {
        page,
        page_size: pageSize.value,
      })
      projectReports.value = result.items
      projectReportsTotal.value = result.total
    } catch (err: any) {
      projectReports.value = []
      projectReportsTotal.value = 0
      projectReportsError.value = err.message || '获取项目报告失败'
    } finally {
      projectReportsLoading.value = false
    }
  }

  async function create(data: ProjectCreateInput): Promise<Project> {
    const project = await createProject(data)
    await fetchProjects()
    return project
  }

  async function update(id: string, data: ProjectUpdateInput): Promise<Project> {
    const project = await updateProject(id, data)
    currentProject.value = project
    return project
  }

  async function remove(id: string): Promise<ProjectDeleteResult> {
    const result = await deleteProject(id)
    if (currentProject.value?.id === id) currentProject.value = null
    await fetchProjects()
    return result
  }

  const totalPages = ref(1)
  function recalcTotalPages() {
    totalPages.value = Math.ceil(total.value / pageSize.value) || 1
  }

  return {
    projects,
    total,
    currentPage,
    pageSize,
    search,
    loading,
    error,
    currentProject,
    currentStats,
    detailLoading,
    detailError,
    projectReports,
    projectReportsTotal,
    projectReportPage,
    projectReportsLoading,
    projectReportsError,
    totalPages,
    fetchProjects,
    setSearch,
    setPage,
    fetchProject,
    fetchStats,
    fetchProjectReports,
    create,
    update,
    remove,
    recalcTotalPages,
  }
})
