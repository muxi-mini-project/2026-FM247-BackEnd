package service

import (
	"2026-FM247-BackEnd/models"
)

type TodoRepository interface {
	CreateTodo(todo *models.Todo) error
	GetTodosByUserID(userID uint) ([]models.Todo, error)
	GetTodoByID(id uint) (*models.Todo, error)
	UpdateTodo(todo *models.Todo) error
	DeleteTodo(id uint) error
}

type TodoService struct {
	todoRepository TodoRepository
}

func NewTodoService(todoRepository TodoRepository) *TodoService {
	return &TodoService{todoRepository: todoRepository}
}

// CreateTodo 创建待办事项
func (s *TodoService) CreateTodo(userID uint, event string) string {
	todo := &models.Todo{
		UserID: userID,
		Event:  event,
	}

	if err := s.todoRepository.CreateTodo(todo); err != nil {
		return err.Error()
	}

	return "创建成功"
}

// UpdateTodo 更新待办事项
func (s *TodoService) UpdateTodo(userID, todoID uint, event string) string {
	todo, err := s.todoRepository.GetTodoByID(todoID)
	if err != nil {
		return "待办事项不存在"
	}

	if todo.UserID != userID {
		return "无权限更新该待办事项"
	}

	// 更新字段
	if event != "" {
		todo.Event = event
	}

	if err := s.todoRepository.UpdateTodo(todo); err != nil {
		return err.Error()
	}

	return "更新成功"
}

func (s *TodoService) GetTodosByUserID(userID uint) ([]TodoInfo, string, bool) {
	todos, err := s.todoRepository.GetTodosByUserID(userID)
	if err != nil {
		return nil, err.Error(), false
	}
	if len(todos) == 0 {
		return nil, "", false
	}
	var todoInfos []TodoInfo
	for _, todo := range todos {
		todoInfos = append(todoInfos, TodoInfo{
			ID:    todo.ID,
			Event: todo.Event,
		})
	}
	return todoInfos, "", true
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
		ID:    todo.ID,
		Event: todo.Event,
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
