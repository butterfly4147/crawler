# Windows 快速上手（PowerShell）

适用：Windows 11 开发环境，建议使用 Docker 以简化依赖。

## 方式一：Docker Compose 一键启动（推荐）
```powershell
cd d:\data\repos\repo_reference\crawler
docker-compose up -d
```
- 端口：
  - Worker：`http://localhost:8080`、`grpc :9090`
  - Master：`http://localhost:8082`、`grpc :9092`
  - etcd：`http://localhost:2379`
- 验证：
  - Worker Greeter 示例：
    ```powershell
    Invoke-WebRequest -Uri "http://localhost:8080/greeter/hello" -Method POST -Body "name=test" -ContentType "application/x-www-form-urlencoded"
    ```
  - MySQL 表：
    ```powershell
    mysql -u root -p123456 -e "USE crawler; SHOW TABLES;"
    ```

## 方式二：本机编译与运行（需安装 Go/etcd/MySQL）
```powershell
cd d:\data\repos\repo_reference\crawler
go build -o crawler .
.\n+crawler master --id=3 --http=:8082 --grpc=:9092
.
crawler worker --id=2 --http=:8080 --grpc=:9090
```
- 依赖：
  - etcd：确保 `2379/2380` 可访问（先启动 etcd）
  - MySQL：`storage.sqlURL` 要与本机数据库一致（默认 `root:123456@tcp(127.0.0.1:3306)/crawler?charset=utf8`）
  - 配置：可直接使用项目根目录的 `config.toml`

## 常见问题
- 首次运行未写入数据：`SQLStorage` 默认批量数（`BatchCount`）≥2，需累积到阈值才 `Flush`；可调低批量或增加示例任务的数据量。
- 无法连接 etcd/MySQL：检查容器健康状态与本机端口占用；Windows 防火墙是否阻断。
- Seeds 未生效：检查 `config.toml [Tasks]` 的 `Name` 是否与 `engine.Store` 中注册的任务一致。