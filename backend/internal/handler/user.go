package handler

import (
	"backend/internal/errcode"
	"backend/internal/model"
	"backend/internal/service"

	"github.com/gin-gonic/gin"
)

// UserHandler 用户相关 HTTP 处理器。
type UserHandler struct {
	service *service.UserService
}

// NewUserHandler 创建 UserHandler，依赖 UserService。
func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

// RegisterRoutes 将用户相关路由注册到指定路由组。
func (h *UserHandler) RegisterRoutes(rg *gin.RouterGroup) {
	users := rg.Group("/users")
	{
		users.POST("", h.Create)
		users.GET("", h.List)
		users.GET("/:id", h.GetByID)
		users.PUT("/:id", h.Update)
		users.PUT("/:id/status", h.UpdateStatus)
		users.DELETE("/:id", h.Delete)
	}
}

// Create 创建用户。
func (h *UserHandler) Create(c *gin.Context) {
	var req model.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, errcode.ErrInvalidParam)
		return
	}

	user, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		Error(c, err)
		return
	}

	Success(c, user)
}

// GetByID 根据 ID 查询用户。
func (h *UserHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	user, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		Error(c, err)
		return
	}

	Success(c, user)
}

// Update 更新用户信息。
func (h *UserHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req model.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, errcode.ErrInvalidParam)
		return
	}

	if err := h.service.Update(c.Request.Context(), id, &req); err != nil {
		Error(c, err)
		return
	}

	Success(c, nil)
}

// UpdateStatus 切换用户启用/禁用状态。
func (h *UserHandler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")

	var req model.UpdateUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, errcode.ErrInvalidParam)
		return
	}

	if err := h.service.UpdateStatus(c.Request.Context(), id, &req); err != nil {
		Error(c, err)
		return
	}

	Success(c, nil)
}

// Delete 删除用户。
func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		Error(c, err)
		return
	}

	Success(c, nil)
}

// List 分页查询用户列表，支持关键词搜索和状态筛选。
func (h *UserHandler) List(c *gin.Context) {
	var req model.ListUsersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		Error(c, errcode.ErrInvalidParam)
		return
	}

	users, total, err := h.service.List(c.Request.Context(), &req)
	if err != nil {
		Error(c, err)
		return
	}

	Success(c, model.PageResult{
		List:  users,
		Total: total,
		Page:  req.Page,
		Size:  req.Size,
	})
}
