package service

import (
	"2026-FM247-BackEnd/models"
	"2026-FM247-BackEnd/repositories"
	"errors"
	"time"
)

type ITodoService interface {
	CreateTodo(todo *models.Todo) error
	GetTodosByUserID(userID uint) ([]models.Todo, error)
	GetTodoByID(id uint) (*models.Todo, error)
	UpdateTodo(todo *models.Todo) error
	DeleteTodo(id uint) error
}

type CreateTodoRequest struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	StartTime   *time.Time `json:"start_time"`
	Deadline    *time.Time `json:"deadline" binding:"required"` // DDL为必填
}

// UpdateTodoRequest 更新待办事项请求
type UpdateTodoRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status" binding:"oneof=pending in_progress completed"`
	StartTime   *time.Time `json:"start_time"`
	Deadline    *time.Time `json:"deadline"`
}

// TodoResponse 待办事项响应
type TodoResponse struct {
	ID          uint       `json:"id"`
	UserID      uint       `json:"user_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	StartTime   *time.Time `json:"start_time"`
	Deadline    *time.Time `json:"deadline"`
	CompletedAt *time.Time `json:"completed_at"`
}

type TodoService struct {
	todoRepository *repository.TodoRepository
}

func NewTodoService(todoRepository *repository.TodoRepository) *TodoService {
	return &TodoService{todoRepository: todoRepository}
}

// CreateTodo 创建待办事项
func (s *TodoService) CreateTodo(userID uint, req CreateTodoRequest) (*models.Todo, error) {
	// 验证DDL时间
	if req.Deadline != nil && req.Deadline.Before(time.Now()) {
		return nil, errors.New("截止时间不能早于当前时间")
	}

	todo := &models.Todo{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Status:      "pending",
		StartTime:   req.StartTime,
		DDL:         req.Deadline,
	}

	if err := s.todoRepository.CreateTodo(todo); err != nil {
		return nil, err
	}

	return todo, nil
}

// UpdateTodo 更新待办事项
func (s *TodoService) UpdateTodo(userID, todoID uint, req UpdateTodoRequest) (*models.Todo, error) {
	todo, err := s.todoRepository.GetTodoByID(todoID)
	if err != nil {
		return nil, errors.New("待办事项不存在或无权限修改")
	}

	// 更新字段
	if req.Title != "" {
		todo.Title = req.Title
	}
	if req.Description != "" {
		todo.Description = req.Description
	}
	if req.Status != "" {
		todo.Status = req.Status
		if req.Status == "completed" {
			now := time.Now()
			todo.CompletedAt = &now
		} else {
			todo.CompletedAt = nil
		}
	}
	if req.StartTime != nil {
		todo.StartTime = req.StartTime
	}
	if req.Deadline != nil {
		// 验证DDL时间
		if req.Deadline.Before(time.Now()) {
			return nil, errors.New("截止时间不能早于当前时间")
		}
		todo.DDL = req.Deadline
	}

	if err := s.todoRepository.UpdateTodo(todo); err != nil {
		return nil, err
	}

	return todo, nil
}

func (s *TodoService) GetTodosByUserID(userID uint, filters map[string]interface{}) ([]models.Todo, error) {
	return s.todoRepository.GetTodosByUserID(userID, filters)
}

func (s *TodoService) GetTodoByID(id uint) (*models.Todo, error) {
	return s.todoRepository.GetTodoByID(id)
}

// DeleteTodo 删除待办事项
func (s *TodoService) DeleteTodo(userID, todoID uint) error {
	return s.todoRepository.DeleteTodo(todoID)
}

// UpdateTodoStatus 更新待办事项状态
func (s *TodoService) UpdateTodoStatus(userID, todoID uint, status string) error {
	// 验证状态
	validStatuses := map[string]bool{
		"pending":     true,
		"in_progress": true,
		"completed":   true,
	}

	if !validStatuses[status] {
		return errors.New("无效的状态值")
	}

	return s.todoRepository.UpdateTodoStatus(todoID, userID, status)
}
