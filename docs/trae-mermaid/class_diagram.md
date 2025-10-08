# 分布式爬虫系统 - 类图

```mermaid
classDiagram
    %% 主控模块
    class Master {
        -ID string
        -leaderID string
        -workNodes []string
        -resources map[string]Resource
        +DeleteResource(id string)
        +AddResource(resource Resource)
        +Campaign() error
        +elect() error
    }

    class WorkerNode {
        -ID string
        -masterID string
        -status NodeStatus
        -tasks []Task
        +RegisterToMaster() error
        +ExecuteTask(task Task) error
        +ReportStatus() error
    }

    %% 爬虫引擎
    class Crawler {
        -name string
        -startUrls []string
        -options Options
        +Run() error
        +Stop()
        +AddTask(task Task)
    }

    class Schedule {
        -requestCh chan Request
        -workerCh chan Request
        -priReqQueue []*Request
        -reqQueue []*Request
        +Schedule()
        +Push(reqs ...*Request)
        +Pull() *Request
        +NRequest() int
        +NWorker() int
    }

    class Scheduler {
        <<interface>>
        +Schedule()
        +Push(reqs ...*Request)
        +Pull() *Request
    }

    class CrawlerStore {
        -list []*Task
        -Hash map[string]*Task
        +Add(task Task)
        +AddJSTask(m map[string]any)
    }

    %% 爬虫任务
    class Task {
        +Property Property
        +Visited map[string]bool
        +Rule RuleTree
        +Options Options
        +NewTask(opts TaskOptions) *Task
    }

    class Property {
        +Name string
        +Namespace string
        +Root string
        +Retry int
        +Interval time.Duration
    }

    class Request {
        +URL *url.URL
        +Method string
        +Depth int
        +Priority int
        +RuleName string
        +TempData *Temp
    }

    class RuleTree {
        +Root func() ([]*Request, error)
        +Trunk map[string]*Rule
    }

    class Rule {
        +ParseFunc func(*Context) (ParseResult, error)
        +ItemFields []string
        +Namespace string
    }

    %% 网络模块
    class Fetcher {
        <<interface>>
        +Get(url string) ([]byte, error)
    }

    class HTTPClient {
        -client *http.Client
        -timeout time.Duration
        +Get(url string) ([]byte, error)
        +Post(url string, data []byte) ([]byte, error)
    }

    %% 存储模块
    class Storage {
        <<interface>>
        +Save(items []any) error
    }

    class MySQLStorage {
        -db *sql.DB
        -options Options
        +Save(items []any) error
        +CreateTable(tableName string) error
    }

    class FileStorage {
        -basePath string
        -format string
        +Save(items []any) error
        +WriteToFile(filename string, data []byte) error
    }

    %% 代理模块
    class ProxyProvider {
        <<interface>>
        +GetProxy() string
        +ReleaseProxy(proxy string)
    }

    class ProxyFunc {
        +RoundRobin(req *Request) (*url.URL, error)
        +Random(req *Request) (*url.URL, error)
    }

    class ProxyPool {
        -proxies []string
        -current int
        -mutex sync.Mutex
        +GetProxy() string
        +ReleaseProxy(proxy string)
    }

    %% 限流模块
    class RateLimiter {
        <<interface>>
        +Allow() bool
        +Wait()
    }

    class TokenBucket {
        -limiter *rate.Limiter
        +Wait(ctx context.Context) error
        +Allow() bool
    }

    %% 关系定义
    Master "1" -- "many" WorkerNode : manages
    Master "1" -- "many" Task : distributes
    WorkerNode "1" -- "1" Crawler : runs
    Crawler "1" -- "1" Schedule : uses
    Crawler "1" -- "1" CrawlerStore : uses
    Schedule ..|> Scheduler : implements
    CrawlerStore "1" -- "many" Task : stores
    Task "1" -- "1" Property : has
    Task "1" -- "1" RuleTree : has
    Task "1" -- "many" Request : generates
    RuleTree "1" -- "many" Rule : contains
    Crawler "1" -- "1" Fetcher : uses
    HTTPClient ..|> Fetcher : implements
    Crawler "1" -- "1" Storage : uses
    MySQLStorage ..|> Storage : implements
    FileStorage ..|> Storage : implements
    Crawler "1" -- "1" ProxyProvider : uses
    ProxyFunc ..|> ProxyProvider : implements
    ProxyPool ..|> ProxyProvider : implements
    Crawler "1" -- "1" RateLimiter : uses
    TokenBucket ..|> RateLimiter : implements
```

## 类图说明

### 主控模块
- **Master**: 主控节点，负责任务分发和集群管理
- **WorkerNode**: 工作节点，执行具体的爬虫任务

### 爬虫引擎
- **Crawler**: 爬虫核心引擎，协调各个组件
- **Schedule**: 任务调度器，管理请求队列
- **CrawlerStore**: 任务存储器，缓存爬虫任务

### 爬虫任务
- **Task**: 爬虫任务实体，包含配置和规则
- **Property**: 任务属性配置
- **Request**: HTTP请求对象
- **RuleTree**: 解析规则树
- **Rule**: 具体的解析规则

### 网络模块
- **Fetcher**: 网络请求接口
- **HTTPClient**: HTTP客户端实现

### 存储模块
- **Storage**: 存储接口
- **MySQLStorage**: MySQL存储实现
- **FileStorage**: 文件存储实现

### 代理模块
- **ProxyProvider**: 代理提供者接口
- **ProxyFunc**: 代理函数实现
- **ProxyPool**: 代理池实现

### 限流模块
- **RateLimiter**: 限流器接口
- **TokenBucket**: 令牌桶限流实现