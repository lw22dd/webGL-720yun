package user

import (
	"time"
	"webGL-720yun/pkg/utils"

	"gorm.io/gorm"
)

// 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 登录响应
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	User         *User  `json:"user"`
}

// 注册请求
type RegisterRequest struct {
	Username   string `json:"username" binding:"required,min=3,max=20"`
	Password   string `json:"password" binding:"required,min=6"`
	Email      string `json:"email" binding:"required,email"`
	Phone      string `json:"phone" binding:"omitempty,len=11"`
	Nickname   string `json:"nickname" binding:"omitempty,max=50"`
	RoleID     uint   `json:"role_id" binding:"required"`
	StudentID  string `json:"student_id" binding:"omitempty"`  // 学号，仅学生角色需要
	ClassID    uint   `json:"class_id" binding:"omitempty"`    // 班级ID，仅学生角色需要
	TeacherIDs []uint `json:"teacher_ids" binding:"omitempty"` // 关联教师ID，仅学生角色需要
}

// 更新用户请求
type UpdateUserRequest struct {
	Email      string `json:"email" binding:"omitempty,email"`
	Phone      string `json:"phone" binding:"omitempty,len=11"`
	Nickname   string `json:"nickname" binding:"omitempty,max=50"`
	Avatar     string `json:"avatar" binding:"omitempty,url"`
	Status     *int   `json:"status" binding:"omitempty,oneof=0 1"`
	ClassID    uint   `json:"class_id" binding:"omitempty"`    // 班级ID，仅学生角色需要
	TeacherIDs []uint `json:"teacher_ids" binding:"omitempty"` // 关联教师ID，仅学生角色需要
}

// 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// 重置密码请求
type ResetPasswordRequest struct {
	Username string `json:"username" binding:"required"`
}

// 刷新令牌请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// 刷新令牌响应
type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

// 用户列表请求
type UserListRequest struct {
	Page     int    `form:"page" binding:"min=1" json:"page"`
	PageSize int    `form:"page_size" binding:"min=1,max=100" json:"page_size"`
	Username string `form:"username" json:"username"`
	Email    string `form:"email" json:"email"`
	RoleID   uint   `form:"role_id" json:"role_id"`
	Status   int    `form:"status" json:"status"`
	ClassID  uint   `form:"class_id" json:"class_id"` // 班级ID，用于查询特定班级的学生
	Keyword  string `form:"keyword" json:"keyword"`
}

// 用户列表响应
type UserListResponse struct {
	utils.PageInfo `json:"page_info"`
	Users          []*User `json:"users"`
}

// 班级请求
type ClassRequest struct {
	Name        string `json:"name" binding:"required,max=100"`
	Description string `json:"description" binding:"omitempty,max=255"`
	TeacherID   uint   `json:"teacher_id" binding:"required"`
}

// 班级响应
type ClassResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	TeacherID   uint      `json:"teacher_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// 班级列表响应
type ClassListResponse struct {
	utils.PageInfo `json:"page_info"`
	Classes        []*Class `json:"classes"`
}

// 通用响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// 用户抽象基类
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Password  string         `gorm:"type:varchar(255);not null" json:"-"`
	Email     string         `gorm:"type:varchar(100);uniqueIndex" json:"email"`
	Phone     string         `gorm:"type:varchar(20)" json:"phone"`
	Nickname  string         `gorm:"type:varchar(50)" json:"nickname"`
	Avatar    string         `gorm:"type:varchar(255)" json:"avatar"`
	RoleID       uint           `gorm:"not null;index" json:"role_id"`
	Role         Role           `gorm:"foreignKey:RoleID" json:"role"`
	IsSuperAdmin bool           `gorm:"default:false" json:"is_super_admin"`
	Status       int            `gorm:"default:1;index" json:"status"` // 1:正常 0:禁用
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// 教师类
type Teacher struct {
	User
	Classes []Class `gorm:"foreignKey:TeacherID" json:"classes"`
}

// 学生类
type Student struct {
	User
	StudentID string    `gorm:"type:varchar(20);uniqueIndex;not null" json:"student_id"` // 学号
	ClassID   uint      `gorm:"not null;index" json:"class_id"`
	Class     Class     `gorm:"foreignKey:ClassID" json:"class"`
	Teachers  []Teacher `gorm:"many2many:sys_student_teachers;" json:"teachers"` // 关联多个教师
}

// 班级类
type Class struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"name"`
	Description string    `gorm:"type:varchar(255)" json:"description"`
	TeacherID   uint      `gorm:"not null" json:"teacher_id"`
	Teacher     User      `gorm:"foreignKey:TeacherID" json:"teacher"`
	Students    []Student `gorm:"foreignKey:ClassID" json:"students"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// 学生-教师关联表
type StudentTeacher struct {
	StudentID uint `gorm:"primaryKey" json:"student_id"`
	TeacherID uint `gorm:"primaryKey" json:"teacher_id"`
}

// 角色模型
type Role struct {
	ID          uint         `gorm:"primaryKey" json:"id"`
	Name        string       `gorm:"type:varchar(50);uniqueIndex;not null" json:"name"`
	Description string       `gorm:"type:varchar(255)" json:"description"`
	Permissions []Permission `gorm:"many2many:sys_role_permissions;" json:"permissions"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

