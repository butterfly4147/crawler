# 熟悉项目的自顶向下路线

面向：具有 2 年 Go 游戏后端经验、首次接手该分布式爬虫项目的开发者。目标是在逐步理解宏观架构的同时，能够在 1~2 周内独立扩展新站点抓取任务并上线。

## 阶段 0：准备环境
- 安装 **Go 1.19+**、**MySQL 5.7**、**etcd 3.5**。
- 执行 `docker-compose up` 以快速启动 MySQL 与 etcd（推荐开发机使用）。
- 阅读 `docs/reading_guide.md`，了解仓库命名、代码规范、调试建议。

## 阶段 1：理解入口与命令
- 路径查阅：`main.go` → `cmd/cmd.go` → `cmd/master` / `cmd/worker`。
- 关注命令行参数：`--config`、`--http`、`--grpc`、`--pprof`、`--cluster` 等，理解如何在容器 / Kubernetes 中注入。
- 运行 `go run . master --config config.toml` 与 `go run . worker --config config.toml`，观察日志输出与 etcd、MySQL 的变化。

## 阶段 2：控制平面（Master）
- 阅读顺序：`master/master.go` → `master/option.go` → `cmd/master/master.go`。
- 核心问题：
  - 如何创建/删除资源（`AddResource`、`DeleteResource`）。
  - Leader 选举如何发生？（etcd session + `/lock` 前缀）。
  - 任务如何绑定到 Worker？（`ResourceSpec.AssignedNode`、`reAssign()` 逻辑）。
- 实践：通过 `curl -X POST localhost:8081/crawler/resource` 新增任务，确认 etcd 键值变化。

## 阶段 3：数据平面（Worker & Engine）
- 阅读顺序：`cmd/worker/worker.go` → `engine/option.go` → `engine/schedule.go`（重点聚焦 `Run`、`Schedule`、`CreateWork`、`HandleResult`）。
- 对照 `config.toml`，理解 Fetcher/Storage/RateLimiter 的注入方式与任务级配置。
- 实践：在非集群模式下运行 Worker，验证 `handleSeeds` 直接读取配置任务；切换 `--cluster=false/true` 对比差异。

## 阶段 4：任务描述与解析链路
- 阅读 `spider` 包：`task.go`、`request.go`、`parse.go` 等文件，掌握 `Task`、`RuleTree`、`ParseResult` 数据结构。
- 结合 `parse/minimal` 与 `parse/doubanbook`，熟悉规则编排、字段提取、`ctx.AddQueue`/`ctx.Output` 的使用方式。
- 了解 JS 规则：查看 `parse/doubangroupjs` 与 `engine.Store.AddJSTask`，理解使用 Otto 动态扩展的范式。

## 阶段 5：存储、限速、代理等横切能力
- `storage/sqlstorage`：掌握动态建表、批量 flush、`engine.GetFields` 反射字段的策略。
- `limiter` 与 `proxy`：看懂多令牌桶组合、代理轮询与请求头策略，评估与现有游戏后端限速/熔断体系的差异。
- 建议输出：把数据写入 MySQL 后，编写一段 SQL 验证表结构与数据质量。

## 阶段 6：运行与运维
- **日志**：查看 `log` 包封装，了解如何自定义 zap 输出、如何在部署环境中采集。
- **可观测性**：默认开启 `net/http/pprof`；必要时可接入 Prometheus（扩展点）。
- **测试**：先运行 `make cover`（单测 + 覆盖率），再考虑针对 parse/limiter 编写表驱动测试。
- **CI/CD**：阅读 `.github/` workflow（若存在），确保 lint、测试与镜像构建符合预期。

## 阶段 7：能力验证
- **新增任务**：参考 `docs/new_task_minimal.md` 与 `example_baidu_home`，从 0 到 1 实现一个新站点的抓取规则。
- **分布式演练**：
  - 启动多实例 Worker，观察 master 在 etcd 上的任务分配。
  - 模拟 Worker 下线（`Ctrl+C`），确认 master 的 `reAssign` 能及时迁移任务。
- **性能评估**：通过配置不同的 `Limits`、代理池长度，记录吞吐与失败率，形成运行手册。

## 资料速查表
- 架构脉络：`docs/codex/architecture_overview.md`。
- 官方课程背景：README 顶部链接。
- 学习辅助：`docs/familiarization_plan.md`（5 日学习计划）、`docs/new_task_minimal.md`（新任务 SOP）。
- 常用命令：`make build | lint | cover | debug`；`go run . master|worker`；`docker-compose up`。

> 建议随笔记同步绘制属于自己的架构图（可参考 `docs/architecture.svg`），并在每个阶段结束时回到主流程图，确认“输入 – 处理 – 输出”的理解是否闭环。

