# 缺陷复现说明

## 缺陷现象
项目里两个有依赖的任务之间会等审批。进度报告把等待空档压没了，例如排期从 8 月 1 日到 8 月 12 日却只算成 5 天，但甘特图仍按原日期展示。文件先不要改，帮我定位为什么两个视图的工期不一致。

## 触发方式
# 1. 进入测试项目根目录；预期：后续命令在目标项目内执行。
cd /workplace/gantt-generator
# 2. target-verify 复现带日历空档的依赖排期；预期：基线会把 8 月 1 日到 8 月 12 日错误压成 5 天。
go test -v ./scheduler -run '^TestCriticalPathKeepsExplicitCalendarGap$' -count=20

## 触发后的实际错误输出

```text
=== RUN   TestCriticalPathKeepsExplicitCalendarGap
    calendar_gap_test.go:24: 计划从 8 月 1 日跨到 8 月 12 日，工期应包含等待空档得到 12，实际为 5
--- FAIL: TestCriticalPathKeepsExplicitCalendarGap (0.00s)
=== RUN   TestCriticalPathKeepsExplicitCalendarGap
    calendar_gap_test.go:24: 计划从 8 月 1 日跨到 8 月 12 日，工期应包含等待空档得到 12，实际为 5
--- FAIL: TestCriticalPathKeepsExplicitCalendarGap (0.00s)
=== RUN   TestCriticalPathKeepsExplicitCalendarGap
    calendar_gap_test.go:24: 计划从 8 月 1 日跨到 8 月 12 日，工期应包含等待空档得到 12，实际为 5
--- FAIL: TestCriticalPathKeepsExplicitCalendarGap (0.00s)
=== RUN   TestCriticalPathKeepsExplicitCalendarGap
    calendar_gap_test.go:24: 计划从 8 月 1 日跨到 8 月 12 日，工期应包含等待空档得到 12，实际为 5
--- FAIL: TestCriticalPathKeepsExplicitCalendarGap (0.00s)
=== RUN   TestCriticalPathKeepsExplicitCalendarGap
    calendar_gap_test.go:24: 计划从 8 月 1 日跨到 8 月 12 日，工期应包含等待空档得到 12，实际为 5
--- FAIL: TestCriticalPathKeepsExplicitCalendarGap (0.00s)
=== RUN   TestCriticalPathKeepsExplicitCalendarGap
    calendar_gap_test.go:24: 计划从 8 月 1 日跨到 8 月 12 日，工期应包含等待空档得到 12，实际为 5
--- FAIL: TestCriticalPathKeepsExplicitCalendarGap (0.00s)
=== RUN   TestCriticalPathKeepsExplicitCalendarGap
    calendar_gap_test.go:24: 计划从 8 月 1 日跨到 8 月 12 日，工期应包含等待空档得到 12，实际为 5
--- FAIL: TestCriticalPathKeepsExplicitCalendarGap (0.00s)
=== RUN   TestCriticalPathKeepsExplicitCalendarGap
    calendar_gap_test.go:24: 计划从 8 月 1 日跨到 8 月 12 日，工期应包含等待空档得到 12，实际为 5
--- FAIL: TestCriticalPathKeepsExplicitCalendarGap (0.00s)
=== RUN   TestCriticalPathKeepsExplicitCalendarGap
    calendar_gap_test.go:24: 计划从 8 月 1 日跨到 8 月 12 日，工期应包含等待空档得到 12，实际为 5
--- FAIL: TestCriticalPathKeepsExplicitCalendarGap (0.00s)
=== RUN   TestCriticalPathKeepsExplicitCalendarGap
    calendar_gap_test.go:24: 计划从 8 月 1 日跨到 8 月 12 日，工期应包含等待空档得到 12，实际为 5
--- FAIL: TestCriticalPathKeepsExplicitCalendarGap (0.00s)
=== RUN   TestCriticalPathKeepsExplicitCalendarGap
    calendar_gap_test.go:24: 计划从 8 月 1 日跨到 8 月 12 日，工期应包含等待空档得到 12，实际为 5
--- FAIL: TestCriticalPathKeepsExplicitCalendarGap (0.00s)
=== RUN   TestCriticalPathKeepsExplicitCalendarGap
    calendar_gap_test.go:24: 计划从 8 月 1 日跨到 8 月 12 日，工期应包含等待空档得到 12，实际为 5
--- FAIL: TestCriticalPathKeepsExplicitCalendarGap (0.00s)
=== RUN   TestCriticalPathKeepsExplicitCalendarGap
    calendar_gap_test.go:24: 计划从 8 月 1 日跨到 8 月 12 日，工期应包含等待空档得到 12，实际为 5
--- FAIL: TestCriticalPathKeepsExplicitCalendarGap (0.00s)
=== RUN   TestCriticalPathKeepsExplicitCalendarGap
    calendar_gap_test.go:24: 计划从 8 月 1 日跨到 8 月 12 日，工期应包含等待空档得到 12，实际为 5
--- FAIL: TestCriticalPathKeepsExplicitCalendarGap (0.00s)
=== RUN   TestCriticalPathKeepsExplicitCalendarGap
    calendar_gap_test.go:24: 计划从 8 月 1 日跨到 8 月 12 日，工期应包含等待空档得到 12，实际为 5
--- FAIL: TestCriticalPathKeepsExplicitCalendarGap (0.00s)
=== RUN   TestCriticalPathKeepsExplicitCalendarGap
    calendar_gap_test.go:24: 计划从 8 月 1 日跨到 8 月 12 日，工期应包含等待空档得到 12，实际为 5
--- FAIL: TestCriticalPathKeepsExplicitCalendarGap (0.00s)
=== RUN   TestCriticalPathKeepsExplicitCalendarGap
    calendar_gap_test.go:24: 计划从 8 月 1 日跨到 8 月 12 日，工期应包含等待空档得到 12，实际为 5
--- FAIL: TestCriticalPathKeepsExplicitCalendarGap (0.00s)
=== RUN   TestCriticalPathKeepsExplicitCalendarGap
    calendar_gap_test.go:24: 计划从 8 月 1 日跨到 8 月 12 日，工期应包含等待空档得到 12，实际为 5
--- FAIL: TestCriticalPathKeepsExplicitCalendarGap (0.00s)
=== RUN   TestCriticalPathKeepsExplicitCalendarGap
    calendar_gap_test.go:24: 计划从 8 月 1 日跨到 8 月 12 日，工期应包含等待空档得到 12，实际为 5
--- FAIL: TestCriticalPathKeepsExplicitCalendarGap (0.00s)
=== RUN   TestCriticalPathKeepsExplicitCalendarGap
    calendar_gap_test.go:24: 计划从 8 月 1 日跨到 8 月 12 日，工期应包含等待空档得到 12，实际为 5
--- FAIL: TestCriticalPathKeepsExplicitCalendarGap (0.00s)
FAIL
FAIL	gantt-generator/scheduler	0.764s
FAIL
```
