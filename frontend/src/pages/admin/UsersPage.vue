<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import {
  CircleCheck,
  CircleX,
  Pencil,
  Power,
  Search,
  ShieldCheck,
  UserCog,
  Users,
} from '@lucide/vue'
import AdminPageShell from '@/components/admin/AdminPageShell.vue'
import apiClient from '@/api/client'

interface AdminUser {
  id: string
  username: string
  email: string
  role: 'admin' | 'user'
  status: 'active' | 'disabled'
}

const search = ref('')
const roleFilter = ref<'' | AdminUser['role']>('')
const users = ref<AdminUser[]>([])
const loading = ref(false)
const error = ref('')
const updatingUserId = ref('')

async function loadUsers() {
  loading.value = true
  error.value = ''
  try {
    users.value = await apiClient.get('/admin/users') as unknown as AdminUser[]
  } catch (err: any) {
    error.value = err?.message || '获取用户列表失败'
  } finally {
    loading.value = false
  }
}

async function toggleUserStatus(user: AdminUser) {
  const status = user.status === 'active' ? 'disabled' : 'active'
  updatingUserId.value = user.id
  error.value = ''
  try {
    await apiClient.put(`/admin/users/${user.id}/status`, { status })
    user.status = status
  } catch (err: any) {
    error.value = err?.message || '更新用户状态失败'
  } finally {
    updatingUserId.value = ''
  }
}

onMounted(loadUsers)

const filteredUsers = computed(() => {
  const keyword = search.value.trim().toLowerCase()
  return users.value.filter((user) => {
    const matchesSearch = !keyword
      || user.username.toLowerCase().includes(keyword)
      || user.email.toLowerCase().includes(keyword)
    const matchesRole = !roleFilter.value || user.role === roleFilter.value
    return matchesSearch && matchesRole
  })
})

const adminCount = computed(() => users.value.filter((user) => user.role === 'admin').length)
</script>

