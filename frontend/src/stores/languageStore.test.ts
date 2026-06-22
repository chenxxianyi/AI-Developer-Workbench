import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it } from 'vitest'
import { i18n } from '@/i18n'
import { useLanguageStore } from './languageStore'

describe('language store', () => {
  beforeEach(() => {
    localStorage.clear()
    setActivePinia(createPinia())
    i18n.global.locale.value = 'zh-CN'
  })

  it('starts in Chinese when no saved preference exists', () => {
    const store = useLanguageStore()

    expect(store.locale).toBe('zh-CN')
    expect(i18n.global.locale.value).toBe('zh-CN')
  })

  it('switches locale and persists the selected language', () => {
    const store = useLanguageStore()

    store.setLocale('en-US')

    expect(store.locale).toBe('en-US')
    expect(i18n.global.locale.value).toBe('en-US')
    expect(localStorage.getItem('ai-workbench-locale')).toBe('en-US')
  })

  it('toggles between Chinese and English', () => {
    const store = useLanguageStore()

    store.toggleLocale()
    expect(store.locale).toBe('en-US')

    store.toggleLocale()
    expect(store.locale).toBe('zh-CN')
  })
})
