import { test, expect, type Page, type Route } from '@playwright/test'

const projectId = 'project-e2e'

const project = {
  id: projectId,
  name: 'E2E 企业官网',
  project_type: 'landing_page',
  description: '自动化测试创建的企业官网项目',
  repo_url: '',
  frontend_stack: 'Vue 3 + TypeScript + Tailwind CSS',
  backend_stack: '',
  status: 'active',
  created_at: '2026-07-15T00:00:00Z',
  updated_at: '2026-07-15T00:00:00Z',
}

const blueprint = {
  id: 'blueprint-e2e',
  project_id: projectId,
  status: 'draft',
  version: 1,
  content: JSON.stringify({
    schema_version: 2,
    app_type: 'landing_page',
    product_positioning: '面向企业客户的品牌官网与线索收集站点',
    tech_stack: 'Vue 3 + TypeScript + Tailwind CSS',
    pages: [
      { name: '首页', route: '/' },
      { name: '解决方案', route: '/solutions' },
      { name: '联系我们', route: '/contact' },
    ],
    features: [
      { id: 'F-001', name: '产品展示', priority: 'must', description: '展示产品能力', acceptance_criteria: ['首页可查看产品能力'] },
      { id: 'F-002', name: '线索表单', priority: 'must', description: '收集合作需求', acceptance_criteria: ['用户可以提交线索表单'] },
    ],
    acceptance_criteria: ['首页正常显示'],
    test_plan: ['验证首页加载'],
  }),
}

