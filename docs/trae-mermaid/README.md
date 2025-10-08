# 分布式爬虫系统 - Mermaid架构图

本目录包含使用Mermaid格式绘制的分布式爬虫系统架构图，提供了系统的多个视角和层面的可视化展示。

## 📁 文件列表

| 文件名 | 描述 | 图表类型 |
|--------|------|----------|
| `class_diagram.md` | 类图 - 展示系统核心类和接口关系 | Class Diagram |
| `component_diagram.md` | 组件图 - 展示高层架构和组件交互 | Component Diagram |
| `sequence_diagram.md` | 序列图 - 展示系统工作流程和时序 | Sequence Diagram |
| `deployment_diagram.md` | 部署图 - 展示部署架构和基础设施 | Deployment Diagram |

## 🎯 图表概览

### 1. 类图 (Class Diagram)
**文件**: `class_diagram.md`

展示分布式爬虫系统的核心类、接口和它们之间的关系，包括：
- **主控模块**: Master、WorkerNode
- **爬虫引擎**: Crawler、Schedule、CrawlerStore
- **爬虫任务**: Task、Property、Request、RuleTree、Rule
- **网络模块**: Fetcher、HTTPClient
- **存储模块**: Storage、MySQLStorage、FileStorage
- **代理模块**: ProxyProvider、ProxyFunc、ProxyPool
- **限流模块**: RateLimiter、TokenBucket

### 2. 组件图 (Component Diagram)
**文件**: `component_diagram.md`

展示系统的高层架构和组件间的交互关系，包括：
- 外部环境（互联网、目标网站）
- 负载均衡层
- 主控集群和工作节点集群
- 爬虫引擎和解析器插件
- 存储层、网络组件
- 配置管理、日志系统
- 服务协调（Etcd）和监控系统

### 3. 序列图 (Sequence Diagram)
**文件**: `sequence_diagram.md`

展示系统的工作流程和交互时序，包含6个主要流程：
- **系统启动流程**: Master和Worker节点启动过程
- **任务分发与执行流程**: 从任务提交到执行完成
- **错误处理与重试流程**: 异常处理和重试机制
- **数据存储流程**: 数据存储的选择和处理
- **集群监控与健康检查**: 监控和健康检查机制
- **系统关闭流程**: 优雅关闭过程

### 4. 部署图 (Deployment Diagram)
**文件**: `deployment_diagram.md`

展示系统的部署架构和基础设施，包括：
- 网络层级（互联网、负载均衡、安全层）
- Kubernetes集群（Master节点、Worker节点、应用Pod）
- 服务协调（Etcd集群）
- 数据存储层（MySQL集群、缓存层、文件存储）
- 监控系统（Prometheus、Grafana、AlertManager）
- 日志系统（ELK Stack）
- 网络安全组件

## 🔧 如何查看图表

### 方法1: GitHub/GitLab在线查看
大多数Git托管平台都支持Mermaid图表的在线渲染，直接在浏览器中打开Markdown文件即可查看。

### 方法2: VS Code插件
安装以下VS Code插件之一：
- **Mermaid Preview** (推荐)
- **Markdown Preview Mermaid Support**
- **Mermaid Markdown Syntax Highlighting**

安装后在VS Code中打开Markdown文件，使用预览功能即可查看图表。

### 方法3: 在线编辑器
使用Mermaid官方在线编辑器：
- 访问 [Mermaid Live Editor](https://mermaid.live/)
- 复制图表代码到编辑器中
- 实时预览和编辑图表

### 方法4: 本地工具
安装Mermaid CLI工具：
```bash
npm install -g @mermaid-js/mermaid-cli
mmdc -i class_diagram.md -o class_diagram.png
```

## 🏗️ 系统架构特点

### 分布式架构
- **主从架构**: Master节点负责任务分发，Worker节点执行任务
- **水平扩展**: Worker节点可动态增减
- **负载均衡**: 任务在多个Worker节点间均匀分配

### 高可用性
- **多Master节点**: 支持Master节点的高可用部署
- **故障转移**: 节点故障时自动重新分配任务
- **健康检查**: 定期检查节点状态

### 可扩展性
- **插件化解析器**: 支持多种网站的解析规则
- **多种存储方式**: 支持MySQL、文件等多种存储
- **代理池**: 支持多种代理策略

### 可观测性
- **全面监控**: Prometheus + Grafana监控体系
- **日志收集**: ELK Stack日志分析
- **告警机制**: 异常情况及时通知

## 🛠️ 技术栈

### 核心技术
- **编程语言**: Go
- **容器化**: Docker + Kubernetes
- **服务发现**: Etcd
- **消息队列**: 内置任务队列

### 存储技术
- **关系数据库**: MySQL
- **缓存**: Redis
- **文件存储**: NFS、S3

### 监控日志
- **监控**: Prometheus + Grafana
- **日志**: Elasticsearch + Logstash + Kibana
- **告警**: AlertManager

### 网络安全
- **负载均衡**: ALB + NLB
- **安全防护**: WAF + 防火墙
- **访问控制**: VPN网关

## 📚 相关文档

- [项目总体架构](../trae/00_top_down_overview.md)
- [分布式框架详解](../trae/01_distributed_framework.md)
- [模块映射关系](../trae/02_module_map.md)
- [学习路径指南](../trae/03_learning_path_for_go_server_dev.md)
- [Windows快速开始](../trae/04_windows_quickstart.md)

## 🤝 贡献指南

如需更新或完善架构图：

1. 修改对应的Markdown文件中的Mermaid代码
2. 确保图表语法正确
3. 更新相关说明文档
4. 提交Pull Request

## 📄 许可证

本文档遵循项目的开源许可证。