package database

import (
	"webGL-720yun/app/models"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

// Init 初始化数据库连接
func Init(config *models.DatabaseConfig) error {
	var err error
	
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
		config.Charset,
	)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return fmt.Errorf("数据库连接失败: %v", err)
	}

	// 设置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("获取数据库实例失败: %v", err)
	}

	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)

	// 自动迁移表结构
	if err := migrate(); err != nil {
		return fmt.Errorf("数据库迁移失败: %v", err)
	}

	// 初始化基础数据
	if err := initBasicData(); err != nil {
		return fmt.Errorf("初始化基础数据失败: %v", err)
	}

	return nil
}

// migrate 自动迁移数据库表
func migrate() error {
	return DB.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Permission{},
		&models.UserSession{},
	)
}

// initBasicData 初始化基础角色和权限数据
func initBasicData() error {
	// 检查是否已经存在数据
	var count int64
	DB.Model(&models.Role{}).Count(&count)
	if count > 0 {
		return nil // 已有数据，跳过初始化
	}

	// 创建基础权限
	permissions := []models.Permission{
		{Name: "user:create", Description: "创建用户", Resource: "user", Action: "create"},
		{Name: "user:read", Description: "查看用户", Resource: "user", Action: "read"},
		{Name: "user:update", Description: "更新用户", Resource: "user", Action: "update"},
		{Name: "user:delete", Description: "删除用户", Resource: "user", Action: "delete"},
		{Name: "profile:read", Description: "查看个人信息", Resource: "profile", Action: "read"},
		{Name: "profile:update", Description: "更新个人信息", Resource: "profile", Action: "update"},
	}

	if err := DB.Create(&permissions).Error; err != nil {
		return err
	}

	// 创建管理员角色
	adminRole := models.Role{
		Name:        models.RoleAdmin,
		Description: "管理员",
		Permissions: permissions, // 拥有所有权限
	}

	// 创建学生角色
	studentRole := models.Role{
		Name:        models.RoleStudent,
		Description: "学生",
		Permissions: permissions[4:], // 只拥有个人信息相关权限
	}

	roles := []models.Role{adminRole, studentRole}
	if err := DB.Create(&roles).Error; err != nil {
		return err
	}

	return nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}