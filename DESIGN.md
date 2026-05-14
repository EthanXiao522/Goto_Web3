# Web3 Infra 学习追踪网站 - 开发设计文档

> 版本：v1.0 | 日期：2026-05-14 | 依据：REQUIREMENTS.md v1.0 | 状态：待确认

---

## 目录

- [1. 系统架构](#1-系统架构)
  - [1.1 三层架构总览](#11-三层架构总览)
  - [1.2 技术栈与模块映射](#12-技术栈与模块映射)
  - [1.3 请求处理流程](#13-请求处理流程)
  - [1.4 部署架构](#14-部署架构)
- [2. 技术栈详情](#2-技术栈详情)
- [3. 项目结构](#3-项目结构)
- [4. 模块设计](#4-模块设计)
  - [4.0 模块功能总览](#40-模块功能总览)
  - [4.1 模块依赖关系](#41-模块依赖关系)
  - [4.2 模块详解](#42-模块详解)
- [5. 数据库详细设计](#5-数据库详细设计)
  - [5.0 实体关系图](#50-实体关系图)
  - [5.1 完整 DDL](#51-完整-ddl)
  - [5.2 关键查询设计](#52-关键查询设计)
- [6. 前端设计](#6-前端设计)
  - [6.0 前端功能总览](#60-前端功能总览)
  - [6.1 模板方案](#61-模板方案)
  - [6.2 JavaScript 模块](#62-javascript-模块)
  - [6.3 API 调用模式](#63-api-调用模式)
- [7. 关键流程设计](#7-关键流程设计)
  - [7.1 任务完成流程](#71-任务完成流程)
  - [7.2 提交保存流程](#72-提交保存流程)
  - [7.3 页面加载流程](#73-页面加载流程)
  - [7.4 数据导入流程](#74-数据导入流程)
- [8. Docker 配置](#8-docker-配置)
  - [8.1 Dockerfile](#81-dockerfile)
  - [8.2 docker-compose.yml](#82-docker-composeyml)
- [9. 开发顺序](#9-开发顺序)
- [10. 错误处理策略](#10-错误处理策略)
- [11. 测试策略](#11-测试策略)
- [附录 A：依赖清单](#附录-a依赖清单)

---

## 1. 系统架构

### 1.1 三层架构总览

```
+--------------------------------------------------------------+
|                    Application Layer                          |
|                                                              |
|  +----------+ +----------+ +----------+ +----------+         |
|  | HTML     | | CSS      | | JS       | | Static   |         |
|  | Go Tmpl  | | 600 lines| | Vanilla  | | marked.js|         |
|  | 10 pages | | Cyberpunk| | 5 modules| | Canvas   |         |
|  | SSR      | | Responsiv| | 400 lines| | SVG Ring |         |
|  +----+-----+ +----------+ +----+-----+ +----------+         |
|       +----------+---------------+                           |
|                  | HTTP / JSON                                |
+------------------+-------------------------------------------+
|                   Business Layer                              |
|                                                              |
|  +------------------------------------------------------+   |
|  |              Gin Router (HTTP Framework)              |   |
|  |  +---------+  +---------------+  +-----------------+  |   |
|  |  |Middleware|  |  Page Routes  |  |   API Routes    |  |   |
|  |  |CORS      |  |  GET /        |  | POST /api/auth  |  |   |
|  |  |Logger    |  |  GET /phases  |  | GET  /api/phases|  |   |
|  |  |JWT Auth  |  |  GET /gantt   |  | PATCH /api/task |  |   |
|  |  |Recovery  |  |  GET /db      |  | PUT  /api/task  |  |   |
|  |  +---------+  +---------------+  +-----------------+  |   |
|  +----------------------+--------------------------------+   |
|                         |                                    |
|  +----------------------+--------------------------------+   |
|  |               Handler Layer (8 handlers)              |   |
|  |  Auth|Phase|Week|Task|Progress|Dashboard|Gantt|Book   |   |
|  +----------------------+--------------------------------+   |
|                         |                                    |
|  +----------------------+--------------------------------+   |
|  |               Service Layer (4 services)              |   |
|  |  AuthService      | TaskService     | ProgressService |   |
|  |  Register/Login   | Toggle/Submit   | GetOverview     |   |
|  |  GenerateJWT      | LazyCreate      | DashboardService|   |
|  +----------------------+--------------------------------+   |
|                         |                                    |
|  +----------------------+--------------------------------+   |
|  |            Repository Layer (6 repositories)          |   |
|  |  User|Phase|Week|Day|Task|UserTask (raw SQL)          |   |
|  +----------------------+--------------------------------+   |
|                         |                                    |
+-------------------------+------------------------------------+
|                    Data Layer                                 |
|  +----------------------+--------------------------------+   |
|  |                  MySQL 8.0 Database                  |   |
|  |  users -- phases -- weeks -- days -- tasks           |   |
|  |    |                                                |   |
|  |    + -- user_tasks (user_id, task_id)               |   |
|  |  Engine: InnoDB | Charset: utf8mb4 | FK: CASCADE    |   |
|  +------------------------------------------------------+   |
|                                                              |
|  +------------------------------------------------------+   |
|  |                   Docker Compose                      |   |
|  |  +--------------+       +--------------------------+  |   |
|  |  | mysql:8.0    |<----->| Go App (alpine)          |  |   |
|  |  | Port:3306    |network| Port:8080                |  |   |
|  |  | Volume       |       | Multi-stage build        |  |   |
|  |  | Healthcheck  |       | templates/ + static/     |  |   |
|  |  +--------------+       +--------------------------+  |   |
|  +------------------------------------------------------+   |
+--------------------------------------------------------------+
```

### 1.2 技术栈与模块映射

```
         Frontend                            Backend
  +-------------------+           +-------------------+
  | Go html/template  |           | Gin v1.9          |
  | Vanilla JS ES6+   |<---API--->| golang-jwt v5.2   |
  | CSS3 Custom Prop  |           | go-sql-driver     |
  | Canvas API        |           | x/crypto/bcrypt   |
  | SVG Graphics      |           | log/slog stdlib   |
  | marked.js (CDN)   |           | encoding/json     |
  +-------------------+           +-------------------+

         Infrastructure
  +------------------------------------------+
  | Docker+Compose | MySQL 8.0 | Alpine      |
  | Git(feature commits) | Volume persist    |
  +------------------------------------------+
```

### 1.3 请求处理流程

```
Browser Request ──► Gin Router
                       │
                       ├── Page Route (/dashboard, /phases, ...)
                       │   └── Handler → Service → Repo → DB
                       │       └── c.HTML(200, template, data)
                       │
                       └── API Route (/api/v1/*)
                           ├── CORS Middleware
                           ├── Logger Middleware
                           ├── Auth Middleware (JWT)
                           └── Handler → Service → Repo → DB
                               └── c.JSON(200, response)
```

### 1.4 部署架构

```
  docker-compose.yml
  +--------------------------------------+
  |                                      |
  |  +------------+   +----------------+ |
  |  |  MySQL 8.0 |   |  Go App :8080  | |
  |  |  Container |   |  Container     | |
  |  |  vol:mysql |   |  build:alpine  | |
  |  |  healthck  |   |  templates+CSS | |
  |  +-----+------+   +-------+--------+ |
  |        +-- internal net --+          |
  +--------------------------------------+
```

---

## 2. 技术栈详情

| 层面 | 选型 | 版本 | 用途 |
|------|------|------|------|
| Go | go.mod | 1.21 | 主语言 |
| Gin | github.com/gin-gonic/gin | v1.9 | HTTP 框架 |
| MySQL Driver | github.com/go-sql-driver/mysql | v1.7 | 数据库连接 |
| JWT | github.com/golang-jwt/jwt/v5 | v5.2 | Token 签发/验证 |
| bcrypt | golang.org/x/crypto | latest | 密码加密 |
| marked.js | CDN | v9+ | 前端 Markdown 渲染 |
| CSS | 手写 | - | 赛博朋克主题 ~600 行 |
| JS | Vanilla ES6+ | - | 5 个模块 ~400 行 |
| Docker | 多阶段构建 | - | golang:1.21-alpine → alpine:3.19 |
| MySQL | Docker Image | 8.0 | utf8mb4 |

---

## 3. 项目结构

```
Goto_Web3/
├── backend/
│   ├── cmd/
│   │   ├── server/main.go              # 应用入口
│   │   └── seed/main.go                # 数据导入入口
│   ├── internal/
│   │   ├── config/config.go            # 环境变量 (PORT/DB_DSN/JWT_SECRET)
│   │   ├── database/
│   │   │   ├── mysql.go                # 连接池 (MaxOpen=25, MaxIdle=5)
│   │   │   └── migrate.go              # 6 表自动创建
│   │   ├── middleware/
│   │   │   ├── auth.go                 # JWT 验证 + user_id 注入 ctx
│   │   │   ├── cors.go                 # CORS 头
│   │   │   └── logger.go              # slog 请求日志
│   │   ├── model/                      # 纯数据结构 (6 files)
│   │   ├── repository/                 # 原生 SQL 查询 (6 files)
│   │   ├── service/                    # 业务逻辑 (4 files)
│   │   ├── handler/                    # HTTP 处理 (8 files)
│   │   ├── router/router.go            # 路由注册
│   │   └── importer/
│   │       ├── parser.go               # MD 状态机解析
│   │       └── seeder.go               # DB 填充 (INSERT IGNORE)
│   ├── templates/ (10 .html)           # Go html/template
│   ├── static/{css,js,lib}/            # 静态资源
│   ├── go.mod / go.sum
│   └── Dockerfile
├── docker-compose.yml
├── sources/                            # 源 Markdown
│   ├── web3_infra_handbook.md
│   └── web3_infra_3month_plan.md
├── REQUIREMENTS.md
├── DESIGN.md
└── .gitignore
```

---

## 4. 模块设计

### 4.0 模块功能总览

| 模块 | 文件数 | 功能说明 |
|------|--------|----------|
| config | 1 | 环境变量加载，默认值 fallback |
| database | 2 | MySQL 连接池 + 6 表 DDL 自动迁移 |
| middleware | 3 | JWT 认证 / CORS 跨域 / slog 请求日志 |
| model | 6 | User, Phase, Week, Day, Task, UserTask |
| repository | 6 | 每表原生 SQL CRUD，参数化查询，手动 Scan |
| service | 4 | Auth(注册/登录/JWT) / Task(完成/提交/懒加载) / Progress(聚合) / Dashboard |
| handler | 8 | Auth, Phase, Week, Task, Progress, Dashboard, Gantt, Handbook |
| router | 1 | 公开路由 + JWT 保护路由 + API 路由分组 |
| importer | 2 | 状态机解析 MD + INSERT IGNORE 幂等导入 |
| templates | 10 | landing/login/register/dashboard/phases/phase_detail/week_detail/task_detail/gantt/handbook |
| static/css | 1 | 赛博朋克主题: 变量/布局/卡片/表单/进度/手风琴/Tab/甘特图/动效/响应式 |
| static/js | 5 | app(认证/粒子/Toast/侧边栏) / phase(手风琴/勾选) / task(Tab/预览/保存) / dashboard(进度环) / gantt(甘特图) |
| cmd/server | 1 | 入口: config → DB → migrate → router → Listen:8080 |
| cmd/seed | 1 | 入口: config → DB → parser → seeder |
| Docker | 2 | Dockerfile(多阶段) + docker-compose(MySQL+App) |

### 4.1 模块依赖关系

```
cmd/server/main.go
    ├── config.Load()
    ├── database.Connect() + Migrate()
    └── router.Setup(handlers)
            └── handler
                └── service
                    └── repository
                        └── database.DB

cmd/seed/main.go
    ├── config.Load()
    ├── database.Connect()
    └── importer.Parse() → importer.Seed()
```


### 4.2 模块详解

#### config

```go
func Load() *Config {
    return &Config{
        Port:      getEnv("PORT", "8080"),
        DBDSN:     getEnv("DB_DSN", "root:pass@tcp(127.0.0.1:3306)/web3_learning?..."),
        JWTSecret: getEnv("JWT_SECRET", "dev-secret"),
    }
}
```

#### database

```go
func Connect(dsn string) error  // sql.Open + Ping, MaxOpenConns=25
func Migrate() error            // 6 × CREATE TABLE IF NOT EXISTS
func Close()                    // DB.Close()
```

#### middleware

```go
// auth.go - 从 Authorization: Bearer <token> 提取验证 JWT
func Auth() gin.HandlerFunc {
    // 解析 token → 验证 → c.Set("user_id", claims.UserID)
}

// cors.go
func CORS() gin.HandlerFunc {
    // Access-Control-Allow-Origin: *
}

// logger.go
func Logger() gin.HandlerFunc {
    // slog.Info("request", "method", c.Request.Method, "path", path, "status", status, "duration", d)
}
```

#### model (示例: Task)

```go
type Task struct {
    ID             uint64    `json:"id"`
    DayID          uint64    `json:"day_id"`
    Content        string    `json:"content"`
    EstimatedHours float64   `json:"estimated_hours"`
    ResourceURLs   string    `json:"resource_urls"`   // JSON string
    SortOrder      int       `json:"sort_order"`
    IsCheckpoint   bool      `json:"is_checkpoint"`
}
```

#### repository (示例: UserRepo)

```go
type UserRepo struct{ db *sql.DB }

func (r *UserRepo) Create(user *model.User) (uint64, error)
func (r *UserRepo) FindByEmail(email string) (*model.User, error)
func (r *UserRepo) FindByID(id uint64) (*model.User, error)
// 禁止 ORM，参数化查询，sql.ErrNoRows → ErrNotFound
```

#### service (示例: TaskService)

```go
type TaskService struct {
    taskRepo     *TaskRepo
    userTaskRepo *UserTaskRepo
}

// ToggleComplete: getOrCreateUserTask → update is_completed + completed_at
func (s *TaskService) ToggleComplete(userID, taskID uint64, completed bool) (*model.UserTask, error)

// UpdateSubmissions: 只更新 non-empty 字段
func (s *TaskService) UpdateSubmissions(userID, taskID uint64, fields map[string]string) (*model.UserTask, error)

// getOrCreateUserTask: INSERT IGNORE (懒加载)
func (s *TaskService) getOrCreateUserTask(userID, taskID uint64) (*model.UserTask, error)
```

#### handler (示例: PhaseHandler)

```go
type PhaseHandler struct {
    phaseRepo       *PhaseRepo
    progressService *ProgressService
}

// API: GET /api/v1/phases → JSON
func (h *PhaseHandler) GetPhases(c *gin.Context) {
    userID := c.GetUint64("user_id")
    phases, _ := h.phaseRepo.GetAllWithProgress(userID)
    c.JSON(200, gin.H{"code": 200, "data": gin.H{"phases": phases}})
}

// Page: GET /phases → HTML
func (h *PhaseHandler) PhaseListPage(c *gin.Context) {
    userID := c.GetUint64("user_id")
    phases, _ := h.phaseRepo.GetAllWithProgress(userID)
    c.HTML(200, "phases.html", gin.H{"User": user, "Phases": phases})
}
```

#### importer

状态机逐行解析：

```
OUTSIDE ──► IN_PHASE      (## Phase N：title)
         ├─► IN_WEEK       (### 第N周：title)
         ├─► IN_DAY        (#### Day N...)
         ├─► IN_TASKS      (- [ ] task content)
         └─► IN_CHECKPOINT (**自检清单：**)
```

---

## 5. 数据库详细设计

### 5.0 实体关系图

```
+----------+                        +----------+
|  users   |                        |  phases  |
+----------+                        +----------+
| id (PK)  |                        | id (PK)  |
| username |                        |phase_num |
| email    |                        | title    |
| pass_hash|                        | goal     |
+----+-----+                        +----+-----+
     | 1                                 | 1
     |                                   |
     |                              +----+-----+
     |                              |  weeks   |
     |                              +----------+
     |                              | id (PK)  |
     |                              |phase(FK) |<-- phase_id
     |                              |week_num  |
     |                              +----+-----+
     |                                   | 1
     |                              +----+-----+
     |                              |   days   |
     |                              +----------+
     |                              | id (PK)  |
     |                              |week(FK)  |<-- week_id
     |                              +----+-----+
     |                                   | 1
     |       +---------------------+----+-----+
     |       |                     |  tasks    |
     |       |                     +-----------+
     |       |                     | id (PK)   |
     |       |                     |day(FK)    |<-- day_id
     |       |                     | content   |
     |       |                     | res_urls  | (JSON)
     |       |                     +-----+-----+
     |  N    |                         N |
     |       +----------+----------------+
     |                  |
+----+------------------+-----+
|         user_tasks          |
+-----------------------------+
| id (PK)                     |
| user_id (FK -> users)       |-- ON DELETE CASCADE
| task_id (FK -> tasks)       |-- ON DELETE CASCADE
| is_completed                |
| --------------------------- |
| learning_links     (TEXT)   |-- 4 submit types
| implementation_plan (TEXT)  |
| implementation_code (TEXT)  |
| experience_summary  (TEXT)  |
| --------------------------- |
| UNIQUE (user_id, task_id)   |-- lazy create
+-----------------------------+

Query path:
phases <- weeks <- days <- tasks <- user_tasks (JOIN on user_id)
```

### 5.1 完整 DDL

```sql
CREATE TABLE users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(64) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE phases (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    phase_number TINYINT UNSIGNED NOT NULL,
    title VARCHAR(255) NOT NULL,
    subtitle VARCHAR(255) DEFAULT '',
    goal TEXT,
    deliverables TEXT,
    sort_order INT NOT NULL DEFAULT 0,
    INDEX idx_phases_sort (sort_order, phase_number)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE weeks (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    phase_id BIGINT UNSIGNED NOT NULL,
    week_number TINYINT UNSIGNED NOT NULL,
    title VARCHAR(255) NOT NULL,
    subtitle VARCHAR(255) DEFAULT '',
    goal TEXT,
    deliverables TEXT,
    sort_order INT NOT NULL DEFAULT 0,
    UNIQUE INDEX idx_weeks_number (week_number),
    INDEX idx_weeks_phase (phase_id, sort_order),
    FOREIGN KEY (phase_id) REFERENCES phases(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE days (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    week_id BIGINT UNSIGNED NOT NULL,
    day_number TINYINT UNSIGNED NOT NULL,
    title VARCHAR(255) NOT NULL,
    sort_order INT NOT NULL DEFAULT 0,
    INDEX idx_days_week (week_id, sort_order),
    FOREIGN KEY (week_id) REFERENCES weeks(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE tasks (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    day_id BIGINT UNSIGNED NOT NULL,
    content TEXT NOT NULL,
    estimated_hours DECIMAL(4,1) DEFAULT 0,
    resource_urls JSON,
    sort_order INT NOT NULL DEFAULT 0,
    is_checkpoint TINYINT(1) DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_tasks_day (day_id, sort_order),
    FOREIGN KEY (day_id) REFERENCES days(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE user_tasks (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    task_id BIGINT UNSIGNED NOT NULL,
    is_completed TINYINT(1) DEFAULT 0,
    completed_at TIMESTAMP NULL,
    learning_links TEXT,
    implementation_plan TEXT,
    implementation_code TEXT,
    experience_summary TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uq_user_task (user_id, task_id),
    INDEX idx_user_tasks_user (user_id, is_completed),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### 5.2 关键查询设计

**Phase 列表（含进度）**

```sql
SELECT p.*,
  (SELECT COUNT(*)
   FROM weeks w JOIN days d ON d.week_id = w.id
   JOIN tasks t ON t.day_id = d.id
   WHERE w.phase_id = p.id AND t.is_checkpoint = 0) AS task_count,
  (SELECT COUNT(*)
   FROM weeks w JOIN days d ON d.week_id = w.id
   JOIN tasks t ON t.day_id = d.id
   JOIN user_tasks ut ON ut.task_id = t.id
   WHERE w.phase_id = p.id AND t.is_checkpoint = 0
   AND ut.user_id = ? AND ut.is_completed = 1) AS completed_count
FROM phases p ORDER BY p.sort_order;
```

**user_tasks 懒加载**

```sql
INSERT IGNORE INTO user_tasks (user_id, task_id) VALUES (?, ?);
```

---

## 6. 前端设计

### 6.0 前端功能总览

| 页面 | 路由 | 关键组件与交互 |
|------|------|----------------|
| Landing | `/` | 粒子星空 Canvas │ 渐变 Hero │ 特性卡片×3 │ 登录/注册入口 |
| Login | `/login` | 邮箱+密码表单 │ 粒子背景 │ 登录后跳转 Dashboard |
| Register | `/register` | 用户名+邮箱+密码 │ 唯一性校验 │ 注册成功自动登录 |
| Dashboard | `/dashboard` | 统计卡片×4 │ SVG 进度环×3 │ 每周进度条×12 │ 最近活动×5 |
| Phase List | `/phases` | Phase 卡片×3 (青/紫/绿) │ 12 周进度表格 │ 点击进入详情 |
| Phase Detail | `/phases/:id` | **核心页**: 四级手风琴 │ 动画复选框 │ 提交指示点 │ 检查点卡片 |
| Week Detail | `/weeks/:id` | 单周 Day 分组 │ 任务列表+完成状态 │ 面包屑导航 |
| Task Detail | `/tasks/:id` | 4 Tab 编辑器 │ 编辑/预览切换 │ marked.js 渲染 │ 自动保存(1s) │ 上/下任务 |
| Gantt | `/gantt` | 自定义 CSS 甘特图 │ Phase 颜色填充 │ 周/天视图切换 │ 里程碑卡片×4 |
| Handbook | `/handbook` | handbook.md 服务端渲染 │ 侧边栏章节目录 │ 结构预览卡片 |

| 公共组件 | 说明 |
|----------|------|
| 侧边栏 | 260px 默认展开，◀ 按钮完全隐藏，左侧边缘滑出 ☰ 展开按钮 |
| Toast | 操作反馈: success(绿)/error(红)，3s 自动消失 |
| 网格背景 | CSS linear-gradient 60px 间距 |
| 粒子背景 | Canvas 100 粒子 + 120px 连线（Landing/Login/Register） |
| 响应式 | >1024px 正常布局 │ 640-1024px 侧边栏默认隐藏 │ <640px 单列 |

### 6.1 模板方案

Go template 使用组合模式：每个页面定义 `{{define "content"}}`，公共 layout 使用 `{{template "content" .}}` 占位。每个 handler 传入共同数据：

```go
c.HTML(200, "page.html", gin.H{
    "Title": "Dashboard",
    "User":  user,
    "Data":  pageData,
})
```

### 6.2 JavaScript 模块

| 文件 | 功能 | 关键函数 |
|------|------|----------|
| `app.js` | 认证管理 / 粒子背景 / Toast / 侧边栏 | `getToken()`, `fetchAPI()`, `showToast()`, `initSidebar()` |
| `phase.js` | 手风琴展开折叠 / 任务勾选 / 进度更新 | `toggleAccordion()`, `toggleTask()`, `updateProgress()` |
| `task.js` | Tab 切换 / 编辑预览切换 / 自动保存 | `switchTab()`, `switchMode()`, `autoSave()` |
| `dashboard.js` | 进度环 SVG 动画 | `animateRings()` |
| `gantt.js` | 甘特图渲染 / 视图切换 | `renderGantt()`, `switchView()` |

### 6.3 API 调用模式

```javascript
async function fetchAPI(path, options = {}) {
    const token = localStorage.getItem('token');
    const res = await fetch(path, {
        ...options,
        headers: {
            'Content-Type': 'application/json',
            ...(token ? {'Authorization': `Bearer ${token}`} : {}),
        },
    });
    if (res.status === 401) {
        localStorage.removeItem('token');
        window.location.href = '/login';
    }
    return res.json();
}

// Usage
const data = await fetchAPI('/api/v1/phases');
await fetchAPI('/api/v1/tasks/42/complete', {
    method: 'PATCH',
    body: JSON.stringify({ completed: true }),
});
```

---

## 7. 关键流程设计

### 7.1 任务完成流程

```
User clicks checkbox
  → JS: toggleTask(id)
  → Optimistic UI (checkbox 动画 + 进度更新)
  → PATCH /api/v1/tasks/:id/complete
    → AuthMiddleware → TaskHandler.ToggleComplete
      → TaskService.ToggleComplete
        → getOrCreateUserTask (INSERT IGNORE)
        → userTaskRepo.Update(is_completed, completed_at)
  → Response → JS 收到 { code:200, data:{user_task} }
    → Success: keep optimistic UI
    → Error: revert checkbox, showToast("error", msg)
```

### 7.2 提交保存流程

```
User types in textarea
  → Debounce 1000ms
  → PUT /api/v1/tasks/:id/submissions
    → AuthMiddleware → TaskHandler.UpdateSubmissions
      → TaskService.UpdateSubmissions
        → getOrCreateUserTask
        → userTaskRepo.UpdateFields (only non-empty)
  → JS show "已保存 · HH:MM:SS"
```

### 7.3 页面加载流程

```
Browser → GET /phases/1
  → PhaseHandler.PhaseDetailPage
    → phaseRepo.FindByID(1) → phase
    → weekRepo.FindByPhase(1) → weeks[]
    → for each week:
      → dayRepo.FindByWeek(weekID) → days[]
      → for each day:
        → taskRepo.FindByDay(dayID) → tasks[]
        → userTaskRepo.FindByUserAndTask(userID, ids) → map
    → c.HTML(200, "phase_detail.html", nestedData)
  → Browser renders + phase.js attaches listeners
```

### 7.4 数据导入流程

```bash
go run cmd/seed/main.go --plan=../web3_infra_3month_plan.md
```

```
main() → config.Load() → database.Connect()
  → importer.Parse(file) → ParsedData
  → importer.Seed(DB, data)
    → INSERT IGNORE phases (×3)
    → INSERT IGNORE weeks  (×12 + phase_id)
    → INSERT IGNORE days   (×84 + week_id)
    → INSERT IGNORE tasks  (×232 + day_id)
  → "Imported: 3 phases, 12 weeks, 84 days, 232 tasks"
```

---

## 8. Docker 配置

### 8.1 Dockerfile

```dockerfile
# Stage 1: Build
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /server cmd/server/main.go

# Stage 2: Runtime
FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /server /server
COPY templates/ /templates/
COPY static/ /static/
EXPOSE 8080
CMD ["/server"]
```

### 8.2 docker-compose.yml

```yaml
services:
  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: ${MYSQL_ROOT_PASSWORD:-web3pass}
      MYSQL_DATABASE: web3_learning
    ports: ["3306:3306"]
    volumes: ["mysql_data:/var/lib/mysql"]
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost"]
      interval: 5s
      retries: 10

  app:
    build: ./backend
    environment:
      PORT: "8080"
      DB_DSN: "root:${MYSQL_ROOT_PASSWORD:-web3pass}@tcp(mysql:3306)/web3_learning?parseTime=true&charset=utf8mb4"
      JWT_SECRET: ${JWT_SECRET:-change-me-in-production}
    ports: ["8080:8080"]
    depends_on:
      mysql:
        condition: service_healthy

volumes:
  mysql_data:
```

---

## 9. 开发顺序

| # | 模块 | 交付物 | 估算 |
|---|------|--------|------|
| 1 | 项目初始化 | go.mod, main.go, config, database, migrate | 30min |
| 2 | 数据导入 | importer/parser.go, importer/seeder.go, cmd/seed | 2h |
| 3 | 用户认证 | model/user, repo, service, handler, middleware/auth | 1.5h |
| 4 | Phase/Week/Day API | 3 组 model+repo+service+handler | 2h |
| 5 | Task API | model+repo+service+handler (完成/提交) | 2h |
| 6 | 进度/Dashboard | progress_service, dashboard_service, handlers | 1.5h |
| 7 | 页面模板 | 10 .html 模板，layout 模式 | 3h |
| 8 | CSS 样式 | style.css（从 prototype 迁移） | 2h |
| 9 | JS 交互 | 5 个 JS 模块 | 3h |
| 10 | Docker | Dockerfile + docker-compose.yml | 1h |

---

## 10. 错误处理策略

| 层级 | 策略 |
|------|------|
| Repository | sql.ErrNoRows → 自定义 ErrNotFound |
| Service | fmt.Errorf("service: %w", err) 包装上下文 |
| Handler | 400(参数错误) / 401(未认证) / 404(不存在) / 500(内部错误) |
| 全局 | Gin Recovery 中间件捕获 panic → 500 |

---

## 11. 测试策略

| 层级 | 类型 | 工具 | 目标 |
|------|------|------|------|
| Repository | 单元测试 | go test + 测试DB | 每个查询方法 |
| Service | 单元测试 | go test + mock | 核心业务逻辑 |
| Handler | 集成测试 | httptest + 测试DB | API 端点 |
| Importer | 单元测试 | go test + fixture MD | 解析准确性 |
| 前端 | 手动测试 | 浏览器 | 14 条验收标准 |

---

## 附录 A：依赖清单

```
github.com/gin-gonic/gin
github.com/go-sql-driver/mysql
github.com/golang-jwt/jwt/v5
golang.org/x/crypto
```

所有前端依赖通过 CDN 引入，无 npm 依赖。
