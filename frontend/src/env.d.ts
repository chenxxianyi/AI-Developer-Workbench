/// <reference types="vite/client" />

declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<{}, {}, any>
  export default component
}

// Vite environment variables
interface ImportMetaEnv {
  /** 后端 API 基地址 */
  readonly VITE_API_BASE_URL: string
  /** 应用标题 */
  readonly VITE_APP_TITLE?: string
  /** 是否启用 Mock AI 模式 */
  readonly VITE_MOCK_AI?: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
