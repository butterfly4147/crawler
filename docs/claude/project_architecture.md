# 分布式爬虫项目架构分析

## 项目概述

这是一个基于Go语言的分布式爬虫系统，采用微服务架构设计，支持任务调度、资源分配、负载均衡等功能。项目来源于极客时间的分布式爬虫课程。

## 核心架构

### 1. 分布式架构模式

项目采用**Master-Worker**分布式架构：

- **Master节点**：负责任务调度、资源分配、负载均衡和集群管理
- **Worker节点**：负责实际的网页抓取、解析和数据处理
- **服务发现**：使用etcd进行服务注册与发现
- **选举机制**：Master节点通过etcd选举实现高可用

### 2. 技术栈

- **微服务框架**: go-micro v4
- **服务注册发现**: etcd
- **RPC通信**: gRPC + HTTP网关
- **配置管理**: TOML格式配置
- **日志系统**: zap
- **限流器**: golang.org/x/time/rate
- **数据库**: MySQL
- **代理**: 支持轮询代理切换

## 核心模块分析

### 2.1 Master模块 (`master/`)

**主要职责：**
- 集群管理：监控Worker节点状态
- 任务调度：分配爬虫任务到Worker节点
- 负载均衡：基于负载情况分配任务
- 选举机制：通过etcd实现Leader选举

**核心组件：**
- `Master`结构体：管理集群状态和资源分配
- `ResourceSpec`：资源规格定义
- `NodeSpec`：节点规格定义
- 选举机制：使用etcd的concurrency包实现

### 2.2 Engine模块 (`engine/`)

**主要职责：**
- 爬虫引擎核心逻辑
- 任务调度和执行
- 请求队列管理
- 结果处理

**核心组件：**
- `Crawler`：爬虫引擎主体
- `Schedule`：调度器接口和实现
- `CrawlerStore`：任务存储管理

### 2.3 Spider模块 (`spider/`)

**主要职责：**
- 定义爬虫任务模型
- 请求和响应处理
- 数据解析规则

**核心组件：**
- `Task`：爬虫任务定义
- `Request`：HTTP请求封装
- `Context`：解析上下文
- `ParseResult`：解析结果

### 2.4 数据采集模块 (`collect/`)

**主要职责：**
- 网页内容抓取
- 代理支持
- 浏览器模拟

### 2.5 存储模块 (`storage/`)

**主要职责：**
- 数据持久化
- 支持SQL存储

### 2.6 解析模块 (`parse/`)

**主要职责：**
- 网页内容解析
- 支持JavaScript动态解析
- 豆瓣小组、图书等特定网站解析

## 分布式框架设计

### 3.1 服务发现与注册

```go
// 使用etcd作为注册中心
reg := etcd.NewRegistry(registry.Addrs(sconfig.RegistryAddress))
service := micro.NewService(
    micro.Registry(reg),
    micro.RegisterTTL(time.Duration(cfg.RegisterTTL)*time.Second),
    micro.RegisterInterval(time.Duration(cfg.RegisterInterval)*time.Second),
)
```

### 3.2 负载均衡策略

Master节点使用**最小负载优先**策略分配任务：

```go
func (m *Master) Assign(r *ResourceSpec) (*NodeSpec, error) {
    candidates := make([]*NodeSpec, 0, len(m.workNodes))
    for _, node := range m.workNodes {
        candidates = append(candidates, node)
    }
    // 按负载排序，选择负载最低的节点
    sort.Slice(candidates, func(i, j int) bool {
        return candidates[i].Payload < candidates[j].Payload
    })
    return candidates[0], nil
}
```

### 3.3 高可用设计

- **Leader选举**：通过etcd实现Master节点选举
- **故障转移**：当Leader失效时自动选举新Leader
- **任务重分配**：Worker节点失效时重新分配任务

### 3.4 限流与熔断

- **限流器**：使用令牌桶算法控制请求频率
- **熔断器**：集成Hystrix实现服务熔断
- **代理轮询**：支持多个代理IP轮换使用

## 数据流分析

### 4.1 任务执行流程

1. **任务定义**：通过配置文件定义爬虫任务
2. **任务分配**：Master节点将任务分配给Worker
3. **网页抓取**：Worker执行HTTP请求获取网页内容
4. **内容解析**：根据规则解析网页内容
5. **数据存储**：将解析结果保存到数据库
6. **链接发现**：从当前页面发现新的链接继续爬取

### 4.2 分布式协作流程

1. **服务注册**：Worker启动时向etcd注册服务
2. **资源监听**：Master监听etcd中的资源变化
3. **负载均衡**：Master根据节点负载分配任务
4. **结果收集**：Worker将爬取结果发送到输出通道

## 配置系统

项目使用TOML格式配置文件，主要配置项包括：

- 日志级别
- 爬虫任务定义
- 代理设置
- 数据库连接
- gRPC服务配置
- 注册中心地址

## 扩展性设计

### 6.1 插件化架构

- **Fetcher接口**：支持不同的网页抓取方式
- **Storage接口**：支持多种数据存储后端
- **Parser接口**：支持自定义解析规则

### 6.2 动态规则

支持JavaScript动态规则，通过otto引擎执行：

```go
func (c *CrawlerStore) AddJSTask(m *spider.TaskModle) {
    // 使用JavaScript定义解析规则
    vm := otto.New()
    vm.Set("AddJsReq", AddJsReqs)
    // ...
}
```

## 部署架构

### 7.1 组件部署

- **etcd集群**：服务注册发现
- **Master集群**：多个Master节点通过选举产生Leader
- **Worker集群**：多个Worker节点执行爬虫任务
- **MySQL数据库**：数据存储

### 7.2 容器化支持

项目提供Dockerfile和docker-compose.yml，支持容器化部署。

## 总结

这个分布式爬虫项目展现了完整的微服务架构设计，具有以下特点：

1. **高可用**：通过etcd选举实现Master节点高可用
2. **可扩展**：Worker节点可水平扩展
3. **负载均衡**：智能的任务分配策略
4. **容错性**：支持任务重试和故障转移
5. **灵活性**：插件化设计支持多种抓取和存储方式

对于有Go语言游戏服务器开发经验的开发者来说，这个项目的并发模型、网络通信和分布式架构设计都很有参考价值。