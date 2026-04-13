# 全景上传与切片架构

## 1. 上传模块时序图

```mermaid
sequenceDiagram
    autonumber
    participant 前端 as 前端
    participant 后端 as 后端
    participant Redis as Redis
    participant MinIO as MinIO

    %% 秒传检查
    前端->>后端: POST /api/v1/upload/init<br/>(filename, file_size, file_hash)
    后端->>Redis: 根据MD5查询file_id
    Redis-->>后端: file_id (如果存在)
    alt 秒传 (文件已存在)
        后端-->>前端: {instant: true, file_id, source_url, thumb_url}
    else 新文件上传
        后端->>Redis: 检查并发上传数量
        alt 并发数达到上限
            后端-->>前端: 错误: 超过同时上传限制
        else 可以上传
            后端->>Redis: 创建上传任务
            后端-->>前端: {upload_id, chunk_size, total_chunks}
        end
    end

    %% 分片上传
    loop 每个分片循环
        前端->>后端: POST /upload/chunk<br/>(upload_id, chunk_index, chunk_hash, chunk_data)
        后端->>后端: 校验分片MD5
        后端->>后端: 保存分片到临时文件
        后端->>Redis: 更新上传进度
        后端-->>前端: {code: 200}
    end

    %% 完成上传
    前端->>后端: POST /api/v1/upload/complete<br/>(upload_id, file_hash)
    后端->>后端: 合并分片为单个文件
    后端->>后端: 校验图片格式和分辨率
    后端->>后端: 生成缩略图
    后端->>MinIO: 上传 source.jpg
    后端->>MinIO: 上传 thumb.jpg
    后端->>Redis: 保存文件信息 (MD5 -> file_id)
    后端->>Redis: 删除上传任务
    后端-->>前端: {file_id, source_url, thumb_url}
```

## 2. 场景创建与切片流水线时序图

```mermaid
sequenceDiagram
    autonumber
    participant 前端 as 前端
    participant 后端 as 后端
    participant Redis as Redis
    participant MinIO as MinIO
    participant 切片Worker as 切片Worker

    %% 创建场景
    前端->>后端: POST /api/v1/resource/scenes<br/>(space_id, title, scene_code, file_id)
    后端->>Redis: 根据file_id获取文件信息
    Redis-->>后端: {source_url, thumb_url, ...}
    后端->>后端: 创建 res_scene 记录
    后端->>Redis: 推送切片任务到 Stream
    后端-->>前端: {scene_id, task_id, slice_status: slicing}

    %% 切片Worker处理
    loop Worker循环 (Redis Stream)
        切片Worker->>Redis: XReadGroup (阻塞)
        Redis-->>切片Worker: SliceTask
        切片Worker->>切片Worker: ProcessSliceTask()
    end

    切片Worker->>MinIO: 下载 source.jpg
    切片Worker->>切片Worker: E2C转换 (6个面)
    切片Worker->>切片Worker: 生成瓦片 (LOD金字塔)
    切片Worker->>MinIO: 上传立方体纹理
    切片Worker->>MinIO: 上传瓦片
    切片Worker->>MinIO: 上传预览图
    切片Worker->>后端: 更新 res_scene<br/>(tile_url, preview_url, slice_status: ready)
    切片Worker->>Redis: XAck 任务确认
    切片Worker-->>前端: WebSocket: slice_complete
```

## 3. 分片上传流程图

```mermaid
flowchart TD
    A[开始: 选择全景图文件] --> B[计算文件MD5]
    B --> C[调用初始化上传API]
    C --> D{文件MD5<br/>是否存在于Redis?}
    D -->|是| E[秒传]
    D -->|否| F{创建上传任务<br/>是否成功?}
    F -->|否| G[错误: 并发数达到上限]
    F -->|是| H[获取分片大小和总分片数]

    E --> I[从Redis获取file_id]
    I --> J[结束: 使用已存在文件]

    H --> J1[将文件拆分为分片]
    J1 --> K{还有分片<br/>需要上传?}
    K -->|是| L[上传第i个分片]
    L --> L1[校验分片MD5]
    L1 --> L2[保存到临时文件]
    L2 --> L3[更新Redis进度]
    L3 --> K
    K -->|否| M[调用完成上传API]

    M --> N{合并分片<br/>是否成功?}
    N -->|否| O[错误: 合并失败]
    N -->|是| P{校验图片<br/>格式和分辨率?}
    P -->|否| Q[错误: 无效图片]
    P -->|是| R[生成缩略图]
    R --> S[上传到MinIO]
    S --> T[保存文件信息到Redis]
    T --> U[删除上传任务]
    U --> V[结束: 上传完成]
```

## 4. 切片流水线流程图

