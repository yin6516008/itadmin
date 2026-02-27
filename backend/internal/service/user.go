// Package service 实现业务逻辑编排，通过接口依赖 repository。
package service

import (
	"context"
	"errors"
	"log/slog"

	"backend/internal/errcode"
	"backend/internal/infra"
	"backend/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRepo 定义用户数据访问层的接口，由 repository.UserRepo 实现。
type UserRepo interface {
	Save(ctx context.Context, user *model.User) error
	FindByID(ctx context.Context, id string) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, req *model.ListUsersRequest) ([]*model.User, int64, error)
}

// UserService 用户业务逻辑服务。
type UserService struct {
	repo UserRepo
}

// NewUserService 创建 UserService，依赖 UserRepo 接口。
func NewUserService(repo UserRepo) *UserService {
	return &UserService{repo: repo}
}

// Create 创建新用户。
func (s *UserService) Create(ctx context.Context, req *model.CreateUserRequest) (*model.User, error) {
	infra.Logger(ctx).Info("creating user", slog.String("email", req.Email))

	user := &model.User{
		ID:     uuid.NewString(),
		Name:   req.Name,
		Email:  req.Email,
		Phone:  req.Phone,
		Status: "active",
	}
	if err := s.repo.Save(ctx, user); err != nil {
		infra.Logger(ctx).Error("failed to create user", slog.String("error", err.Error()))
		return nil, errcode.ErrInternal
	}
	return user, nil
}

// GetByID 根据 ID 查询用户。
func (s *UserService) GetByID(ctx context.Context, id string) (*model.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrNotFound
		}
		infra.Logger(ctx).Error("failed to get user", slog.String("id", id), slog.String("error", err.Error()))
		return nil, errcode.ErrInternal
	}
	return user, nil
}

// Update 更新用户信息。
func (s *UserService) Update(ctx context.Context, id string, req *model.UpdateUserRequest) error {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errcode.ErrNotFound
		}
		infra.Logger(ctx).Error("failed to find user for update", slog.String("id", id), slog.String("error", err.Error()))
		return errcode.ErrInternal
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}

	if err := s.repo.Update(ctx, user); err != nil {
		infra.Logger(ctx).Error("failed to update user", slog.String("id", id), slog.String("error", err.Error()))
		return errcode.ErrInternal
	}
	return nil
}

// UpdateStatus 切换用户状态（启用/禁用）。
func (s *UserService) UpdateStatus(ctx context.Context, id string, req *model.UpdateUserStatusRequest) error {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errcode.ErrNotFound
		}
		infra.Logger(ctx).Error("failed to find user for status update", slog.String("id", id), slog.String("error", err.Error()))
		return errcode.ErrInternal
	}

	user.Status = req.Status
	if err := s.repo.Update(ctx, user); err != nil {
		infra.Logger(ctx).Error("failed to update user status", slog.String("id", id), slog.String("error", err.Error()))
		return errcode.ErrInternal
	}
	infra.Logger(ctx).Info("user status updated", slog.String("id", id), slog.String("status", req.Status))
	return nil
}

// Delete 删除用户。
func (s *UserService) Delete(ctx context.Context, id string) error {
	if _, err := s.repo.FindByID(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errcode.ErrNotFound
		}
		infra.Logger(ctx).Error("failed to find user for delete", slog.String("id", id), slog.String("error", err.Error()))
		return errcode.ErrInternal
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		infra.Logger(ctx).Error("failed to delete user", slog.String("id", id), slog.String("error", err.Error()))
		return errcode.ErrInternal
	}
	return nil
}

// List 分页查询用户列表，支持关键词搜索和状态筛选。
func (s *UserService) List(ctx context.Context, req *model.ListUsersRequest) ([]*model.User, int64, error) {
	users, total, err := s.repo.List(ctx, req)
	if err != nil {
		infra.Logger(ctx).Error("failed to list users", slog.String("error", err.Error()))
		return nil, 0, errcode.ErrInternal
	}
	return users, total, nil
}
