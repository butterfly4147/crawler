# 分布式爬虫系统 UML 架构图

本文件夹包含了分布式爬虫系统的完整UML架构图，使用PlantUML格式编写。

## 文件说明

### 1. 类图 (class_diagram.puml)
展示系统的核心类和它们之间的关系，包括：
- **主控模块**: Master、NodeSpec、ResourceSpec
- **爬虫引擎**: Crawler、Schedule、CrawlerStore
- **爬虫任务**: Task、TaskConfig、Request、ParseResult
- **存储模块**: Storage接口和SQLStorage实现
- **代理模块**: ProxyFunc代理函数
- **限流模块**: RateLimiter和MultiLimiter

### 2. 组件图 (component_diagram.puml)
展示系统的高层架构和组件间的交互，包括：
- **Master节点**: 资源管理、任务调度、节点监控、领导选举
- **Worker节点**: 爬虫引擎、任务执行、数据解析、请求调度
- **解析器插件**: 豆瓣图书、豆瓣小组、JavaScript、通用解析器
- **存储组件**: SQL存储、文件存储、缓存存储
- **网络组件**: HTTP客户端、代理管理、限流器、用户代理管理
- **基础设施**: etcd集群、MySQL数据库

### 3. 序列图 (sequence_diagram.puml)
展示爬虫系统的完整工作流程，包括：
- **系统启动阶段**: 服务注册、领导选举
- **任务分发阶段**: 任务创建、资源分配、任务分发
- **爬取执行阶段**: 请求处理、数据解析、结果存储
- **异常处理**: 错误重试、故障转移
- **任务完成阶段**: 状态更新、资源清理
- **系统监控**: 健康检查、负载均衡
- **系统关闭**: 优雅停机

### 4. 部署图 (deployment_diagram.puml)
展示系统的部署架构和基础设施，包括：
- **Kubernetes集群**: Master Pod、Worker Pod集群
- **etcd集群**: 三节点高可用配置
- **数据存储层**: MySQL主从架构
- **监控系统**: Prometheus + Grafana
- **日志系统**: ELK Stack
- **配置管理**: ConfigMap + Secret
- **负载均衡**: Nginx/HAProxy

## 如何查看UML图

### 方法1: 使用PlantUML在线编辑器
1. 访问 [PlantUML在线编辑器](http://www.plantuml.com/plantuml/uml/)
2. 复制对应的.puml文件内容
3. 粘贴到编辑器中查看图形

### 方法2: 使用VS Code插件
1. 安装PlantUML插件
2. 打开.puml文件
3. 使用快捷键 `Alt+D` 预览图形

### 方法3: 使用命令行工具
```bash
# 安装PlantUML
npm install -g node-plantuml

# 生成PNG图片
puml generate class_diagram.puml
puml generate component_diagram.puml
puml generate sequence_diagram.puml
puml generate deployment_diagram.puml
```

### 方法4: 使用Docker
```bash
# 生成所有图片
docker run --rm -v $(pwd):/work -w /work plantuml/plantuml:latest *.puml
```

## 系统架构特点

### 分布式架构
- Master-Worker模式，支持水平扩展
- 基于etcd的服务发现和配置管理
- 自动故障转移和负载均衡

### 高可用性
- etcd集群保证配置和状态的高可用
- MySQL主从架构保证数据的高可用
- 多Worker节点保证处理能力的高可用

### 可扩展性
- 插件化的解析器架构
- 可配置的存储后端
- 灵活的代理和限流策略

### 可观测性
- 结构化日志记录
- Prometheus指标监控
- 分布式链路追踪支持

## 技术栈

- **编程语言**: Go 1.18+
- **微服务框架**: go-micro v4
- **服务发现**: etcd v3
- **数据库**: MySQL 8.0
- **容器化**: Docker + Kubernetes
- **监控**: Prometheus + Grafana
- **日志**: ELK Stack
- **负载均衡**: Nginx/HAProxy

## 更新说明

本UML图基于项目当前代码结构生成，如果代码发生重大变更，请及时更新对应的UML图文件。

更新步骤：
1. 分析代码变更
2. 修改对应的.puml文件
3. 重新生成图片
4. 更新本README文档