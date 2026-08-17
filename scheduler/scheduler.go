package scheduler

import (
	"fmt"
	"sort"
	"time"

	"gantt-generator/models"
)

const dateLayout = "2006-01-02"

// CalculateCriticalPath 计算关键路径
func CalculateCriticalPath(project *models.Project) (*models.CriticalPathResult, error) {
	if err := models.ValidateProjectTasks(project); err != nil {
		return nil, err
	}
	if len(project.Tasks) == 0 {
		return &models.CriticalPathResult{
			CriticalPath:  []string{},
			TotalDuration: 0,
			TaskInfo:      make(map[string]models.TaskSchedule),
		}, nil
	}

	// 1. 解析日期并计算任务持续时间
	taskMap := make(map[string]*models.Task)
	taskSchedule := make(map[string]models.TaskSchedule)
	taskIDs := []string{}

	for i := range project.Tasks {
		task := &project.Tasks[i]
		taskMap[task.ID] = task
		taskIDs = append(taskIDs, task.ID)

		startDate, err := time.Parse(dateLayout, task.StartDate)
		if err != nil {
			return nil, fmt.Errorf("任务 %s 开始日期格式错误: %w", task.ID, err)
		}
		endDate, err := time.Parse(dateLayout, task.EndDate)
		if err != nil {
			return nil, fmt.Errorf("任务 %s 结束日期格式错误: %w", task.ID, err)
		}
		if endDate.Before(startDate) {
			return nil, fmt.Errorf("任务 %s 结束日期早于开始日期", task.ID)
		}
		duration := int(endDate.Sub(startDate).Hours()/24) + 1 // 包含首尾两天

		taskSchedule[task.ID] = models.TaskSchedule{
			TaskID:    task.ID,
			Duration:  duration,
			StartDate: startDate,
			EndDate:   endDate,
		}
	}

	// 验证依赖是否有效
	for _, task := range project.Tasks {
		for _, depID := range task.Dependencies {
			if _, exists := taskMap[depID]; !exists {
				return nil, fmt.Errorf("任务 %s 依赖不存在的任务: %s", task.ID, depID)
			}
		}
	}

	// 2. 拓扑排序并检测环
	topologicalOrder, err := topologicalSort(taskIDs, func(id string) []string {
		return taskMap[id].Dependencies
	})
	if err != nil {
		return nil, fmt.Errorf("拓扑排序失败: %w", err)
	}

	// 3. 正向计算：最早开始时间(ES)和最早完成时间(EF)
	for _, id := range topologicalOrder {
		sched := taskSchedule[id]
		task := taskMap[id]

		// 找到所有前置依赖中最大的最早完成时间
		maxEF := 0
		for _, depID := range task.Dependencies {
			depSched := taskSchedule[depID]
			if depSched.EarliestFinish > maxEF {
				maxEF = depSched.EarliestFinish
			}
		}

		sched.EarliestStart = maxEF
		sched.EarliestFinish = maxEF + sched.Duration
		taskSchedule[id] = sched
	}

	// 4. 找到项目总工期（所有任务中最大的EF）
	totalDuration := 0
	for _, id := range taskIDs {
		if taskSchedule[id].EarliestFinish > totalDuration {
			totalDuration = taskSchedule[id].EarliestFinish
		}
	}

	// 5. 反向计算：最晚完成时间(LF)和最晚开始时间(LS)
	// 先反转拓扑顺序
	reverseOrder := make([]string, len(topologicalOrder))
	copy(reverseOrder, topologicalOrder)
	for i, j := 0, len(reverseOrder)-1; i < j; i, j = i+1, j-1 {
		reverseOrder[i], reverseOrder[j] = reverseOrder[j], reverseOrder[i]
	}

	// 先初始化所有任务的LF为总工期
	for _, id := range taskIDs {
		sched := taskSchedule[id]
		sched.LatestFinish = totalDuration
		taskSchedule[id] = sched
	}

	// 构建后继关系
	successors := make(map[string][]string)
	for _, id := range taskIDs {
		task := taskMap[id]
		for _, depID := range task.Dependencies {
			successors[depID] = append(successors[depID], id)
		}
	}

	for _, id := range reverseOrder {
		sched := taskSchedule[id]

		// 找到所有后继任务中最小的最晚开始时间
		minLS := totalDuration
		for _, succID := range successors[id] {
			succSched := taskSchedule[succID]
			if succSched.LatestStart < minLS {
				minLS = succSched.LatestStart
			}
		}

		// 如果有后继任务，使用后继的最小LS；否则保持总工期
		if len(successors[id]) > 0 {
			sched.LatestFinish = minLS
		}
		sched.LatestStart = sched.LatestFinish - sched.Duration
		taskSchedule[id] = sched
	}

	// 6. 计算松弛时间(Slack)和识别关键路径
	for _, id := range taskIDs {
		sched := taskSchedule[id]
		sched.Slack = sched.LatestStart - sched.EarliestStart
		sched.IsCritical = (sched.Slack == 0)
		taskSchedule[id] = sched
	}

	// 7. 构建关键路径
	criticalPath := buildCriticalPath(taskIDs, taskMap, taskSchedule)

	return &models.CriticalPathResult{
		CriticalPath:  criticalPath,
		TotalDuration: totalDuration,
		TaskInfo:      taskSchedule,
	}, nil
}

