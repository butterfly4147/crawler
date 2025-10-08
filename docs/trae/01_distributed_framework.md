# 分布式框架拆解（自顶向下）

本章聚焦 Master-Worker 架构、选举与任务分配、Worker 引擎执行与结果存储，帮助你快速建立心智模型。

## Master 职责
- 选举与高可用：基于 `etcd/clientv3/concurrency.Election` 选出 Leader（见 `master/master.go`）。Leader 负责资源管理与任务分配。
- 资源模型：`ResourceSpec` 表示任务资源，`NodeSpec` 表示 Worker 节点信息与负载。
- 任务分配：`Master.Assign` 按 `Payload`（负载）选择最空闲节点；通过 `CrawlerMaster` RPC 将资源下发，Worker 侧监听 `etcd` 的 `RESOURCEPATH`。
- 对外服务：
  - gRPC：`proto/crawler/crawler_grpc.pb.go` 注册 `CrawlerMaster`。
  - HTTP 网关：`proto/crawler/crawler.pb.gw.go` 将 RPC 映射到 REST（`POST/DELETE /crawler/resource`）。
  - 命令行：`cmd/master/master.go` 提供 `--http --grpc --id --config --podip` 等参数。

## Worker 职责
- 引擎初始化：`cmd/worker/worker.go` 构建 `Fetcher/Storage/Scheduler/Seeds/Logger` 后创建引擎（`engine.NewEngine`）。
- 运行流程：`engine.Crawler.Run` 启动 Seeds 处理、资源加载与监听（集群）、调度循环与工作协程、结果处理。
- 调度与并发：`Schedule` 持有 `requestCh/workerCh` 与内存队列；`CreateWork` 从 `Pull` 获取请求，执行 `Check/Fetch/Parse`，递归 `Push` 新请求。
- 去重与重试：`Visited` 哈希做去重；`SetFailure` 首次失败重入队，避免误判导致任务中断。
- 存储：默认 `SQLStorage`，批量缓冲 `dataDocker`，按 `Rule.ItemFields` 构造表结构与插入数据（见 `storage/sqlstorage/*.go`）。

## 配置与扩展点
- `config.toml`：定义 Tasks（Name/WaitTime/Reload/MaxDepth/Fetcher/Limits）、Fetcher（timeout/proxy）、Storage（sqlURL）。
- 插件化接口：
  - Fetcher：`collect/collect.go` 定义抓取器接口，可替换为浏览器/HTTP 客户端等实现。
  - Storage：`spider.Storage` 接口可扩展到 ES/ClickHouse/文件等。
  - Parser/Rule：`spider.RuleTree` 支持自定义解析与 JS 规则（`otto` 引擎，`parse/doubangroupjs`）。
- 运行参数：
  - Worker：`--id --http --grpc`（默认 `8080/9090`）。
  - Master：`--id --http --grpc --podip`（默认 `8082/9092`）。

## 服务发现与选举
- Worker 注册与发现：使用 go-micro etcd 注册；Master 维护 `workNodes`。
- Master 选举：`concurrency.NewElection` 参与选举，成为 Leader 后加载资源并开始分配。
- 资源下发：通过 RPC 写入 `etcd` 的 `RESOURCEPATH`；Worker `watchResource/loadResource` 处理新增或删除任务。

## 典型调用链
1. Master 启动，参与选举，成为 Leader。
2. Master 暴露 `CrawlerMaster.AddResource/DeleteResource`（HTTP 对应 `/crawler/resource`）。
3. Worker 启动，引擎 `Run` 监听 `RESOURCEPATH`，根据分配的资源 `runTasks/deleteTasks`。
4. Worker 执行任务：获取请求 -> 抓取 -> 解析 -> 产出数据 -> 存储。