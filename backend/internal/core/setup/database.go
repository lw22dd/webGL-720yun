package setup

import (
	"fmt"
	"log"
	"time"

	"webGL-720yun/config"
	"webGL-720yun/internal/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type Database struct {
	DB *gorm.DB
}

func NewDatabase(cfg *config.DatabaseConfig) (*Database, error) {
	var steps []string
	var successSteps []string
	var failedSteps []string

	logConfig := logger.Config{
		SlowThreshold:             time.Second,
		LogLevel:                  logger.Error,
		IgnoreRecordNotFoundError: true,
		Colorful:                  false,
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.Charset,
	)

	steps = append(steps, "连接数据库")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: false,
		},
		Logger: logger.New(
			log.New(log.Writer(), "", log.LstdFlags),
			logConfig,
		),
	})

	if err != nil {
		failedSteps = append(failedSteps, fmt.Sprintf("连接数据库: 失败 - %v", err))
		printInitResult(steps, successSteps, failedSteps)
		return nil, fmt.Errorf("数据库连接失败: %v", err)
	}
	successSteps = append(successSteps, "连接数据库: 成功")

	steps = append(steps, "设置连接池")
	sqlDB, err := db.DB()
	if err != nil {
		failedSteps = append(failedSteps, fmt.Sprintf("设置连接池: 失败 - %v", err))
		printInitResult(steps, successSteps, failedSteps)
		return nil, fmt.Errorf("获取数据库实例失败: %v", err)
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	successSteps = append(successSteps, "设置连接池: 成功")

	steps = append(steps, "自动迁移数据库表结构")
	if err := autoMigrate(db); err != nil {
		failedSteps = append(failedSteps, fmt.Sprintf("自动迁移数据库表结构: 失败 - %v", err))
		printInitResult(steps, successSteps, failedSteps)
		return nil, fmt.Errorf("自动迁移数据库表结构失败: %v", err)
	}
	successSteps = append(successSteps, "自动迁移数据库表结构: 成功")

	steps = append(steps, "插入测试数据")
	if err := SeedTestData(db); err != nil {
		failedSteps = append(failedSteps, fmt.Sprintf("插入测试数据: 失败 - %v", err))
		printInitResult(steps, successSteps, failedSteps)
		return nil, fmt.Errorf("插入测试数据失败: %v", err)
	}
	successSteps = append(successSteps, "插入测试数据: 成功")

	printInitResult(steps, successSteps, failedSteps)

	return &Database{DB: db}, nil
}

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.Permission{},
		&model.Class{},
		&model.Student{},
		&model.Teacher{},
		&model.StudentTeacher{},
		&model.UserSession{},
		&model.ResSpace{},
		&model.ResScene{},
		&model.ResHotspot{},
	)
}

func printInitResult(steps []string, successSteps []string, failedSteps []string) {
	fmt.Println("\n=== 数据库初始化结果 ===")
	for _, step := range successSteps {
		fmt.Printf("✅ %s\n", step)
	}
	for _, step := range failedSteps {
		fmt.Printf("❌ %s\n", step)
	}
	fmt.Printf("\n总步骤: %d, 成功: %d, 失败: %d\n", len(steps), len(successSteps), len(failedSteps))
<<<<<<< HEAD
	fmt.Println("========================\n")
=======
	fmt.Println("========================")
>>>>>>> 8146554307dc850256079e5aa35fb05bb5b6a503
}

func (d *Database) GetDB() *gorm.DB {
	return d.DB
}

func (d *Database) Close() error {
	if d.DB != nil {
		sqlDB, err := d.DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
