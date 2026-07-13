<template>
  <div class="min-h-screen flex items-center justify-center bg-[var(--color-background)]">
    <div class="w-full max-w-md bg-[var(--color-surface)] rounded-[var(--radius-panel)] border border-[var(--color-border)] p-8 shadow-[var(--shadow-md)]">
      <h1 class="text-2xl font-bold text-center text-[var(--color-text-primary)] mb-6">登录 AI Workbench</h1>
      <div class="space-y-4">
        <div><label class="block text-sm font-medium mb-1">用户名</label><input v-model="username" class="w-full px-3 py-2.5 border rounded-lg bg-white" placeholder="输入用户名" @keyup.enter="login" /></div>
        <div><label class="block text-sm font-medium mb-1">密码</label><input v-model="password" type="password" class="w-full px-3 py-2.5 border rounded-lg bg-white" placeholder="输入密码" @keyup.enter="login" /></div>
        <p v-if="error" class="text-sm text-[var(--color-danger)]">{{ error }}</p>
        <button @click="login" :disabled="loading" class="w-full py-2.5 rounded-lg bg-[var(--color-accent)] text-white font-medium disabled:opacity-50">{{ loading ? '登录中...' : '登录' }}</button>
      </div>
      <p class="text-center text-sm text-[var(--color-text-muted)] mt-6">
        没有账号？<RouterLink to="/register" class="text-[var(--color-accent)] hover:underline">注册</RouterLink>
      </p>
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import { useRouter, RouterLink } from 'vue-router'
import { useAuthStore } from '@/stores/authStore'
const router = useRouter(); const auth = useAuthStore()
const username = ref(''); const password = ref(''); const loading = ref(false); const error = ref('')
async function login() {
  if (!username.value || !password.value) { error.value = '请输入用户名和密码'; return }
  loading.value = true; error.value = ''
  try { await auth.login(username.value, password.value); router.push('/dashboard') }
  catch { error.value = '用户名或密码错误' }
  finally { loading.value = false }
}
</script>