<template>
  <AdminPageShell
    :icon="Users"
    title="用户管理"
    description="查看用户身份、权限角色和当前账号状态。"
    badge-text="权限管理"
  >
    <section class="mb-5 rounded-lg border border-border bg-surface p-4 shadow-sm">
      <div class="flex flex-col gap-3 lg:flex-row lg:items-center lg:justify-between">
        <div class="relative w-full lg:max-w-md">
          <Search :size="17" class="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-text-muted" />
          <input
            v-model="search"
            type="search"
            class="min-h-10 w-full rounded-lg border border-border bg-surface-muted pl-10 pr-4 text-sm text-text-primary placeholder:text-text-muted focus-visible:border-accent focus-visible:ring-2 focus-visible:ring-accent focus:outline-none"
            placeholder="搜索用户名或邮箱"
            aria-label="搜索用户"
          />
        </div>

        <div class="flex flex-col gap-3 sm:flex-row sm:items-center">
          <select
            v-model="roleFilter"
            class="min-h-10 rounded-lg border border-border bg-surface px-3 text-sm text-text-primary focus-visible:border-accent focus-visible:ring-2 focus-visible:ring-accent focus:outline-none"
            aria-label="按角色筛选用户"
          >
            <option value="">全部角色</option>
            <option value="admin">管理员</option>
            <option value="user">普通用户</option>
          </select>
          <span class="text-sm text-text-muted">
            {{ users.length }} 个用户，{{ adminCount }} 个管理员
          </span>
        </div>
      </div>
    </section>

    <p v-if="error" role="alert" class="mb-4 rounded-lg border border-danger/20 bg-danger/10 px-4 py-3 text-sm text-danger">{{ error }}</p>

    <section class="overflow-hidden rounded-lg border border-border bg-surface shadow-sm">
      <div class="flex items-center justify-between gap-3 border-b border-border px-5 py-4">
        <div>
          <h2 class="font-semibold text-text-primary">用户列表</h2>
          <p class="mt-0.5 text-xs text-text-muted">管理账号角色与访问状态</p>
        </div>
        <span class="rounded-full border border-border bg-surface-muted px-2.5 py-1 text-xs font-semibold text-text-secondary">
          {{ filteredUsers.length }} 项
        </span>
      </div>

      <div v-if="loading" class="flex min-h-64 items-center justify-center text-sm text-text-muted">正在读取数据库用户...</div>

      <div v-else-if="filteredUsers.length" class="overflow-x-auto">
        <table class="w-full min-w-[820px] text-sm">
          <thead class="bg-surface-muted/70 text-xs text-text-muted">
            <tr>
              <th class="px-5 py-3 text-left font-semibold">用户</th>
              <th class="px-5 py-3 text-left font-semibold">邮箱</th>
              <th class="px-5 py-3 text-left font-semibold">角色</th>
              <th class="px-5 py-3 text-left font-semibold">状态</th>
              <th class="px-5 py-3 text-right font-semibold">操作</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="user in filteredUsers"
              :key="user.id"
              class="border-t border-border transition-colors duration-200 hover:bg-surface-muted/50"
            >
              <td class="px-5 py-4">
                <div class="flex items-center gap-3">
                  <div class="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-accent-soft text-accent">
                    <UserCog :size="18" />
                  </div>
                  <div>
                    <p class="font-semibold text-text-primary">{{ user.username }}</p>
                    <p class="mt-0.5 font-mono text-xs text-text-muted">ID {{ user.id }}</p>
                  </div>
                </div>
              </td>
              <td class="px-5 py-4 text-text-secondary">{{ user.email }}</td>
              <td class="px-5 py-4">
                <span
                  :class="[
                    'inline-flex items-center gap-1.5 rounded-full border px-2.5 py-1 text-xs font-semibold',
                    user.role === 'admin'
                      ? 'border-danger/20 bg-danger/10 text-danger'
                      : 'border-border bg-surface-muted text-text-secondary',
                  ]"
                >
                  <ShieldCheck v-if="user.role === 'admin'" :size="13" />
                  <Users v-else :size="13" />
                  {{ user.role === 'admin' ? '管理员' : '普通用户' }}
                </span>
              </td>
              <td class="px-5 py-4">
                <span
                  :class="[
                    'inline-flex items-center gap-1.5 rounded-full border px-2.5 py-1 text-xs font-semibold',
                    user.status === 'active'
                      ? 'border-success/20 bg-success/10 text-success'
                      : 'border-border bg-surface-muted text-text-muted',
                  ]"
                >
                  <CircleCheck v-if="user.status === 'active'" :size="13" />
                  <CircleX v-else :size="13" />
                  {{ user.status === 'active' ? '正常' : '已停用' }}
                </span>
              </td>
              <td class="px-5 py-4">
                <div class="flex items-center justify-end gap-1">
                  <button
                    type="button"
                    title="编辑用户"
                    aria-label="编辑用户"
                    class="flex h-9 w-9 items-center justify-center rounded-lg text-text-muted transition-smooth hover:bg-accent-soft hover:text-accent focus-visible:ring-2 focus-visible:ring-accent focus:outline-none"
                  >
                    <Pencil :size="16" />
                  </button>
                  <button
                    type="button"
                    :title="user.status === 'active' ? '停用用户' : '启用用户'"
                    :aria-label="user.status === 'active' ? '停用用户' : '启用用户'"
                    :disabled="updatingUserId === user.id"
                    class="flex h-9 w-9 items-center justify-center rounded-lg text-text-muted transition-smooth hover:bg-surface-muted hover:text-text-primary disabled:cursor-wait disabled:opacity-50 focus-visible:ring-2 focus-visible:ring-accent focus:outline-none"
                    @click="toggleUserStatus(user)"
                  >
                    <Power :size="16" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div v-else class="flex min-h-64 flex-col items-center justify-center px-6 text-center">
        <Search :size="30" class="mb-3 text-text-muted" />
        <p class="text-sm font-medium text-text-primary">没有匹配的用户</p>
        <p class="mt-1 text-xs text-text-muted">调整搜索关键词或角色筛选。</p>
      </div>
    </section>
  </AdminPageShell>
</template>
