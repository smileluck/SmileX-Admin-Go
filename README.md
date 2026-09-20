# SmileX Admin

<div align="center">

**DDD 驱动的全栈后台管理系统**

Gin · GORM · Wire · Vue3 · TypeScript · Naive UI

多数据库（MySQL / PostgreSQL / SQLite）· RBAC · JWT · 动态菜单 · 微服务可演进

</div>

---

## ✨ 特性

### 后端（DDD + Kratos 兼容分层）
- 🏗️ **DDD 四层架构**：`biz`（领域层）/ `data`（基础设施）/ `service`（应用层）/ `server`（传输层），依赖倒置，仓储接口定义在领域层
- 🧩 **四个限界上下文**：auth / user / role / permission，跨上下文仅通过最小接口或领域事件通信
- 🔌 **多数据库**：配置一行切换 MySQL / PostgreSQL / SQLite，自动建表 + 种子数据
- 🔐 **RBAC + JWT**：用户-角色-权限三级模型，API 路径通配匹配，access + refresh 双令牌
- 🎯 **google/wire 依赖注入**：与 Kratos 完全同构，**切换微服务时业务层零改动**
- 📜 **proto 契约先行**：`api/` 下已定义接口契约，平滑生成 gRPC/HTTP 代码

### 前端（Vue3 + Naive UI）
- 🎨 **现代化 UI**：Naive UI 按需加载，路由级代码分割
- 🧭 **动态菜单路由**：菜单完全由后台配置，`addRoute()` 运行时注册，改菜单不改代码
- 🔘 **按钮级权限**：`v-permission` 自定义指令
- 🔄 **Token 静默续期**：401 自动刷新重放，用户无感
- 📦 **两种部署**：单二进制（后端托管 SPA）/ Nginx 分离部署

## 🚀 快速开始

```bash
# 零依赖体验（SQLite）：configs/config.yaml 里 db.driver 改为 sqlite 即可
git clone https://github.com/yourname/SmileX-Admin-Gin.git
cd SmileX-Admin-Gin
make web-install web-build run
# 打开 http://localhost:8080；首次播种的 admin 初始密码随机生成，见启动日志打印（见下方"安全提醒"）
```

MySQL 方式：

```bash
mysql -uroot -p < migrations/mysql.sql   # 建库
# 修改 configs/config.yaml 的 db.mysql 连接信息后
make web-build run
```

前后端热更新开发（推荐，一条命令）：

```bash
go install github.com/air-verse/air@latest  # 首次需安装 air
make dev    # 后端 air 热加载 :8080 + 前端 Vite 热更新 :5173（/api 自动代理到后端）
# 打开 http://localhost:5173
```

或分开启动：

```bash
make run        # 终端1：后端 :8080
make web-dev    # 终端2：前端 :5173（/api 自动代理到后端）
```

Windows 环境说明：

- 安装 make：`choco install make` 或 `scoop install make`（MinGW 用户可直接用 `mingw32-make`），CMD / PowerShell / Git Bash 均可运行
- 所有 make 目标均已适配：`make build` 产物为 `bin\server.exe`，`make dev` 自动使用 `.air.windows.toml`
- `make dev` 在同一控制台并行启动后端 air 与前端 Vite，Ctrl+C 一并停止（与 macOS/Linux 行为一致）

## 🐳 Docker 部署（docker compose）

一条命令拉起完整服务（app + MySQL + Redis），首次启动自动建表 + 种子数据：

```bash
cp .env.example .env          # 并填写 MYSQL 密码与 JWT_SECRET
docker compose up -d          # 构建镜像并启动
# 打开 http://localhost:8080
```

说明：

- **多阶段构建**：Node 构建前端 dist → Go 编译单二进制（含 CGO，三种数据库驱动均可用）→ Alpine 运行时，镜像内已含 SPA 托管，无需 Nginx
- **配置注入**：所有配置项支持 `APP_` 前缀环境变量覆盖（`.` → `_`，如 `APP_DB_MYSQL_HOST` 覆盖 `db.mysql.host`），优先级高于镜像内置的 config.yaml；改完 `docker compose up -d` 重建即生效
- **数据持久化**：named volume——`app-data`（上传/导出/sqlite）、`app-logs`（应用日志）、`mysql-data`、`redis-data`
- **换数据库**：`.env` 里 `DB_DRIVER=postgres` 并放开 docker-compose.yml 中的 postgres 服务；`DB_DRIVER=sqlite` 则零外部依赖（数据落 `app-data` 卷）
- **默认生产语义**：release 模式、登录验证码开启、日志同时输出控制台（`docker logs`）
- **端口冲突**：宿主机 8080 已被占用（如本地开发后端）时，在 `.env` 中改 `APP_PORT` 换对外端口
- 常用命令：`make docker-up` / `make docker-down` / `make docker-logs` / `make docker-rebuild`


## 📖 API（/api/v1）

