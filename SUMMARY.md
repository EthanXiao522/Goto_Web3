# Goto Web3 — 项目开发总结

> 版本：v1.1 | 日期：2026-05-15

---

## 目录

- [1. 做了什么](#1-做了什么)
- [2. 输出了什么](#2-输出了什么)
- [3. 技术要点](#3-技术要点)
- [4. 经验总结](#4-经验总结)

---

## 1. 做了什么

本次迭代（v1.0 → v1.1）围绕 Dashboard、学习任务、学习计划书三大模块进行了全面优化和功能扩展：

### Dashboard 优化
- 4 张统计卡片调整为单行展示（`stats-grid-four`），各分配独立配色（teal/emerald/violet/amber）
- 恢复「每周进度」柱状图，柱高=完成率×1.2px，按 Phase 着色
- 甘特图移至每周进度下方，统一 bar-fill + bar-text + 网格线风格
- Landing/Demo/Dashboard 三页面甘特图风格统一

### 学习任务页（全新）
- 新增 `/tasks` 页面，233 个任务以树形可展开表格呈现（阶段→周→天→任务）
- 层级箭头展开/收起，子行自动缩进
- 内联编辑功能：点击「✎」弹出编辑框，Enter 保存 / Escape 取消
- 任务内容修改自动同步到 `sources/web3_infra_3month_plan.md` 源文档
- 侧边栏主菜单新增「📋 学习任务」入口

### 学习计划书重构
- 切换为服务端 Markdown 渲染（gomarkdown），替代客户端 marked.js
- 过滤 `[TOC]` 占位符，移除侧边栏目录
- 顶部新增📄源文档链接，点击在新标签页查看完整渲染
- 新增 `/handbook/source` 页面，将 md 文件渲染为 HTML 展示
- 修复 blockquote 内列表换行格式（源文件增加空 `>` 分隔行）

### 测试体系完善
- 新增 `service_test.go`（ProgressService + TaskService，4 tests）
- 新增 `middleware_test.go`（Auth + CORS 中间件，5 tests）
- 扩展 `handler_test.go`（TaskHandler + PhaseHandler + ProgressHandler，+9 tests）
- 扩展 `repository_test.go`（UpdateContent + UpdateFields + Batch + WithProgress，+6 tests）
- 单元测试从 15 增长到 40，功能测试 38 项，全部通过（88/88）

### 文档更新
- REQUIREMENTS.md v1.1：新增 FR-2.6~2.8、FR-4.6、FR-5.1~5.4、FR-6.1~6.6，验收标准 14→19 条
- DESIGN.md v1.1：更新架构图、模块表、前端页面表、测试策略、依赖清单
- TEST_REPORT.md v1.1：完整重写，40 单元 + 38 功能 = 88 测试，19 项验收标准

---

## 2. 输出了什么

| 类别 | 文件 | 说明 |
|------|------|------|
| 需求 | REQUIREMENTS.md v1.1 | 7 大功能模块，19 项验收标准 |
| 设计 | DESIGN.md v1.1 | 三层架构，9 模块，12 页面，6 表 |
| 测试 | TEST_REPORT.md v1.1 | 88 项测试，19/19 验收通过 |
| README | README.md v1.1 | 系统概述 + 部署说明 + 功能 + 代码说明 |
| 总结 | SUMMARY.md v1.1 | 本文档 |
| 后端新代码 | handler/task.go | UpdateContent + syncTaskToMd（含 regexp 回退匹配） |
| 后端新代码 | handler/page.go | LearningTasks + HandbookSource + Handbook 重构 |
| 后端新代码 | repository/task.go | UpdateContent 方法 |
| 后端新代码 | router/router.go | /tasks, /handbook/source, PUT tasks/:id/content |
| 前端新模板 | learning_tasks.html | 树形可展开表格 + 内联编辑 JS |
| 前端新模板 | handbook_source.html | md→HTML 全量渲染页 |
| 前端更新 | handbook.html | 服务端渲染 + 源文档链接 bar |
| 前端更新 | _sidebar.html | 新增「学习任务」菜单项 |
| 前端更新 | dashboard.html | 4 卡单行 + 周柱状图 + 甘特图 |
| 前端更新 | landing.html | 甘特图统一风格 + 移除阶段跳转 |
| 前端更新 | demo.html | 同步 Dashboard 布局 |
| 前端更新 | style.css | +180 行（stats-grid-four + task-tree + source-bar） |
| 依赖新增 | gomarkdown/markdown | 服务端 Markdown 渲染 |
| 测试新增 | service_test.go | 4 tests |
| 测试新增 | middleware_test.go | 5 tests |
| 测试扩展 | handler_test.go | +9 tests |
| 测试扩展 | repository_test.go | +6 tests |
| 源文件修复 | web3_infra_3month_plan.md | 3 处 blockquote 列表空行修复 |

---

## 3. 技术要点

### 3.1 Markdown 同步策略
任务内容编辑后同步到 md 源文件采用**两阶段匹配**：
1. **精确匹配**：在文件中查找旧内容字符串，直接替换
2. **回退匹配**：正则去除 markdown 链接后逐行匹配，找到对应行后替换

已知局限：DB 不存储 URL（仅链接文本），同步时 md 文件中的 `[text](url)` 会退化为纯文本。

### 3.2 树形表格渲染
学习任务页一次性加载全部 233 个任务（~165KB），使用 CSS 类控制展开/收起：
- `.tree-children { display: none }` / `.tree-children.open { display: block }`
- 层级缩进通过 `padding-left` 递进实现（phase:16px, week:40px, day:68px, task:96px）
- 内联编辑通过 DOM 操作替换 `<span>` 为 `<input>`，保存后还原

### 3.3 服务端 Markdown 渲染
从客户端 marked.js 切换到服务端 gomarkdown 的原因：
- `html/template` 在 `<script>` 标签内转义 markdown 特殊字符，导致解析异常
- 服务端渲染无编码问题，直接输出 `template.HTML`

gomarkdown 使用 `AutoHeadingIDs` 扩展自动生成标题锚点，`CommonExtensions` 支持表格/代码块/删除线等。

### 3.4 测试架构
测试采用**共享数据库**模式（非 mock），所有测试连接同一测试库：
- `setupXxx(t)` 辅助函数封装 config.Load + database.Connect + Migrate
- `t.Cleanup(func() { database.Close() })` 保证资源清理
- 测试间数据通过 `uniqueName()` 时间戳隔离用户，任务数据依赖种子库

---

## 4. 经验总结

### 做得好的
1. **服务端渲染避坑**：`html/template` 对 `<script>` 标签内容会进行转义，即使 `type="text/plain"` 也不可靠，服务端渲染是更稳妥的方案
2. **渐进式匹配**：md 同步采用「精确→回退」两级策略，兼顾性能和容错
3. **CSS 变量体系**：20+ 个自定义属性，配色切换仅需修改 `:root`，无需改动各处硬编码
4. **测试驱动修复**：先写测试再现问题，修复后验证，最后扩展测试覆盖新增方法

### 可以改进的
1. **md 同步的 URL 保留**：当前 DB 不存储完整 markdown 链接，同步时会丢失 URL。可考虑解析时保留原始 markdown 文本
2. **学习任务页性能**：233 个任务一次性渲染 165KB，未来可考虑虚拟滚动或分页
3. **测试隔离**：当前测试共用数据库，顺序敏感。可考虑每个测试用例独立事务或测试数据库
4. **前端 JS 模块化**：内联 JS 散落在模板中（handbook、learning_tasks），可统一抽取到 JS 文件

### 关键决策记录
| 决策 | 选项 | 选择 | 原因 |
|------|------|------|------|
| Markdown 渲染 | marked.js vs gomarkdown | gomarkdown | 服务端渲染避免 html/template 转义问题 |
| md 同步方式 | 行号映射 vs 内容匹配 | 内容匹配 | 行号会随编辑漂移，内容匹配更健壮 |
| 测试数据库 | Mock vs 真实 DB | 真实 DB | 原生 SQL 无 ORM，mock 成本高于真实 DB |
| 配色方案 | 赛博朋克(青紫绿) vs 暖色(金红蓝) | 暖色(金红蓝) | 更符合 Web3 金色主题，视觉对比度更好 |
| 学习任务页 | 手风琴 vs 树形表格 | 树形表格 | 扁平结构便于一次性浏览全部任务，展开收起满足层级需求 |
