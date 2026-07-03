/**
 * M1 v0.2 E2E Smoke Test
 *
 * Requires: Mock Mode backend running on localhost:8080, frontend dev server on :5173.
 * Run: npm run test:e2e
 */
import { test, expect } from '@playwright/test'

test.describe('v0.2 regression', () => {
  test('Dashboard loads and shows mock mode', async ({ page }) => {
    await page.goto('/')
    // Dashboard should load with key elements.
    await expect(page.locator('text=AI Developer Workbench')).toBeVisible({ timeout: 10_000 })
    // Mock mode indicator should be present.
    await expect(page.locator('text=演示')).toBeVisible({ timeout: 10_000 })
  })

  test('Reports page loads with empty state', async ({ page }) => {
    await page.goto('/reports')
    await expect(page.locator('text=报告列表')).toBeVisible({ timeout: 10_000 })
  })

  test('Settings page shows Mock Mode status', async ({ page }) => {
    await page.goto('/settings')
    await expect(page.locator('text=系统状态')).toBeVisible({ timeout: 10_000 })
    await expect(page.locator('text=演示模式')).toBeVisible({ timeout: 10_000 })
  })

  test('Tool pages are accessible', async ({ page }) => {
    const tools = [
      '/tools/ui-review',
      '/tools/project-doctor',
      '/tools/agent-config',
      '/tools/api-doc',
      '/tools/db-schema',
    ]
    for (const path of tools) {
      await page.goto(path)
      // Each tool page should have a title.
      await expect(page.locator('h1, h2').first()).toBeVisible({ timeout: 10_000 })
    }
  })

  test('Report detail direct URL shows error state', async ({ page }) => {
    await page.goto('/reports/nonexistent-id')
    await expect(page.locator('text=报告未找到')).toBeVisible({ timeout: 10_000 })
  })
})
