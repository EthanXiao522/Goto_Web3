# Web3 Infra 学习追踪网站 - 需求文档

> 版本：v1.0 | 日期：2026-05-14 | 状态：待确认

---

## 1. 项目概述

### 1.1 项目背景

将以下两份 Markdown 文档转化为具有科技感的学习追踪 Web 应用：

| 文件 | 行数 | 内容 |
|------|------|------|
| `web3_infra_handbook.md` | 288 行 | Web3 基础设施工程师转型指南，含技术路线、系统设计、面试题、学习资源 |
| `web3_infra_3month_plan.md` | 975 行 | 3 阶段 12 周冲刺计划，含 84 天、约 232 个具体任务 |

### 1.2 目标用户

- 从传统后端（Java/Go）转型 Web3 基础设施方向的工程师
- 需要系统化追踪学习进度、记录学习成果

### 1.3 核心价值

- 将静态 Markdown 计划转化为可交互的任务追踪系统
- 每个用户独立进度，支持个性化学习路径
- 甘特图可视化整体时间线，进度一目了然
- 每个任务可记录学习链接、实现计划、代码和经验总结

---

## 2. 技术选型

| 层面 | 选型 | 版本/说明 |
|------|------|-----------|
| 后端语言 | Go | 1.21+ |
| Web 框架 | Gin | 路由、中间件、JSON |
| 前端渲染 | Go html/template | 服务端渲染，无 SPA 框架 |
| 前端交互 | 原生 JavaScript | 勾选框、手风琴、Tab 切换、甘特图 |
| 前端 Markdown | marked.js | CDN 引入，Markdown 渲染+预览 |
| 数据库 | MySQL | 8.0，字符集 utf8mb4 |
| 认证 | JWT | golang-jwt/jwt/v5，7 天过期，bcrypt 密码加密 |
| 甘特图 | 自定义 SVG/CSS 实现 | 周级别色块 + 完成率填充 |
| 部署 | Docker + Docker Compose | MySQL + Go 后端（含前端）二合一镜像 |

---

## 3. 功能需求

### 3.1 用户系统

**FR-1.1 注册**
- 输入：用户名、邮箱、密码（≥6 位）
- 密码使用 bcrypt 加密存储
- 用户名和邮箱全局唯一

**FR-1.2 登录**
- 输入：邮箱 + 密码
- 返回：JWT Token（7 天有效）+ 用户信息
- Token 存储在浏览器 localStorage

**FR-1.3 用户隔离**
- 每个用户有完全独立的任务完成状态和提交内容
- 用户不可见其他用户的数据

**FR-1.4 登录状态持久化**
- 页面刷新后自动恢复登录状态（从 localStorage 读取 Token）
- Token 过期自动跳转登录页

### 3.2 任务系统

**FR-2.1 任务层级（4 级树形结构）**
```
Phase (3 个)
└── Week (12 周，每 Phase 4 周)
    └── Day (84 天，每周 6-7 天)
        └── Task (约 232 个具体任务)
```

**FR-2.2 任务展示**
- Phase 详情页：四级手风琴结构，逐层展开/折叠
- 每层显示完成进度（如 "18/22 完成"）
- 任务项显示：复选框 + 内容描述 + 预估时间 + 提交类型指示点

**FR-2.3 任务完成**
- 点击复选框标记完成/取消完成
- 完成状态实时更新进度条
- 已完成任务文字显示删除线 + 降低透明度
- 完成后可进入详情页填写提交内容

**FR-2.4 进度统计**
- Phase 级别：进度环 + 百分比 + 完成数/总数
- Week 级别：进度条 + 百分比
- Day 级别：完成数/总任务数
- Dashboard 级别：总完成数、总任务数、百分比、预估剩余时间

**FR-2.5 检查点任务**
- Phase 末尾的自检清单与普通任务同等对待
- 可勾选完成、可提交内容
- 在 Phase 详情页底部以高亮卡片展示

### 3.3 任务提交内容

**FR-3.1 四种提交类型**

