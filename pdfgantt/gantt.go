package pdfgantt

import (
	"bytes"
	"fmt"
	"sort"
	"time"

	"github.com/jung-kurt/gofpdf"

	"gantt-generator/models"
	"gantt-generator/scheduler"
)

const (
	pageWidth      = 297.0 // A4横向宽度 mm
	pageHeight     = 210.0 // A4横向高度 mm
	marginLeft     = 15.0
	marginRight    = 15.0
	marginTop      = 20.0
	marginBottom   = 20.0
	taskNameWidth  = 55.0
	assigneeWidth  = 30.0
	rowHeight      = 8.0
	headerHeight   = 12.0
	ganttStartX    = marginLeft + taskNameWidth + assigneeWidth
	dayWidth       = 3.5 // 每天的宽度 mm
)

// GenerateGanttPDF 生成甘特图PDF
func GenerateGanttPDF(project *models.Project, result *models.CriticalPathResult) ([]byte, error) {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.SetAutoPageBreak(true, marginBottom)
	pdf.AddPage()

	// 设置字体 - 使用内置字体，避免中文依赖
	pdf.SetFont("Arial", "B", 14)

	// 标题
	title := fmt.Sprintf("Project Gantt - %s", project.Name)
	pdf.CellFormat(pageWidth-marginLeft-marginRight, 15, title, "", 1, "C", false, 0, "")
	pdf.Ln(2)

	// 计算项目时间范围
	projStart, projEnd, err := scheduler.GetProjectStartEndDate(project)
	if err != nil {
		return nil, fmt.Errorf("计算项目时间范围失败: %w", err)
	}
	totalDays := scheduler.DaysBetween(projStart, projEnd)
	ganttWidth := float64(totalDays) * dayWidth

	// 确保宽度足够
	availWidth := pageWidth - marginLeft - marginRight - taskNameWidth - assigneeWidth
	if ganttWidth < availWidth {
		ganttWidth = availWidth
	}

	// 绘制时间刻度表头
	drawTimelineHeader(pdf, projStart, totalDays, ganttWidth)

	// 绘制任务列表
	// 按开始日期排序任务
	sortedTasks := sortTasksByDate(project.Tasks)

	y := pdf.GetY()
	for i, task := range sortedTasks {
		// 交替行背景
		if i%2 == 0 {
			pdf.SetFillColor(245, 245, 245)
			pdf.Rect(marginLeft, y, taskNameWidth+assigneeWidth+ganttWidth, rowHeight, "F")
		}

		// 任务名称
		pdf.SetFont("Arial", "", 9)
		pdf.SetXY(marginLeft, y)
		taskName := truncateString(task.Name, 25)
		pdf.CellFormat(taskNameWidth, rowHeight, taskName, "LR", 0, "L", false, 0, "")

		// 负责人
		assignee := truncateString(task.Assignee, 12)
		pdf.CellFormat(assigneeWidth, rowHeight, assignee, "LR", 0, "L", false, 0, "")

		// 绘制甘特条
		sched := result.TaskInfo[task.ID]
		drawGanttBar(pdf, y, task, sched, projStart, totalDays, ganttWidth)

		y += rowHeight

		// 换页
		if y > pageHeight-marginBottom-rowHeight {
			pdf.AddPage()
			drawTimelineHeader(pdf, projStart, totalDays, ganttWidth)
			y = pdf.GetY()
		}
	}

	// 图例
	pdf.Ln(5)
	drawLegend(pdf)

	// 关键路径信息
	pdf.Ln(3)
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(0, 8, fmt.Sprintf("Total Duration: %d days", result.TotalDuration), "", 1, "L", false, 0, "")
	if len(result.CriticalPath) > 0 {
		pdf.SetFont("Arial", "", 9)
		criticalNames := getCriticalTaskNames(project, result.CriticalPath)
		pdf.MultiCell(0, 6, fmt.Sprintf("Critical Path: %s", criticalNames), "", "L", false)
	}

	// 输出到buffer
	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return nil, fmt.Errorf("生成PDF失败: %w", err)
	}
	return buf.Bytes(), nil
}

