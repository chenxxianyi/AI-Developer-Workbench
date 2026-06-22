<script setup lang="ts">
import { Languages } from '@lucide/vue'
import { supportedLocales, type AppLocale } from '@/i18n'
import { useLanguageStore } from '@/stores/languageStore'

const languageStore = useLanguageStore()

function selectLocale(locale: AppLocale) {
  languageStore.setLocale(locale)
}
</script>

<template>
  <div
    class="inline-flex items-center gap-1 rounded-lg border border-border bg-surface-muted p-1"
    :aria-label="$t('landing.nav.language')"
  >
    <Languages :size="16" class="mx-1 text-text-muted" aria-hidden="true" />
    <button
      v-for="localeOption in supportedLocales"
      :key="localeOption.code"
      type="button"
      class="min-w-9 rounded-md px-2 py-1 text-sm font-semibold transition-smooth cursor-pointer focus-visible:outline-accent"
      :class="languageStore.locale === localeOption.code
        ? 'bg-surface text-accent shadow-sm'
        : 'text-text-secondary hover:text-text-primary'"
      :aria-pressed="languageStore.locale === localeOption.code"
      @click="selectLocale(localeOption.code)"
    >
      {{ localeOption.shortLabel }}
    </button>
  </div>
</template>