| 类型 | 格式 | 说明 |
|------|------|------|
| 学习链接 | Markdown | 参考文章、视频、文档等资源链接 |
| 落地计划 | Markdown | 任务的实现计划和设计思路 |
| 落地代码 | Markdown（含代码块） | 实际编写的代码 |
| 经验总结 | Markdown | 学习心得、踩坑记录 |

**FR-3.2 提交编辑器**
- 4 个 Tab 切换（学习链接/落地计划/落地代码/经验总结）
- 每个 Tab 内含编辑/预览切换
- 编辑模式：textarea
- 预览模式：marked.js 渲染 Markdown
- 自动保存提示（"已保存 · HH:MM:SS"）

**FR-3.3 提交内容指示**
- 任务列表项上显示彩色小圆点表示已提交的内容类型
  - 青色：学习链接
  - 紫色：落地计划
  - 绿色：落地代码
  - 黄色：经验总结

### 3.4 甘特图

**FR-4.1 周级别视图（默认）**
- X 轴：12 周（W1-W12）
- Y 轴：3 个 Phase 分组
- 每周一个色块，按 Phase 着色（青/紫/绿）
- 色块填充高度 = 该周任务完成率
- 悬停显示：周标题、任务数、完成率
- 点击色块跳转对应 Phase 详情页

**FR-4.2 天级别视图（切换）**
- 展开为 84 天的色块
- 按周分组，保持 Phase 颜色

**FR-4.3 甘特图附加信息**
- 图例（Phase 颜色说明）
- 关键里程碑卡片（Week 4/8/10/12）
- 12 周数据明细表（任务数、完成率、预计小时）

**FR-4.4 今天标记**
- 红色竖线标注"今天"位置

### 3.5 学习手册

**FR-5.1 手册渲染**
- 将 `web3_infra_handbook.md` 后端解析为 HTML
- 支持：标题、表格、代码块、引用块、链接、列表
- 响应式排版

**FR-5.2 侧边栏目录**
- 左侧导航显示手册章节链接
- 点击跳转到对应章节

### 3.6 Dashboard 仪表盘

**FR-6.1 统计卡片**
- 总任务数（232）
- 已完成数
- 完成率百分比
- 预估剩余时间

**FR-6.2 Phase 进度环**
- 3 个 SVG 圆环，颜色分别对应 Phase
- 动画填充到当前百分比
- 显示完成数/总任务数

**FR-6.3 每周进度条**
- 12 周的进度条列表
- 按 Phase 着色

**FR-6.4 最近活动**
- 最近 5 条活动记录（完成任务、提交内容等）
- 显示时间戳

### 3.7 页面导航

**FR-7.1 侧边栏**
- 默认展开（260px 宽），显示 Logo + 导航菜单 + 用户信息
- 右上角 ◀ 按钮：点击完全隐藏侧边栏
- 隐藏后：鼠标移到屏幕最左侧边缘，滑出 ☰ 展开按钮
- 移动端：侧边栏默认隐藏，顶部 ☰ 按钮切换

**FR-7.2 面包屑导航**
- Task 详情页显示：任务计划 > Phase N > Week N > Day N > 任务详情
- 每级可点击跳转

---

## 4. 非功能需求

### 4.1 性能
- 页面首次加载 < 2 秒
- API 响应 < 500ms（不含数据库查询）
- 粒子背景限制 100 个粒子，保留 60fps

### 4.2 安全
- 密码 bcrypt cost ≥ 12
- JWT Secret 从环境变量读取
- 所有 SQL 使用参数化查询（防注入）
- 日志中不输出密码、Token
- CORS 配置白名单

### 4.3 可用性
- 基本响应式适配（1024px / 640px 断点）
- 移动端可查看和勾选任务
- 甘特图等复杂视图优先桌面端
- 侧边栏可完全隐藏，最大化内容区域

### 4.4 兼容性
- Chrome / Firefox / Edge 最新版
- 不要求 IE 支持

### 4.5 数据完整性
- 任务导入幂等（重复运行不产生重复数据）
- user_tasks 懒加载创建（首次交互时生成记录）
- 外键约束保证数据一致性

---

## 5. 页面清单与路由

