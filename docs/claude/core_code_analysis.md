# 核心代码分析

## 1. 项目入口和命令结构

### 1.1 主程序入口 (`main.go`)

```go
package main

import (
    "github.com/dreamerjackson/crawler/cmd"
    _ "net/http/pprof"
)

func main() {
    cmd.Execute()
}
```

**分析：**
- 使用cobra命令行框架
- 集成pprof性能分析
- 委托给cmd模块执行具体命令

### 1.2 命令定义 (`cmd/cmd.go`)

```go
func Execute() {
    var rootCmd = &cobra.Command{Use: "crawler"}

    rootCmd.AddCommand(
        master.MasterCmd,    // Master服务命令
        worker.WorkerCmd,    // Worker服务命令
        versionCmd,          // 版本命令
    )

    _ = rootCmd.Execute()
}
```

**分析：**
- 支持三个子命令：master、worker、version
- 采用模块化设计，各命令在独立包中实现

## 2. Master服务核心逻辑

### 2.1 Master选举机制 (`master/master.go`)

```go
func (m *Master) Campaign() {
    s, err := concurrency.NewSession(m.etcdCli, concurrency.WithTTL(5))
    e := concurrency.NewElection(s, "/crawler/election")

    // 选举过程
    go m.elect(e, leaderCh)

    // 监听选举变化
    leaderChange := e.Observe(context.Background())

    for {
        select {
        case err := <-leaderCh:
            if err == nil {
                m.BecomeLeader()  // 成为Leader
            }
        case resp := <-leaderChange:
            m.leaderID = string(resp.Kvs[0].Value)  // 更新Leader ID
        }
    }
}
```

**关键点：**
- 使用etcd的concurrency包实现分布式选举
- 通过Session TTL保证节点故障时自动重新选举
- 支持Leader故障时的自动切换

### 2.2 负载均衡算法

```go
func (m *Master) Assign(r *ResourceSpec) (*NodeSpec, error) {
    candidates := make([]*NodeSpec, 0, len(m.workNodes))

    for _, node := range m.workNodes {
        candidates = append(candidates, node)
    }

    // 按负载排序，选择负载最低的节点
    sort.Slice(candidates, func(i, j int) bool {
        return candidates[i].Payload < candidates[j].Payload
    })

    if len(candidates) > 0 {
        return candidates[0], nil
    }

    return nil, errors.New("no worker nodes")
}
```

**算法特点：**
- 最小负载优先策略
- 时间复杂度：O(n log n)
- 保证任务均匀分布

## 3. Worker服务核心逻辑

### 3.1 爬虫引擎初始化 (`cmd/worker/worker.go`)

```go
func Run() {
    // 配置加载
    cfg, err := config.NewConfig(config.WithReader(json.NewReader(reader.WithEncoder(enc))))
    err = cfg.Load(file.NewSource(file.WithPath("config.toml")))

    // 组件初始化
    f := &collect.BrowserFetch{Timeout: timeout, Logger: logger, Proxy: p}
    storage, err = sqlstorage.New(sqlstorage.WithSQLURL(sqlURL))

    // 引擎创建
    s, err := engine.NewEngine(
        engine.WithFetcher(f),
        engine.WithLogger(logger),
        engine.WithWorkCount(5),
        engine.WithSeeds(seeds),
        engine.WithScheduler(engine.NewSchedule()),
        engine.WithStorage(storage),
    )

    // 启动服务
    go s.Run(id, cluster)
    RunGRPCServer(logger, sconfig)
}
```

**架构特点：**
- 依赖注入模式，通过Option函数配置组件
- 支持插件化架构，可替换Fetcher、Storage等组件
- 并发控制，通过WorkCount控制worker数量

### 3.2 任务解析配置

```go
func ParseTaskConfig(logger *zap.Logger, f spider.Fetcher, s spider.Storage, cfgs []spider.TaskConfig) []*spider.Task {
    tasks := make([]*spider.Task, 0, 1000)

    for _, cfg := range cfgs {
        t := spider.NewTask(
            spider.WithName(cfg.Name),
            spider.WithReload(cfg.Reload),
            spider.WithCookie(cfg.Cookie),
            spider.WithLogger(logger),
            spider.WithStorage(s),
        )

        // 限流器配置
        var limits []limiter.RateLimiter
        for _, lcfg := range cfg.Limits {
            l := rate.NewLimiter(limiter.Per(lcfg.EventCount, time.Duration(lcfg.EventDur)*time.Second), lcfg.Bucket)
            limits = append(limits, l)
        }

        // Fetcher配置
        switch cfg.Fetcher {
        case "browser":
            t.Fetcher = f
        }

        tasks = append(tasks, t)
    }
    return tasks
}
```

**配置特性：**
- 支持多种Fetcher类型（目前只有browser）
- 灵活的限流器配置
- 任务级别的配置隔离

## 4. 爬虫引擎核心

