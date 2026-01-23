package repository

import (
	"2026-FM247-BackEnd/models"
	"fmt"
	"time"

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
	if todo.DDL != nil && todo.DDL.Before(todo.CreatedAt) {
		return fmt.Errorf("截止日期不能早于创建日期")
	}

	return r.db.Create(todo).Error
}

func (r *TodoRepository) GetTodosByUserID(userID uint, filters map[string]interface{}) ([]models.Todo, error) {
	if userID == 0 {
		return nil, fmt.Errorf("无效的用户ID")
	}

	var todos []models.Todo
	query := r.db.Where("user_id = ?", userID)

	// 应用过滤器
	if status, ok := filters["status"]; ok {
		query = query.Where("status = ?", status)
	}

	if startDate, ok := filters["start_date"]; ok {
		query = query.Where("DATE(created_at) >= ?", startDate)
	}

	if endDate, ok := filters["end_date"]; ok {
		query = query.Where("DATE(created_at) <= ?", endDate)
	}

	// 排序
	query = query.Order("created_at DESC")

	err := query.Find(&todos).Error
	return todos, err
}

func (r *TodoRepository) GetTodoByID(id uint) (*models.Todo, error) {

	if id == 0 {
		return nil, fmt.Errorf("无效的ID")
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

// UpdateTodoStatus 更新待办事项状态
func (r *TodoRepository) UpdateTodoStatus(id uint, userID uint, status string) error {
	updates := map[string]interface{}{
		"status": status,
	}

	if status == "completed" {
		now := time.Now()
		updates["completed_at"] = &now
	} else if status == "pending" || status == "in_progress" {
		updates["completed_at"] = nil
	}

	return r.db.Model(&models.Todo{}).
		Where("id = ? AND user_id = ?", id, userID).
		Updates(updates).Error
}
