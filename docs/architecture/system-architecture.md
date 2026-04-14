# 全景系统架构文档

## 目录

1. [上传模块时序图](#1-上传模块时序图)
2. [场景创建与切片流水线时序图](#2-场景创建与切片流水线时序图)
3. [分片上传流程图](#3-分片上传流程图)
4. [切片流水线流程图](#4-切片流水线流程图)
5. [系统架构图](#5-系统架构图)
6. [数据流架构图](#6-数据流架构图)
7. [Redis数据结构](#7-redis数据结构)
8. [场景切片状态机](#8-场景切片状态机)
9. [E2C转换原理图](#9-e2c转换原理图)
10. [LOD瓦片金字塔](#10-lod瓦片金字塔)

---

## 1. 上传模块时序图

```mermaid
sequenceDiagram
    participant Client as 前端客户端
    participant API as Gin API Server
    participant Redis as Redis
    participant MinIO as MinIO Storage

    Client->>API: POST /api/upload/init<br/>(FileMD5, FileSize, SpaceID)
    API->>Redis: 检查用户并发上传数
    API->>Redis: 检查MD5是否存在
    alt 文件已存在
        API-->>Client: 返回已有文件信息<br/>(Instant: true, FileID)
    else 文件不存在
        API->>Redis: 创建上传任务
        API->>Redis: 增加用户并发计数
        API-->>Client: 返回UploadID, ChunkSize, TotalChunks
    end

    loop 分片上传
        Client->>API: POST /api/upload/chunk<br/>(UploadID, ChunkIndex, ChunkData)
        API->>Redis: 检查分片是否已上传
        alt 分片未上传
            API->>MinIO: 上传分片到 temp/{UploadID}/chunk_{Index}
            API->>Redis: 标记分片已上传
            API->>Redis: 更新已上传字节数
        end
        API-->>Client: 返回UploadedChunks列表
    end

    Client->>API: POST /api/upload/complete<br/>(UploadID)
    API->>Redis: 获取已上传分片列表
    alt 分片完整
        API->>MinIO: 按顺序合并所有分片
        API->>API: 验证文件格式和分辨率
        API->>MinIO: 生成缩略图
        API->>MinIO: 上传源文件到 spaces/{SpaceName}/sources/{FileID}/
        API->>MinIO: 上传缩略图到 spaces/{SpaceName}/previews/{FileID}/
        API->>Redis: 保存MD5映射
        API->>Redis: 保存文件信息
        API->>Redis: 清理临时分片
        API->>Redis: 减少用户并发计数
        API-->>Client: 返回FileID, SourceURL, ThumbURL
    else 分片不完整
        API-->>Client: 返回错误 (已上传X/Y分片)
    end
```

---

## 2. 场景创建与切片流水线时序图

```mermaid
sequenceDiagram
    participant Client as 前端客户端
    participant API as Gin API Server
    participant DB as MySQL
    participant Queue as Redis Stream
    participant Worker as Slice Worker
    participant MinIO as MinIO Storage
    participant WS as WebSocket Hub

    Client->>API: POST /api/scenes<br/>(SpaceID, Title, FileID)
    API->>DB: 创建场景记录 (Status: pending)
    API->>Queue: 推送切片任务<br/>(SliceTask: SceneID, FileID, SpaceName)
    API-->>Client: 返回SceneID, TaskID

    Worker->>Queue: XReadGroup 阻塞获取任务
    Worker->>MinIO: 下载源文件 from spaces/{SpaceName}/sources/{FileID}/

    par 并行处理
        Worker->>Worker: 生成快速预览 (1024x512)
        Worker->>MinIO: 上传预览图
        Worker->>DB: 更新Scene: PreviewURL, SliceStatus=slicing
        Worker->>WS: 通知客户端 进度25%
    end

    par 并行E2C转换 (6个面)
        Worker->>Worker: 提取正面 PX (Right)
        Worker->>Worker: 提取背面 NX (Left)
        Worker->>Worker: 提取顶面 PY (Top)
        Worker->>Worker: 提取底面 NY (Bottom)
        Worker->>Worker: 提取右面 PZ (Front)
        Worker->>Worker: 提取左面 NZ (Back)
    end
    Worker->>WS: 通知客户端 进度45%

    par 并行瓦片生成 (每个面)
        Worker->>Worker: Level 0: 512x512 -> 4个256x256瓦片
        Worker->>Worker: Level 1: 256x256 -> 4个瓦片
        Worker->>Worker: Level 2: 128x128 -> 4个瓦片
        Worker->>Worker: ... 继续直到 < 256
    end
    Worker->>WS: 通知客户端 进度70%

    par 并行上传瓦片 (最多10并发)
        Worker->>MinIO: 上传cubemap面图
        Worker->>MinIO: 上传所有Level瓦片
    end
    Worker->>WS: 通知客户端 进度95%

    Worker->>DB: 更新Scene: SliceStatus=ready, TileURL, IsConverted=true
    Worker->>Queue: XAck 确认任务完成
    Worker->>WS: 通知客户端 切片完成
    WS-->>Client: 接收完成通知<br/>(TileURL, PreviewURL)
```

---

## 3. 分片上传流程图

```mermaid
flowchart TD
    A[开始上传] --> B[前端计算文件MD5]
    B --> C[调用InitUpload接口]
    C --> D{检查MD5是否存在?}

    D -->|是| E[秒传成功<br/>返回已有文件信息]
    D -->|否| F[创建上传任务<br/>返回UploadID]

    F --> G[将文件分片<br/>每片5MB]
    G --> H{还有未上传分片?}

    H -->|是| I[上传分片到MinIO<br/>temp/{UploadID}/chunk_{Index}]
    I --> J{上传成功?}

    J -->|否| K[重试最多3次]
    K --> I
    J -->|是| L[更新Redis分片状态]
    L --> H

    H -->|否| M[调用CompleteUpload]
    M --> N[按顺序合并分片]
    N --> O[验证文件格式和分辨率]
    O --> P{验证通过?}

    P -->|否| Q[返回错误<br/>清理临时文件]
    P -->|是| R[生成缩略图]
    R --> S[上传源文件到正式目录]
    S --> T[上传缩略图]
    T --> U[保存MD5映射到Redis]
    U --> V[清理临时分片]
    V --> W([上传完成])

    E --> W
```

---

## 4. 切片流水线流程图

```mermaid
flowchart TD
    A([开始切片任务]) --> B[下载源文件]
    B --> C[生成快速预览<br/>1024x512]
    C --> D[上传预览到MinIO]
    D --> E[更新数据库<br/>PreviewURL]
    E --> F[通知客户端<br/>进度25%]

    F --> G[创建6个Cubemap目录]

    G --> H{并行处理6个面}

    H -->|正面 PX| I1[E2C转换<br/>提取正面像素]
    H -->|背面 NX| I2[E2C转换<br/>提取背面像素]
    H -->|顶面 PY| I3[E2C转换<br/>提取顶面像素]
    H -->|底面 NY| I4[E2C转换<br/>提取底面像素]
    H -->|右面 PZ| I5[E2C转换<br/>提取右面像素]
    H -->|左面 NZ| I6[E2C转换<br/>提取左面像素]

    I1 --> J[等待所有面完成]
    I2 --> J
    I3 --> J
    I4 --> J
    I5 --> J
    I6 --> J

    J --> K[通知客户端<br/>进度45%]

    K --> L[每个面生成LOD瓦片]

    subgraph LOD生成 [每个面独立生成]
        L --> M1[Level 0: 原尺寸]
        M1 --> M2[Level 1: 1/2尺寸]
        M2 --> M3[Level 2: 1/4尺寸]
        M3 --> M4[Level 3: 1/8尺寸]
        M4 --> M5[...]
        M5 --> M6[直到 < 256px]
    end

    M6 --> N[收集所有瓦片文件]
    N --> O[通知客户端<br/>进度70%]

    O --> P[并发上传瓦片<br/>最大10并发]

    subgraph 上传阶段
        P --> P1[上传Cubemap面图]
        P --> P2[上传Level 0 瓦片]
        P --> P3[上传Level 1 瓦片]
        P --> P4[上传Level 2 瓦片]
        P --> P5[上传Level 3 瓦片]
    end

    P1 --> Q[等待所有上传完成]
    P2 --> Q
    P3 --> Q
    P4 --> Q
    P5 --> Q

    Q --> R[更新数据库<br/>SliceStatus=ready]
    R --> S[通知客户端<br/>进度100%]
    S --> T([任务完成])

    style H fill:#f9f,stroke:#333
    style L fill:#ff9,stroke:#333
    style P fill:#9f9,stroke:#333
```

---

## 5. 系统架构图

```mermaid
graph TB
    subgraph Frontend ["前端 (Vue 3 + TypeScript)"]
        UI["静态资源服务"]
        Viewer["Pannellum Viewer<br/>全景浏览组件"]
        Uploader["分片上传组件"]
        Editor["场景编辑器"]
    end

    subgraph Backend ["后端 (Go + Gin)"]
        API["Gin HTTP Server<br/>:8080"]

        subgraph Core ["核心模块"]
            Auth["认证鉴权<br/>JWT + RBAC"]
            Router["路由分发"]
        end

        subgraph Resource ["资源管理"]
            SpaceAPI["空间管理 API"]
            SceneAPI["场景管理 API"]
            HotspotAPI["热点管理 API"]
            UploadAPI["上传管理 API"]
        end

        subgraph Upload ["分片上传"]
            UploadSVC["UploadService<br/>分片管理"]
            UploadRepo["UploadRepository<br/>Redis操作"]
        end

        subgraph Slice ["切片流水线"]
            Queue["SliceQueue<br/>Redis Stream"]
            Processor["SliceProcessor<br/>核心处理器"]
            Worker["WorkerPool<br/>Goroutine池"]
        end

        subgraph Panorana ["全景处理"]
            E2C["E2C Converter<br/>等距圆柱转立方体"]
            Tiler["Tile Generator<br/>LOD瓦片生成"]
        end

        subgraph WebSocket ["实时通信"]
            Hub["WebSocket Hub<br/>消息广播"]
        end
    end

    subgraph Storage ["存储层"]
        Redis["Redis<br/>任务队列 + 缓存"]
        MinIO["MinIO S3<br/>对象存储"]
        MySQL["MySQL<br/>关系数据"]
    end

    UI --> |HTTP/REST| API
    Viewer --> |WebSocket| Hub
    Uploader --> |分片上传| UploadAPI

    API --> Auth
    API --> Router
    Router --> SpaceAPI
    Router --> SceneAPI
    Router --> HotspotAPI
    Router --> UploadAPI

    UploadAPI --> UploadSVC
    UploadSVC --> UploadRepo
    UploadRepo --> Redis
    UploadSVC --> MinIO

    SceneAPI --> Processor
    Processor --> Queue
    Queue --> Worker
    Worker --> Processor

    Processor --> E2C
    Processor --> Tiler
    Processor --> MinIO
    Processor --> Hub
    Processor --> MySQL

    E2C --> |并行处理| Worker
    Tiler --> |并行处理| Worker

    Redis --> |Stream| Queue
    Redis --> |Set| UploadRepo
```

---

## 6. 数据流架构图

```mermaid
flowchart LR
    subgraph Upload ["上传数据流"]
        direction TB
        A1[前端分片] --> A2[InitUpload]
        A2 --> A3[Chunk Upload]
        A3 --> A4[CompleteUpload]
        A4 --> A5[MinIO存储]
    end

    subgraph Slice ["切片数据流"]
        direction TB
        B1[源文件] --> B2[下载到本地]
        B2 --> B3[E2C转换]
        B3 --> B4[Cubemap 6面]
        B4 --> B5[LOD瓦片生成]
        B5 --> B6[上传到MinIO]
    end

    subgraph Serve ["服务数据流"]
        direction TB
        C1[客户端请求] --> C2[API Gateway]
        C2 --> C3[获取Scene信息]
        C3 --> C4[返回TileURL]
        C4 --> C5[前端加载瓦片]
    end

    subgraph Notify ["通知数据流"]
        direction TB
        D1[切片进度] --> D2[WebSocket]
        D2 --> D3[实时推送]
        D3 --> D4[前端更新]
    end

    style Upload fill:#e1f5fe
    style Slice fill:#fff3e0
    style Serve fill:#e8f5e9
    style Notify fill:#fce4ec
```

```mermaid
flowchart TD
    subgraph Client ["客户端"]
        Browser["浏览器"]
        Mobile["移动端"]
    end

    subgraph CDN ["CDN / 静态资源"]
        Static["Nginx<br/>静态资源"]
    end

    subgraph API ["API Layer"]
        Gateway["API Gateway<br/>负载均衡"]
    end

    subgraph UploadFlow ["上传流程"]
        Init["InitUpload"]
        Chunk["UploadChunk"]
        Complete["CompleteUpload"]
    end

    subgraph SliceFlow ["切片流程"]
        Download["Download Source"]
        Preview["Generate Preview"]
        E2C["E2C Convert"]
        Tile["Generate Tiles"]
        Upload["Upload Tiles"]
    end

    subgraph Storage ["存储"]
        TempBucket["temp/<br/>分片暂存"]
        SourceBucket["spaces/sources/<br/>源文件"]
        PreviewBucket["spaces/previews/<br/>预览图"]
        TileBucket["spaces/tiles/<br/>瓦片"]
    end

    subgraph Cache ["缓存层"]
        UploadTask["upload:task:{ID}<br/>任务信息"]
        UploadChunks["upload:chunks:{ID}<br/>分片Set"]
        UserUploadCount["upload:user:{ID}:count<br/>用户计数"]
        FileMD5["file:md5:{MD5}<br/>MD5映射"]
    end

    subgraph Queue ["消息队列"]
        SliceStream["slice:tasks<br/>Redis Stream"]
    end

    Browser --> |上传| Gateway
    Gateway --> Init
    Gateway --> Chunk
    Gateway --> Complete

    Init --> UploadTask
    Chunk --> TempBucket
    Chunk --> UploadChunks
    Complete --> SourceBucket
    Complete --> PreviewBucket
    Complete --> FileMD5

    SourceBucket --> Download
    Download --> Preview
    Download --> E2C
    E2C --> Tile
    Tile --> Upload

    Upload --> TileBucket

    SliceStream --> |消费者| Preview
    SliceStream --> |消费者| E2C
    SliceStream --> |消费者| Tile
    SliceStream --> |消费者| Upload

    style UploadFlow fill:#b3e5fc
    style SliceFlow fill:#ffccbc
    style Storage fill:#c8e6c9
    style Cache fill:#fff9c4
    style Queue fill:#e1bee7
```

---

## 7. Redis数据结构

```mermaid
erDiagram
    REDIS_DATASETS ||--|| UPLOAD_TASK : stores
    REDIS_DATASETS ||--|| UPLOAD_CHUNKS : stores
    REDIS_DATASETS ||--|| USER_UPLOAD_COUNT : tracks
    REDIS_DATASETS ||--|| FILE_MD5 : maps
    REDIS_DATASETS ||--|| SLICE_STREAM : manages

    UPLOAD_TASK {
        string key "upload:task:{UploadID}"
        hash value "TaskID, UserID, SpaceID, FileName,<br/>FileSize, FileMD5, TotalChunks,<br/>ChunkSize, Status, UploadedBytes"
        int ttl "48小时"
    }

    UPLOAD_CHUNKS {
        string key "upload:chunks:{UploadID}"
        set value "ChunkIndex列表"
        int ttl "48小时"
    }

    USER_UPLOAD_COUNT {
        string key "upload:user:{UserID}:count"
        int value "当前并发上传数"
        int ttl "无"
    }

    FILE_MD5 {
        string key "file:md5:{MD5}"
        string value "FileID"
        int ttl "永不过期"
    }

    FILE_INFO {
        string key "file:info:{FileID}"
        hash value "FileID, SpaceID, SpaceName,<br/>SourceURL, ThumbURL, FileSize,<br/>Width, Height, CreatedAt"
        int ttl "30天"
    }

    SLICE_STREAM {
        stream key "slice:tasks"
        entry field "data: JSON序列化任务"
        consumer_group "slice_workers"
    }
```

```mermaid
mindmap
    root((Redis数据结构))
        上传任务
            upload:task:{UploadID}
            Hash存储
                TaskID
                UserID
                SpaceID
                FileName
                FileSize
                TotalChunks
                Status
            TTL: 48小时
        分片状态
            upload:chunks:{UploadID}
            Set存储
            每个元素是分片索引
            TTL: 48小时
        用户上传计数
            upload:user:{UserID}:count
            String整数
            INCR/DECR操作
        MD5映射
            file:md5:{MD5}
            String存储FileID
            秒传依据
        切片任务流
            slice:tasks
            Stream类型
            XADD添加任务
            XReadGroup消费
            XAck确认完成
        用户会话
            user:session:{UserID}
            Hash存储
            access_token
            refresh_token
        在线状态
            user:online:{UserID}
            String "1"
            TTL管理
```

---

## 8. 场景切片状态机

```mermaid
stateDiagram-v2
    [*] --> Pending : 创建场景

    Pending --> Slicing : 开始切片
    Slicing --> Ready : 切片完成
    Slicing --> Failed : 切片失败
    Failed --> Slicing : 重试

    Ready --> [*] : 删除场景

    note right of Pending
        初始状态
        等待Worker接收任务
    end

    note right of Slicing
        处理中状态
        包含多个阶段:
        - 下载源文件 (5%)
        - 生成预览 (15%)
        - E2C转换 (30%)
        - 瓦片生成 (50%)
        - 上传瓦片 (75%)
        - 完成 (95%)
    end

    note right of Ready
        切片完成
        可供前端加载
        TileURL已设置
        IsConverted=true
    end

    note right of Failed
        切片失败
        可重试3次
        保存错误信息
    end
```

---

## 9. E2C转换原理图

```mermaid
flowchart LR
    subgraph Equirectangular ["等距圆柱投影 (EQR)<br/>8192 x 4096"]
        direction TB
        E[全景球体<br/>ERP展开图]
    end

    subgraph CubeFaces ["立方体6面展开"]
        direction LR
        PX["+X 正面<br/>Right"]
        NX["-X 背面<br/>Left"]
        PY["+Y 顶面<br/>Top"]
        NY["-Y 底面<br/>Bottom"]
        PZ["+Z 右面<br/>Front"]
        NZ["-Z 左面<br/>Back"]
    end

    subgraph Cube3D ["立方体视图"]
        direction LR
        C3D["立方体贴图<br/>Cubemap"]
    end

    E --> |"球面坐标<br/>(θ, φ) → 面像素"| PX
    E --> |"球面坐标<br/>(θ, φ) → 面像素"| NX
    E --> |"球面坐标<br/>(θ, φ) → 面像素"| PY
    E --> |"球面坐标<br/>(θ, φ) → 面像素"| NY
    E --> |"球面坐标<br/>(θ, φ) → 面像素"| PZ
    E --> |"球面坐标<br/>(θ, φ) → 面像素"| NZ

    PX & NX & PY & NY & PZ & NZ --> C3D

    style Equirectangular fill:#e3f2fd
    style CubeFaces fill:#fff3e0
    style Cube3D fill:#e8f5e9
```

```mermaid
graph LR
    subgraph Input ["输入: ERP全景图"]
        A["球面上的点 P"]
        B["球坐标 (θ, φ)"]
        C["θ: 方位角 -180° ~ 180°"]
        D["φ: 仰角 -90° ~ 90°"]
    end

    subgraph Transform ["E2C坐标变换"]
        E["计算立方体面"]
        F["计算面内坐标"]
        G["(u, v) ∈ [-1, 1]"]
    end

    subgraph Output ["输出: Cubemap面"]
        H["面索引: PX|NX|PY|NY|PZ|NZ"]
        I["像素位置 (x, y)"]
        J["面尺寸: N x N"]
    end

    A --> B
    B --> C
    C --> D
    D --> E
    E --> F
    F --> G
    G --> H
    H --> I
    I --> J

    style Input fill:#bbdefb
    style Transform fill:#c8e6c9
    style Output fill:#ffe0b2
```

```mermaid
graph TD
    subgraph 面索引计算
        A["输入点 P(θ, φ)"] --> B{cos(φ) ≥ sin(θ)?}
        B -->|是 且 cos(φ) ≥ -sin(θ)| PX["PX (+X)<br/>u = 1/z<br/>v = y/|z|"]
        B -->|是 且 cos(φ) < -sin(θ)| NX["NX (-X)<br/>u = -1/z<br/>v = y/|z|"]
        B -->|否| C{sin(φ) ≥ 0?}
        C -->|是| PY["PY (+Y)<br/>u = x/|y|<br/>v = -1/z"]
        C -->|否| NY["NY (-Y)<br/>u = x/|y|<br/>v = 1/z"]
        D{"cos(φ) ≥ cos(θ)?"} -->|是| PZ["PZ (+Z)<br/>u = x/|z|<br/>v = y/|z|"]
        D -->|否| NZ["NZ (-Z)<br/>u = -x/|z|<br/>v = y/|z|"]
    end

    style PX fill:#ffcdd2
    style NX fill:#f8bbd0
    style PY fill:#fff9c4
    style NY fill:#ffe082
    style PZ fill:#c8e6c9
    style NZ fill:#b2dfdb
```

---

## 10. LOD瓦片金字塔

```mermaid
graph TB
    subgraph LOD0 ["Level 0 - 原始分辨率"]
        L0_1["512x512<br/>Cubemap面"]
        L0_2["┌────┬────┐"]
        L0_3["│tile│tile│"]
        L0_4["├────┼────┤"]
        L0_5["│tile│tile│"]
        L0_6["└────┴────┘"]
        L0_7["4 tiles (256x256)"]
    end

    subgraph LOD1 ["Level 1 - 1/2分辨率"]
        L1_1["256x256"]
        L1_2["┌────┬────┐"]
        L1_3["│tile│tile│"]
        L1_4["├────┼────┤"]
        L1_5["│tile│tile│"]
        L1_6["└────┴────┘"]
        L1_7["4 tiles (128x128)"]
    end

    subgraph LOD2 ["Level 2 - 1/4分辨率"]
        L2_1["128x128"]
        L2_2["┌────┬────┐"]
        L2_3["│tile│tile│"]
        L2_4["├────┼────┤"]
        L2_5["│tile│tile│"]
        L2_6["└────┴────┘"]
        L2_7["4 tiles (64x64)"]
    end

    subgraph LOD3 ["Level 3 - 1/8分辨率"]
        L3_1["64x64"]
        L3_2["┌────┬────┐"]
        L3_3["│tile│tile│"]
        L3_4["├────┼────┤"]
        L3_5["│tile│tile│"]
        L3_6["└────┴────┘"]
        L3_7["4 tiles (32x32)"]
    end

    L0_1 --> L1_1
    L1_1 --> L2_1
    L2_1 --> L3_1

    style LOD0 fill:#e3f2fd
    style LOD1 fill:#e8f5e9
    style LOD2 fill:#fff3e0
    style LOD3 fill:#fce4ec
```

```mermaid
mindmap
    root((LOD瓦片<br/>金字塔))
        存储结构
            MinIO路径
            spaces/{SpaceName}/tiles/{SceneCode}/cubemap/{face}.jpg
            spaces/{SpaceName}/tiles/{SceneCode}/lod/{face}/level_{N}/tile_{row}_{col}.jpg
        瓦片规格
            标准瓦片大小
            256 x 256 像素
            JPEG格式
            质量85%
        LOD层级
            Level 0: 512x512 (2x2)
            Level 1: 256x256 (2x2)
            Level 2: 128x128 (2x2)
            Level 3: 64x64 (2x2)
            Level 4: 32x32 (2x2)
            直到 < 256px
        加载策略
            远距离 → 低层级
            近距离 → 高层级
            渐进式加载
            视锥体剔除
        前端渲染
            Pannellum支持
            WebGL渲染
            球面坐标映射
```

```mermaid
graph LR
    subgraph 加载流程
        A["用户缩放"] --> B{计算所需Level}
        B --> C{获取视野范围}
        C --> D[请求可见瓦片]
        D --> E{检查缓存}
        E -->|命中| F[直接使用]
        E -->|未命中| G[加载瓦片]
        G --> H[放入缓存]
        H --> F
        F --> I[渲染显示]
    end

    subgraph 瓦片命名
        J["Path: lod/px/level_2/tile_1_3.jpg"]
        K["px: 正面"]
        L["level_2: 第2级"]
        M["tile_1_3.jpg: 第1行第3列"]
    end

    style 加载流程 fill:#e8f5e9
    style 瓦片命名 fill:#fff3e0
```

---

## 附录: 相关文件索引

| 模块 | 文件路径 | 说明 |
|-----|---------|-----|
| 分片上传 | [upload_service.go](file:///d:/lwdd/code/毕设/webGL-720yun/backend/internal/resource/upload/upload_service.go) | 分片上传核心服务 |
| 分片仓库 | [upload_repository.go](file:///d:/lwdd/code/毕设/webGL-720yun/backend/internal/resource/upload/upload_repository.go) | Redis操作封装 |
| 切片队列 | [queue.go](file:///d:/lwdd/code/毕设/webGL-720yun/backend/internal/slice/queue.go) | Redis Stream任务队列 |
| 切片处理器 | [processor.go](file:///d:/lwdd/code/毕设/webGL-720yun/backend/internal/slice/processor.go) | 切片流水线核心 |
| E2C转换 | [e2c.go](file:///d:/lwdd/code/毕设/webGL-720yun/backend/pkg/panorama/e2c.go) | 等距圆柱转立方体 |
| 立方体转换器 | [converter.go](file:///d:/lwdd/code/毕设/webGL-720yun/backend/pkg/panorama/converter.go) | 6面并行转换 |
| Redis服务 | [redis.go](file:///d:/lwdd/code/毕设/webGL-720yun/backend/pkg/redis/redis.go) | Redis操作封装 |
| 场景模型 | [resource.go](file:///d:/lwdd/code/毕设/webGL-720yun/backend/internal/model/resource.go) | 数据模型定义 |
