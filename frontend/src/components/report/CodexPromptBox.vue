<script setup lang="ts">
/**
 * CodexPromptBox — displays and copies the AI fix prompt.
 */
import { ref } from 'vue'
import { Copy, Check } from '@lucide/vue'
import { copyToClipboard } from '@/utils/clipboard'

const props = defineProps<{
  prompt: string
}>()

const copied = ref(false)

async function handleCopy() {
  try {
    await copyToClipboard(props.prompt)
    copied.value = true
    setTimeout(() => { copied.value = false }, 2000)
  } catch {
    // Error handled silently — feedback via button state.
  }
}
</script>

<template>
  <div v-if="prompt" class="rounded-lg border border-accent/30 bg-accent-soft/50 p-5">
    <div class="flex items-center justify-between mb-3">
      <h3 class="text-lg font-semibold text-text-primary">AI 修复 Prompt</h3>
      <button
        class="inline-flex items-center gap-1.5 rounded-md border px-3 py-1.5 text-sm font-medium transition-smooth"
        :class="
          copied
            ? 'border-emerald-300 bg-emerald-50 text-emerald-700 dark:bg-emerald-900/20 dark:text-emerald-300'
            : 'border-accent/30 bg-surface text-accent hover:bg-accent-soft'
        "
        @click="handleCopy"
        aria-label="复制 Prompt 到剪贴板"
      >
        <Check v-if="copied" :size="15" />
        <Copy v-else :size="15" />
        {{ copied ? '已复制' : '复制 Prompt' }}
      </button>
    </div>
    <pre class="whitespace-pre-wrap text-sm text-text-secondary bg-background/70 rounded-md p-4 border border-border">{{ prompt }}</pre>
  </div>
</template>
