/**
 * M1 v0.2 E2E regression tests.
 *
 * The suite mocks backend API responses so it can validate the P0 frontend
 * report workflow without requiring Docker or a running Gin server.
 */
import { test, expect, type Page } from '@playwright/test'

const scoredReport = {
  id: 'report-1',
  tool_type: 'ui_review',
  title: 'UI Mock 可访问性审查',
  input_mode: 'code',
  status: 'succeeded',
  summary: '移动端上传区域和表单标签需要优先修复。',
  total_score: 82,
  grade: 'B',
  input_data: {},
  report_data: {
    scores: [
      { name: '可访问性', score: 78, max_score: 100, comment: '键盘操作需要补齐。' },
      { name: '响应式', score: 86, max_score: 100, comment: '移动端布局整体稳定。' },
    ],
    issues: [
      {
        title: '上传区域缺少键盘操作',
        severity: 'high',
        category: 'accessibility',
        problem: '上传区域只响应鼠标点击。',
        suggestion: '增加 Enter 和 Space 键触发。',
        action: '为上传区域添加 tabindex、role 和键盘事件。',
      },
      {
        title: '表单字段缺少显式 label',
        severity: 'medium',
        category: 'form',
        problem: '部分输入框依赖占位符说明。',
        suggestion: '补充 label 与描述关联。',
        action: '为所有字段添加 label 和 aria-describedby。',
      },
    ],
    recommendations: ['先修复键盘操作。', '再统一表单标签和错误提示。'],
    codex_prompt: '请修复 UIReviewPage 的上传键盘操作和表单 label。',
  },
  generated_files: [
    {
      id: 'file-1',
      filename: 'PATCH_SUMMARY.md',
      language: 'markdown',
      mime_type: 'text/markdown',
      size_bytes: 128,
    },
  ],
  created_at: '2026-07-07T07:00:00Z',
  updated_at: '2026-07-07T07:00:00Z',
}

const nonScoredReport = {
  id: 'report-2',
  tool_type: 'agent_config',
  title: 'Agent 配置生成',
  input_mode: 'manual',
  status: 'succeeded',
  summary: '已生成项目 AI 协作配置文件。',
  total_score: null,
  grade: null,
  input_data: {},
  report_data: {
    generated_files_content: {
      'AGENTS.md': '# Agent Guide',
    },
    recommendations: ['把 AGENTS.md 放到仓库根目录。'],
    codex_prompt: '请应用生成的 Agent 配置文件。',
  },
  generated_files: [
    {
      id: 'file-2',
      filename: 'AGENTS.md',
      language: 'markdown',
      mime_type: 'text/markdown',
      size_bytes: 64,
    },
  ],
  created_at: '2026-07-07T06:00:00Z',
  updated_at: '2026-07-07T06:00:00Z',
}

const projectDoctorReport = makeReport({
  id: 'report-3',
  tool_type: 'project_doctor',
  title: '项目诊断 Mock 报告',
  summary: '项目结构、测试和部署配置需要补齐。',
})

const apiDocReport = {
  ...nonScoredReport,
  id: 'report-4',
  tool_type: 'api_doc',
  title: 'API 文档 Mock 报告',
  summary: '已生成 Markdown 和 OpenAPI 文档草稿。',
  report_data: {
    modules: [{ name: 'reports', endpoints: [] }],
    recommendations: ['补充错误码说明。'],
    codex_prompt: '请根据 API 文档补齐前端客户端。',
  },
  generated_files: [
    {
      id: 'file-4',
      filename: 'API_DOCUMENTATION.md',
      language: 'markdown',
      mime_type: 'text/markdown',
      size_bytes: 96,
    },
  ],
}

const dbSchemaReport = makeReport({
  id: 'report-5',
  tool_type: 'db_schema',
  title: '数据库结构 Mock 报告',
  summary: 'reports 表需要增加组合索引并补充约束说明。',
  parent_report_id: 'report-0',
})

const reports = [scoredReport, nonScoredReport, projectDoctorReport, apiDocReport, dbSchemaReport]
const runReports = new Map([
  ['ui-review', scoredReport],
  ['project-doctor', projectDoctorReport],
  ['agent-config', nonScoredReport],
  ['api-doc', apiDocReport],
  ['db-schema', dbSchemaReport],
])

