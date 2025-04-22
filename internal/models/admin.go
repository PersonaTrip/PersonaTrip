package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Admin 管理员模型
type Admin struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"size:50;not null;unique"`
	Email     string    `json:"email" gorm:"size:100;not null;unique"`
	Password  string    `json:"-" gorm:"size:100;not null"` // 不在JSON中返回密码
	Role      string    `json:"role" gorm:"size:20;not null;default:admin"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SetPassword 设置加密后的密码
func (a *Admin) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	a.Password = string(hashedPassword)
	return nil
}

// CheckPassword 检查密码是否正确
func (a *Admin) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(password))
	return err == nil
}

// AdminResponse 是管理员信息的响应格式
type AdminResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

// ToResponse 将Admin转换为AdminResponse
func (a *Admin) ToResponse() AdminResponse {
	return AdminResponse{
		ID:       a.ID,
		Username: a.Username,
		Email:    a.Email,
		Role:     a.Role,
	}
}

// AdminLoginRequest 是管理员登录的请求格式
type AdminLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// AdminCreateRequest 是创建管理员的请求格式
type AdminCreateRequest struct {
	Username string `json:"username" binding:"required,min=4,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
	Role     string `json:"role" binding:"omitempty,oneof=admin super_admin"`
}

// AdminUpdateRequest 是更新管理员信息的请求格式
type AdminUpdateRequest struct {
	Password string `json:"password" binding:"omitempty,min=6"`
	Email    string `json:"email" binding:"omitempty,email"`
	Role     string `json:"role" binding:"omitempty,oneof=admin super_admin"`
}
