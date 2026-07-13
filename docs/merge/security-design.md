# A-07 文件、构建和预览安全边界

## 项目工作区根目录

- 根目录：`<DATA_DIR>/workspaces/<project_id>/`
- 所有文件操作必须限定在该目录内
- 路径解析后必须校验前缀匹配

## 路径归一化规则

```go
func SafeResolve(workspaceRoot, userPath string) (string, error) {
    // 1. 清理路径（去除 ../、./、//）
    clean := filepath.Clean(userPath)
    // 2. 禁止绝对路径
    if filepath.IsAbs(clean) {
        return "", ErrAbsolutePath
    }
    // 3. 拼接并再次清理
    full := filepath.Clean(filepath.Join(workspaceRoot, clean))
    // 4. 校验前缀
    if !strings.HasPrefix(full, workspaceRoot) {
        return "", ErrPathTraversal
    }
    return full, nil
}
```

## 文件操作安全

### 读取

- 仅允许读取项目工作区内的文件
- 文件大小限制：文本文件 10MB，二进制文件不加载到内存
- 禁止读取文件名含 `..`、绝对路径、符号链接指向外部
- 返回相对路径，不暴露服务器绝对路径

### 写入

- 仅允许写入项目工作区
- 文件名禁止包含：`..`、`/`、`\`、`:`、`*`、`?`、`"`、`<`、`>`、`|`
- 文件大小限制：单文件 50MB

### 导出（ZIP）

- 流式打包，不全部加载到内存
- 排除文件/目录：`.env`、`.env.local`、`node_modules/`、`.git/`、`*_test.go`、`secrets/`
- 文件名安全编码（UTF-8）
- 导出大小限制：总计 500MB

## 构建安全

### 超时限制

```go
const (
    BuildTimeout    = 10 * time.Minute  // npm install + build
    InstallTimeout  = 5 * time.Minute   // npm install only
)
```

### 进程隔离

- 子进程执行，非当前进程内执行
- 构建完成后强制清理子进程树
- `context.WithTimeout` 控制超时

### 输出限制

- stdout/stderr 累计最大 10MB
- 超出截断并记录警告

### 并发控制

- 全局构建并发数：3
- 超出的请求排队等待，5 分钟超时

### 文件类型限制

- 仅构建以下类型项目：`package.json`（Node.js）存在
- 不执行任意 shell 脚本

## 预览安全

### 会话管理

```go
type PreviewSession struct {
    ID        string
    ProjectID string
    UserID    string
    Token     string    // JWT
    ExpiresAt time.Time // 2 小时
}
```

### 预览文件服务

- 只能访问 `workspaces/<project_id>/dist/` 内的文件
- 支持 SPA History Fallback（`index.html` 兜底）
- Cookie `preview_token` 鉴权
- 会话过期（2 小时）后拒绝访问

### Cookie 属性

```text
preview_token=<jwt>; Path=/api/v1/preview; HttpOnly; SameSite=Strict; Max-Age=7200
```

## 安全评审覆盖矩阵

| 操作 | 路径穿越 | 越权 | 资源限制 | 文件类型 | 会话管理 |
|------|:-------:|:----:|:-------:|:-------:|:-------:|
| 文件读取 | ✅ | ✅ | ✅ | ✅ | — |
| 文件写入 | ✅ | ✅ | ✅ | ✅ | — |
| ZIP 导出 | ✅ | ✅ | ✅ | ✅ | — |
| 项目构建 | — | ✅ | ✅ | — | — |
| 预览服务 | ✅ | ✅ | — | ✅ | ✅ |
