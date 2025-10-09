# 项目 UML 架构图（GPT-5）

本目录包含使用 PlantUML 编写的架构图，帮助你快速熟悉该分布式爬虫项目的核心模块与运行流程。

包含内容：
- 组件总览：`component_overview.puml`
- 核心类关系：`class_diagram.puml`
- 部署架构：`deployment_diagram.puml`
- 主从协作时序：`sequence_master_worker.puml`

渲染方式（任选其一）：
- VSCode 安装 PlantUML 插件，直接预览 `.puml` 文件。
- 使用命令行（需要 Java 与 Graphviz）：
  - 将 `plantuml.jar` 放到本机，执行：
  - `java -jar plantuml.jar docs/uml-gpt5/*.puml`

图示约定：
- 以包名/文件夹映射模块；以结构体/接口映射类关系；箭头表示依赖或调用；粗体标题为场景名。