| # | 页面 | 路由 | 需要认证 | 说明 |
|---|------|------|----------|------|
| 1 | 首页 | `/` | 否 | 粒子星空背景 + 渐变标题 + 特性卡片 + 注册/登录入口 |
| 2 | 注册页 | `/register` | 否 | 用户名 + 邮箱 + 密码表单 |
| 3 | 登录页 | `/login` | 否 | 邮箱 + 密码表单，登录后跳转 Dashboard |
| 4 | Dashboard | `/dashboard` | 是 | 统计卡片 + Phase 进度环 + 每周进度条 + 最近活动 |
| 5 | Phase 列表 | `/phases` | 是 | 3 个 Phase 卡片 + 12 周进度表格 |
| 6 | Phase 详情 | `/phases/:id` | 是 | 四级手风琴展开 Phase→Week→Day→Task |
| 7 | Week 详情 | `/weeks/:id` | 是 | 聚焦某一周的完整任务列表 |
| 8 | Task 详情 | `/tasks/:id` | 是 | 任务信息 + 完成切换 + 4 个 Tab 提交编辑器 |
| 9 | 甘特图 | `/gantt` | 是 | 12 周甘特图 + 里程碑 + 明细表 |
| 10 | 学习手册 | `/handbook` | 是 | handbook.md 渲染页面 |

---

## 6. 数据库设计

### 6.1 ER 关系
```
users 1 ──── N user_tasks N ──── 1 tasks
                                    │
                                    N:1
                                    │
phases 1 ──── N weeks 1 ──── N days 1
```

### 6.2 表结构

#### users
| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGINT UNSIGNED | PK, AUTO_INCREMENT | |
| username | VARCHAR(64) | UNIQUE, NOT NULL | |
| email | VARCHAR(255) | UNIQUE, NOT NULL | |
| password_hash | VARCHAR(255) | NOT NULL | bcrypt |
| created_at | TIMESTAMP | DEFAULT NOW() | |
| updated_at | TIMESTAMP | ON UPDATE NOW() | |

#### phases
| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGINT UNSIGNED | PK | |
| phase_number | TINYINT | NOT NULL | 1/2/3 |
| title | VARCHAR(255) | NOT NULL | "Phase 1：基础夯实" |
| subtitle | VARCHAR(255) | | "第 1-4 周" |
| goal | TEXT | | 阶段目标 |
| deliverables | TEXT | | 阶段产出 |
| sort_order | INT | DEFAULT 0 | |

#### weeks
| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGINT UNSIGNED | PK | |
| phase_id | BIGINT UNSIGNED | FK→phases, NOT NULL | |
| week_number | TINYINT | UNIQUE, NOT NULL | 1-12 |
| title | VARCHAR(255) | NOT NULL | |
| subtitle | VARCHAR(255) | | |
| goal | TEXT | | |
| deliverables | TEXT | | |
| sort_order | INT | DEFAULT 0 | |

#### days
| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGINT UNSIGNED | PK | |
| week_id | BIGINT UNSIGNED | FK→weeks, NOT NULL | |
| day_number | TINYINT | NOT NULL | 1-7 |
| title | VARCHAR(255) | NOT NULL | |
| sort_order | INT | DEFAULT 0 | |

#### tasks
| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGINT UNSIGNED | PK | |
| day_id | BIGINT UNSIGNED | FK→days, NOT NULL | |
| content | TEXT | NOT NULL | 任务描述 |
| estimated_hours | DECIMAL(4,1) | DEFAULT 0 | |
| resource_urls | JSON | | 内置 URL 列表 |
| sort_order | INT | DEFAULT 0 | |
| is_checkpoint | TINYINT(1) | DEFAULT 0 | |

#### user_tasks
| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGINT UNSIGNED | PK | |
| user_id | BIGINT UNSIGNED | FK→users, NOT NULL | |
| task_id | BIGINT UNSIGNED | FK→tasks, NOT NULL | |
| is_completed | TINYINT(1) | DEFAULT 0 | |
| completed_at | TIMESTAMP | NULL | 点击完成时设置 |
| learning_links | TEXT | | Markdown |
| implementation_plan | TEXT | | Markdown |
| implementation_code | TEXT | | Markdown |
| experience_summary | TEXT | | Markdown |
| created_at | TIMESTAMP | DEFAULT NOW() | |
| updated_at | TIMESTAMP | ON UPDATE NOW() | |
| UNIQUE(user_id, task_id) | | | 懒加载创建 |

