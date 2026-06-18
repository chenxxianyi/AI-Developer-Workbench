/**
 * Tool Store
 * Tool metadata from backend API
 */

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { ToolType, ToolMeta } from '@/types/tool'
import { getTools } from '@/api/tools'

export const useToolStore = defineStore('tool', () => {
  // State
  const tools = ref<ToolMeta[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  // Actions
  async function fetchTools() {
    loading.value = true
    error.value = null

    try {
      tools.value = await getTools()
    } catch (err: any) {
      error.value = err.message || '获取工具列表失败'
    } finally {
      loading.value = false
    }
  }

  // Getters
  function getToolByType(type: ToolType): ToolMeta | undefined {
    return tools.value.find((tool) => tool.tool_type === type)
  }

  const toolCount = computed(() => tools.value.length)

  return {
    // State
    tools,
    loading,
    error,

    // Actions
    fetchTools,

    // Getters
    getToolByType,
    toolCount,
  }
})