概览
- 这些补充图专门体现分布式：主从选举、注册发现、资源分配/重分配、限流与熔断、HTTP-Gateway 到 gRPC 的转换，以及 Kubernetes 部署拓扑。
- 对应代码位置：
  - 选举与领导者：`master/master.go` 的 `Campaign()` 使用 Etcd `/crawler/election`；`IsLeader()` 判断是否为 Leader。
  - 资源路径：`master/master.go` 的 `RESOURCEPATH = "/resources"`；`AddResource`/`DeleteResource` 写入/删除 Etcd；`engine/schedule.go` 的 `watchResource()`/`loadResource()` 按 worker 监听并拉取任务。
  - 注册发现：worker 服务名为 `go.micro.server.worker`（`cmd/worker/worker.go`）；`Master.WatchWorker()`/`updateWorkNodes()` 通过 Etcd Registry 监控并维护节点；`Assign()`/`reAssign()` 基于负载进行分配与重分配。
  - 网关与调用防护：`cmd/master/master.go` 中 `RunHTTPServer()` 提供 gRPC-Gateway；`RunGRPCServer()` 通过 `hystrix.NewClientWrapper()` 与 `ratelimit` 包装 gRPC 客户端调用。
  - K8s 暴露端口：`kubernetes/crawl-master-service.yaml` 与 `kubernetes/crawl-worker-service.yaml` 暴露 HTTP 与 gRPC 端口，支撑水平扩展。

如何阅读图
- `distributed_overview.puml`：从组件与集群视角展示 Master/Worker/Etcd 的关系、HTTP->gRPC、熔断/限流包裹、资源路径及服务发现。
- `distributed_flows.puml`：包含三页时序——Leader 选举、添加资源并下发到 Worker、Worker 节点变化触发重分配。

渲染方式
- 若你已安装 PlantUML，本目录下运行：
  - `plantuml distributed_overview.puml` 生成总览图 PNG/SVG
  - `plantuml distributed_flows.puml` 生成三页时序图
- 或在 VSCode + PlantUML 扩展中直接预览。

在现有图中的体现建议
- 组件图：明确标注 Etcd 的两个作用（选举与注册发现）与具体路径 `"/crawler/election"`、`"/resources"`。
- 部署图：将 Master/Worker 作为 Pod，旁挂 Service；标注 HTTP 与 gRPC 端口；Etcd 作为外部或集群内服务。
- 类图：在 `Master` 与 `Crawler`（Engine）旁加注释：前者负责选举/分配，后者监听/执行。
- 时序图：补充 Hystrix/RateLimit 作为调用包裹，体现故障隔离与背压。