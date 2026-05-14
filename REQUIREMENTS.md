# Web3 Infra 学习追踪网站 - 需求文档

> 版本：v1.0 | 日期：2026-05-14 | 状态：待确认

---

## 目录

- [1. 项目概述](#1-项目概述)
  - [1.1 项目背景](#11-项目背景)
  - [1.2 目标用户](#12-目标用户)
  - [1.3 核心价值](#13-核心价值)
- [2. 技术选型](#2-技术选型)
- [3. 功能需求](#3-功能需求)
  - [3.1 用户系统](#31-用户系统)
  - [3.2 任务系统](#32-任务系统)
  - [3.3 任务提交内容](#33-任务提交内容)
  - [3.4 甘特图](#34-甘特图)
  - [3.5 学习手册](#35-学习手册)
  - [3.6 Dashboard 仪表盘](#36-dashboard-仪表盘)
  - [3.7 页面导航](#37-页面导航)
- [4. 非功能需求](#4-非功能需求)
  - [4.1 性能](#41-性能)
  - [4.2 安全](#42-安全)
  - [4.3 可用性](#43-可用性)
  - [4.4 兼容性](#44-兼容性)
  - [4.5 数据完整性](#45-数据完整性)
- [5. 页面清单与路由](#5-页面清单与路由)
- [6. 数据库设计](#6-数据库设计)
  - [6.1 ER 关系](#61-er-关系)
  - [6.2 表结构](#62-表结构)
- [7. API 接口规格](#7-api-接口规格)
  - [7.1 响应格式](#71-响应格式)
  - [7.2 接口列表](#72-接口列表)
  - [7.3 错误码](#73-错误码)
- [8. UI 设计规范](#8-ui-设计规范)
  - [8.1 配色方案](#81-配色方案)
  - [8.2 字体](#82-字体)
  - [8.3 视觉效果](#83-视觉效果)
  - [8.4 响应式断点](#84-响应式断点)
- [9. 数据导入](#9-数据导入)
  - [9.1 解析规则](#91-解析规则)
  - [9.2 执行方式](#92-执行方式)
- [10. 验收标准](#10-验收标准)
- [11. 项目结构](#11-项目结构)

---

## 1. 项目概述

### 1.1 项目背景

将两份 Markdown 文档转化为具有科技感的学习追踪 Web 应用：

| 文件 | 行数 | 内容 |
|------|------|------|
| `web3_infra_handbook.md` | 288 行 | Web3 基础设施工程师转型指南，含技术路线、系统设计、面试题 |
| `web3_infra_3month_plan.md` | 975 行 | 3 阶段 12 周冲刺计划，含约 232 个具体任务 |

### 1.2 目标用户

从传统后端（Java/Go）转型 Web3 基础设施方向的工程师，需要系统化追踪学习进度、记录学习成果。

### 1.3 核心价值

- 将静态 Markdown 计划转化为可交互的任务追踪系统
- 每个用户独立进度，支持个性化学习路径
- 甘特图可视化整体时间线，进度一目了然
- 每个任务可记录学习链接、实现计划、代码和经验总结

---

## 2. 技术选型

| 层面 | 选型 | 说明 |
|------|------|------|
| 后端语言 | Go 1.21+ | 主语言 |
| Web 框架 | Gin v1.9 | 路由、中间件、JSON 响应 |
| 前端渲染 | Go html/template | 服务端渲染 HTML，无 SPA 框架 |
| 前端交互 | 原生 JavaScript ES6+ | 勾选框、手风琴、Tab 切换、甘特图 |
| 前端 Markdown | marked.js (CDN) | Markdown 渲染和预览 |
| 数据库 | MySQL 8.0 | utf8mb4 字符集 |
| 认证 | JWT (golang-jwt v5) | 7 天过期，bcrypt 密码加密 |
| 容器化 | Docker + Docker Compose | MySQL + Go 后端（含前端）二合一镜像 |
| 部署 | 多阶段构建 (golang→alpine) | 一个镜像包含 API + 全部前端页面 |

---

## 3. 功能需求

### 3.1 用户系统

| 编号 | 功能 | 说明 |
|------|------|------|
| FR-1.1 | 注册 | 用户名 + 邮箱 + 密码（≥6 位），bcrypt 加密，用户名和邮箱全局唯一 |
| FR-1.2 | 登录 | 邮箱 + 密码，返回 JWT Token（7 天有效）+ 用户信息 |
| FR-1.3 | 用户隔离 | 每个用户独立任务完成状态和提交内容，不可见其他用户数据 |
| FR-1.4 | 登录持久化 | 页面刷新后从 localStorage 恢复 Token，过期自动跳转登录页 |

### 3.2 任务系统

| 编号 | 功能 | 说明 |
|------|------|------|
| FR-2.1 | 四级层级 | Phase(3) → Week(12) → Day(84) → Task(~232) |
| FR-2.2 | 手风琴展示 | 逐层展开/折叠，每层显示完成进度 |
| FR-2.3 | 任务完成 | 点击复选框标记完成/取消，完成项显示删除线和降低透明度 |
| FR-2.4 | 进度统计 | Phase 级别: 进度环 + 百分比；Week 级别: 进度条；Day 级别: 完成数 |
| FR-2.5 | 检查点任务 | Phase 末尾自检清单与普通任务同等对待，可勾选、可提交内容 |

### 3.3 任务提交内容

| 编号 | 功能 | 格式 | 说明 |
|------|------|------|------|
| FR-3.1 | 学习链接 | Markdown | 参考文章、视频、文档等资源 |
| FR-3.2 | 落地计划 | Markdown | 实现计划和设计思路 |
| FR-3.3 | 落地代码 | Markdown | 实际编写的代码（含代码块） |
| FR-3.4 | 经验总结 | Markdown | 学习心得、踩坑记录 |
| FR-3.5 | 编辑器 | textarea | 4 Tab 切换 + 编辑/预览切换 |
| FR-3.6 | 状态指示 | 彩色圆点 | 任务列表显示提交类型指示点（青/紫/绿/黄） |

### 3.4 甘特图

| 编号 | 功能 | 说明 |
|------|------|------|
| FR-4.1 | 周级视图 | 12 周甘特图，按 Phase 着色（青/紫/绿），色块填充高度=完成率 |
| FR-4.2 | 天级视图 | 展开为 84 天色块，按周分组 |
| FR-4.3 | 交互 | 悬停显示详情，点击跳转对应 Week |
| FR-4.4 | 附加信息 | 图例、关键里程碑卡片（Week 4/8/10/12）、12 周数据明细表 |
| FR-4.5 | 今日标记 | 红色竖线标注当前位置 |

### 3.5 学习手册

| 编号 | 功能 | 说明 |
|------|------|------|
| FR-5.1 | 内容渲染 | 服务端解析 handbook.md → HTML |
| FR-5.2 | 侧边栏目录 | 手册章节导航，点击跳转 |

### 3.6 Dashboard 仪表盘

| 编号 | 功能 | 说明 |
|------|------|------|
| FR-6.1 | 统计卡片 | 总任务数 / 已完成 / 完成率 / 预估剩余时间 |
| FR-6.2 | Phase 进度环 | 3 个 SVG 圆环，颜色对应 Phase，动画填充 |
| FR-6.3 | 每周进度条 | 12 周进度条列表，按 Phase 着色 |
| FR-6.4 | 最近活动 | 最近 5 条活动记录，含时间戳 |

### 3.7 页面导航

| 编号 | 功能 | 说明 |
|------|------|------|
| FR-7.1 | 侧边栏 | 默认展开 260px，◀ 按钮完全隐藏，左侧边缘滑出 ☰ 展开按钮 |
| FR-7.2 | 面包屑 | Task 详情页显示完整路径，每级可点击 |

---

## 4. 非功能需求

### 4.1 性能

- 页面首次加载 < 2 秒
- API 响应 < 500ms
- 粒子背景限制 100 个粒子，保持 60fps

### 4.2 安全

- 密码 bcrypt cost ≥ 12
- JWT Secret 从环境变量读取
- 所有 SQL 参数化查询
- 日志中不输出密码、Token
- CORS 配置白名单

### 4.3 可用性

- 响应式适配（1024px / 640px 断点）
- 移动端可查看和勾选任务
- 甘特图等复杂视图优先桌面端
- 侧边栏可完全隐藏

### 4.4 兼容性

Chrome / Firefox / Edge 最新版，不要求 IE 支持。

### 4.5 数据完整性

- 任务导入幂等（INSERT IGNORE）
- user_tasks 懒加载创建
- 外键约束保证一致性

---

## 5. 页面清单与路由

| # | 页面 | 路由 | 认证 | 核心功能 |
|---|------|------|------|----------|
| 1 | 首页 | `/` | 否 | 粒子背景 + 渐变标题 + 特性卡片 |
| 2 | 登录 | `/login` | 否 | 邮箱 + 密码表单 |
| 3 | 注册 | `/register` | 否 | 用户名 + 邮箱 + 密码表单 |
| 4 | Dashboard | `/dashboard` | 是 | 统计卡片 + 进度环 + 活动列表 |
| 5 | Phase 列表 | `/phases` | 是 | Phase 卡片 + 周进度表格 |
| 6 | Phase 详情 | `/phases/:id` | 是 | 四级手风琴展开（核心交互页） |
| 7 | Week 详情 | `/weeks/:id` | 是 | 聚焦单周完整任务列表 |
| 8 | Task 详情 | `/tasks/:id` | 是 | 4 Tab 提交编辑器 |
| 9 | 甘特图 | `/gantt` | 是 | 12 周甘特图 + 里程碑 |
| 10 | 学习手册 | `/handbook` | 是 | handbook.md 渲染 |

---

## 6. 数据库设计

### 6.1 ER 关系

```
    users 1───N user_tasks N───1 tasks
                                     │
                                   N:1
                                     │
    phases 1───N weeks 1───N days 1──┘
```

### 6.2 表结构

#### users

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGINT UNSIGNED | PK, AUTO_INC | |
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
| title | VARCHAR(255) | NOT NULL | |
| subtitle | VARCHAR(255) | | |
| goal | TEXT | | |
| deliverables | TEXT | | |
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
| content | TEXT | NOT NULL | |
| estimated_hours | DECIMAL(4,1) | DEFAULT 0 | |
| resource_urls | JSON | | 内嵌 URL |
| sort_order | INT | DEFAULT 0 | |
| is_checkpoint | TINYINT(1) | DEFAULT 0 | |

#### user_tasks

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGINT UNSIGNED | PK | |
| user_id | BIGINT UNSIGNED | FK→users, NOT NULL | |
| task_id | BIGINT UNSIGNED | FK→tasks, NOT NULL | |
| is_completed | TINYINT(1) | DEFAULT 0 | |
| completed_at | TIMESTAMP | NULL | |
| learning_links | TEXT | | Markdown |
| implementation_plan | TEXT | | Markdown |
| implementation_code | TEXT | | Markdown |
| experience_summary | TEXT | | Markdown |
| created_at | TIMESTAMP | DEFAULT NOW() | |
| updated_at | TIMESTAMP | ON UPDATE NOW() | |

唯一约束: (user_id, task_id)，懒加载创建。

---

## 7. API 接口规格

### 7.1 响应格式

```
成功: { "code": 200, "data": {...} }
错误: { "code": 4xx/5xx, "message": "error description" }
```

### 7.2 接口列表

#### 公开接口

| 方法 | 路径 | 请求 | 响应 |
|------|------|------|------|
| POST | `/api/v1/auth/register` | `{username,email,password}` | `{token,user}` |
| POST | `/api/v1/auth/login` | `{email,password}` | `{token,user}` |

#### 需认证接口 (Authorization: Bearer token)

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/auth/me` | 当前用户信息 |
| GET | `/api/v1/phases` | Phase 列表（含进度） |
| GET | `/api/v1/phases/:id` | Phase 详情（含 Weeks） |
| GET | `/api/v1/weeks/:id` | Week 详情（含 Days+Tasks+完成状态） |
| GET | `/api/v1/tasks/:id` | Task 详情（含 user_task 提交内容） |
| PATCH | `/api/v1/tasks/:id/complete` | `{completed:bool}` → 切换完成 |
| PUT | `/api/v1/tasks/:id/submissions` | 更新四种提交内容 |
| GET | `/api/v1/progress/overview` | 总体进度统计 |
| GET | `/api/v1/dashboard` | Dashboard 聚合数据 |

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
| 侧边栏/卡片 | `#111827` / `#1a1f35` | 表面色 |
| 边框 | `#2a3558` / `#2a4068` | 普通/发光 |
| 青色 | `#00f0ff` | Phase 1、主强调 |
| 紫色 | `#7c3aed` | Phase 2 |
| 绿色 | `#00ff88` | Phase 3、成功 |
| 黄色 | `#ffb700` | 警告、检查点 |
| 红色 | `#ff3860` | 错误、今天标记 |
| 文字 | `#e2e8f0` / `#94a3b8` / `#64748b` | 主/次/禁用 |

### 8.2 字体

| 用途 | 字体 | 备选 |
|------|------|------|
| 标题 | Orbitron | sans-serif |
| 正文 | Inter | system-ui, sans-serif |
| 代码 | JetBrains Mono | Fira Code, monospace |

### 8.3 视觉效果

- Canvas 粒子星空背景（100 粒子，连线距离 120px 内）
- CSS 网格叠加层（60px 间距）
- 卡片悬停：边框发光 + 上移 4px
- 复选框：青色填充 + 缩放动画 + 发光
- SVG 进度环：stroke-dasharray 动画
- 侧边栏：完全隐藏/展开，0.3s 过渡

### 8.4 响应式断点

- `>1024px`: 正常布局（侧边栏 260px + 内容区）
- `640-1024px`: 侧边栏默认隐藏，顶部 ☰ 按钮
- `<640px`: 单列布局，padding 缩小

---

## 9. 数据导入

### 9.1 解析规则

- 状态机逐行解析 `web3_infra_3month_plan.md`
- 检测 `## Phase` / `### 第N周` / `#### Day` 标题
- 检测 `- [ ]` 行作为任务
- 提取 `[text](url)` 链接存入 resource_urls (JSON)
- 提取 `（Nh）` 格式小时数
- 检测 `**自检清单：**` 标记后续任务 is_checkpoint=1

### 9.2 执行方式

```bash
cd backend && go run cmd/seed/main.go --plan=../web3_infra_3month_plan.md
```

幂等导入：INSERT IGNORE，预期导入 3 Phase + 12 Week + 84 Day + ~232 Task。

---

## 10. 验收标准

| # | 功能 | 验收条件 |
|---|------|----------|
| 1 | 注册 | 成功返回 Token，重复用户名/邮箱返回 409 |
| 2 | 登录 | 正确密码返回 Token，错误返回 401 |
| 3 | Phase 列表 | 3 个 Phase 卡片，进度统计正确 |
| 4 | Phase 详情 | 四级手风琴展开，勾选任务更新进度 |
| 5 | Week 详情 | 显示该周全部 Task，完成状态正确 |
| 6 | Task 详情 | 4 Tab 编辑器，Markdown 预览正常 |
| 7 | 任务完成 | 勾选后进度实时更新，取消回退 |
| 8 | 提交内容 | 四种类型独立保存，自动保存提示 |
| 9 | 甘特图 | 12 周色块 + 完成率填充 + 里程碑 |
| 10 | Dashboard | 统计卡片 + 进度环 + 活动列表数据正确 |
| 11 | 学习手册 | handbook.md 正确渲染为 HTML |
| 12 | 用户隔离 | 用户 A 完成状态不影响用户 B |
| 13 | 侧边栏 | 隐藏/展开动画 + 左侧边缘滑出按钮 |
| 14 | Docker | docker-compose up 一键启动全站 |

---

## 11. 项目结构

```
Goto_Web3/
├── backend/
│   ├── cmd/server/main.go           # 应用入口
│   ├── cmd/seed/main.go             # 数据导入入口
│   ├── internal/
│   │   ├── config/config.go         # 环境变量
│   │   ├── database/                # 连接池 + 迁移
│   │   ├── middleware/              # JWT/CORS/日志
│   │   ├── model/                   # 6 个数据结构
│   │   ├── repository/              # 6 个数据库查询
│   │   ├── service/                 # 4 个业务服务
│   │   ├── handler/                 # 8 个 HTTP 处理器
│   │   ├── router/router.go         # 路由注册
│   │   └── importer/                # MD 解析器 + 填充器
│   ├── templates/ (10 .html)        # Go 模板
│   ├── static/{css,js,lib}/         # 静态资源
│   ├── go.mod / go.sum
│   └── Dockerfile
├── docker-compose.yml
├── sources/                         # 源 Markdown
│   ├── web3_infra_handbook.md
│   └── web3_infra_3month_plan.md
├── REQUIREMENTS.md
├── DESIGN.md
└── .gitignore
```
