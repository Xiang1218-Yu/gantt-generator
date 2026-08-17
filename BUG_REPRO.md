# 缺陷复现说明

## 缺陷现象
进度页面的编辑框只提交进度和状态。把“开发”从 35 改成 60 后，任务名称、起止日期和前置关系都被清空，关键路径接口随后提示日期不合法。请修复这个任务局部更新把排期信息冲掉的问题。

## 触发方式
# 1. 进入测试项目根目录；预期：后续命令在目标项目内执行。
cd /workplace/gantt-generator
# 2. target-verify 检查只提交进度和状态的任务编辑；预期：未修复时排期字段保留断言失败。
go test -v ./handlers -run '^TestPartialTaskUpdatePreservesExistingSchedule$' -count=20
# 3. 修复后执行全量回归；预期：所有包通过且不引入新的失败。
go test ./...

## 触发后的实际错误输出

```text
=== RUN   TestPartialTaskUpdatePreservesExistingSchedule
    partial_task_update_test.go:57: 局部进度更新不应清空排期字段，实际任务为 {ID:build Name: StartDate: EndDate: Dependencies:[] Assignee: Status:进行中 Progress:60}
--- FAIL: TestPartialTaskUpdatePreservesExistingSchedule (0.00s)
=== RUN   TestPartialTaskUpdatePreservesExistingSchedule
    partial_task_update_test.go:57: 局部进度更新不应清空排期字段，实际任务为 {ID:build Name: StartDate: EndDate: Dependencies:[] Assignee: Status:进行中 Progress:60}
--- FAIL: TestPartialTaskUpdatePreservesExistingSchedule (0.00s)
=== RUN   TestPartialTaskUpdatePreservesExistingSchedule
    partial_task_update_test.go:57: 局部进度更新不应清空排期字段，实际任务为 {ID:build Name: StartDate: EndDate: Dependencies:[] Assignee: Status:进行中 Progress:60}
--- FAIL: TestPartialTaskUpdatePreservesExistingSchedule (0.00s)
=== RUN   TestPartialTaskUpdatePreservesExistingSchedule
    partial_task_update_test.go:57: 局部进度更新不应清空排期字段，实际任务为 {ID:build Name: StartDate: EndDate: Dependencies:[] Assignee: Status:进行中 Progress:60}
--- FAIL: TestPartialTaskUpdatePreservesExistingSchedule (0.00s)
=== RUN   TestPartialTaskUpdatePreservesExistingSchedule
    partial_task_update_test.go:57: 局部进度更新不应清空排期字段，实际任务为 {ID:build Name: StartDate: EndDate: Dependencies:[] Assignee: Status:进行中 Progress:60}
--- FAIL: TestPartialTaskUpdatePreservesExistingSchedule (0.00s)
=== RUN   TestPartialTaskUpdatePreservesExistingSchedule
    partial_task_update_test.go:57: 局部进度更新不应清空排期字段，实际任务为 {ID:build Name: StartDate: EndDate: Dependencies:[] Assignee: Status:进行中 Progress:60}
--- FAIL: TestPartialTaskUpdatePreservesExistingSchedule (0.00s)
=== RUN   TestPartialTaskUpdatePreservesExistingSchedule
    partial_task_update_test.go:57: 局部进度更新不应清空排期字段，实际任务为 {ID:build Name: StartDate: EndDate: Dependencies:[] Assignee: Status:进行中 Progress:60}
--- FAIL: TestPartialTaskUpdatePreservesExistingSchedule (0.00s)
=== RUN   TestPartialTaskUpdatePreservesExistingSchedule
    partial_task_update_test.go:57: 局部进度更新不应清空排期字段，实际任务为 {ID:build Name: StartDate: EndDate: Dependencies:[] Assignee: Status:进行中 Progress:60}
--- FAIL: TestPartialTaskUpdatePreservesExistingSchedule (0.00s)
=== RUN   TestPartialTaskUpdatePreservesExistingSchedule
    partial_task_update_test.go:57: 局部进度更新不应清空排期字段，实际任务为 {ID:build Name: StartDate: EndDate: Dependencies:[] Assignee: Status:进行中 Progress:60}
--- FAIL: TestPartialTaskUpdatePreservesExistingSchedule (0.00s)
=== RUN   TestPartialTaskUpdatePreservesExistingSchedule
    partial_task_update_test.go:57: 局部进度更新不应清空排期字段，实际任务为 {ID:build Name: StartDate: EndDate: Dependencies:[] Assignee: Status:进行中 Progress:60}
--- FAIL: TestPartialTaskUpdatePreservesExistingSchedule (0.00s)
=== RUN   TestPartialTaskUpdatePreservesExistingSchedule
    partial_task_update_test.go:57: 局部进度更新不应清空排期字段，实际任务为 {ID:build Name: StartDate: EndDate: Dependencies:[] Assignee: Status:进行中 Progress:60}
--- FAIL: TestPartialTaskUpdatePreservesExistingSchedule (0.00s)
=== RUN   TestPartialTaskUpdatePreservesExistingSchedule
    partial_task_update_test.go:57: 局部进度更新不应清空排期字段，实际任务为 {ID:build Name: StartDate: EndDate: Dependencies:[] Assignee: Status:进行中 Progress:60}
--- FAIL: TestPartialTaskUpdatePreservesExistingSchedule (0.00s)
=== RUN   TestPartialTaskUpdatePreservesExistingSchedule
    partial_task_update_test.go:57: 局部进度更新不应清空排期字段，实际任务为 {ID:build Name: StartDate: EndDate: Dependencies:[] Assignee: Status:进行中 Progress:60}
--- FAIL: TestPartialTaskUpdatePreservesExistingSchedule (0.00s)
=== RUN   TestPartialTaskUpdatePreservesExistingSchedule
    partial_task_update_test.go:57: 局部进度更新不应清空排期字段，实际任务为 {ID:build Name: StartDate: EndDate: Dependencies:[] Assignee: Status:进行中 Progress:60}
--- FAIL: TestPartialTaskUpdatePreservesExistingSchedule (0.00s)
=== RUN   TestPartialTaskUpdatePreservesExistingSchedule
    partial_task_update_test.go:57: 局部进度更新不应清空排期字段，实际任务为 {ID:build Name: StartDate: EndDate: Dependencies:[] Assignee: Status:进行中 Progress:60}
--- FAIL: TestPartialTaskUpdatePreservesExistingSchedule (0.00s)
=== RUN   TestPartialTaskUpdatePreservesExistingSchedule
    partial_task_update_test.go:57: 局部进度更新不应清空排期字段，实际任务为 {ID:build Name: StartDate: EndDate: Dependencies:[] Assignee: Status:进行中 Progress:60}
--- FAIL: TestPartialTaskUpdatePreservesExistingSchedule (0.00s)
=== RUN   TestPartialTaskUpdatePreservesExistingSchedule
    partial_task_update_test.go:57: 局部进度更新不应清空排期字段，实际任务为 {ID:build Name: StartDate: EndDate: Dependencies:[] Assignee: Status:进行中 Progress:60}
--- FAIL: TestPartialTaskUpdatePreservesExistingSchedule (0.00s)
=== RUN   TestPartialTaskUpdatePreservesExistingSchedule
    partial_task_update_test.go:57: 局部进度更新不应清空排期字段，实际任务为 {ID:build Name: StartDate: EndDate: Dependencies:[] Assignee: Status:进行中 Progress:60}
--- FAIL: TestPartialTaskUpdatePreservesExistingSchedule (0.00s)
=== RUN   TestPartialTaskUpdatePreservesExistingSchedule
    partial_task_update_test.go:57: 局部进度更新不应清空排期字段，实际任务为 {ID:build Name: StartDate: EndDate: Dependencies:[] Assignee: Status:进行中 Progress:60}
--- FAIL: TestPartialTaskUpdatePreservesExistingSchedule (0.00s)
=== RUN   TestPartialTaskUpdatePreservesExistingSchedule
    partial_task_update_test.go:57: 局部进度更新不应清空排期字段，实际任务为 {ID:build Name: StartDate: EndDate: Dependencies:[] Assignee: Status:进行中 Progress:60}
--- FAIL: TestPartialTaskUpdatePreservesExistingSchedule (0.00s)
FAIL
FAIL	gantt-generator/handlers	0.024s
FAIL
```