### 4.1 调度器实现 (`engine/schedule.go`)

```go
func (s *Schedule) Schedule() {
    var ch chan *spider.Request
    var req *spider.Request

    for {
        // 优先级队列处理
        if req == nil && len(s.priReqQueue) > 0 {
            req = s.priReqQueue[0]
            s.priReqQueue = s.priReqQueue[1:]
            ch = s.workerCh
        }

        // 普通队列处理
        if req == nil && len(s.reqQueue) > 0 {
            req = s.reqQueue[0]
            s.reqQueue = s.reqQueue[1:]
            ch = s.workerCh
        }

        select {
        case r := <-s.requestCh:  // 接收新请求
            if r.Priority > 0 {
                s.priReqQueue = append(s.priReqQueue, r)
            } else {
                s.reqQueue = append(s.reqQueue, r)
            }
        case ch <- req:  // 分发请求给worker
            req = nil
            ch = nil
        }
    }
}
```

**调度策略：**
- 双队列设计：优先级队列 + 普通队列
- 非阻塞调度，避免goroutine阻塞
- 支持请求优先级

### 4.2 Worker执行逻辑

```go
func (c *Crawler) CreateWork() {
    defer func() {
        if err := recover(); err != nil {
            c.Logger.Error("worker panic", zap.Any("err", err))
        }
    }()

    for {
        req := c.scheduler.Pull()

        // 请求校验
        if err := req.Check(); err != nil {
            continue
        }

        // 去重检查
        if !req.Task.Reload && c.HasVisited(req) {
            continue
        }

        c.StoreVisited(req)

        // 执行抓取
        body, err := req.Fetch()
        if err != nil {
            c.SetFailure(req)
            continue
        }

        // 解析内容
        rule := req.Task.Rule.Trunk[req.RuleName]
        ctx := &spider.Context{Body: body, Req: req}
        result, err := rule.ParseFunc(ctx)

        // 处理结果
        if len(result.Requesrts) > 0 {
            go c.scheduler.Push(result.Requesrts...)
        }

        c.out <- result
    }
}
```

**执行流程：**
1. 从调度器获取请求
2. 请求校验和去重检查
3. 执行HTTP抓取
4. 解析网页内容
5. 发现新链接并加入队列
6. 输出解析结果

## 5. 分布式协调

### 5.1 资源监听

```go
func (c *Crawler) watchResource() {
    watch := c.etcdCli.Watch(context.Background(), master.RESOURCEPATH,
        clientv3.WithPrefix(), clientv3.WithPrevKV())

    for w := range watch {
        for _, ev := range w.Events {
            switch ev.Type {
            case clientv3.EventTypePut:
                spec, _ := master.Decode(ev.Kv.Value)
                c.runTasks(spec.Name)  // 启动任务
            case clientv3.EventTypeDelete:
                spec, _ := master.Decode(ev.PrevKv.Value)
                c.deleteTasks(spec.Name)  // 停止任务
            }
        }
    }
}
```

**协调机制：**
- 通过etcd Watch机制监听资源变化
- 支持任务的动态添加和删除
- 保证集群状态一致性

### 5.2 服务注册发现

```go
func (m *Master) WatchWorker() chan *registry.Result {
    watch, err := m.registry.Watch(registry.WatchService(ServiceName))

    ch := make(chan *registry.Result)
    go func() {
        for {
            res, err := watch.Next()
            if err != nil {
                continue
            }
            ch <- res
        }
    }()
    return ch
}
```

**发现机制：**
- 使用go-micro的服务发现功能
- 实时监控Worker节点状态变化
- 支持节点的动态加入和离开

## 6. 关键设计模式

### 6.1 Option模式

```go
type Option func(*options)

func WithLogger(logger *zap.Logger) Option {
    return func(o *options) {
        o.Logger = logger
    }
}

func WithWorkCount(workCount int) Option {
    return func(o *options) {
        o.WorkCount = workCount
    }
}
```

**优势：**
- 灵活的配置方式
- 支持默认值和自定义配置
- 良好的可扩展性

### 6.2 接口隔离

```go
type Fetcher interface {
    Get(url *Request) ([]byte, error)
}

type Storage interface {
    Save(data *spider.DataCell) error
}

type Scheduler interface {
    Schedule()
    Push(...*spider.Request)
    Pull() *spider.Request
}
```

**设计原则：**
- 依赖接口而非实现
- 支持组件替换
- 便于测试和扩展

## 总结

这个项目的代码设计体现了以下优秀实践：

1. **清晰的架构分层**：Master-Worker架构明确
2. **完善的错误处理**：panic恢复、错误重试
3. **高效的并发模型**：goroutine + channel
4. **灵活的配置系统**：Option模式 + 配置文件
5. **强大的扩展能力**：插件化架构 + 接口设计

对于有Go语言经验的开发者，这个项目是学习分布式系统设计和微服务架构的优秀案例。