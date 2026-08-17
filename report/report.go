package report

import (
	"fmt"
	"time"

	"gantt-generator/models"
)

// GenerateProgressReport 生成进度报告
func GenerateProgressReport(project *models.Project, result *models.CriticalPathResult) (*models.ProgressReport, error) {
	report := &models.ProgressReport{
		ProjectName: project.Name,
		ReportDate:  time.Now().Format("2006-01-02 15:04:05"),
		TotalTasks:  len(project.Tasks),
		TaskDetails: []models.TaskProgressDetail{},
		CriticalPath: result.CriticalPath,
		TotalDuration: result.TotalDuration,
	}

	// 统计各状态任务数
	totalProgress := 0
	now := time.Now()
	_ = totalProgress

	for _, task := range project.Tasks {
		// 统计状态
		switch task.Status {
		case models.StatusCompleted:
			report.CompletedTasks++
		case models.StatusInProgress:
			report.InProgressTasks++
		case models.StatusNotStarted:
			report.NotStartedTasks++
		case models.StatusDelayed:
			report.DelayedTasks++
		}

		// 累加进度
		totalProgress += task.Progress

		// 检查是否延期（已过结束日期但未完成）
		isDelayed := false
		if task.Status == models.StatusDelayed {
			isDelayed = true
		} else {
			endDate, err := time.Parse("2006-01-02", task.EndDate)
			if err == nil && now.After(endDate) && task.Status != models.StatusCompleted {
				isDelayed = true
			}
		}

		// 是否为关键任务
		sched := result.TaskInfo[task.ID]

		detail := models.TaskProgressDetail{
			ID:         task.ID,
			Name:       task.Name,
			Assignee:   task.Assignee,
			Status:     task.Status,
			Progress:   task.Progress,
			StartDate:  task.StartDate,
			EndDate:    task.EndDate,
			IsCritical: sched.IsCritical,
			IsDelayed:  isDelayed,
		}
		report.TaskDetails = append(report.TaskDetails, detail)
	}

	// 计算整体进度（加权平均：按任务工期加权）
	if len(project.Tasks) > 0 {
		weightedProgress := 0
		totalWeight := 0
		for _, task := range project.Tasks {
			sched := result.TaskInfo[task.ID]
			weight := sched.Duration
			if weight <= 0 {
				weight = 1
			}
			weightedProgress += task.Progress * weight
			totalWeight += weight
		}
		if totalWeight > 0 {
			report.OverallProgress = weightedProgress / totalWeight
		} else {
			report.OverallProgress = totalProgress / len(project.Tasks)
		}
	}

	return report, nil
}

// FormatReportAsText 将报告格式化为纯文本
func FormatReportAsText(report *models.ProgressReport) string {
	text := fmt.Sprintf("=== 项目进度报告 ===\n")
	text += fmt.Sprintf("项目名称: %s\n", report.ProjectName)
	text += fmt.Sprintf("报告日期: %s\n\n", report.ReportDate)

	text += fmt.Sprintf("--- 项目概览 ---\n")
	text += fmt.Sprintf("总工期: %d 天\n", report.TotalDuration)
	text += fmt.Sprintf("总任务数: %d\n", report.TotalTasks)
	text += fmt.Sprintf("整体进度: %d%%\n", report.OverallProgress)
	text += fmt.Sprintf("  - 已完成: %d 个\n", report.CompletedTasks)
	text += fmt.Sprintf("  - 进行中: %d 个\n", report.InProgressTasks)
	text += fmt.Sprintf("  - 未开始: %d 个\n", report.NotStartedTasks)
	text += fmt.Sprintf("  - 延期: %d 个\n\n", report.DelayedTasks)

	text += fmt.Sprintf("--- 关键路径 ---\n")
	if len(report.CriticalPath) > 0 {
		for i, id := range report.CriticalPath {
			taskName := id
			for _, td := range report.TaskDetails {
				if td.ID == id {
					taskName = td.Name
					break
				}
			}
			if i > 0 {
				text += " → "
			}
			text += taskName
		}
		text += "\n\n"
	} else {
		text += "无\n\n"
	}

	text += fmt.Sprintf("--- 任务详情 ---\n")
	for i, td := range report.TaskDetails {
		text += fmt.Sprintf("\n任务 %d: %s\n", i+1, td.Name)
		text += fmt.Sprintf("  ID: %s\n", td.ID)
		text += fmt.Sprintf("  负责人: %s\n", td.Assignee)
		text += fmt.Sprintf("  时间: %s ~ %s\n", td.StartDate, td.EndDate)
		text += fmt.Sprintf("  状态: %s\n", td.Status)
		text += fmt.Sprintf("  进度: %d%%\n", td.Progress)
		if td.IsCritical {
			text += fmt.Sprintf("  [关键任务]\n")
		}
		if td.IsDelayed {
			text += fmt.Sprintf("  [⚠️ 已延期]\n")
		}
	}

	return text
}
