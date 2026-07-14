<script setup lang="ts">
import { computed, ref } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { ArrowLeft, CheckCircle2, Loader2, UserPlus, Zap } from '@lucide/vue'
import { useAuthStore } from '@/stores/authStore'

const router = useRouter()
const auth = useAuthStore()

const username = ref('')
const email = ref('')
const password = ref('')
const confirmPassword = ref('')
const loading = ref(false)
const error = ref('')

const canSubmit = computed(() =>
  username.value.trim().length >= 2
  && email.value.trim().includes('@')
  && password.value.length >= 6
  && password.value === confirmPassword.value,
)

function validateForm(): string {
  if (username.value.trim().length < 2) return '用户名至少需要 2 个字符'
  if (!email.value.trim().includes('@')) return '请输入有效邮箱'
  if (password.value.length < 6) return '密码至少需要 6 位'
  if (password.value !== confirmPassword.value) return '两次输入的密码不一致'
  return ''
}

async function register() {
  const validationError = validateForm()
  if (validationError) {
    error.value = validationError
    return
  }

  loading.value = true
  error.value = ''
  try {
    await auth.register(username.value.trim(), email.value.trim(), password.value)
    await router.push('/dashboard')
  } catch (err: any) {
    error.value = err.message || '注册失败，请稍后重试'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen bg-background px-4 py-8 text-text-primary">
    <div class="mx-auto flex min-h-[calc(100vh-4rem)] max-w-6xl items-center justify-center">
      <div class="grid w-full overflow-hidden rounded-lg border border-border bg-surface shadow-md lg:grid-cols-[0.92fr_1.08fr]">
        <section class="hidden border-r border-border bg-surface-muted/70 p-8 lg:flex lg:flex-col lg:justify-between">
          <RouterLink to="/" class="inline-flex items-center gap-3 text-text-primary transition-smooth hover:text-accent">
            <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-accent text-white">
              <Zap :size="22" />
            </div>
            <span class="text-lg font-semibold">AI Workbench</span>
          </RouterLink>

          <div>
            <div class="mb-5 flex h-12 w-12 items-center justify-center rounded-lg bg-accent-soft text-accent">
              <UserPlus :size="24" />
            </div>
            <h1 class="text-3xl font-bold text-text-primary">创建工作台账号</h1>
            <p class="mt-3 max-w-sm text-sm leading-6 text-text-secondary">
              注册后可以创建项目、运行 AI 工具，并沉淀项目报告和配置资产。
            </p>
          </div>

          <div class="grid gap-3 text-sm text-text-secondary">
            <div class="flex items-center gap-2">
              <CheckCircle2 :size="16" class="text-success" />
              <span>项目、工具和报告统一归档</span>
            </div>
            <div class="flex items-center gap-2">
              <CheckCircle2 :size="16" class="text-success" />
              <span>注册成功后自动进入工作台</span>
            </div>
          </div>
        </section>

        <section class="p-6 sm:p-8">
          <RouterLink
            to="/login"
            class="mb-6 inline-flex items-center gap-2 rounded-lg text-sm font-medium text-text-secondary transition-smooth hover:text-accent focus-visible:ring-2 focus-visible:ring-accent focus:outline-none"
          >
            <ArrowLeft :size="16" />
            返回登录
          </RouterLink>

          <div class="mb-6 lg:hidden">
            <div class="mb-4 flex h-11 w-11 items-center justify-center rounded-lg bg-accent text-white">
              <UserPlus :size="22" />
            </div>
            <h1 class="text-2xl font-bold text-text-primary">注册 AI Workbench</h1>
            <p class="mt-1 text-sm text-text-secondary">创建账号后自动进入工作台。</p>
          </div>

          <div class="hidden lg:block">
            <h2 class="text-2xl font-bold text-text-primary">注册账号</h2>
            <p class="mt-1 text-sm text-text-secondary">填写基础信息即可开始使用。</p>
          </div>

          <form class="mt-6 space-y-4" @submit.prevent="register">
            <div>
              <label for="register-username" class="mb-2 block text-sm font-medium text-text-secondary">用户名</label>
              <input
                id="register-username"
                v-model="username"
                type="text"
                autocomplete="username"
                class="w-full rounded-lg border border-border bg-surface-muted px-4 py-2.5 text-text-primary placeholder:text-text-muted focus-visible:border-accent focus-visible:ring-2 focus-visible:ring-accent focus:outline-none"
                placeholder="输入用户名"
              />
            </div>

            <div>
              <label for="register-email" class="mb-2 block text-sm font-medium text-text-secondary">邮箱</label>
              <input
                id="register-email"
                v-model="email"
                type="email"
                autocomplete="email"
                class="w-full rounded-lg border border-border bg-surface-muted px-4 py-2.5 text-text-primary placeholder:text-text-muted focus-visible:border-accent focus-visible:ring-2 focus-visible:ring-accent focus:outline-none"
                placeholder="name@example.com"
              />
            </div>

            <div class="grid gap-4 sm:grid-cols-2">
              <div>
                <label for="register-password" class="mb-2 block text-sm font-medium text-text-secondary">密码</label>
                <input
                  id="register-password"
                  v-model="password"
                  type="password"
                  autocomplete="new-password"
                  class="w-full rounded-lg border border-border bg-surface-muted px-4 py-2.5 text-text-primary placeholder:text-text-muted focus-visible:border-accent focus-visible:ring-2 focus-visible:ring-accent focus:outline-none"
                  placeholder="至少 6 位"
                />
              </div>

              <div>
                <label for="register-confirm-password" class="mb-2 block text-sm font-medium text-text-secondary">确认密码</label>
                <input
                  id="register-confirm-password"
                  v-model="confirmPassword"
                  type="password"
                  autocomplete="new-password"
                  class="w-full rounded-lg border border-border bg-surface-muted px-4 py-2.5 text-text-primary placeholder:text-text-muted focus-visible:border-accent focus-visible:ring-2 focus-visible:ring-accent focus:outline-none"
                  placeholder="再次输入密码"
                />
              </div>
            </div>

            <p v-if="error" role="alert" class="rounded-lg border border-danger/20 bg-danger/10 px-4 py-3 text-sm text-danger">
              {{ error }}
            </p>

            <button
              type="submit"
              :disabled="loading || !canSubmit"
              class="inline-flex min-h-11 w-full items-center justify-center gap-2 rounded-lg bg-accent px-4 text-sm font-semibold text-white transition-smooth hover:bg-accent/80 disabled:cursor-not-allowed disabled:bg-surface-muted disabled:text-text-muted focus-visible:ring-2 focus-visible:ring-accent focus:outline-none"
            >
              <Loader2 v-if="loading" :size="18" class="animate-spin" />
              <UserPlus v-else :size="18" />
              {{ loading ? '注册中...' : '注册并进入工作台' }}
            </button>
          </form>

          <p class="mt-6 text-center text-sm text-text-muted">
            已有账号？
            <RouterLink to="/login" class="font-medium text-accent transition-smooth hover:text-accent/80">
              去登录
            </RouterLink>
          </p>
        </section>
      </div>
    </div>
  </div>
</template>
