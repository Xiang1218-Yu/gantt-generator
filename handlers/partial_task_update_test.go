package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"gantt-generator/models"
	"gantt-generator/scheduler"
	"gantt-generator/storage"
)

func TestPartialTaskUpdatePreservesExistingSchedule(t *testing.T) {
	workDir := t.TempDir()
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(originalDir) })
	if err := os.Chdir(workDir); err != nil {
		t.Fatal(err)
	}

	store, err := storage.NewStore()
	if err != nil {
		t.Fatal(err)
	}
	project := &models.Project{
		ID:   "project-1",
		Name: "发布计划",
		Tasks: []models.Task{
			{ID: "design", Name: "设计", StartDate: "2026-08-01", EndDate: "2026-08-02", Status: models.StatusCompleted, Progress: 100},
			{ID: "build", Name: "开发", StartDate: "2026-08-03", EndDate: "2026-08-05", Dependencies: []string{"design"}, Assignee: "李四", Status: models.StatusInProgress, Progress: 35},
		},
	}
	if err := store.CreateProject(project); err != nil {
		t.Fatal(err)
	}

	handler := NewHandler(store)
	request := httptest.NewRequest(http.MethodPut, "/api/projects/project-1/tasks/build", bytes.NewBufferString(`{"progress":60,"status":"进行中"}`))
	response := httptest.NewRecorder()
	handler.UpdateTask(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("更新任务返回 %d: %s", response.Code, response.Body.String())
	}

	saved, err := store.GetProject("project-1")
	if err != nil {
		t.Fatal(err)
	}
	task := saved.Tasks[1]
	if task.Name != "开发" || task.StartDate != "2026-08-03" || task.EndDate != "2026-08-05" {
		t.Fatalf("局部进度更新不应清空排期字段，实际任务为 %+v", task)
	}
	if len(task.Dependencies) != 1 || task.Dependencies[0] != "design" {
		t.Fatalf("局部进度更新不应丢失前置关系，实际依赖为 %#v", task.Dependencies)
	}
	if task.Progress != 60 {
		t.Fatalf("进度没有更新，得到 %d", task.Progress)
	}
	if _, err := scheduler.CalculateCriticalPath(saved); err != nil {
		t.Fatalf("更新进度后任务仍应能计算关键路径: %v", err)
	}
	if _, err := os.Stat(filepath.Join("data", "projects.json")); err != nil {
		t.Fatalf("更新后的项目未落盘: %v", err)
	}
}
