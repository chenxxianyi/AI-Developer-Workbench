// @vitest-environment jsdom

import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'
import CodeInput from './CodeInput.vue'

describe('CodeInput', () => {
  it('binds value via v-model', async () => {
    const wrapper = mount(CodeInput, {
      props: { modelValue: 'initial', id: 'code' },
    })

    const textarea = wrapper.get('textarea#code')
    expect((textarea.element as HTMLTextAreaElement).value).toBe('initial')

    await textarea.setValue('changed')
    expect(wrapper.emitted('update:modelValue')![0]).toEqual(['changed'])
  })

  it('respects rows and placeholder', () => {
    const wrapper = mount(CodeInput, {
      props: { modelValue: '', rows: 12, placeholder: '粘贴代码...' },
    })

    const textarea = wrapper.get('textarea')
    expect(textarea.attributes('rows')).toBe('12')
    expect(textarea.attributes('placeholder')).toBe('粘贴代码...')
  })

  it('binds aria-describedby', () => {
    const wrapper = mount(CodeInput, {
      props: { modelValue: '', describedBy: 'code-help' },
    })

    expect(wrapper.get('textarea').attributes('aria-describedby')).toBe('code-help')
  })
})
