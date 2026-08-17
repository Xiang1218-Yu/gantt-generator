package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"gantt-generator/models"
	"gantt-generator/storage"
)

func TestProjectUpdateRejectsDuplicateTaskIDs(t *testing.T) {
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
	if err := store.CreateProject(&models.Project{ID: "project-1", Name: "数据迁移"}); err != nil {
		t.Fatal(err)
	}

	handler := NewHandler(store)
	body := `{
        "name":"数据迁移",
        "tasks":[
            {"id":"approval","name":"业务审批","start_date":"2026-08-01","end_date":"2026-08-02","status":"未开始","progress":0},
            {"id":"approval","name":"安全审批","start_date":"2026-08-03","end_date":"2026-08-04","status":"未开始","progress":0}
        ]
    }`
	request := httptest.NewRequest(http.MethodPut, "/api/projects/project-1", bytes.NewBufferString(body))
	response := httptest.NewRecorder()
	handler.UpdateProject(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("重复任务ID应被拒绝，实际状态码为 %d，响应为 %s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), "重复") {
		t.Fatalf("响应应说明重复ID，实际为 %s", response.Body.String())
	}
}
