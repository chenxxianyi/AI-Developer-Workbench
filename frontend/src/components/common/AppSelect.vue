<template>
  <div>
    <label v-if="label" :for="inputId" class="block text-sm font-medium text-[var(--color-text-primary)] mb-1.5">
      {{ label }}<span v-if="required" class="text-[var(--color-danger)] ml-0.5">*</span>
    </label>
    <select
      :id="inputId"
      :value="modelValue"
      :disabled="disabled"
      :class="selectClasses"
      @change="emit('update:modelValue', ($event.target as HTMLSelectElement).value)"
    >
      <option v-if="placeholder" value="" disabled>{{ placeholder }}</option>
      <option v-for="opt in options" :key="opt.value" :value="opt.value">{{ opt.label }}</option>
    </select>
    <span v-if="error" class="block text-xs text-[var(--color-danger)] mt-1">{{ error }}</span>
  </div>
</template>

<script setup lang="ts">
import { computed, useId } from 'vue'
const props = withDefaults(defineProps<{
  modelValue?: string; options: { label: string; value: string }[]; label?: string; placeholder?: string; disabled?: boolean; required?: boolean; error?: string
}>(), { options: () => [] })
const emit = defineEmits<{ 'update:modelValue': [value: string] }>()
const inputId = useId()
const selectClasses = computed(() =>
  `w-full px-3 py-2 text-sm rounded-[var(--radius-md)] border bg-[var(--color-surface)] transition-colors duration-200 appearance-none
  focus:outline-none focus:ring-2 focus:ring-[var(--color-accent)]/30 focus:border-[var(--color-accent)]
  disabled:opacity-50 ${props.error ? 'border-[var(--color-danger)]' : 'border-[var(--color-border)]'} text-[var(--color-text-primary)]`
)
</script>
