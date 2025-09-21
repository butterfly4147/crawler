# 项目熟悉 5 天计划（可按需压缩）

Day 1: 快速通读与本地运行
- 阅读 docs/reading_guide.md 指定路线，理解 main/cmd/worker/master/engine 的边界
- 准备依赖：MySQL 5.7、etcd 3.5（推荐直接 docker-compose up）
- 本地运行 worker + master，确认服务启动成功（HTTP/GRPC 端口连通）

Day 2: Engine 工作流与调度
- 深入 engine/schedule.go（Schedule/Push/Pull/CreateWork/HandleResult）
- 打断点观察 Seeds -> Root -> scheduler 队列与 Visited 去重、失败重试（SetFailure）

Day 3: 解析规则与数据流
- 阅读 parse/doubanbook、parse/doubangroup、parse/doubangroupjs（JS 规则）
- 掌握 spider.RuleTree/Rule/Request/Context/ParseResult 的协作模型

Day 4: 存储与限速
- 阅读 storage/sqlstorage，了解 DataCell 的落库方式与表命名（Task 名作为表名）
- 理解 limiter.Multi 与 rate.Limiter，掌握在 config.toml 中配置 [Tasks].Limits 的方法

Day 5: 动手扩展
- 参考 docs/new_task_minimal.md 新增并运行一个任务（本次已内置 example_baidu_home）
- 思考如何按需扩展：代理池、告警、分布式调度、数据清洗/导出

里程碑检查
- 能清晰叙述一次端到端任务从 Root 到入库的完整路径
- 能独立新增一个任务并通过配置启停
- 能定位常见问题：代理/超时、限速、任务名不匹配、抓取体积过小（<6000 字节被视为失败）