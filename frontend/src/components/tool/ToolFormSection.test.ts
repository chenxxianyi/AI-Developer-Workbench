// @vitest-environment jsdom

import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'
import ToolFormSection from './ToolFormSection.vue'

describe('ToolFormSection', () => {
  it('renders a label associated with idFor', () => {
    const wrapper = mount(ToolFormSection, {
      props: { label: '标题', required: true, idFor: 'field-title' },
      slots: { default: '<input id="field-title" />' },
    })

    const label = wrapper.get('label[for="field-title"]')
    expect(label.text()).toContain('标题')
    expect(label.text()).toContain('*')
  })

  it('renders the optional marker when optional is true', () => {
    const wrapper = mount(ToolFormSection, {
      props: { label: '页面类型', optional: true, idFor: 'p' },
      slots: { default: '<input id="p" />' },
    })

    expect(wrapper.get('label').text()).toContain('(可选)')
  })

  it('renders a span label (not for-bound) when asLabel is true', () => {
    const wrapper = mount(ToolFormSection, {
      props: { label: '分析深度', asLabel: true, required: true },
    })

    expect(wrapper.find('label').exists()).toBe(false)
    expect(wrapper.get('span').text()).toContain('分析深度')
  })

  it('renders help text with the given helpId', () => {
    const wrapper = mount(ToolFormSection, {
      props: { label: 'ZIP', help: '按 Enter 或空格选择文件', helpId: 'zip-help', idFor: 'zip' },
      slots: { default: '<input id="zip" />' },
    })

    const help = wrapper.get('p#zip-help')
    expect(help.text()).toBe('按 Enter 或空格选择文件')
  })
})
