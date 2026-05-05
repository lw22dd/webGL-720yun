package database

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

// DB 是数据库连接包装器
type DB struct {
	*gorm.DB
}

// NewMySQL 创建新的 MySQL 数据库连接
func NewMySQL(cfg *config.DatabaseConfig) (*DB, error) {
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
		return nil, fmt.Errorf("数据库连接失败: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取数据库实例失败: %w", err)
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)

	return &DB{DB: db}, nil
}

// AutoMigrate 自动迁移数据库表结构
func (d *DB) AutoMigrate() error {
	return d.DB.AutoMigrate(
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

// Close 关闭数据库连接
func (d *DB) Close() error {
	if d.DB != nil {
		sqlDB, err := d.DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
