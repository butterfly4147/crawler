# 最小可运行的新增任务示例

本示例已在代码中内置：parse/minimal/minimal.go，任务名 example_baidu_home，用于抓取 https://baidu.com 的 <title> 并入库。

一、任务结构速览
- spider.Task
  - Options：Name/Reload/WaitTime/MaxDepth/Limit/Fetcher/Storage
  - Rule（RuleTree）：
    - Root() -> []*Request 初始请求
    - Trunk[ruleName] -> ParseFunc(ctx) 解析函数
- ParseResult：Requesrts（后续请求）+ Items（结果项，建议 *spider.DataCell）

二、如何从零新增一个任务（已为你完成）
1) 新建规则文件 parse/minimal/minimal.go，定义 ExampleTask：
   - Root: 只发起 1 个 GET https://baidu.com，RuleName = parse_title
   - Trunk: parse_title 解析 <title>，输出 DataCell
2) 在 engine/schedule.go 的 init() 中注册：Store.Add(minimal.ExampleTask)
3) 在 config.toml 的 Tasks 中加入同名项（我们已替换占位项）：
   {Name = "example_baidu_home", WaitTime = 1, Reload = true, MaxDepth = 1, Fetcher = "browser", Limits=[{EventCount=1, EventDur=1, Bucket=1}]}
   注意：必须配置 Limits，否则 Request.Fetch() 中的限速器为 nil 会 panic。

三、运行与验证
- 构建：go build -o crawler .
- 启动依赖（推荐）：docker-compose up -d（MySQL/etcd）
- 启动 Worker：.\crawler worker --id=1 --http=:8080 --grpc=:9090
- 启动 Master：.\crawler master --id=2 --http=:8082 --grpc=:9092
- 观察日志：应看到 example_baidu_home 的拉取与 parse_title 的输出；HandleResult 会尝试存储到 MySQL。

四、常见坑
- 页面体积过小：engine.CreateWork 中对 body 长度有校验（<6000 视为失败并重试），通常 https://baidu.com 首页体积充足，可通过此校验。
- 任务名不匹配：config.toml 的 Name 必须与 Store 中注册的 Task.Name 一致。
- 限速未配置：必须配置 Limits；或在 Task.Options.Limit 中内置一个 limiter.Multi。

五、举一反三
- 你可以将 Root 扩展为多个种子 URL，并在 ParseFunc 中通过正则/DOM 选择器生成后续请求（Depth+1，RuleName 指向下个规则）。
- 如果站点需要 Cookie 或代理，分别在 [Tasks] 的 Cookie 字段与 [fetcher].proxy 中配置。