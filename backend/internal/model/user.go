package model

import (
	"time"
	"webGL-720yun/pkg/utils"

	"gorm.io/gorm"
)

const (
	UserStatusActive   = 1
	UserStatusInactive = 0
)

const (
	RoleAdmin   = "admin"
	RoleTeacher = "teacher"
	RoleStudent = "student"
)

type User struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Username    string         `gorm:"type:varchar(50);uniqueIndex;not null;comment:用户名" json:"username"`
	Password    string         `gorm:"type:varchar(255);not null;comment:密码" json:"-"`
	Email       string         `gorm:"type:varchar(100);uniqueIndex;comment:邮箱" json:"email"`
	Phone       string         `gorm:"type:varchar(20);comment:手机号" json:"phone"`
	Nickname    string         `gorm:"type:varchar(50);comment:昵称" json:"nickname"`
	Avatar      string         `gorm:"type:varchar(255);comment:头像" json:"avatar"`
	RoleID      uint           `gorm:"not null;index;comment:角色ID" json:"role_id"`
	Role        Role           `gorm:"foreignKey:RoleID" json:"role"`
	IsSuperAdmin bool          `gorm:"default:false;comment:是否超级管理员" json:"is_super_admin"`
	Status      int            `gorm:"default:1;index;comment:状态(1启用0禁用)" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (User) TableName() string {
	return "sys_users"
}

func (u *User) ResetPassword() (string, error) {
	tempPassword, err := utils.GenerateRandomPassword(8)
	if err != nil {
		return "", err
	}
	hashedPassword, err := utils.HashPassword(tempPassword)
	if err != nil {
		return "", err
	}
	u.Password = hashedPassword
	return tempPassword, nil
}

func (u *User) ChangePassword(oldPassword, newPassword string) error {
	if !utils.CheckPassword(oldPassword, u.Password) {
		return utils.ErrOldPasswordIncorrect
	}
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}
	u.Password = hashedPassword
	return nil
}

func (u *User) ToTeacher() *Teacher {
	return &Teacher{User: *u}
}

func (u *User) ToStudent() *Student {
	return &Student{User: *u}
}

type Teacher struct {
	User
	Classes []Class `gorm:"foreignKey:TeacherID" json:"classes"`
}

func (Teacher) TableName() string {
	return "sys_users"
}

type Student struct {
	User
	StudentID string    `gorm:"type:varchar(20);uniqueIndex;not null;comment:学号" json:"student_id"`
	ClassID   uint      `gorm:"not null;index;comment:班级ID" json:"class_id"`
	Class     Class     `gorm:"foreignKey:ClassID" json:"class"`
	Teachers  []Teacher `gorm:"many2many:sys_student_teachers;" json:"teachers"`
}

func (Student) TableName() string {
	return "sys_users"
}

type UserSession struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"not null;index;comment:用户ID" json:"user_id"`
	User         User      `gorm:"foreignKey:UserID" json:"user"`
	AccessToken  string    `gorm:"type:text;not null;comment:访问令牌" json:"-"`
	RefreshToken string    `gorm:"type:text;not null;comment:刷新令牌" json:"-"`
	IP           string    `gorm:"type:varchar(50);comment:IP地址" json:"ip"`
	UserAgent    string    `gorm:"type:varchar(500);comment:用户代理" json:"user_agent"`
	ExpiresAt    time.Time `gorm:"comment:过期时间" json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (UserSession) TableName() string {
	return "sys_user_sessions"
}
