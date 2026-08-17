package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"

	"gantt-generator/models"
	"gantt-generator/pdfgantt"
	"gantt-generator/report"
	"gantt-generator/scheduler"
	"gantt-generator/storage"
)

// Handler HTTP处理器
type Handler struct {
	store *storage.Store
}

// NewHandler 创建新的HTTP处理器
func NewHandler(store *storage.Store) *Handler {
	return &Handler{store: store}
}

// generateID 生成随机ID
func generateID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// writeJSON 写JSON响应
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// writeError 写错误响应
func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// parseJSON 解析JSON请求体
func parseJSON(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}

// ========== 项目相关接口 ==========

// GetProjects 获取所有项目
func (h *Handler) GetProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := h.store.GetAllProjects()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, projects)
}

// GetProject 获取单个项目
func (h *Handler) GetProject(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/projects/")
	if id == "" {
		writeError(w, http.StatusBadRequest, "缺少项目ID")
		return
	}
	project, err := h.store.GetProject(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, project)
}

// CreateProject 创建项目
func (h *Handler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var project models.Project
	if err := parseJSON(r, &project); err != nil {
		writeError(w, http.StatusBadRequest, "无效的请求数据: "+err.Error())
		return
	}
	if project.Name == "" {
		writeError(w, http.StatusBadRequest, "项目名称不能为空")
		return
	}
	if project.ID == "" {
		project.ID = generateID()
	}
	if err := h.store.CreateProject(&project); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, project)
}

// UpdateProject 更新项目
func (h *Handler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/projects/")
	if id == "" {
		writeError(w, http.StatusBadRequest, "缺少项目ID")
		return
	}
	var project models.Project
	if err := parseJSON(r, &project); err != nil {
		writeError(w, http.StatusBadRequest, "无效的请求数据: "+err.Error())
		return
	}
	project.ID = id
	if err := h.store.UpdateProject(&project); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, project)
}

// DeleteProject 删除项目
func (h *Handler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/projects/")
	if id == "" {
		writeError(w, http.StatusBadRequest, "缺少项目ID")
		return
	}
	if err := h.store.DeleteProject(id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "删除成功"})
}

// ========== 任务相关接口 ==========

// AddTask 添加任务
func (h *Handler) AddTask(w http.ResponseWriter, r *http.Request) {
	// 路径: /api/projects/{id}/tasks
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 3 {
		writeError(w, http.StatusBadRequest, "URL格式错误")
		return
	}
	projectID := pathParts[2]

	var task models.Task
	if err := parseJSON(r, &task); err != nil {
		writeError(w, http.StatusBadRequest, "无效的请求数据: "+err.Error())
		return
	}
	if task.Name == "" {
		writeError(w, http.StatusBadRequest, "任务名称不能为空")
		return
	}
	if task.ID == "" {
		task.ID = generateID()
	}
	if task.Status == "" {
		task.Status = models.StatusNotStarted
	}
	if task.Dependencies == nil {
		task.Dependencies = []string{}
	}

	if err := h.store.AddTask(projectID, &task); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, task)
}

// UpdateTask 更新任务
//
// 支持局部更新：进度页面只会提交 progress 和 status，请求体中未出现的字段
// 必须保留原任务的值。否则一次"把开发从 35 改成 60"的提交会把任务名称、起止
// 日期、前置关系一并清空，关键路径接口随后因日期为空而报"日期不合法"。
//
// 实现方式：用指针字段解码请求体，以区分"字段未提供"（nil，保留原值）与
// "字段显式置空/零值"（非 nil，覆盖），再把出现的字段合并到原任务上落盘。
func (h *Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	// 路径: /api/projects/{projectId}/tasks/{taskId}
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 4 {
		writeError(w, http.StatusBadRequest, "URL格式错误")
		return
	}
	projectID := pathParts[2]
	taskID := pathParts[4]

	var req taskUpdateRequest
	if err := parseJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "无效的请求数据: "+err.Error())
		return
	}

	// 取出原任务作为合并基底，局部更新只覆盖请求中出现的字段。
	project, err := h.store.GetProject(projectID)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	var existing *models.Task
	for i := range project.Tasks {
		if project.Tasks[i].ID == taskID {
			existing = &project.Tasks[i]
			break
		}
	}
	if existing == nil {
		writeError(w, http.StatusNotFound, "任务不存在: "+taskID)
		return
	}

	merged := *existing
	merged.ID = taskID
	req.applyTo(&merged)

	if merged.Name == "" {
		writeError(w, http.StatusBadRequest, "任务名称不能为空")
		return
	}
	if merged.Status == "" {
		merged.Status = models.StatusNotStarted
	}
	if merged.Dependencies == nil {
		merged.Dependencies = []string{}
	}

	if err := h.store.UpdateTask(projectID, &merged); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, merged)
}

