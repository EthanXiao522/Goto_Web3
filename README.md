# Goto Web3 — Web3 Infra 学习追踪平台

> 版本：v1.1 | 日期：2026-05-15 | 测试：88/88 通过

---

## 目录

- [1. 系统概述](#1-系统概述)
- [2. 技术框架](#2-技术框架)
- [3. 部署说明](#3-部署说明)
- [4. 功能模块](#4-功能模块)
  - [4.1 用户系统](#41-用户系统)
  - [4.2 Dashboard 仪表盘](#42-dashboard-仪表盘)
  - [4.3 学习任务](#43-学习任务)
  - [4.4 任务详情与提交](#44-任务详情与提交)
  - [4.5 甘特图](#45-甘特图)
  - [4.6 学习计划书](#46-学习计划书)
  - [4.7 首页与游客模式](#47-首页与游客模式)
- [5. 代码说明](#5-代码说明)
  - [5.1 项目结构](#51-项目结构)
  - [5.2 后端架构](#52-后端架构)
  - [5.3 前端架构](#53-前端架构)
  - [5.4 数据库设计](#54-数据库设计)
  - [5.5 API 接口](#55-api-接口)
  - [5.6 测试体系](#56-测试体系)
- [6. 开发流程](#6-开发流程)

---

## 1. 系统概述

Goto Web3 是一个面向 Web3 基础设施工程师的 **学习任务追踪平台**。系统基于一份 3 个月（12 周）的详细学习计划，将学习内容拆解为 **3 阶段 × 12 周 × 84 天 × 232 个具体任务**，提供进度追踪、任务管理、甘特图可视化和学习手册渲染等功能。

**核心数据流：** `web3_infra_3month_plan.md` → Seed 导入 → MySQL → Gin API → Go html/template SSR → 浏览器

**核心特性：**
- 四级任务层级（阶段→周→天→任务），手风琴展开/折叠
- 232 个可追踪任务，含完成勾选、4 种提交内容编辑
- 12 周甘特图时间线，按阶段着色，周/日粒度切换
- 树形任务总览页，支持内联编辑并自动同步到源 Markdown 文档
- 服务端 Markdown 渲染，学习计划书实时呈现
- JWT 认证 + 用户数据隔离 + bcrypt 密码加密
- 赛博朋克主题 UI，响应式适配

---

## 2. 技术框架

| 层级 | 技术 | 版本 |
|------|------|------|
| 后端框架 | Gin | v1.9 |
| 语言 | Go | 1.21+ |
| 模板引擎 | Go html/template | 标准库 |
| 数据库 | MySQL (InnoDB, utf8mb4) | 8.0 |
| 数据库驱动 | go-sql-driver/mysql | 最新 |
| 认证 | golang-jwt + bcrypt | v5 |
| Markdown 渲染 | gomarkdown/markdown | 最新 |
| 前端交互 | Vanilla JavaScript | ES6+ |
| CSS | 自定义属性 + 赛博朋克主题 | ~2400 行 |
| 部署 | Docker + Docker Compose | 多阶段构建 |

**无 ORM、无 SPA 框架、无 npm 依赖**。所有 SQL 原生参数化查询，所有前端 JS 原生实现。

---

## 3. 部署说明

### 环境要求

- Go 1.21+
- MySQL 8.0
- Docker & Docker Compose (可选)

### 快速启动

```bash
# 1. 克隆项目
cd Goto_Web3

# 2. 配置环境变量
export DB_DSN="username:password@tcp(127.0.0.1:3306)/goto_web3?parseTime=true&charset=utf8mb4"
export JWT_SECRET="your-secret-key"
export PORT=8080

# 3. 导入学习计划数据
cd backend && go run cmd/seed/main.go ../sources/web3_infra_3month_plan.md

# 4. 启动服务
go run cmd/server/main.go
# 访问 http://localhost:8080
```

### Docker 部署

```bash
docker-compose up -d
# MySQL 8.0 + Go App 一键启动，访问 http://localhost:8080
```

### 架构说明：
 本项目是 Go SSR（服务端渲染），没有独立的前端服务。go run cmd/server/main.go 一个命令同时启动：

- API 服务（/api/v1/*）
- 页面渲染（/, /dashboard, /tasks, /handbook 等）
- 静态资源（/static/css/style.css, /static/js/app.js）

#### 日常开发命令：
```
# 启动:根目录执行
cd backend && go run cmd/server/main.go
cd ../

# 导入/更新种子数据:根目录执行
cd backend && go run cmd/seed/main.go ../sources/web3_infra_3month_plan.md
cd ../

# 运行测试:根目录执行
cd backend && go test -v ./test/
cd ../
```

### 目录结构

```
├── backend/
│   ├── cmd/server/main.go        # HTTP 服务入口
│   ├── cmd/seed/main.go          # 数据导入入口
│   └── internal/
│       ├── config/               # 环境变量
│       ├── database/             # 连接池 + DDL 迁移
│       ├── middleware/            # JWT / CORS / Logger
│       ├── model/                # 6 个数据结构
│       ├── repository/           # 原生 SQL 查询 (6 repo)
│       ├── service/              # 业务逻辑 (4 service)
│       ├── handler/              # HTTP 处理 (9 handler)
│       ├── router/               # 路由注册
│       └── importer/             # Markdown 解析 + 导入
├── frontend/
│   ├── templates/ (12 .html)     # SSR 模板
│   └── static/css/               # 样式 (~2400 行)
├── sources/                      # 源 Markdown 文件
├── REQUIREMENTS.md               # 需求文档
├── DESIGN.md                     # 设计文档
└── TEST_REPORT.md                # 测试报告
```

---

## 4. 功能模块

### 4.1 用户系统

- **注册/登录**：用户名 + 邮箱 + 密码，bcrypt cost≥12 加密存储
- **JWT 认证**：7 天有效期，Cookie 传递，中间件拦截
- **用户隔离**：每个用户独立任务完成状态和提交内容
- **登录跳转**：已登录用户访问 `/login` 或 `/register` 自动重定向到 Dashboard
- **个人信息**：支持修改用户名和邮箱

### 4.2 Dashboard 仪表盘

核心数据看板，登录后首页，包含 5 个区块：

| 区块 | 内容 |
|------|------|
| Dashboard进度卡 | 4 卡片单行展示（总任务数/已完成/完成阶段/总体进度），独立 teal/emerald/violet/amber 配色 |
| 阶段进度 | 3 个 SVG 圆环，金/红/蓝对应 Phase 1/2/3，显示百分比和任务数 |
| 每周进度 | 12 周柱状图，柱高=完成率，按 Phase 着色，顶部显示完成数和百分比 |
| 甘特图 | 12 周网格线，bar-fill 填充 + 百分比文字，中文阶段标签 |
| 最近活动 | 最近任务活动列表，支持分页展开（5条→更多） |

### 4.3 学习任务

**路由：** `/tasks`（需登录）  
**模板：** [learning_tasks.html](frontend/templates/learning_tasks.html)

- 树形可展开表格：阶段 → 周 → 天 → 任务，四层层级
- 每行点击箭头展开/收起子层级，子行自动缩进
- 显示完成状态（✓/○）、预估学时、资源链接标记、检查点高亮
- 内联编辑：点击「✎」弹出编辑框，支持 Enter 保存 / Escape 取消
- 编辑后自动同步到 `sources/web3_infra_3month_plan.md` 源文档
- 数据量：3 Phase × 12 Week × 84 Day × 233 Task

### 4.4 任务详情与提交

**路由：** `/tasks/:id`（需登录）

- **4 种提交类型：** 学习链接、落地计划、落地代码、经验总结
- **编辑器：** Tab 切换 + Markdown 编辑/预览切换（marked.js）
- **完成勾选：** 复选框一键标记完成/取消
- **提交指示：** 任务列表显示彩色圆点标记提交类型

### 4.5 甘特图

**路由：** `/gantt`（需登录）

- 12 周时间线，按 Phase 着色（金/红/蓝）
- 周/日粒度切换，色块填充高度=完成率
- 图例、里程碑卡片、12 周数据明细表
- 首页（Landing/Demo/Dashboard）甘特图统一风格：bar-fill + bar-text + 网格线

### 4.6 学习计划书

**路由：** `/handbook`（需登录）  
**渲染引擎：** gomarkdown（服务端）

- 读取 `sources/web3_infra_3month_plan.md`，服务端渲染为 HTML
- 顶部「📄 源文档」链接，点击在新标签页查看完整渲染
- 自动过滤 `[TOC]` 占位符
- Blockquote 内列表正确渲染为 `<ul>` 结构
- 源文档页 `/handbook/source` 单独展示完整渲染

### 4.7 首页与游客模式

- **首页 `/`**：粒子背景 + 渐变标题 + 特性卡片 + 功能预览区（Dashboard/阶段进度/每周进度/甘特图）
- **阶段进度卡片**：首页预览区不可点击跳转
- **游客模式 `/demo`**：示例数据预览完整功能，引导注册

---

## 5. 代码说明

### 5.1 项目结构

```
Goto_Web3/
├── backend/                        # Go 后端 (37 Go files, ~3869 lines)
│   ├── cmd/
│   │   ├── server/main.go          # HTTP 服务入口
│   │   └── seed/main.go            # 数据导入命令
│   ├── internal/
│   │   ├── config/config.go        # PORT / DB_DSN / JWT_SECRET
│   │   ├── database/
│   │   │   ├── mysql.go            # sql.Open + Ping, MaxOpen=25
│   │   │   └── migrate.go          # 6 表 CREATE TABLE IF NOT EXISTS
│   │   ├── middleware/
│   │   │   ├── auth.go             # JWT 解析 + user_id 注入 ctx
│   │   │   ├── cors.go             # Access-Control-Allow-Origin: *
│   │   │   └── logger.go           # slog 请求日志
│   │   ├── model/                  # User/Phase/Week/Day/Task/UserTask
│   │   ├── repository/             # UserRepo/PhaseRepo/WeekRepo/DayRepo/TaskRepo/UserTaskRepo
│   │   ├── service/                # AuthService/TaskService/ProgressService
│   │   ├── handler/
│   │   │   ├── auth.go             # Register/Login/Me/UpdateProfile/CheckUsername
│   │   │   ├── phase.go            # GetPhases/GetPhaseDetail/GetWeekDetail
│   │   │   ├── task.go             # GetTaskDetail/ToggleComplete/UpdateSubmissions/UpdateContent
│   │   │   ├── progress.go         # GetDashboard/GetOverview
│   │   │   └── page.go             # Landing/Dashboard/PhaseDetail/WeekDetail/TaskDetail/
│   │   │                           #   LearningTasks/Handbook/HandbookSource/Profile/Demo
│   │   ├── router/router.go        # 27 条路由注册
│   │   └── importer/               # MD 状态机解析 + INSERT IGNORE 导入
│   └── test/                       # 集成测试 (5 files, 40 tests)
├── frontend/
│   ├── templates/ (12 .html)       # SSR 模板
│   │   ├── landing.html            # 首页（功能预览）
│   │   ├── dashboard.html          # Dashboard 仪表盘
│   │   ├── phases.html             # Phase 列表
│   │   ├── phase_detail.html       # Phase 详情（四级手风琴）
│   │   ├── week_detail.html        # Week 详情
│   │   ├── task_detail.html        # Task 详情（4 Tab 编辑器）
│   │   ├── gantt.html              # 甘特图页面
│   │   ├── learning_tasks.html     # 学习任务树形表格（NEW）
│   │   ├── handbook.html           # 学习计划书（gomarkdown 渲染）
│   │   ├── handbook_source.html    # 源文档渲染页（NEW）
│   │   ├── auth_login.html         # 登录
│   │   ├── auth_register.html      # 注册
│   │   ├── demo.html               # 游客模式
│   │   ├── profile.html            # 个人信息
│   │   ├── _sidebar.html           # 侧边栏组件
│   │   ├── _topnav.html            # 顶栏组件
│   │   ├── app-layout.html         # 主布局（含侧边栏）
│   │   └── layout.html             # 基础布局
│   └── static/
│       ├── css/style.css           # 赛博朋克主题 (~2400 行)
│       ├── js/app.js               # 认证/粒子/Toast/侧边栏
│       ├── js/phase.js             # 手风琴/勾选/进度更新
│       ├── js/task.js              # Tab/预览/自动保存
│       ├── js/dashboard.js         # 进度环动画
│       └── js/gantt.js             # 甘特图渲染
├── sources/
│   └── web3_infra_3month_plan.md   # 学习计划源文档
├── docker-compose.yml              # MySQL 8.0 + Go App
├── Dockerfile                      # 多阶段构建 (golang→alpine)
├── REQUIREMENTS.md                 # 需求文档 v1.1
├── DESIGN.md                       # 设计文档 v1.1
└── TEST_REPORT.md                  # 测试报告 v1.1
```

### 5.2 后端架构

**三层架构：** Handler → Service → Repository

```
Browser Request → Gin Router
    ├── Page Route (SSR)  → Handler → Service → Repo → DB
    │                       └── c.HTML(200, template, data)
    └── API Route (JSON)  → CORS → Logger → Auth(JWT)
                            └── Handler → Service → Repo → DB
                                └── c.JSON(200, response)
```

**关键设计原则：**
- **无 ORM**：全部使用原生 SQL + 参数化查询 + 手动 Scan
- **懒加载**：user_tasks 表使用 INSERT IGNORE，首次操作时自动创建
- **用户隔离**：所有查询带 user_id 条件，JWT 中间件注入
- **幂等导入**：Seed 命令可重复执行不报错

### 5.3 前端架构

- **SSR 渲染**：Go html/template，无 SPA 框架
- **模板组合**：`app-layout` + `sidebar` + `topnav` 组件化复用
- **CSS 变量**：`--bg-primary`, `--accent-cyan`, `--border-color` 等 20+ 变量
- **配色方案**：Phase 1/2/3 对应 金(#ffd700)/红(#ff4466)/蓝(#4488ff)
- **响应式断点**：1024px（侧边栏隐藏）、640px（单列）
- **无 npm**：前端 JS 零依赖，marked.js 通过 CDN 按需加载

### 5.4 数据库设计

**6 表 ER 关系：**

```
users 1──N user_tasks N──1 tasks
                                    │
                                  N:1
                                    │
phases 1──N weeks 1──N days 1──────┘
```

| 表 | 记录数 | 核心字段 |
|----|--------|----------|
| users | 按注册量 | username, email, password_hash |
| phases | 3 | phase_number(1-3), title, goal, deliverables |
| weeks | 12 | phase_id(FK), week_number(1-12), title, goal |
| days | 84 | week_id(FK), day_number(1-7), title |
| tasks | 232 | day_id(FK), content, estimated_hours, resource_urls, is_checkpoint |
| user_tasks | 懒加载 | user_id(FK)+task_id(FK) UNIQUE, is_completed, 4 种提交字段 |

### 5.5 API 接口

**公开接口（无需认证）：**

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/auth/register` | 用户注册 |
| POST | `/api/v1/auth/login` | 用户登录 |
| GET | `/api/v1/auth/check-username` | 检查用户名可用性 |

**需认证接口（JWT Cookie）：**

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/auth/me` | 当前用户信息 |
| PUT | `/api/v1/auth/profile` | 修改个人信息 |
| GET | `/api/v1/phases` | Phase 列表（含进度） |
| GET | `/api/v1/phases/:id` | Phase 详情（含 Weeks） |
| GET | `/api/v1/weeks/:id` | Week 详情（含 Days+Tasks） |
| GET | `/api/v1/tasks/:id` | Task 详情（含 UserTask） |
| PATCH | `/api/v1/tasks/:id/complete` | 切换任务完成状态 |
| PUT | `/api/v1/tasks/:id/content` | 更新任务内容（自动同步 md） |
| PUT | `/api/v1/tasks/:id/submissions` | 更新 4 种提交内容 |
| GET | `/api/v1/dashboard` | Dashboard 聚合数据 |
| GET | `/api/v1/progress` | 总体进度统计 |

**页面路由（需认证）：** `/dashboard`, `/tasks`, `/phases/:id`, `/weeks/:id`, `/tasks/:id`, `/handbook`, `/handbook/source`, `/profile`

### 5.6 测试体系

| 类型 | 文件 | 测试数 | 覆盖范围 |
|------|------|--------|----------|
| Repository 单元测试 | repository_test.go | 14 | 6 repo 全部 CRUD + UpdateContent + UpdateFields + Batch |
| Handler 集成测试 | handler_test.go | 14 | Auth(5) + Task(5) + Phase(4) |
| Service 单元测试 | service_test.go | 4 | ProgressService(2) + TaskService(2) |
| Middleware 单元测试 | middleware_test.go | 5 | Auth(4) + CORS(1) |
| Importer 单元测试 | importer_test.go | 2 | 真实文件 + Fixture |
| API 功能测试 | curl | 38 | 页面/API/鉴权/布局 |
| **合计** | **5 files** | **88** | **全部通过** |

```bash
cd backend && go test -v ./test/   # 40 unit tests
```

---

## 6. 开发流程

本项目遵循标准化 7 步开发流程：

| 步骤 | 文档 | 状态 |
|------|------|------|
| 1. 需求分析 | [REQUIREMENTS.md](REQUIREMENTS.md) v1.1 | ✅ 已确认 |
| 2. 系统设计 | [DESIGN.md](DESIGN.md) v1.1 | ✅ 已确认 |
| 3. 编码实现 | 每个功能 git commit | ✅ 已完成 |
| 4. 系统测试 | [TEST_REPORT.md](TEST_REPORT.md) v1.1 | ✅ 已确认 |
| 5. README | README.md v1.1 | ✅ 本文档 |
| 6. 总结 | [SUMMARY.md](SUMMARY.md) | → 下一步 |
| 7. 权限确认 | 需求/设计/测试 三确认点 | → 待确认 |