// drawTimelineHeader 绘制时间刻度表头
func drawTimelineHeader(pdf *gofpdf.Fpdf, projStart time.Time, totalDays int, ganttWidth float64) {
	y := pdf.GetY()

	// 表头背景
	pdf.SetFillColor(80, 120, 180)
	pdf.Rect(marginLeft, y, taskNameWidth+assigneeWidth+ganttWidth, headerHeight, "F")
	pdf.SetTextColor(255, 255, 255)

	// 任务名称表头
	pdf.SetFont("Arial", "B", 10)
	pdf.SetXY(marginLeft, y)
	pdf.CellFormat(taskNameWidth, headerHeight, "Task", "LTRB", 0, "C", true, 0, "")

	// 负责人表头
	pdf.CellFormat(assigneeWidth, headerHeight, "Owner", "LTRB", 0, "C", true, 0, "")

	// 甘特图区域 - 绘制月份刻度
	x := ganttStartX
	currentMonth := projStart.Month()
	currentYear := projStart.Year()
	monthStartX := x
	monthDays := 0

	for d := 0; d < totalDays; d++ {
		date := projStart.AddDate(0, 0, d)
		if date.Month() != currentMonth || d == totalDays-1 {
			if d == totalDays-1 {
				monthDays++
			}
			// 绘制月份标签
			monthLabel := fmt.Sprintf("%d/%02d", currentYear, currentMonth)
			monthWidth := float64(monthDays) * dayWidth
			pdf.SetXY(monthStartX, y)
			pdf.CellFormat(monthWidth, headerHeight/2, monthLabel, "LTRB", 0, "C", true, 0, "")

			// 绘制日期刻度
			drawDayTicks(pdf, monthStartX, y+headerHeight/2, monthDays, totalDays, d-monthDays+1)

			currentMonth = date.Month()
			currentYear = date.Year()
			monthStartX = x + float64(d)*dayWidth
			monthDays = 0
		}
		monthDays++
	}
	// 处理最后一个月
	if monthDays > 0 {
		monthWidth := float64(monthDays) * dayWidth
		pdf.SetXY(monthStartX, y)
		pdf.CellFormat(monthWidth, headerHeight/2, fmt.Sprintf("%d/%02d", currentYear, currentMonth), "LTRB", 0, "C", true, 0, "")
		drawDayTicks(pdf, monthStartX, y+headerHeight/2, monthDays, totalDays, totalDays-monthDays)
	}

	pdf.SetTextColor(0, 0, 0)
	pdf.SetY(y + headerHeight)
}

// drawDayTicks 绘制日期刻度
func drawDayTicks(pdf *gofpdf.Fpdf, x, y float64, monthDays, totalDays, startDay int) {
	pdf.SetFont("Arial", "", 7)
	for d := 0; d < monthDays; d++ {
		// 每5天或第1天显示日期
		absDay := startDay + d
		if d%5 == 0 || absDay == 1 || d == monthDays-1 {
			dayX := x + float64(d)*dayWidth
			pdf.SetXY(dayX, y)
			pdf.CellFormat(dayWidth*5, headerHeight/2, fmt.Sprintf("%d", absDay), "LTRB", 0, "L", true, 0, "")
		}
	}
}