---

## 7. API 接口规格

### 7.1 响应格式

成功：`{ "code": 200, "data": {...} }`
错误：`{ "code": 4xx/5xx, "message": "error description" }`

### 7.2 接口列表

#### 认证接口（公开）
```
POST /api/v1/auth/register
  Request:  { "username": "str", "email": "str", "password": "str" }
  Response: { "code": 200, "data": { "token": "jwt...", "user": {...} } }

POST /api/v1/auth/login
  Request:  { "email": "str", "password": "str" }
  Response: { "code": 200, "data": { "token": "jwt...", "user": {...} } }
```

#### 认证接口（需 JWT）
```
GET /api/v1/auth/me
  Response: { "code": 200, "data": { "user": {...} } }
```

#### Phase 接口（需 JWT）
```
GET /api/v1/phases
  Response: { "code": 200, "data": { "phases": [{id, phase_number, title, subtitle, week_count, task_count, completed_count}] } }

GET /api/v1/phases/:id
  Response: { "code": 200, "data": { "phase": {...}, "weeks": [{id, week_number, title, day_count, task_count, completed_count}] } }
```

#### Week 接口（需 JWT）
```
GET /api/v1/weeks/:id
  Response: { "code": 200, "data": { "week": {...}, "phase": {...}, "days": [{id, day_number, title, task_count, completed_count, tasks: [{id, content, estimated_hours, is_completed, has_links, has_plan, has_code, has_summary}]}] } }
```

#### Task 接口（需 JWT）
```
GET /api/v1/tasks/:id
  Response: { "code": 200, "data": { "task": {...}, "day": {...}, "week": {...}, "phase": {...}, "user_task": {...} } }

PATCH /api/v1/tasks/:id/complete
  Request:  { "completed": true|false }
  Response: { "code": 200, "data": { "user_task": {...} } }

PUT /api/v1/tasks/:id/submissions
  Request:  { "learning_links"?: "str", "implementation_plan"?: "str", "implementation_code"?: "str", "experience_summary"?: "str" }
  Response: { "code": 200, "data": { "user_task": {...} } }
```

#### 进度接口（需 JWT）
```
GET /api/v1/progress/overview
  Response: { "code": 200, "data": { "total_tasks": 232, "completed_tasks": N, "percentage": X.X, "phases": [{id, completed, total, percentage}], "weeks": [{id, completed, total, percentage}] } }

GET /api/v1/dashboard
  Response: { "code": 200, "data": { "overview": {...}, "recent_completions": [...], "recent_submissions": [...] } }
```

### 7.3 错误码
| 状态码 | 含义 |
|--------|------|
| 400 | 请求参数错误 |
| 401 | Token 无效或过期 |
| 404 | 资源不存在 |
| 409 | 数据冲突（如重复注册） |
| 500 | 服务器内部错误 |

---

## 8. UI 设计规范

### 8.1 配色方案
| 用途 | 色值 | 说明 |
|------|------|------|
| 背景 | `#0a0e1a` | 深蓝黑 |
| 表面 | `#111827` | 侧边栏、卡片 |
| 卡片 | `#1a1f35` | 卡片背景 |
| 边框 | `#2a3558` | 普通边框 |
| 边框发光 | `#2a4068` | 悬停边框 |
| 青色 | `#00f0ff` | Phase 1、主强调色 |
| 紫色 | `#7c3aed` | Phase 2 |
| 绿色 | `#00ff88` | Phase 3、成功 |
| 黄色 | `#ffb700` | 警告、检查点 |
| 红色 | `#ff3860` | 错误、今天标记 |
| 主文字 | `#e2e8f0` | |
| 次要文字 | `#94a3b8` | |
| 禁用文字 | `#64748b` | |

