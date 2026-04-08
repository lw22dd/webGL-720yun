package model

import "time"

type Role struct {
	ID          uint         `gorm:"primaryKey" json:"id"`
	Name        string       `gorm:"type:varchar(50);uniqueIndex;not null;comment:角色名称" json:"name"`
	Description string       `gorm:"type:varchar(255);comment:角色描述" json:"description"`
	Permissions []Permission `gorm:"many2many:sys_role_permissions;" json:"permissions"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

func (Role) TableName() string {
	return "sys_roles"
}

type Permission struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(100);uniqueIndex;not null;comment:权限名称" json:"name"`
	Description string    `gorm:"type:varchar(255);comment:权限描述" json:"description"`
	Resource    string    `gorm:"type:varchar(50);comment:资源" json:"resource"`
	Action      string    `gorm:"type:varchar(50);comment:动作" json:"action"`
	Roles       []Role    `gorm:"many2many:sys_role_permissions;" json:"-"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Permission) TableName() string {
	return "sys_permissions"
}
