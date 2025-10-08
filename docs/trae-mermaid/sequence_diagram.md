# 分布式爬虫系统 - 序列图

## 1. 系统启动流程

```mermaid
sequenceDiagram
    participant User as 用户
    participant Master as Master节点
    participant Etcd as Etcd集群
    participant Worker as Worker节点
    participant Config as 配置管理

    User->>Master: 启动Master节点
    Master->>Config: 加载配置文件
    Config-->>Master: 返回配置信息
    Master->>Etcd: 注册服务信息
    Etcd-->>Master: 确认注册成功
    Master->>Etcd: 开始领导选举
    Etcd-->>Master: 选举结果
    
    User->>Worker: 启动Worker节点
    Worker->>Config: 加载配置文件
    Config-->>Worker: 返回配置信息
    Worker->>Etcd: 注册到集群
    Etcd-->>Worker: 确认注册成功
    Worker->>Master: 向Master注册
    Master-->>Worker: 确认注册成功
    
    Master->>Master: 初始化任务队列
    Master->>Worker: 发送心跳检测
    Worker-->>Master: 响应心跳
```

## 2. 任务分发与执行流程

```mermaid
sequenceDiagram
    participant User as 用户
    participant Master as Master节点
    participant Worker as Worker节点
    participant Scheduler as 任务调度器
    participant Engine as 爬虫引擎
    participant Parser as 解析器
    participant Storage as 存储系统

    User->>Master: 提交爬虫任务
    Master->>Master: 解析任务配置
    Master->>Scheduler: 创建任务实例
    Scheduler-->>Master: 任务创建成功
    
    Master->>Worker: 分发任务
    Worker-->>Master: 确认接收任务
    
    Worker->>Engine: 启动爬虫引擎
    Engine->>Scheduler: 获取待处理请求
    Scheduler-->>Engine: 返回请求列表
    
    loop 处理请求队列
        Engine->>Engine: 发送HTTP请求
        Engine->>Parser: 解析响应内容
        Parser-->>Engine: 返回解析结果
        Engine->>Scheduler: 添加新发现的链接
        Engine->>Storage: 保存提取的数据
        Storage-->>Engine: 确认保存成功
    end
    
    Engine->>Worker: 报告任务进度
    Worker->>Master: 上报执行状态
    Master-->>User: 返回任务状态
```

## 3. 错误处理与重试流程

```mermaid
sequenceDiagram
    participant Engine as 爬虫引擎
    participant HTTPClient as HTTP客户端
    participant ProxyPool as 代理池
    participant RateLimiter as 限流器
    participant Scheduler as 任务调度器
    participant Logger as 日志系统

    Engine->>RateLimiter: 检查请求频率
    RateLimiter-->>Engine: 允许请求
    
    Engine->>ProxyPool: 获取代理
    ProxyPool-->>Engine: 返回代理地址
    
    Engine->>HTTPClient: 发送HTTP请求
    HTTPClient-->>Engine: 请求失败(超时/错误)
    
    Engine->>Logger: 记录错误日志
    Engine->>Engine: 检查重试次数
    
    alt 重试次数未超限
        Engine->>ProxyPool: 释放失败代理
        Engine->>ProxyPool: 获取新代理
        ProxyPool-->>Engine: 返回新代理
        Engine->>HTTPClient: 重新发送请求
        HTTPClient-->>Engine: 请求成功
    else 重试次数超限
        Engine->>Scheduler: 标记请求失败
        Engine->>Logger: 记录最终失败日志
    end
```

## 4. 数据存储流程

```mermaid
sequenceDiagram
    participant Parser as 解析器
    participant Engine as 爬虫引擎
    participant StorageInterface as 存储接口
    participant MySQL as MySQL存储
    participant FileStorage as 文件存储
    participant Logger as 日志系统

    Parser->>Engine: 返回解析的数据项
    Engine->>StorageInterface: 调用存储接口
    
    alt 结构化数据
        StorageInterface->>MySQL: 保存到数据库
        MySQL-->>StorageInterface: 确认保存成功
    else 非结构化数据
        StorageInterface->>FileStorage: 保存到文件
        FileStorage-->>StorageInterface: 确认保存成功
    end
    
    StorageInterface-->>Engine: 返回存储结果
    
    alt 存储成功
        Engine->>Logger: 记录成功日志
    else 存储失败
        Engine->>Logger: 记录错误日志
        Engine->>Engine: 加入重试队列
    end
```

## 5. 集群监控与健康检查

```mermaid
sequenceDiagram
    participant Master as Master节点
    participant Worker as Worker节点
    participant Etcd as Etcd集群
    participant Prometheus as Prometheus
    participant Grafana as Grafana
    participant AlertManager as 告警管理

    loop 定期健康检查
        Master->>Worker: 发送心跳请求
        Worker-->>Master: 响应心跳
        
        alt Worker响应正常
            Master->>Etcd: 更新节点状态
        else Worker无响应
            Master->>Etcd: 标记节点异常
            Master->>Master: 重新分配任务
        end
    end
    
    loop 指标收集
        Worker->>Prometheus: 上报性能指标
        Master->>Prometheus: 上报集群状态
        Prometheus->>Grafana: 提供监控数据
    end
    
    alt 发现异常
        Prometheus->>AlertManager: 触发告警
        AlertManager->>AlertManager: 发送告警通知
    end
```

## 6. 系统关闭流程

```mermaid
sequenceDiagram
    participant User as 用户
    participant Master as Master节点
    participant Worker as Worker节点
    participant Engine as 爬虫引擎
    participant Storage as 存储系统
    participant Etcd as Etcd集群

    User->>Master: 发送关闭信号
    Master->>Worker: 通知停止接收新任务
    Worker-->>Master: 确认停止接收
    
    Worker->>Engine: 等待当前任务完成
    Engine->>Storage: 保存未完成的状态
    Storage-->>Engine: 确认保存完成
    Engine-->>Worker: 任务清理完成
    
    Worker->>Master: 报告关闭准备就绪
    Master->>Etcd: 注销Worker节点
    Etcd-->>Master: 确认注销成功
    
    Master->>Etcd: 注销Master节点
    Etcd-->>Master: 确认注销成功
    
    Master-->>User: 系统关闭完成
```

## 序列图说明

### 系统启动流程
展示了Master和Worker节点的启动过程，包括配置加载、服务注册、领导选举等关键步骤。

### 任务分发与执行流程
描述了从用户提交任务到任务执行完成的完整流程，包括任务解析、分发、执行和结果存储。

### 错误处理与重试流程
展示了系统如何处理网络错误、代理失效等异常情况，以及重试机制的工作原理。

### 数据存储流程
说明了解析后的数据如何根据类型选择不同的存储方式，以及存储失败时的处理机制。

### 集群监控与健康检查
描述了集群的健康监控机制，包括心跳检测、指标收集和异常告警。

### 系统关闭流程
展示了系统优雅关闭的过程，确保正在执行的任务能够正常完成并保存状态。