// drawGanttBar 绘制甘特条
func drawGanttBar(pdf *gofpdf.Fpdf, y float64, task models.Task, sched models.TaskSchedule, projStart time.Time, totalDays int, ganttWidth float64) {
	// 计算任务条位置
	taskStartOffset := scheduler.DaysBetween(projStart, sched.StartDate) - 1
	barX := ganttStartX + float64(taskStartOffset)*dayWidth
	barWidth := float64(sched.Duration) * dayWidth

	// 确保不超出范围
	if barX < ganttStartX {
		barWidth -= ganttStartX - barX
		barX = ganttStartX
	}
	if barX+barWidth > ganttStartX+ganttWidth {
		barWidth = ganttStartX + ganttWidth - barX
	}

	// 设置颜色
	var r, g, b int
	switch task.Status {
	case models.StatusCompleted:
		r, g, b = 80, 180, 80 // 绿色
	case models.StatusInProgress:
		r, g, b = 80, 140, 220 // 蓝色
	case models.StatusDelayed:
		r, g, b = 220, 80, 80 // 红色
	default:
		r, g, b = 160, 160, 160 // 灰色
	}

	// 关键任务用橙色边框
	barY := y + 1.5
	barH := rowHeight - 3.0

	pdf.SetDrawColor(0, 0, 0)
	if sched.IsCritical {
		pdf.SetDrawColor(255, 140, 0) // 橙色
		pdf.SetLineWidth(0.5)
	}

	// 绘制任务条背景
	pdf.SetFillColor(r, g, b)
	pdf.Rect(barX, barY, barWidth, barH, "DF")

	// 绘制进度条
	if task.Progress > 0 && task.Progress < 100 {
		progressWidth := barWidth * float64(task.Progress) / 100.0
		dr, dg, db := r-30, g-30, b-30
		if dr < 0 {
			dr = 0
		}
		if dg < 0 {
			dg = 0
		}
		if db < 0 {
			db = 0
		}
		pdf.SetFillColor(dr, dg, db)
		pdf.Rect(barX, barY, progressWidth, barH, "F")
	}

	pdf.SetLineWidth(0.2)
	pdf.SetDrawColor(0, 0, 0)
}

// drawLegend 绘制图例
func drawLegend(pdf *gofpdf.Fpdf) {
	pdf.SetFont("Arial", "B", 9)
	pdf.CellFormat(0, 7, "Legend:", "", 1, "L", false, 0, "")

	legends := []struct {
		color [3]int
		label string
	}{
		{[3]int{160, 160, 160}, "Not Started"},
		{[3]int{80, 140, 220}, "In Progress"},
		{[3]int{80, 180, 80}, "Completed"},
		{[3]int{220, 80, 80}, "Delayed"},
	}

	x := pdf.GetX()
	y := pdf.GetY()
	for i, leg := range legends {
		px := x + float64(i)*50.0
		pdf.SetFillColor(leg.color[0], leg.color[1], leg.color[2])
		pdf.Rect(px, y+1, 8, 5, "DF")
		pdf.SetXY(px+10, y)
		pdf.SetFont("Arial", "", 9)
		pdf.CellFormat(40, 7, leg.label, "", 0, "L", false, 0, "")
	}
	pdf.Ln(7)

	// 关键路径标记
	pdf.SetDrawColor(255, 140, 0)
	pdf.SetLineWidth(0.5)
	pdf.Rect(x, y+8, 8, 5, "D")
	pdf.SetLineWidth(0.2)
	pdf.SetDrawColor(0, 0, 0)
	pdf.SetXY(x+10, y+7)
	pdf.CellFormat(0, 7, "Orange border = Critical Path", "", 1, "L", false, 0, "")
}

// sortTasksByDate 按开始日期排序任务
func sortTasksByDate(tasks []models.Task) []models.Task {
	sorted := make([]models.Task, len(tasks))
	copy(sorted, tasks)
	sort.Slice(sorted, func(i, j int) bool {
		ti, _ := time.Parse("2006-01-02", sorted[i].StartDate)
		tj, _ := time.Parse("2006-01-02", sorted[j].StartDate)
		if !ti.Equal(tj) {
			return ti.Before(tj)
		}
		return sorted[i].ID < sorted[j].ID
	})
	return sorted
}

// truncateString 截断字符串
func truncateString(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) > maxLen {
		return string(runes[:maxLen-2]) + ".."
	}
	return s
}

// getCriticalTaskNames 获取关键路径任务名称
func getCriticalTaskNames(project *models.Project, criticalPath []string) string {
	nameMap := make(map[string]string)
	for _, t := range project.Tasks {
		nameMap[t.ID] = t.Name
	}
	result := ""
	for i, id := range criticalPath {
		if i > 0 {
			result += " -> "
		}
		result += nameMap[id]
	}
	return result
}
