<script setup lang="ts">
/**
 * ToolFormSection — label + control wrapper for tool form fields.
 * Unifies the "label row + optional required/optional marker + help text + control" pattern.
 *
 * Pass `error` to render a field-level error message tied to the control via
 * `aria-describedby` (the error node id is `${helpId ?? idFor}-error`). When
 * `error` is set, the section also exposes `data-error="true"` so focusFirstError
 * can locate the first invalid field.
 */
const props = withDefaults(
  defineProps<{
    label: string
    required?: boolean
    optional?: boolean
    /** label[for] association id. */
    idFor?: string
    /** help text shown under the control; also linked via aria-describedby when helpId is set. */
    help?: string
    /** id used for aria-describedby on the control. */
    helpId?: string
    /** render the label as a <span> instead of <label> (use for radio/checkbox groups). */
    asLabel?: boolean
    /** field-level error message; when set, the section is marked invalid. */
    error?: string
  }>(),
  {
    required: false,
    optional: false,
    asLabel: false,
    error: '',
  },
)

function describedBy(): string | undefined {
  const ids: string[] = []
  if (props.helpId) ids.push(props.helpId)
  if (props.error) ids.push(errorId())
  return ids.length ? ids.join(' ') : undefined
}

function errorId(): string {
  return `${props.helpId ?? props.idFor ?? 'field'}-error`
}
</script>

<template>
  <div class="mb-4" :data-error="error ? 'true' : undefined">
    <label
      v-if="!asLabel && idFor"
      :for="idFor"
      class="block text-sm font-medium text-text-secondary mb-2"
    >
      {{ label }}
      <span v-if="required" aria-hidden="true"> *</span>
      <span v-if="optional" class="text-text-muted font-normal"> (可选)</span>
    </label>
    <span v-else class="block text-sm font-medium text-text-secondary mb-2">
      {{ label }}
      <span v-if="required" aria-hidden="true"> *</span>
      <span v-if="optional" class="text-text-muted font-normal"> (可选)</span>
    </span>

    <slot :aria-describedby="describedBy()" :aria-invalid="error ? 'true' : undefined" />

    <p
      v-if="help"
      :id="helpId"
      class="text-xs text-text-secondary mt-1"
    >
      {{ help }}
    </p>

    <p
      v-if="error"
      :id="errorId()"
      role="alert"
      class="text-xs text-danger mt-1"
    >
      {{ error }}
    </p>
  </div>
</template>
