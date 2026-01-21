package repository

import (
	"2026-FM247-BackEnd/models"
	"fmt"

	"gorm.io/gorm"
)

type TodoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) *TodoRepository {
	return &TodoRepository{db: db}
}

func (r *TodoRepository) CreateTodo(todo *models.Todo) error {
	if todo.Title == "" {
		return fmt.Errorf("任务不能为空")
	}
	return r.db.Create(todo).Error
}

func (r *TodoRepository) GetTodosByUserID(userID uint) ([]models.Todo, error) {

	if userID == 0 {
		return nil, fmt.Errorf("无效的用户ID")
	}

	var todos []models.Todo
	err := r.db.Where("user_id = ?", userID).Find(&todos).Error
	return todos, err
}

func (r *TodoRepository) GetTodoByID(id uint) (*models.Todo, error) {

	if id == 0 {
		return nil, fmt.Errorf("无效的图书ID")
	}

	var todo models.Todo
	err := r.db.First(&todo, id).Error
	return &todo, err
}

func (r *TodoRepository) UpdateTodo(todo *models.Todo) error {
	var existing models.Todo
	if err := r.db.First(&existing, todo.ID).Error; err != nil {
		return fmt.Errorf("任务不存在")
	}

	return r.db.Save(todo).Error
}

func (r *TodoRepository) DeleteTodo(id uint) error {

	var todo models.Todo
	if err := r.db.First(&todo, id).Error; err != nil {
		return fmt.Errorf("任务不存在")
	}

	return r.db.Delete(&models.Todo{}, id).Error
}
