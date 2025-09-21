# 代码阅读说明（建议路线）

本文帮助你在 1~2 小时内理清本项目的入口、运行形态与核心模块。

- 入口
  - main: d:/data/repos/repo_reference/crawler/main.go（调用 Cobra 执行子命令）
  - 命令: d:/data/repos/repo_reference/crawler/cmd/cmd.go（注册 master、worker、version 子命令）
- Worker 进程
  - 配置加载/日志初始化/依赖构建：d:/data/repos/repo_reference/crawler/cmd/worker/worker.go
  - Engine 构建：engine.NewEngine + WithSeeds/WithScheduler/WithFetcher/WithStorage
  - 对外服务：gRPC（go-micro）+ HTTP 网关（grpc-gateway）
- Master 进程
  - d:/data/repos/repo_reference/crawler/cmd/master/master.go 启动；master/master.go 负责注册/监听资源（etcd）
- Engine 内核
  - 调度器与工作流：d:/data/repos/repo_reference/crawler/engine/schedule.go
  - 选读：engine/option.go（依赖注入）
- 规则与任务
  - 固定规则：parse/doubanbook、parse/doubangroup
  - JS 规则：parse/doubangroupjs
  - 最小示例：parse/minimal（本次新增）
- 抓取与存储
  - 抓取：collect/collect.go（BrowserFetch 模拟浏览器）
  - 存储：storage/sqlstorage

推荐阅读顺序
1) main.go -> cmd/cmd.go -> cmd/worker/worker.go（理解启动参数和配置）
2) engine/option.go -> engine/schedule.go（理解 Engine 的协作：Seeds -> Scheduler -> CreateWork -> HandleResult）
3) parse/* 了解 RuleTree（Root/Trunk）与 ParseResult 的数据流
4) collect/collect.go（Fetcher）与 storage/sqlstorage（结果入库）

关键跳转点
- Seeds 来源：worker.ParseTaskConfig 读取 [Tasks] 并构建 []*spider.Task，随后 engine.NewEngine.WithSeeds 注入。
- Root 执行：engine.Crawler.handleSeeds 中为每个 Task 复制预置规则（Store.Hash[Name]）并执行 Root 生成初始请求。
- Worker 流程：Schedule.Pull -> Request.Fetch -> Rule.ParseFunc(ctx) -> scheduler.Push(新请求) -> out 通道 -> HandleResult 存储。

调试建议
- 首次运行：在 cmd/worker/worker.go 的 engine.NewEngine 与 s.Run 之间打断点，检查 Seeds 数量、WorkCount、Scheduler 是否工作。
- 观察调度：engine/schedule.go 的 CreateWork/HandleResult；检查 Visited 去重、失败重试逻辑（SetFailure）。
- 网络：collect.BrowserFetch.User-Agent 随机、可选 Proxy；可在 config.toml[fetcher] 调整 timeout/proxy。

运行方式（PowerShell）
- 开发机直跑（需本地 MySQL/etcd）：
  - go build -o crawler .
  - .\crawler worker --id=1 --http=:8080 --grpc=:9090
  - .\crawler master --id=2 --http=:8082 --grpc=:9092
- 一键：docker-compose up -d（包含 mysql 与 etcd）

更多细节见 docs/familiarization_plan.md 与 docs/new_task_minimal.md。