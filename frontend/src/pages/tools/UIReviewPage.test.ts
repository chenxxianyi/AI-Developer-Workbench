// @vitest-environment jsdom

import { mount } from '@vue/test-utils'
import { nextTick } from 'vue'
import { describe, expect, it, vi } from 'vitest'
import UIReviewPage from './UIReviewPage.vue'
import { runUIReview } from '@/api/tools'

vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: vi.fn(),
  }),
  useRoute: () => ({
    query: {},
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

  it('submits a frontend project ZIP as the code source in code mode', async () => {
    vi.mocked(runUIReview).mockResolvedValue({
      id: 'report-1',
      tool_type: 'ui_review',
      title: 'ZIP 审查',
      input_mode: 'code',
      status: 'succeeded',
      summary: 'done',
      total_score: 80,
      grade: 'B',
      input_data: {},
      report_data: {},
      generated_files: [],
      created_at: '2026-06-23T00:00:00Z',
      updated_at: '2026-06-23T00:00:00Z',
    } as any)

    const wrapper = mount(UIReviewPage)

    await wrapper.get('input[placeholder="输入分析标题..."]').setValue('ZIP 审查')
    await wrapper.get('button[data-testid="review-mode-code"]').trigger('click')
    await wrapper.get('button[data-testid="code-source-project-zip"]').trigger('click')

    expect(wrapper.text()).toContain('上传前端项目 ZIP')
    expect(wrapper.text()).toContain('系统只做静态读取，不执行代码')

    const zip = new File(['zip-bytes'], 'frontend.zip', { type: 'application/zip' })
    Object.defineProperty(wrapper.get('input[data-testid="project-zip-input"]').element, 'files', {
      value: [zip],
    })
    await wrapper.get('input[data-testid="project-zip-input"]').trigger('change')

    await wrapper.get('button[data-testid="ui-review-submit"]').trigger('click')

    const formData = vi.mocked(runUIReview).mock.calls[0][0] as FormData
    expect(formData.get('review_mode')).toBe('code')
    expect(formData.get('code_source')).toBe('project_zip')
    expect(formData.get('project_zip')).toBe(zip)
  })

  it('opens the ZIP file picker when the project ZIP dropzone is clicked', async () => {
    const wrapper = mount(UIReviewPage)

    await wrapper.get('button[data-testid="review-mode-code"]').trigger('click')
    await wrapper.get('button[data-testid="code-source-project-zip"]').trigger('click')

    const input = wrapper.get('input[data-testid="project-zip-input"]').element as HTMLInputElement
    const click = vi.spyOn(input, 'click').mockImplementation(() => {})

    await wrapper.get('[data-testid="project-zip-upload-zone"]').trigger('click')

    expect(click).toHaveBeenCalledOnce()
  })
})
