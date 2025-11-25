package user

import (
	"time"
	"webGL-720yun/pkg/utils"

	"gorm.io/gorm"
)

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	User         *User  `json:"user"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"omitempty,len=11"`
	Nickname string `json:"nickname" binding:"omitempty,max=50"`
	RoleID   uint   `json:"role_id" binding:"required"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Email    string `json:"email" binding:"omitempty,email"`
	Phone    string `json:"phone" binding:"omitempty,len=11"`
	Nickname string `json:"nickname" binding:"omitempty,max=50"`
	Avatar   string `json:"avatar" binding:"omitempty,url"`
	Status   int    `json:"status" binding:"omitempty,oneof=0 1"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// RefreshTokenRequest 刷新令牌请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshTokenResponse 刷新令牌响应
type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

// UserListRequest 用户列表请求
type UserListRequest struct {
	Page     int    `form:"page" binding:"min=1" json:"page"`
	PageSize int    `form:"page_size" binding:"min=1,max=100" json:"page_size"`
	Username string `form:"username" json:"username"`
	Email    string `form:"email" json:"email"`
	RoleID   uint   `form:"role_id" json:"role_id"`
	Status   int    `form:"status" json:"status"`
	Keyword  string `form:"keyword" json:"keyword"`
}

// UserListResponse 用户列表响应
type UserListResponse struct {
	utils.PageInfo `json:"page_info"`
	Users          []*User `json:"users"`
}

// Response 通用响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// User 用户模型
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Password  string         `gorm:"type:varchar(255);not null" json:"-"`
	Email     string         `gorm:"type:varchar(100);uniqueIndex" json:"email"`
	Phone     string         `gorm:"type:varchar(20);uniqueIndex" json:"phone"`
	Nickname  string         `gorm:"type:varchar(50)" json:"nickname"`
	Avatar    string         `gorm:"type:varchar(255)" json:"avatar"`
	RoleID    uint           `gorm:"not null" json:"role_id"`
	Role      Role           `gorm:"foreignKey:RoleID" json:"role"`
	Status    int            `gorm:"default:1" json:"status"` // 1:正常 0:禁用
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Role 角色模型
type Role struct {
	ID          uint         `gorm:"primaryKey" json:"id"`
	Name        string       `gorm:"type:varchar(50);uniqueIndex;not null" json:"name"`
	Description string       `gorm:"type:varchar(255)" json:"description"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

// Permission 权限模型
type Permission struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"name"`
	Description string    `gorm:"type:varchar(255)" json:"description"`
	Resource    string    `gorm:"type:varchar(50)" json:"resource"`
	Action      string    `gorm:"type:varchar(50)" json:"action"`
	Roles       []Role    `gorm:"many2many:role_permissions;" json:"-"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// UserSession 用户会话模型
type UserSession struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"not null;index" json:"user_id"`
	User         User      `gorm:"foreignKey:UserID" json:"user"`
	AccessToken  string    `gorm:"type:text;not null" json:"-"`
	RefreshToken string    `gorm:"type:text;not null" json:"-"`
	IP           string    `json:"ip"`
	UserAgent    string    `json:"user_agent"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// 用户状态常量
const (
	UserStatusActive   = 1 // 正常
	UserStatusInactive = 0 // 禁用
)

// 角色名称常量
const (
	RoleAdmin   = "admin"
	RoleStudent = "student"
)

// TableName 设置表名
func (User) TableName() string {
	return "users"
}

func (Role) TableName() string {
	return "roles"
}

func (Permission) TableName() string {
	return "permissions"
}

func (UserSession) TableName() string {
	return "user_sessions"
}
