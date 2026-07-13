<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between"><h1 class="text-2xl font-bold">AI 模型管理</h1><button class="px-4 py-2 rounded-lg bg-[var(--color-accent)] text-white text-sm">+ 添加模型</button></div>
    <div class="bg-[var(--color-surface)] rounded-lg border overflow-hidden"><table class="w-full text-sm"><thead><tr class="border-b bg-[var(--color-surface-muted)]"><th class="px-4 py-3 text-left">名称</th><th class="px-4 py-3 text-left">Provider</th><th class="px-4 py-3 text-left">状态</th><th class="px-4 py-3 text-right">操作</th></tr></thead><tbody><tr v-for="m in models" :key="m.id" class="border-b hover:bg-[var(--color-surface-muted)]/50"><td class="px-4 py-3 font-medium">{{ m.name }}</td><td class="px-4 py-3 text-[var(--color-text-muted)]">{{ m.provider }}</td><td class="px-4 py-3"><span :class="m.status === 'active' ? 'text-green-600' : 'text-red-500'">{{ m.status === 'active' ? '启用' : '禁用' }}</span></td><td class="px-4 py-3 text-right"><button class="text-[var(--color-accent)] hover:underline text-xs mr-2">编辑</button><button class="text-[var(--color-danger)] hover:underline text-xs">删除</button></td></tr></tbody></table></div>
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
const models = ref([{ id: '1', name: 'GPT-4.1', provider: 'openai', status: 'active' }, { id: '2', name: 'Claude', provider: 'anthropic', status: 'active' }])
</script>
