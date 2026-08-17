package scheduler

import (
	"testing"

	"gantt-generator/models"
)

func TestCriticalPathKeepsExplicitCalendarGap(t *testing.T) {
	project := &models.Project{
		ID:   "calendar-gap",
		Name: "等待外部审批",
		Tasks: []models.Task{
			{ID: "spec", Name: "需求确认", StartDate: "2026-08-01", EndDate: "2026-08-02"},
			{ID: "approval", Name: "上线审批", StartDate: "2026-08-10", EndDate: "2026-08-12", Dependencies: []string{"spec"}},
		},
	}

	result, err := CalculateCriticalPath(project)
	if err != nil {
		t.Fatal(err)
	}
	if result.TotalDuration != 12 {
		t.Fatalf("计划从 8 月 1 日跨到 8 月 12 日，工期应包含等待空档得到 12，实际为 %d", result.TotalDuration)
	}
	if got := result.TaskInfo["approval"].EarliestStart; got != 9 {
		t.Fatalf("审批任务应保留计划的第 9 天开始偏移，实际为 %d", got)
	}
}
