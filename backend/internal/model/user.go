// Package model 定义业务实体、请求参数和响应 DTO。
package model

import (
	"time"

	"gorm.io/gorm"
)

// ——— 实体 ———

// User 用户实体。
type User struct {
	ID        string         `json:"id" gorm:"primaryKey;size:36"`
	Name      string         `json:"name" gorm:"size:50;not null"`
	Email     string         `json:"email" gorm:"size:100;uniqueIndex;not null"`
	Phone     string         `json:"phone" gorm:"size:20"`
	Status    string         `json:"status" gorm:"size:20;default:active;not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// ——— Request ———

// CreateUserRequest 创建用户请求参数。
type CreateUserRequest struct {
	Name  string `json:"name" binding:"required,max=50"`
	Email string `json:"email" binding:"required,email"`
	Phone string `json:"phone" binding:"max=20"`
}

// UpdateUserRequest 更新用户请求参数。
type UpdateUserRequest struct {
	Name  string `json:"name" binding:"max=50"`
	Email string `json:"email" binding:"omitempty,email"`
	Phone string `json:"phone" binding:"max=20"`
}

// UpdateUserStatusRequest 切换用户状态请求参数。
type UpdateUserStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=active inactive"`
}

// ListUsersRequest 用户列表查询参数。
type ListUsersRequest struct {
	Page    int    `form:"page" binding:"min=1"`
	Size    int    `form:"size" binding:"min=1,max=100"`
	Keyword string `form:"keyword"`
	Status  string `form:"status"`
}

// ——— Response ———

// PageResult 通用分页响应结构。
type PageResult struct {
	List  any   `json:"list"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Size  int   `json:"size"`
}