// topologicalSort 拓扑排序
func topologicalSort(ids []string, getDeps func(string) []string) ([]string, error) {
	visited := make(map[string]bool)
	tempMarked := make(map[string]bool)
	result := []string{}
	var dfs func(id string) error

	dfs = func(id string) error {
		if tempMarked[id] {
			return fmt.Errorf("检测到循环依赖，涉及任务: %s", id)
		}
		if visited[id] {
			return nil
		}
		tempMarked[id] = true
		for _, dep := range getDeps(id) {
			if err := dfs(dep); err != nil {
				return err
			}
		}
		tempMarked[id] = false
		visited[id] = true
		result = append(result, id)
		return nil
	}

	// 按ID排序以获得确定性结果
	sortedIds := make([]string, len(ids))
	copy(sortedIds, ids)
	sort.Strings(sortedIds)

	for _, id := range sortedIds {
		if !visited[id] {
			if err := dfs(id); err != nil {
				return nil, err
			}
		}
	}
	return result, nil
}

// buildCriticalPath 构建关键路径任务ID列表
func buildCriticalPath(taskIDs []string, taskMap map[string]*models.Task, taskSchedule map[string]models.TaskSchedule) []string {
	// 找出所有关键任务
	criticalTasks := []string{}
	for _, id := range taskIDs {
		if taskSchedule[id].IsCritical {
			criticalTasks = append(criticalTasks, id)
		}
	}

	// 按开始时间排序
	sort.Slice(criticalTasks, func(i, j int) bool {
		si := taskSchedule[criticalTasks[i]]
		sj := taskSchedule[criticalTasks[j]]
		if si.EarliestStart != sj.EarliestStart {
			return si.EarliestStart < sj.EarliestStart
		}
		return si.EarliestFinish < sj.EarliestFinish
	})

	return criticalTasks
}

// GetProjectStartEndDate 获取项目的最早开始和最晚结束日期
func GetProjectStartEndDate(project *models.Project) (time.Time, time.Time, error) {
	if len(project.Tasks) == 0 {
		now := time.Now()
		return now, now, nil
	}

	var minStart, maxEnd time.Time
	initialized := false

	for _, task := range project.Tasks {
		start, err := time.Parse(dateLayout, task.StartDate)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
		end, err := time.Parse(dateLayout, task.EndDate)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}

		if !initialized {
			minStart = start
			maxEnd = end
			initialized = true
		} else {
			if start.Before(minStart) {
				minStart = start
			}
			if end.After(maxEnd) {
				maxEnd = end
			}
		}
	}

	return minStart, maxEnd, nil
}

// DaysBetween 计算两个日期之间的天数（包含首尾）
func DaysBetween(start, end time.Time) int {
	return int(end.Sub(start).Hours()/24) + 1
}

// ParseDate 解析日期
func ParseDate(dateStr string) (time.Time, error) {
	return time.Parse(dateLayout, dateStr)
}

// FormatDate 格式化日期
func FormatDate(t time.Time) string {
	return t.Format(dateLayout)
}
