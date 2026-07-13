<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between"><h1 class="text-2xl font-bold">用户管理</h1></div>
    <div class="bg-[var(--color-surface)] rounded-lg border overflow-hidden"><table class="w-full text-sm"><thead><tr class="border-b bg-[var(--color-surface-muted)]"><th class="px-4 py-3 text-left">用户名</th><th class="px-4 py-3 text-left">邮箱</th><th class="px-4 py-3 text-left">角色</th><th class="px-4 py-3 text-left">状态</th><th class="px-4 py-3 text-right">操作</th></tr></thead><tbody><tr v-for="u in users" :key="u.id" class="border-b hover:bg-[var(--color-surface-muted)]/50"><td class="px-4 py-3 font-medium">{{ u.username }}</td><td class="px-4 py-3 text-[var(--color-text-muted)]">{{ u.email }}</td><td class="px-4 py-3"><AppBadge :variant="u.role === 'admin' ? 'danger' : 'default'">{{ u.role }}</AppBadge></td><td class="px-4 py-3">{{ u.status }}</td><td class="px-4 py-3 text-right"><button class="text-[var(--color-accent)] hover:underline text-xs">编辑</button></td></tr></tbody></table></div>
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import AppBadge from '@/components/common/AppBadge.vue'
const users = ref([{ id: '1', username: 'admin', email: 'admin@workbench.local', role: 'admin', status: 'active' }])
</script>