```mermaid
flowchart TD
    A[开始: 从Redis Stream<br/>获取切片任务] --> B[从MinIO下载<br/>source.jpg]
    B --> C{源文件<br/>是否有效?}
    C -->|否| D[错误: 无效源文件]
    C -->|是| E[创建临时目录]
    E --> F[阶段: 下载中<br/>进度: 0-15%]

    F --> G[创建立方体纹理目录]
    G --> H[E2C投影转换]
    H --> I[提取6个立方体面<br/>px, nx, py, ny, pz, nz]
    I --> J[阶段: E2C转换<br/>进度: 15-30%]

    J --> K[创建瓦片目录]
    K --> L{每个LOD级别<br/>0-4}
    L -->|级别 0| M0[生成256x256瓦片]
    L -->|级别 1| M1[生成128x128瓦片]
    L -->|级别 2| M2[生成64x64瓦片]
    L -->|级别 3| M3[生成32x32瓦片]
    L -->|级别 4| M4[生成16x16瓦片]
    M0 --> L
    M1 --> L
    M2 --> L
    M3 --> L
    M4 --> L
    L -->|所有级别完成| N[阶段: 生成瓦片<br/>进度: 30-70%]

    N --> O[从正面生成<br/>预览图片]
    O --> P[阶段: 上传中<br/>进度: 70-90%]

    P --> Q{上传立方体纹理<br/>到MinIO?}
    Q -->|否| R[错误: 上传失败]
    Q -->|是| S[上传瓦片到MinIO]
    S --> T[上传预览图到MinIO]

    T --> U[更新数据库res_scene<br/>tile_url, preview_url<br/>slice_status: ready]
    U --> V[阶段: 完成<br/>进度: 100%]
    V --> W[发送WebSocket<br/>slice_complete通知]
    W --> X[ACK任务确认<br/>从Stream中删除]
    X --> Y[结束: 切片完成]
```

## 5. 系统架构图

```mermaid
flowchart TB
    subgraph 前端
        A1[Vue.js 应用]
        A2[Panolens.js 全景查看器]
        A3[WebSocket客户端]
    end

    subgraph API网关["API网关 / Gin路由"]
        B1[认证中间件]
        B2[跨域中间件]
        B3[限流中间件]
    end

    subgraph 上传模块["上传模块"]
        C1[初始化上传处理器]
        C2[分片上传处理器]
        C3[完成上传处理器]
        C4[上传服务]
        C5[上传仓库]
    end

    subgraph 资源模块["资源模块"]
        D1[空间处理器]
        D2[场景处理器]
        D3[热点处理器]
        D4[场景服务]
    end

    subgraph 切片模块["切片模块"]
        E1[切片队列<br/>Redis Stream]
        E2[切片处理器]
        E3[Worker池]
        E4[E2C转换器]
        E5[瓦片生成器]
    end

    subgraph 存储层
        F1[MinIO<br/>源文件存储]
        F2[MinIO<br/>瓦片存储]
        F3[Redis<br/>上传任务]
        F4[Redis<br/>文件信息]
        F5[Redis<br/>切片流]
        F6[MySQL<br/>场景数据]
    end

    A1 -->|HTTP请求| B1
    A1 -->|WebSocket| A3
    A3 -->|切片进度消息| E3

    B1 --> B2 --> C1
    B1 --> B2 --> C2
    B1 --> B2 --> C3
    B1 --> B2 --> D1
    B1 --> B2 --> D2
    B1 --> B2 --> D3

    C1 --> C4 --> C5
    C2 --> C4 --> C5
    C3 --> C4 --> C5

    D2 --> D4
    D4 -->|推送切片任务| E1

    C5 --> F3
    C5 --> F4
    C5 --> F1

    E1 --> E3
    E3 --> E2
    E2 --> E4
    E2 --> E5

    E4 --> F1
    E5 --> F2
    E2 --> F6
    E2 --> A3

    F1 --> MinIO_S3[(MinIO对象存储)]
    F2 --> MinIO_S3
    F3 --> Redis[(Redis)]
    F4 --> Redis
    F5 --> Redis
    F6 --> MySQL[(MySQL数据库)]
```

## 6. 数据流架构图

