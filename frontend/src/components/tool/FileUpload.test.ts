// @vitest-environment jsdom

import { mount, flushPromises } from '@vue/test-utils'
import { describe, expect, it, vi } from 'vitest'
import FileUpload from './FileUpload.vue'

function makeFile(name: string, type = 'application/zip'): File {
  return new File(['bytes'], name, { type })
}

describe('FileUpload', () => {
  it('opens the native picker on click', async () => {
    const wrapper = mount(FileUpload, {
      props: { modelValue: null, testid: 'project-zip', inputId: 'project-zip' },
    })

    const input = wrapper.get('input[data-testid="project-zip-input"]').element as HTMLInputElement
    const click = vi.spyOn(input, 'click').mockImplementation(() => {})

    await wrapper.get('[data-testid="project-zip-upload-zone"]').trigger('click')

    expect(click).toHaveBeenCalledOnce()
  })

  it('opens the picker on Enter and Space', async () => {
    const wrapper = mount(FileUpload, {
      props: { modelValue: null, testid: 'pz' },
    })

    const input = wrapper.get('input[data-testid="pz-input"]').element as HTMLInputElement
    const click = vi.spyOn(input, 'click').mockImplementation(() => {})

    const zone = wrapper.get('[data-testid="pz-upload-zone"]')
    await zone.trigger('keydown', { key: 'Enter' })
    await zone.trigger('keydown', { key: ' ' })

    expect(click).toHaveBeenCalledTimes(2)
  })

  it('emits update:modelValue when a file is selected via the input', async () => {
    const wrapper = mount(FileUpload, {
      props: { modelValue: null, testid: 'pz' },
    })

    const input = wrapper.get('input[data-testid="pz-input"]')
    const file = makeFile('frontend.zip')
    Object.defineProperty(input.element, 'files', { value: [file] })

    await input.trigger('change')

    expect(wrapper.emitted('update:modelValue')![0]).toEqual([file])
  })

  it('accepts a dropped file and emits update:modelValue', async () => {
    const wrapper = mount(FileUpload, {
      props: { modelValue: null, testid: 'pz' },
    })

    const file = makeFile('dropped.zip')
    // jsdom lacks DataTransfer; use a plain object carrying files.
    const dt = { files: [file] }

    await wrapper.get('[data-testid="pz-upload-zone"]').trigger('drop', { dataTransfer: dt })

    expect(wrapper.emitted('update:modelValue')![0]).toEqual([file])
  })

  it('clears the file when the remove button is clicked', async () => {
    const file = makeFile('selected.zip')
    const wrapper = mount(FileUpload, {
      props: { modelValue: file },
    })

    expect(wrapper.text()).toContain('selected.zip')

    await wrapper.get('button[aria-label="移除已上传文件"]').trigger('click')

    expect(wrapper.emitted('update:modelValue')![0]).toEqual([null])
  })

  it('supports Ctrl+V paste in pasteable mode and emits update:modelValue + paste', async () => {
    const wrapper = mount(FileUpload, {
      props: { modelValue: null, pasteable: true, preview: true, testid: 'shot' },
    })

    const image = makeFile('clipboard.png', 'image/png')
    class MockFileReader {
      onload: ((event: ProgressEvent<FileReader>) => void) | null = null
      readAsDataURL(file: File) {
        this.onload?.({ target: { result: `data:${file.type};base64,pasted-screenshot` } } as ProgressEvent<FileReader>)
      }
    }
    vi.stubGlobal('FileReader', MockFileReader)

    const pasteEvent = new Event('paste', { cancelable: true }) as ClipboardEvent
    Object.defineProperty(pasteEvent, 'clipboardData', {
      value: {
        items: [{ type: image.type, getAsFile: () => image }],
      },
    })

    wrapper.get('[data-testid="shot-upload-zone"]').element.dispatchEvent(pasteEvent)
    await flushPromises()

    expect(pasteEvent.defaultPrevented).toBe(true)
    expect(wrapper.emitted('update:modelValue')![0]).toEqual([image])
    expect(wrapper.emitted('paste')![0]).toEqual([image])
    vi.unstubAllGlobals()
  })

  it('ignores paste when not in pasteable mode', async () => {
    const wrapper = mount(FileUpload, {
      props: { modelValue: null, pasteable: false, testid: 'pz' },
    })

    const image = makeFile('clipboard.png', 'image/png')
    const pasteEvent = new Event('paste', { cancelable: true }) as ClipboardEvent
    Object.defineProperty(pasteEvent, 'clipboardData', {
      value: {
        items: [{ type: image.type, getAsFile: () => image }],
      },
    })

    wrapper.get('[data-testid="pz-upload-zone"]').element.dispatchEvent(pasteEvent)
    await flushPromises()

    expect(pasteEvent.defaultPrevented).toBe(false)
    expect(wrapper.emitted('update:modelValue')).toBeUndefined()
  })

  it('renders the preview image when preview is true and a file is set', async () => {
    const wrapper = mount(FileUpload, {
      props: { modelValue: null, preview: true },
    })

    const image = makeFile('shot.png', 'image/png')
    class MockFileReader {
      onload: ((event: ProgressEvent<FileReader>) => void) | null = null
      readAsDataURL(file: File) {
        this.onload?.({ target: { result: `data:${file.type};base64,pasted-screenshot` } } as ProgressEvent<FileReader>)
      }
    }
    vi.stubGlobal('FileReader', MockFileReader)

    await wrapper.setProps({ modelValue: image })

    expect(wrapper.find('img').exists()).toBe(true)
    vi.unstubAllGlobals()
  })
})
