# Gantt Generator

## 项目说明

这是一个使用 Go 编写的项目排期与甘特图服务。项目以本地 JSON 文件保存项目和任务数据，提供项目管理、任务依赖、关键路径、进度报告以及 PDF 甘特图等能力，同时包含用于查看排期数据的 Web 页面。

主要模块如下：

- `models/`：项目、任务、排期和报告的数据模型。
- `storage/`：基于 JSON 文件的数据持久化。
- `scheduler/`：任务依赖分析、日期计算和关键路径计算。
- `handlers/`：HTTP API 路由和请求处理。
- `report/`：项目进度报告生成。
- `pdfgantt/`：PDF 甘特图生成。
- `web/`：前端静态页面。
- `data/`：运行时项目数据目录，默认使用 `data/projects.json`。

## 环境要求

- Go 1.26.5 或兼容的 Go 1.26 工具链。
- Docker（仅在使用 Benzhi Docker 辅助文件时需要）。

## 标准构建、运行和测试命令

在项目根目录执行：

```bash
# 下载依赖
go mod download

# 构建全部 Go 包
go build ./...

# 运行全部测试
go test ./...

# 启动服务
go run .
```

服务默认监听 `http://localhost:8080`。启动后可以通过 Web 页面或 HTTP API 访问项目数据。

常用接口包括：

- `GET /api/projects`：列出项目。
- `POST /api/projects`：创建项目。
- `GET /api/projects/{id}`：读取项目。
- `PUT /api/projects/{id}`：更新项目。
- `GET /api/projects/{id}/critical-path`：计算关键路径。
- `GET /api/projects/{id}/report`：生成进度报告。
- `GET /api/projects/{id}/gantt-pdf`：生成 PDF 甘特图。

## Benzhi Docker 构建

使用评测专用的 `benzhi.Dockerfile` 构建镜像：

```bash
./build_benzhi_docker.sh [镜像名] [目标平台]
```

参数均可省略：镜像名默认为 `my-project`，目标平台默认为 `linux/amd64`。构建完成后启动交互式容器：

```bash
docker run -it my-project:latest
```

