package models

import (
	"fmt"
	"time"
)

const taskDateLayout = "2006-01-02"

// ValidateProjectTasks 校验项目任务图在存储和调度之前是否保持可识别的一致状态。
func ValidateProjectTasks(project *Project) error {
	if project == nil {
		return fmt.Errorf("项目不能为空")
	}

	tasksByID := make(map[string]Task, len(project.Tasks))
	for _, task := range project.Tasks {
		if task.ID == "" {
			return fmt.Errorf("任务ID不能为空")
		}
		if _, exists := tasksByID[task.ID]; exists {
			return fmt.Errorf("项目中存在重复的任务ID: %s", task.ID)
		}
		if task.Name == "" {
			return fmt.Errorf("任务 %s 名称不能为空", task.ID)
		}
		if err := validateTaskDates(task); err != nil {
			return err
		}
		if task.Progress < 0 || task.Progress > 100 {
			return fmt.Errorf("任务 %s 进度必须在 0 到 100 之间", task.ID)
		}
		tasksByID[task.ID] = task
	}

	graph := make(map[string][]string, len(tasksByID))
	for _, task := range project.Tasks {
		seenDependency := make(map[string]struct{}, len(task.Dependencies))
		for _, dependencyID := range task.Dependencies {
			if dependencyID == task.ID {
				return fmt.Errorf("任务 %s 不能依赖自己", task.ID)
			}
			if _, exists := tasksByID[dependencyID]; !exists {
				return fmt.Errorf("任务 %s 依赖不存在的任务: %s", task.ID, dependencyID)
			}
			if _, repeated := seenDependency[dependencyID]; repeated {
				return fmt.Errorf("任务 %s 重复声明依赖: %s", task.ID, dependencyID)
			}
			seenDependency[dependencyID] = struct{}{}
			graph[task.ID] = append(graph[task.ID], dependencyID)
		}
	}

	visiting := make(map[string]bool, len(graph))
	visited := make(map[string]bool, len(graph))
	var visit func(string) error
	visit = func(id string) error {
		if visiting[id] {
			return fmt.Errorf("任务依赖存在循环，涉及任务: %s", id)
		}
		if visited[id] {
			return nil
		}
		visiting[id] = true
		for _, dependencyID := range graph[id] {
			if err := visit(dependencyID); err != nil {
				return err
			}
		}
		visiting[id] = false
		visited[id] = true
		return nil
	}
	for taskID := range tasksByID {
		if err := visit(taskID); err != nil {
			return err
		}
	}
	return nil
}

func validateTaskDates(task Task) error {
	start, err := time.Parse(taskDateLayout, task.StartDate)
	if err != nil {
		return fmt.Errorf("任务 %s 开始日期格式错误: %w", task.ID, err)
	}
	end, err := time.Parse(taskDateLayout, task.EndDate)
	if err != nil {
		return fmt.Errorf("任务 %s 结束日期格式错误: %w", task.ID, err)
	}
	if end.Before(start) {
		return fmt.Errorf("任务 %s 结束日期早于开始日期", task.ID)
	}
	return nil
}
