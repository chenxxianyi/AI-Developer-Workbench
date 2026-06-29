import { createI18n } from 'vue-i18n'

export const DEFAULT_LOCALE = 'zh-CN'
export const LOCALE_STORAGE_KEY = 'ai-workbench-locale'

export type AppLocale = 'zh-CN' | 'en-US'

export const supportedLocales: Array<{ code: AppLocale; label: string; shortLabel: string }> = [
  { code: 'zh-CN', label: '中文', shortLabel: '中' },
  { code: 'en-US', label: 'English', shortLabel: 'EN' },
]

export const messages = {
  'zh-CN': {
    landing: {
      nav: {
        home: '首页',
        tools: '工具',
        workflow: '工作流',
        features: '特性',
        start: '开始使用',
        language: '切换语言',
      },
      hero: {
        title: '更高质量地交付 AI 生成项目',
        description: '在一个开发者工作台中完成 UI 质量审查、项目结构检查、AGENTS.md 生成、API 文档构建和数据库结构优化。',
        enter: '进入工作台',
        learnMore: '了解更多',
      },
      stats: {
        tools: '5 核心工具',
        analysis: '智能分析',
        export: '一键导出',
      },
      tools: {
        title: '五大核心工具',
        description: '专为 AI Coding 开发者设计，覆盖 UI、项目结构、Agent 配置、API 文档和数据库设计全流程',
        items: {
          uiReview: {
            name: 'UI 质量审查',
            subtitle: 'UI 质量审查',
            description: '根据截图和前端代码审查 UI 质量，识别模板化痕迹，评估设计一致性，提供改进建议。',
          },
          projectDoctor: {
            name: '项目诊断',
            subtitle: '项目结构检查',
            description: '静态分析项目 ZIP，检查工程结构、依赖管理、代码规范和潜在风险，生成健康报告。',
          },
          agentConfig: {
            name: 'Agent 配置生成',
            subtitle: 'AI Agent 配置生成',
            description: '根据项目特征生成 AGENTS.md、TASK_PLAN.md 等配置文件，优化 AI Coding 效果。',
          },
          apiDoc: {
            name: 'API 文档生成',
            subtitle: 'API 文档生成',
            description: '从代码或项目 ZIP 生成 Markdown/OpenAPI 文档，支持多种后端框架和输出格式。',
          },
          dbSchema: {
            name: '数据库结构审查',
            subtitle: '数据库结构审查',
            description: '审查 SQL、GORM、Prisma 等数据库定义，评估表结构、索引、性能和安全问题。',
          },
        },
      },
      workflow: {
        title: '统一工作流',
        description: '从输入到报告，每个工具遵循一致的交互闭环',
        steps: {
          choose: {
            title: '选择工具',
            description: '从 Dashboard 选择合适的分析工具，填写表单或上传文件',
          },
          analyze: {
            title: '智能分析',
            description: '后端进行安全处理和 AI 分析，生成结构化报告和评分',
          },
          report: {
            title: '查看报告',
            description: '查看评分、问题列表和改进建议，了解具体优化方向',
          },
          export: {
            title: '复制与导出',
            description: '复制 Codex Prompt 到 AI Coding 工具，下载生成文件或报告',
          },
        },
      },
      features: {
        title: '专为开发者设计',
        kicker: '质量交付系统',
        items: {
          product: {
            title: '真实产品感',
            description: '避免模板化 SaaS 风格，克制、专业的开发者体验',
          },
          safety: {
            title: '安全优先',
            description: '静态分析不执行代码，完整的安全处理和输入校验',
          },
          output: {
            title: '结构化输出',
            description: '评分、问题、建议清晰分离，支持 Markdown 和 OpenAPI',
          },
          integration: {
            title: '一键集成',
            description: 'Codex Prompt 直接复制到 Claude Code、Cursor 等 AI 工具',
          },
        },
      },
      sampleReport: {
        title: '示例报告',
        promptReady: 'Codex Prompt 已就绪',
        qualityLabel: '质量',
        score: '总分',
        issues: '发现问题',
        high: '高：检测到 AI 模板化风险',
        medium: '中：按钮间距不一致',
        low: '低：缺少悬停状态',
      },
      cta: {
        title: '开始提升你的 AI 项目质量',
        description: '立即进入工作台，体验智能分析工具带来的效率提升',
        button: '进入 Dashboard',
        command: 'ai-workbench --ship',
      },
      footer: {
        tagline: 'MVP 0.1.0 · 专为 AI Coding 开发者设计',
      },
    },
  },
  'en-US': {
    landing: {
      nav: {
        home: 'Home',
        tools: 'Tools',
        workflow: 'Workflow',
        features: 'Features',
        start: 'Get Started',
        language: 'Switch language',
      },
      hero: {
        title: 'Build better AI-generated projects',
        description: 'Review UI quality, inspect project structure, generate AGENTS.md, build API docs, and improve database schemas in one developer workbench.',
        enter: 'Enter Workbench',
        learnMore: 'Learn More',
      },
      stats: {
        tools: '5 Core Tools',
        analysis: 'Smart Analysis',
        export: 'One-click Export',
      },
      tools: {
        title: 'Five Core Tools',
        description: 'Designed for AI coding developers, covering UI, project structure, agent configuration, API docs, and database design workflows.',
        items: {
          uiReview: {
            name: 'UI Review',
            subtitle: 'UI Quality Review',
            description: 'Review UI quality from screenshots and frontend code, identify template-like patterns, evaluate design consistency, and provide improvement suggestions.',
          },
          projectDoctor: {
            name: 'Project Doctor',
            subtitle: 'Project Structure Check',
            description: 'Statically analyze project ZIP files, inspect engineering structure, dependency management, code conventions, and potential risks, then generate a health report.',
          },
          agentConfig: {
            name: 'Agent Config Studio',
            subtitle: 'AI Agent Config Generator',
            description: 'Generate AGENTS.md, TASK_PLAN.md, and related configuration files based on project characteristics to improve AI coding outcomes.',
          },
          apiDoc: {
            name: 'API Doc Builder',
            subtitle: 'API Documentation Generator',
            description: 'Generate Markdown or OpenAPI documentation from source code or project ZIP files, with support for multiple backend frameworks and output formats.',
          },
          dbSchema: {
            name: 'DB Schema Review',
            subtitle: 'Database Schema Review',
            description: 'Review SQL, GORM, Prisma, and other database definitions to evaluate table structure, indexes, performance, and security concerns.',
          },
        },
      },
      workflow: {
        title: 'Unified Workflow',
        description: 'From input to report, every tool follows a consistent interaction loop.',
        steps: {
          choose: {
            title: 'Choose a Tool',
            description: 'Select the right analysis tool from Dashboard, then fill in the form or upload files.',
          },
          analyze: {
            title: 'Smart Analysis',
            description: 'The backend safely processes inputs and runs AI analysis to generate structured reports and scores.',
          },
          report: {
            title: 'Review the Report',
            description: 'Inspect scores, issue lists, and recommendations to understand the next optimization steps.',
          },
          export: {
            title: 'Copy and Export',
            description: 'Copy Codex prompts to AI coding tools, or download generated files and reports.',
          },
        },
      },
      features: {
        title: 'Designed for Developers',
        kicker: 'Quality delivery system',
        items: {
          product: {
            title: 'Real Product Feel',
            description: 'Avoid template-like SaaS patterns with a restrained, professional developer experience.',
          },
          safety: {
            title: 'Safety First',
            description: 'Static analysis does not execute code, with complete safety handling and input validation.',
          },
          output: {
            title: 'Structured Output',
            description: 'Scores, issues, and suggestions are clearly separated, with Markdown and OpenAPI support.',
          },
          integration: {
            title: 'One-click Integration',
            description: 'Copy Codex prompts directly into Claude Code, Cursor, and other AI tools.',
          },
        },
      },
      sampleReport: {
        title: 'Sample Report',
        promptReady: 'Codex Prompt ready',
        qualityLabel: 'Quality',
        score: 'Total Score',
        issues: 'Issues Found',
        high: 'High: AI template risk detected',
        medium: 'Medium: Inconsistent button spacing',
        low: 'Low: Missing hover states',
      },
      cta: {
        title: 'Start improving your AI project quality',
        description: 'Enter the workbench now and experience the efficiency of intelligent analysis tools.',
        button: 'Enter Dashboard',
        command: 'ai-workbench --ship',
      },
      footer: {
        tagline: 'MVP 0.1.0 · Designed for AI Coding Developers',
      },
    },
  },
} as const

export const i18n = createI18n({
  legacy: false,
  locale: DEFAULT_LOCALE,
  fallbackLocale: DEFAULT_LOCALE,
  messages,
})
