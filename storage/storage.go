package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"gantt-generator/models"
)

const (
	dataDir     = "data"
	projectFile = "projects.json"
)

// Store JSON存储管理器
type Store struct {
	mu       sync.RWMutex
	filePath string
}

// NewStore 创建新的存储管理器
func NewStore() (*Store, error) {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("创建数据目录失败: %w", err)
	}
	filePath := filepath.Join(dataDir, projectFile)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		if err := os.WriteFile(filePath, []byte("[]"), 0644); err != nil {
			return nil, fmt.Errorf("初始化数据文件失败: %w", err)
		}
	}
	return &Store{filePath: filePath}, nil
}

// loadProjects 从文件加载所有项目
func (s *Store) loadProjects() ([]models.Project, error) {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return nil, fmt.Errorf("读取数据文件失败: %w", err)
	}
	var projects []models.Project
	if len(data) == 0 {
		return projects, nil
	}
	if err := json.Unmarshal(data, &projects); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %w", err)
	}
	return projects, nil
}

// saveProjects 保存所有项目到文件
func (s *Store) saveProjects(projects []models.Project) error {
	data, err := json.MarshalIndent(projects, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化JSON失败: %w", err)
	}
	if err := os.WriteFile(s.filePath, data, 0644); err != nil {
		return fmt.Errorf("写入数据文件失败: %w", err)
	}
	return nil
}

// GetAllProjects 获取所有项目
func (s *Store) GetAllProjects() ([]models.Project, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.loadProjects()
}

// GetProject 根据ID获取项目
func (s *Store) GetProject(id string) (*models.Project, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	projects, err := s.loadProjects()
	if err != nil {
		return nil, err
	}
	for i := range projects {
		if projects[i].ID == id {
			return &projects[i], nil
		}
	}
	return nil, fmt.Errorf("项目不存在: %s", id)
}

// CreateProject 创建新项目
func (s *Store) CreateProject(project *models.Project) error {
	if err := models.ValidateProjectTasks(project); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	projects, err := s.loadProjects()
	if err != nil {
		return err
	}
	now := time.Now().Format(time.RFC3339)
	project.CreatedAt = now
	project.UpdatedAt = now
	if project.Tasks == nil {
		project.Tasks = []models.Task{}
	}
	projects = append(projects, *project)
	return s.saveProjects(projects)
}

// UpdateProject 更新项目
func (s *Store) UpdateProject(project *models.Project) error {
	if err := models.ValidateProjectTasks(project); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	projects, err := s.loadProjects()
	if err != nil {
		return err
	}
	found := false
	for i := range projects {
		if projects[i].ID == project.ID {
			project.CreatedAt = projects[i].CreatedAt
			project.UpdatedAt = time.Now().Format(time.RFC3339)
			projects[i] = *project
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("项目不存在: %s", project.ID)
	}
	return s.saveProjects(projects)
}

// DeleteProject 删除项目
func (s *Store) DeleteProject(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	projects, err := s.loadProjects()
	if err != nil {
		return err
	}
	newProjects := []models.Project{}
	found := false
	for i := range projects {
		if projects[i].ID != id {
			newProjects = append(newProjects, projects[i])
		} else {
			found = true
		}
	}
	if !found {
		return fmt.Errorf("项目不存在: %s", id)
	}
	return s.saveProjects(newProjects)
}

// AddTask 向项目添加任务
func (s *Store) AddTask(projectID string, task *models.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	projects, err := s.loadProjects()
	if err != nil {
		return err
	}
	found := false
	for i := range projects {
		if projects[i].ID == projectID {
			if projects[i].Tasks == nil {
				projects[i].Tasks = []models.Task{}
			}
			projects[i].Tasks = append(projects[i].Tasks, *task)
			projects[i].UpdatedAt = time.Now().Format(time.RFC3339)
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("项目不存在: %s", projectID)
	}
	for i := range projects {
		if projects[i].ID == projectID {
			if err := models.ValidateProjectTasks(&projects[i]); err != nil {
				return err
			}
			break
		}
	}
	return s.saveProjects(projects)
}

// UpdateTask 更新项目中的任务
func (s *Store) UpdateTask(projectID string, task *models.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	projects, err := s.loadProjects()
	if err != nil {
		return err
	}
	projectFound := false
	taskFound := false
	for i := range projects {
		if projects[i].ID == projectID {
			projectFound = true
			for j := range projects[i].Tasks {
				if projects[i].Tasks[j].ID == task.ID {
					projects[i].Tasks[j] = *task
					projects[i].UpdatedAt = time.Now().Format(time.RFC3339)
					taskFound = true
					break
				}
			}
			break
		}
	}
	if !projectFound {
		return fmt.Errorf("项目不存在: %s", projectID)
	}
	if !taskFound {
		return fmt.Errorf("任务不存在: %s", task.ID)
	}
	return s.saveProjects(projects)
}

// DeleteTask 从项目中删除任务
func (s *Store) DeleteTask(projectID string, taskID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	projects, err := s.loadProjects()
	if err != nil {
		return err
	}
	projectFound := false
	taskFound := false
	for i := range projects {
		if projects[i].ID == projectID {
			projectFound = true
			newTasks := []models.Task{}
			for j := range projects[i].Tasks {
				if projects[i].Tasks[j].ID != taskID {
					newTasks = append(newTasks, projects[i].Tasks[j])
				} else {
					taskFound = true
				}
			}
			projects[i].Tasks = newTasks
			projects[i].UpdatedAt = time.Now().Format(time.RFC3339)
			break
		}
	}
	if !projectFound {
		return fmt.Errorf("项目不存在: %s", projectID)
	}
	if !taskFound {
		return fmt.Errorf("任务不存在: %s", taskID)
	}
	return s.saveProjects(projects)
}
