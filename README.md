# WebGL-720yun 高性能全景展示平台

## 🚀 项目简介

WebGL-720yun 是一个基于 WebGL 技术构建的高性能 720° 全景图展示与管理平台。项目旨在提供从全景图上传、自动化切片处理（E2C & LOD）到高性能 Web 端展示的完整解决方案。

## 🛠 技术栈

### 后端 (Go)

- **框架**: [Gin](https://github.com/gin-gonic/gin) (Web 路由)
- **ORM**: [GORM](https://gorm.io/) (数据库操作)
- **存储**: [MinIO](https://min.io/) (对象存储), [Redis](https://redis.io/) (缓存与任务队列)
- **认证**: [JWT-Go](https://github.com/golang-jwt/jwt)
- **图像处理**: `imaging` & 自研 E2C 转换库

### 前端 (Vue 3)

- **框架**: [Vue 3](https://vuejs.org/) + [TypeScript](https://www.typescriptlang.org/)
- **构建**: [Vite](https://vitejs.dev/)
- **全景渲染**: [Three.js](https://threejs.org/) & [PhotoSphereViewer](https://photo-sphere-viewer.js.org/)
- **地图**: [@amap/amap-jsapi-loader](https://lbs.amap.com/), [@antv/l7](https://l7.antv.vision/)
- **样式**: [Tailwind CSS](https://tailwindcss.com/)
- **状态**: [Pinia](https://pinia.vuejs.org/)

## 📂 项目结构

```text
webGL-720yun/
├── backend/                  # 后端服务
│   ├── cmd/                 # 核心启动逻辑
│   ├── config/              # 配置文件与结构体
│   ├── internal/            # 内部业务逻辑
│   │   ├── core/            # 核心设施 (路由、中间件、初始化)
│   │   ├── model/           # GORM 模型定义
│   │   ├── resource/        # 空间、场景、热点、上传业务
│   │   ├── slice/           # 切片流水线 (Queue, Worker, Processor)
│   │   └── user/            # 用户认证与管理
│   ├── pkg/                 # 通用工具包 (E2C 算法、MinIO、JWT 等)
│   └── main.go              # 程序入口
├── frontend/                 # 前端应用
│   ├── src/
│   │   ├── components/      # UI 组件 (全景播放器、地图、场景列表)
│   │   ├── stores/          # 状态管理
│   │   ├── services/        # API 请求与 WebSocket
│   │   ├── views/           # 页面视图
│   │   └── utils/           # 地图加载、坐标转换等工具
│   └── vite.config.ts       # 构建配置
├── docker-compose.yml        # 基础设施编排 (MySQL, Redis, MinIO)
└── docs/                     # 架构文档与优化报告
```

## ⚙️ 快速开始

### 1. 基础设施启动

确保已安装 Docker 和 Docker Compose，然后在根目录运行：

```bash
docker-compose up -d
```

这将启动 MySQL, Redis 和 MinIO 容器。

### 2. 后端开发

```bash
cd backend
go mod tidy
# 根据 config.dev.yaml.example 创建并配置 config.dev.yaml
go run main.go
```

### 3. 前端开发

```bash
cd frontend
npm install
npm run dev
```

## 📖 架构图

项目详细的设计文档、时序图及数据流图请参考：

- [系统架构详述](./docs/architecture/system-architecture.md)
- [性能优化报告](./docs/architecture/performance-optimization-report.md)

