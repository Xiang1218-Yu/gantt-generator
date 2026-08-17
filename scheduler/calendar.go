package scheduler

import (
	"fmt"
	"time"

	"gantt-generator/models"
)

// calendarOffset 返回日期相对项目最早计划开始日的零基偏移。
func calendarOffset(projectStart, date time.Time) (int, error) {
	if projectStart.IsZero() || date.IsZero() {
		return 0, fmt.Errorf("项目计划日期不能为空")
	}
	offset := int(date.Sub(projectStart).Hours() / 24)
	if offset < 0 {
		return 0, fmt.Errorf("任务日期早于项目起点")
	}
	return offset, nil
}

// applyPlannedOffsets 保存每个任务的计划日历位置。关键路径计算既要遵守依赖，
// 也不能压缩业务排期中明确存在的等待空档。
func applyPlannedOffsets(schedule map[string]models.TaskSchedule, projectStart time.Time) error {
	for id, value := range schedule {
		offset, err := calendarOffset(projectStart, value.StartDate)
		if err != nil {
			return fmt.Errorf("任务 %s 的计划开始时间无效: %w", id, err)
		}
		value.PlannedStart = offset
		value.PlannedFinish = offset + value.Duration
		schedule[id] = value
	}
	return nil
}
