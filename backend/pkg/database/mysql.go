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

var DB *gorm.DB

// Init 初始化数据库连接
func Init(config *config.DatabaseConfig) error {
	var err error
	var steps []string
	var successSteps []string
	var failedSteps []string

	// 设置GORM日志级别为Error，只输出错误信息
	logConfig := logger.Config{
		SlowThreshold:             time.Second,  // 慢SQL阈值
		LogLevel:                  logger.Error, // 日志级别
		IgnoreRecordNotFoundError: true,         // 忽略记录未找到错误
		Colorful:                  false,        // 禁用彩色输出
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

	// 步骤3: 执行初始化SQL脚本
	steps = append(steps, "执行初始化SQL脚本")
	if err := executeInitSQL(); err != nil {
		failedSteps = append(failedSteps, fmt.Sprintf("执行初始化SQL脚本: 失败 - %v", err))
		printInitResult(steps, successSteps, failedSteps)
		return fmt.Errorf("执行初始化SQL脚本失败: %v", err)
	}
	successSteps = append(successSteps, "执行初始化SQL脚本: 成功")

	// 打印最终结果
	printInitResult(steps, successSteps, failedSteps)

	return nil
}

// executeInitSQL 执行初始化SQL脚本
func executeInitSQL() error {
	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取当前工作目录失败: %v", err)
	}

	// 根据当前工作目录计算init.sql的正确路径
	var sqlFile string

	// 检查当前工作目录是否包含app/user（测试环境）
	if strings.Contains(currentDir, "app/user") {
		// 测试环境：从app/user向上两级目录
		sqlFile = currentDir + "/../../pkg/database/init.sql"
	} else if strings.HasSuffix(currentDir, "backend") {
		// 正常运行环境：直接使用相对路径
		sqlFile = currentDir + "/pkg/database/init.sql"
	} else {
		// 其他情况：假设当前目录是项目根目录
		sqlFile = currentDir + "/backend/pkg/database/init.sql"
	}

	// 检查文件是否存在
	if _, err := os.Stat(sqlFile); os.IsNotExist(err) {
		return fmt.Errorf("init.sql文件不存在: %v, 尝试的路径: %s", err, sqlFile)
	}

	// 读取SQL文件内容
	content, err := os.ReadFile(sqlFile)
	if err != nil {
		return fmt.Errorf("读取SQL文件失败: %v, 路径: %s", err, sqlFile)
	}

	// 使用database/sql包的Exec方法执行SQL文件
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("获取数据库连接失败: %v", err)
	}

	// 将SQL内容按分号分割成多个语句
	sqlStatements := splitSQL(string(content))

	// 逐个执行SQL语句
	for _, stmt := range sqlStatements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" || strings.HasPrefix(strings.TrimSpace(stmt), "--") {
			continue // 跳过空语句和注释
		}

		_, err = sqlDB.Exec(stmt)
		if err != nil {
			return fmt.Errorf("执行SQL语句失败: %v, 语句: %s", err, stmt)
		}
	}

	return nil
}

// splitSQL 将SQL内容按分号分割成多个语句
func splitSQL(sql string) []string {
	var statements []string
	var currentStmt strings.Builder
	inString := false
	stringChar := byte(0)

	for i := 0; i < len(sql); i++ {
		char := sql[i]

		// 处理字符串开始/结束
		if (char == '\'' || char == '"' || char == '`') && (i == 0 || sql[i-1] != '\\') {
			if !inString {
				inString = true
				stringChar = char
			} else if char == stringChar {
				inString = false
				stringChar = 0
			}
		}

		// 处理分号分隔符
		if char == ';' && !inString {
			statements = append(statements, currentStmt.String())
			currentStmt.Reset()
		} else {
			currentStmt.WriteByte(char)
		}
	}

	// 添加最后一个语句
	if currentStmt.Len() > 0 {
		statements = append(statements, currentStmt.String())
	}

	return statements
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

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}
