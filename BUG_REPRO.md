# 缺陷复现说明

## 缺陷现象
导入排期后，概览里明明有两项审批任务，关键路径却只剩一项；后面更新其中一项还会影响另一项。不要修改代码，帮我查清楚这份项目数据为什么会被这样处理。

## 触发方式
# 1. 进入测试项目根目录；预期：后续命令在目标项目内执行。
cd /workplace/gantt-generator
# 2. target-verify 复现重复任务 ID 的导入；预期：基线会接受重复 ID，接口返回成功。
go test -v ./handlers -run '^TestProjectUpdateRejectsDuplicateTaskIDs$' -count=20

## 触发后的实际错误输出

```text
=== RUN   TestProjectUpdateRejectsDuplicateTaskIDs
    duplicate_task_id_test.go:47: 重复任务ID应被拒绝，实际状态码为 200，响应为 {"id":"project-1","name":"数据迁移","tasks":[{"id":"approval","name":"业务审批","start_date":"2026-08-01","end_date":"2026-08-02","dependencies":null,"assignee":"","status":"未开始","progress":0},{"id":"approval","name":"安全审批","start_date":"2026-08-03","end_date":"2026-08-04","dependencies":null,"assignee":"","status":"未开始","progress":0}],"created_at":"2026-08-17T11:59:32+08:00","updated_at":"2026-08-17T11:59:32+08:00"}
--- FAIL: TestProjectUpdateRejectsDuplicateTaskIDs (0.00s)
=== RUN   TestProjectUpdateRejectsDuplicateTaskIDs
    duplicate_task_id_test.go:47: 重复任务ID应被拒绝，实际状态码为 200，响应为 {"id":"project-1","name":"数据迁移","tasks":[{"id":"approval","name":"业务审批","start_date":"2026-08-01","end_date":"2026-08-02","dependencies":null,"assignee":"","status":"未开始","progress":0},{"id":"approval","name":"安全审批","start_date":"2026-08-03","end_date":"2026-08-04","dependencies":null,"assignee":"","status":"未开始","progress":0}],"created_at":"2026-08-17T11:59:32+08:00","updated_at":"2026-08-17T11:59:32+08:00"}
--- FAIL: TestProjectUpdateRejectsDuplicateTaskIDs (0.00s)
=== RUN   TestProjectUpdateRejectsDuplicateTaskIDs
    duplicate_task_id_test.go:47: 重复任务ID应被拒绝，实际状态码为 200，响应为 {"id":"project-1","name":"数据迁移","tasks":[{"id":"approval","name":"业务审批","start_date":"2026-08-01","end_date":"2026-08-02","dependencies":null,"assignee":"","status":"未开始","progress":0},{"id":"approval","name":"安全审批","start_date":"2026-08-03","end_date":"2026-08-04","dependencies":null,"assignee":"","status":"未开始","progress":0}],"created_at":"2026-08-17T11:59:32+08:00","updated_at":"2026-08-17T11:59:32+08:00"}
--- FAIL: TestProjectUpdateRejectsDuplicateTaskIDs (0.00s)
=== RUN   TestProjectUpdateRejectsDuplicateTaskIDs
    duplicate_task_id_test.go:47: 重复任务ID应被拒绝，实际状态码为 200，响应为 {"id":"project-1","name":"数据迁移","tasks":[{"id":"approval","name":"业务审批","start_date":"2026-08-01","end_date":"2026-08-02","dependencies":null,"assignee":"","status":"未开始","progress":0},{"id":"approval","name":"安全审批","start_date":"2026-08-03","end_date":"2026-08-04","dependencies":null,"assignee":"","status":"未开始","progress":0}],"created_at":"2026-08-17T11:59:32+08:00","updated_at":"2026-08-17T11:59:32+08:00"}
--- FAIL: TestProjectUpdateRejectsDuplicateTaskIDs (0.00s)
=== RUN   TestProjectUpdateRejectsDuplicateTaskIDs
    duplicate_task_id_test.go:47: 重复任务ID应被拒绝，实际状态码为 200，响应为 {"id":"project-1","name":"数据迁移","tasks":[{"id":"approval","name":"业务审批","start_date":"2026-08-01","end_date":"2026-08-02","dependencies":null,"assignee":"","status":"未开始","progress":0},{"id":"approval","name":"安全审批","start_date":"2026-08-03","end_date":"2026-08-04","dependencies":null,"assignee":"","status":"未开始","progress":0}],"created_at":"2026-08-17T11:59:32+08:00","updated_at":"2026-08-17T11:59:32+08:00"}
--- FAIL: TestProjectUpdateRejectsDuplicateTaskIDs (0.00s)
=== RUN   TestProjectUpdateRejectsDuplicateTaskIDs
    duplicate_task_id_test.go:47: 重复任务ID应被拒绝，实际状态码为 200，响应为 {"id":"project-1","name":"数据迁移","tasks":[{"id":"approval","name":"业务审批","start_date":"2026-08-01","end_date":"2026-08-02","dependencies":null,"assignee":"","status":"未开始","progress":0},{"id":"approval","name":"安全审批","start_date":"2026-08-03","end_date":"2026-08-04","dependencies":null,"assignee":"","status":"未开始","progress":0}],"created_at":"2026-08-17T11:59:32+08:00","updated_at":"2026-08-17T11:59:32+08:00"}
--- FAIL: TestProjectUpdateRejectsDuplicateTaskIDs (0.00s)
=== RUN   TestProjectUpdateRejectsDuplicateTaskIDs
    duplicate_task_id_test.go:47: 重复任务ID应被拒绝，实际状态码为 200，响应为 {"id":"project-1","name":"数据迁移","tasks":[{"id":"approval","name":"业务审批","start_date":"2026-08-01","end_date":"2026-08-02","dependencies":null,"assignee":"","status":"未开始","progress":0},{"id":"approval","name":"安全审批","start_date":"2026-08-03","end_date":"2026-08-04","dependencies":null,"assignee":"","status":"未开始","progress":0}],"created_at":"2026-08-17T11:59:32+08:00","updated_at":"2026-08-17T11:59:32+08:00"}
--- FAIL: TestProjectUpdateRejectsDuplicateTaskIDs (0.00s)
=== RUN   TestProjectUpdateRejectsDuplicateTaskIDs
    duplicate_task_id_test.go:47: 重复任务ID应被拒绝，实际状态码为 200，响应为 {"id":"project-1","name":"数据迁移","tasks":[{"id":"approval","name":"业务审批","start_date":"2026-08-01","end_date":"2026-08-02","dependencies":null,"assignee":"","status":"未开始","progress":0},{"id":"approval","name":"安全审批","start_date":"2026-08-03","end_date":"2026-08-04","dependencies":null,"assignee":"","status":"未开始","progress":0}],"created_at":"2026-08-17T11:59:32+08:00","updated_at":"2026-08-17T11:59:32+08:00"}
--- FAIL: TestProjectUpdateRejectsDuplicateTaskIDs (0.00s)
=== RUN   TestProjectUpdateRejectsDuplicateTaskIDs
    duplicate_task_id_test.go:47: 重复任务ID应被拒绝，实际状态码为 200，响应为 {"id":"project-1","name":"数据迁移","tasks":[{"id":"approval","name":"业务审批","start_date":"2026-08-01","end_date":"2026-08-02","dependencies":null,"assignee":"","status":"未开始","progress":0},{"id":"approval","name":"安全审批","start_date":"2026-08-03","end_date":"2026-08-04","dependencies":null,"assignee":"","status":"未开始","progress":0}],"created_at":"2026-08-17T11:59:32+08:00","updated_at":"2026-08-17T11:59:32+08:00"}
--- FAIL: TestProjectUpdateRejectsDuplicateTaskIDs (0.00s)
=== RUN   TestProjectUpdateRejectsDuplicateTaskIDs
    duplicate_task_id_test.go:47: 重复任务ID应被拒绝，实际状态码为 200，响应为 {"id":"project-1","name":"数据迁移","tasks":[{"id":"approval","name":"业务审批","start_date":"2026-08-01","end_date":"2026-08-02","dependencies":null,"assignee":"","status":"未开始","progress":0},{"id":"approval","name":"安全审批","start_date":"2026-08-03","end_date":"2026-08-04","dependencies":null,"assignee":"","status":"未开始","progress":0}],"created_at":"2026-08-17T11:59:32+08:00","updated_at":"2026-08-17T11:59:32+08:00"}
--- FAIL: TestProjectUpdateRejectsDuplicateTaskIDs (0.00s)
=== RUN   TestProjectUpdateRejectsDuplicateTaskIDs
    duplicate_task_id_test.go:47: 重复任务ID应被拒绝，实际状态码为 200，响应为 {"id":"project-1","name":"数据迁移","tasks":[{"id":"approval","name":"业务审批","start_date":"2026-08-01","end_date":"2026-08-02","dependencies":null,"assignee":"","status":"未开始","progress":0},{"id":"approval","name":"安全审批","start_date":"2026-08-03","end_date":"2026-08-04","dependencies":null,"assignee":"","status":"未开始","progress":0}],"created_at":"2026-08-17T11:59:32+08:00","updated_at":"2026-08-17T11:59:32+08:00"}
--- FAIL: TestProjectUpdateRejectsDuplicateTaskIDs (0.00s)
=== RUN   TestProjectUpdateRejectsDuplicateTaskIDs
    duplicate_task_id_test.go:47: 重复任务ID应被拒绝，实际状态码为 200，响应为 {"id":"project-1","name":"数据迁移","tasks":[{"id":"approval","name":"业务审批","start_date":"2026-08-01","end_date":"2026-08-02","dependencies":null,"assignee":"","status":"未开始","progress":0},{"id":"approval","name":"安全审批","start_date":"2026-08-03","end_date":"2026-08-04","dependencies":null,"assignee":"","status":"未开始","progress":0}],"created_at":"2026-08-17T11:59:32+08:00","updated_at":"2026-08-17T11:59:32+08:00"}
--- FAIL: TestProjectUpdateRejectsDuplicateTaskIDs (0.00s)
=== RUN   TestProjectUpdateRejectsDuplicateTaskIDs
    duplicate_task_id_test.go:47: 重复任务ID应被拒绝，实际状态码为 200，响应为 {"id":"project-1","name":"数据迁移","tasks":[{"id":"approval","name":"业务审批","start_date":"2026-08-01","end_date":"2026-08-02","dependencies":null,"assignee":"","status":"未开始","progress":0},{"id":"approval","name":"安全审批","start_date":"2026-08-03","end_date":"2026-08-04","dependencies":null,"assignee":"","status":"未开始","progress":0}],"created_at":"2026-08-17T11:59:32+08:00","updated_at":"2026-08-17T11:59:32+08:00"}
--- FAIL: TestProjectUpdateRejectsDuplicateTaskIDs (0.00s)
=== RUN   TestProjectUpdateRejectsDuplicateTaskIDs
    duplicate_task_id_test.go:47: 重复任务ID应被拒绝，实际状态码为 200，响应为 {"id":"project-1","name":"数据迁移","tasks":[{"id":"approval","name":"业务审批","start_date":"2026-08-01","end_date":"2026-08-02","dependencies":null,"assignee":"","status":"未开始","progress":0},{"id":"approval","name":"安全审批","start_date":"2026-08-03","end_date":"2026-08-04","dependencies":null,"assignee":"","status":"未开始","progress":0}],"created_at":"2026-08-17T11:59:32+08:00","updated_at":"2026-08-17T11:59:32+08:00"}
--- FAIL: TestProjectUpdateRejectsDuplicateTaskIDs (0.00s)
=== RUN   TestProjectUpdateRejectsDuplicateTaskIDs
    duplicate_task_id_test.go:47: 重复任务ID应被拒绝，实际状态码为 200，响应为 {"id":"project-1","name":"数据迁移","tasks":[{"id":"approval","name":"业务审批","start_date":"2026-08-01","end_date":"2026-08-02","dependencies":null,"assignee":"","status":"未开始","progress":0},{"id":"approval","name":"安全审批","start_date":"2026-08-03","end_date":"2026-08-04","dependencies":null,"assignee":"","status":"未开始","progress":0}],"created_at":"2026-08-17T11:59:32+08:00","updated_at":"2026-08-17T11:59:32+08:00"}
--- FAIL: TestProjectUpdateRejectsDuplicateTaskIDs (0.00s)
=== RUN   TestProjectUpdateRejectsDuplicateTaskIDs
    duplicate_task_id_test.go:47: 重复任务ID应被拒绝，实际状态码为 200，响应为 {"id":"project-1","name":"数据迁移","tasks":[{"id":"approval","name":"业务审批","start_date":"2026-08-01","end_date":"2026-08-02","dependencies":null,"assignee":"","status":"未开始","progress":0},{"id":"approval","name":"安全审批","start_date":"2026-08-03","end_date":"2026-08-04","dependencies":null,"assignee":"","status":"未开始","progress":0}],"created_at":"2026-08-17T11:59:32+08:00","updated_at":"2026-08-17T11:59:32+08:00"}
--- FAIL: TestProjectUpdateRejectsDuplicateTaskIDs (0.00s)
=== RUN   TestProjectUpdateRejectsDuplicateTaskIDs
    duplicate_task_id_test.go:47: 重复任务ID应被拒绝，实际状态码为 200，响应为 {"id":"project-1","name":"数据迁移","tasks":[{"id":"approval","name":"业务审批","start_date":"2026-08-01","end_date":"2026-08-02","dependencies":null,"assignee":"","status":"未开始","progress":0},{"id":"approval","name":"安全审批","start_date":"2026-08-03","end_date":"2026-08-04","dependencies":null,"assignee":"","status":"未开始","progress":0}],"created_at":"2026-08-17T11:59:32+08:00","updated_at":"2026-08-17T11:59:32+08:00"}
--- FAIL: TestProjectUpdateRejectsDuplicateTaskIDs (0.00s)
=== RUN   TestProjectUpdateRejectsDuplicateTaskIDs
    duplicate_task_id_test.go:47: 重复任务ID应被拒绝，实际状态码为 200，响应为 {"id":"project-1","name":"数据迁移","tasks":[{"id":"approval","name":"业务审批","start_date":"2026-08-01","end_date":"2026-08-02","dependencies":null,"assignee":"","status":"未开始","progress":0},{"id":"approval","name":"安全审批","start_date":"2026-08-03","end_date":"2026-08-04","dependencies":null,"assignee":"","status":"未开始","progress":0}],"created_at":"2026-08-17T11:59:32+08:00","updated_at":"2026-08-17T11:59:32+08:00"}
--- FAIL: TestProjectUpdateRejectsDuplicateTaskIDs (0.00s)
=== RUN   TestProjectUpdateRejectsDuplicateTaskIDs
    duplicate_task_id_test.go:47: 重复任务ID应被拒绝，实际状态码为 200，响应为 {"id":"project-1","name":"数据迁移","tasks":[{"id":"approval","name":"业务审批","start_date":"2026-08-01","end_date":"2026-08-02","dependencies":null,"assignee":"","status":"未开始","progress":0},{"id":"approval","name":"安全审批","start_date":"2026-08-03","end_date":"2026-08-04","dependencies":null,"assignee":"","status":"未开始","progress":0}],"created_at":"2026-08-17T11:59:32+08:00","updated_at":"2026-08-17T11:59:32+08:00"}
--- FAIL: TestProjectUpdateRejectsDuplicateTaskIDs (0.00s)
=== RUN   TestProjectUpdateRejectsDuplicateTaskIDs
    duplicate_task_id_test.go:47: 重复任务ID应被拒绝，实际状态码为 200，响应为 {"id":"project-1","name":"数据迁移","tasks":[{"id":"approval","name":"业务审批","start_date":"2026-08-01","end_date":"2026-08-02","dependencies":null,"assignee":"","status":"未开始","progress":0},{"id":"approval","name":"安全审批","start_date":"2026-08-03","end_date":"2026-08-04","dependencies":null,"assignee":"","status":"未开始","progress":0}],"created_at":"2026-08-17T11:59:32+08:00","updated_at":"2026-08-17T11:59:32+08:00"}
--- FAIL: TestProjectUpdateRejectsDuplicateTaskIDs (0.00s)
FAIL
FAIL	gantt-generator/handlers	0.586s
FAIL
```
