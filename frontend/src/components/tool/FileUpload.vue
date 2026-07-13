<script setup lang="ts">
/**
 * FileUpload — accessible drag/drop/click/keyboard file input.
 *
 * Renders a dropzone with role="button" + tabindex="0", opens the native
 * picker on click / Enter / Space, accepts drag-drop, and shows either an
 * image preview (preview=true) or a file name chip once a file is selected.
 *
 * `testid` is transparently propagated:
 *   - dropzone -> data-testid="${testid}-upload-zone"
 *   - input    -> data-testid="${testid}-input"
 *
 * `pasteable` enables Ctrl+V image paste (matches the UI Review screenshot flow).
 */
import { ref, computed, watch } from 'vue'
import { Upload, FileText, AlertCircle, X } from '@lucide/vue'

const props = withDefaults(
  defineProps<{
    modelValue: File | null
    accept?: string
    /** dropzone main text (e.g. "上传前端项目 ZIP"). */
    emptyText?: string
    /** dropzone secondary text (e.g. "支持 .zip，最大 20MB"). */
    hint?: string
    /** help paragraph id; the control gains aria-describedby=helpId. */
    helpId?: string
    /** id for the <input>, used by external <label for>. */
    inputId?: string
    /** testid stem; see file header. */
    testid?: string
    /** accent color key for the dropzone border/icon. */
    accent?: 'accent' | 'success' | 'orange' | 'teal' | 'purple'
    /** image mode: show preview + support Ctrl+V paste. */
    pasteable?: boolean
    /** render an <img> preview instead of a file-name chip. */
    preview?: boolean
    /** aria-label for the remove button. */
    removeLabel?: string
  }>(),
  {
    accept: '',
    emptyText: '点击、拖拽或按 Enter / 空格上传文件',
    hint: '',
    testid: '',
    accent: 'accent',
    pasteable: false,
    preview: false,
    removeLabel: '移除已上传文件',
  },
)

const emit = defineEmits<{
  'update:modelValue': [file: File | null]
  paste: [file: File]
}>()

const inputRef = ref<HTMLInputElement | null>(null)
const previewUrl = ref<string | null>(null)

const fileName = computed(() => props.modelValue?.name ?? '')

// Generate a preview URL whenever modelValue changes (covers the case where the
// parent sets the file directly rather than through setFile, e.g. in tests).
// Uses FileReader.readAsDataURL so the result is a portable data URL (matches the
// screenshot preview behavior expected by callers and existing tests).
watch(
  () => props.modelValue,
  (file) => {
    previewUrl.value = null
    if (!file || !props.preview) return
    const reader = new FileReader()
    reader.onload = (e) => {
      previewUrl.value = (e.target?.result as string) ?? null
    }
    reader.readAsDataURL(file)
  },
  { immediate: true },
)

function setFile(file: File | null) {
  emit('update:modelValue', file)
}

function openPicker() {
  inputRef.value?.click()
}

function onSelect(event: Event) {
  const target = event.target as HTMLInputElement
  if (target.files?.length) {
    setFile(target.files[0])
  }
}

function onKeydown(event: KeyboardEvent) {
  if (event.key !== 'Enter' && event.key !== ' ') return
  event.preventDefault()
  openPicker()
}

function onDrop(event: DragEvent) {
  if (!event.dataTransfer?.files?.length) return
  event.preventDefault()
  setFile(event.dataTransfer.files[0])
}

function onDragover(event: DragEvent) {
  // allow drop
  event.preventDefault()
}

function onPaste(event: ClipboardEvent) {
  if (!props.pasteable) return
  const items = event.clipboardData?.items
  if (!items?.length) return
  for (let i = 0; i < items.length; i += 1) {
    const item = items[i]
    if (!item.type.startsWith('image/')) continue
    const file = item.getAsFile()
    if (!file) return
    event.preventDefault()
    setFile(file)
    emit('paste', file)
    return
  }
}

function clear() {
  setFile(null)
}

const accentBorderHover: Record<string, string> = {
  accent: 'hover:border-accent focus-visible:ring-accent focus:border-accent',
  success: 'hover:border-success focus-visible:ring-success focus:border-success',
  orange: 'hover:border-orange-500 focus-visible:ring-orange-500 focus:border-orange-500',
  teal: 'hover:border-teal-500 focus-visible:ring-teal-500 focus:border-teal-500',
  purple: 'hover:border-purple-500 focus-visible:ring-purple-500 focus:border-purple-500',
}
const accentIcon: Record<string, string> = {
  accent: 'text-accent',
  success: 'text-success',
  orange: 'text-orange-500',
  teal: 'text-teal-500',
  purple: 'text-purple-500',
}
</script>

<template>
  <div>
    <!-- Dropzone / preview -->
    <div v-if="!modelValue">
      <div
        :data-testid="testid ? `${testid}-upload-zone` : undefined"
        role="button"
        tabindex="0"
        :aria-describedby="helpId"
        :class="[
          'border-2 border-dashed border-border/80 rounded-lg p-8 text-center focus-visible:ring-2 focus-visible:border-2 focus:outline-none transition-smooth cursor-pointer bg-surface-muted/40',
          accentBorderHover[accent],
        ]"
        @click="openPicker"
        @keydown="onKeydown"
        @drop="onDrop"
        @dragover="onDragover"
        @paste="onPaste"
      >
        <Upload :size="32" :class="['mx-auto mb-2', accentIcon[accent]]" />
        <p class="text-text-primary font-medium">{{ emptyText }}</p>
        <p v-if="hint" class="text-text-secondary text-sm mt-1">{{ hint }}</p>
        <p v-if="pasteable" class="text-text-muted text-xs mt-1">支持 Ctrl+V 粘贴</p>
      </div>
      <input
        ref="inputRef"
        :id="inputId"
        :data-testid="testid ? `${testid}-input` : undefined"
        type="file"
        :accept="accept"
        class="hidden"
        @change="onSelect"
      />
    </div>

    <!-- Image preview -->
    <div v-else-if="preview && previewUrl" class="relative">
      <img
        :src="previewUrl!"
        class="w-full rounded-lg border border-border"
        style="max-height: 300px; object-fit: contain;"
      />
      <button
        type="button"
        :aria-label="removeLabel"
        class="absolute top-2 right-2 p-1 bg-danger text-white rounded hover:bg-danger/80 transition-smooth focus-visible:ring-2 focus-visible:ring-danger focus:outline-none"
        @click="clear"
      >
        <X :size="16" />
      </button>
    </div>

    <!-- File chip -->
    <div v-else class="flex items-center justify-between gap-3 p-3 bg-surface-muted rounded-lg">
      <div class="flex items-center gap-2 min-w-0">
        <FileText :size="18" :class="accentIcon[accent]" />
        <span class="text-text-primary truncate">{{ fileName }}</span>
      </div>
      <button
        type="button"
        :aria-label="removeLabel"
        class="p-1 text-danger hover:text-danger/80 rounded focus-visible:ring-2 focus-visible:ring-danger focus:outline-none"
        @click="clear"
      >
        <AlertCircle :size="16" />
      </button>
    </div>
  </div>
</template>
