package service

import (
	"2026-FM247-BackEnd/models"
	"time"
)

type TodoRepository interface {
	CreateTodo(todo *models.Todo) error
	GetTodosByUserID(userID uint) ([]models.Todo, error)
	GetTodoByID(id uint) (*models.Todo, error)
	UpdateTodo(todo *models.Todo) error
	DeleteTodo(id uint) error
	UpdateTodoStatus(id uint, userID uint, status string) error
}

type TodoService struct {
	todoRepository TodoRepository
}

func NewTodoService(todoRepository TodoRepository) *TodoService {
	return &TodoService{todoRepository: todoRepository}
}

// CreateTodo 创建待办事项
func (s *TodoService) CreateTodo(userID uint, title string, description string, startTime *time.Time, deadline *time.Time) string {
	// 验证DDL时间
	if deadline != nil && deadline.Before(time.Now()) {
		return "截止时间不能早于当前时间"
	}

	todo := &models.Todo{
		UserID:      userID,
		Title:       title,
		Description: description,
		Status:      "pending",
		StartTime:   startTime,
		DDL:         deadline,
	}

	if err := s.todoRepository.CreateTodo(todo); err != nil {
		return err.Error()
	}

	return "创建成功"
}

// UpdateTodo 更新待办事项
func (s *TodoService) UpdateTodo(userID, todoID uint, title string, description string, startTime *time.Time, deadline *time.Time) string {
	todo, err := s.todoRepository.GetTodoByID(todoID)
	if err != nil {
		return "待办事项不存在"
	}

	if todo.UserID != userID {
		return "无权限更新该待办事项"
	}

	// 更新字段
	if title != "" {
		todo.Title = title
	}
	if description != "" {
		todo.Description = description
	}
	if startTime != nil {
		todo.StartTime = startTime
	}
	if deadline != nil {
		// 验证DDL时间
		if deadline.Before(time.Now()) {
			return "截止时间不能早于当前时间"
		}
		todo.DDL = deadline
	}

	if err := s.todoRepository.UpdateTodo(todo); err != nil {
		return err.Error()
	}

	return "更新成功"
}

func (s *TodoService) GetTodosByUserID(userID uint) ([]TodoInfo, string) {
	todos, err := s.todoRepository.GetTodosByUserID(userID)
	if err != nil {
		return nil, err.Error()
	}
	var todoInfos []TodoInfo
	for _, todo := range todos {
		todoInfos = append(todoInfos, TodoInfo{
			ID:          todo.ID,
			Title:       todo.Title,
			Description: todo.Description,
			Status:      todo.Status,
			StartTime:   todo.StartTime,
			Deadline:    todo.DDL,
			CompletedAt: todo.CompletedAt,
		})
	}
	return todoInfos, ""
}

func (s *TodoService) GetTodoByID(userID, id uint) (TodoInfo, string) {
	todo, err := s.todoRepository.GetTodoByID(id)
	if err != nil {
		return TodoInfo{}, err.Error()
	}
	if todo.UserID != userID {
		return TodoInfo{}, "无权限访问该待办事项"
	}
	todoinfo := TodoInfo{
		ID:          todo.ID,
		Title:       todo.Title,
		Description: todo.Description,
		Status:      todo.Status,
		StartTime:   todo.StartTime,
		Deadline:    todo.DDL,
		CompletedAt: todo.CompletedAt,
	}
	return todoinfo, ""
}

// DeleteTodo 删除待办事项
func (s *TodoService) DeleteTodo(userID, todoID uint) string {
	todo, err := s.todoRepository.GetTodoByID(todoID)
	if err != nil {
		return "待办事项不存在"
	}
	if todo.UserID != userID {
		return "无权限删除该待办事项"
	}
	if err := s.todoRepository.DeleteTodo(todoID); err != nil {
		return err.Error()
	}
	return "删除成功"
}

// UpdateTodoStatus 更新待办事项状态
func (s *TodoService) UpdateTodoStatus(userID, todoID uint, status string) string {
	// 验证状态
	validStatuses := map[string]bool{
		"pending":     true,
		"in_progress": true,
		"completed":   true,
	}

	if !validStatuses[status] {
		return "无效的状态值"
	}

	err := s.todoRepository.UpdateTodoStatus(todoID, userID, status)
	if err != nil {
		return err.Error()
	}
	return "更新成功"
}
