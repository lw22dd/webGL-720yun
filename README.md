# WebGL-720yun 项目文档

## 项目概述

这是一个基于WebGL技术的720云全景展示平台，支持全景图的上传、展示和交互。项目采用前后端分离架构，使用Docker容器化部署，确保开发环境的一致性和便捷性。

插入默认管理员用户 admin
-- 密码：admin123（已加密）
## 技术栈

### 后端技术栈

| 技术/框架 | 版本 | 用途 | 备注 |
|---------|------|------|------|
| Go | 1.25.3 | 后端开发语言 | 容器化环境使用 |
| Gin | v1.11.0 | Web框架 | 路由处理、中间件支持 |
| GORM | v1.31.1 | ORM框架 | 数据库操作 |
| MySQL | 8.0 | 关系型数据库 | 存储用户数据、全景图元数据 |
| Redis | 7.0 | 缓存数据库 | 会话管理、缓存热点数据 |
| JWT | v5.3.0 | 身份认证 | 用户登录、权限验证 |
| Viper | v1.21.0 | 配置管理 | 读取配置文件 |
| Logrus | v1.9.3 | 日志管理 | 系统日志记录 |
| Zap | v1.27.0 | 高性能日志 | Gin中间件日志 |

### 前端技术栈

| 技术/框架 | 版本 | 用途 | 备注 |
|---------|------|------|------|
| Vue | 3.5.12 | 前端框架 | 组件化开发 |
| TypeScript | ~5.6.2 | 类型安全 | 静态类型检查 |
| Vite | 5.4.9 | 构建工具 | 快速开发、热更新 |
| Element Plus | 2.12.0 | UI组件库 | 表单、按钮等组件 |
| Pinia | 3.0.3 | 状态管理 | 全局状态管理 |
| Vue Router | 4.5.1 | 路由管理 | 单页应用路由 |
| Axios | 1.12.2 | HTTP客户端 | API请求 |
| Tailwind CSS | 4.1.18 | CSS框架 | 样式开发 |
| Sass | 1.95.0 | CSS预处理器 | 高级样式语法 |

### 容器化技术

| 技术 | 版本 | 用途 | 备注 |
|-----|------|------|------|
| Docker | 最新 | 容器化平台 | 服务隔离、环境一致性 |
| Docker Compose | 最新 | 容器编排 | 多服务管理、一键部署 |
| Nginx | alpine | Web服务器 | 静态资源服务、API反向代理 |

## 项目结构

```
webGL-720yun/
├── backend/                      # 后端代码
│   ├── app/                     # 业务模块
│   │   └── user/                # 用户模块
│   ├── cmd/                     # 入口文件
│   │   └── main.go              # 主程序入口
│   ├── config/                  # 配置文件
│   │   ├── config.dev.yaml      # 开发环境配置
│   │   ├── config.prod.yaml     # 生产环境配置
│   │   └── config.go            # 配置加载逻辑
│   ├── pkg/                     # 公共包
│   │   ├── database/            # 数据库连接
│   │   ├── logger/              # 日志配置
│   │   ├── middleware/          # 中间件
│   │   ├── services/            # 服务层
│   │   └── utils/               # 工具函数
│   ├── route/                   # 路由配置
│   ├── Dockerfile               # 后端Dockerfile
│   ├── go.mod                   # Go模块依赖
│   └── go.sum                   # 依赖版本锁定
├── frontend/                     # 前端代码
│   ├── src/                     # 源码目录
│   │   ├── apis/                # API请求封装
│   │   ├── components/          # Vue组件
│   │   ├── models/              # TypeScript类型定义
│   │   ├── router/              # 路由配置
│   │   ├── stores/              # Pinia状态管理
│   │   ├── utils/               # 工具函数
│   │   ├── App.vue              # 根组件
│   │   ├── main.ts              # 入口文件
│   │   └── style.css            # 全局样式
│   ├── public/                  # 静态资源
│   ├── .gitignore               # Git忽略文件
│   ├── Dockerfile               # 前端Dockerfile
│   ├── nginx.conf               # Nginx配置
│   ├── package.json             # 前端依赖
│   ├── tsconfig.json            # TypeScript配置
│   └── vite.config.ts           # Vite配置
├── docker-compose.yml           # Docker Compose配置
└── README.md                    # 项目文档
```

## 快速开始（使用Docker）

### 前提条件

