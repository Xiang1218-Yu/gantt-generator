package models

import "time"

// TaskStatus 任务状态枚举
type TaskStatus string

const (
	StatusNotStarted TaskStatus = "未开始"
	StatusInProgress TaskStatus = "进行中"
	StatusCompleted  TaskStatus = "已完成"
	StatusDelayed    TaskStatus = "延期"
)

// Task 任务模型
type Task struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	StartDate   string     `json:"start_date"`
	EndDate     string     `json:"end_date"`
	Dependencies []string   `json:"dependencies"` // 前置依赖任务ID列表
	Assignee    string     `json:"assignee"`
	Status      TaskStatus `json:"status"`
	Progress    int        `json:"progress"` // 0-100
}

// Project 项目模型
type Project struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Tasks     []Task `json:"tasks"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// CriticalPathResult 关键路径计算结果
type CriticalPathResult struct {
	CriticalPath  []string        `json:"critical_path"`  // 关键路径上的任务ID列表
	TotalDuration int             `json:"total_duration"` // 总工期（天）
	TaskInfo      map[string]TaskSchedule `json:"task_info"` // 每个任务的调度信息
}

// TaskSchedule 任务调度信息
type TaskSchedule struct {
	TaskID    string `json:"task_id"`
	EarliestStart  int `json:"earliest_start"`  // 最早开始时间（天，从项目起点计）
	EarliestFinish int `json:"earliest_finish"` // 最早完成时间
	LatestStart    int `json:"latest_start"`    // 最晚开始时间
	LatestFinish   int `json:"latest_finish"`   // 最晚完成时间
	Slack          int `json:"slack"`           // 松弛时间
	IsCritical     bool `json:"is_critical"`    // 是否在关键路径上
	Duration       int `json:"duration"`        // 任务持续天数
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
}

// ProgressReport 进度报告
type ProgressReport struct {
	ProjectName     string              `json:"project_name"`
	ReportDate      string              `json:"report_date"`
	TotalTasks      int                 `json:"total_tasks"`
	CompletedTasks  int                 `json:"completed_tasks"`
	InProgressTasks int                 `json:"in_progress_tasks"`
	NotStartedTasks int                 `json:"not_started_tasks"`
	DelayedTasks    int                 `json:"delayed_tasks"`
	OverallProgress int                 `json:"overall_progress"` // 整体进度百分比
	TotalDuration   int                 `json:"total_duration"`
	TaskDetails     []TaskProgressDetail `json:"task_details"`
	CriticalPath    []string            `json:"critical_path"`
}

// TaskProgressDetail 任务进度详情
type TaskProgressDetail struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Assignee    string     `json:"assignee"`
	Status      TaskStatus `json:"status"`
	Progress    int        `json:"progress"`
	StartDate   string     `json:"start_date"`
	EndDate     string     `json:"end_date"`
	IsCritical  bool       `json:"is_critical"`
	IsDelayed   bool       `json:"is_delayed"`
}