test.describe('v0.2 regression', () => {
  test('Dashboard and Settings show Mock Mode from system API', async ({ page }) => {
    await mockApi(page)

    await page.goto('/dashboard')
    await expect(page.getByText(/服务正常/)).toBeVisible()
    await expect(page.getByRole('heading', { name: '生成工作流' })).toBeVisible()

    await page.goto('/settings')
    await expect(page.getByRole('heading', { name: '系统状态', level: 1 })).toBeVisible()
    await expect(page.getByText('真实 AI 服务')).toBeVisible()
    await expect(page.getByText('AI 服务商')).toBeVisible()
    await expect(page.getByText('mock-text')).toBeVisible()
  })

  test('Reports page supports status/tool filtering and non-scored reports', async ({ page }) => {
    await mockApi(page)

    await page.goto('/reports')
    await expect(page.getByRole('heading', { name: '报告列表' })).toBeVisible()
    await expect(page.getByRole('heading', { name: scoredReport.title, level: 3 })).toBeVisible()
    await expect(page.getByRole('heading', { name: nonScoredReport.title, level: 3 })).toBeVisible()
    await expect(page.getByText('无需评分').first()).toBeVisible()

    await page.getByLabel('按状态筛选').selectOption('succeeded')
    await expect(page).toHaveURL(/status=succeeded/)

    await page.getByLabel('按工具类型筛选').selectOption('ui_review')
    await expect(page).toHaveURL(/tool_type=ui_review/)
    await expect(page.getByRole('heading', { name: scoredReport.title, level: 3 })).toBeVisible()
    await expect(page.getByRole('heading', { name: nonScoredReport.title, level: 3 })).toHaveCount(0)
  })

  test('Report detail supports hierarchy, downloads, delete cancel, and direct 404', async ({ page }) => {
    const api = await mockApi(page)

    await page.goto('/reports/report-1')
    await expect(page.getByRole('heading', { name: scoredReport.title })).toBeVisible()
    await expect(page.getByText('结论摘要')).toBeVisible()
    await expect(page.getByText('最严重问题')).toBeVisible()
    await expect(page.getByText('行动计划')).toBeVisible()
    await expect(page.getByText('AI 修复 Prompt')).toBeVisible()
    await expect(page.getByText('生成文件')).toBeVisible()

    await expect(page.getByRole('button', { name: '下载 Markdown 报告' })).toBeVisible()
    const markdownResponse = await page.evaluate(async () => {
      const response = await fetch('http://localhost:8080/api/reports/report-1/export?format=markdown')
      return {
        disposition: response.headers.get('content-disposition'),
        text: await response.text(),
      }
    })
    expect(markdownResponse.disposition).toContain('ui_review_report.md')
    expect(markdownResponse.text).toContain(scoredReport.title)

    const generatedFileLink = page.locator('a[download="PATCH_SUMMARY.md"]')
    await expect(generatedFileLink).toHaveAttribute('href', '/api/reports/report-1/files/PATCH_SUMMARY.md')

    await page.getByRole('button', { name: '删除报告' }).click()
    await expect(page.getByRole('dialog', { name: '删除报告' })).toBeVisible()
    await page.getByRole('button', { name: '取消' }).click()
    await expect(page.getByRole('dialog', { name: '删除报告' })).toHaveCount(0)
    expect(api.deleteRequests).toEqual([])

    await page.goto('/reports/missing-report')
    await expect(page.getByRole('heading', { name: '报告未找到' })).toBeVisible()
  })

  test('Project detail scopes history, trend, artifacts, and tool entry to one project', async ({ page }) => {
    await mockApi(page)
    await page.goto('/projects/project-1')

    await expect(page.getByRole('heading', { name: 'AI Developer Workbench', level: 1 })).toBeVisible()
    await expect(page.getByRole('heading', { name: '质量趋势（近 30 天）' })).toBeVisible()
    await expect(page.getByRole('heading', { name: '最新产物' })).toBeVisible()
    await expect(page.getByText('migration.sql')).toBeVisible()
    await expect(page.getByRole('heading', { name: '项目报告' })).toBeVisible()
    await expect(page.getByText('UI Mock 可访问性审查')).toBeVisible()

    await expect(page.locator('a[href*="project_id=project-1"]', { hasText: '数据库结构审查' })).toHaveCount(1)
    await expect(page.getByRole('link', { name: '基于 UI Mock 可访问性审查 复查' })).toHaveAttribute('href', /parent_report_id=report-1/)
    await expect(page.getByRole('link', { name: '对比 数据库结构 Mock 报告' })).toHaveAttribute('href', '/reports/report-0/compare/report-5')
    await expect(page.locator('a[href*="project_id=project-1"]', { hasText: 'Agent 配置生成' })).toHaveCount(2)
  })

  test('Tool pages are accessible', async ({ page }) => {
    await mockApi(page)

    const tools = [
      '/tools/ui-review',
      '/tools/project-doctor',
      '/tools/agent-config',
      '/tools/api-doc',
      '/tools/db-schema',
    ]

    for (const path of tools) {
      await page.goto(path)
      await expect(page.locator('h1, h2').first()).toBeVisible()
    }
  })

  test('Mock tool run APIs return reports that open in detail view', async ({ page }) => {
    await mockApi(page)
    await page.goto('/dashboard')

    for (const [toolSlug, report] of runReports) {
      const result = await page.evaluate(async (path) => {
        const response = await fetch(path, { method: 'POST', body: new FormData() })
        return response.json()
      }, `/api/tools/${toolSlug}/run`)

      expect(result.code).toBe(0)
      expect(result.data.id).toBe(report.id)

      await page.goto(`/reports/${report.id}`)
      await expect(page.getByRole('heading', { name: report.title })).toBeVisible()
    }
  })
})

