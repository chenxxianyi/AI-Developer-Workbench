<script setup lang="ts">
/**
 * CodeInput — monospace textarea for code/schema input.
 */
const props = withDefaults(
  defineProps<{
    modelValue: string
    placeholder?: string
    rows?: number
    id?: string
    /** aria-describedby association. */
    describedBy?: string
  }>(),
  {
    placeholder: '',
    rows: 8,
  },
)

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

function onInput(event: Event) {
  emit('update:modelValue', (event.target as HTMLTextAreaElement).value)
}
</script>

<template>
  <textarea
    :id="id"
    :value="modelValue"
    :rows="rows"
    :placeholder="placeholder"
    :aria-describedby="describedBy"
    class="w-full px-4 py-2 bg-surface-muted border border-border/80 rounded-lg focus-visible:ring-2 focus-visible:ring-accent focus-visible:border-accent focus:outline-none font-mono text-sm text-text-primary placeholder:text-text-muted"
    @input="onInput"
  ></textarea>
</template>
