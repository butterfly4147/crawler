# 学习路径（面向有 2 年 Go 服务器经验）

目标：在 3~7 天内从“能跑通”到“能改造”，最终“能扩展与部署”。

## 阶段一：跑通与观察（0.5~1 天）
- 准备环境：Docker Desktop、MySQL、etcd。
- 直接 `docker-compose up -d`，观察端口：Worker `8080/9090`，Master `8082/9092`，etcd `2379/2380`。
- 打断点：`engine/schedule.go` 的 `CreateWork/HandleResult`、`SetFailure`、`handleSeeds`，确认 Seeds 注入与调度。
- 用浏览器或 `Invoke-WebRequest` 访问 HTTP 网关接口，触发 `CrawlerMaster.AddResource` 流程。

## 阶段二：理解架构与数据流（1 天）
- 通读：`cmd/master/master.go`、`cmd/worker/worker.go`、`engine/schedule.go`、`spider/*`、`storage/sqlstorage/*`。
- 画图：Master 选举与资源分配、Worker 监听与任务运行、引擎工作流（Pull/Fetch/Parse/Push/Out）。
- 自查：去重逻辑（`Visited`）、失败重试（`failures`）、批量存储（`dataDocker`）。

## 阶段三：动手扩展（1~2 天）
- 新增解析器：在 `parse/minimal` 基础上添加一个你熟悉的网站解析规则（`RuleTree`），并在 `engine.Store.Add` 注册。
- 自定义 Fetcher：实现一个轻量 HTTP Fetcher（或接入浏览器驱动），替换 `collect.BrowserFetch`。
- 扩展存储：实现 `spider.Storage` 的新后端（如 ES/ClickHouse），在 `cmd/worker` 中通过 `WithStorage` 注入。

## 阶段四：性能与稳定性（1~2 天）
- 限流与熔断：调整 `config.toml [Tasks].Limits`、go-micro Hystrix 包装；验证 QPS 与故障处理。
- 观测性：接入 pprof（已集成），评估 goroutine 与阻塞点；必要时添加 Prometheus 指标。
- 容错演练：模拟 Worker 故障，确认 Master 的任务重分配与恢复（参考 `CHANGELOG.md` v0.3.8）。

## 阶段五：部署与治理（1 天）
- 容器化：理解 `Dockerfile` 双阶段构建；按需裁剪镜像。
- K8s：阅读 `kubernetes/*`，在本地或云端部署，使用 `--podip` 生成分布式 ID；通过 Service/Ingress 暴露网关。
- 配置治理：把 `config.toml` 外置成 ConfigMap（K8s），或通过环境变量覆盖。

## 练手题（建议）
- 增加一个“新闻站点列表页 → 详情页”的双层解析任务，包含去重、限流与代理轮换。
- 为 SQLStorage 增加“表存在即跳过建表”的缓存与 TTL 过期策略（现有有 `Table` 缓存，可完善策略）。
- 在 `auth` 中实现简单的签名校验与白名单路径。