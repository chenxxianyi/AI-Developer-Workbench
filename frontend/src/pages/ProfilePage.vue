<template>
  <div class="max-w-2xl mx-auto py-8">
    <h1 class="text-2xl font-bold text-[var(--color-text-primary)] mb-6">个人资料</h1>
    <div class="bg-[var(--color-surface)] rounded-[var(--radius-panel)] border border-[var(--color-border)] p-6 space-y-4">
      <div class="grid grid-cols-2 gap-4">
        <div><label class="block text-sm font-medium text-[var(--color-text-primary)] mb-1">用户名</label><input class="w-full px-3 py-2 border rounded-lg" aria-label="用户名" disabled :value="user?.username" /></div>
        <div><label class="block text-sm font-medium text-[var(--color-text-primary)] mb-1">邮箱</label><input class="w-full px-3 py-2 border rounded-lg" aria-label="邮箱" v-model="email" /></div>
      </div>
      <div class="flex justify-end gap-3 pt-4">
        <button class="px-4 py-2 text-sm rounded-lg border" @click="$router.push('/dashboard')">取消</button>
        <button class="px-4 py-2 text-sm rounded-lg bg-[var(--color-accent)] text-white" @click="save">保存</button>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/authStore'
const auth = useAuthStore()
const user = ref(auth.user)
const email = ref(auth.user?.email || '')
onMounted(async () => { if (auth.isLoggedIn) await auth.fetchProfile(); user.value = auth.user; email.value = auth.user?.email || '' })
async function save() { await auth.updateProfile({ email: email.value }) }
</script>
