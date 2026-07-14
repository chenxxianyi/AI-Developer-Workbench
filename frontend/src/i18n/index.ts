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
        studio: '生成工作室',
        tools: '质量工具',
        workflow: '工作流',
        start: '开始使用',
        language: '切换语言',
      },
      hero: {
        title: '更高质量地交付 AI 生成项目',
        eyebrow: 'AI 项目生成与审查',
        titleLead: '从想法到',
        titleAccent: '可运行项目',
        titleTail: '更简单',
        description: '输入需求，生成蓝图、代码和预览；再用质量工具完成交付前检查。',
        studioCta: '进入生成工作室',
        toolsCta: '查看质量工具',
        proof: {
          blueprint: '蓝图可确认',
          preview: '生成可预览',
          quality: '质量可检查',
        },
        preview: {
          title: '生成工作室',
          subtitle: '智能运营平台 · Vue 3 + Go',
          running: '生成中',
          pipeline: '项目流水线',
          current: '当前阶段',
          generating: '生成项目结构与核心代码',
          log: '正在装配页面、接口与配置文件…',
          files: '项目文件',
          pages: '页面',
          checks: '检查项',
        },
      },
      overview: {
        kicker: '核心能力',
        title: '生成与审查，一个工作台',
        description: '创建新项目，或检查已有项目。',
      },
      studio: {
        badge: 'Generation Studio',
        title: '生成工作室',
        description: '确认需求和蓝图，生成代码，预览并导出项目。',
        cta: '创建新项目',
        stages: {
          project: '创建项目',
          requirements: '梳理需求',
          blueprint: '确认蓝图',
          generation: '生成代码',
          preview: '在线预览',
          delivery: '文件交付',
        },
        outputs: {
          source: '完整源码',
          preview: '可访问预览',
          blueprint: '产品与技术蓝图',
          package: '项目文件包',
        },
      },
      toolbox: {
        badge: 'Quality Tools',
        title: '质量工具箱',
        description: '五个独立工具，快速发现 UI、工程和文档问题。',
        cta: '展开全部工具',
      },
      tools: {
        kicker: '质量工具',
        title: '五个工具，解决交付问题',
        description: '按需使用，结果清晰可导出。',
        open: '打开工具',
        items: {
          uiReview: {
            name: 'UI 质量审查',
            description: '检查视觉、响应式和设计一致性。',
          },
          projectDoctor: {
            name: '项目诊断',
            description: '检查项目结构、依赖和工程风险。',
          },
          agentConfig: {
            name: 'Agent 配置生成',
            description: '生成 AI Coding 协作配置。',
          },
          apiDoc: {
            name: 'API 文档生成',
            description: '从代码生成 Markdown 和 OpenAPI 文档。',
          },
          dbSchema: {
            name: '数据库结构审查',
            description: '检查表结构、索引和性能问题。',
          },
        },
      },
      workflow: {
        kicker: '生成工作流',
        title: '四步完成交付',
        description: '需求、蓝图、生成、交付，过程清晰可控。',
        cta: '开始一条新流程',
        steps: {
          brief: {
            title: '定义需求',
            description: '填写目标、功能和技术栈。',
          },
          blueprint: {
            title: '确认蓝图',
            description: '确认页面结构与技术方案。',
          },
          generate: {
            title: '生成代码',
            description: '查看进度与实时日志。',
          },
          ship: {
            title: '预览交付',
            description: '预览结果并导出文件。',
          },
        },
      },
      features: {
        kicker: '为真实交付设计',
        title: '比“一键生成”更可控',
        description: '我们把生成、检查和交付放进同一套项目上下文，重点不是制造更多内容，而是让结果更容易被理解、验证和继续开发。',
        items: {
          workspace: {
            title: '项目化工作区',
            description: '需求、蓝图、生成、预览、文件和报告围绕同一个项目持续沉淀。',
          },
          safety: {
            title: '静态分析优先',
            description: '诊断上传项目时不执行源代码，并提供输入校验和清晰的任务状态。',
          },
          control: {
            title: '过程透明',
            description: '关键阶段、生成进度与日志都可见，避免黑盒式等待和结果失控。',
          },
          output: {
            title: '可继续开发的产物',
            description: '交付源码、预览、结构化报告和配置文件，而不仅是一段聊天回复。',
          },
        },
      },
      cta: {
        kicker: 'Start building',
        title: '开始生成你的下一个项目',
        description: '从需求开始，交付可运行的代码。',
        primary: '创建新项目',
        secondary: '进入工作台',
      },
      footer: {
        tagline: 'AI 项目生成与质量工作台',
      },
    },
  },
  'en-US': {
    landing: {
      nav: {
        home: 'Home',
        studio: 'Generation Studio',
        tools: 'Quality Tools',
        workflow: 'Workflow',
        start: 'Get Started',
        language: 'Switch language',
      },
      hero: {
        title: 'Build better AI-generated projects',
        eyebrow: 'AI project generation and review',
        titleLead: 'From idea to',
        titleAccent: 'working product',
        titleTail: 'made simple',
        description: 'Turn a brief into blueprints, code, and previews, then run a focused quality review.',
        studioCta: 'Open Generation Studio',
        toolsCta: 'Explore Quality Tools',
        proof: {
          blueprint: 'Reviewable blueprint',
          preview: 'Live preview',
          quality: 'Quality checks',
        },
        preview: {
          title: 'Generation Studio',
          subtitle: 'Operations platform · Vue 3 + Go',
          running: 'Generating',
          pipeline: 'Project pipeline',
          current: 'Current stage',
          generating: 'Generating structure and core code',
          log: 'Assembling pages, APIs, and configuration…',
          files: 'Files',
          pages: 'Pages',
          checks: 'Checks',
        },
      },
      overview: {
        kicker: 'Core capabilities',
        title: 'Build and review in one workbench',
        description: 'Create a new project or review an existing one.',
      },
      studio: {
        badge: 'Generation Studio',
        title: 'Generation Studio',
        description: 'Confirm the brief, generate code, preview, and export.',
        cta: 'Create Project',
        stages: {
          project: 'Create project',
          requirements: 'Define requirements',
          blueprint: 'Confirm blueprint',
          generation: 'Generate code',
          preview: 'Online preview',
          delivery: 'Deliver files',
        },
        outputs: {
          source: 'Complete source',
          preview: 'Accessible preview',
          blueprint: 'Product and technical blueprint',
          package: 'Project package',
        },
      },
      toolbox: {
        badge: 'Quality Tools',
        title: 'Quality Toolbox',
        description: 'Five focused tools for UI, engineering, and documentation issues.',
        cta: 'View all tools',
      },
      tools: {
        kicker: 'Quality tools',
        title: 'Five tools for delivery quality',
        description: 'Use what you need and export clear results.',
        open: 'Open Tool',
        items: {
          uiReview: {
            name: 'UI Quality Review',
            description: 'Check visuals, responsiveness, and consistency.',
          },
          projectDoctor: {
            name: 'Project Doctor',
            description: 'Check structure, dependencies, and engineering risks.',
          },
          agentConfig: {
            name: 'Agent Config Studio',
            description: 'Generate AI coding collaboration files.',
          },
          apiDoc: {
            name: 'API Doc Builder',
            description: 'Generate Markdown and OpenAPI docs from code.',
          },
          dbSchema: {
            name: 'Database Schema Review',
            description: 'Check schemas, indexes, and performance.',
          },
        },
      },
      workflow: {
        kicker: 'Generation workflow',
        title: 'Ship in four steps',
        description: 'Brief, blueprint, generation, and delivery—clearly connected.',
        cta: 'Start a New Workflow',
        steps: {
          brief: {
            title: 'Define the brief',
            description: 'Add goals, features, and stack.',
          },
          blueprint: {
            title: 'Confirm the blueprint',
            description: 'Confirm pages and technical direction.',
          },
          generate: {
            title: 'Generate code',
            description: 'Follow progress and live logs.',
          },
          ship: {
            title: 'Preview and ship',
            description: 'Preview results and export files.',
          },
        },
      },
      features: {
        kicker: 'Designed for real delivery',
        title: 'More control than one-click generation',
        description: 'Generation, review, and delivery share the same project context. The goal is not more content, but results that are easier to understand, validate, and continue developing.',
        items: {
          workspace: {
            title: 'Project-based workspace',
            description: 'Requirements, blueprints, generation, previews, files, and reports stay connected to one project.',
          },
          safety: {
            title: 'Static analysis first',
            description: 'Uploaded projects are reviewed without executing source code, with clear validation and task states.',
          },
          control: {
            title: 'Transparent process',
            description: 'Critical stages, generation progress, and logs remain visible instead of hiding behind a black box.',
          },
          output: {
            title: 'Artifacts you can keep building',
            description: 'Receive source code, previews, structured reports, and configuration files—not just a chat response.',
          },
        },
      },
      cta: {
        kicker: 'Start building',
        title: 'Start your next AI project',
        description: 'Start with a brief and ship working code.',
        primary: 'Create New Project',
        secondary: 'Open Workbench',
      },
      footer: {
        tagline: 'AI project generation and quality workbench',
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
