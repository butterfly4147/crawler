# 分布式爬虫项目运行演示 SOP 文档

## 项目概述
这是一个基于Go语言开发的分布式爬虫系统，采用Master-Worker架构，支持动态任务调度和数据存储。

## 前置条件

### 1. 环境要求
- **操作系统**: Windows 11
- **Go版本**: 1.19+
- **数据库**: MySQL 5.7+
- **服务发现**: etcd 3.5+

### 2. 依赖服务安装
```powershell
# 安装MySQL（如果未安装）
# 下载并安装MySQL 5.7+，设置root密码为: 123456

# 创建数据库
mysql -u root -p123456 -e "CREATE DATABASE IF NOT EXISTS crawler DEFAULT CHARSET utf8;"
```

## 快速启动步骤

### 第一步：编译项目
```powershell
# 进入项目目录
cd d:\data\repos\repo_reference\crawler

# 编译项目
go build -o crawler.exe .
```

### 第二步：启动服务（按顺序执行）

#### 1. 启动etcd服务注册中心
```powershell
# 新开终端1
etcd --data-dir=default.etcd --listen-client-urls=http://127.0.0.1:2379 --advertise-client-urls=http://127.0.0.1:2379
```

#### 2. 启动Master服务
```powershell
# 新开终端2
.\crawler.exe master --id=2 --http=:8082 --grpc=:9092 --pprof=:9982
```

#### 3. 启动Worker服务
```powershell
# 新开终端3
.\crawler.exe worker --id=1 --http=:8080 --grpc=:9090
```

### 第三步：验证服务状态

#### 1. 检查服务是否正常启动
```powershell
# 测试Worker服务HTTP接口
Invoke-WebRequest -Uri "http://localhost:8080/greeter/hello" -Method POST -Body "name=test" -ContentType "application/x-www-form-urlencoded"

# 预期返回: {"greeting":"Hello test"}
```

#### 2. 检查数据库表是否创建
```powershell
mysql -u root -p123456 -e "USE crawler; SHOW TABLES;"
# 应该看到 example_baidu_home 表
```

## 演示功能操作

### 1. 触发爬虫任务
```powershell
# 方法1：通过Master服务API添加任务
Invoke-WebRequest -Uri "http://localhost:8082/crawler/resource" -Method POST -Body '{"name":"example_baidu_home"}' -ContentType "application/json"

# 预期返回: {"id":"go.micro.server.worker-1", "Address":"192.168.0.108:9090"}
```

### 2. 查看爬虫结果
```powershell
# 查看数据库中的爬虫数据
mysql -u root -p123456 -e "USE crawler; SELECT * FROM example_baidu_home ORDER BY id DESC LIMIT 5;"
```

### 3. 批量触发演示
```powershell
# 由于系统使用批量存储机制(BatchCount=2)，需要触发至少2次任务才能看到数据入库
# 第一次触发
Invoke-WebRequest -Uri "http://localhost:8082/crawler/resource" -Method POST -Body '{"name":"example_baidu_home"}' -ContentType "application/json"

# 第二次触发
Invoke-WebRequest -Uri "http://localhost:8082/crawler/resource" -Method POST -Body '{"name":"example_baidu_home"}' -ContentType "application/json"

# 查看结果
mysql -u root -p123456 -e "USE crawler; SELECT COUNT(*) as total_records FROM example_baidu_home;"
```

## 预期演示结果

### 1. 服务日志输出
- **etcd**: 显示服务注册信息
- **Master**: 显示任务调度和资源分配日志
- **Worker**: 显示爬虫执行和数据存储日志

### 2. 数据库结果
```
+----+-----------------------------+--------------------+---------------------+
| id | title                       | URL                | Time                |
+----+-----------------------------+--------------------+---------------------+
|  2 | 百度一下，你就知道          | https://baidu.com/ | 2025-09-21 18:25:07 |
|  1 | 百度一下，你就知道          | 2025-09-21 18:21:55 |
+----+-----------------------------+--------------------+---------------------+
```

### 3. 关键日志示例
```
# Worker日志 - 任务接收
{"level":"INFO","msg":"receive create resource","spec":{"Name":"example_baidu_home"}}

# Worker日志 - 数据解析
{"level":"INFO","msg":"get result: &{map[Data:map[title:百度一下，你就知道] Rule:parse_title Task:example_baidu_home]}"}

# Worker日志 - 数据库插入
{"level":"DEBUG","msg":"insert table","sql":"INSERT INTO example_baidu_home(title,URL,Time) VALUES (?,?,?),(?,?,?);"}
```

## 故障排除

### 1. 常见问题
- **数据库连接失败**: 检查MySQL服务是否启动，密码是否正确
- **etcd连接失败**: 确保etcd服务正常运行在2379端口
- **端口占用**: 检查8080、8082、9090、9092端口是否被占用

### 2. 调试命令
```powershell
# 检查端口占用
netstat -ano | findstr :8080
netstat -ano | findstr :8082

# 检查MySQL连接
mysql -u root -p123456 -e "SELECT 1;"

# 检查etcd状态
curl http://127.0.0.1:2379/health
```

### 3. 重置环境
```powershell
# 停止所有服务 (Ctrl+C)
# 清理etcd数据
Remove-Item -Recurse -Force default.etcd

# 清理数据库
mysql -u root -p123456 -e "DROP DATABASE IF EXISTS crawler; CREATE DATABASE crawler DEFAULT CHARSET utf8;"
```

## 扩展功能

### 1. 添加新的爬虫任务
编辑 `config.toml` 文件，在 `[[Tasks]]` 部分添加新任务配置。

### 2. 监控服务状态
- Master服务监控: http://localhost:8082
- Worker服务监控: http://localhost:8080
- 性能分析: http://localhost:9982/debug/pprof/

### 3. 集群部署
参考 `docker-compose.yml` 文件进行容器化部署。

## 注意事项

1. **批量存储机制**: 系统默认BatchCount=2，需要至少2条数据才会触发数据库写入
2. **服务启动顺序**: 必须按照 etcd → Master → Worker 的顺序启动
3. **数据库权限**: 确保MySQL用户有创建表和插入数据的权限
4. **防火墙设置**: 确保相关端口未被防火墙阻止

## 联系支持
如遇到问题，请检查：
1. 服务日志输出
2. 数据库连接状态  
3. 网络端口状态
4. 配置文件格式

---
**文档版本**: v1.0  
**更新时间**: 2025-09-21  
**适用环境**: Windows 11 + PowerShell