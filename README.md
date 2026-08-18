# gantt-generator

一个使用 Go 编写的轻量级项目排期与甘特图服务。服务提供项目、任务、关键路径、进度报告和 PDF 甘特图接口，项目数据以本地 JSON 文件保存。

## 目录结构

- `models/`：项目、任务、排期和报告的数据模型。
- `storage/`：基于 JSON 文件的项目数据读写。
- `scheduler/`：任务依赖拓扑、关键路径和日期计算。
- `handlers/`：HTTP API 处理逻辑。
- `report/`：进度报告生成与文本格式化。
- `pdfgantt/`：PDF 甘特图渲染。
- `web/`：前端静态页面。

## 运行

要求 Go 1.23 或更高版本。

```bash
go build ./...
go test ./...
go run .
```

服务默认监听 `http://localhost:8080`，数据保存到运行目录下的 `data/projects.json`。

## 常用接口

- `GET /api/projects`：列出项目。
- `POST /api/projects`：创建项目。
- `GET /api/projects/{id}/critical-path`：计算关键路径。
- `GET /api/projects/{id}/report`：获取项目进度报告。
- `GET /api/projects/{id}/gantt-pdf`：生成 PDF 甘特图。
