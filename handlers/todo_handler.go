package handler

import (
	"2026-FM247-BackEnd/service"
	"2026-FM247-BackEnd/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

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
		FailWithMessage(c, "请先登录")
		return
	}
	// 2. 绑定请求参数
	var req service.CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		FailWithMessage(c, "参数绑定失败: "+err.Error())
		return
	}
	// 3. 创建待办事项
	todo, err := h.todoService.CreateTodo(claims.UserID, req)
	if err != nil {
		FailWithMessage(c, "创建失败: "+err.Error())
		return
	}
	// 4. 返回结果
	OkWithData(c, todo)
}

// GetTodos 获取用户的待办事项列表
// @Router /api/todos [get]
func (h *TodoHandler) GetTodos(c *gin.Context) {
	// 1. 验证登录
	claims, err := utils.GetClaimsFromContext(c)
	if err != nil {
		FailWithMessage(c, "请先登录")
		return
	}
	// 2. 获取查询参数
	userID := claims.UserID
	filters := make(map[string]interface{})
	status := c.Query("status")
	if status != "" {
		filters["status"] = status
	}
	startDate := c.Query("start_date")
	if startDate != "" {
		filters["start_date"] = startDate
	}
	endDate := c.Query("end_date")
	if endDate != "" {
		filters["end_date"] = endDate
	}
	// 3. 获取待办事项列表
	todos, err := h.todoService.GetTodosByUserID(userID, filters)
	if err != nil {
		FailWithMessage(c, "获取失败: "+err.Error())
		return
	}
	// 4. 返回结果
	OkWithData(c, todos)
}

// GetTodo 获取单个待办事项
// @Router /api/todos/:id [get]
func (h *TodoHandler) GetTodo(c *gin.Context) {
	// 1. 验证登录
	_, err := utils.GetClaimsFromContext(c)
	if err != nil {
		FailWithMessage(c, "请先登录")
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
	todo, err := h.todoService.GetTodoByID(uint(todoID))
	if err != nil {
		FailWithMessage(c, "获取失败: "+err.Error())
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
		FailWithMessage(c, "请先登录")
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
	var req service.UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		FailWithMessage(c, "参数绑定失败: "+err.Error())
		return
	}
	// 4. 更新待办事项
	todo, err := h.todoService.UpdateTodo(claims.UserID, uint(todoID), req)
	if err != nil {
		FailWithMessage(c, "更新失败: "+err.Error())
		return
	}
	// 5. 返回结果
	OkWithData(c, todo)
}

// DeleteTodo 删除待办事项
// @Router /api/todos/:id [delete]
func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	// 1. 验证登录
	claims, err := utils.GetClaimsFromContext(c)
	if err != nil {
		FailWithMessage(c, "请先登录")
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
	err = h.todoService.DeleteTodo(claims.UserID, uint(todoID))
	if err != nil {
		FailWithMessage(c, "删除失败: "+err.Error())
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
		FailWithMessage(c, "请先登录")
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
	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		FailWithMessage(c, "参数绑定失败: "+err.Error())
		return
	}
	// 4. 更新待办事项状态
	if err := h.todoService.UpdateTodoStatus(claims.UserID, uint(todoID), req.Status); err != nil {
		FailWithMessage(c, "更新失败: "+err.Error())
		return
	}

	// 5. 返回更新后的待办事项
	todo, err := h.todoService.GetTodoByID(uint(todoID))
	if err != nil {
		FailWithMessage(c, "获取失败: "+err.Error())
		return
	}

	OkWithData(c, todo)
}
