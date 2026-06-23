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
})
