# 分布式爬虫系统架构图

## 系统架构概览

```mermaid
graph TB
    subgraph "外部依赖"
        A[etcd集群<br/>服务注册发现]
        B[MySQL数据库<br/>数据存储]
        C[代理服务器池<br/>反爬虫]
    end

    subgraph "Master集群"
        M1[Master节点1]
        M2[Master节点2]
        M3[Master节点3]

        subgraph "Master核心功能"
            MF1[Leader选举]
            MF2[任务调度]
            MF3[负载均衡]
            MF4[节点监控]
        end
    end

    subgraph "Worker集群"
        W1[Worker节点1<br/>爬虫引擎]
        W2[Worker节点2<br/>爬虫引擎]
        W3[Worker节点3<br/>爬虫引擎]
        W4[Worker节点N...]

        subgraph "Worker核心组件"
            WF1[调度器<br/>Schedule]
            WF2[采集器<br/>Fetcher]
            WF3[解析器<br/>Parser]
            WF4[存储器<br/>Storage]
        end
    end

    subgraph "配置与监控"
        CFG[配置文件<br/>config.toml]
        LOG[日志系统<br/>zap]
        MON[性能监控<br/>pprof]
    end

    %% 连接关系
    A -.-> M1
    A -.-> M2
    A -.-> M3
    A -.-> W1
    A -.-> W2
    A -.-> W3
    A -.-> W4

    M1 --> W1
    M1 --> W2
    M1 --> W3
    M1 --> W4

    W1 --> B
    W2 --> B
    W3 --> B
    W4 --> B

    W1 -.-> C
    W2 -.-> C
    W3 -.-> C
    W4 -.-> C

    CFG --> M1
    CFG --> W1
    LOG --> M1
    LOG --> W1
    MON --> M1
    MON --> W1

    %% 选举关系
    M1 -.-> MF1
    M2 -.-> MF1
    M3 -.-> MF1

    %% Master功能关联
    MF1 --> MF2
    MF2 --> MF3
    MF3 --> MF4

    %% Worker组件关联
    W1 -.-> WF1
    W1 -.-> WF2
    W1 -.-> WF3
    W1 -.-> WF4

    classDef external fill:#e1f5fe
    classDef master fill:#f3e5f5
    classDef worker fill:#e8f5e8
    classDef config fill:#fff3e0
    classDef function fill:#fce4ec

    class A,B,C external
    class M1,M2,M3 master
    class W1,W2,W3,W4 worker
    class CFG,LOG,MON config
    class MF1,MF2,MF3,MF4,WF1,WF2,WF3,WF4 function
```

## 数据流程图

```mermaid
sequenceDiagram
    participant C as 配置系统
    participant M as Master节点
    participant E as etcd
    participant W as Worker节点
    participant T as 目标网站
    participant S as 存储系统

    %% 系统启动阶段
    Note right of C: 系统启动阶段
    C->>M: 加载配置(config.toml)
    M->>E: 注册Master服务
    W->>E: 注册Worker服务
    E->>M: 通知Worker节点上线

    %% 任务调度阶段
    Note right of M: 任务调度阶段
    M->>M: Leader选举
    M->>W: 分配爬虫任务

    %% 爬虫执行阶段
    Note right of W: 爬虫执行阶段
    loop 任务执行循环
        W->>W: 从调度器获取请求
        W->>T: 发送HTTP请求
        T->>W: 返回网页内容
        W->>W: 解析网页内容
        W->>W: 提取新链接
        W->>W: 新链接加入队列
        W->>S: 存储解析结果
    end

    %% 监控与协调阶段
    Note right of M: 监控与协调阶段
    W->>E: 定期心跳
    E->>M: 节点状态变化
    M->>M: 重新分配任务(如果需要)
```

## 组件交互图

```mermaid
graph LR
    subgraph "Master节点"
        M_API[gRPC API]
        M_ELECT[选举模块]
        M_SCHED[调度模块]
        M_MON[监控模块]
    end

    subgraph "Worker节点"
        W_API[gRPC API]
        W_ENGINE[爬虫引擎]
        W_SCHED[本地调度器]
        W_FETCH[采集器]
        W_PARSE[解析器]
        W_STORE[存储器]
    end

    subgraph "外部服务"
        ETCD[etcd]
        DB[(MySQL)]
        PROXY[代理池]
    end

    %% Master内部连接
    M_API --> M_ELECT
    M_API --> M_SCHED
    M_ELECT --> M_MON
    M_SCHED --> M_MON

    %% Worker内部连接
    W_API --> W_ENGINE
    W_ENGINE --> W_SCHED
    W_SCHED --> W_FETCH
    W_FETCH --> W_PARSE
    W_PARSE --> W_STORE
    W_PARSE --> W_SCHED

    %% 跨节点通信
    M_SCHED -.-> W_ENGINE
    M_MON -.-> W_API

    %% 外部依赖
    M_ELECT -.-> ETCD
    M_MON -.-> ETCD
    W_ENGINE -.-> ETCD
    W_FETCH -.-> PROXY
    W_STORE -.-> DB

    classDef master fill:#f3e5f5
    classDef worker fill:#e8f5e8
    classDef external fill:#e1f5fe

    class M_API,M_ELECT,M_SCHED,M_MON master
    class W_API,W_ENGINE,W_SCHED,W_FETCH,W_PARSE,W_STORE worker
    class ETCD,DB,PROXY external
```

## 核心模块层次图

```mermaid
graph TD
    A[分布式爬虫系统]

    A --> B[控制平面 Control Plane]
    A --> C[数据平面 Data Plane]
    A --> D[基础设施 Infrastructure]

    B --> B1[Master服务]
    B1 --> B11[Leader选举]
    B1 --> B12[任务调度]
    B1 --> B13[负载均衡]
    B1 --> B14[节点监控]

    C --> C1[Worker服务]
    C1 --> C11[爬虫引擎]
    C11 --> C111[调度器]
    C11 --> C112[采集器]
    C11 --> C113[解析器]
    C11 --> C114[存储器]

    D --> D1[服务发现 etcd]
    D --> D2[数据存储 MySQL]
    D --> D3[网络代理]
    D --> D4[配置管理]
    D --> D5[日志系统]

    classDef control fill:#f3e5f5
    classDef data fill:#e8f5e8
    classDef infra fill:#e1f5fe

    class B,B1,B11,B12,B13,B14 control
    class C,C1,C11,C111,C112,C113,C114 data
    class D,D1,D2,D3,D4,D5 infra
```

## 关键特性说明

### 🎯 高可用设计
- **多Master选举**：通过etcd实现自动Leader选举
- **故障转移**：Leader失效时自动切换
- **服务发现**：实时监控Worker节点状态

### ⚡ 性能优化
- **并发控制**：可配置的Worker数量
- **限流机制**：令牌桶算法控制请求频率
- **代理轮询**：多代理IP避免被封

### 🔧 扩展性设计
- **插件架构**：支持自定义Fetcher、Parser、Storage
- **动态配置**：热加载配置变更
- **水平扩展**：Worker节点可无限扩展

### 🛡️ 容错机制
- **错误重试**：失败请求自动重试
- **panic恢复**：Worker异常自动恢复
- **状态同步**：通过etcd保证状态一致性

这些架构图从不同角度展示了系统的设计，帮助你更好地理解各个组件之间的关系和协作方式。