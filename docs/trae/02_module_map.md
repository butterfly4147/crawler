# 模块地图（文件夹 → 职责 → 关键类型）

帮助你在代码树中快速定位实现与边界，结合你的游戏服务器经验，从调度、并发、RPC、存储逐层掌握。

## 入口与命令
- `main.go`：进程入口，调用 `cmd.Execute()`。
- `cmd/cmd.go`：Cobra 命令注册（`master/worker/version`）。
- `cmd/master/master.go`：Master 启动与服务注册（HTTP/GRPC），读取配置、依赖构建、选举与 RPC 网关。
- `cmd/worker/worker.go`：Worker 引擎构建与服务注册（HTTP/GRPC），Seeds 注入与运行。

## 分布式与协调
- `master/master.go`：`Master` 结构维护 `workNodes/resources`，实现选举、分配、RPC 处理。
- `engine/schedule.go`：引擎内核（请求队列、并发工作、去重/重试、资源监听），包含 `Crawler/Store/Schedule`。
- `engine/option.go`：依赖注入（Option 模式）。

## 爬虫任务与解析
- `spider/task.go`：任务属性、配置结构体。
- `spider/parse.go`：规则树 `RuleTree` 与解析节点 `Rule`（`Root/Trunk/ItemFields`）。
- `spider/request.go`：请求结构、`Fetch/Check/Unique`、JS 解析辅助。
- `parse/*`：具体站点示例（`doubanbook/doubangroup/doubangroupjs/minimal`）。

## 数据采集与限流
- `collect/collect.go`：抓取器接口（如 `BrowserFetch`）。
- `limiter/limiter.go`：多限流器组合（令牌桶），`spider.Request.Fetch` 中调用 `Task.Limit.Wait`。
- `proxy/proxy.go`：代理轮询（Round-Robin），支持多代理 URL 切换。

## 存储层
- `storage/sqlstorage/*.go`：批量缓冲与表结构推断；调用 `sqldb` 执行建表/插入。
- `sqldb/*.go`：MySQL 访问封装（`CreateTable/Insert/DropTable`）。

## RPC 与 API
- `proto/crawler/crawler.proto`：定义 `CrawlerMaster` 服务（`AddResource/DeleteResource`），并通过 `grpc-gateway` 暴露 `/crawler/resource`。
- `proto/crawler/*`：gRPC、go-micro、gateway 生成代码（客户端/服务端/HTTP 映射）。

## 其他
- `auth/auth.go`：微服务处理器包装，提取与校验上下文元数据（可用于鉴权）。
- `extensions/randomua.go`：随机 UA。
- `log/*`：日志封装。
- `kubernetes/*`：K8s 部署示例。