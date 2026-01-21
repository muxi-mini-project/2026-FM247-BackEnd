package service

import (
	"2026-FM247-BackEnd/models"
	"2026-FM247-BackEnd/repositories"
)

type ITodoService interface {
	CreateTodo(todo *models.Todo) error
	GetTodosByUserID(userID uint) ([]models.Todo, error)
	GetTodoByID(id uint) (*models.Todo, error)
	UpdateTodo(todo *models.Todo) error
	DeleteTodo(id uint) error
}

type TodoService struct {
	todoRepository *repository.TodoRepository
}

func NewTodoService(todoRepository *repository.TodoRepository) *TodoService {
	return &TodoService{todoRepository: todoRepository}
}

func (s *TodoService) CreateTodo(todo *models.Todo) error {
	return s.todoRepository.CreateTodo(todo)
}

func (s *TodoService) GetTodosByUserID(userID uint) ([]models.Todo, error) {
	return s.todoRepository.GetTodosByUserID(userID)
}

func (s *TodoService) GetTodoByID(id uint) (*models.Todo, error) {
	return s.todoRepository.GetTodoByID(id)
}

func (s *TodoService) UpdateTodo(todo *models.Todo) error {
	return s.todoRepository.UpdateTodo(todo)
}

func (s *TodoService) DeleteTodo(id uint) error {
	return s.todoRepository.DeleteTodo(id)
}
