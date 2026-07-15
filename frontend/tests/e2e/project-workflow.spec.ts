import { test, expect, type Page, type Route } from '@playwright/test'

const projectId = 'project-e2e'

const project = {
  id: projectId,
  name: 'E2E 企业官网',
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
    product_positioning: '面向企业客户的品牌官网与线索收集站点',
    tech_stack: 'Vue 3 + TypeScript + Tailwind CSS',
    pages: [
      { name: '首页', route: '/' },
      { name: '解决方案', route: '/solutions' },
      { name: '联系我们', route: '/contact' },
    ],
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
    await page.getByLabel('项目目标').fill('上线一个可以展示产品能力并收集销售线索的企业官网。')
    await page.getByLabel('目标用户').fill('中小企业管理者')
    await page.getByLabel('核心功能（每行一个）').fill('产品展示\n线索表单\n客户案例')
    await page.getByRole('button', { name: /保存并生成蓝图/ }).click()

    await expect(page).toHaveURL(new RegExp(`/projects/${projectId}/blueprint$`))
    await expect(page.getByRole('heading', { name: '蓝图评审' })).toBeVisible()
    await page.getByRole('button', { name: /生成蓝图/ }).click()
    await expect(page.getByText('面向企业客户的品牌官网与线索收集站点')).toBeVisible()
    await expect(page.getByText('Vue 3 + TypeScript + Tailwind CSS')).toBeVisible()
    await page.getByRole('button', { name: /确认蓝图/ }).click()

    await expect(page).toHaveURL(new RegExp(`/projects/${projectId}/generation$`))
    await expect(page.getByRole('heading', { name: '准备开始生成' })).toBeVisible()
    await page.getByRole('button', { name: /开始生成/ }).click()
    await expect(page.getByText('代码生成已完成')).toBeVisible()
    await expect(page.getByRole('progressbar')).toHaveAttribute('aria-valuenow', '100')
    await page.getByRole('button', { name: /查看生成结果/ }).click()

    await expect(page).toHaveURL(new RegExp(`/projects/${projectId}/preview$`))
    await expect(page.getByRole('heading', { name: '尚未构建预览' })).toBeVisible()
    await page.getByRole('button', { name: /构建并预览/ }).click()
    await expect(page.locator('iframe[title="项目在线预览"]')).toHaveAttribute('src', /^data:text\/html/)

    await page.getByRole('link', { name: /^文件$/ }).click()
    await expect(page).toHaveURL(new RegExp(`/projects/${projectId}/files$`))
    await expect(page.getByText('浏览生成后的目录与源码，并导出完整项目文件。')).toBeVisible()
    await page.getByRole('button', { name: 'App.vue' }).click()
    await expect(page.getByText('<template>')).toBeVisible()
    await expect(page.locator('pre code')).toContainText('E2E 企业官网')
  })
})

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
      await route.fulfill(apiError(404, 40401, 'not found', '暂未保存需求'))
      return
    }

    if (method === 'PUT' && path === `/api/projects/${projectId}/requirements`) {
      await route.fulfill(apiOk({ id: 'requirements-e2e', project_id: projectId }))
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
      await route.fulfill(apiOk({ ...blueprint, status: 'confirmed' }))
      return
    }

    if (method === 'POST' && path === '/api/tasks') {
      await route.fulfill(apiOk({ id: 'task-e2e' }))
      return
    }

    if (method === 'POST' && path === `/api/projects/${projectId}/build`) {
      await route.fulfill(apiOk({ preview_url: 'data:text/html,%3Ch1%3EE2E%20Preview%3C%2Fh1%3E' }))
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