| 方法 | 路径 | 说明 |
|---|---|---|
| POST | /auth/login | 登录 |
| POST | /auth/refresh | 刷新令牌 |
| GET | /auth/profile | 当前用户信息（含权限） |
| GET | /menus | 当前用户可见菜单树 |
| GET/POST | /users | 用户列表 / 创建 |
| GET/PUT/DELETE | /users/:id | 用户详情 / 更新 / 删除 |
| PUT | /users/:id/roles | 分配角色 |
| PUT | /users/:id/password | 重置密码 |
| GET/POST | /roles | 角色列表 / 创建 |
| PUT | /roles/:id/permissions | 绑定权限 |
| GET/POST | /permissions | 权限列表 / 创建（支持 type=menu 筛选） |
| GET/POST | /ip-blacklist | IP 黑名单列表 / 新增 |
| DELETE | /ip-blacklist/:id | 解封指定黑名单 |

## 🏛️ 架构与依赖方向

```
server ──▶ service ──▶ biz ◀── data
 (Gin)     (用例编排)  (领域层)  (GORM/PG/SQLite/JWT)
                     ▲
              （仓储接口在 biz，data 实现 —— 依赖倒置）
```

- 依赖倒置：仓储接口定义在 `internal/biz/*/repo.go`，`internal/data/*` 实现
- 限界上下文：auth / user / role / permission 各自独立，**跨上下文只允许依赖对方 biz 层接口**（如 `auth.UserReader`），或走 `pkg/eventbus` 领域事件
- PO（持久化对象）与领域实体分离，转换在 `internal/data/model`

## 📈 微服务升级路径

单体优先，架构预留拆分能力，切换到 Kratos 时 **biz / data / service 三层零改动**：

1. **proto 契约先行**：`api/` 下已定义 auth / admin 的 proto（首期手写实现与其对齐）
2. **拆服务**：把目标上下文的 `biz` + `data` + `service` 平移到 Kratos 工程，`internal/server` 层由 proto 生成的 HTTP/gRPC server 代替
3. **DI 复用**：wire Provider 图与 Kratos 完全同构
4. **服务间通信**：gRPC（契约即 `api/proto`）
5. **异步解耦**：`pkg/eventbus` 进程内总线替换为 MQ（NSQ/Kafka），业务代码不变
6. **可观测**：接入 OpenTelemetry

## 🎨 前端设计要点（web/）

- **动态菜单**：登录后拉取 `/auth/profile` + `/menus`，菜单树动态 `addRoute()` 注册路由，侧边栏递归渲染——菜单在"菜单管理"页面配置，无需改前端代码
- **按钮级权限**：`v-permission="['menu:user']"` 指令控制显隐
- **token 静默续期**：axios 拦截器捕获 401，自动用 refresh token 换新并重放请求
- **部署两种方式**：① 后端托管 `web/dist`（单二进制，SPA history fallback）；② Nginx 独立部署 + 反代 `/api`
- 新增页面：在 `web/src/router/index.ts` 的 `viewModules` 登记菜单 code → 组件映射，再到"菜单管理"里建菜单即可

## 🗂️ 目录结构

```
cmd/server/        # 入口 + wire 注入
api/               # proto 契约
configs/           # 配置
internal/biz/      # 领域层（实体/值对象/仓储接口/Usecase）
internal/data/     # 基础设施（GORM PO/仓储实现/JWT/多数据库工厂）
internal/service/  # 应用层（薄用例）
internal/server/   # 传输层（Gin 路由/中间件/SPA 静态托管）
internal/conf/     # 配置加载
pkg/               # eventbus / response / logger / pagination
web/               # 前端（Vue3 + TS + Naive UI）
migrations/        # 三方言建表 SQL
```

## 🛠️ 技术栈

| 层 | 技术 |
|---|---|
| Web | Gin、JWT、CORS、RBAC 中间件 |
| ORM | GORM（MySQL / PostgreSQL / SQLite） |
| DI | google/wire |
| 前端 | Vue3、TypeScript、Vite、Naive UI、Pinia |
| 日志 | zap |

## 🧰 常用命令

```bash
make build       # 编译后端
make run         # 运行后端
make wire        # 重新生成 DI（改 Provider 后执行）
make test        # 测试
make web-install # 安装前端依赖
make web-dev     # 前端开发（热更新）
make web-build   # 前端构建（产物 web/dist，由后端静态托管）
make docker-up   # Docker 构建并启动完整服务（app + MySQL + Redis）
make docker-down # 停止并移除容器
```

## ⚠️ 安全提醒

默认账号 `admin`，初始密码按以下优先级确定（仅首次播种空库时生效一次）：

1. 环境变量 `APP_SEED_ADMIN_PASSWORD`
2. 配置 `seed.adminPassword`（configs/config.yaml）
3. 均未设置时随机生成 12 位密码，仅在启动日志打印一次

密钥与敏感配置：

- 仓库不保存任何可用密钥：`jwt.secret` 与 `seed.adminPassword` 默认为空。
- 本地开发（debug 模式）secret 留空可正常启动（仅告警）；**release 模式下空或长度不足 32 位的
  jwt secret 会拒绝启动**，须以环境变量 `APP_JWT_SECRET` 注入强随机值（`openssl rand -base64 32`）。
- 更换 jwt secret 后所有已签发 token 失效需重新登录；若未单独配置 `agent.cryptoKey`，
  已保存的供应商 API Key 需重新保存一次。
- 跨域访问 API 需配置 `server.corsOrigins` 白名单（默认同源模式，不下发 CORS 头）。

## License

MIT
