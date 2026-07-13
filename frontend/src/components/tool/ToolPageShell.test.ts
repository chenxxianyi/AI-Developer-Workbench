// @vitest-environment jsdom

import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'
import { Eye } from '@lucide/vue'
import ToolPageShell from './ToolPageShell.vue'

describe('ToolPageShell', () => {
  it('renders the title, description and step text', () => {
    const wrapper = mount(ToolPageShell, {
      props: { icon: Eye, title: 'UI 质量审查', description: '描述', stepText: 'A → B' },
    })

    expect(wrapper.get('h1').text()).toBe('UI 质量审查')
    expect(wrapper.text()).toContain('描述')
    expect(wrapper.text()).toContain('A → B')
  })

  it('renders the error message when provided', () => {
    const wrapper = mount(ToolPageShell, {
      props: { icon: Eye, title: 't', description: 'd', error: '出错啦' },
    })

    expect(wrapper.text()).toContain('出错啦')
  })

  it('disables the submit button while loading or when canSubmit is false', async () => {
    const wrapper = mount(ToolPageShell, {
      props: { icon: Eye, title: 't', description: 'd', canSubmit: false },
    })

    const submit = wrapper.findAll('button[type="button"]').find((b) => b.text().includes('开始分析'))!
    expect((submit.element as HTMLButtonElement).disabled).toBe(true)

    await wrapper.setProps({ canSubmit: true })
    expect((submit.element as HTMLButtonElement).disabled).toBe(false)

    await wrapper.setProps({ loading: true })
    expect((submit.element as HTMLButtonElement).disabled).toBe(true)
  })

  it('emits submit when the submit button is clicked and not disabled', async () => {
    const wrapper = mount(ToolPageShell, {
      props: { icon: Eye, title: 't', description: 'd', canSubmit: true },
    })

    const submit = wrapper.findAll('button[type="button"]').find((b) => b.text().includes('开始分析'))!
    await submit.trigger('click')

    expect(wrapper.emitted('submit')).toHaveLength(1)
  })

  it('does not emit submit when disabled', async () => {
    const wrapper = mount(ToolPageShell, {
      props: { icon: Eye, title: 't', description: 'd', canSubmit: false },
    })

    const submit = wrapper.findAll('button[type="button"]').find((b) => b.text().includes('开始分析'))!
    await submit.trigger('click')

    expect(wrapper.emitted('submit')).toBeUndefined()
  })

  it('shows the loading spinner and hint while loading', () => {
    const wrapper = mount(ToolPageShell, {
      props: { icon: Eye, title: 't', description: 'd', loading: true, loadingHint: 'AI 正在分析...' },
      slots: { result: '<div data-testid="r"></div>' },
    })

    expect(wrapper.text()).toContain('AI 正在分析...')
    // result slot not rendered while loading
    expect(wrapper.find('[data-testid="r"]').exists()).toBe(false)
  })

  it('renders the empty slot when there is no result and no loading', () => {
    const wrapper = mount(ToolPageShell, {
      props: { icon: Eye, title: 't', description: 'd' },
      slots: { empty: '<div data-testid="empty">empty</div>' },
    })

    expect(wrapper.find('[data-testid="empty"]').exists()).toBe(true)
  })

  it('renders the result slot when provided and not loading', () => {
    const wrapper = mount(ToolPageShell, {
      props: { icon: Eye, title: 't', description: 'd' },
      slots: { result: '<div data-testid="r">result</div>' },
    })

    expect(wrapper.find('[data-testid="r"]').exists()).toBe(true)
    expect(wrapper.find('[data-testid="empty"]').exists()).toBe(false)
  })

  it('emits back when the back button is clicked', async () => {
    const wrapper = mount(ToolPageShell, {
      props: { icon: Eye, title: 't', description: 'd' },
    })

    const back = wrapper.findAll('button[type="button"]').find((b) => b.text().includes('返回工作台'))!
    await back.trigger('click')

    expect(wrapper.emitted('back')).toHaveLength(1)
  })

  it('hides the back button when backLabel is empty', () => {
    const wrapper = mount(ToolPageShell, {
      props: { icon: Eye, title: 't', description: 'd', backLabel: '' },
    })

    expect(wrapper.text()).not.toContain('返回工作台')
  })
})
