# UI Review 粘贴截图上传设计

## 目标

让 UI Review 的“上传截图”区域支持用户截图后直接按 `Ctrl+V` 粘贴图片，减少保存文件再选择上传的步骤。

## 交互范围

仅上传截图区域响应粘贴事件。用户点击上传框使其获得焦点后，可以按 `Ctrl+V` 粘贴剪贴板中的图片。页面不做全局粘贴监听，避免在代码输入框、标题输入框等位置粘贴内容时被误拦截。

## 行为设计

- 上传框文案调整为“点击、拖拽或 Ctrl+V 粘贴截图”。
- 上传框可聚焦，并在聚焦时显示与 hover 类似的强调边框。
- 粘贴事件从 `clipboardData.items` 中查找第一个 `image/*` 文件。
- 找到图片时阻止默认粘贴行为，并复用现有文件预览流程设置 `screenshotFile` 和 `screenshotPreview`。
- 剪贴板没有图片时不改变当前状态，也不打断用户。

## 文件与测试

- 修改 `frontend/src/pages/tools/UIReviewPage.vue`。
- 新增组件测试 `frontend/src/pages/tools/UIReviewPage.test.ts`，验证粘贴图片会显示预览，并验证上传区文案包含 `Ctrl+V`。