```mermaid
flowchart LR
    subgraph 上传流程["1. 上传流程"]
        A1[选择文件] --> A2[计算MD5]
        A2 --> A3[初始化API]
        A3 --> A4{MD5<br/>是否存在?}
        A4 -->|是| A5[返回<br/>file_id]
        A4 -->|否| A6[拆分<br/>分片]
        A6 --> A7[上传<br/>分片]
        A7 --> A8[完成<br/>上传API]
        A8 --> A9[合并并<br/>校验]
        A9 --> A10[上传到<br/>MinIO]
        A10 --> A11[保存到<br/>Redis]
    end

    subgraph 场景流程["2. 场景创建流程"]
        B1[创建<br/>场景API] --> B2[从Redis<br/>获取文件信息]
        B2 --> B3[保存到<br/>MySQL]
        B3 --> B4[推送到<br/>Redis Stream]
    end

    subgraph 切片流程["3. 切片流水线"]
        C1[Worker<br/>取出任务] --> C2[下载<br/>源文件]
        C2 --> C3[E2C<br/>转换]
        C3 --> C4[生成<br/>瓦片]
        C4 --> C5[上传到<br/>MinIO]
        C5 --> C6[更新<br/>数据库状态]
        C6 --> C7[WebSocket<br/>通知]
    end

    subgraph 查看器流程["4. 查看器加载流程"]
        D1[从API加载<br/>场景信息] --> D2[获取瓦片URL]
        D2 --> D3[从MinIO<br/>加载瓦片]
        D3 --> D4[渲染<br/>全景图]
    end

    上传流程 -->|file_id| 场景流程
    场景流程 -->|task| 切片流程
    切片流程 -->|ready| 查看器流程
```

## 7. Redis数据结构

```mermaid
erDiagram
    Redis上传任务 ||--o| Redis文件信息 : "映射关系"
    Redis文件信息 {
        string file_id 主键
        string source_url 源文件URL
        string thumb_url 缩略图URL
        int file_size 文件大小
        int width 图片宽度
        int height 图片高度
        timestamp created_at 创建时间
    }
    Redis上传任务 {
        string upload_id 主键
        string file_hash 文件MD5
        int total_chunks 总分片数
        int uploaded_chunks 已上传分片数
        string status 状态
        timestamp created_at 创建时间
    }
    RedisMD5索引 {
        string file_md5 主键
        string file_id 外键
    }
    Redis切片流 {
        string stream_key 流键
        string task_id 任务ID
        uint scene_id 场景ID
        string file_id 文件ID
        uint user_id 用户ID
        timestamp created_at 创建时间
    }
```

## 8. 场景切片状态机

```mermaid
stateDiagram-v2
    [*] --> pending: 创建场景
    pending --> slicing: 推送任务到Redis

    slicing --> downloading: Worker获取任务
    downloading --> e2c: 下载完成
    e2c --> tiles: E2C转换完成
    tiles --> uploading: 瓦片生成完成
    uploading --> ready: 上传完成

    slicing --> failed: 错误发生
    downloading --> failed: 错误发生
    e2c --> failed: 错误发生
    tiles --> failed: 错误发生
    uploading --> failed: 错误发生

    ready --> [*]
    failed --> [*]
```

## 9. E2C转换原理图

```mermaid
flowchart LR
    subgraph 等距圆柱投影["等距圆柱投影 (ERP)"]
        A1[2:1全景图<br/>4096x2048]
    end

    subgraph 球面坐标["球面坐标映射"]
        B1[经度: -π to +π]
        B2[纬度: -π/2 to +π/2]
    end

    subgraph 立方体面["立方体六个面"]
        C1[PX<br/>右面 +X]
        C2[NX<br/>左面 -X]
        C3[PY<br/>上面 +Y]
        C4[NY<br/>下面 -Y]
        C5[PZ<br/>前面 +Z]
        C6[NZ<br/>后面 -Z]
    end

    A1 --> B1
    B1 --> B2
    B2 --> C1
    B2 --> C2
    B2 --> C3
    B2 --> C4
    B2 --> C5
    B2 --> C6

    style C1 fill:#ff9999
    style C2 fill:#99ff99
    style C3 fill:#9999ff
    style C4 fill:#ffff99
    style C5 fill:#ff99ff
    style C6 fill:#99ffff
```

## 10. LOD瓦片金字塔

```mermaid
flowchart TB
    subgraph LOD0["LOD 0 - 最高清晰度 (256x256)"]
        L0_1[0,0]
        L0_2[1,0]
        L0_3[0,1]
        L0_4[1,1]
    end

    subgraph LOD1["LOD 1 - 较低清晰度 (128x128)"]
        L1_1[0,0]
        L1_2[1,0]
    end

    subgraph LOD2["LOD 2 - 低清晰度 (64x64)"]
        L2_1[0,0]
        L2_2[1,0]
    end

    subgraph LOD3["LOD 3 - 很低清晰度 (32x32)"]
        L3_1[0,0]
        L3_2[1,0]
    end

    subgraph LOD4["LOD 4 - 最低清晰度 (16x16)"]
        L4_1[0,0]
        L4_2[1,0]
    end

    L0_1 --> L1_1
    L0_2 --> L1_1
    L0_3 --> L1_2
    L0_4 --> L1_2
    L1_1 --> L2_1
    L1_2 --> L2_2
    L2_1 --> L3_1
    L2_2 --> L3_2
    L3_1 --> L4_1
    L3_2 --> L4_2
```
