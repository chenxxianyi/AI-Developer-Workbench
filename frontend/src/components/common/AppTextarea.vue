<template>
  <div :class="wrapperClass">
    <label v-if="label" :for="inputId" class="block text-sm font-medium text-[var(--color-text-primary)] mb-1.5">
      {{ label }}
      <span v-if="required" class="text-[var(--color-danger)]">*</span>
    </label>
    <textarea
      :id="inputId"
      :value="modelValue"
      :placeholder="placeholder"
      :disabled="disabled"
      :rows="rows"
      :class="textareaClasses"
      @input="onInput"
    />
    <span v-if="error" class="text-xs text-[var(--color-danger)] mt-1">{{ error }}</span>
  </div>
</template>

<script setup lang="ts">
import { computed, useId } from 'vue'
const props = withDefaults(defineProps<{ modelValue?: string; label?: string; placeholder?: string; disabled?: boolean; required?: boolean; error?: string; rows?: number }>(), { rows: 4 })
const emit = defineEmits<{ 'update:modelValue': [value: string] }>()
const inputId = useId()
const textareaClasses = computed(() =>
  `w-full px-3 py-2 text-sm rounded-[var(--radius-md)] border resize-y transition-colors duration-200
  focus:outline-none focus:ring-2 focus:ring-[var(--color-accent)]/30 focus:border-[var(--color-accent)]
  disabled:opacity-50 ${props.error ? 'border-[var(--color-danger)]' : 'border-[var(--color-border)]'}
  placeholder:text-[var(--color-text-muted)] bg-[var(--color-surface)] text-[var(--color-text-primary)]`
)
const wrapperClass = 'flex flex-col'
function onInput(e: Event) { emit('update:modelValue', (e.target as HTMLTextAreaElement).value) }
</script>