// taskUpdateRequest 任务更新请求体。
// 全部使用指针字段：nil 表示请求未提供该字段（保留原值），非 nil 表示显式给出（覆盖）。
type taskUpdateRequest struct {
	Name         *string            `json:"name"`
	StartDate    *string            `json:"start_date"`
	EndDate      *string            `json:"end_date"`
	Assignee     *string            `json:"assignee"`
	Status       *models.TaskStatus `json:"status"`
	Progress     *int               `json:"progress"`
	Dependencies *[]string          `json:"dependencies"`
}

// applyTo 把请求中提供的字段覆盖到目标任务上，未提供的字段保持不变。
func (r *taskUpdateRequest) applyTo(t *models.Task) {
	if r == nil {
		return
	}
	if r.Name != nil {
		t.Name = *r.Name
	}
	if r.StartDate != nil {
		t.StartDate = *r.StartDate
	}
	if r.EndDate != nil {
		t.EndDate = *r.EndDate
	}
	if r.Assignee != nil {
		t.Assignee = *r.Assignee
	}
	if r.Status != nil {
		t.Status = *r.Status
	}
	if r.Progress != nil {
		t.Progress = *r.Progress
	}
	if r.Dependencies != nil {
		deps := *r.Dependencies
		if deps == nil {
			deps = []string{}
		}
		t.Dependencies = deps
	}
}

// DeleteTask 删除任务
func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 4 {
		writeError(w, http.StatusBadRequest, "URL格式错误")
		return
	}
	projectID := pathParts[2]
	taskID := pathParts[4]

	if err := h.store.DeleteTask(projectID, taskID); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "删除成功"})
}

// ========== 分析功能接口 ==========

// CalculateCriticalPath 计算关键路径
func (h *Handler) CalculateCriticalPath(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 3 {
		writeError(w, http.StatusBadRequest, "URL格式错误")
		return
	}
	projectID := pathParts[2]

	project, err := h.store.GetProject(projectID)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	result, err := scheduler.CalculateCriticalPath(project)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, result)
}

// GenerateGanttPDF 生成甘特图PDF
func (h *Handler) GenerateGanttPDF(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 3 {
		writeError(w, http.StatusBadRequest, "URL格式错误")
		return
	}
	projectID := pathParts[2]

	project, err := h.store.GetProject(projectID)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	result, err := scheduler.CalculateCriticalPath(project)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	pdfBytes, err := pdfgantt.GenerateGanttPDF(project, result)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=gantt_"+project.Name+".pdf")
	w.WriteHeader(http.StatusOK)
	w.Write(pdfBytes)
}

// GetProgressReport 获取进度报告
func (h *Handler) GetProgressReport(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 3 {
		writeError(w, http.StatusBadRequest, "URL格式错误")
		return
	}
	projectID := pathParts[2]

	project, err := h.store.GetProject(projectID)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	result, err := scheduler.CalculateCriticalPath(project)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	progressReport, err := report.GenerateProgressReport(project, result)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// 支持文本格式
	format := r.URL.Query().Get("format")
	if format == "text" {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(report.FormatReportAsText(progressReport)))
		return
	}

	writeJSON(w, http.StatusOK, progressReport)
}
