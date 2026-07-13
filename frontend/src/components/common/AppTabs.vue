<template>
  <div>
    <div class="flex border-b border-[var(--color-border)]" role="tablist">
      <button
        v-for="tab in tabs"
        :key="tab.value"
        role="tab"
        :aria-selected="modelValue === tab.value"
        :class="tabClass(tab.value)"
        @click="emit('update:modelValue', tab.value)"
      >
        {{ tab.label }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
export interface Tab { label: string; value: string }
const props = defineProps<{ tabs: Tab[]; modelValue: string }>()
const emit = defineEmits<{ 'update:modelValue': [value: string] }>()
function tabClass(val: string) {
  const active = val === props.modelValue
  return `px-4 py-2.5 text-sm font-medium transition-colors duration-200 border-b-2 -mb-[1px]
    ${active ? 'border-[var(--color-accent)] text-[var(--color-accent)]' : 'border-transparent text-[var(--color-text-muted)] hover:text-[var(--color-text-secondary)]'}`
}
</script>
