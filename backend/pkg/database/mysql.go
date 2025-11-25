package database

import (
	"webGL-720yun/app/models"
	"fmt"
	"log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

var DB *gorm.DB

// Init 初始化数据库连接
func Init(config *models.DatabaseConfig) error {
	var err error
	var steps []string
	var successSteps []string
	var failedSteps []string

	// 设置GORM日志级别为Error，只输出错误信息
	logConfig := logger.Config{
		SlowThreshold:             time.Second, // 慢SQL阈值
		LogLevel:                  logger.Error, // 日志级别
		IgnoreRecordNotFoundError: true,        // 忽略记录未找到错误
		Colorful:                  false,       // 禁用彩色输出
	}

	// 构建DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
		config.Charset,
	)

	// 步骤1: 连接数据库
	steps = append(steps, "连接数据库")
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
		Logger: logger.New(
			log.New(log.Writer(), "", log.LstdFlags), // io writer
			logConfig,
		),
	})

	if err != nil {
		failedSteps = append(failedSteps, fmt.Sprintf("连接数据库: 失败 - %v", err))
		printInitResult(steps, successSteps, failedSteps)
		return fmt.Errorf("数据库连接失败: %v", err)
	}
	successSteps = append(successSteps, "连接数据库: 成功")

	// 步骤2: 设置连接池
	steps = append(steps, "设置连接池")
	sqlDB, err := DB.DB()
	if err != nil {
		failedSteps = append(failedSteps, fmt.Sprintf("设置连接池: 失败 - %v", err))
		printInitResult(steps, successSteps, failedSteps)
		return fmt.Errorf("获取数据库实例失败: %v", err)
	}

	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	successSteps = append(successSteps, "设置连接池: 成功")

	// 步骤3: 自动迁移表结构
	steps = append(steps, "自动迁移表结构")
	if err := migrate(); err != nil {
		failedSteps = append(failedSteps, fmt.Sprintf("自动迁移表结构: 失败 - %v", err))
		printInitResult(steps, successSteps, failedSteps)
		return fmt.Errorf("数据库迁移失败: %v", err)
	}
	successSteps = append(successSteps, "自动迁移表结构: 成功")

	// 步骤4: 初始化基础数据
	steps = append(steps, "初始化基础数据")
	if err := initBasicData(); err != nil {
		failedSteps = append(failedSteps, fmt.Sprintf("初始化基础数据: 失败 - %v", err))
		printInitResult(steps, successSteps, failedSteps)
		return fmt.Errorf("初始化基础数据失败: %v", err)
	}
	successSteps = append(successSteps, "初始化基础数据: 成功")

	// 打印最终结果
	printInitResult(steps, successSteps, failedSteps)

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

// printInitResult 打印初始化结果
func printInitResult(steps []string, successSteps []string, failedSteps []string) {
	fmt.Println("\n=== 数据库初始化结果 ===")
	for _, step := range successSteps {
		fmt.Printf("✅ %s\n", step)
	}
	for _, step := range failedSteps {
		fmt.Printf("❌ %s\n", step)
	}
	fmt.Printf("\n总步骤: %d, 成功: %d, 失败: %d\n", len(steps), len(successSteps), len(failedSteps))
	fmt.Println("========================\n")
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