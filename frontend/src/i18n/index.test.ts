import { describe, expect, it } from 'vitest'
import { DEFAULT_LOCALE, messages, supportedLocales } from './index'

describe('i18n messages', () => {
  it('defaults to Chinese and provides Chinese and English landing copy', () => {
    expect(DEFAULT_LOCALE).toBe('zh-CN')
    expect(supportedLocales.map((locale) => locale.code)).toEqual(['zh-CN', 'en-US'])

    expect(messages['zh-CN'].landing.hero.title).toBe('更高质量地交付 AI 生成项目')
    expect(messages['en-US'].landing.hero.title).toBe('Build better AI-generated projects')
    expect(messages['zh-CN'].landing.nav.home).toBe('首页')
    expect(messages['en-US'].landing.nav.home).toBe('Home')
    expect(messages['zh-CN'].landing.nav.start).toBe('开始使用')
    expect(messages['en-US'].landing.nav.start).toBe('Get Started')
  })
})
