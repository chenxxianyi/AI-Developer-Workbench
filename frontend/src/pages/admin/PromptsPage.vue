<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between"><h1 class="text-2xl font-bold">Prompt 模板管理</h1><button class="px-4 py-2 rounded-lg bg-[var(--color-accent)] text-white text-sm">+ 新建模板</button></div>
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div v-for="p in prompts" :key="p.id" class="bg-[var(--color-surface)] rounded-lg border p-4">
        <div class="flex items-center justify-between mb-2"><h3 class="font-semibold">{{ p.name }}</h3><AppBadge :variant="p.status === 'active' ? 'success' : 'default'">{{ p.status === 'active' ? '启用' : '禁用' }}</AppBadge></div>
        <p class="text-xs text-[var(--color-text-muted)] mb-2">类型: {{ p.type }} · v{{ p.version }}</p>
        <div class="flex gap-2"><button class="text-xs text-[var(--color-accent)] hover:underline">编辑</button><button class="text-xs text-[var(--color-danger)] hover:underline">删除</button></div>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import AppBadge from '@/components/common/AppBadge.vue'
const prompts = ref([{ id: '1', name: 'UI Review Default', type: 'ui_review', status: 'active', version: 1 }])
</script>
