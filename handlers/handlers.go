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
		writeError(w, http.StatusBadRequest, err.Error())
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
		writeError(w, http.StatusBadRequest, err.Error())
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
func (h *Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	// 路径: /api/projects/{projectId}/tasks/{taskId}
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 4 {
		writeError(w, http.StatusBadRequest, "URL格式错误")
		return
	}
	projectID := pathParts[2]
	taskID := pathParts[4]

	var task models.Task
	if err := parseJSON(r, &task); err != nil {
		writeError(w, http.StatusBadRequest, "无效的请求数据: "+err.Error())
		return
	}
	task.ID = taskID

	if err := h.store.UpdateTask(projectID, &task); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, task)
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
