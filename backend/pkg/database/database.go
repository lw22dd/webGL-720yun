package database

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"webGL-720yun/config"

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
	successSteps = append(steps, "设置连接池: 成功")

	steps = append(steps, "执行初始化SQL脚本")
	if err := executeInitSQL(db); err != nil {
		failedSteps = append(failedSteps, fmt.Sprintf("执行初始化SQL脚本: 失败 - %v", err))
		printInitResult(steps, successSteps, failedSteps)
		return nil, fmt.Errorf("执行初始化SQL脚本失败: %v", err)
	}
	successSteps = append(steps, "执行初始化SQL脚本: 成功")

	printInitResult(steps, successSteps, failedSteps)

	return &Database{DB: db}, nil
}

func executeInitSQL(db *gorm.DB) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取当前工作目录失败: %v", err)
	}

	var sqlFile string
	var found bool

	possiblePaths := []string{
		currentDir + "/pkg/database/init.sql",
		currentDir + "/../pkg/database/init.sql",
		currentDir + "/../../pkg/database/init.sql",
		currentDir + "/backend/pkg/database/init.sql",
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			sqlFile = path
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("init.sql文件不存在，尝试的路径: %v", possiblePaths)
	}

	if _, err := os.Stat(sqlFile); os.IsNotExist(err) {
		return fmt.Errorf("init.sql文件不存在: %v, 尝试的路径: %s", err, sqlFile)
	}

	fmt.Printf("正在执行SQL脚本: %s\n", sqlFile)

	content, err := os.ReadFile(sqlFile)
	if err != nil {
		return fmt.Errorf("读取SQL文件失败: %v, 路径: %s", err, sqlFile)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取数据库连接失败: %v", err)
	}

	sqlScript := string(content)
	lines := strings.Split(sqlScript, "\n")
	var currentStmt strings.Builder
	var statements []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "--") {
			continue
		}

		currentStmt.WriteString(line)
		currentStmt.WriteString("\n")

		if strings.HasSuffix(line, ";") {
			statements = append(statements, currentStmt.String())
			currentStmt.Reset()
		}
	}

	if currentStmt.Len() > 0 {
		statements = append(statements, currentStmt.String())
	}

	fmt.Printf("SQL语句数量: %d\n", len(statements))

	for i, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}

		shortStmt := stmt
		if len(shortStmt) > 50 {
			shortStmt = shortStmt[:50] + "..."
		}
		fmt.Printf("正在执行SQL语句 %d: %s\n", i+1, shortStmt)

		_, err = sqlDB.Exec(stmt)
		if err != nil {
			return fmt.Errorf("执行SQL语句失败: %v, 语句: %s", err, stmt)
		}
	}

	fmt.Printf("SQL脚本执行成功，共执行 %d 条语句\n", len(statements))

	return nil
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
	fmt.Println("========================\n")
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
