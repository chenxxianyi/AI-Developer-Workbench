// @vitest-environment jsdom

import { mount } from '@vue/test-utils'
import { describe, expect, it, vi } from 'vitest'
import ProjectDoctorPage from './ProjectDoctorPage.vue'

vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: vi.fn(),
  }),
}))

vi.mock('@/api/tools', () => ({
  runProjectDoctor: vi.fn(),
}))

describe('ProjectDoctorPage accessibility improvements', () => {
  it('binds labels to fields and exposes semantic upload/depth controls', () => {
    const wrapper = mount(ProjectDoctorPage)

    expect(wrapper.get('label[for="project-doctor-title"]').text()).toContain('标题')
    expect(wrapper.get('#project-doctor-title').attributes('required')).toBeDefined()
    expect(wrapper.get('label[for="project-doctor-name"]').text()).toContain('项目名称')
    expect(wrapper.get('label[for="project-doctor-zip"]').text()).toContain('项目 ZIP 文件')
    expect(wrapper.get('label[for="project-doctor-tech-stack"]').text()).toContain('技术栈')
    expect(wrapper.get('label[for="project-doctor-description"]').text()).toContain('项目描述')

    const uploadZone = wrapper.get('[data-testid="project-zip-upload-zone"]')
    expect(uploadZone.attributes('role')).toBe('button')
    expect(uploadZone.attributes('tabindex')).toBe('0')
    expect(uploadZone.attributes('aria-describedby')).toBe('project-doctor-zip-help')
    expect(wrapper.text()).toContain('按 Enter 或空格选择文件')

    const depthGroup = wrapper.get('[role="radiogroup"]')
    expect(depthGroup.attributes('aria-label')).toBe('分析深度')
    expect(wrapper.findAll('[role="radio"]')).toHaveLength(3)
    expect(wrapper.find('[role="radio"][aria-checked="true"]').text()).toContain('标准')
  })

  it('renders issue severity as localized text instead of raw color-only status', async () => {
    const wrapper = mount(ProjectDoctorPage)

    ;(wrapper.vm as any).result = {
      id: 'report-1',
      title: 'Project Doctor',
      type: 'project_doctor',
      status: 'completed',
      summary: '存在高优先级问题',
      total_score: 72,
      grade: 'C',
      report_data: {
        scores: [],
        issues: [
          {
            title: '依赖风险',
            severity: 'high',
            problem: '存在过期依赖',
            suggestion: '升级依赖并补充锁文件审查',
          },
        ],
        recommendations: [],
      },
      generated_files: [],
      created_at: '2026-06-23T00:00:00Z',
      updated_at: '2026-06-23T00:00:00Z',
    }
    await wrapper.vm.$nextTick()

    expect(wrapper.text()).toContain('高')
    expect(wrapper.text()).not.toContain('high')
  })
})
