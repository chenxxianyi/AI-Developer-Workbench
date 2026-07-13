<template>
  <div :class="wrapperClass">
    <label v-if="label" :for="inputId" class="block text-sm font-medium text-[var(--color-text-primary)] mb-1.5">
      {{ label }}
      <span v-if="required" class="text-[var(--color-danger)] ml-0.5">*</span>
    </label>
    <div class="relative">
      <input
        :id="inputId"
        :type="type"
        :value="modelValue"
        :placeholder="placeholder"
        :disabled="disabled"
        :class="inputClasses"
        @input="onInput"
      />
      <span v-if="error" class="block text-xs text-[var(--color-danger)] mt-1">{{ error }}</span>
      <span v-else-if="hint" class="block text-xs text-[var(--color-text-muted)] mt-1">{{ hint }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, useId } from 'vue'

const props = withDefaults(defineProps<{
  modelValue?: string | number
  label?: string
  placeholder?: string
  type?: string
  disabled?: boolean
  required?: boolean
  error?: string
  hint?: string
}>(), {
  type: 'text',
})

const emit = defineEmits<{ 'update:modelValue': [value: string] }>()
const inputId = useId()

const inputClasses = computed(() =>
  `w-full px-3 py-2 text-sm rounded-[var(--radius-md)] border bg-[var(--color-surface)] transition-colors duration-200
  focus:outline-none focus:ring-2 focus:ring-[var(--color-accent)]/30 focus:border-[var(--color-accent)]
  disabled:opacity-50 disabled:cursor-not-allowed
  ${props.error ? 'border-[var(--color-danger)]' : 'border-[var(--color-border)]'}
  placeholder:text-[var(--color-text-muted)] text-[var(--color-text-primary)]`
)
const wrapperClass = 'flex flex-col gap-0'

function onInput(e: Event) {
  emit('update:modelValue', (e.target as HTMLInputElement).value)
}
</script>