async function mockApi(page: Page) {
  await page.addInitScript(() => {
    window.localStorage.setItem('auth_token', 'e2e-token')
  })

  const deleteRequests: string[] = []

  await page.route(/https?:\/\/(?:localhost|127\.0\.0\.1)(?::\d+)?\/api\/.*/, async (route) => {
    const request = route.request()
    const url = new URL(request.url())
    const path = url.pathname
    const method = request.method()

    if (path === '/api/auth/profile') {
      await route.fulfill(apiOk({ id: 'user-e2e', username: 'e2e', email: 'e2e@example.com', role: 'admin' }))
      return
    }

    if (path === '/api/health') {
      await route.fulfill(apiOk({ status: 'ok', timestamp: '2026-07-07T07:00:00Z' }))
      return
    }

    if (path === '/api/system/status') {
      await route.fulfill(apiOk({
        healthy: true,
        provider: 'openai',
        text_model: 'mock-text',
        vision_model: 'mock-vision',
        upload_limits: {
          image_max_bytes: 20_971_520,
          zip_max_bytes: 104_857_600,
          zip_max_files: 120,
          zip_max_total_bytes: 307_200_000,
        },
      }))
      return
    }

    if (path === '/api/dashboard/stats') {
      await route.fulfill(apiOk({
        total_reports: reports.length,
        tool_usage: {
          ui_review: 1,
          project_doctor: 0,
          agent_config: 1,
          api_doc: 0,
          db_schema: 0,
        },
        average_score: 82,
        recent_reports: reports.map((report) => ({
          id: report.id,
          tool_type: report.tool_type,
          title: report.title,
          status: report.status,
          total_score: report.total_score,
          grade: report.grade,
          summary: report.summary,
          created_at: report.created_at,
        })),
      }))
      return
    }

    if (method === 'GET' && path === '/api/projects') {
      await route.fulfill(apiOk({
        items: [{
          id: 'project-1',
          name: 'AI Developer Workbench',
          description: '本地 AI 质量工作台',
          repo_url: 'https://example.com/workbench',
          report_count: 2,
          average_score: 82,
          created_at: '2026-07-07T05:00:00Z',
          updated_at: '2026-07-07T07:00:00Z',
        }],
        total: 1,
        page: 1,
        page_size: 100,
      }))
      return
    }

    if (method === 'GET' && path === '/api/projects/project-1') {
      await route.fulfill(apiOk({
        id: 'project-1',
        name: 'AI Developer Workbench',
        description: '本地 AI 质量工作台',
        repo_url: 'https://example.com/workbench',
        frontend_stack: 'Vue 3 + TypeScript',
        backend_stack: 'Go + Gin',
        database: 'MySQL',
        ui_style: '简洁现代',
        coding_rules: '保持测试覆盖。',
        created_at: '2026-07-07T05:00:00Z',
        updated_at: '2026-07-07T07:00:00Z',
      }))
      return
    }

    if (method === 'GET' && path === '/api/projects/project-1/stats') {
      await route.fulfill(apiOk({
        total_reports: 2,
        average_score: 82,
        tool_usage: {
          ui_review: 1,
          project_doctor: 0,
          agent_config: 0,
          api_doc: 0,
          db_schema: 1,
        },
        high_severity_count: 1,
        recent_reports: reports,
        quality_trend: [{
          bucket: '2026-07-07',
          report_count: 2,
          average_score: 82,
          high_severity_count: 1,
        }],
        latest_artifacts: [{
          tool_type: 'db_schema',
          report_id: 'report-1',
          report_title: 'Schema review',
          filename: 'migration.sql',
          mime_type: 'text/x-sql',
          created_at: '2026-07-07T07:00:00Z',
        }],
      }))
      return
    }

    if (method === 'GET' && path === '/api/projects/project-1/reports') {
      await route.fulfill(apiOk({
        items: reports,
        total: reports.length,
        page: Number(url.searchParams.get('page') || 1),
        page_size: Number(url.searchParams.get('page_size') || 10),
      }))
      return
    }

    if (path === '/api/tools') {
      await route.fulfill(apiOk([
        tool('ui_review', 'accent', 1),
        tool('project_doctor', 'success', 0),
        tool('agent_config', 'warning', 1),
        tool('api_doc', 'danger', 0),
        tool('db_schema', 'accent', 0),
      ]))
      return
    }

    const runMatch = path.match(/^\/api\/tools\/([^/]+)\/run$/)
    if (runMatch && method === 'POST') {
      const report = runReports.get(runMatch[1])
      if (report) {
        await route.fulfill(apiOk(report))
        return
      }

      await route.fulfill(apiError(400, 40004, 'invalid tool type', runMatch[1]))
      return
    }

    if (method === 'GET' && path === '/api/reports') {
      const status = url.searchParams.get('status')
      const toolType = url.searchParams.get('tool_type')
      const filteredReports = reports.filter((report) => {
        if (status && report.status !== status) return false
        if (toolType && report.tool_type !== toolType) return false
        return true
      })

      await route.fulfill(apiOk({
        items: filteredReports,
        total: filteredReports.length,
        page: Number(url.searchParams.get('page') || 1),
        page_size: Number(url.searchParams.get('page_size') || 10),
      }))
      return
    }

    if (method === 'GET' && path === '/api/reports/report-1/export') {
      await route.fulfill({
        status: 200,
        contentType: 'text/markdown; charset=utf-8',
        headers: {
          'content-disposition': 'attachment; filename="ui_review_report.md"; filename*=UTF-8\'\'ui_review_report.md',
          'access-control-expose-headers': 'Content-Disposition',
        },
        body: '# UI Mock 可访问性审查\n\n移动端上传区域和表单标签需要优先修复。\n',
      })
      return
    }

    if (method === 'GET' && path === '/api/reports/report-1/files/PATCH_SUMMARY.md') {
      await route.fulfill({
        status: 200,
        contentType: 'text/markdown; charset=utf-8',
        headers: {
          'content-disposition': 'attachment; filename="PATCH_SUMMARY.md"; filename*=UTF-8\'\'PATCH_SUMMARY.md',
          'access-control-expose-headers': 'Content-Disposition',
        },
        body: '# Patch Summary\n\nAdd keyboard upload support.\n',
      })
      return
    }

    const reportMatch = path.match(/^\/api\/reports\/([^/]+)$/)
    if (reportMatch && method === 'GET') {
      const report = reports.find((item) => item.id === reportMatch[1])
      if (report) {
        await route.fulfill(apiOk(report))
        return
      }

      await route.fulfill(apiError(404, 40401, 'report not found', 'Report not found'))
      return
    }

    if (reportMatch && method === 'DELETE') {
      deleteRequests.push(reportMatch[1])
      await route.fulfill(apiOk(null))
      return
    }

    await route.fulfill(apiError(500, 50001, 'unhandled mock route', `${method} ${path}`))
  })

  return { deleteRequests }
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

function tool(tool_type: string, color: string, usage_count: number) {
  return {
    tool_type,
    name: tool_type,
    description: tool_type,
    icon: 'wrench',
    color,
    usage_count,
  }
}

function makeReport(input: {
  id: string
  tool_type: string
  title: string
  summary: string
  parent_report_id?: string
}) {
  return {
    id: input.id,
    tool_type: input.tool_type,
    title: input.title,
    input_mode: 'mock',
    status: 'succeeded',
    summary: input.summary,
    total_score: 78,
    grade: 'B',
    input_data: {},
    parent_report_id: input.parent_report_id,
    report_data: {
      scores: [
        { name: '完整性', score: 78, max_score: 100, comment: '主体结构可用，仍有补强空间。' },
      ],
      issues: [
        {
          title: '缺少关键验收证据',
          severity: 'medium',
          category: 'quality',
          problem: '当前报告缺少部分自动化验证说明。',
          suggestion: '补齐测试和验收命令。',
          action: '增加对应测试并记录命令输出。',
        },
      ],
      recommendations: ['优先补齐自动化验证。'],
      codex_prompt: '请补齐该工具的验收测试和文档。',
    },
    generated_files: [],
    created_at: '2026-07-07T05:00:00Z',
    updated_at: '2026-07-07T05:00:00Z',
  }
}
