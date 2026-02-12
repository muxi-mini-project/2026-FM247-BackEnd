package handler

import (
	"2026-FM247-BackEnd/service"
	"2026-FM247-BackEnd/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type TodoService interface {
	CreateTodo(userID uint, title string, description string, startTime *time.Time, deadline *time.Time) string
	UpdateTodo(userID, todoID uint, title string, description string, startTime *time.Time, deadline *time.Time) string
	GetTodosByUserID(userID uint) ([]service.TodoInfo, string)
	GetTodoByID(userID, id uint) (service.TodoInfo, string)
	DeleteTodo(userID, todoID uint) string
	UpdateTodoStatus(userID, todoID uint, status string) string
}

type TodoHandler struct {
	todoService *service.TodoService
}

func NewTodoHandler(todoService *service.TodoService) *TodoHandler {
	return &TodoHandler{todoService: todoService}
}

// CreateTodo 创建待办事项
// @Router /api/todos [post]
func (h *TodoHandler) CreateTodo(c *gin.Context) {
	// 1. 验证登录
	claims, err := utils.GetClaimsFromContext(c)
	if err != nil {
		FailWithMessage(c, "无法获取用户信息: "+err.Error())
		return
	}
	// 2. 绑定请求参数
	var req CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		FailWithMessage(c, "请求参数有误: "+err.Error())
		return
	}
	// 3. 创建待办事项
	msg := h.todoService.CreateTodo(claims.UserID, req.Title, req.Description, req.StartTime, req.Deadline)
	if msg != "创建成功" {
		FailWithMessage(c, "创建失败: "+msg)
		return
	}
	// 4. 返回结果
	OkWithMessage(c, msg)
}

// GetTodos 获取用户的待办事项列表
// @Router /api/todos [get]
func (h *TodoHandler) GetTodos(c *gin.Context) {
	//验证登录
	claims, err := utils.GetClaimsFromContext(c)
	if err != nil {
		FailWithMessage(c, "无法获取用户信息: "+err.Error())
		return
	}
	//获取待办事项列表
	todos, msg := h.todoService.GetTodosByUserID(claims.UserID)
	if msg != "" {
		FailWithMessage(c, "获取失败: "+msg)
		return
	}
	// 4. 返回结果
	OkWithData(c, todos)
}

// GetTodoByID 获取单个待办事项
// @Router /api/todos/:id [get]
func (h *TodoHandler) GetTodoByID(c *gin.Context) {
	// 1. 验证登录
	claims, err := utils.GetClaimsFromContext(c)
	if err != nil {
		FailWithMessage(c, "无法获取用户信息: "+err.Error())
		return
	}
	// 2. 获取待办事项ID
	todoIDStr := c.Param("id")
	todoID, err := strconv.ParseUint(todoIDStr, 10, 64)
	if err != nil {
		FailWithMessage(c, "无效的待办事项ID")
		return
	}
	// 3. 获取待办事项
	todo, msg := h.todoService.GetTodoByID(claims.UserID, uint(todoID))
	if msg != "" {
		FailWithMessage(c, "获取失败: "+msg)
		return
	}
	// 4. 返回结果
	OkWithData(c, todo)
}

// UpdateTodo 更新待办事项
// @Router /api/todos/:id [put]
func (h *TodoHandler) UpdateTodo(c *gin.Context) {
	// 1. 验证登录
	claims, err := utils.GetClaimsFromContext(c)
	if err != nil {
		FailWithMessage(c, "无法获取用户信息: "+err.Error())
		return
	}
	// 2. 获取待办事项ID
	todoIDStr := c.Param("id")
	todoID, err := strconv.ParseUint(todoIDStr, 10, 64)
	if err != nil {
		FailWithMessage(c, "无效的待办事项ID")
		return
	}
	// 3. 绑定请求参数
	var req UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		FailWithMessage(c, "请求参数有误: "+err.Error())
		return
	}
	// 4. 更新待办事项
	msg := h.todoService.UpdateTodo(claims.UserID, uint(todoID), req.Title, req.Description, req.StartTime, req.Deadline)
	if msg != "更新成功" {
		FailWithMessage(c, "更新失败: "+msg)
		return
	}
	// 5. 返回结果
	OkWithMessage(c, msg)
}

// DeleteTodo 删除待办事项
// @Router /api/todos/:id [delete]
func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	// 1. 验证登录
	claims, err := utils.GetClaimsFromContext(c)
	if err != nil {
		FailWithMessage(c, "无法获取用户信息: "+err.Error())
		return
	}
	// 2. 获取待办事项ID
	todoIDStr := c.Param("id")
	todoID, err := strconv.ParseUint(todoIDStr, 10, 64)
	if err != nil {
		FailWithMessage(c, "无效的待办事项ID")
		return
	}
	// 3. 删除待办事项
	msg := h.todoService.DeleteTodo(claims.UserID, uint(todoID))
	if msg != "删除成功" {
		FailWithMessage(c, "删除失败: "+msg)
		return
	}
	// 4. 返回结果
	OkWithMessage(c, "删除成功")
}

// UpdateTodoStatus 更新待办事项状态
// @Router /api/todos/:id/status [patch]
func (h *TodoHandler) UpdateTodoStatus(c *gin.Context) {
	// 1. 验证登录
	claims, err := utils.GetClaimsFromContext(c)
	if err != nil {
		FailWithMessage(c, "无法获取用户信息: "+err.Error())
		return
	}
	// 2. 获取待办事项ID
	todoIDStr := c.Param("id")
	todoID, err := strconv.ParseUint(todoIDStr, 10, 64)
	if err != nil {
		FailWithMessage(c, "无效的待办事项ID")
		return
	}
	// 3. 绑定请求参数
	var req UpdateTodoStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		FailWithMessage(c, "请求参数有误: "+err.Error())
		return
	}
	// 4. 更新待办事项状态
	msg := h.todoService.UpdateTodoStatus(claims.UserID, uint(todoID), req.Status)
	if msg != "更新成功" {
		FailWithMessage(c, "更新失败: "+msg)
		return
	}

	// 5. 返回更新后的待办事项
	todo, msg := h.todoService.GetTodoByID(claims.UserID, uint(todoID))
	if msg != "" {
		FailWithMessage(c, "获取失败: "+msg)
		return
	}

	OkWithData(c, todo)
}