### 8.2 字体
| 用途 | 字体 | 备选 |
|------|------|------|
| 标题 | Orbitron | sans-serif |
| 正文 | Inter | system-ui, sans-serif |
| 代码 | JetBrains Mono | Fira Code, monospace |

### 8.3 视觉效果
- Canvas 粒子星空背景（100 粒子，连线距离 120px 内）
- CSS 网格叠加层（60px 间距）
- 卡片悬停：边框发光 + 轻微上移 (translateY -4px)
- 复选框勾选：青色填充 + 缩放动画 + 阴影发光
- SVG 进度环：stroke-dasharray 动画
- 侧边栏：完全隐藏/展开，过渡动画 0.3s

### 8.4 响应式断点
```
> 1024px: 正常布局（侧边栏 260px + 内容区）
640-1024px: 侧边栏默认隐藏，顶部显示 ☰ 按钮
< 640px: 统计卡片单列，内容 padding 缩小
```

---

## 9. 数据导入

### 9.1 解析规则
- 状态机逐行解析 `web3_infra_3month_plan.md`
- 检测 `## Phase` / `### 第N周` / `#### Day` 标题
- 检测 `- [ ]` 行作为任务
- 提取任务中的 `[text](url)` 链接存入 resource_urls
- 提取 `（Nh）` 格式的小时数存入 estimated_hours
- 检测 `**自检清单：**` 标记后续任务为 is_checkpoint=1

### 9.2 执行方式
```bash
cd backend && go run cmd/seed/main.go --plan ../web3_infra_3month_plan.md
```
- 幂等：INSERT IGNORE 基于唯一约束去重
- 数据量：3 Phase + 12 Week + 84 Day + ~232 Task

---

## 10. 验收标准

| # | 功能 | 验收标准 |
|---|------|----------|
| 1 | 注册 | 注册成功返回 Token，重复用户名/邮箱返回 409 |
| 2 | 登录 | 正确密码返回 Token，错误密码返回 401 |
| 3 | Phase 列表 | 正确显示 3 个 Phase 及其进度统计 |
| 4 | Phase 详情 | 四级手风琴展开，勾选任务更新进度 |
| 5 | Week 详情 | 显示该周全部 Day 和 Task，正确显示完成状态 |
| 6 | Task 详情 | 显示任务信息、参考链接、4 个 Tab 编辑器 |
| 7 | 任务完成 | 勾选后进度实时更新，取消勾选回退 |
| 8 | 提交内容 | 4 种类型独立保存，Markdown 预览正常 |
| 9 | 甘特图 | 12 周色块 + 完成率填充 + 里程碑卡片 |
| 10 | Dashboard | 统计卡片 + 进度环 + 活动列表数据正确 |
| 11 | 学习手册 | handbook 内容正确渲染为 HTML |
| 12 | 用户隔离 | 用户 A 的完成状态不影响用户 B |
| 13 | 侧边栏 | 隐藏/展开动画流畅，左侧边缘滑出展开 |
| 14 | Docker | docker-compose up 一键启动全站 |

---

## 11. 项目结构

```
Goto_Web3/
├── backend/
│   ├── cmd/server/main.go
│   ├── cmd/seed/main.go
│   ├── internal/
│   │   ├── config/config.go
│   │   ├── database/{mysql,migrate}.go
│   │   ├── middleware/{auth,cors}.go
│   │   ├── model/{user,phase,week,day,task,user_task}.go
│   │   ├── repository/{user,phase,week,day,task,user_task}_repo.go
│   │   ├── service/{auth,task,progress,dashboard}_service.go
│   │   ├── handler/{auth,phase,week,task,progress,dashboard}_handler.go
│   │   ├── router/router.go
│   │   └── importer/{parser,seeder}.go
│   ├── templates/  (10 个 .html)
│   ├── static/{css,js,lib}/
│   ├── go.mod
│   └── Dockerfile
├── docker-compose.yml
├── REQUIREMENTS.md
├── DESIGN.md
├── DEPLOYMENT.md
├── TEST_REPORT.md
├── SUMMARY.md
└── README.md
```