test.describe('项目生成工作流', () => {
  test('从首页创建项目、填写需求、生成蓝图、生成代码、预览并查看文件', async ({ page }) => {
    await mockAuth(page)
    await mockTaskStream(page)
    await mockWorkflowApi(page)

    await page.goto('/')
    await expect(page.getByRole('heading', { name: /从想法到/ })).toBeVisible()

    await page.getByRole('link', { name: /进入生成工作室/ }).click()
    await expect(page).toHaveURL(/\/projects\/new$/)
    await expect(page.getByRole('heading', { name: '创建新项目' })).toBeVisible()

    await page.getByRole('button', { name: /继续/ }).click()
    await page.getByLabel('项目名称').fill(project.name)
    await page.getByLabel('项目描述').fill(project.description)
    await page.getByRole('button', { name: /创建项目/ }).click()

    await expect(page.getByRole('heading', { name: '项目创建成功' })).toBeVisible()
    await page.getByRole('link', { name: /填写项目需求/ }).click()

    await expect(page).toHaveURL(new RegExp(`/projects/${projectId}/requirements$`))
    await expect(page.getByRole('heading', { name: '告诉 AI 你想做什么' })).toBeVisible()
    await page.getByLabel('你的产品想法').fill('上线一个可以展示产品能力并收集销售线索的企业官网。')
    await page.getByRole('button', { name: /让 AI 帮我整理/ }).click()
    await expect(page.getByText('AI 已经整理好主要内容')).toBeVisible()
    await page.getByRole('button', { name: /^继续/ }).click()
    await expect(page.getByRole('heading', { name: '谁会使用？通常怎么使用？' })).toBeVisible()
    await page.getByRole('button', { name: /^继续/ }).click()
    await expect(page.getByRole('heading', { name: '确认这次需要完成的功能' })).toBeVisible()
    await page.getByRole('button', { name: /^继续/ }).click()
    await expect(page.getByRole('heading', { name: '选择喜欢的外观和使用设备' })).toBeVisible()
    await page.getByRole('button', { name: /现在生成完成效果/ }).click()
    await page.getByRole('button', { name: /整理需求并查看方案/ }).click()

    await expect(page).toHaveURL(new RegExp(`/projects/${projectId}/blueprint$`))
    await expect(page.getByRole('heading', { name: '产品与技术蓝图' })).toBeVisible()
    await page.getByRole('button', { name: /生成完整蓝图/ }).click()
    await expect(page.getByText('面向企业客户的品牌官网与线索收集站点')).toBeVisible()
    await expect(page.getByText('Vue 3 + TypeScript + Tailwind CSS')).toBeVisible()
    await page.getByRole('button', { name: /确认并进入生成/ }).click()

    await expect(page).toHaveURL(new RegExp(`/projects/${projectId}/generation$`))
    await expect(page.getByRole('heading', { name: '准备开始生成' })).toBeVisible()
    await page.getByRole('button', { name: /开始生成/ }).click()
    await expect(page.getByText('代码生成已完成')).toBeVisible()
    await expect(page.getByRole('progressbar')).toHaveAttribute('aria-valuenow', '100')
    await page.getByRole('button', { name: /查看生成结果/ }).click()

    await expect(page).toHaveURL(new RegExp(`/projects/${projectId}/preview$`))
    await expect(page.getByRole('heading', { name: '尚未构建预览' })).toBeVisible()
    await page.getByRole('button', { name: /构建并预览/ }).click()
    await expect(page.locator('iframe[title="项目在线预览"]')).toHaveAttribute(
      'src',
      `/api/projects/${projectId}/preview/`,
    )
    await expect(page.frameLocator('iframe[title="项目在线预览"]').getByRole('heading', { name: 'E2E Preview' })).toBeVisible()

    await page.getByRole('link', { name: /^文件$/ }).click()
    await expect(page).toHaveURL(new RegExp(`/projects/${projectId}/files$`))
    await expect(page.getByText('浏览生成后的目录与源码，并导出完整项目文件。')).toBeVisible()
    await page.getByRole('button', { name: 'App.vue' }).click()
    await expect(page.getByText('<template>')).toBeVisible()
    await expect(page.locator('pre code')).toContainText('E2E 企业官网')
  })

  test('需求更新后旧蓝图会被明确阻止，并可根据最新需求重新生成', async ({ page }) => {
    await mockAuth(page)
    await mockStaleBlueprintApi(page)

    await page.goto(`/projects/${projectId}/blueprint`)
    await expect(page.getByRole('heading', { name: '需求已更新，需要重新生成方案' })).toBeVisible()
    await expect(page.getByText('需求已经更新，这份蓝图已失效，请根据最新需求重新生成')).toBeVisible()
    await expect(page.getByText('0 / 5')).toBeVisible()
    await expect(page.getByRole('button', { name: /确认并进入生成/ })).toBeDisabled()

    await page.getByRole('button', { name: /根据最新需求重新生成/ }).click()
    await expect(page.getByRole('heading', { name: 'AI 正在根据最新需求生成蓝图' })).toBeVisible()
    await expect(page.getByText('蓝图已生成，请逐项评审后确认')).toBeVisible()
    await expect(page.getByRole('heading', { name: '需求已更新，需要重新生成方案' })).toHaveCount(0)
    await expect(page.getByRole('button', { name: /确认并进入生成/ })).toBeEnabled()
  })
})

