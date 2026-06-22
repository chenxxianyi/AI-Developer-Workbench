import { defineStore } from 'pinia'
import { ref } from 'vue'
import { DEFAULT_LOCALE, LOCALE_STORAGE_KEY, i18n, supportedLocales, type AppLocale } from '@/i18n'

function isSupportedLocale(locale: string | null): locale is AppLocale {
  return supportedLocales.some((supportedLocale) => supportedLocale.code === locale)
}

function getSavedLocale(): AppLocale {
  if (typeof localStorage === 'undefined') {
    return DEFAULT_LOCALE
  }

  const savedLocale = localStorage.getItem(LOCALE_STORAGE_KEY)
  return isSupportedLocale(savedLocale) ? savedLocale : DEFAULT_LOCALE
}

export const useLanguageStore = defineStore('language', () => {
  const locale = ref<AppLocale>(getSavedLocale())

  function applyLocale(nextLocale: AppLocale) {
    locale.value = nextLocale
    i18n.global.locale.value = nextLocale

    if (typeof localStorage !== 'undefined') {
      localStorage.setItem(LOCALE_STORAGE_KEY, nextLocale)
    }
  }

  function setLocale(nextLocale: AppLocale) {
    applyLocale(nextLocale)
  }

  function toggleLocale() {
    applyLocale(locale.value === 'zh-CN' ? 'en-US' : 'zh-CN')
  }

  applyLocale(locale.value)

  return {
    locale,
    setLocale,
    toggleLocale,
  }
})
