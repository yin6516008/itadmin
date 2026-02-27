package handler

import (
	"backend/internal/errcode"

	"github.com/gin-gonic/gin"
)

// AuthHandler 认证相关 HTTP 处理器（开发阶段使用 mock 实现）。
type AuthHandler struct{}

// NewAuthHandler 创建 AuthHandler。
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

// RegisterRoutes 注册认证相关路由。
func (h *AuthHandler) RegisterRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	{
		auth.POST("/login", h.Login)
	}
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type loginResponse struct {
	Token string   `json:"token"`
	User  authUser `json:"user"`
}

type authUser struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Avatar      string   `json:"avatar"`
	Permissions []string `json:"permissions"`
}

// Login 开发阶段 mock 登录，任意邮箱密码均可登录，返回全权限。
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, errcode.ErrInvalidParam)
		return
	}

	Success(c, loginResponse{
		Token: "dev-token-" + req.Email,
		User: authUser{
			ID:          "1",
			Name:        "管理员",
			Avatar:      "",
			Permissions: []string{"*"},
		},
	})
}
