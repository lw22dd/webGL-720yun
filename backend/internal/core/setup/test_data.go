package setup

import (
	"webGL-720yun/internal/model"
	"webGL-720yun/pkg/utils"

	"gorm.io/gorm"
)

func SeedTestData(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := seedRoles(tx); err != nil {
			return err
		}

		if err := seedPermissions(tx); err != nil {
			return err
		}

		if err := seedAdminUser(tx); err != nil {
			return err
		}

		if err := seedTestClasses(tx); err != nil {
			return err
		}

		return nil
	})
}

func seedRoles(db *gorm.DB) error {
	roles := []model.Role{
		{Name: model.RoleAdmin, Description: "系统管理员，拥有所有权限"},
		{Name: model.RoleTeacher, Description: "教师，可以管理学生和查看报告"},
		{Name: model.RoleStudent, Description: "学生，可以查看自己的数据"},
	}

	for _, role := range roles {
		if err := db.FirstOrCreate(&role, model.Role{Name: role.Name}).Error; err != nil {
			return err
		}
	}

	return nil
}

func seedPermissions(db *gorm.DB) error {
	var adminRole, teacherRole, studentRole model.Role

	if err := db.Where("name = ?", model.RoleAdmin).First(&adminRole).Error; err != nil {
		return err
	}
	if err := db.Where("name = ?", model.RoleTeacher).First(&teacherRole).Error; err != nil {
		return err
	}
	if err := db.Where("name = ?", model.RoleStudent).First(&studentRole).Error; err != nil {
		return err
	}

	permissions := []model.Permission{
		{Name: "user:read", Description: "查看用户", Resource: "/api/v1/user", Action: "GET"},
		{Name: "user:write", Description: "编辑用户", Resource: "/api/v1/user", Action: "POST"},
		{Name: "user:delete", Description: "删除用户", Resource: "/api/v1/user", Action: "DELETE"},
		{Name: "class:read", Description: "查看班级", Resource: "/api/v1/class", Action: "GET"},
		{Name: "class:write", Description: "编辑班级", Resource: "/api/v1/class", Action: "POST"},
	}

	for i := range permissions {
		if err := db.FirstOrCreate(&permissions[i], model.Permission{Name: permissions[i].Name}).Error; err != nil {
			return err
		}
	}

	if err := db.Model(&adminRole).Association("Permissions").Replace(permissions); err != nil {
		return err
	}

	teacherPerms := []model.Permission{permissions[0], permissions[3]}
	if err := db.Model(&teacherRole).Association("Permissions").Replace(teacherPerms); err != nil {
		return err
	}

	studentPerms := []model.Permission{permissions[0]}
	if err := db.Model(&studentRole).Association("Permissions").Replace(studentPerms); err != nil {
		return err
	}

	return nil
}

func seedAdminUser(db *gorm.DB) error {
	var adminRole model.Role
	if err := db.Where("name = ?", model.RoleAdmin).First(&adminRole).Error; err != nil {
		return err
	}

	var count int64
	db.Model(&model.User{}).Where("role_id = ?", adminRole.ID).Count(&count)
	if count > 0 {
		return nil
	}

	hashedPassword, err := utils.HashPassword("admin123")
	if err != nil {
		return err
	}

	admin := model.User{
		Username:    "admin",
		Password:    hashedPassword,
		Email:       "admin@example.com",
		Nickname:    "系统管理员",
		RoleID:      adminRole.ID,
		Status:      model.UserStatusActive,
		IsSuperAdmin: true,
	}

	return db.Create(&admin).Error
}

func seedTestClasses(db *gorm.DB) error {
	var teacherRole model.Role
	if err := db.Where("name = ?", model.RoleTeacher).First(&teacherRole).Error; err != nil {
		return err
	}

	var count int64
	db.Model(&model.User{}).Where("role_id = ?", teacherRole.ID).Count(&count)
	if count > 0 {
		return nil
	}

	hashedPassword, err := utils.HashPassword("teacher123")
	if err != nil {
		return err
	}

	teachers := []model.User{
		{
			Username: "teacher1",
			Password: hashedPassword,
			Email:    "teacher1@example.com",
			Nickname: "张老师",
			RoleID:   teacherRole.ID,
			Status:   model.UserStatusActive,
		},
		{
			Username: "teacher2",
			Password: hashedPassword,
			Email:    "teacher2@example.com",
			Nickname: "李老师",
			RoleID:   teacherRole.ID,
			Status:   model.UserStatusActive,
		},
	}

	for _, teacher := range teachers {
		if err := db.Create(&teacher).Error; err != nil {
			return err
		}
	}

	classes := []model.Class{
		{Name: "计算机科学1班", Description: "2024级计算机科学1班", TeacherID: teachers[0].ID},
		{Name: "计算机科学2班", Description: "2024级计算机科学2班", TeacherID: teachers[1].ID},
		{Name: "软件工程1班", Description: "2024级软件工程1班", TeacherID: teachers[0].ID},
	}

	for _, class := range classes {
		if err := db.Create(&class).Error; err != nil {
			return err
		}
	}

	return nil
}
