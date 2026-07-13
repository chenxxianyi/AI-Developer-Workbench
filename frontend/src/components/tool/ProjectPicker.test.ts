// @vitest-environment jsdom

import { flushPromises, mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import ProjectPicker from './ProjectPicker.vue'
import { listProjects } from '@/api/projects'

vi.mock('@/api/projects', () => ({
  listProjects: vi.fn(),
}))

const mockListProjects = vi.mocked(listProjects)

describe('ProjectPicker', () => {
  beforeEach(() => {
    mockListProjects.mockResolvedValue({
      items: [
        { id: 'alpha', name: 'Alpha', description: '', repo_url: '', report_count: 0, average_score: null, created_at: '', updated_at: '' },
        { id: 'beta', name: 'Beta', description: '', repo_url: '', report_count: 1, average_score: 80, created_at: '', updated_at: '' },
      ],
      total: 2,
      page: 1,
      page_size: 100,
    })
  })

  it('loads projects and emits the selected optional project id', async () => {
    const wrapper = mount(ProjectPicker, {
      props: { modelValue: '', inputId: 'tool-project', helpId: 'tool-project-help' },
    })
    await flushPromises()

    const select = wrapper.get('select#tool-project')
    expect(select.findAll('option')).toHaveLength(3)
    await select.setValue('beta')

    expect(wrapper.emitted('update:modelValue')).toEqual([['beta']])
    expect(select.attributes('aria-describedby')).toBe('tool-project-help')
  })

  it('reloads from the API with a search query', async () => {
    const wrapper = mount(ProjectPicker, {
      props: { modelValue: '', inputId: 'tool-project' },
    })
    await flushPromises()
    await wrapper.get('input[type="search"]').setValue('Alpha')
    await flushPromises()

    expect(mockListProjects).toHaveBeenLastCalledWith({
      search: 'Alpha',
      page: 1,
      page_size: 100,
    })
  })
})