async function mockStaleBlueprintApi(page: Page) {
  const mustFeatures = ['产品展示', '线索表单', '客户案例', '响应式布局', '内容导航']
  const staleBlueprint = {
    id: 'blueprint-stale',
    project_id: projectId,
    status: 'superseded',
    version: 1,
    content: JSON.stringify({
      product_positioning: '旧版企业网站',
      tech_stack: 'Vue 3',
      pages: [{ name: '首页', route: '/' }],
    }),
  }
  const regeneratedBlueprint = {
    ...staleBlueprint,
    id: 'blueprint-new',
    status: 'generated',
    version: 2,
    content: JSON.stringify({
      schema_version: 2,
      app_type: 'landing_page',
      product_positioning: '根据最新需求生成的企业网站',
      tech_stack: 'Vue 3 + TypeScript',
      pages: [{ name: '首页', route: '/' }],
      features: mustFeatures.map((name, index) => ({ id: `F-${index + 1}`, name, priority: 'must', acceptance_criteria: [`可以使用${name}`] })),
      acceptance_criteria: mustFeatures.map((name) => `可以使用${name}`),
      test_plan: ['检查核心功能'],
    }),
  }

  await page.route(/https?:\/\/(?:localhost|127\.0\.0\.1)(?::\d+)?\/api\/.*/, async (route) => {
    const request = route.request()
    const path = new URL(request.url()).pathname
    if (request.method() === 'GET' && path === '/api/auth/profile') {
      await route.fulfill(apiOk({ id: 'user-e2e', username: 'e2e', email: 'e2e@example.com', role: 'admin' }))
      return
    }
    if (request.method() === 'GET' && path === `/api/projects/${projectId}`) {
      await route.fulfill(apiOk(project))
      return
    }
    if (request.method() === 'GET' && path === `/api/projects/${projectId}/requirements`) {
      await route.fulfill(apiOk({ content: JSON.stringify({ schema_version: 2, must_have_features: mustFeatures, goal: '企业官网', target_users: ['客户'], acceptance_criteria: ['可以浏览首页'] }) }))
      return
    }
    if (request.method() === 'GET' && path === `/api/projects/${projectId}/blueprint`) {
      await route.fulfill(apiOk(staleBlueprint))
      return
    }
    if (request.method() === 'POST' && path === `/api/projects/${projectId}/blueprint/generate`) {
      await new Promise((resolve) => setTimeout(resolve, 250))
      await route.fulfill(apiOk(regeneratedBlueprint))
      return
    }
    await route.fulfill(apiError(501, 50100, 'unmocked api', `${request.method()} ${path}`))
  })
}

async function mockAuth(page: Page) {
  await page.addInitScript(() => {
    window.localStorage.setItem('auth_token', 'e2e-token')
  })
}

async function mockTaskStream(page: Page) {
  await page.addInitScript(() => {
    class MockEventSource {
      static CONNECTING = 0
      static OPEN = 1
      static CLOSED = 2
      readonly CONNECTING = 0
      readonly OPEN = 1
      readonly CLOSED = 2
      readyState = 1
      url: string
      withCredentials = false
      onopen: ((event: Event) => void) | null = null
      onmessage: ((event: MessageEvent) => void) | null = null
      onerror: ((event: Event) => void) | null = null

      constructor(url: string | URL) {
        this.url = String(url)
        window.setTimeout(() => this.onopen?.(new Event('open')), 0)
        window.setTimeout(() => {
          this.onmessage?.(new MessageEvent('message', {
            data: JSON.stringify({ status: 'running', progress: 45, message: '正在生成页面结构' }),
          }))
        }, 20)
        window.setTimeout(() => {
          this.onmessage?.(new MessageEvent('message', {
            data: JSON.stringify({ type: 'task_completed', status: 'completed', progress: 100, message: '生成完成' }),
          }))
        }, 60)
      }

      addEventListener() {}
      removeEventListener() {}
      dispatchEvent() { return true }
      close() { this.readyState = 2 }
    }

    Object.defineProperty(window, 'EventSource', {
      configurable: true,
      writable: true,
      value: MockEventSource,
    })
  })
}

