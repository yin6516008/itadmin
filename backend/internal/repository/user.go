// Package repository 实现数据访问层，通过 GORM 操作数据库。
package repository

import (
	"context"

	"backend/internal/model"

	"gorm.io/gorm"
)

// UserRepo 用户数据访问层。
type UserRepo struct {
	db *gorm.DB
}

// NewUserRepo 创建 UserRepo，依赖 GORM 数据库连接。
func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

// Save 创建一条用户记录。
func (r *UserRepo) Save(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// FindByID 根据 ID 查询用户，未找到时返回 gorm.ErrRecordNotFound。
func (r *UserRepo) FindByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail 根据邮箱查询用户。
func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Update 更新用户记录。
func (r *UserRepo) Update(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// Delete 根据 ID 软删除用户记录。
func (r *UserRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&model.User{}, "id = ?", id).Error
}

// List 分页查询用户列表，支持关键词搜索和状态筛选，按创建时间倒序。
func (r *UserRepo) List(ctx context.Context, req *model.ListUsersRequest) ([]*model.User, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.User{})

	if req.Keyword != "" {
		kw := "%" + req.Keyword + "%"
		query = query.Where("name ILIKE ? OR email ILIKE ?", kw, kw)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (req.Page - 1) * req.Size
	var users []*model.User
	if err := query.Order("created_at DESC").Offset(offset).Limit(req.Size).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}
