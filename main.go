package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"gantt-generator/handlers"
	"gantt-generator/storage"
)

func main() {
	// 初始化存储
	store, err := storage.NewStore()
	if err != nil {
		log.Fatalf("初始化存储失败: %v", err)
	}

	// 创建处理器
	h := handlers.NewHandler(store)

	// 注册路由
	mux := http.NewServeMux()

	// 静态文件服务
	fileServer := http.FileServer(http.Dir("web"))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 如果是API请求，不处理
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}
		// 根路径返回index.html
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "web/index.html")
			return
		}
		// 其他静态文件
		if _, err := os.Stat("web" + r.URL.Path); err == nil {
			fileServer.ServeHTTP(w, r)
		} else {
			// SPA fallback: 返回index.html
			http.ServeFile(w, r, "web/index.html")
		}
	})

	// ========== 项目API ==========
	// GET /api/projects - 获取所有项目
	mux.HandleFunc("/api/projects", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.GetProjects(w, r)
		case http.MethodPost:
			h.CreateProject(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// GET/PUT/DELETE /api/projects/{id}
	mux.HandleFunc("/api/projects/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/api/projects/")
		// 检查是否包含子路径: {id}/tasks, {id}/critical-path, {id}/gantt-pdf, {id}/report
		parts := strings.Split(path, "/")

		if len(parts) == 1 && parts[0] != "" {
			// /api/projects/{id}
			switch r.Method {
			case http.MethodGet:
				h.GetProject(w, r)
			case http.MethodPut:
				h.UpdateProject(w, r)
			case http.MethodDelete:
				h.DeleteProject(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		if len(parts) >= 2 {
			_ = parts[0]
			subPath := parts[1]

			switch subPath {
			case "tasks":
				// /api/projects/{id}/tasks
				if len(parts) == 2 {
					switch r.Method {
					case http.MethodPost:
						h.AddTask(w, r)
					default:
						http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
					}
				} else if len(parts) == 3 {
					// /api/projects/{id}/tasks/{taskId}
					switch r.Method {
					case http.MethodPut:
						h.UpdateTask(w, r)
					case http.MethodDelete:
						h.DeleteTask(w, r)
					default:
						http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
					}
				}
			case "critical-path":
				if r.Method == http.MethodGet {
					h.CalculateCriticalPath(w, r)
				} else {
					http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				}
			case "gantt-pdf":
				if r.Method == http.MethodGet {
					h.GenerateGanttPDF(w, r)
				} else {
					http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				}
			case "report":
				if r.Method == http.MethodGet {
					h.GetProgressReport(w, r)
				} else {
					http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				}
			default:
				http.NotFound(w, r)
			}
			return
		}

		http.NotFound(w, r)
	})

	// 启动服务器
	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}
	addr := fmt.Sprintf(":%s", port)

	fmt.Printf("🚀 工程项目进度管理系统启动中...\n")
	fmt.Printf("📡 服务地址: http://localhost%s\n", addr)
	fmt.Printf("📂 数据存储目录: ./data\n")
	fmt.Printf("📖 前端页面: http://localhost%s/\n", addr)
	fmt.Println("-----------------------------")
	fmt.Println("API 接口说明:")
	fmt.Println("  GET    /api/projects              获取所有项目")
	fmt.Println("  POST   /api/projects              创建新项目")
	fmt.Println("  GET    /api/projects/{id}         获取项目详情")
	fmt.Println("  PUT    /api/projects/{id}         更新项目")
	fmt.Println("  DELETE /api/projects/{id}         删除项目")
	fmt.Println("  POST   /api/projects/{id}/tasks   添加任务")
	fmt.Println("  PUT    /api/projects/{id}/tasks/{taskId}   更新任务")
	fmt.Println("  DELETE /api/projects/{id}/tasks/{taskId}   删除任务")
	fmt.Println("  GET    /api/projects/{id}/critical-path    计算关键路径")
	fmt.Println("  GET    /api/projects/{id}/gantt-pdf        生成甘特图PDF")
	fmt.Println("  GET    /api/projects/{id}/report           获取进度报告 (format=text获取文本)")
	fmt.Println("-----------------------------")

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
