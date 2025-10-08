# 分布式爬虫系统 - 组件图

```mermaid
graph TB
    %% 外部服务
    subgraph "外部环境"
        Internet[互联网]
        TargetSites[目标网站]
    end

    %% 负载均衡层
    subgraph "负载均衡层"
        LB[负载均衡器]
    end

    %% 主控集群
    subgraph "主控集群"
        Master1[Master节点1]
        Master2[Master节点2]
        Master3[Master节点3]
    end

    %% 工作节点集群
    subgraph "工作节点集群"
        Worker1[Worker节点1]
        Worker2[Worker节点2]
        Worker3[Worker节点3]
        WorkerN[Worker节点N]
    end

    %% 爬虫引擎组件
    subgraph "爬虫引擎"
        direction TB
        CrawlerEngine[爬虫引擎]
        TaskScheduler[任务调度器]
        RequestQueue[请求队列]
        ResponseProcessor[响应处理器]
    end

    %% 解析器插件
    subgraph "解析器插件"
        DoubanBookParser[豆瓣图书解析器]
        DoubanGroupParser[豆瓣小组解析器]
        JSParser[JavaScript解析器]
        MinimalParser[最小解析器]
    end

    %% 存储层
    subgraph "存储层"
        MySQL[(MySQL数据库)]
        FileStorage[文件存储]
        StorageInterface[存储接口]
    end

    %% 网络组件
    subgraph "网络组件"
        HTTPClient[HTTP客户端]
        ProxyPool[代理池]
        RateLimiter[限流器]
    end

    %% 配置管理
    subgraph "配置管理"
        ConfigManager[配置管理器]
        ConfigFile[配置文件]
    end

    %% 日志系统
    subgraph "日志系统"
        Logger[日志记录器]
        LogFile[日志文件]
    end

    %% 服务发现与协调
    subgraph "服务协调"
        Etcd[(Etcd集群)]
        ServiceDiscovery[服务发现]
        LeaderElection[领导选举]
    end

    %% 监控系统
    subgraph "监控系统"
        Prometheus[Prometheus]
        Grafana[Grafana]
        AlertManager[告警管理]
    end

    %% 连接关系
    Internet --> LB
    LB --> Master1
    LB --> Master2
    LB --> Master3

    Master1 --> Worker1
    Master1 --> Worker2
    Master2 --> Worker3
    Master3 --> WorkerN

    Worker1 --> CrawlerEngine
    Worker2 --> CrawlerEngine
    Worker3 --> CrawlerEngine
    WorkerN --> CrawlerEngine

    CrawlerEngine --> TaskScheduler
    TaskScheduler --> RequestQueue
    RequestQueue --> ResponseProcessor

    CrawlerEngine --> DoubanBookParser
    CrawlerEngine --> DoubanGroupParser
    CrawlerEngine --> JSParser
    CrawlerEngine --> MinimalParser

    ResponseProcessor --> StorageInterface
    StorageInterface --> MySQL
    StorageInterface --> FileStorage

    CrawlerEngine --> HTTPClient
    HTTPClient --> ProxyPool
    HTTPClient --> RateLimiter
    HTTPClient --> TargetSites

    Master1 --> Etcd
    Master2 --> Etcd
    Master3 --> Etcd
    Worker1 --> Etcd
    Worker2 --> Etcd
    Worker3 --> Etcd
    WorkerN --> Etcd

    Etcd --> ServiceDiscovery
    Etcd --> LeaderElection

    ConfigManager --> ConfigFile
    CrawlerEngine --> ConfigManager

    CrawlerEngine --> Logger
    Logger --> LogFile

    Master1 --> Prometheus
    Worker1 --> Prometheus
    Prometheus --> Grafana
    Prometheus --> AlertManager

    %% 样式定义
    classDef masterNode fill:#e1f5fe,stroke:#01579b,stroke-width:2px
    classDef workerNode fill:#f3e5f5,stroke:#4a148c,stroke-width:2px
    classDef storage fill:#e8f5e8,stroke:#1b5e20,stroke-width:2px
    classDef network fill:#fff3e0,stroke:#e65100,stroke-width:2px
    classDef parser fill:#fce4ec,stroke:#880e4f,stroke-width:2px
    classDef external fill:#f5f5f5,stroke:#424242,stroke-width:2px

    class Master1,Master2,Master3 masterNode
    class Worker1,Worker2,Worker3,WorkerN workerNode
    class MySQL,FileStorage,Etcd storage
    class HTTPClient,ProxyPool,RateLimiter network
    class DoubanBookParser,DoubanGroupParser,JSParser,MinimalParser parser
    class Internet,TargetSites external
```

## 组件图说明

### 外部环境
- **互联网**: 系统运行的网络环境
- **目标网站**: 需要爬取的目标站点

### 负载均衡层
- **负载均衡器**: 分发请求到不同的Master节点

### 主控集群
- **Master节点**: 主控节点集群，负责任务分发和集群管理
- 支持多节点部署，通过Etcd进行领导选举

### 工作节点集群
- **Worker节点**: 工作节点集群，执行具体的爬虫任务
- 可水平扩展，动态加入和退出集群

### 爬虫引擎
- **爬虫引擎**: 核心爬虫逻辑
- **任务调度器**: 管理任务的调度和分发
- **请求队列**: 缓存待处理的HTTP请求
- **响应处理器**: 处理HTTP响应和数据提取

### 解析器插件
- **豆瓣图书解析器**: 专门解析豆瓣图书页面
- **豆瓣小组解析器**: 专门解析豆瓣小组页面
- **JavaScript解析器**: 处理需要JS渲染的页面
- **最小解析器**: 基础解析器实现

### 存储层
- **MySQL数据库**: 结构化数据存储
- **文件存储**: 非结构化数据存储
- **存储接口**: 统一的存储抽象层

### 网络组件
- **HTTP客户端**: 发送HTTP请求
- **代理池**: 管理代理服务器
- **限流器**: 控制请求频率

### 配置管理
- **配置管理器**: 管理系统配置
- **配置文件**: 存储配置信息

### 日志系统
- **日志记录器**: 记录系统日志
- **日志文件**: 存储日志信息

### 服务协调
- **Etcd集群**: 分布式键值存储，用于服务发现和配置管理
- **服务发现**: 自动发现集群中的服务
- **领导选举**: Master节点的领导选举机制

### 监控系统
- **Prometheus**: 指标收集和存储
- **Grafana**: 监控数据可视化
- **告警管理**: 系统异常告警