package user

import (
	"time"
	"webGL-720yun/internal/model"
	"webGL-720yun/pkg/utils"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	TokenType    string      `json:"token_type"`
	User         *model.User `json:"user"`
}

type RegisterRequest struct {
	ID         uint   `json:"id" binding:"omitempty"` // 学生角色时必填（作为学号），其他角色自增
	Username   string `json:"username" binding:"required,min=3,max=20"`
	Password   string `json:"password" binding:"required,min=6"`
	Email      string `json:"email" binding:"omitempty,email"`
	Phone      string `json:"phone" binding:"omitempty,len=11"`
	Nickname   string `json:"nickname" binding:"omitempty,max=50"`
	RoleID     uint   `json:"role_id" binding:"required"`
	ClassID    uint   `json:"class_id" binding:"omitempty"`
	TeacherIDs []uint `json:"teacher_ids" binding:"omitempty"`
}

type UpdateUserRequest struct {
	Email      string `json:"email" binding:"omitempty,email"`
	Phone      string `json:"phone" binding:"omitempty,len=11"`
	Nickname   string `json:"nickname" binding:"omitempty,max=50"`
	Avatar     string `json:"avatar" binding:"omitempty,url"`
	Status     *int   `json:"status" binding:"omitempty,oneof=0 1"`
	ClassID    uint   `json:"class_id" binding:"omitempty"`
	TeacherIDs []uint `json:"teacher_ids" binding:"omitempty"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type ResetPasswordRequest struct {
	Username string `json:"username" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type UserListRequest struct {
	Page     int    `form:"page" binding:"min=1" json:"page"`
	PageSize int    `form:"page_size" binding:"min=1,max=100" json:"page_size"`
	Username string `form:"username" json:"username"`
	Email    string `form:"email" json:"email"`
	RoleID   uint   `form:"role_id" json:"role_id"`
	Status   int    `form:"status" json:"status"`
	ClassID  uint   `form:"class_id" json:"class_id"`
	Keyword  string `form:"keyword" json:"keyword"`
}

type UserListResponse struct {
	utils.PageInfo `json:"page_info"`
	Users          []*model.User `json:"users"`
}

type ClassRequest struct {
	Name        string `json:"name" binding:"required,max=100"`
	Description string `json:"description" binding:"omitempty,max=255"`
	TeacherID   uint   `json:"teacher_id" binding:"required"`
}

type ClassResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	TeacherID   uint      `json:"teacher_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ClassListResponse struct {
	utils.PageInfo `json:"page_info"`
	Classes        []*model.Class `json:"classes"`
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type BatchRegisterResponse struct {
	SuccessCount int                      `json:"success_count"`
	FailedCount  int                      `json:"failed_count"`
	Results      []BatchRegisterResult    `json:"results"`
	Errors       []map[string]interface{} `json:"errors"`
}

type BatchRegisterResult struct {
	Index     int    `json:"index"`
	Username  string `json:"username"`
	StudentID string `json:"student_id"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}

type StudentExcelData struct {
	Index      int
	Name       string
	StudentID  string
	Email      string
	Phone      string
	ClassName  string
	ClassID    uint
	TeacherIDs []uint
}
