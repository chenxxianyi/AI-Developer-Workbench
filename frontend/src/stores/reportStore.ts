/**
 * Report Store
 * Report list with pagination, filtering, sorting, and detail cache
 */

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Report, ReportListParams, ReportStatus } from '@/types/report'
import type { ToolType } from '@/types/tool'
import { listReports, getReport, deleteReport } from '@/api/reports'
import type { PaginatedData } from '@/types/api'

export const useReportStore = defineStore('report', () => {
  // State - List
  const reports = ref<Report[]>([])
  const total = ref(0)
  const currentPage = ref(1)
  const pageSize = ref(10)

  const toolFilter = ref<ToolType | ''>('')
  const statusFilter = ref<ReportStatus | ''>('')
  const sort = ref<'newest' | 'oldest' | 'score_desc' | 'score_asc'>('newest')

  const loading = ref(false)
  const error = ref<string | null>(null)

  // State - Detail
  const currentReport = ref<Report | null>(null)
  const detailLoading = ref(false)
  const detailError = ref<string | null>(null)

  // Actions - List
  async function fetchReports() {
    loading.value = true
    error.value = null

    try {
      const params: ReportListParams = {
        page: currentPage.value,
        page_size: pageSize.value,
        sort: sort.value,
      }

      if (toolFilter.value) {
        params.tool_type = toolFilter.value
      }

      if (statusFilter.value) {
        params.status = statusFilter.value
      }

      const response: PaginatedData<Report> = await listReports(params)
      reports.value = response.items
      total.value = response.total

      // Handle empty page: decrement and re-fetch
      if (response.items.length === 0 && currentPage.value > 1) {
        currentPage.value--
        await fetchReports()
      }
    } catch (err: any) {
      error.value = err.message || '获取报告列表失败'
      reports.value = []
      total.value = 0
    } finally {
      loading.value = false
    }
  }

  function setFilter(filter: ToolType | '') {
    toolFilter.value = filter
    currentPage.value = 1
    fetchReports()
  }

  function setStatusFilter(filter: ReportStatus | '') {
    statusFilter.value = filter
    currentPage.value = 1
    fetchReports()
  }

  function setSort(newSort: 'newest' | 'oldest' | 'score_desc' | 'score_asc') {
    sort.value = newSort
    fetchReports()
  }

  function setPage(page: number) {
    currentPage.value = page
    fetchReports()
  }

  // Actions - Detail
  async function fetchReport(id: string) {
    detailLoading.value = true
    detailError.value = null

    try {
      currentReport.value = await getReport(id)
    } catch (err: any) {
      detailError.value = err.message || '获取报告详情失败'
      currentReport.value = null
    } finally {
      detailLoading.value = false
    }
  }

  async function deleteReportById(id: string) {
    try {
      await deleteReport(id)

      // Update list if current report is deleted
      if (currentReport.value?.id === id) {
        currentReport.value = null
      }

      // Refresh list
      await fetchReports()
    } catch (err: any) {
      error.value = err.message || '删除报告失败'
      throw err
    }
  }

  // Getters
  const hasNextPage = computed(() => {
    const totalPages = Math.ceil(total.value / pageSize.value)
    return currentPage.value < totalPages
  })

  const hasPrevPage = computed(() => currentPage.value > 1)

  const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

  return {
    // State - List
    reports,
    total,
    currentPage,
    pageSize,
    toolFilter,
    statusFilter,
    sort,
    loading,
    error,

    // State - Detail
    currentReport,
    detailLoading,
    detailError,

    // Actions - List
    fetchReports,
    setFilter,
    setStatusFilter,
    setSort,
    setPage,

    // Actions - Detail
    fetchReport,
    deleteReportById,

    // Getters
    hasNextPage,
    hasPrevPage,
    totalPages,
  }
})