// 权限模型
type Permission struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"name"`
	Description string    `gorm:"type:varchar(255)" json:"description"`
	Resource    string    `gorm:"type:varchar(50)" json:"resource"`
	Action      string    `gorm:"type:varchar(50)" json:"action"`
	Roles       []Role    `gorm:"many2many:sys_role_permissions;" json:"-"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// 用户会话模型
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
	RoleTeacher = "teacher"
	RoleStudent = "student"
)

// 设置表名
func (User) TableName() string {
	return "sys_users"
}

func (Teacher) TableName() string {
	return "sys_users"
}

func (Student) TableName() string {
	return "sys_users"
}

func (Role) TableName() string {
	return "sys_roles"
}

func (Permission) TableName() string {
	return "sys_permissions"
}

func (UserSession) TableName() string {
	return "sys_user_sessions"
}

func (Class) TableName() string {
	return "sys_classes"
}

func (StudentTeacher) TableName() string {
	return "sys_student_teachers"
}

// 重置密码方法
func (u *User) ResetPassword() (string, error) {
	// 生成临时密码
	tempPassword := utils.GenerateRandomPassword(8)
	// 加密临时密码
	hashedPassword, err := utils.HashPassword(tempPassword)
	if err != nil {
		return "", err
	}
	// 更新密码
	u.Password = hashedPassword
	return tempPassword, nil
}

// 修改密码方法
func (u *User) ChangePassword(oldPassword, newPassword string) error {
	// 验证原密码
	if !utils.CheckPassword(oldPassword, u.Password) {
		return utils.ErrOldPasswordIncorrect
	}
	// 加密新密码
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}
	// 更新密码
	u.Password = hashedPassword
	return nil
}

// 将User转换为Teacher
func (u *User) ToTeacher() *Teacher {
	return &Teacher{User: *u}
}

// 将User转换为Student
func (u *User) ToStudent() *Student {
	return &Student{User: *u}
}

// 批量注册相关结构体

// 批量注册响应
type BatchRegisterResponse struct {
	SuccessCount int                      `json:"success_count"` // 成功数量
	FailedCount  int                      `json:"failed_count"`  // 失败数量
	Results      []BatchRegisterResult    `json:"results"`       // 详细结果
	Errors       []map[string]interface{} `json:"errors"`        // 错误信息
}

// 单个学生注册结果
type BatchRegisterResult struct {
	Index     int    `json:"index"`      // 行号
	Username  string `json:"username"`   // 用户名
	StudentID string `json:"student_id"` // 学号
	Status    string `json:"status"`     // 状态：success/failed
	Message   string `json:"message"`    // 消息
}

// Excel解析后的学生数据
type StudentExcelData struct {
	Index      int    // 行号
	Name       string // 姓名
	StudentID  string // 学号
	Email      string // 邮箱
	Phone      string // 手机号
	ClassName  string // 班级名称
	ClassID    uint   // 班级ID
	TeacherIDs []uint // 关联教师ID
}