async function mockWorkflowApi(page: Page) {
  let hasBlueprint = false
  let requirementContent = ''

  await page.route(/https?:\/\/(?:localhost|127\.0\.0\.1)(?::\d+)?\/api\/.*/, async (route) => {
    const request = route.request()
    const url = new URL(request.url())
    const path = url.pathname
    const method = request.method()

    if (method === 'GET' && path === '/api/auth/profile') {
      await route.fulfill(apiOk({ id: 'user-e2e', username: 'e2e', email: 'e2e@example.com', role: 'admin' }))
      return
    }

    if (method === 'POST' && path === '/api/projects') {
      await route.fulfill(apiOk(project))
      return
    }

    if (method === 'GET' && path === `/api/projects/${projectId}`) {
      await route.fulfill(apiOk(project))
      return
    }

    if (method === 'GET' && path === `/api/projects/${projectId}/requirements`) {
      await route.fulfill(requirementContent
        ? apiOk({ id: 'requirements-e2e', project_id: projectId, content: requirementContent })
        : apiError(404, 40401, 'not found', '暂未保存需求'))
      return
    }

    if (method === 'PUT' && path === `/api/projects/${projectId}/requirements`) {
      requirementContent = (request.postDataJSON() as { content: string }).content
      await route.fulfill(apiOk({ id: 'requirements-e2e', project_id: projectId }))
      return
    }

    if (method === 'POST' && path === `/api/projects/${projectId}/requirements/assist`) {
      const current = (request.postDataJSON() as any).current_spec
      await route.fulfill(apiOk({
        spec: {
          ...current,
          schema_version: 2,
          app_type: 'landing_page',
          target_users: ['中小企业管理者'],
          primary_scenarios: ['查看产品能力', '提交合作需求'],
          must_have_features: ['产品展示', '线索表单'],
          should_have_features: ['客户案例'],
          screens: ['首页'],
          interaction_rules: ['表单提交后显示成功反馈'],
          data_and_storage: { persistence: '不用保存', backend_required: false },
          visual_preferences: { style: '简洁现代', primary_color: '蓝色科技' },
          responsive_targets: ['desktop', 'mobile'],
          non_functional_requirements: ['电脑和手机上都能正常使用'],
          acceptance_criteria: ['首页可以查看产品能力', '用户可以提交线索表单'],
          out_of_scope: [],
        },
        inferred_fields: ['target_users', 'must_have_features', 'acceptance_criteria'],
        questions: [],
        ready: true,
        summary: { product: current.goal, users: '中小企业管理者', features: ['产品展示', '线索表单'], success: ['首页可以查看产品能力'] },
      }))
      return
    }

    if (method === 'GET' && path === `/api/projects/${projectId}/blueprint`) {
      await route.fulfill(hasBlueprint ? apiOk(blueprint) : apiError(404, 40401, 'not found', '暂未生成蓝图'))
      return
    }

    if (method === 'POST' && path === `/api/projects/${projectId}/blueprint/generate`) {
      hasBlueprint = true
      await route.fulfill(apiOk(blueprint))
      return
    }

    if (method === 'POST' && path === `/api/projects/${projectId}/blueprint/confirm`) {
      blueprint.status = 'confirmed'
      await route.fulfill(apiOk({ ...blueprint, status: 'confirmed' }))
      return
    }

    if (method === 'POST' && path === '/api/tasks') {
      await route.fulfill(apiOk({ id: 'task-e2e' }))
      return
    }

    if (method === 'POST' && path === `/api/projects/${projectId}/build`) {
      await route.fulfill(apiOk({ preview_url: `/api/projects/${projectId}/preview/` }))
      return
    }

    if (method === 'GET' && path === `/api/projects/${projectId}/build`) {
      await route.fulfill(apiOk({ ready: false, preview_url: '' }))
      return
    }

    if (method === 'GET' && path === `/api/projects/${projectId}/preview/`) {
      await route.fulfill({
        status: 200,
        contentType: 'text/html; charset=utf-8',
        body: '<!doctype html><html><body><h1>E2E Preview</h1></body></html>',
      })
      return
    }

    if (method === 'GET' && path === `/api/projects/${projectId}/files`) {
      await route.fulfill(apiOk([
        { name: 'src', is_dir: true },
        { name: 'App.vue', is_dir: false },
        { name: 'package.json', is_dir: false },
      ]))
      return
    }

    if (method === 'GET' && path === `/api/projects/${projectId}/files/content`) {
      await route.fulfill(apiOk({
        path: url.searchParams.get('path') || 'App.vue',
        content: '<template>\n  <main>E2E 企业官网</main>\n</template>',
        size: 48,
        is_binary: false,
      }))
      return
    }

    await route.fulfill(apiError(501, 50100, 'unmocked api', `${method} ${path}`))
  })
}

function apiOk(data: unknown) {
  return {
    status: 200,
    contentType: 'application/json',
    body: JSON.stringify({ code: 0, message: 'success', data }),
  }
}

function apiError(status: number, code: number, message: string, error: string) {
  return {
    status,
    contentType: 'application/json',
    body: JSON.stringify({ code, message, error }),
  }
}
