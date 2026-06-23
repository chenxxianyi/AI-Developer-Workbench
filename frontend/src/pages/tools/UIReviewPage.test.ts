// @vitest-environment jsdom

import { mount } from '@vue/test-utils'
import { nextTick } from 'vue'
import { describe, expect, it, vi } from 'vitest'
import UIReviewPage from './UIReviewPage.vue'

vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: vi.fn(),
  }),
}))

vi.mock('@/api/tools', () => ({
  runUIReview: vi.fn(),
}))

class MockFileReader {
  onload: ((event: ProgressEvent<FileReader>) => void) | null = null

  readAsDataURL(file: File) {
    this.onload?.({
      target: {
        result: `data:${file.type};base64,pasted-screenshot`,
      },
    } as ProgressEvent<FileReader>)
  }
}

describe('UIReviewPage screenshot paste upload', () => {
  it('loads a screenshot preview from an image pasted into the upload zone', async () => {
    vi.stubGlobal('FileReader', MockFileReader)

    const wrapper = mount(UIReviewPage)
    const pastedImage = new File(['image-bytes'], 'clipboard.png', { type: 'image/png' })

    const uploadZone = wrapper.get('[data-testid="screenshot-upload-zone"]')

    expect(uploadZone.text()).toContain('Ctrl+V')

    const pasteEvent = new Event('paste', { cancelable: true }) as ClipboardEvent
    Object.defineProperty(pasteEvent, 'clipboardData', {
      value: {
        items: [
          {
            type: pastedImage.type,
            getAsFile: () => pastedImage,
          },
        ],
      },
    })

    uploadZone.element.dispatchEvent(pasteEvent)
    await nextTick()

    expect(pasteEvent.defaultPrevented).toBe(true)
    expect(wrapper.get('img').attributes('src')).toBe('data:image/png;base64,pasted-screenshot')

    vi.unstubAllGlobals()
  })

  it('presents the approved focused task flow and result preview guidance', () => {
    const wrapper = mount(UIReviewPage)

    expect(wrapper.get('h1').text()).toBe('UI 质量审查')
    expect(wrapper.text()).toContain('填写输入 → 上传素材 → 开始分析 → 查看结果')
    expect(wrapper.text()).not.toContain('返回工作台')
    expect(wrapper.text()).toContain('仅上传截图，快速评估视觉质量')
    expect(wrapper.text()).toContain('仅粘贴前端代码，审查结构与样式')
    expect(wrapper.text()).toContain('截图 + 代码，获得更完整建议')
    expect(wrapper.text()).toContain('将输出哪些内容')
    expect(wrapper.text()).toContain('评分维度')
    expect(wrapper.text()).toContain('问题优先级')
    expect(wrapper.text()).toContain('改进建议')
  })
})