- 安装 [Docker](https://www.docker.com/get-started)
- 安装 [Docker Compose](https://docs.docker.com/compose/install/)

### 环境分离配置

项目使用多文件Docker Compose配置，实现开发环境和生产环境的分离：

| 配置文件 | 用途 | 特点 |
|---------|------|------|
| `docker-compose.yml` | 基础配置 | 包含所有服务的通用配置 |
| `docker-compose.dev.yml` | 开发环境 | 挂载源码、调试模式、热重载 |
| `docker-compose.prod.yml` | 生产环境 | 标准端口、性能优化、安全性配置 |

### 开发环境启动

```bash
# 克隆项目
git clone <项目仓库地址>
cd webGL-720yun

# 开发环境启动（推荐）
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d

# 开发环境启动并构建（首次运行或修改Dockerfile后）
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d --build

# 开发环境访问地址：
# - 前端：http://localhost:8080
# - 后端API：http://localhost:7000
# - MySQL：localhost:3307（容器内部3306）
# - Redis：localhost:6379
```

### 生产环境启动

```bash
# 生产环境启动
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d

# 生产环境启动并构建
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build

# 生产环境访问地址：
# - 前端：http://localhost 或 https://localhost
# - 后端API：http://localhost/api 或 https://localhost/api
```

### 停止服务

```bash
# 停止开发环境服务
docker-compose -f docker-compose.yml -f docker-compose.dev.yml down

# 停止生产环境服务
docker-compose -f docker-compose.yml -f docker-compose.prod.yml down

# 停止并删除数据卷（谨慎使用）
docker-compose -f docker-compose.yml -f docker-compose.dev.yml down -v
```

### 查看日志

```bash
# 查看开发环境日志
docker-compose -f docker-compose.yml -f docker-compose.dev.yml logs -f

# 查看生产环境日志
docker-compose -f docker-compose.yml -f docker-compose.prod.yml logs -f

# 查看特定服务日志
docker-compose -f docker-compose.yml -f docker-compose.dev.yml logs -f frontend
```

## 开发指南（本地开发）

### 后端开发

#### 前提条件

- 安装 Go 1.25.3 或更高版本
- 安装 MySQL 8.0（可选，也可以使用Docker中的MySQL）
- 安装 Redis 7.0（可选，也可以使用Docker中的Redis）

#### 开发流程

1. **配置环境**
   ```bash
   # 进入后端目录
   cd backend
   
   # 复制配置文件（如果不存在）
   cp config/config.dev.yaml.example config/config.dev.yaml
   
   # 修改配置文件，根据需要调整数据库和Redis连接信息
   # 注意：如果使用Docker中的服务，host应为容器名称（mysql、redis）
   # 如果使用本地服务，host应为localhost
   ```

2. **安装依赖**
   ```bash
   go mod tidy
   ```

3. **运行服务**
   ```bash
   # 运行主程序
   go run cmd/main.go
   
   # 或使用空气（air）实现热重载（推荐）
   # 安装 air：go install github.com/cosmtrek/air@latest
   air
   ```

4. **测试API**
   ```bash
   # 使用curl测试
   curl http://localhost:7000/health
   
   # 或使用Postman等工具
   ```

### 前端开发

#### 前提条件

- 安装 Node.js 18 或更高版本
- 安装 npm 或 yarn 或 pnpm

#### 开发流程

1. **配置环境**
   ```bash
   # 进入前端目录
   cd frontend
   ```

2. **安装依赖**
   ```bash
   npm install
   ```

3. **运行开发服务器**
   ```bash
   npm run dev
   
   # 服务将在 http://localhost:5173 启动
   ```

4. **构建生产版本**
   ```bash
   npm run build
   
   # 构建产物将生成在 dist 目录
   ```

5. **预览生产构建**
   ```bash
   npm run preview
   ```

## 部署说明

### 开发环境部署

使用Docker Compose一键部署，详见「快速开始」章节。

### 生产环境部署

1. **准备生产环境配置**
   - 修改 `backend/config/config.prod.yaml`，配置生产环境的数据库、Redis等信息
   - 修改 `frontend/.env.production`，配置生产环境的API地址

2. **构建生产镜像**
   ```bash
   # 构建后端镜像
docker build -t webgl-720yun-backend:prod ./backend

# 构建前端镜像
docker build -t webgl-720yun-frontend:prod ./frontend
   ```

3. **部署到生产服务器**
   - 使用Docker Compose或Kubernetes部署
   - 配置反向代理（如Nginx）处理HTTPS
   - 配置负载均衡（可选，用于高可用）
   - 配置监控和日志收集

