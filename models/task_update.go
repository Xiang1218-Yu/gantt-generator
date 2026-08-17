package models

import (
	"fmt"
	"strings"
	"time"
)

// TaskPatch 表示任务编辑接口允许更新的字段。使用指针可以区分“没有提交”与“明确清空”。
type TaskPatch struct {
	Name         *string     `json:"name"`
	StartDate    *string     `json:"start_date"`
	EndDate      *string     `json:"end_date"`
	Dependencies *[]string   `json:"dependencies"`
	Assignee     *string     `json:"assignee"`
	Status       *TaskStatus `json:"status"`
	Progress     *int        `json:"progress"`
}

// ApplyTaskPatch 将部分更新应用到已有任务，并保留请求中没有出现的排期字段。
func ApplyTaskPatch(task *Task, patch TaskPatch) error {
	if task == nil {
		return fmt.Errorf("任务不能为空")
	}

	if patch.Name != nil {
		name := strings.TrimSpace(*patch.Name)
		if name == "" {
			return fmt.Errorf("任务名称不能为空")
		}
		task.Name = name
	}
	if patch.StartDate != nil {
		task.StartDate = *patch.StartDate
	}
	if patch.EndDate != nil {
		task.EndDate = *patch.EndDate
	}
	if patch.Dependencies != nil {
		task.Dependencies = append([]string(nil), (*patch.Dependencies)...)
	}
	if patch.Assignee != nil {
		task.Assignee = *patch.Assignee
	}
	if patch.Status != nil {
		if !isKnownTaskStatus(*patch.Status) {
			return fmt.Errorf("未知任务状态: %s", *patch.Status)
		}
		task.Status = *patch.Status
	}
	if patch.Progress != nil {
		if *patch.Progress < 0 || *patch.Progress > 100 {
			return fmt.Errorf("任务进度必须在 0 到 100 之间")
		}
		task.Progress = *patch.Progress
	}

	if _, err := time.Parse("2006-01-02", task.StartDate); err != nil {
		return fmt.Errorf("开始日期格式错误: %w", err)
	}
	if _, err := time.Parse("2006-01-02", task.EndDate); err != nil {
		return fmt.Errorf("结束日期格式错误: %w", err)
	}
	start, _ := time.Parse("2006-01-02", task.StartDate)
	end, _ := time.Parse("2006-01-02", task.EndDate)
	if end.Before(start) {
		return fmt.Errorf("结束日期不能早于开始日期")
	}
	return nil
}

func isKnownTaskStatus(status TaskStatus) bool {
	switch status {
	case StatusNotStarted, StatusInProgress, StatusCompleted, StatusDelayed:
		return true
	default:
		return false
	}
}
